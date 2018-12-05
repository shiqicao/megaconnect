package flowmanager

import (
	"context"
	"io"
	"net"
	"testing"
	"time"

	mgrpc "github.com/megaspacelab/megaconnect/grpc"
	"github.com/megaspacelab/megaconnect/protos"
	w "github.com/megaspacelab/megaconnect/workflow"

	"github.com/golang/protobuf/proto"
	"github.com/golang/protobuf/ptypes/empty"
	"github.com/stretchr/testify/suite"
	"go.uber.org/zap"
	"google.golang.org/grpc"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

const chainID = "Example"

type FlowManagerSuite struct {
	suite.Suite

	log        *zap.Logger
	stateStore *MemStateStore
	fm         *FlowManager

	server    *grpc.Server
	conn      *grpc.ClientConn
	mblockAPI mgrpc.MBlockApiClient
	ctx       context.Context
	ctxCancel context.CancelFunc
}

func (s *FlowManagerSuite) SetupTest() {
	log, err := zap.NewDevelopment()
	s.Require().NoError(err)
	s.log = log
	s.stateStore = NewMemStateStore()
	s.fm = NewFlowManager(s.stateStore, log)

	s.server = grpc.NewServer()
	s.fm.Register(s.server)

	lis, err := net.Listen("tcp", "localhost:")
	s.Require().NoError(err)

	s.conn, err = grpc.Dial(lis.Addr().String(), grpc.WithInsecure())
	s.Require().NoError(err)
	s.mblockAPI = mgrpc.NewMBlockApiClient(s.conn)

	s.ctx, s.ctxCancel = context.WithCancel(context.Background())

	go s.server.Serve(lis)
}

func (s *FlowManagerSuite) TearDownTest() {
	if s.conn != nil {
		s.conn.Close()
	}
	if s.server != nil {
		s.server.Stop()
	}
}

func (s *FlowManagerSuite) workflow1() *w.WorkflowDecl {
	wf := w.NewWorkflowDecl(w.NewId("TestWorkflow"), 1)
	wf.AddDecl(w.NewEventDecl(
		w.NewId("TestEvent"),
		w.NewObjType(w.NewIdToTy().Put(
			"balance",
			w.IntType,
		)),
	))
	wf.AddDecl(w.NewEventDecl(
		w.NewId("TestEvent2"),
		w.NewObjType(w.NewIdToTy().Put(
			"balance",
			w.IntType,
		)),
	))
	wf.AddDecl(w.NewEventDecl(
		w.NewId("TestEvent3"),
		w.NewObjType(w.NewIdToTy().Put(
			"balance",
			w.IntType,
		)),
	))
	wf.AddDecl(w.NewMonitorDecl(
		w.NewId("TestMonitor"),
		w.TrueConst,
		w.NewIdToExpr().Put(
			"balance",
			w.NewIntConstFromI64(100),
		),
		w.NewFire(
			"TestEvent",
			w.NewObjLit(w.NewIdToExpr().Put(
				"balance",
				w.NewIntConstFromI64(1),
			)),
		),
		chainID,
	))
	wf.AddDecl(w.NewActionDecl(
		w.NewId("TestAction"),
		w.NewEVar("TestEvent"),
		[]w.Stmt{
			w.NewFire(
				"TestEvent2",
				w.NewProps(w.NewVar("TestEvent")),
			),
		},
	))
	wf.AddDecl(w.NewActionDecl(
		w.NewId("TestAction2"),
		w.NewEBinOp(
			w.AndEOp,
			w.NewEVar("TestEvent"),
			w.NewEVar("TestEvent2"),
		),
		[]w.Stmt{
			w.NewFire(
				"TestEvent3",
				w.NewProps(w.NewVar("TestEvent2")),
			),
		},
	))
	return wf
}

func (s *FlowManagerSuite) TestDeployAndUndeployWorkflow() {
	ctx := context.Background()
	config, sub, err := s.fm.GetChainConfig(chainID)
	s.Require().NoError(err)
	s.Require().NotNil(sub)
	s.Require().NotNil(config)
	s.Require().Empty(config.Monitors)

	// deploy a valid workflow should be successful
	wf := s.workflow1()
	bin, err := w.EncodeWorkflow(wf)
	s.Require().NoError(err)
	id, err := s.fm.DeployWorkflow(ctx, &mgrpc.DeployWorkflowRequest{
		Payload: bin,
	})
	s.Require().Equal(id.WorkflowId, []byte(wf.Name().String()))
	s.Require().NoError(err)
	s.Require().Empty(sub.Patches(), "Workflow activated too early")

	// deploy a workflow twice should fail gracefully
	_, err = s.fm.DeployWorkflow(ctx, &mgrpc.DeployWorkflowRequest{
		Payload: bin,
	})
	s.Require().Error(err)

	// deploy a workflow without payload should fail gracefully
	_, err = s.fm.DeployWorkflow(ctx, &mgrpc.DeployWorkflowRequest{
		Payload: nil,
	})
	s.Require().Error(err)

	mblock, err := s.fm.FinalizeAndCommitMBlock()
	s.Require().NoError(err)
	s.Require().NotNil(mblock)

	s.Require().Len(mblock.Events, 1)

	body := mblock.Events[0].Body
	s.Require().IsType(body, new(protos.MEvent_DeployWorkflow))

	actualwf, err := w.
		NewByteDecoder(body.(*protos.MEvent_DeployWorkflow).DeployWorkflow.Payload).
		DecodeWorkflow()
	s.Require().NoError(err)
	s.Require().True(actualwf.Equal(wf))

	select {
	case patch := <-sub.Patches():
		s.Require().NotNil(patch)
		s.Require().Len(patch.AddMonitors, 1)
		s.Require().Empty(patch.RemoveMonitors)
	default:
		s.FailNow("Config wasn't updated after workflow deployment")
	}

	// undeploy a deployed workflow should be successful
	_, err = s.fm.UndeployWorkflow(ctx, &mgrpc.UndeployWorkflowRequest{
		WorkflowId: id.WorkflowId,
	})
	s.Require().NoError(err)
	s.Require().Empty(sub.Patches(), "Workflow deactivated too early")

	// undeploy a non-existing workflow should not fail gracefully
	_, err = s.fm.UndeployWorkflow(ctx, &mgrpc.UndeployWorkflowRequest{
		WorkflowId: []byte("invalid"),
	})
	s.Require().Error(err)

	// undeploy a workflow without workflow id should fail gracefully
	_, err = s.fm.UndeployWorkflow(ctx, &mgrpc.UndeployWorkflowRequest{
		WorkflowId: nil,
	})
	s.Require().Error(err)

	mblock, err = s.fm.FinalizeAndCommitMBlock()
	s.Require().NoError(err)
	s.Require().NotNil(mblock)
	s.Require().Len(mblock.Events, 1)

	body = mblock.Events[0].Body
	s.Require().IsType(body, new(protos.MEvent_UndeployWorkflow))
	s.Require().Equal(GetWorkflowID(wf).UnsafeBytes(), body.(*protos.MEvent_UndeployWorkflow).UndeployWorkflow.WorkflowId)
}

func (s *FlowManagerSuite) TestReportBlockEvents() {
	wf := s.workflow1()
	_, err := s.fm.deployWorkflow(wf)
	s.Require().NoError(err)

	_, err = s.fm.FinalizeAndCommitMBlock()
	s.Require().NoError(err)

	config, sub, err := s.fm.GetChainConfig(chainID)
	s.Require().NoError(err)
	s.Require().NotNil(sub)
	s.Require().NotNil(config)
	s.Require().Len(config.Monitors, 1)

	block := &protos.Block{Chain: chainID}

	// This should be ignored
	s.fm.ReportBlockEvents(config.MonitorsVersion-1, block, nil)

	actualBlock, err := s.stateStore.LastReportedBlockByChain(chainID)
	s.Require().NoError(err)
	s.Require().Nil(actualBlock)

	eventPayload, err := w.EncodeObjConst(
		w.NewObjConst(w.ObjFields{
			"balance": w.NewIntConstFromI64(1),
		}),
	)
	s.Require().NoError(err)

	// This should be processed, and it should trigger a chain of events
	monitorEvent := &protos.MonitorEvent{
		WorkflowId:  GetWorkflowID(wf).UnsafeBytes(),
		MonitorName: wf.MonitorDecls()[0].Name().Id(),
		EventName:   "TestEvent",
		Payload:     eventPayload,
	}
	s.fm.ReportBlockEvents(config.MonitorsVersion, block, []*protos.MonitorEvent{monitorEvent})

	actualBlock, err = s.stateStore.LastReportedBlockByChain(chainID)
	s.Require().NoError(err)
	s.Require().True(proto.Equal(block, actualBlock))

	// Finalize mega block and verify its contents
	mblock, err := s.fm.FinalizeAndCommitMBlock()
	s.Require().NoError(err)
	s.Require().NotNil(mblock)
	s.Require().Len(mblock.Events, 4)

	body := mblock.Events[0].Body
	s.Require().IsType(body, new(protos.MEvent_ReportBlock))
	s.Require().True(proto.Equal(body.(*protos.MEvent_ReportBlock).ReportBlock, block))

	body = mblock.Events[1].Body
	s.Require().IsType(body, new(protos.MEvent_MonitorEvent))
	s.Require().True(proto.Equal(body.(*protos.MEvent_MonitorEvent).MonitorEvent, monitorEvent))

	actionEvent := &protos.ActionEvent{
		WorkflowId: GetWorkflowID(wf).UnsafeBytes(),
		ActionName: "TestAction",
		EventName:  "TestEvent2",
		Payload:    eventPayload,
	}
	body = mblock.Events[2].Body
	s.Require().IsType(body, new(protos.MEvent_ActionEvent))
	s.Require().True(proto.Equal(body.(*protos.MEvent_ActionEvent).ActionEvent, actionEvent))

	actionEvent.ActionName = "TestAction2"
	actionEvent.EventName = "TestEvent3"
	body = mblock.Events[3].Body
	s.Require().IsType(body, new(protos.MEvent_ActionEvent))
	s.Require().True(proto.Equal(body.(*protos.MEvent_ActionEvent).ActionEvent, actionEvent))
}

func (s *FlowManagerSuite) TestLatestMBlock() {
	// No mblocks yet.
	_, err := s.mblockAPI.LatestMBlock(s.ctx, new(empty.Empty))
	s.Require().Error(err)
	s.Require().Equal(codes.NotFound, status.Convert(err).Code())

	// Finalize an mblock and get again.
	mblock, err := s.fm.FinalizeAndCommitMBlock()
	s.Require().NoError(err)
	mblock2, err := s.mblockAPI.LatestMBlock(s.ctx, new(empty.Empty))
	s.Require().NoError(err)
	s.Require().True(proto.Equal(mblock, mblock2))
}

func (s *FlowManagerSuite) TestMBlockByHeight() {
	// Finalize a couple of mblocks.
	mblock, err := s.fm.FinalizeAndCommitMBlock()
	s.Require().NoError(err)
	mblock2, err := s.fm.FinalizeAndCommitMBlock()
	s.Require().NoError(err)

	// Retrieve those mblocks.
	actualMBlock, err := s.mblockAPI.MBlockByHeight(s.ctx, &mgrpc.MBlockByHeightRequest{Height: mblock.Height})
	s.Require().NoError(err)
	s.Require().True(proto.Equal(mblock, actualMBlock))

	actualMBlock, err = s.mblockAPI.MBlockByHeight(s.ctx, &mgrpc.MBlockByHeightRequest{Height: mblock2.Height})
	s.Require().NoError(err)
	s.Require().True(proto.Equal(mblock2, actualMBlock))

	// Out of range.
	_, err = s.mblockAPI.MBlockByHeight(s.ctx, &mgrpc.MBlockByHeightRequest{Height: mblock2.Height + 1})
	s.Require().Error(err)
}

func (s *FlowManagerSuite) waitForSubsLen(n int, timeoutSecs int) {
	checkSub := func() bool {
		s.fm.lock.Lock()
		defer s.fm.lock.Unlock()
		return len(s.fm.mblockSubs) == n
	}

	for i := 0; !checkSub(); i++ {
		if i >= timeoutSecs {
			s.FailNow("Timeout waiting for subscription")
		}
		time.Sleep(time.Second)
	}
}

func (s *FlowManagerSuite) TestSubscribeMBlock() {
	stream, err := s.mblockAPI.SubscribeMBlock(s.ctx, new(empty.Empty))
	s.Require().NoError(err)

	mblockChan := make(chan *protos.MBlock)
	errChan := make(chan error, 1)

	go func() {
		defer close(mblockChan)
		defer close(errChan)

		for {
			mblock, err := stream.Recv()
			if err != nil {
				if err != io.EOF {
					errChan <- err
				}
				return
			}
			mblockChan <- mblock
		}
	}()

	// Wait until subscription has been established on the server side.
	s.waitForSubsLen(1, 2)

	// Finalize a few mblocks.
	mblocks := make([]*protos.MBlock, 3)
	for i := 0; i < len(mblocks); i++ {
		mblocks[i], err = s.fm.FinalizeAndCommitMBlock()
		s.Require().NoError(err)
	}

	// Verify all mblocks were received.
	for i := 0; i < len(mblocks); i++ {
		select {
		case mblock, ok := <-mblockChan:
			s.Require().True(ok)
			s.Require().True(proto.Equal(mblock, mblocks[i]))
		case err := <-errChan:
			s.FailNowf("Subscription errored", "%v", err)
		case <-time.After(time.Second):
			s.FailNow("Timeout waiting for MBlock")
		}
	}
}

func (s *FlowManagerSuite) TestSubscribeMBlock_ClientCancel() {
	_, err := s.mblockAPI.SubscribeMBlock(s.ctx, new(empty.Empty))
	s.Require().NoError(err)

	// Wait until subscription has been established on the server side.
	s.waitForSubsLen(1, 2)

	// Cancel client side streaming.
	s.ctxCancel()

	// Finalize a few mblocks.
	for i := 0; i < 3; i++ {
		_, err = s.fm.FinalizeAndCommitMBlock()
		s.Require().NoError(err)
	}

	// The dead subscription should eventually be removed from fm.
	s.waitForSubsLen(0, 2)
}

func (s *FlowManagerSuite) TestSubscribeMBlock_ClientStall() {
	_, err := s.mblockAPI.SubscribeMBlock(s.ctx, new(empty.Empty))
	s.Require().NoError(err)

	// Wait until subscription has been established on the server side.
	s.waitForSubsLen(1, 2)

	// Finalize enough mblocks to overflow the buffer.
	for i := 0; i < mblockSubBufferSize+5; i++ {
		_, err = s.fm.FinalizeAndCommitMBlock()
		s.Require().NoError(err)
	}

	// The stalled subscription should eventually be removed from fm.
	s.waitForSubsLen(0, 2)
}

func TestFlowManager(t *testing.T) {
	suite.Run(t, new(FlowManagerSuite))
}

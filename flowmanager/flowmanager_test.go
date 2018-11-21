package flowmanager

import (
	"context"
	"testing"

	"github.com/golang/protobuf/proto"
	mgrpc "github.com/megaspacelab/megaconnect/grpc"
	"github.com/megaspacelab/megaconnect/protos"
	w "github.com/megaspacelab/megaconnect/workflow"
	"github.com/stretchr/testify/suite"
	"go.uber.org/zap"
)

const chainID = "Example"

type FlowManagerSuite struct {
	suite.Suite

	log        *zap.Logger
	stateStore *MemStateStore
	fm         *FlowManager
}

func (s *FlowManagerSuite) SetupTest() {
	log, err := zap.NewDevelopment()
	s.Require().NoError(err)
	s.log = log
	s.stateStore = NewMemStateStore()
	s.fm = NewFlowManager(s.stateStore, log)
}

func (s *FlowManagerSuite) workflow1() *w.WorkflowDecl {
	wf := w.NewWorkflowDecl(w.NewId("TestWorkflow"), 1)
	wf.AddChild(w.NewEventDecl(
		w.NewId("TestEvent"),
		w.NewObjType(w.NewIdToTy().Put(
			"balance",
			w.IntType,
		)),
	))
	wf.AddChild(w.NewEventDecl(
		w.NewId("TestEvent2"),
		w.NewObjType(w.NewIdToTy().Put(
			"balance",
			w.IntType,
		)),
	))
	wf.AddChild(w.NewEventDecl(
		w.NewId("TestEvent3"),
		w.NewObjType(w.NewIdToTy().Put(
			"balance",
			w.IntType,
		)),
	))
	wf.AddChild(w.NewMonitorDecl(
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
	wf.AddChild(w.NewActionDecl(
		w.NewId("TestAction"),
		w.NewEVar("TestEvent"),
		[]w.Stmt{
			w.NewFire(
				"TestEvent2",
				w.NewProps(w.NewVar("TestEvent")),
			),
		},
	))
	wf.AddChild(w.NewActionDecl(
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
	err := s.fm.deployWorkflow(wf)
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

func TestFlowManager(t *testing.T) {
	suite.Run(t, new(FlowManagerSuite))
}

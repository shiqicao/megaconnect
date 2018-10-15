package flowmanager

//go:generate moq -out chainmanagerserver_moq_test.go -pkg flowmanager ../grpc ChainManagerServer

import (
	"bytes"
	"context"
	"io"
	"net"
	"testing"
	"time"

	"github.com/megaspacelab/megaconnect/common"
	mgrpc "github.com/megaspacelab/megaconnect/grpc"

	"github.com/golang/protobuf/proto"
	"github.com/golang/protobuf/ptypes/empty"
	"github.com/stretchr/testify/suite"
	"go.uber.org/zap"
	"google.golang.org/grpc"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

type OrchestratorSuite struct {
	suite.Suite

	log *zap.Logger
	ctx context.Context

	server *grpc.Server
	conn   *grpc.ClientConn

	fm         *FlowManager
	orch       *Orchestrator
	orchClient mgrpc.OrchestratorClient
	cm         *fakeChainManager
}

func (s *OrchestratorSuite) SetupTest() {
	log, err := zap.NewDevelopment()
	s.Require().NoError(err)
	s.log = log
	s.ctx = context.Background()

	lis, err := net.Listen("tcp", "localhost:")
	s.Require().NoError(err)

	s.server = grpc.NewServer()
	s.conn, err = grpc.Dial(lis.Addr().String(), grpc.WithInsecure())
	s.Require().NoError(err)

	s.fm = NewFlowManager(s.log)
	s.orch = NewOrchestrator(s.fm, s.log)
	s.orchClient = mgrpc.NewOrchestratorClient(s.conn)
	s.orch.Register(s.server)

	s.cm = &fakeChainManager{
		chainID:    "test",
		cmID:       &mgrpc.InstanceId{Id: []byte("cm1"), Instance: 10},
		listenPort: lis.Addr().(*net.TCPAddr).Port,
		sessionID:  []byte("session1"),

		ChainManagerServerMock: ChainManagerServerMock{
			SetMonitorsFunc: func(stream mgrpc.ChainManager_SetMonitorsServer) error {
				return stream.SendAndClose(new(empty.Empty))
			},
			UpdateMonitorsFunc: func(stream mgrpc.ChainManager_UpdateMonitorsServer) error {
				return stream.SendAndClose(new(empty.Empty))
			},
		},
	}
	mgrpc.RegisterChainManagerServer(s.server, s.cm)

	go s.server.Serve(lis)
}

func (s *OrchestratorSuite) TearDownTest() {
	if s.conn != nil {
		s.conn.Close()
	}
	if s.server != nil {
		s.server.Stop()
	}
}

func (s *OrchestratorSuite) registerCM() (*mgrpc.RegisterChainManagerResponse, error) {
	return s.orchClient.RegisterChainManager(s.ctx, &mgrpc.RegisterChainManagerRequest{
		ChainId:        s.cm.chainID,
		ChainManagerId: s.cm.cmID,
		ListenPort:     int32(s.cm.listenPort),
		SessionId:      s.cm.sessionID,
	})
}

func (s *OrchestratorSuite) TestRegisterChainManager() {
	ims := IndexedMonitors{1: &mgrpc.Monitor{Id: 1}}
	s.fm.SetChainConfig(s.cm.chainID, ims, nil)

	reg, err := s.registerCM()
	s.Require().NoError(err)

	s.True(reg.Lease.RemainingSeconds > 0)
	s.Equal(1, len(reg.Monitors.Monitors))
	s.True(proto.Equal(ims[1], reg.Monitors.Monitors[0]))
}

func (s *OrchestratorSuite) TestRenewLease() {
	s.fm.SetChainConfig(s.cm.chainID, nil, nil)

	reg, err := s.registerCM()
	s.Require().NoError(err)

	lease2, err := s.orchClient.RenewLease(s.ctx, &mgrpc.RenewLeaseRequest{
		LeaseId: reg.Lease.Id,
	})
	s.Require().NoError(err)
	s.True(lease2.RemainingSeconds > 0)
	s.False(bytes.Equal(reg.Lease.Id, lease2.Id))
}

func (s *OrchestratorSuite) TestReportBlockEvents() {
	ims := IndexedMonitors{1: &mgrpc.Monitor{Id: 1}}
	s.fm.SetChainConfig(s.cm.chainID, ims, nil)

	reg, err := s.registerCM()
	s.Require().NoError(err)

	stream, err := s.orchClient.ReportBlockEvents(s.ctx)
	s.Require().NoError(err)

	err = stream.Send(&mgrpc.ReportBlockEventsRequest{
		MsgType: &mgrpc.ReportBlockEventsRequest_Preflight_{
			Preflight: &mgrpc.ReportBlockEventsRequest_Preflight{
				LeaseId:           reg.Lease.Id,
				MonitorSetVersion: reg.Monitors.Version,
			},
		},
	})
	s.Require().NoError(err)

	err = stream.Send(&mgrpc.ReportBlockEventsRequest{
		MsgType: &mgrpc.ReportBlockEventsRequest_Block{
			Block: &mgrpc.Block{
				Hash: bytes.Repeat([]byte{1}, common.HashSize),
				Height: &mgrpc.BigInt{
					Bytes:    []byte{10},
					Negative: false,
				},
				ParentHash: bytes.Repeat([]byte{2}, common.HashSize),
			},
		},
	})
	s.Require().NoError(err)

	err = stream.Send(&mgrpc.ReportBlockEventsRequest{
		MsgType: &mgrpc.ReportBlockEventsRequest_Event{
			Event: &mgrpc.Event{
				MonitorId: reg.Monitors.Monitors[0].Id,
			},
		},
	})
	s.Require().NoError(err)

	_, err = stream.CloseAndRecv()
	s.Require().NoError(err)
}

func (s *OrchestratorSuite) TestUpdateMonitors() {
	s.fm.SetChainConfig(s.cm.chainID, nil, nil)

	reg, err := s.registerCM()
	s.Require().NoError(err)
	s.Empty(reg.Monitors.Monitors)
	monitorsVersion := reg.Monitors.Version

	setReqs := make(chan *mgrpc.SetMonitorsRequest)
	updateReqs := make(chan *mgrpc.UpdateMonitorsRequest)

	s.cm.SetMonitorsFunc = func(stream mgrpc.ChainManager_SetMonitorsServer) error {
		msg, err := stream.Recv()
		s.Require().NoError(err)

		preflight := msg.GetPreflight()
		s.Require().NotNil(preflight)
		s.True(bytes.Equal(preflight.SessionId, s.cm.sessionID))
		s.True(preflight.MonitorSetVersion > monitorsVersion)
		monitorsVersion = preflight.MonitorSetVersion

		for {
			msg, err = stream.Recv()
			if err == io.EOF {
				return stream.SendAndClose(new(empty.Empty))
			}
			s.Require().NoError(err)
			setReqs <- msg
		}
	}
	s.cm.UpdateMonitorsFunc = func(stream mgrpc.ChainManager_UpdateMonitorsServer) error {
		msg, err := stream.Recv()
		s.Require().NoError(err)

		preflight := msg.GetPreflight()
		s.Require().NotNil(preflight)
		s.True(bytes.Equal(preflight.SessionId, s.cm.sessionID))
		s.Equal(monitorsVersion, preflight.PreviousMonitorSetVersion)
		s.True(preflight.MonitorSetVersion > preflight.PreviousMonitorSetVersion)
		monitorsVersion = preflight.MonitorSetVersion

		for {
			msg, err = stream.Recv()
			if err == io.EOF {
				return stream.SendAndClose(new(empty.Empty))
			}
			s.Require().NoError(err)
			updateReqs <- msg
		}
	}

	assertNoMoreReqs := func() {
		select {
		case req := <-setReqs:
			s.FailNowf("Got unexpected set req", "Req: %v", req)
		case req := <-updateReqs:
			s.FailNowf("Got unexpected update req", "Req: %v", req)
		default:
		}
	}

	// Scenario 1 - single addition
	ims := IndexedMonitors{1: &mgrpc.Monitor{Id: 1}}
	s.fm.SetChainConfig(s.cm.chainID, ims, nil)
	select {
	case req := <-updateReqs:
		s.True(proto.Equal(ims[1], req.GetAddMonitor().Monitor))
	case <-time.After(5 * time.Second):
		s.FailNow("Timeout waiting for update req")
	}
	assertNoMoreReqs()

	// Scenario 2 - single addition & single removal
	ims = IndexedMonitors{2: &mgrpc.Monitor{Id: 2}}
	s.fm.SetChainConfig(s.cm.chainID, ims, nil)
	select {
	case req := <-updateReqs:
		s.True(proto.Equal(ims[2], req.GetAddMonitor().Monitor))
		req = <-updateReqs
		s.Equal(int64(1), req.GetRemoveMonitor().MonitorId)
	case <-time.After(5 * time.Second):
		s.FailNow("Timeout waiting for update req")
	}
	assertNoMoreReqs()

	// Scenario 3 - forced reset
	ims = IndexedMonitors{3: &mgrpc.Monitor{Id: 3}}
	s.cm.UpdateMonitorsFunc = func(stream mgrpc.ChainManager_UpdateMonitorsServer) error {
		return status.Error(codes.Aborted, "Injected error")
	}
	s.fm.SetChainConfig(s.cm.chainID, ims, nil)
	select {
	case req := <-setReqs:
		s.True(proto.Equal(ims[3], req.GetMonitor()))
	case <-time.After(5 * time.Second):
		s.FailNow("Timeout waiting for set req")
	}
	assertNoMoreReqs()
}

func (s *OrchestratorSuite) TestRenewLease_Expired() {
	s.orch.leaseDuration = 1 * time.Second
	s.fm.SetChainConfig(s.cm.chainID, nil, nil)

	reg, err := s.registerCM()
	s.Require().NoError(err)

	time.Sleep(2 * time.Second)

	_, err = s.orchClient.RenewLease(s.ctx, &mgrpc.RenewLeaseRequest{
		LeaseId: reg.Lease.Id,
	})
	s.Require().Error(err)
	status, ok := status.FromError(err)
	s.Require().True(ok)
	s.Equal(codes.Aborted, status.Code())
}

func (s *OrchestratorSuite) TestUpdateMonitors_Fail() {
	s.fm.SetChainConfig(s.cm.chainID, nil, nil)

	_, err := s.registerCM()
	s.Require().NoError(err)
	cm := s.orch.chainToLease[s.cm.chainID].chainManager

	s.cm.UpdateMonitorsFunc = func(stream mgrpc.ChainManager_UpdateMonitorsServer) error {
		return status.Error(codes.FailedPrecondition, "Injected error")
	}

	ims := IndexedMonitors{1: &mgrpc.Monitor{Id: 1}}
	s.fm.SetChainConfig(s.cm.chainID, ims, nil)

	select {
	case <-cm.done:
	case <-time.After(5 * time.Second):
		s.FailNow("Timeout waiting for closing cm")
	}
}

type fakeChainManager struct {
	ChainManagerServerMock

	chainID    string
	cmID       *mgrpc.InstanceId
	listenPort int
	sessionID  []byte
}

func TestOrchestrator(t *testing.T) {
	suite.Run(t, new(OrchestratorSuite))
}

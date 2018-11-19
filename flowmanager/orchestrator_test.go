package flowmanager

//go:generate moq -out mock_chainmanagerserver_test.go -pkg flowmanager ../grpc ChainManagerServer

import (
	"bytes"
	"context"
	"io"
	"net"
	"testing"
	"time"

	"github.com/megaspacelab/megaconnect/common"
	mgrpc "github.com/megaspacelab/megaconnect/grpc"
	"github.com/megaspacelab/megaconnect/protos"

	"github.com/golang/protobuf/proto"
	"github.com/golang/protobuf/ptypes/empty"
	"github.com/google/uuid"
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

	stateStore *MemStateStore
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

	s.stateStore = NewMemStateStore()
	s.fm = NewFlowManager(s.stateStore, s.log)
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
	resp, err := s.orchClient.RegisterChainManager(s.ctx, &mgrpc.RegisterChainManagerRequest{
		ChainId:        s.cm.chainID,
		ChainManagerId: s.cm.cmID,
		ListenPort:     int32(s.cm.listenPort),
		SessionId:      s.cm.sessionID,
	})
	if err != nil {
		return nil, err
	}
	s.cm.leaseID = resp.Lease.Id
	return resp, err
}

func (s *OrchestratorSuite) unregisterCM(leaseID []byte) error {
	if leaseID == nil {
		leaseID = s.cm.leaseID[:]
	}
	_, err := s.orchClient.UnregsiterChainManager(s.ctx, &mgrpc.UnregisterChainManagerRequest{
		ChainId: s.cm.chainID,
		LeaseId: leaseID,
		Message: "Unhealthy",
	})
	return err
}

func (s *OrchestratorSuite) addMonitor(name string, bytes []byte) *mgrpc.Monitor {
	monitorID := mgrpc.MonitorID{MonitorName: name}
	monitor := &mgrpc.Monitor{
		WorkflowId:  monitorID.WorkflowID.UnsafeBytes(),
		MonitorName: monitorID.MonitorName,
		Monitor:     bytes,
	}

	if s.stateStore.committed.monitorByID == nil {
		s.stateStore.committed.monitorByID = make(map[mgrpc.MonitorID]*mgrpc.Monitor)
	}
	s.stateStore.committed.monitorByID[monitorID] = monitor

	if s.stateStore.committed.monitorsByChain == nil {
		s.stateStore.committed.monitorsByChain = make(map[string][]*mgrpc.MonitorID)
	}
	s.stateStore.committed.monitorsByChain[s.cm.chainID] = []*mgrpc.MonitorID{&monitorID}

	s.fm.lock.Lock()
	defer s.fm.lock.Unlock()

	patch := &ChainConfigPatch{
		AddMonitors: []*mgrpc.Monitor{monitor},
		chainConfigBase: chainConfigBase{
			MonitorsVersion: s.fm.chainConfigVersions[s.cm.chainID] + 1,
		},
	}
	s.fm.publishChainConfigPatchWithLock(s.cm.chainID, patch)

	return monitor
}

func (s *OrchestratorSuite) removeMonitor(name string) {
	monitorID := mgrpc.MonitorID{MonitorName: name}

	if s.stateStore.committed.monitorByID[monitorID] != nil {
		delete(s.stateStore.committed.monitorByID, monitorID)
	}
	if mids := s.stateStore.committed.monitorsByChain[s.cm.chainID]; mids != nil {
		filtered := mids[:0]
		for _, mid := range mids {
			if *mid != monitorID {
				filtered = append(filtered, mid)
			}
		}
		s.stateStore.committed.monitorsByChain[s.cm.chainID] = mids
	}

	s.fm.lock.Lock()
	defer s.fm.lock.Unlock()

	patch := &ChainConfigPatch{
		RemoveMonitors: []*mgrpc.MonitorID{&monitorID},
		chainConfigBase: chainConfigBase{
			MonitorsVersion: s.fm.chainConfigVersions[s.cm.chainID] + 1,
		},
	}
	s.fm.publishChainConfigPatchWithLock(s.cm.chainID, patch)
}

func (s *OrchestratorSuite) TestRegisterUnregisterChainManager() {
	monitor := s.addMonitor("1", []byte("1"))

	reg, err := s.registerCM()
	s.Require().NoError(err)

	s.True(reg.Lease.RemainingSeconds > 0)
	s.Equal(1, len(reg.Monitors.Monitors))
	s.True(proto.Equal(monitor, reg.Monitors.Monitors[0]))

	// test with non valid leaseID should fail
	err = s.unregisterCM([]byte("jibberish"))
	s.Require().Error(err, "Invalid LeaseId")

	// temper with wrong valid leaseID should fail
	newLeaseID := uuid.New()
	err = s.unregisterCM(newLeaseID[:])
	s.Require().Error(err, "Lease doesn't match, doesn't exist or has already expired")

	// send coorect leaseID should succeed
	err = s.unregisterCM(nil)
	s.Require().NoError(err)
}

func (s *OrchestratorSuite) TestRenewLease() {
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
	monitor := s.addMonitor("1", []byte("1"))

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
			Block: &protos.Block{
				Hash: bytes.Repeat([]byte{1}, common.HashSize),
				Height: &protos.BigInt{
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
			Event: &protos.MonitorEvent{
				WorkflowId:  monitor.WorkflowId,
				MonitorName: monitor.MonitorName,
			},
		},
	})
	s.Require().NoError(err)

	_, err = stream.CloseAndRecv()
	s.Require().NoError(err)
}

func (s *OrchestratorSuite) TestUpdateMonitors() {
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
	monitor := s.addMonitor("1", []byte("1"))
	select {
	case req := <-updateReqs:
		s.True(proto.Equal(monitor, req.GetAddMonitor().Monitor))
	case <-time.After(5 * time.Second):
		s.FailNow("Timeout waiting for update req")
	}
	assertNoMoreReqs()

	// Scenario 2 - single removal
	s.removeMonitor("1")
	select {
	case req := <-updateReqs:
		s.Equal("1", string(req.GetRemoveMonitor().MonitorName))
	case <-time.After(5 * time.Second):
		s.FailNow("Timeout waiting for update req")
	}
	assertNoMoreReqs()

	// Scenario 3 - forced reset
	s.cm.UpdateMonitorsFunc = func(stream mgrpc.ChainManager_UpdateMonitorsServer) error {
		return status.Error(codes.Aborted, "Injected error")
	}
	monitor = s.addMonitor("2", []byte("2"))
	select {
	case req := <-setReqs:
		s.True(proto.Equal(monitor, req.GetMonitor()))
	case <-time.After(5 * time.Second):
		s.FailNow("Timeout waiting for set req")
	}
	assertNoMoreReqs()
}

func (s *OrchestratorSuite) TestRenewLease_Expired() {
	s.orch.leaseDuration = 1 * time.Second

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
	_, err := s.registerCM()
	s.Require().NoError(err)
	cm := s.orch.chainToLease[s.cm.chainID].chainManager

	s.cm.UpdateMonitorsFunc = func(stream mgrpc.ChainManager_UpdateMonitorsServer) error {
		return status.Error(codes.FailedPrecondition, "Injected error")
	}

	s.addMonitor("1", []byte("1"))

	select {
	case <-cm.done:
	case <-time.After(5 * time.Second):
		s.FailNow("Timeout waiting for closing cm")
	}
}

type fakeChainManager struct {
	ChainManagerServerMock

	chainID    string
	leaseID    []byte
	cmID       *mgrpc.InstanceId
	listenPort int
	sessionID  []byte
}

func TestOrchestrator(t *testing.T) {
	suite.Run(t, new(OrchestratorSuite))
}

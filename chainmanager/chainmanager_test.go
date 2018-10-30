// Copyright 2018 @ MegaSpace

// The MegaSpace source code is distributed in the hope that it will be useful,
// but WITHOUT ANY WARRANTY; without even the implied warranty of
// MERCHANTABILITY or FITNESS FOR A PARTICULAR PURPOSE. See the
// GNU Lesser General Public License for more details.

// You should have received a copy of the GNU Lesser General Public License
// along with the MegaSpace source code. If not, see <http://www.gnu.org/licenses/>.

package chainmanager_test

//go:generate moq -out mock_orchestratorserver_test.go -pkg chainmanager ../grpc OrchestratorServer

import (
	"context"
	"encoding/hex"
	"fmt"
	"math/big"
	"net"
	"sync"
	"syscall"
	"testing"
	"time"

	"github.com/golang/protobuf/ptypes/empty"
	"github.com/google/uuid"
	. "github.com/megaspacelab/megaconnect/chainmanager"
	mcli "github.com/megaspacelab/megaconnect/chainmanager/cli"
	"github.com/megaspacelab/megaconnect/common"
	"github.com/megaspacelab/megaconnect/connector"
	"github.com/megaspacelab/megaconnect/connector/example"
	mgrpc "github.com/megaspacelab/megaconnect/grpc"
	"github.com/stretchr/testify/suite"
	"go.uber.org/zap"
	"google.golang.org/grpc"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
	cli "gopkg.in/urfave/cli.v2"
)

const (
	blockInterval      = 100 * time.Millisecond
	leaseRenewalBuffer = 5 * time.Second
)

var badHash = common.Hash{1}

type TestConnector struct {
	connector.Connector

	// How long we want QueryAccountBalance to take, in seconds
	queryAccountBalanceMinDuration time.Duration

	// set healthy state for test connector
	healthy bool

	// concurrency
	lock sync.Mutex
}

func (t *TestConnector) QueryAccountBalance(addr string, height *big.Int) (*big.Int, error) {
	t.lock.Lock()
	defer t.lock.Unlock()
	time.Sleep(t.queryAccountBalanceMinDuration)
	return t.Connector.QueryAccountBalance(addr, height)
}

func (t *TestConnector) SetAccountBalanceMinDuration(value time.Duration) {
	t.lock.Lock()
	defer t.lock.Unlock()
	t.queryAccountBalanceMinDuration = value
}

func (t *TestConnector) IsHealthy() (bool, error) {
	t.lock.Lock()
	defer t.lock.Unlock()
	return t.healthy, nil
}

func (t *TestConnector) SetHealthy(value bool) {
	t.lock.Lock()
	defer t.lock.Unlock()
	t.healthy = value
}

type ChainManagerSuite struct {
	suite.Suite

	log      *zap.Logger
	ctx      context.Context
	monitors []*mgrpc.Monitor

	server     *grpc.Server
	conn       *grpc.ClientConn
	listenAddr *net.TCPAddr

	orch      *fakeOrchestrator
	connector *TestConnector
	cm        *ChainManager
	cmClient  mgrpc.ChainManagerClient
}

func (s *ChainManagerSuite) SetupTest() {
	log, err := zap.NewDevelopment()
	s.Require().NoError(err)
	s.log = log
	s.ctx = context.Background()

	for i, hexStr := range []string{
		"04546573740401060a47657442616c616e636502022a307846426231623733433466304244613466363764634132363663653645663432663532306642423938080b626c6f636b4865696768740103457468060a47657442616c616e636502022a307846426231623733433466304244613466363764634132363663653645663432663532306642423938040a080b626c6f636b486569676874010101010103457468010b626c6f636b486569676874070608476574426c6f636b00010345746806686569676874",
		"04546573740401060a47657442616c616e636502022a307838373665616266343431623265653562356230353534666435303261386530363030393530636661080b626c6f636b4865696768740103457468060a47657442616c616e636502022a307838373665616266343431623265653562356230353534666435303261386530363030393530636661040a080b626c6f636b486569676874010101010103457468010b626c6f636b486569676874070608476574426c6f636b00010345746806686569676874",
	} {
		monitor, err := parseMonitor(int64(i), hexStr)
		s.Require().NoError(err)
		s.monitors = append(s.monitors, monitor)
	}

	lis, err := net.Listen("tcp", "localhost:")
	s.Require().NoError(err)
	s.listenAddr = lis.Addr().(*net.TCPAddr)

	s.server = grpc.NewServer()
	s.conn, err = grpc.Dial(s.listenAddr.String(), grpc.WithInsecure())
	s.Require().NoError(err)

	s.orch = &fakeOrchestrator{
		leaseSeconds: 60,

		OrchestratorServerMock: OrchestratorServerMock{
			RegisterChainManagerFunc: func(
				ctx context.Context,
				req *mgrpc.RegisterChainManagerRequest,
			) (*mgrpc.RegisterChainManagerResponse, error) {
				s.orch.lock.Lock()
				defer s.orch.lock.Unlock()

				s.orch.sessionID = req.SessionId
				leaseID := uuid.New()
				s.orch.leaseID = &leaseID
				return &mgrpc.RegisterChainManagerResponse{
					Lease: &mgrpc.Lease{
						Id:               leaseID[:],
						RemainingSeconds: s.orch.leaseSeconds,
					},
					Monitors: &mgrpc.MonitorSet{
						Monitors: s.monitors[:1], // include one monitor by default
					},
				}, nil
			},
			UnregsiterChainManagerFunc: func(
				in1 context.Context,
				in2 *mgrpc.UnregisterChainManagerRequest,
			) (*empty.Empty, error) {
				if in2.ChainId == "" {
					return nil, status.Error(codes.InvalidArgument, "Missing ChainId")
				}
				s.orch.lock.Lock()
				defer s.orch.lock.Unlock()
				s.orch.unregisterCallCount++
				s.log.Debug("Unregistered Chain Manager")
				return &empty.Empty{}, nil
			},
			ReportBlockEventsFunc: func(stream mgrpc.Orchestrator_ReportBlockEventsServer) error {
				msg, err := stream.Recv()
				s.Require().NoError(err)
				preflight := msg.GetPreflight()
				s.Require().NotNil(preflight)
				leaseID, err := uuid.FromBytes(preflight.LeaseId)
				s.Require().NoError(err)

				s.orch.lock.Lock()
				defer s.orch.lock.Unlock()

				s.Equal(s.orch.leaseID, &leaseID)
				s.Condition(func() bool { return preflight.MonitorSetVersion <= s.orch.monitorsVersion })

				if preflight.MonitorSetVersion < s.orch.monitorsVersion {
					return status.Error(codes.Aborted, "Stale MonitorSetVersion")
				}

				msg, err = stream.Recv()
				s.Require().NoError(err)
				block := msg.GetBlock()
				s.Require().NotNil(block)
				s.log.Debug("Received block", zap.Stringer("block", block))

				height := block.Height.BigInt()
				if s.orch.blockHeight != nil {
					s.orch.blockHeight.Add(s.orch.blockHeight, big.NewInt(1))
					s.Equal(s.orch.blockHeight, height)
				}
				s.orch.blockHeight = height

				s.orch.receivedBlocks++
				return stream.SendAndClose(new(empty.Empty))
			},
			RenewLeaseFunc: func(ctx context.Context, req *mgrpc.RenewLeaseRequest) (*mgrpc.Lease, error) {
				leaseID, err := uuid.FromBytes(req.LeaseId)
				s.Require().NoError(err)

				s.orch.lock.Lock()
				defer s.orch.lock.Unlock()

				s.Equal(s.orch.leaseID, &leaseID)

				leaseID = uuid.New()
				s.orch.leaseID = &leaseID
				s.orch.leaseRenewals++
				return &mgrpc.Lease{
					Id:               leaseID[:],
					RemainingSeconds: s.orch.leaseSeconds,
				}, nil
			},
		},
	}
	mgrpc.RegisterOrchestratorServer(s.server, s.orch)

	exampleConnector, err := example.New(s.log, blockInterval)
	s.Require().NoError(err)

	s.connector = &TestConnector{
		Connector:                      exampleConnector,
		queryAccountBalanceMinDuration: 0,
		healthy:                        true,
	}

	s.cm = New(
		"test-cm",
		s.listenAddr.String(),
		s.connector,
		s.log,
	)
	s.cm.Register(s.server)
	s.cmClient = mgrpc.NewChainManagerClient(s.conn)

	go s.server.Serve(lis)
}

func (s *ChainManagerSuite) TearDownTest() {
	if s.conn != nil {
		s.conn.Close()
	}
	if s.server != nil {
		s.server.Stop()
	}
}

func (s *ChainManagerSuite) TestStartStop() {
	err := s.cm.Start(s.listenAddr.Port)
	s.Require().NoError(err)

	time.Sleep(2 * blockInterval)

	err = s.cm.Stop()
	s.NoError(err)

	s.Condition(func() bool { return s.orch.receivedBlocks > 0 })
}

func (s *ChainManagerSuite) TestSetMonitors() {
	err := s.cm.Start(s.listenAddr.Port)
	s.Require().NoError(err)

	time.Sleep(2 * blockInterval)

	s.orch.lock.Lock()
	receivedBlocks := s.orch.receivedBlocks
	s.Require().Condition(func() bool { return receivedBlocks > 0 })
	s.Require().NotNil(s.orch.blockHeight)
	height := new(big.Int).Set(s.orch.blockHeight)
	s.orch.monitorsVersion++
	s.orch.lock.Unlock()

	time.Sleep(2 * blockInterval)

	stream, err := s.cmClient.SetMonitors(s.ctx)
	s.Require().NoError(err)

	err = stream.Send(&mgrpc.SetMonitorsRequest{
		MsgType: &mgrpc.SetMonitorsRequest_Preflight_{
			Preflight: &mgrpc.SetMonitorsRequest_Preflight{
				MonitorSetVersion: s.orch.monitorsVersion,
				SessionId:         s.orch.sessionID,
				ResumeAfter: &mgrpc.BlockSpec{
					Hash:   badHash.Bytes(),
					Height: mgrpc.NewBigInt(height),
				},
			},
		},
	})
	s.Require().NoError(err)

	for _, monitor := range s.monitors {
		err = stream.Send(&mgrpc.SetMonitorsRequest{
			MsgType: &mgrpc.SetMonitorsRequest_Monitor{
				Monitor: monitor,
			},
		})
		s.Require().NoError(err)
	}

	_, err = stream.CloseAndRecv()
	s.Require().NoError(err)

	time.Sleep(2 * blockInterval)

	err = s.cm.Stop()
	s.NoError(err)

	// Verify that new blocks have been received after monitors update
	s.Condition(func() bool { return s.orch.receivedBlocks > receivedBlocks })
}

func (s *ChainManagerSuite) TestUpdateMonitors() {
	err := s.cm.Start(s.listenAddr.Port)
	s.Require().NoError(err)

	s.orch.lock.Lock()
	receivedBlocks := s.orch.receivedBlocks
	s.orch.monitorsVersion++
	s.orch.lock.Unlock()

	stream, err := s.cmClient.UpdateMonitors(s.ctx)
	s.Require().NoError(err)

	err = stream.Send(&mgrpc.UpdateMonitorsRequest{
		MsgType: &mgrpc.UpdateMonitorsRequest_Preflight_{
			Preflight: &mgrpc.UpdateMonitorsRequest_Preflight{
				MonitorSetVersion:         s.orch.monitorsVersion,
				PreviousMonitorSetVersion: s.orch.monitorsVersion - 1,
				SessionId:                 s.orch.sessionID,
			},
		},
	})
	s.Require().NoError(err)

	// Add & remove monitor
	err = stream.Send(&mgrpc.UpdateMonitorsRequest{
		MsgType: &mgrpc.UpdateMonitorsRequest_AddMonitor_{
			AddMonitor: &mgrpc.UpdateMonitorsRequest_AddMonitor{
				Monitor: s.monitors[1],
			},
		},
	})
	s.Require().NoError(err)
	err = stream.Send(&mgrpc.UpdateMonitorsRequest{
		MsgType: &mgrpc.UpdateMonitorsRequest_RemoveMonitor_{
			RemoveMonitor: &mgrpc.UpdateMonitorsRequest_RemoveMonitor{
				MonitorId: s.monitors[0].Id,
			},
		},
	})
	s.Require().NoError(err)

	_, err = stream.CloseAndRecv()
	s.Require().NoError(err)

	time.Sleep(2 * blockInterval)

	err = s.cm.Stop()
	s.NoError(err)

	// Verify that new blocks have been received after monitors update
	s.Condition(func() bool { return s.orch.receivedBlocks > receivedBlocks })
}

func (s *ChainManagerSuite) TestRenewLease() {
	s.orch.leaseSeconds = uint32((leaseRenewalBuffer + time.Second).Seconds())

	err := s.cm.Start(s.listenAddr.Port)
	s.Require().NoError(err)

	time.Sleep(2 * time.Second)

	err = s.cm.Stop()
	s.NoError(err)

	s.Condition(func() bool { return s.orch.receivedBlocks > 0 })
	s.Condition(func() bool { return s.orch.leaseRenewals > 0 })
}

func (s *ChainManagerSuite) TestHealth() {
	// set connector to be not healthy to begin with and launching the chain manager would fail
	s.connector.SetHealthy(false)
	err := s.cm.Start(s.listenAddr.Port)
	s.EqualError(err, "Initialization error, not healthy")
	err = s.cm.Stop()
	s.NoError(err)

	// launch it again, this time, set the health to be healthy 0.2 seconds after start is called,
	// chain manager should start just fine
	s.connector.SetHealthy(false)
	go func() {
		time.Sleep(200 * time.Millisecond)
		s.connector.SetHealthy(true)
	}()
	err = s.cm.Start(s.listenAddr.Port)
	s.NoError(err)
	err = s.cm.Stop()
	s.NoError(err)

	// launch chain manager with a healthy connector, after it is running, make connector go bad
	// then within grace period, make it be good again, things should run just fine.
	s.connector.SetHealthy(true)
	prevCallCount := s.orch.unregisterCallCount
	err = s.cm.Start(s.listenAddr.Port)
	s.NoError(err)
	time.Sleep(200 * time.Millisecond)

	go func() {
		s.connector.SetHealthy(false) // setting the health to be false
		time.Sleep(s.connector.Metadata().HealthCheckGracePeriod / 2)
		s.connector.SetHealthy(true) // within grace period, set the health back to true
	}()

	time.Sleep(s.connector.Metadata().HealthCheckGracePeriod) // sleep for a little over grace period
	time.Sleep(2 * time.Second)

	s.True(s.orch.unregisterCallCount-prevCallCount == 0) // unregister should never been called
	err = s.cm.Stop()                                     // calling stop on connector, it should be closed without error
	s.NoError(err)

}

func (s *ChainManagerSuite) TestHealthFail() {
	// launch chain manager with a healthy connector, after it is running, make connector go bad
	// chain manager should call UnregisterChainManager with orch and removes itself gracefully
	prevCallCount := s.orch.unregisterCallCount
	s.connector.SetHealthy(true)
	err := s.cm.Start(s.listenAddr.Port)
	s.NoError(err)
	time.Sleep(200 * time.Millisecond)

	go func() { // setting the health to be false
		s.connector.SetHealthy(false)
	}()

	time.Sleep(s.connector.Metadata().HealthCheckGracePeriod) // sleep for a little over grace period
	time.Sleep(time.Second)

	s.orch.lock.Lock()
	s.True(s.orch.unregisterCallCount-prevCallCount == 1) // unregister should have been called
	s.orch.lock.Unlock()
}

func (s *ChainManagerSuite) TestRenewLeaseWhileProcessingBlock() {
	s.orch.leaseSeconds = uint32((leaseRenewalBuffer + time.Second).Seconds())

	err := s.cm.Start(s.listenAddr.Port)
	s.Require().NoError(err)

	s.connector.SetAccountBalanceMinDuration(2 * time.Second)

	time.Sleep(2 * time.Second)

	err = s.cm.Stop()
	s.NoError(err)

	s.Condition(func() bool { return s.orch.leaseRenewals > 1 })
}

func (s *ChainManagerSuite) TestRunner() {
	runner := mcli.NewRunner(func(ctx *cli.Context, logger *zap.Logger) (connector.Connector, error) {
		return example.New(logger, blockInterval)
	})
	runner.WithFlag(&cli.StringFlag{Name: "not-used"})

	// Launch it in background
	done := make(chan common.Nothing)
	go func() {
		(*cli.App)(runner).Run([]string{
			"test-runner",
			"--orch-addr", fmt.Sprintf("localhost:%d", s.listenAddr.Port),
			"--listen-addr", "localhost:0", // Bind to localhost only. This prevents a nasty firewall popup on Mac.
		})
		close(done)
	}()

	time.Sleep(2 * time.Second)

	// Signal the runner to stop
	// TODO - this doesn't seem to work when debugging. why?
	syscall.Kill(syscall.Getpid(), syscall.SIGTERM)

	select {
	case <-done:
	case <-time.After(5 * time.Second):
		s.FailNow("Timeout waiting for runner")
	}

	s.Condition(func() bool {
		s.orch.lock.Lock()
		defer s.orch.lock.Unlock()
		return s.orch.receivedBlocks > 0
	})
}

func parseMonitor(id int64, hexStr string) (*mgrpc.Monitor, error) {
	monitorRaw, err := hex.DecodeString(hexStr)
	if err != nil {
		return nil, err
	}
	return &mgrpc.Monitor{
		Id:      id,
		Monitor: monitorRaw,
	}, nil
}

type fakeOrchestrator struct {
	OrchestratorServerMock

	leaseSeconds        uint32
	sessionID           []byte
	leaseID             *uuid.UUID
	monitorsVersion     uint32
	receivedBlocks      int
	blockHeight         *big.Int
	leaseRenewals       int
	unregisterCallCount int

	lock sync.Mutex
}

func TestChainManager(t *testing.T) {
	suite.Run(t, new(ChainManagerSuite))
}

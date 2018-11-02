// Copyright 2018 @ MegaSpace

// The MegaSpace source code is distributed in the hope that it will be useful,
// but WITHOUT ANY WARRANTY; without even the implied warranty of
// MERCHANTABILITY or FITNESS FOR A PARTICULAR PURPOSE. See the
// GNU Lesser General Public License for more details.

// You should have received a copy of the GNU Lesser General Public License
// along with the MegaSpace source code. If not, see <http://www.gnu.org/licenses/>.

package chainmanager

import (
	"bytes"
	"container/ring"
	"context"
	"encoding/hex"
	"errors"
	"io"
	"sync"
	"time"

	"github.com/megaspacelab/megaconnect/common"
	"github.com/megaspacelab/megaconnect/connector"
	mgrpc "github.com/megaspacelab/megaconnect/grpc"
	"github.com/megaspacelab/megaconnect/unsafe"
	wf "github.com/megaspacelab/megaconnect/workflow"

	"github.com/golang/protobuf/ptypes/empty"
	"github.com/google/uuid"
	"go.uber.org/zap"
	"google.golang.org/grpc"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

const (
	blockChanSize        = 100
	blockCacheSize       = blockChanSize
	leaseRenewalBuffer   = 5 * time.Second
	interpreterCacheSize = 1000
)

// ChainManager manages and interacts with a chain through connector.
// Itself is being orchestrated by the Orchestrator.
type ChainManager struct {
	// Static members that don't change after construction.
	id        string
	orchAddr  string
	connector connector.Connector
	logger    *zap.Logger

	// Dynamic members.
	ctx               context.Context
	cancel            context.CancelFunc
	orchClient        mgrpc.OrchestratorClient
	orchConn          *grpc.ClientConn
	sessionID         uuid.UUID
	instance          uint32
	leaseID           []byte
	leaseRenewalTimer *time.Timer
	monitors          map[string]*mgrpc.Monitor
	monitorsVersion   uint32
	blockCache        *ring.Ring

	healthCheckTimer  *time.Timer // timer to check health of connector
	healthLastUpdated time.Time   // last known time of the connector health condition change
	healthLastState   bool        // last state of connector, true is healthy, false is not healthy

	// Concurrency control.
	running bool

	// Locks for concurrency control.
	// When acquiring lock, should always follow sequence:
	// 1, lock
	// 2, leaseLock

	// Lock guarding general stuffs like
	// running, blockCache, monitors, etc.
	lock sync.Mutex

	// Lock guarding lease.
	leaseLock sync.Mutex
}

// New constructs an instance of ChainManager.
func New(
	id string,
	orchAddr string,
	conn connector.Connector,
	logger *zap.Logger,
) *ChainManager {
	return &ChainManager{
		id:         id,
		orchAddr:   orchAddr,
		connector:  conn,
		logger:     logger,
		monitors:   make(map[string]*mgrpc.Monitor),
		blockCache: ring.New(blockCacheSize),
	}
}

// Start would start an ChainManager loop
func (e *ChainManager) Start(listenPort int) error {
	e.lock.Lock()
	defer e.lock.Unlock()
	e.leaseLock.Lock()
	defer e.leaseLock.Unlock()

	if e.running {
		return errors.New("ChainManager is already running")
	}

	e.logger.Debug("Connecting to Orchestrator", zap.String("orchAddr", e.orchAddr))
	e.ctx, e.cancel = context.WithCancel(context.Background())
	orchConn, err := grpc.DialContext(e.ctx, e.orchAddr, grpc.WithInsecure())
	if err != nil {
		return err
	}
	e.orchConn = orchConn
	e.orchClient = mgrpc.NewOrchestratorClient(orchConn)
	e.sessionID = uuid.New()
	e.instance++

	err = e.connector.Start()
	if err != nil {
		return err
	}
	metadata := e.connector.Metadata()

	// TODO Add healthy endpoint in CM and start async for HA
	// initial health check, wait till connector is healthy
	for {
		e.healthLastState, err = e.connector.IsHealthy()
		if e.healthLastState && err == nil {
			break
		}
		if err != nil {
			e.logger.Debug("health check error", zap.Error(err))
		}
		e.logger.Info("Connector is not healthy, waiting for it to turn healthy...")
		time.Sleep(metadata.HealthCheckInterval)
	}

	// health check is finished and connector is healthy, register chainmanager and start monitoring health repeatedly
	e.logger.Debug("Register with Orchestrator", zap.String("chain ID", metadata.ChainID))
	resp, err := e.orchClient.RegisterChainManager(e.ctx, &mgrpc.RegisterChainManagerRequest{
		ChainId:        metadata.ChainID,
		ChainManagerId: &mgrpc.InstanceId{Id: unsafe.StringToBytes(e.id), Instance: e.instance},
		ListenPort:     int32(listenPort),
		SessionId:      e.sessionID[:],
	})
	if err != nil {
		return err
	}

	// register with orchestrator
	e.logger.Debug("Registered with Orchestrator", zap.Stringer("lease", resp.Lease))
	e.updateLeaseWithLeaseLock(resp.Lease)

	e.logger.Debug("Start monitoring health of connected blockchain", zap.String("name", metadata.ConnectorID))
	e.checkHealthWithLock()

	e.monitorsVersion = resp.Monitors.GetVersion()
	for _, m := range resp.Monitors.GetMonitors() {
		e.logger.Debug("Monitor received", zap.String("monitor", hex.EncodeToString(m.Monitor)))
		e.monitors[string(m.Id)] = m
	}

	resumeAfter, err := convertBlockSpec(resp.ResumeAfter)
	if err != nil {
		return err
	}

	sub, err := e.connector.SubscribeBlock(resumeAfter)
	if err != nil {
		return err
	}

	go func() {
		defer e.logger.Info("Block subscription stopped")
		for {
			select {
			case block, ok := <-sub.Blocks():
				if !ok {
					return
				}
				err := e.processNewBlock(block)
				if err != nil {
					e.logger.Fatal("Failed to process block", zap.Error(err))
					panic(err)
				}
			case err := <-sub.Err():
				if err != nil {
					e.logger.Fatal("Block subscription failed", zap.Error(err))
					panic(err)
				}
				return
			}
		}
	}()

	e.running = true
	return nil
}

// Stop would stop an ChainManager loop
func (e *ChainManager) Stop() error {
	e.lock.Lock()
	defer e.lock.Unlock()
	e.leaseLock.Lock()
	defer e.leaseLock.Unlock()

	if !e.running {
		return errors.New("event manager is not running")
	}

	e.running = false

	err := e.connector.Stop()
	if err != nil {
		e.logger.Error("connect stop with err", zap.Error(err))
	}

	if e.healthCheckTimer != nil {
		e.healthCheckTimer.Stop()
		e.healthCheckTimer = nil
	}

	if e.leaseRenewalTimer != nil {
		e.leaseRenewalTimer.Stop()
		e.leaseRenewalTimer = nil
	}

	e.cancel()

	err = e.orchConn.Close()
	if err != nil {
		e.logger.Error("orchConn closed with err", zap.Error(err))
	}

	e.logger.Info("ChainManager stopped", zap.String("ConnectorID", e.connector.Metadata().ConnectorID))

	return nil
}

func (e *ChainManager) updateLeaseWithLeaseLock(lease *mgrpc.Lease) {
	e.leaseID = lease.Id
	timeout := time.Duration(lease.RemainingSeconds)*time.Second - leaseRenewalBuffer
	e.leaseRenewalTimer = time.AfterFunc(timeout, e.renewLease)
}

func (e *ChainManager) renewLease() {
	e.leaseLock.Lock()
	defer e.leaseLock.Unlock()

	if !e.running || e.leaseID == nil {
		e.logger.Warn("Skipping lease renewal on stopped ChainManager")
		return
	}

	resp, err := e.orchClient.RenewLease(e.ctx, &mgrpc.RenewLeaseRequest{
		LeaseId: e.leaseID,
	})
	if err != nil {
		e.logger.Fatal("Failed to renew lease. Shutting down")
		e.leaseID = nil

		// Perform stop in another go routine to respect lock order.
		// Otherwise, we risk deadlocks.
		go func() {
			e.Stop()
			panic("Failed to renew lease")
		}()
		return
	}

	e.updateLeaseWithLeaseLock(resp)
}

func (e *ChainManager) checkHealth() {
	e.lock.Lock()
	defer e.lock.Unlock()
	e.checkHealthWithLock()
}

func (e *ChainManager) checkHealthWithLock() {
	// set current state
	connectorState, err := e.connector.IsHealthy()
	currentState := connectorState && err == nil

	// if it is first run or there is a state change, update the state and time
	if e.healthCheckTimer == nil || (currentState != e.healthLastState) {
		current := time.Now()
		e.logger.Info("Health state changed",
			zap.Bool("new state", currentState),
			zap.Stringer("time", current))

		e.healthLastState = currentState
		e.healthLastUpdated = current
	}

	// if health condition is not healthy (pass grace period), unregisters the chain manager
	metadata := e.connector.Metadata()
	timeout := e.healthLastUpdated.Add(metadata.HealthCheckGracePeriod)
	needsUnregister := !e.healthLastState && time.Now().After(timeout)

	if needsUnregister {
		e.logger.Info("Connector is not healthy, unregister chain Manager",
			zap.String("ChainID", metadata.ChainID))

		e.orchClient.UnregsiterChainManager(e.ctx, &mgrpc.UnregisterChainManagerRequest{
			ChainId: metadata.ChainID,
			LeaseId: e.leaseID,
			Message: "Connected blockchain is not healthy",
		})

		go func() {
			e.Stop()

			// TODO add panic
			//panic("ChainManager self destruct due to connector health")
		}()
	} else {
		e.healthCheckTimer = time.AfterFunc(metadata.HealthCheckInterval, e.checkHealth)
	}
}

func (e *ChainManager) processNewBlock(block common.Block) error {
	e.logger.Debug("Received new block", zap.Stringer("height", block.Height()))

	e.lock.Lock()
	defer e.lock.Unlock()
	if !e.running {
		e.logger.Warn("Skipping block on stopped ChainManager")
		return nil
	}

	e.cacheBlockWithLock(block)
	return e.processBlockWithLock(block)
}

func (e *ChainManager) processBlockWithLock(block common.Block) error {
	e.logger.Debug("Processing block", zap.Stringer("height", block.Height()))

	cache, err := wf.NewFuncCallCache(interpreterCacheSize)
	if err != nil {
		e.logger.Error("Cache creation failed", zap.Error(err))
	}
	interpreter := wf.NewInterpreter(wf.NewEnv(e.connector, block, nil), cache, wf.NewResolver(nil, wf.Libs), e.logger)
	events := []*mgrpc.Event{}
	for id, monitor := range e.monitors {
		e.logger.Debug("Processing monitor", zap.Stringer("height", block.Height()), zap.String("monitor", hex.EncodeToString(monitor.Monitor)))
		md, err := wf.NewByteDecoder(monitor.Monitor).DecodeMonitorDecl()
		if err != nil {
			return err
		}

		e.logger.Debug("Evaluating condition", zap.Stringer("height", block.Height()), zap.Stringer("condition", md.Condition()))
		eventValue, err := interpreter.EvalMonitor(md)
		if err != nil {
			return err
		}
		e.logger.Debug("Evaluated condition", zap.Stringer("height", block.Height()), zap.Bool("condition", eventValue != nil))
		if eventValue == nil {
			continue
		}
		event := mgrpc.Event{}
		event.MonitorId = unsafe.StringToBytes(id)
		payload, err := wf.EncodeObjConst(eventValue.Payload())
		if err != nil {
			return err
		}
		event.Payload = payload
		events = append(events, &event)
	}

	err = e.reportBlockEventsWithLock(block, events)
	if err != nil {
		s, ok := status.FromError(err)
		if !ok || s.Code() != codes.Aborted {
			return err
		}

		// Aborted means the monitor set is stale. No need to take aggreviated actions just yet.
		// Should simply wait to receive the new monitor set.
		e.logger.Warn("Ignoring Aborted response when reporting block",
			zap.Stringer("height", block.Height()), zap.Stringer("status", s.Proto()))
		return nil
	}

	return nil
}

func (e *ChainManager) reportBlockEventsWithLock(block common.Block, events []*mgrpc.Event) error {
	e.leaseLock.Lock()
	defer e.leaseLock.Unlock()

	stream, err := e.orchClient.ReportBlockEvents(e.ctx)
	if err != nil {
		return err
	}

	err = stream.Send(&mgrpc.ReportBlockEventsRequest{
		MsgType: &mgrpc.ReportBlockEventsRequest_Preflight_{
			Preflight: &mgrpc.ReportBlockEventsRequest_Preflight{
				LeaseId:           e.leaseID,
				MonitorSetVersion: e.monitorsVersion,
			},
		},
	})
	if err != nil {
		return err
	}

	hash := block.Hash()
	parentHash := block.ParentHash()

	err = stream.Send(&mgrpc.ReportBlockEventsRequest{
		MsgType: &mgrpc.ReportBlockEventsRequest_Block{
			Block: &mgrpc.Block{
				Hash:       hash.Bytes(),
				ParentHash: parentHash.Bytes(),
				Height: &mgrpc.BigInt{
					Bytes:    block.Height().Bytes(),
					Negative: block.Height().Sign() < 0,
				},
			},
		},
	})
	if err != nil {
		return err
	}

	for _, event := range events {
		err = stream.Send(&mgrpc.ReportBlockEventsRequest{
			MsgType: &mgrpc.ReportBlockEventsRequest_Event{
				Event: event,
			},
		})
		if err != nil {
			return err
		}
	}

	_, err = stream.CloseAndRecv()
	return err
}

func (e *ChainManager) cacheBlockWithLock(block common.Block) {
	// After the update, blockCache points to the newest entry, and blockCache.Next() is the oldest.
	e.blockCache = e.blockCache.Next()
	e.blockCache.Value = block
}

func (e *ChainManager) findCachedBlockWithLock(bs *connector.BlockSpec) *ring.Ring {
	// Find by hash first.
	if bs.Hash != nil {
		// Traverse in reverse order.
		for cache, i := e.blockCache, 0; cache.Value != nil && i < blockCacheSize; cache, i = cache.Prev(), i+1 {
			if cache.Value.(common.Block).Hash() == *bs.Hash {
				return cache
			}
		}
	}

	// Now find by height.
	if bs.Height != nil {
		// Traverse in reverse order.
		// Return the first (latest) block whose height is <= bs.Height.
		for cache, i := e.blockCache, 0; cache.Value != nil && i < blockCacheSize; cache, i = cache.Prev(), i+1 {
			if cache.Value.(common.Block).Height().Cmp(bs.Height) <= 0 {
				return cache
			}
		}
	}

	return nil
}

// SetMonitors resets ChainManager's current MonitorSet.
func (e *ChainManager) SetMonitors(stream mgrpc.ChainManager_SetMonitorsServer) error {
	msg, err := stream.Recv()
	if err != nil {
		return err
	}

	preflight := msg.GetPreflight()
	if preflight == nil {
		return status.Error(codes.InvalidArgument, "Missing Preflight")
	}

	e.logger.Debug("Received SetMonitorsRequest.Preflight", zap.Stringer("preflight", preflight))

	resumeAfter, err := convertBlockSpec(preflight.ResumeAfter)
	if err != nil {
		return status.Errorf(codes.InvalidArgument, "Invalid ResumeAfter. %v", err)
	}

	e.lock.Lock()
	defer e.lock.Unlock()

	if !e.running {
		return status.Error(codes.FailedPrecondition, "ChainManager not running")
	}

	if !bytes.Equal(preflight.SessionId, e.sessionID[:]) {
		return status.Error(codes.Aborted, "SessionId mismatch")
	}
	if preflight.MonitorSetVersion <= e.monitorsVersion {
		return status.Error(codes.Aborted, "MonitorSetVersion too low")
	}

	monitors := make(map[string]*mgrpc.Monitor)

	for {
		msg, err = stream.Recv()
		if err == io.EOF {
			break
		}
		if err != nil {
			return err
		}

		monitor := msg.GetMonitor()
		if monitor == nil {
			return status.Error(codes.InvalidArgument, "Wrong message type. Expecting Monitor")
		}

		e.logger.Debug("Received SetMonitorsRequest.Monitor", zap.Stringer("monitor", monitor))
		monitors[string(monitor.Id)] = monitor
	}

	e.monitors = monitors
	e.monitorsVersion = preflight.MonitorSetVersion

	err = stream.SendAndClose(&empty.Empty{})
	e.replayWithLock(resumeAfter)

	return err
}

// UpdateMonitors patches ChainManager's current MonitorSet with additions and removals.
func (e *ChainManager) UpdateMonitors(stream mgrpc.ChainManager_UpdateMonitorsServer) error {
	msg, err := stream.Recv()
	if err != nil {
		return err
	}

	preflight := msg.GetPreflight()
	if preflight == nil {
		return status.Error(codes.InvalidArgument, "Missing Preflight")
	}

	e.logger.Debug("Received UpdateMonitorsRequest.Preflight", zap.Stringer("preflight", preflight))

	resumeAfter, err := convertBlockSpec(preflight.ResumeAfter)
	if err != nil {
		return status.Errorf(codes.InvalidArgument, "Invalid ResumeAfter. %v", err)
	}

	e.lock.Lock()
	defer e.lock.Unlock()

	if !e.running {
		return status.Error(codes.FailedPrecondition, "ChainManager not running")
	}

	if !bytes.Equal(preflight.SessionId, e.sessionID[:]) {
		return status.Error(codes.Aborted, "SessionId mismatch")
	}
	if preflight.PreviousMonitorSetVersion != e.monitorsVersion {
		return status.Error(codes.Aborted, "MonitorSetVersion mismatch")
	}

	for {
		msg, err = stream.Recv()
		if err == io.EOF {
			break
		}
		if err != nil {
			return err
		}

		switch m := msg.MsgType.(type) {
		case *mgrpc.UpdateMonitorsRequest_AddMonitor_:
			e.logger.Debug("Received UpdateMonitorsRequest.AddMonitor", zap.Stringer("m", m.AddMonitor))
			e.monitors[string(m.AddMonitor.Monitor.Id)] = m.AddMonitor.Monitor
		case *mgrpc.UpdateMonitorsRequest_RemoveMonitor_:
			e.logger.Debug("Received UpdateMonitorsRequest.RemoveMonitor", zap.Stringer("m", m.RemoveMonitor))
			delete(e.monitors, unsafe.BytesToString(m.RemoveMonitor.MonitorId))
		default:
			return status.Error(codes.InvalidArgument, "Wrong message type. Expecting AddMonitor or RemoveMonitor")
		}
	}

	e.monitorsVersion = preflight.MonitorSetVersion

	err = stream.SendAndClose(&empty.Empty{})
	e.replayWithLock(resumeAfter)

	return err
}

func (e *ChainManager) replayWithLock(resumeAfter *connector.BlockSpec) {
	if resumeAfter == nil {
		return
	}

	cache := e.findCachedBlockWithLock(resumeAfter)
	if cache == nil {
		return
	}

	for ; cache != e.blockCache; cache = cache.Next() {
		block := cache.Next().Value.(common.Block)
		e.logger.Debug("Replaying block", zap.Stringer("height", block.Height()))

		if err := e.processBlockWithLock(block); err != nil {
			e.logger.Fatal("Failed to replay block", zap.Error(err))
			panic(err)
		}
	}
}

// Register registers itself as ChainManagerServer to the gRPC server.
func (e *ChainManager) Register(server *grpc.Server) {
	mgrpc.RegisterChainManagerServer(server, e)
}

func convertBlockSpec(bs *mgrpc.BlockSpec) (*connector.BlockSpec, error) {
	if bs == nil {
		return nil, nil
	}
	hash, err := common.BytesToHash(bs.Hash)
	if err != nil {
		return nil, err
	}

	return &connector.BlockSpec{
		Hash:   hash,
		Height: bs.Height.BigInt(),
	}, nil
}

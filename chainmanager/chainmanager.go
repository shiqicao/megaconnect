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
	"errors"
	"io"
	"sync"
	"time"

	"github.com/megaspacelab/megaconnect/common"
	"github.com/megaspacelab/megaconnect/connector"
	mgrpc "github.com/megaspacelab/megaconnect/grpc"
	"github.com/megaspacelab/megaconnect/workflow"

	"github.com/golang/protobuf/ptypes/empty"
	"github.com/google/uuid"
	"go.uber.org/zap"
	"google.golang.org/grpc"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

const (
	blockChanSize                    = 100
	blockCacheSize                   = blockChanSize
	leaseRenewalBuffer time.Duration = 5 * time.Second
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
	monitors          map[int64]*mgrpc.Monitor
	monitorsVersion   uint32
	blockCache        *ring.Ring

	// Concurrency control.
	running bool
	lock    sync.Mutex
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
		monitors:   make(map[int64]*mgrpc.Monitor),
		blockCache: ring.New(blockCacheSize),
	}
}

// Start would start an ChainManager loop
func (e *ChainManager) Start(listenPort int) error {
	e.lock.Lock()
	defer e.lock.Unlock()
	if e.running {
		return errors.New("event manage is already running")
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

	resp, err := e.orchClient.RegisterChainManager(e.ctx, &mgrpc.RegisterChainManagerRequest{
		ChainId:        e.connector.ChainName(),
		ChainManagerId: &mgrpc.InstanceId{Id: []byte(e.id), Instance: e.instance},
		ListenPort:     int32(listenPort),
		SessionId:      e.sessionID[:],
	})
	if err != nil {
		return err
	}

	e.logger.Debug("Registered with Orchestrator", zap.Stringer("lease", resp.Lease))
	e.updateLeaseWithLock(resp.Lease)

	err = e.connector.Start()
	if err != nil {
		return err
	}

	e.monitorsVersion = resp.Monitors.Version
	for _, m := range resp.Monitors.Monitors {
		e.monitors[m.Id] = m
	}

	resumeAfter, err := common.BytesToHash(resp.ResumeAfterBlockHash)
	if err != nil {
		return err
	}

	blocks := make(chan common.Block, blockChanSize)
	_, err = e.connector.SubscribeBlock(resumeAfter, blocks)
	if err != nil {
		return err
	}

	go func() {
		for block := range blocks {
			err := e.processNewBlock(block)
			if err != nil {
				e.logger.Fatal("Failed to process block", zap.Error(err))
				panic(err)
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
	return e.stopWithLock()
}

func (e *ChainManager) stopWithLock() error {
	if !e.running {
		return errors.New("event manager is not running")
	}
	e.running = false

	err := e.connector.Stop()
	if err != nil {
		e.logger.Error("connect stop with err", zap.Error(err))
	}

	e.leaseRenewalTimer.Stop()
	e.cancel()

	err = e.orchConn.Close()
	if err != nil {
		e.logger.Error("orchConn closed with err", zap.Error(err))
	}

	return nil
}

func (e *ChainManager) updateLeaseWithLock(lease *mgrpc.Lease) {
	e.leaseID = lease.Id
	timeout := time.Duration(lease.RemainingSeconds)*time.Second - leaseRenewalBuffer
	e.leaseRenewalTimer = time.AfterFunc(timeout, e.renewLease)
}

func (e *ChainManager) renewLease() {
	e.lock.Lock()
	defer e.lock.Unlock()
	if !e.running {
		e.logger.Warn("Skipping lease renewal on stopped ChainManager")
		return
	}

	resp, err := e.orchClient.RenewLease(e.ctx, &mgrpc.RenewLeaseRequest{
		LeaseId: e.leaseID,
	})
	if err != nil {
		e.logger.Fatal("Failed to renew lease. Shutting down")
		e.stopWithLock()
		panic("Failed to renew lease")
	}

	e.updateLeaseWithLock(resp)
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

	interpreter := workflow.New(workflow.NewEnv(e.connector, block))
	events := []*mgrpc.Event{}
	for _, monitor := range e.monitors {
		cond, err := workflow.NewByteDecoder(monitor.Condition).DecodeExpr()
		if err != nil {
			return err
		}
		e.logger.Debug("Evaluating", zap.Stringer("condition", cond))
		result, err := interpreter.EvalExpr(cond)
		if err != nil {
			return err
		}
		e.logger.Debug("Evaluated", zap.Stringer("condition", result))
		if result.Equal(workflow.GetBoolConst(false)) {
			continue
		}

		event := mgrpc.Event{}
		event.MonitorId = monitor.Id
		event.EvaluationsResults = make([][]byte, 0, len(monitor.Evaluations))
		for _, expr := range monitor.Evaluations {
			expr, err := workflow.NewByteDecoder(expr).DecodeExpr()
			if err != nil {
				return err
			}
			e.logger.Debug("Evaluating", zap.Stringer("expr", expr))
			result, err := interpreter.EvalExpr(expr)
			if err != nil {
				return err
			}
			e.logger.Debug("Evaluated", zap.Stringer("result", result))
			encoded, err := workflow.EncodeExpr(result)
			if err != nil {
				e.logger.Error("Failed to encode result", zap.Error(err))
				return err
			}
			event.EvaluationsResults = append(event.EvaluationsResults, encoded)
		}
		events = append(events, &event)
	}

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

func (e *ChainManager) findCachedBlockWithLock(hash *common.Hash) *ring.Ring {
	// Traverse in reverse order.
	for cache, i := e.blockCache, 0; cache.Value != nil && i < blockCacheSize; cache, i = cache.Prev(), i+1 {
		if cache.Value.(common.Block).Hash() == *hash {
			return cache
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

	resumeAfter, err := common.BytesToHash(preflight.ResumeAfterBlockHash)
	if err != nil {
		return status.Errorf(codes.InvalidArgument, "Invalid ResumeAfterBlockHash. %v", err)
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

	monitors := make(map[int64]*mgrpc.Monitor)

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
		monitors[monitor.Id] = monitor
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

	resumeAfter, err := common.BytesToHash(preflight.ResumeAfterBlockHash)
	if err != nil {
		return status.Errorf(codes.InvalidArgument, "Invalid ResumeAfterBlockHash. %v", err)
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
			e.monitors[m.AddMonitor.Monitor.Id] = m.AddMonitor.Monitor
		case *mgrpc.UpdateMonitorsRequest_RemoveMonitor_:
			e.logger.Debug("Received UpdateMonitorsRequest.RemoveMonitor", zap.Stringer("m", m.RemoveMonitor))
			delete(e.monitors, m.RemoveMonitor.MonitorId)
		default:
			return status.Error(codes.InvalidArgument, "Wrong message type. Expecting AddMonitor or RemoveMonitor")
		}
	}

	e.monitorsVersion = preflight.MonitorSetVersion

	err = stream.SendAndClose(&empty.Empty{})
	e.replayWithLock(resumeAfter)

	return err
}

func (e *ChainManager) replayWithLock(resumeAfter *common.Hash) {
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
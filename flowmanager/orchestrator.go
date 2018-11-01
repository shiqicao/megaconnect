package flowmanager

import (
	"bytes"
	"context"
	"errors"
	"io"
	"net"
	"reflect"
	"sync"
	"time"

	"github.com/golang/protobuf/ptypes/empty"
	"github.com/google/uuid"
	mgrpc "github.com/megaspacelab/megaconnect/grpc"
	"go.uber.org/zap"
	"google.golang.org/grpc"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/peer"
	"google.golang.org/grpc/status"
)

const (
	leaseDuration time.Duration = 60 * time.Second
	cmConnTimeout time.Duration = 5 * time.Second
)

// Orchestrator orchestrates ChainManagers, by telling them what to monitor and receiving reports from them.
// This type implements the mgrpc.OrchestratorServer interface.
type Orchestrator struct {
	// Static members.
	flowManager   *FlowManager
	log           *zap.Logger
	leaseDuration time.Duration

	// Dynamic members.
	leases       map[leaseID]*lease
	chainToLease map[string]*lease

	// Concurrency control.
	lock sync.Mutex
}

type leaseID uuid.UUID
type lease struct {
	id         leaseID
	expiration time.Time
	timer      *time.Timer

	chainManager *chainManagerProxy
}

type chainManagerProxy struct {
	// Static members.
	id        *mgrpc.InstanceId
	chainID   string
	sessionID []byte
	cmClient  mgrpc.ChainManagerClient
	conn      *grpc.ClientConn
	log       *zap.Logger

	// Dynamic members.
	config *ChainConfig

	// Signaling.
	done   chan struct{}
	closed bool
}

// NewOrchestrator creates a new Orchestrator.
func NewOrchestrator(fm *FlowManager, log *zap.Logger) *Orchestrator {
	return &Orchestrator{
		leases:        make(map[leaseID]*lease),
		chainToLease:  make(map[string]*lease),
		flowManager:   fm,
		log:           log,
		leaseDuration: leaseDuration,
	}
}

// Register registers this Orchestrator to the gRPC server.
func (o *Orchestrator) Register(server *grpc.Server) {
	mgrpc.RegisterOrchestratorServer(server, o)
}

// RegisterChainManager is invoked by ChainManagers to register themselves with this Orchestrator.
func (o *Orchestrator) RegisterChainManager(
	ctx context.Context,
	req *mgrpc.RegisterChainManagerRequest,
) (*mgrpc.RegisterChainManagerResponse, error) {
	o.log.Debug("Received RegisterChainManager request", zap.Stringer("req", req))

	if req.ChainId == "" {
		return nil, status.Error(codes.InvalidArgument, "Missing ChainId")
	}
	if req.ChainManagerId == nil {
		return nil, status.Error(codes.InvalidArgument, "Missing ChainManagerId")
	}
	if req.ListenPort == 0 {
		return nil, status.Error(codes.InvalidArgument, "Missing ListenPort")
	}

	peer, ok := peer.FromContext(ctx)
	if !ok {
		return nil, status.Error(codes.Unknown, "Failed to get peer info")
	}
	peerAddr, ok := peer.Addr.(*net.TCPAddr)
	if !ok {
		return nil, status.Errorf(codes.Unknown, "Unknown peer addr type %v", reflect.TypeOf(peer.Addr))
	}
	rpcAddr := *peerAddr
	rpcAddr.Port = int(req.ListenPort)

	// Front-load connection to minimize work inside lock.
	o.log.Debug("Connecting to ChainManager", zap.Stringer("rpcAddr", &rpcAddr))
	conn, err := grpc.DialContext(ctx, rpcAddr.String(),
		grpc.WithInsecure(), grpc.WithBlock(), grpc.WithTimeout(cmConnTimeout))
	if err != nil {
		return nil, err
	}
	cmClient := mgrpc.NewChainManagerClient(conn)

	chainConfig := o.flowManager.GetChainConfig(req.ChainId)
	if chainConfig == nil {
		conn.Close()
		return nil, status.Error(codes.InvalidArgument, "Invalid ChainId")
	}

	o.lock.Lock()
	defer o.lock.Unlock()

	lease := o.chainToLease[req.ChainId]
	if lease != nil && lease.valid() && !isNewerInstance(lease.chainManager.id, req.ChainManagerId) {
		conn.Close()
		return nil, status.Errorf(codes.FailedPrecondition,
			"Another live ChainManager already registered for %s", req.ChainId)
	}

	if lease != nil {
		o.expireLeaseWithLock(lease)
	}

	cm := &chainManagerProxy{
		id:        req.ChainManagerId,
		chainID:   req.ChainId,
		sessionID: req.SessionId,
		cmClient:  cmClient,
		conn:      conn,
		log:       o.log,
		config:    chainConfig,
		done:      make(chan struct{}),
	}
	lease = o.newLease(cm)

	o.leases[lease.id] = lease
	o.chainToLease[req.ChainId] = lease

	go o.chainManagerEventLoop(cm)

	o.log.Info("Registered new ChainManager",
		zap.Stringer("id", cm.id), zap.String("chainID", cm.chainID))

	return &mgrpc.RegisterChainManagerResponse{
		Lease:       lease.toRPCLease(),
		ResumeAfter: chainConfig.ResumeAfter,
		Monitors: &mgrpc.MonitorSet{
			Monitors: chainConfig.Monitors.Monitors(),
			Version:  chainConfig.MonitorsVersion,
		},
	}, nil
}

// UnregsiterChainManager is invoked by ChainManagers to register themselves with this Orchestrator.
func (o *Orchestrator) UnregsiterChainManager(
	ctx context.Context,
	req *mgrpc.UnregisterChainManagerRequest,
) (*empty.Empty, error) {
	o.log.Debug("Received UnregisterChainManager request", zap.Stringer("req", req))

	if req.ChainId == "" {
		return nil, status.Error(codes.InvalidArgument, "Missing ChainId")
	}
	if req.LeaseId == nil {
		return nil, status.Error(codes.InvalidArgument, "Missing LeaseId")
	}
	reqLeaseID, err := leaseIDFromBytes(req.LeaseId)
	if err != nil {
		return nil, status.Error(codes.InvalidArgument, "Invalid LeaseId")
	}
	o.lock.Lock()
	defer o.lock.Unlock()

	lease := o.chainToLease[req.ChainId]
	if lease == nil || lease.id != reqLeaseID {
		return nil, status.Error(codes.Aborted, "Lease doesn't match, doesn't exist or has already expired")
	}

	// expires lease and shutdown chain manager
	o.expireLeaseWithLock(lease)

	return &empty.Empty{}, nil
}

// RenewLease renews the lease between a ChainManager and this Orchestrator.
func (o *Orchestrator) RenewLease(
	ctx context.Context,
	req *mgrpc.RenewLeaseRequest,
) (*mgrpc.Lease, error) {
	o.log.Debug("Received RenewLease request", zap.Stringer("req", req))

	lid, err := leaseIDFromBytes(req.LeaseId)
	if err != nil {
		return nil, status.Error(codes.InvalidArgument, "Invalid LeaseId")
	}

	o.lock.Lock()
	defer o.lock.Unlock()

	lease := o.leases[lid]
	if lease == nil {
		return nil, status.Error(codes.Aborted, "Lease doesn't exist or has already expired")
	}
	cm := lease.chainManager
	if cm.closed {
		panic("ChainManager closed while lease is live")
	}

	lease.close(false)
	delete(o.leases, lease.id)

	lease = o.newLease(cm)
	o.leases[lease.id] = lease
	o.chainToLease[cm.chainID] = lease

	return lease.toRPCLease(), nil
}

// ReportBlockEvents reports a new block on the monitored chain, as well as all fired events from this block.
func (o *Orchestrator) ReportBlockEvents(stream mgrpc.Orchestrator_ReportBlockEventsServer) error {
	o.log.Debug("Received ReportBlockEvents request")

	msg, err := stream.Recv()
	if err == io.EOF {
		return status.Error(codes.InvalidArgument, "Missing Preflight")
	}
	if err != nil {
		return err
	}

	preflight := msg.GetPreflight()
	if preflight == nil {
		return status.Error(codes.InvalidArgument, "Wrong message type. Expecting Preflight")
	}
	o.log.Debug("Received ReportBlockEvents.Preflight", zap.Stringer("preflight", preflight))

	lid, err := leaseIDFromBytes(preflight.LeaseId)
	if err != nil {
		return status.Error(codes.InvalidArgument, "Invalid LeaseId")
	}

	lease := func() *lease {
		o.lock.Lock()
		defer o.lock.Unlock()
		return o.leases[lid]
	}()
	if lease == nil {
		return status.Error(codes.FailedPrecondition, "Lease doesn't exist or has already expired")
	}

	cm := lease.chainManager
	if preflight.MonitorSetVersion != cm.config.MonitorsVersion {
		return status.Error(codes.Aborted, "Wrong MonitorSetVersion")
	}

	msg, err = stream.Recv()
	if err == io.EOF {
		return status.Error(codes.InvalidArgument, "Missing Block")
	}
	if err != nil {
		return err
	}

	block := msg.GetBlock()
	if block == nil {
		return status.Error(codes.InvalidArgument, "Wrong message type. Expecting Block")
	}
	o.log.Debug("Received Block", zap.Stringer("block", block))

	var events []*mgrpc.Event

	for {
		msg, err = stream.Recv()
		if err == io.EOF {
			break
		}
		if err != nil {
			return err
		}

		event := msg.GetEvent()
		if event == nil {
			return status.Error(codes.InvalidArgument, "Wrong message type. Expecting Event")
		}
		o.log.Debug("Received Event", zap.Stringer("event", event))
		events = append(events, event)
	}

	o.flowManager.ReportBlockEvents(cm.chainID, cm.config.MonitorsVersion, block, events)
	return stream.SendAndClose(&empty.Empty{})
}

// isNewerInstance checks if inst2 is a newer instance of inst1.
func isNewerInstance(inst1 *mgrpc.InstanceId, inst2 *mgrpc.InstanceId) bool {
	return inst1 != nil && inst2 != nil && bytes.Equal(inst1.Id, inst2.Id) && inst1.Instance < inst2.Instance
}

func (o *Orchestrator) newLease(cm *chainManagerProxy) *lease {
	l := &lease{
		id:           leaseID(uuid.New()),
		expiration:   time.Now().Add(o.leaseDuration),
		chainManager: cm,
	}

	l.timer = time.AfterFunc(o.leaseDuration, func() {
		o.log.Warn("Lease expired", zap.Stringer("id", l.id), zap.Time("expiration", l.expiration))
		o.expireLease(l.id)
	})

	return l
}

func (o *Orchestrator) expireLease(lid leaseID) bool {
	o.lock.Lock()
	defer o.lock.Unlock()

	lease := o.leases[lid]
	if lease == nil {
		return false
	}
	if lease.id != lid {
		panic(errors.New("lease out of sync"))
	}

	o.expireLeaseWithLock(lease)
	return true
}

func (o *Orchestrator) expireLeaseForChainManager(cm *chainManagerProxy) bool {
	o.lock.Lock()
	defer o.lock.Unlock()

	lease := o.chainToLease[cm.chainID]
	if lease == nil || lease.chainManager != cm {
		return false
	}

	o.expireLeaseWithLock(lease)
	return true
}

func (o *Orchestrator) expireLeaseWithLock(lease *lease) {
	delete(o.leases, lease.id)
	delete(o.chainToLease, lease.chainManager.chainID)
	lease.close(true)
}

func leaseIDFromBytes(bs []byte) (leaseID, error) {
	id, err := uuid.FromBytes(bs)
	return leaseID(id), err
}

func (lid leaseID) String() string {
	return uuid.UUID(lid).String()
}

func (l *lease) valid() bool {
	return time.Now().Before(l.expiration)
}

func (l *lease) toRPCLease() *mgrpc.Lease {
	return &mgrpc.Lease{
		Id:               l.id[:],
		RemainingSeconds: uint32(time.Until(l.expiration).Seconds()),
	}
}

func (l *lease) close(closeCM bool) {
	l.timer.Stop()
	if closeCM {
		l.chainManager.close()
	}
}

func (o *Orchestrator) chainManagerEventLoop(cm *chainManagerProxy) {
	defer cm.conn.Close()

	for {
		select {
		case <-cm.done:
			return
		case <-cm.config.Outdated:
			o.log.Info("ChainConfig outdated", zap.String("chainID", cm.chainID))
			patch, err := o.flowManager.GetChainConfigPatch(cm.chainID, cm.config)
			if err != nil {
				o.log.Error("Failed to get ChainConfigPatch. Closing chainManager",
					zap.Error(err), zap.String("chainID", cm.chainID))
				o.expireLeaseForChainManager(cm)
				return
			}
			err = cm.patchConfig(patch)
			if err != nil {
				o.log.Error("Failed to update monitors. Closing chainManager",
					zap.Error(err), zap.String("chainID", cm.chainID))
				o.expireLeaseForChainManager(cm)
				return
			}
		}
	}
}

func (cm *chainManagerProxy) close() {
	if cm.closed {
		return
	}
	cm.closed = true
	close(cm.done)
}

func (cm *chainManagerProxy) patchConfig(patch *ChainConfigPatch) error {
	if cm.closed {
		return nil
	}

	oldVersion := cm.config.MonitorsVersion
	cm.config = patch.Apply(cm.config)

	ctx := context.Background()

	err := func() error {
		stream, err := cm.cmClient.UpdateMonitors(ctx)
		if err != nil {
			return err
		}

		err = stream.Send(&mgrpc.UpdateMonitorsRequest{
			MsgType: &mgrpc.UpdateMonitorsRequest_Preflight_{
				Preflight: &mgrpc.UpdateMonitorsRequest_Preflight{
					SessionId:                 cm.sessionID,
					PreviousMonitorSetVersion: oldVersion,
					MonitorSetVersion:         cm.config.MonitorsVersion,
					ResumeAfter:               cm.config.ResumeAfter,
				},
			},
		})
		if err != nil {
			return err
		}

		for _, m := range patch.AddMonitors {
			err = stream.Send(&mgrpc.UpdateMonitorsRequest{
				MsgType: &mgrpc.UpdateMonitorsRequest_AddMonitor_{
					AddMonitor: &mgrpc.UpdateMonitorsRequest_AddMonitor{
						Monitor: m,
					},
				},
			})
			if err != nil {
				return err
			}
		}

		for _, m := range patch.RemoveMonitors {
			err = stream.Send(&mgrpc.UpdateMonitorsRequest{
				MsgType: &mgrpc.UpdateMonitorsRequest_RemoveMonitor_{
					RemoveMonitor: &mgrpc.UpdateMonitorsRequest_RemoveMonitor{
						MonitorId: []byte(m),
					},
				},
			})
			if err != nil {
				return err
			}
		}

		_, err = stream.CloseAndRecv()
		return err
	}()

	if err == nil {
		return nil
	}

	s := status.Convert(err)
	if s.Code() != codes.Aborted {
		return err
	}

	// Fallback to sending full set.
	cm.log.Warn("Failed to patch monitors. Falling back to full reset", zap.String("chainID", cm.chainID))
	stream, err := cm.cmClient.SetMonitors(ctx)
	if err != nil {
		return err
	}

	err = stream.Send(&mgrpc.SetMonitorsRequest{
		MsgType: &mgrpc.SetMonitorsRequest_Preflight_{
			Preflight: &mgrpc.SetMonitorsRequest_Preflight{
				SessionId:         cm.sessionID,
				MonitorSetVersion: cm.config.MonitorsVersion,
				ResumeAfter:       cm.config.ResumeAfter,
			},
		},
	})
	if err != nil {
		return err
	}

	for _, m := range cm.config.Monitors {
		err = stream.Send(&mgrpc.SetMonitorsRequest{
			MsgType: &mgrpc.SetMonitorsRequest_Monitor{
				Monitor: m,
			},
		})
		if err != nil {
			return err
		}
	}

	_, err = stream.CloseAndRecv()
	return err
}

package flowmanager

import (
	"sync"

	"github.com/megaspacelab/megaconnect/grpc"

	"go.uber.org/zap"
)

// FlowManager manages megaflows and compiles them into monitors for each chain.
// It also aggregates reports from all chains. A disctinction from Orchestrator is that FlowManager deals with "chains",
// whereas Orchestrator deals with "chain managers".
type FlowManager struct {
	chainConfigs map[string]*ChainConfig

	lock sync.Mutex
	log  *zap.Logger
}

// NewFlowManager creates a new FlowManager.
func NewFlowManager(log *zap.Logger) *FlowManager {
	return &FlowManager{
		chainConfigs: make(map[string]*ChainConfig),
		log:          log,
	}
}

// GetChainConfig returns the ChainConfig for chainID.
func (fm *FlowManager) GetChainConfig(chainID string) *ChainConfig {
	fm.lock.Lock()
	defer fm.lock.Unlock()
	return fm.chainConfigs[chainID]
}

// SetChainConfig updates the ChainConfig for chainID.
func (fm *FlowManager) SetChainConfig(
	chainID string,
	monitors IndexedMonitors,
	resumeAfter *grpc.BlockSpec,
) *ChainConfig {
	fm.lock.Lock()
	defer fm.lock.Unlock()

	old := fm.chainConfigs[chainID]
	new := &ChainConfig{
		Monitors:    monitors,
		ResumeAfter: resumeAfter,
		Outdated:    make(chan struct{}),
	}
	fm.chainConfigs[chainID] = new

	if old != nil {
		new.MonitorsVersion = old.MonitorsVersion + 1
		close(old.Outdated)
	}
	return new
}

// ReportBlockEvents takes report of a new block and associated events.
func (fm *FlowManager) ReportBlockEvents(
	chainID string,
	monitorsVersion uint32,
	block *grpc.Block,
	events []*grpc.Event,
) {
	// TODO - implement. is locking needed?
	fm.log.Info("Received new block and events",
		zap.String("chainID", chainID),
		zap.Uint32("monitorsVersion", monitorsVersion),
		zap.Stringer("block", block),
		zap.Int("numEvents", len(events)),
	)
}

// ChainConfig captures configs about a chain.
type ChainConfig struct {
	// Monitors are the active monitors for this chain.
	Monitors IndexedMonitors

	// MonitorsVersion is used for Orchestrator and ChainManager to agree on the set of active monitors.
	MonitorsVersion uint32

	// ResumeAfter specifies the LAST block BEFORE this ChainConfig became effective.
	ResumeAfter *grpc.BlockSpec

	// Outdated is closed when this ChainConfig becomes outdated.
	Outdated chan struct{}
}

// TODO - what should it be and how is it generated?
type MonitorID int64

// IndexedMonitors is a bunch of monitors indexed by ID.
type IndexedMonitors map[MonitorID]*grpc.Monitor

// Monitors returns all monitors contained in this IndexedMonitors.
func (im IndexedMonitors) Monitors() []*grpc.Monitor {
	monitors := make([]*grpc.Monitor, 0, len(im))
	for _, m := range im {
		monitors = append(monitors, m)
	}
	return monitors
}

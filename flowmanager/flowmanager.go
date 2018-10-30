package flowmanager

import (
	"errors"
	"fmt"
	"sync"

	"github.com/megaspacelab/megaconnect/grpc"
	"github.com/megaspacelab/megaconnect/workflow"

	"go.uber.org/zap"
)

// FlowManager manages workflows and compiles them into monitors for each chain.
// It also aggregates reports from all chains. A disctinction from Orchestrator is that FlowManager deals with "chains",
// whereas Orchestrator deals with "chain managers".
type FlowManager struct {
	// States that go out to ChainManagers.
	chainConfigs map[string]*ChainConfig

	// Internal states.
	workflows    map[string]*workflow.WorkflowDecl
	actionsIndex map[string][]*workflow.ActionDecl // Actions indexed by monitor ids

	lock sync.Mutex
	log  *zap.Logger
}

// NewFlowManager creates a new FlowManager.
func NewFlowManager(log *zap.Logger) *FlowManager {
	return &FlowManager{
		chainConfigs: make(map[string]*ChainConfig),
		workflows:    make(map[string]*workflow.WorkflowDecl),
		actionsIndex: make(map[string][]*workflow.ActionDecl),
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

	return fm.setChainConfigWithLock(chainID, monitors, resumeAfter)
}

func (fm *FlowManager) setChainConfigWithLock(
	chainID string,
	monitors IndexedMonitors,
	resumeAfter *grpc.BlockSpec,
) *ChainConfig {
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
	fm.log.Info("Received new block and events",
		zap.String("chainID", chainID),
		zap.Uint32("monitorsVersion", monitorsVersion),
		zap.Stringer("block", block),
		zap.Int("numEvents", len(events)),
	)

	fm.lock.Lock()
	defer fm.lock.Unlock()

	if config := fm.chainConfigs[chainID]; config == nil || config.MonitorsVersion != monitorsVersion {
		fm.log.Info("Skipping stale block reports", zap.Uint32("expectedVersion", config.MonitorsVersion))
		return
	}

	for _, event := range events {
		fm.processEventWithLock(block, event)
	}
}

func (fm *FlowManager) processEventWithLock(
	block *grpc.Block,
	event *grpc.Event,
) {
	payload, err := workflow.NewByteDecoder(event.Payload).DecodeObjConst()
	if err != nil {
		fm.log.Error("Failed to decode event payload", zap.Error(err))
	}
	fm.log.Debug("Processing event", zap.ByteString("MonitorID", event.MonitorId), zap.Stringer("EventPayload", payload))
	for _, action := range fm.actionsIndex[string(event.MonitorId)] {
		fm.log.Debug("Evaluating action", zap.String("action", action.Name()))
		// TODO - evaluate
	}
}

// DeployWorkflow deploys a workflow.
// TODO - implement pending and activation. For now, this is immediately activated.
func (fm *FlowManager) DeployWorkflow(wf *workflow.WorkflowDecl) error {
	fm.lock.Lock()
	defer fm.lock.Unlock()

	id := WorkflowID(wf)
	if _, ok := fm.workflows[id]; ok {
		return errors.New("Workflow already exists")
	}
	fm.workflows[id] = wf

	// Deploy monitors.

	// This is a buffer for all affected chain configs.
	chainConfigs := make(map[string]*ChainConfig)
	eventMonitorIDs := make(map[string][]string)

	for _, m := range wf.MonitorDecls() {
		chain := m.Chain()
		if _, ok := fm.chainConfigs[chain]; !ok {
			return fmt.Errorf("Unsupported chain %s", chain)
		}

		cc := chainConfigs[chain]
		if cc == nil {
			// Copy the original set of monitors.
			cc = &ChainConfig{
				Monitors: make(map[string]*grpc.Monitor),
			}
			for k, v := range fm.chainConfigs[chain].Monitors {
				cc.Monitors[k] = v
			}
			chainConfigs[chain] = cc
		}

		mid := MonitorID(id, m)
		if _, ok := cc.Monitors[mid]; ok {
			panic("Duplicate monitor id")
		}

		monitor, err := Monitor(mid, m)
		if err != nil {
			return err
		}

		cc.Monitors[mid] = monitor

		eventName := m.EventName()
		eventMonitorIDs[eventName] = append(eventMonitorIDs[eventName], mid)
	}

	// Apply the new configs.
	for chain, config := range chainConfigs {
		fm.setChainConfigWithLock(chain, config.Monitors, config.ResumeAfter)
	}

	// Deploy actions.
	for _, a := range wf.ActionDecls() {
		for _, e := range a.TriggerEvents() {
			for _, mid := range eventMonitorIDs[e] {
				fm.actionsIndex[mid] = append(fm.actionsIndex[mid], a)
			}
		}
	}

	return nil
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

// IndexedMonitors is a bunch of monitors indexed by ID.
type IndexedMonitors map[string]*grpc.Monitor

// Monitors returns all monitors contained in this IndexedMonitors.
func (im IndexedMonitors) Monitors() []*grpc.Monitor {
	monitors := make([]*grpc.Monitor, 0, len(im))
	for _, m := range im {
		monitors = append(monitors, m)
	}
	return monitors
}

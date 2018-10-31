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
// The returned ChainConfig is a copy, so caller is free to mutate its non-pointer contents.
// In particular, the contents of the monitors map can be modified (not the pointer target though).
func (fm *FlowManager) GetChainConfig(chainID string) *ChainConfig {
	fm.lock.Lock()
	defer fm.lock.Unlock()
	return fm.chainConfigs[chainID].copy()
}

// GetChainConfigPatch returns a patch that can be applied to an existing ChainConfig.
// This is to avoid making a full copy with every config update.
func (fm *FlowManager) GetChainConfigPatch(chainID string, old *ChainConfig) (*ChainConfigPatch, error) {
	fm.lock.Lock()
	defer fm.lock.Unlock()

	cc := fm.chainConfigs[chainID]
	if cc == nil {
		return nil, fmt.Errorf("Non-existent chain %s", chainID)
	}

	patch := &ChainConfigPatch{
		chainConfig: cc.chainConfig,
	}

	var oldMonitors IndexedMonitors
	if old != nil {
		oldMonitors = old.Monitors
	}

	for k, v := range cc.Monitors {
		if _, ok := oldMonitors[k]; !ok {
			patch.AddMonitors = append(patch.AddMonitors, v)
		}
	}

	for k := range oldMonitors {
		if _, ok := cc.Monitors[k]; !ok {
			patch.RemoveMonitors = append(patch.RemoveMonitors, k)
		}
	}

	return patch, nil
}

// SetChainConfig updates the ChainConfig for chainID.
// Note that no copy is made for the passed parameters, so caller should relinquish ownership of them afterwards.
func (fm *FlowManager) SetChainConfig(
	chainID string,
	monitors IndexedMonitors,
	resumeAfter *grpc.BlockSpec,
) {
	fm.lock.Lock()
	defer fm.lock.Unlock()

	fm.setChainConfigWithLock(chainID, monitors, resumeAfter)
}

func (fm *FlowManager) setChainConfigWithLock(
	chainID string,
	monitors IndexedMonitors,
	resumeAfter *grpc.BlockSpec,
) {
	cc := fm.chainConfigs[chainID]
	if cc != nil {
		cc.MonitorsVersion++
		close(cc.Outdated)
	} else {
		cc = new(ChainConfig)
		fm.chainConfigs[chainID] = cc
	}

	cc.Monitors = monitors
	cc.ResumeAfter = resumeAfter
	cc.Outdated = make(chan struct{})
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
		fm.log.Error("Failed to decode event payload", zap.Error(err), zap.ByteString("payload", event.Payload))
		return
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

	chainMonitors := make(map[string]IndexedMonitors)
	eventMonitorIDs := make(map[string][]string)

	// Pre-validate all monitors before making any state changes.
	for _, m := range wf.MonitorDecls() {
		chain := m.Chain()
		if _, ok := fm.chainConfigs[chain]; !ok {
			return fmt.Errorf("Unsupported chain %s", chain)
		}

		monitors := chainMonitors[chain]
		if monitors == nil {
			monitors = make(IndexedMonitors)
			chainMonitors[chain] = monitors
		}

		mid := MonitorID(id, m)
		if _, ok := monitors[mid]; ok {
			panic("Duplicate monitor id")
		}

		monitor, err := Monitor(mid, m)
		if err != nil {
			return err
		}

		monitors[mid] = monitor

		eventName := m.EventName()
		eventMonitorIDs[eventName] = append(eventMonitorIDs[eventName], mid)
	}

	// Deploy monitors.
	for chain, ms := range chainMonitors {
		monitors := fm.chainConfigs[chain].Monitors
		if monitors == nil {
			monitors = make(IndexedMonitors)
		}

		for mid, monitor := range ms {
			if _, ok := monitors[mid]; ok {
				panic("Duplicate monitor id")
			}
			monitors[mid] = monitor
		}

		fm.setChainConfigWithLock(chain, monitors, nil) // TODO - resumeAfter
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

type chainConfig struct {
	// MonitorsVersion is used for Orchestrator and ChainManager to agree on the set of active monitors.
	MonitorsVersion uint32

	// ResumeAfter specifies the LAST block BEFORE this ChainConfig became effective.
	ResumeAfter *grpc.BlockSpec

	// Outdated is closed when this ChainConfig becomes outdated.
	Outdated chan struct{}
}

// ChainConfig captures configs about a chain.
type ChainConfig struct {
	chainConfig

	// Monitors are the active monitors for this chain.
	Monitors IndexedMonitors
}

func (cc *ChainConfig) copy() *ChainConfig {
	if cc == nil {
		return nil
	}
	cp := *cc
	cp.Monitors = cp.Monitors.Copy()
	return &cp
}

// ChainConfigPatch represents a patch that can be applied to a ChainConfig.
type ChainConfigPatch struct {
	chainConfig

	AddMonitors    []*grpc.Monitor
	RemoveMonitors []string
}

// Apply applies this patch to cc.
// cc will be patched in place, unless it's nil, in which case a new ChainConfig is created.
func (ccp *ChainConfigPatch) Apply(cc *ChainConfig) *ChainConfig {
	if ccp == nil {
		return cc
	}

	if cc == nil {
		cc = &ChainConfig{
			chainConfig: ccp.chainConfig,
		}
	} else {
		cc.chainConfig = ccp.chainConfig
	}

	for _, m := range ccp.RemoveMonitors {
		delete(cc.Monitors, m)
	}

	if cc.Monitors == nil && len(ccp.AddMonitors) > 0 {
		cc.Monitors = make(IndexedMonitors, len(ccp.AddMonitors))
	}

	for _, m := range ccp.AddMonitors {
		cc.Monitors[string(m.Id)] = m
	}

	return cc
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

// Copy makes a shallow copy of the IndexedMonitors.
func (im IndexedMonitors) Copy() IndexedMonitors {
	if im == nil {
		return nil
	}

	cp := make(IndexedMonitors, len(im))
	for k, v := range im {
		cp[k] = v
	}
	return cp
}

// Copyright 2018 @ MegaSpace

// The MegaSpace source code is distributed in the hope that it will be useful,
// but WITHOUT ANY WARRANTY; without even the implied warranty of
// MERCHANTABILITY or FITNESS FOR A PARTICULAR PURPOSE. See the
// GNU Lesser General Public License for more details.

// You should have received a copy of the GNU Lesser General Public License
// along with the MegaSpace source code. If not, see <http://www.gnu.org/licenses/>.

package flowmanager

import (
	"errors"
	"fmt"
	"sync"

	"github.com/megaspacelab/megaconnect/grpc"
	"github.com/megaspacelab/megaconnect/unsafe"
	"github.com/megaspacelab/megaconnect/workflow"

	"go.uber.org/zap"
)

var (
	dummyResolver = workflow.NewResolver(nil, nil)
)

// FlowManager manages workflows and compiles them into monitors for each chain.
// It also aggregates reports from all chains. A disctinction from Orchestrator is that FlowManager deals with "chains",
// whereas Orchestrator deals with "chain managers".
type FlowManager struct {
	// States that go out to ChainManagers.
	chainConfigs map[string]*ChainConfig

	// Internal states.
	workflows    map[string]*workflow.WorkflowDecl
	monitorInfos map[string]*monitorInfo
	actionsIndex map[string][]*actionState // Actions indexed by event ids

	lock sync.Mutex
	log  *zap.Logger
}

// NewFlowManager creates a new FlowManager.
func NewFlowManager(log *zap.Logger) *FlowManager {
	return &FlowManager{
		chainConfigs: make(map[string]*ChainConfig),
		workflows:    make(map[string]*workflow.WorkflowDecl),
		monitorInfos: make(map[string]*monitorInfo),
		actionsIndex: make(map[string][]*actionState),
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
		chainConfigBase: cc.chainConfigBase,
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
		info := fm.monitorInfos[unsafe.BytesToString(event.MonitorId)]
		if info == nil {
			fm.log.Error("Missing monitor info", zap.Binary("monitorID", event.MonitorId))
			continue
		}

		payload, err := workflow.NewByteDecoder(event.Payload).DecodeObjConst()
		if err != nil {
			fm.log.Error("Failed to decode event payload", zap.Error(err), zap.Binary("payload", event.Payload))
			continue
		}

		fm.processEventWithLock(info.workflowID, info.eventName, payload)
	}
}

func (fm *FlowManager) processEventWithLock(
	workflowID string,
	eventName string,
	eventPayload *workflow.ObjConst,
) {
	fm.log.Debug("Processing event",
		zap.String("workflowID", workflowID),
		zap.String("eventName", eventName),
		zap.String("eventPayload", workflow.PrintNode(eventPayload)),
	)

	var fires []*workflow.FireEventResult
	eventID := EventID(workflowID, eventName)

	for _, action := range fm.actionsIndex[eventID] {
		fm.log.Debug("Evaluating action", zap.Stringer("action", action.decl.Name()))
		action.observe(eventName, eventPayload)

		env := workflow.NewEnv(action)
		intp := workflow.NewInterpreter(env, nil, dummyResolver, fm.log)

		results, err := intp.EvalAction(action.decl)
		if err != nil {
			fm.log.Error("Failed to evaluate action",
				zap.Error(err), zap.Stringer("action", action.decl.Name()))
			continue
		}
		if results == nil {
			continue
		}

		fm.log.Debug("Action triggered", zap.Stringer("action", action.decl.Name()))

		// Reset action state upon triggering
		action.reset()

		for _, r := range results {
			switch result := r.(type) {
			case *workflow.FireEventResult:
				fires = append(fires, result)
			}
		}
	}

	// Recursively process fired events.
	// This means all fired events are inserted into the event stream right after their triggering event.
	for _, fire := range fires {
		fm.processEventWithLock(workflowID, fire.EventName(), fire.Payload())
	}
}

// DeployWorkflow deploys a workflow.
// TODO - implement pending and activation. For now, this is immediately activated.
func (fm *FlowManager) DeployWorkflow(wf *workflow.WorkflowDecl) error {
	fm.lock.Lock()
	defer fm.lock.Unlock()

	wfid := WorkflowID(wf)
	if _, ok := fm.workflows[wfid]; ok {
		return errors.New("Workflow already exists")
	}
	fm.workflows[wfid] = wf

	chainMonitors := make(map[string]IndexedMonitors)
	monitorInfos := make(map[string]*monitorInfo)

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

		mid := MonitorID(wfid, m)
		if _, ok := monitors[mid]; ok {
			panic("Duplicate monitor id")
		}

		monitor, err := Monitor(mid, m)
		if err != nil {
			return err
		}

		monitors[mid] = monitor
		monitorInfos[mid] = &monitorInfo{
			eventName:  m.EventName(),
			workflowID: wfid,
		}
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
	for k, v := range monitorInfos {
		if _, ok := fm.monitorInfos[k]; ok {
			panic("Duplicate monitor id")
		}
		fm.monitorInfos[k] = v
	}

	// Deploy actions.
	for _, a := range wf.ActionDecls() {
		action := &actionState{decl: a}

		for _, e := range a.TriggerEvents() {
			eid := EventID(wfid, e)
			fm.actionsIndex[eid] = append(fm.actionsIndex[eid], action)
		}
	}

	return nil
}

type monitorInfo struct {
	workflowID string
	eventName  string
}

type actionState struct {
	decl           *workflow.ActionDecl
	observedEvents map[string]*workflow.ObjConst
}

func (s *actionState) Lookup(event string) *workflow.ObjConst {
	return s.observedEvents[event]
}

func (s *actionState) observe(event string, props *workflow.ObjConst) {
	if s.observedEvents == nil {
		s.observedEvents = make(map[string]*workflow.ObjConst)
	}
	s.observedEvents[event] = props
}

func (s *actionState) reset() {
	s.observedEvents = nil
}

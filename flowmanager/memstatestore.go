package flowmanager

import (
	"errors"

	"github.com/megaspacelab/megaconnect/grpc"
	"github.com/megaspacelab/megaconnect/protos"
	"github.com/megaspacelab/megaconnect/workflow"

	"github.com/golang/protobuf/proto"
)

var (
	ErrFinalized = errors.New("Pending state already finalized")
)

type eventID struct {
	wfid      WorkflowID
	eventName string
}

type actionID struct {
	wfid       WorkflowID
	actionName string
}

type eventStates map[string]*workflow.ObjConst

type memState struct {
	workflowByID        map[WorkflowID]*workflow.WorkflowDecl
	lastBlockByChain    map[string]*protos.Block
	monitorByID         map[grpc.MonitorID]*grpc.Monitor
	monitorsByChain     map[string][]*grpc.MonitorID
	actionByID          map[actionID]*workflow.ActionDecl
	eventStatesByAction map[actionID]eventStates
	actionsByEvent      map[eventID][]*actionID
}

func newMemState() memState {
	return memState{
		workflowByID:        make(map[WorkflowID]*workflow.WorkflowDecl),
		lastBlockByChain:    make(map[string]*protos.Block),
		monitorByID:         make(map[grpc.MonitorID]*grpc.Monitor),
		monitorsByChain:     make(map[string][]*grpc.MonitorID),
		actionByID:          make(map[actionID]*workflow.ActionDecl),
		eventStatesByAction: make(map[actionID]eventStates),
		actionsByEvent:      make(map[eventID][]*actionID),
	}
}

type committedState struct {
	memState
	mblocks []*protos.MBlock
}

type pendingState struct {
	memState
	mblock    *protos.MBlock
	finalized bool
}

// MemStateStore is an in-memory implementation of StateStore.
type MemStateStore struct {
	committed committedState
	pending   pendingState
}

func NewMemStateStore() *MemStateStore {
	return &MemStateStore{
		committed: committedState{
			memState: newMemState(),
		},
		pending: pendingState{
			memState: newMemState(),
		},
	}
}

func (s *MemStateStore) LastReportedBlockByChain(chain string) (*protos.Block, error) {
	if block, ok := s.pending.lastBlockByChain[chain]; ok {
		return block, nil
	}
	return s.committed.lastBlockByChain[chain], nil
}

func (s *MemStateStore) WorkflowByID(id WorkflowID) (*workflow.WorkflowDecl, error) {
	if wf, ok := s.pending.workflowByID[id]; ok {
		return wf, nil
	}
	return s.committed.workflowByID[id], nil
}

func (s *MemStateStore) monitorByID(id *grpc.MonitorID) *grpc.Monitor {
	if monitor, ok := s.pending.monitorByID[*id]; ok {
		return monitor
	}
	return s.committed.monitorByID[*id]
}

func (s *MemStateStore) MonitorsByChain(chain string) ([]*grpc.Monitor, error) {
	var monitors []*grpc.Monitor

	for _, id := range s.committed.monitorsByChain[chain] {
		if monitor := s.monitorByID(id); monitor != nil {
			monitors = append(monitors, monitor)
		}
	}

	for _, id := range s.pending.monitorsByChain[chain] {
		if monitor := s.monitorByID(id); monitor != nil {
			monitors = append(monitors, monitor)
		}
	}

	return monitors, nil
}

func copyEventStates(states eventStates) eventStates {
	// Always return a non-nil map, even if empty.
	cp := make(eventStates, len(states))
	for k, v := range states {
		cp[k] = v
	}

	return cp
}

func (s *MemStateStore) actionByID(aid *actionID) *workflow.ActionDecl {
	if action, ok := s.pending.actionByID[*aid]; ok {
		return action
	}
	return s.committed.actionByID[*aid]
}

func (s *MemStateStore) actionsByEvent(eid *eventID) []*actionID {
	if aids, ok := s.pending.actionsByEvent[*eid]; ok {
		return aids
	}
	return s.committed.actionsByEvent[*eid]
}

func (s *MemStateStore) eventStatesByAction(aid *actionID) eventStates {
	if states, ok := s.pending.eventStatesByAction[*aid]; ok {
		return states
	}
	return s.committed.eventStatesByAction[*aid]
}

func (s *MemStateStore) ActionStatesByEvent(wfid WorkflowID, eventName string) ([]*ActionState, error) {
	aids := s.actionsByEvent(&eventID{wfid, eventName})
	if aids == nil {
		return nil, nil
	}

	var actionStates []*ActionState

	for _, aid := range aids {
		if action := s.actionByID(aid); action != nil {
			// Always return a copy of the states because the caller may mutate map contents.
			actionStates = append(actionStates, &ActionState{
				Action: action,
				State:  copyEventStates(s.eventStatesByAction(aid)),
			})
		}
	}

	return actionStates, nil
}

func (s *MemStateStore) MBlockByHeight(height uint64) (*protos.MBlock, error) {
	if height >= uint64(len(s.committed.mblocks)) {
		return nil, errors.New("Height out of range")
	}

	return s.committed.mblocks[height], nil
}

func (s *MemStateStore) lastMBlock() *protos.MBlock {
	n := len(s.committed.mblocks)
	if n == 0 {
		return nil
	}
	return s.committed.mblocks[n-1]
}

func (s *MemStateStore) pendingMBlock() *protos.MBlock {
	if s.pending.mblock == nil {
		parent := s.lastMBlock()
		if parent == nil {
			// Genisis
			s.pending.mblock = &protos.MBlock{}
		} else {
			s.pending.mblock = &protos.MBlock{
				Height:     parent.Height + 1,
				ParentHash: parent.Hash,
			}
		}
	}
	return s.pending.mblock
}

func (s *MemStateStore) PutBlockReport(block *protos.Block) error {
	if s.pending.finalized {
		return ErrFinalized
	}

	mblock := s.pendingMBlock()
	mblock.Events = append(mblock.Events, &protos.MEvent{
		Body: &protos.MEvent_ReportBlock{
			ReportBlock: proto.Clone(block).(*protos.Block),
		},
	})

	s.pending.lastBlockByChain[block.Chain] = block
	return nil
}

func (s *MemStateStore) PutMonitorEvent(event *protos.MonitorEvent) error {
	if s.pending.finalized {
		return ErrFinalized
	}

	mblock := s.pendingMBlock()
	mblock.Events = append(mblock.Events, &protos.MEvent{
		Body: &protos.MEvent_MonitorEvent{
			MonitorEvent: proto.Clone(event).(*protos.MonitorEvent),
		},
	})
	return nil
}

func (s *MemStateStore) PutActionEvent(event *protos.ActionEvent) error {
	if s.pending.finalized {
		return ErrFinalized
	}

	mblock := s.pendingMBlock()
	mblock.Events = append(mblock.Events, &protos.MEvent{
		Body: &protos.MEvent_ActionEvent{
			ActionEvent: proto.Clone(event).(*protos.ActionEvent),
		},
	})
	return nil
}

func (s *MemStateStore) PutActionEventState(
	wfid WorkflowID,
	actionName string,
	eventName string,
	state *workflow.ObjConst,
) error {
	if s.pending.finalized {
		return ErrFinalized
	}

	aid := actionID{wfid, actionName}
	states := s.pending.eventStatesByAction[aid]
	if states == nil {
		states = copyEventStates(s.committed.eventStatesByAction[aid])
		s.pending.eventStatesByAction[aid] = states
	}
	states[eventName] = state
	return nil
}

func (s *MemStateStore) ResetActionState(wfid WorkflowID, actionName string) error {
	if s.pending.finalized {
		return ErrFinalized
	}

	s.pending.eventStatesByAction[actionID{wfid, actionName}] = nil
	return nil
}

func (s *MemStateStore) PutWorkflow(id WorkflowID, wf *workflow.WorkflowDecl) error {
	if s.pending.finalized {
		return ErrFinalized
	}

	wfPayload, err := workflow.EncodeWorkflow(wf)
	if err != nil {
		return err
	}

	mblock := s.pendingMBlock()
	mblock.Events = append(mblock.Events, &protos.MEvent{
		Body: &protos.MEvent_DeployWorkflow{
			DeployWorkflow: &protos.DeployWorkflow{
				WorkflowId: id.Bytes(),
				Payload:    wfPayload,
			},
		},
	})

	s.pending.workflowByID[id] = wf

ForMonitor:
	for _, m := range wf.MonitorDecls() {
		bs, err := workflow.EncodeMonitorDecl(m)
		if err != nil {
			return err
		}

		mid := &grpc.MonitorID{id, m.Name().Id()}

		s.pending.monitorByID[*mid] = &grpc.Monitor{
			WorkflowId:  id.Bytes(),
			MonitorName: m.Name().Id(),
			Monitor:     bs,
		}

		mids, ok := s.pending.monitorsByChain[m.Chain()]
		if !ok {
			cmids := s.committed.monitorsByChain[m.Chain()]
			mids = cmids[:len(cmids):len(cmids)] // This will trigger a copy-on-append
		}
		for _, v := range mids {
			if *v == *mid {
				// Skip adding if already existing.
				continue ForMonitor
			}
		}
		s.pending.monitorsByChain[m.Chain()] = append(mids, mid)
	}

ForAction:
	for _, action := range wf.ActionDecls() {
		aid := &actionID{id, action.Name().Id()}

		s.pending.actionByID[*aid] = action
		s.pending.eventStatesByAction[*aid] = nil

		for _, event := range action.TriggerEvents() {
			eid := &eventID{id, event}

			aids, ok := s.pending.actionsByEvent[*eid]
			if !ok {
				caids := s.committed.actionsByEvent[*eid]
				aids = caids[:len(caids):len(caids)] // This will trigger a copy-on-append
			}
			for _, v := range aids {
				if *v == *aid {
					// Skip adding if already existing.
					continue ForAction
				}
			}
			s.pending.actionsByEvent[*eid] = append(aids, aid)
		}
	}

	return nil
}

func (s *MemStateStore) DeleteWorkflow(id WorkflowID) error {
	if s.pending.finalized {
		return ErrFinalized
	}

	wf, err := s.WorkflowByID(id)
	if err != nil {
		return err
	}
	if wf == nil {
		return nil
	}

	mblock := s.pendingMBlock()
	mblock.Events = append(mblock.Events, &protos.MEvent{
		Body: &protos.MEvent_UndeployWorkflow{
			UndeployWorkflow: &protos.UndeployWorkflow{
				WorkflowId: id.Bytes(),
			},
		},
	})

	s.pending.workflowByID[id] = nil

	for _, m := range wf.MonitorDecls() {
		mid := &grpc.MonitorID{id, m.Name().Id()}

		s.pending.monitorByID[*mid] = nil

		mids, ok := s.pending.monitorsByChain[m.Chain()]
		if !ok {
			cmids := s.committed.monitorsByChain[m.Chain()]
			mids = append(cmids[:0:0], cmids...) // Copy
		}

		// Delete mid from mids.
		filtered := mids[:0]
		for _, v := range mids {
			if *v != *mid {
				filtered = append(filtered, v)
			}
		}

		s.pending.monitorsByChain[m.Chain()] = filtered
	}

	for _, event := range wf.EventDecls() {
		s.pending.actionsByEvent[eventID{id, event.Name().Id()}] = nil
	}
	for _, action := range wf.ActionDecls() {
		aid := &actionID{id, action.Name().Id()}

		s.pending.actionByID[*aid] = nil
		s.pending.eventStatesByAction[*aid] = nil

		for _, event := range action.TriggerEvents() {
			eid := &eventID{id, event}

			aids, ok := s.pending.actionsByEvent[*eid]
			if !ok {
				caids := s.committed.actionsByEvent[*eid]
				aids = append(caids[:0:0], caids...) // Copy
			}

			// Delete aid from aids.
			filtered := aids[:0]
			for _, v := range aids {
				if *v != *aid {
					filtered = append(filtered, v)
				}
			}

			s.pending.actionsByEvent[*eid] = filtered
		}
	}

	return nil
}

func (s *MemStateStore) FinalizeMBlock() (*protos.MBlock, error) {
	block := s.pendingMBlock()
	s.pending.finalized = true
	return block, nil
}

func (s *MemStateStore) CommitPendingMBlock(block *protos.MBlock) error {
	if block != s.pendingMBlock() {
		return errors.New("Wrong mblock to commit")
	}

	for k, v := range s.pending.workflowByID {
		if v == nil {
			delete(s.committed.workflowByID, k)
		} else {
			s.committed.workflowByID[k] = v
		}
	}

	for k, v := range s.pending.lastBlockByChain {
		if v == nil {
			delete(s.committed.lastBlockByChain, k)
		} else {
			s.committed.lastBlockByChain[k] = v
		}
	}

	for k, v := range s.pending.monitorByID {
		if v == nil {
			delete(s.committed.monitorByID, k)
		} else {
			s.committed.monitorByID[k] = v
		}
	}

	for k, v := range s.pending.monitorsByChain {
		if v == nil {
			delete(s.committed.monitorsByChain, k)
		} else {
			s.committed.monitorsByChain[k] = v
		}
	}

	for k, v := range s.pending.actionByID {
		if v == nil {
			delete(s.committed.actionByID, k)
		} else {
			s.committed.actionByID[k] = v
		}
	}

	for k, v := range s.pending.eventStatesByAction {
		if v == nil {
			delete(s.committed.eventStatesByAction, k)
		} else {
			s.committed.eventStatesByAction[k] = v
		}
	}

	for k, v := range s.pending.actionsByEvent {
		if v == nil {
			delete(s.committed.actionsByEvent, k)
		} else {
			s.committed.actionsByEvent[k] = v
		}
	}

	s.committed.mblocks = append(s.committed.mblocks, block)

	return s.ResetPending()
}

func (s *MemStateStore) ResetPending() error {
	s.pending = pendingState{memState: newMemState()}
	return nil
}

package flowmanager

import (
	"github.com/megaspacelab/megaconnect/grpc"
	"github.com/megaspacelab/megaconnect/protos"
	"github.com/megaspacelab/megaconnect/workflow"
)

// TODO - unify/cleanup proto types with decl types.

// StateStore is the data access layer of FlowManager.
type StateStore interface {
	// Read operations.
	LastReportedBlockByChain(chain string) (*protos.Block, error)
	WorkflowByID(id WorkflowID) (*workflow.WorkflowDecl, error)
	MonitorsByChain(chain string) ([]*grpc.Monitor, error)
	ActionStatesByEvent(wfid WorkflowID, eventName string) ([]*ActionState, error)
	MBlockByHeight(height uint64) (*protos.MBlock, error)

	// Write operations.
	PutBlockReport(block *protos.Block) error
	PutMonitorEvent(event *protos.MonitorEvent) error
	PutActionEvent(event *protos.ActionEvent) error
	PutActionEventState(wfid WorkflowID, actionName string, eventName string, state *workflow.ObjConst) error
	ResetActionState(wfid WorkflowID, actionName string) error
	PutWorkflow(id WorkflowID, wf *workflow.WorkflowDecl) error
	DeleteWorkflow(id WorkflowID) error
	FinalizeMBlock() (*protos.MBlock, error)
	CommitPendingMBlock(block *protos.MBlock) error
	ResetPending() error
}

// ActionState captures the state of an action.
type ActionState struct {
	Action *workflow.ActionDecl
	State  map[string]*workflow.ObjConst
}

// Lookup looks up action state for the event.
func (s *ActionState) Lookup(event string) *workflow.ObjConst {
	return s.State[event]
}

// Observe reports a new event to this action.
// It returns whether the observation triggered a state change.
func (s *ActionState) Observe(event string, payload *workflow.ObjConst) bool {
	// This implies last touch attribution model.
	s.State[event] = payload
	return true
}

// Copyright 2018 @ MegaSpace

// The MegaSpace source code is distributed in the hope that it will be useful,
// but WITHOUT ANY WARRANTY; without even the implied warranty of
// MERCHANTABILITY or FITNESS FOR A PARTICULAR PURPOSE. See the
// GNU Lesser General Public License for more details.

// You should have received a copy of the GNU Lesser General Public License
// along with the MegaSpace source code. If not, see <http://www.gnu.org/licenses/>.

package flowmanager

import (
	"context"
	"fmt"
	"sort"
	"sync"

	"github.com/golang/protobuf/ptypes/empty"
	"github.com/megaspacelab/megaconnect/common"
	mgrpc "github.com/megaspacelab/megaconnect/grpc"
	"github.com/megaspacelab/megaconnect/protos"
	"github.com/megaspacelab/megaconnect/unsafe"
	"github.com/megaspacelab/megaconnect/workflow"
	"go.uber.org/zap"
	"google.golang.org/grpc"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

const (
	chainConfigSubBufferSize = 3
)

var (
	dummyResolver = workflow.NewResolver(nil, nil)
)

// FlowManager manages workflows and compiles them into monitors for each chain.
// It also aggregates reports from all chains. A disctinction from Orchestrator is that FlowManager deals with "chains",
// whereas Orchestrator deals with "chain managers".
type FlowManager struct {
	// In-memory states for chain communications.
	chainConfigVersions map[string]uint32
	chainConfigSubs     map[string][]*ChainConfigSub

	// Workflow changes don't enter state store right away, because we don't want them to take effect immediately.
	// This is a temp area to hold them. Non-nil values imply deploy. Nil values imply undeploy.
	pendingWorkflows map[WorkflowID]*workflow.WorkflowDecl

	// All persistent states go through stateStore.
	stateStore StateStore

	lock sync.Mutex
	log  *zap.Logger
}

// NewFlowManager creates a new FlowManager.
func NewFlowManager(stateStore StateStore, log *zap.Logger) *FlowManager {
	return &FlowManager{
		chainConfigVersions: make(map[string]uint32),
		chainConfigSubs:     make(map[string][]*ChainConfigSub),
		pendingWorkflows:    make(map[WorkflowID]*workflow.WorkflowDecl),
		stateStore:          stateStore,
		log:                 log,
	}
}

// Register registers this FlowManager to the gRPC server.
func (fm *FlowManager) Register(server *grpc.Server) {
	mgrpc.RegisterWorkflowManagerServer(server, fm)
}

// DeployWorkflow is invoked to deploy a workflow
func (fm *FlowManager) DeployWorkflow(
	ctx context.Context,
	req *mgrpc.DeployWorkflowRequest,
) (*mgrpc.DeployWorkflowResponse, error) {
	fm.log.Debug("Received DeployWorkflow request")

	if req.Payload == nil {
		return nil, status.Error(codes.InvalidArgument, "Missing Payload")
	}

	workflow, err := workflow.NewByteDecoder(req.Payload).DecodeWorkflow()
	if err != nil {
		return nil, err
	}

	err = fm.deployWorkflow(workflow)
	if err != nil {
		return nil, err
	}

	id := GetWorkflowID(workflow).Bytes()

	resp := &mgrpc.DeployWorkflowResponse{
		WorkflowId: id,
	}
	fm.log.Debug("DeployWorkflow successful", zap.ByteString("Workflow ID", id))

	return resp, nil
}

// UndeployWorkflow is invoked to undeploy a workflow
func (fm *FlowManager) UndeployWorkflow(
	ctx context.Context,
	req *mgrpc.UndeployWorkflowRequest,
) (*empty.Empty, error) {
	fm.log.Debug("Received UndeployWorkflow request", zap.ByteString("Workflow ID", req.WorkflowId))

	if req.WorkflowId == nil {
		return nil, status.Error(codes.InvalidArgument, "Missing Workflow ID")
	}

	id := common.ImmutableBytes(req.WorkflowId)

	err := fm.undeployWorkflow(id)
	if err != nil {
		return nil, err
	}
	fm.log.Debug("UndeployWorkflow successful", zap.ByteString("Workflow ID", req.WorkflowId))

	return &empty.Empty{}, nil
}

// GetChainConfig returns the current ChainConfig and a subscription for future patches for chainID.
func (fm *FlowManager) GetChainConfig(chainID string) (*ChainConfig, *ChainConfigSub, error) {
	fm.lock.Lock()
	defer fm.lock.Unlock()

	monitors, err := fm.stateStore.MonitorsByChain(chainID)
	if err != nil {
		return nil, nil, err
	}

	var resumeAfter *mgrpc.BlockSpec
	lastBlock, err := fm.stateStore.LastReportedBlockByChain(chainID)
	if err != nil {
		return nil, nil, err
	}
	if lastBlock != nil {
		resumeAfter = &mgrpc.BlockSpec{
			Hash:   lastBlock.Hash,
			Height: lastBlock.Height,
		}
	}

	version, ok := fm.chainConfigVersions[chainID]
	if !ok {
		version = 1
		fm.chainConfigVersions[chainID] = version
	}

	cc := &ChainConfig{
		Monitors: monitors,
		chainConfigBase: chainConfigBase{
			MonitorsVersion: version,
			ResumeAfter:     resumeAfter,
		},
	}

	sub := &ChainConfigSub{
		patches: make(chan *ChainConfigPatch, chainConfigSubBufferSize),
		cancel:  make(chan common.Nothing),
	}

	// TODO - active GC for subs
	fm.chainConfigSubs[chainID] = append(fm.chainConfigSubs[chainID], sub)

	return cc, sub, nil
}

// ReportBlockEvents takes report of a new block and associated events.
func (fm *FlowManager) ReportBlockEvents(
	monitorsVersion uint32,
	block *protos.Block,
	events []*protos.MonitorEvent,
) {
	fm.log.Info("Received new block and events",
		zap.String("chainID", block.Chain),
		zap.Uint32("monitorsVersion", monitorsVersion),
		zap.Stringer("block", block),
		zap.Int("numEvents", len(events)),
	)

	fm.lock.Lock()
	defer fm.lock.Unlock()

	if version, ok := fm.chainConfigVersions[block.Chain]; !ok || version != monitorsVersion {
		fm.log.Info("Skipping stale block reports", zap.Uint32("expectedVersion", version))
		return
	}

	err := fm.stateStore.PutBlockReport(block)
	if err != nil {
		fm.log.Error("Failed to put block report", zap.Error(err))
		return
	}

	for _, event := range events {
		payload, err := workflow.NewByteDecoder(event.Payload).DecodeObjConst()
		if err != nil {
			fm.log.Error("Failed to decode event payload", zap.Error(err), zap.Binary("payload", event.Payload))
			continue
		}

		err = fm.stateStore.PutMonitorEvent(event)
		if err != nil {
			fm.log.Error("Failed to put monitor event", zap.Error(err))
			continue
		}

		fm.processEventWithLock(WorkflowID(unsafe.BytesToString(event.WorkflowId)), event.EventName, payload)
	}
}

func (fm *FlowManager) processEventWithLock(
	workflowID WorkflowID,
	eventName string,
	eventPayload *workflow.ObjConst,
) {
	fm.log.Debug("Processing event",
		zap.Stringer("workflowID", workflowID),
		zap.String("eventName", eventName),
		zap.String("eventPayload", workflow.PrintNode(eventPayload)),
	)

	var fires []*workflow.FireEventResult
	actions, err := fm.stateStore.ActionStatesByEvent(workflowID, eventName)
	if err != nil {
		fm.log.Error("Failed to retrieve action states for event",
			zap.Error(err), zap.Stringer("workflowID", workflowID), zap.String("event", eventName))
		return
	}

	for _, action := range actions {
		actionName := action.Action.Name().Id()
		fm.log.Debug("Evaluating action", zap.String("action", actionName))

		if !action.Observe(eventName, eventPayload) {
			fm.log.Debug("No state change. Skipping evaluation", zap.String("action", actionName))
			continue
		}

		err = fm.stateStore.PutActionEventState(workflowID, actionName, eventName, eventPayload)
		if err != nil {
			fm.log.Error("Failed to update action state", zap.String("action", actionName))
			continue
		}

		env := workflow.NewEnv(action)
		intp := workflow.NewInterpreter(env, nil, dummyResolver, fm.log)

		results, err := intp.EvalAction(action.Action)
		if err != nil {
			fm.log.Error("Failed to evaluate action",
				zap.Error(err), zap.String("action", actionName))
			continue
		}
		if results == nil {
			continue
		}

		fm.log.Debug("Action triggered", zap.String("action", actionName))

		// Reset action state upon triggering
		fm.stateStore.ResetActionState(workflowID, actionName)

		for _, r := range results {
			switch result := r.(type) {
			case *workflow.FireEventResult:
				payload, err := workflow.EncodeObjConst(result.Payload())
				if err != nil {
					fm.log.Error("Failed to encode fire event payload",
						zap.Error(err), zap.String("event", result.EventName()),
						zap.String("payload", workflow.PrintNode(result.Payload())))
					continue
				}

				err = fm.stateStore.PutActionEvent(&protos.ActionEvent{
					WorkflowId: workflowID.UnsafeBytes(),
					ActionName: actionName,
					EventName:  result.EventName(),
					Payload:    payload,
				})
				if err != nil {
					fm.log.Error("Failed to put action event",
						zap.Error(err), zap.String("event", result.EventName()))
					continue
				}

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

// doDeployWorkflow deploys a workflow.
func (fm *FlowManager) deployWorkflow(wf *workflow.WorkflowDecl) error {
	fm.lock.Lock()
	defer fm.lock.Unlock()

	wfid := GetWorkflowID(wf)
	if _, ok := fm.pendingWorkflows[wfid]; !ok {
		old, err := fm.stateStore.WorkflowByID(wfid)
		if err != nil {
			return err
		}
		if old == nil {
			fm.pendingWorkflows[wfid] = wf
			return nil
		}
	}

	return fmt.Errorf("Workflow already exists %v", wfid)
}

// undeployWorkflow undeploys a workflow.
func (fm *FlowManager) undeployWorkflow(wfid WorkflowID) error {
	fm.lock.Lock()
	defer fm.lock.Unlock()

	if wf, ok := fm.pendingWorkflows[wfid]; ok {
		if wf != nil {
			// The workflow was marked as to be deployed. Simply undo it.
			delete(fm.pendingWorkflows, wfid)
		}
		// Otherwise, it's already marked as to be deleted. Do nothing.
		return nil
	}

	old, err := fm.stateStore.WorkflowByID(wfid)
	if err != nil {
		return err
	}
	if old == nil {
		return fmt.Errorf("Workflow doesn't exists %v", wfid)
	}

	fm.pendingWorkflows[wfid] = nil
	return nil
}

// FinalizeAndCommitMBlock finalizes the pending mega block and commits it to state store.
func (fm *FlowManager) FinalizeAndCommitMBlock() (*protos.MBlock, error) {
	fm.lock.Lock()
	defer fm.lock.Unlock()

	block, configPatches, err := fm.finalizeMBlockWithLock()
	if err != nil {
		return nil, err
	}

	// Commit the pending block. From this point on there's no turning back.
	err = fm.stateStore.CommitPendingMBlock(block)
	if err != nil {
		return nil, err
	}

	// Finally, publish all chain config patches.
	for chain, patch := range configPatches {
		fm.publishChainConfigPatchWithLock(chain, patch)
	}

	return block, nil
}

func (fm *FlowManager) finalizeMBlockWithLock() (*protos.MBlock, map[string]*ChainConfigPatch, error) {
	configPatches := make(map[string]*ChainConfigPatch)
	getOrCreatePatch := func(chain string) *ChainConfigPatch {
		patch := configPatches[chain]
		if patch == nil {
			patch = new(ChainConfigPatch)
			configPatches[chain] = patch
		}
		return patch
	}

	// Process all workflow changes right before finalization, so they are always at the end of a block.
	wfids := make([]string, 0, len(fm.pendingWorkflows))
	for id := range fm.pendingWorkflows {
		wfids = append(wfids, string(id))
	}
	sort.Strings(wfids) // Make the order deterministic.

	var err error
	for _, id := range wfids {
		wfid := WorkflowID(id)
		wf := fm.pendingWorkflows[wfid]
		if wf != nil {
			// Deploy
			for _, m := range wf.MonitorDecls() {
				monitor, err := Monitor(wfid, m)
				if err != nil {
					return nil, nil, err
				}

				patch := getOrCreatePatch(m.Chain())
				patch.AddMonitors = append(patch.AddMonitors, monitor)
			}

			err = fm.stateStore.PutWorkflow(wfid, wf)
		} else {
			// Undeploy
			wf, err := fm.stateStore.WorkflowByID(wfid)
			if err != nil {
				return nil, nil, err
			}
			if wf == nil {
				fm.log.Error("Missing workflow during finalization", zap.Stringer("wfid", wfid))
				// Unexpected, but not fatal. Simply skip this request.
				continue
			}

			for _, m := range wf.MonitorDecls() {
				mid := &mgrpc.MonitorID{
					WorkflowID:  wfid,
					MonitorName: m.Name().Id(),
				}
				patch := getOrCreatePatch(m.Chain())
				patch.RemoveMonitors = append(patch.RemoveMonitors, mid)
			}

			err = fm.stateStore.DeleteWorkflow(wfid)
		}
		if err != nil {
			return nil, nil, err
		}
	}

	// Reset pending zone.
	fm.pendingWorkflows = make(map[WorkflowID]*workflow.WorkflowDecl)

	// Fill in missing info from chain config patches.
	for chain, patch := range configPatches {
		block, err := fm.stateStore.LastReportedBlockByChain(chain)
		if err != nil {
			return nil, nil, err
		}

		// It's possible there are no previous block reports for this chain.
		if block != nil {
			patch.ResumeAfter = &mgrpc.BlockSpec{
				Hash:   block.Hash,
				Height: block.Height,
			}
		}

		patch.MonitorsVersion = fm.chainConfigVersions[chain] + 1
	}

	block, err := fm.stateStore.FinalizeMBlock()
	if err != nil {
		return nil, nil, err
	}

	return block, configPatches, nil
}

func (fm *FlowManager) publishChainConfigPatchWithLock(chainID string, patch *ChainConfigPatch) {
	fm.chainConfigVersions[chainID] = patch.MonitorsVersion

	subs := fm.chainConfigSubs[chainID]
	if len(subs) == 0 {
		return
	}

	filteredSubs := subs[:0]
	for _, sub := range subs {
		select {
		case <-sub.cancel:
			fm.log.Info("Subscription canceled", zap.String("chainID", chainID))
			close(sub.patches)
			continue
		default:
		}

		select {
		case sub.patches <- patch:
			fm.log.Debug("Config patch published", zap.String("chainID", chainID))
		default:
			fm.log.Warn("Patch channel overflow", zap.String("chainID", chainID))
			close(sub.patches)
			continue
		}

		filteredSubs = append(filteredSubs, sub)
	}
	fm.chainConfigSubs[chainID] = filteredSubs
}

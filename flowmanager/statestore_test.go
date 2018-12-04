package flowmanager

import (
	"math/big"

	"github.com/megaspacelab/megaconnect/protos"
	w "github.com/megaspacelab/megaconnect/workflow"

	"github.com/golang/protobuf/proto"
	"github.com/stretchr/testify/suite"
	"go.uber.org/zap"
)

// StateStoreSuite contains a common set of tests for all StateStore implementations.
type StateStoreSuite struct {
	suite.Suite

	log        *zap.Logger
	stateStore StateStore
}

func (s *StateStoreSuite) SetupTest() {
	var err error
	s.log, err = zap.NewDevelopment()
	s.Require().NoError(err)
}

func (s *StateStoreSuite) workflow1() *w.WorkflowDecl {
	wf := w.NewWorkflowDecl(w.NewId("TestWorkflow"), 1)
	wf.AddDecl(w.NewEventDecl(
		w.NewId("TestEvent"),
		w.NewObjType(w.NewIdToTy().Put(
			"balance",
			w.IntType,
		)),
	))
	wf.AddDecl(w.NewEventDecl(
		w.NewId("TestEvent2"),
		w.NewObjType(w.NewIdToTy().Put(
			"balance",
			w.IntType,
		)),
	))
	wf.AddDecl(w.NewEventDecl(
		w.NewId("TestEvent3"),
		w.NewObjType(w.NewIdToTy().Put(
			"balance",
			w.IntType,
		)),
	))
	wf.AddDecl(w.NewMonitorDecl(
		w.NewId("TestMonitor"),
		w.TrueConst,
		w.NewIdToExpr().Put(
			"balance",
			w.NewIntConstFromI64(100),
		),
		w.NewFire(
			"TestEvent",
			w.NewObjLit(w.NewIdToExpr().Put(
				"balance",
				w.NewIntConstFromI64(1),
			)),
		),
		chainID,
	))
	wf.AddDecl(w.NewActionDecl(
		w.NewId("TestAction"),
		w.NewEVar("TestEvent"),
		[]w.Stmt{
			w.NewFire(
				"TestEvent2",
				w.NewProps(w.NewVar("TestEvent")),
			),
		},
	))
	wf.AddDecl(w.NewActionDecl(
		w.NewId("TestAction2"),
		w.NewEBinOp(
			w.AndEOp,
			w.NewEVar("TestEvent"),
			w.NewEVar("TestEvent2"),
		),
		[]w.Stmt{
			w.NewFire(
				"TestEvent3",
				w.NewProps(w.NewVar("TestEvent2")),
			),
		},
	))
	return wf
}

func (s *StateStoreSuite) TestPutAndGetAndDeleteWorkflow() {
	wf := s.workflow1()
	id := GetWorkflowID(wf)

	actualwf, err := s.stateStore.WorkflowByID(id)
	s.Require().NoError(err)
	s.Require().Nil(actualwf)

	err = s.stateStore.PutWorkflow(id, wf)
	s.Require().NoError(err)

	actualwf, err = s.stateStore.WorkflowByID(id)
	s.Require().NoError(err)
	s.Require().True(wf.Equal(actualwf))

	monitors, err := s.stateStore.MonitorsByChain(chainID)
	s.Require().NoError(err)

	monitorDecls := wf.MonitorDecls()
	s.Require().Len(monitors, len(monitorDecls))
	for i := 0; i < len(monitors); i++ {
		s.Require().Equal(id.UnsafeBytes(), monitors[i].WorkflowId)
		s.Require().Equal(monitorDecls[i].Name().Id(), monitors[i].MonitorName)

		decl, err := w.NewByteDecoder(monitors[i].Monitor).DecodeMonitorDecl()
		s.Require().NoError(err)
		s.Require().True(decl.Equal(monitorDecls[i]))
	}

	err = s.stateStore.DeleteWorkflow(id)
	s.Require().NoError(err)

	actualwf, err = s.stateStore.WorkflowByID(id)
	s.Require().NoError(err)
	s.Require().Nil(actualwf)

	monitors, err = s.stateStore.MonitorsByChain(chainID)
	s.Require().NoError(err)
	s.Require().Empty(monitors)

	// Shouldn't error when deleting non-existent.
	err = s.stateStore.DeleteWorkflow(id)
	s.Require().NoError(err)
}

func (s *StateStoreSuite) TestLastReportedBlockByChain() {
	actualBlock, err := s.stateStore.LastReportedBlockByChain(chainID)
	s.Require().NoError(err)
	s.Require().Nil(actualBlock)

	block := &protos.Block{
		Chain:  chainID,
		Height: protos.NewBigInt(big.NewInt(10)),
		Hash:   []byte{1, 2, 3},
	}

	err = s.stateStore.PutBlockReport(block)
	s.Require().NoError(err)

	actualBlock, err = s.stateStore.LastReportedBlockByChain(chainID)
	s.Require().NoError(err)
	s.Require().True(proto.Equal(actualBlock, block))
}

func (s *StateStoreSuite) TestActionStates() {
	wf := s.workflow1()
	id := GetWorkflowID(wf)
	err := s.stateStore.PutWorkflow(id, wf)
	s.Require().NoError(err)

	// Retrieve empty states.
	event := wf.EventDecls()[0]
	states, err := s.stateStore.ActionStatesByEvent(id, event.Name().Id())
	s.Require().NoError(err)
	s.Require().Len(states, 2)
	s.Require().True(states[0].Action.Equal(wf.ActionDecls()[0]))
	s.Require().Empty(states[0].State)
	s.Require().True(states[1].Action.Equal(wf.ActionDecls()[1]))
	s.Require().Empty(states[1].State)

	event = wf.EventDecls()[1]
	states, err = s.stateStore.ActionStatesByEvent(id, event.Name().Id())
	s.Require().NoError(err)
	s.Require().Len(states, 1)
	s.Require().True(states[0].Action.Equal(wf.ActionDecls()[1]))
	s.Require().Empty(states[0].State)

	// Put some state.
	value := w.NewObjConst(map[string]w.Const{
		"balance": w.NewIntConstFromI64(100),
	})
	s.Require().True(states[0].Observe(event.Name().Id(), value))

	action := wf.ActionDecls()[1]
	err = s.stateStore.PutActionEventState(id, action.Name().Id(), event.Name().Id(), value)
	s.Require().NoError(err)

	// Retrieve non-empty states.
	states, err = s.stateStore.ActionStatesByEvent(id, event.Name().Id())
	s.Require().NoError(err)
	s.Require().Len(states, 1)
	s.Require().True(states[0].Action.Equal(action))
	s.Require().Len(states[0].State, 1)
	s.Require().True(states[0].Lookup(event.Name().Id()).Equal(value))

	// Reset states.
	err = s.stateStore.ResetActionState(id, action.Name().Id())
	s.Require().NoError(err)

	// Retrieve empty states.
	states, err = s.stateStore.ActionStatesByEvent(id, event.Name().Id())
	s.Require().NoError(err)
	s.Require().Len(states, 1)
	s.Require().True(states[0].Action.Equal(action))
	s.Require().Empty(states[0].State)
}

func (s *StateStoreSuite) TestFinalizeAndCommitMBlock() {
	wf := s.workflow1()
	id := GetWorkflowID(wf)
	err := s.stateStore.PutWorkflow(id, wf)
	s.Require().NoError(err)

	block := &protos.Block{
		Chain:  chainID,
		Height: protos.NewBigInt(big.NewInt(10)),
	}
	err = s.stateStore.PutBlockReport(block)
	s.Require().NoError(err)

	monitorEvent := &protos.MonitorEvent{
		WorkflowId:  id.UnsafeBytes(),
		MonitorName: wf.MonitorDecls()[0].Name().Id(),
		EventName:   wf.MonitorDecls()[0].EventName(),
		Payload:     []byte{1, 1, 1},
	}
	err = s.stateStore.PutMonitorEvent(monitorEvent)
	s.Require().NoError(err)

	actionEvent := &protos.ActionEvent{
		WorkflowId: id.UnsafeBytes(),
		ActionName: wf.ActionDecls()[0].Name().Id(),
		EventName:  wf.EventDecls()[0].Name().Id(),
		Payload:    []byte{2, 2, 2},
	}
	err = s.stateStore.PutActionEvent(actionEvent)
	s.Require().NoError(err)

	mblock, err := s.stateStore.FinalizeMBlock()
	s.Require().NoError(err)

	s.Require().Len(mblock.Events, 4)

	err = s.stateStore.CommitPendingMBlock(mblock)
	s.Require().NoError(err)

	mblock2, err := s.stateStore.LatestMBlock()
	s.Require().NoError(err)
	s.Require().True(proto.Equal(mblock, mblock2))

	mblock2, err = s.stateStore.MBlockByHeight(mblock.Height)
	s.Require().NoError(err)
	s.Require().True(proto.Equal(mblock, mblock2))
}

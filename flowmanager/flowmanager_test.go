package flowmanager

import (
	"testing"

	"github.com/megaspacelab/megaconnect/grpc"
	w "github.com/megaspacelab/megaconnect/workflow"

	"github.com/stretchr/testify/suite"
	"go.uber.org/zap"
)

const chainID = "Example"

type FlowManagerSuite struct {
	suite.Suite

	log *zap.Logger
	fm  *FlowManager
}

func (s *FlowManagerSuite) SetupTest() {
	log, err := zap.NewDevelopment()
	s.Require().NoError(err)
	s.log = log
	s.fm = NewFlowManager(log)
}

func (s *FlowManagerSuite) workflow1() *w.WorkflowDecl {
	wf := w.NewWorkflowDecl(w.NewId("TestWorkflow"), 1)
	wf.AddChild(w.NewEventDecl(
		w.NewId("TestEvent"),
		w.NewObjType(w.NewIdToTy().Put(
			"balance",
			w.IntType,
		)),
	))
	wf.AddChild(w.NewEventDecl(
		w.NewId("TestEvent2"),
		w.NewObjType(w.NewIdToTy().Put(
			"balance",
			w.IntType,
		)),
	))
	wf.AddChild(w.NewMonitorDecl(
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
	wf.AddChild(w.NewActionDecl(
		w.NewId("TestAction"),
		w.NewEVar("TestEvent"),
		[]w.Stmt{
			w.NewFire(
				"TestEvent2",
				w.NewProps(w.NewVar("TestEvent")),
			),
		},
	))
	wf.AddChild(w.NewActionDecl(
		w.NewId("TestAction2"),
		w.NewEBinOp(
			w.AndEOp,
			w.NewEVar("TestEvent"),
			w.NewEVar("TestEvent2"),
		),
		[]w.Stmt{},
	))
	return wf
}

func (s *FlowManagerSuite) TestDeployWorkflow() {
	s.fm.SetChainConfig(chainID, nil, nil)
	config := s.fm.GetChainConfig(chainID)
	s.Require().NotNil(config)

	err := s.fm.DeployWorkflow(s.workflow1())
	s.Require().NoError(err)

	select {
	case <-config.Outdated:
	default:
		s.FailNow("Config wasn't updated after workflow deployment")
	}

	config = s.fm.GetChainConfig(chainID)
	s.Require().NotNil(config)
	s.Require().Len(config.Monitors, 1)
}

func (s *FlowManagerSuite) TestReportBlockEvents() {
	s.fm.SetChainConfig(chainID, nil, nil)
	config := s.fm.GetChainConfig(chainID)
	s.Require().NotNil(config)

	wf := s.workflow1()
	err := s.fm.DeployWorkflow(wf)
	s.Require().NoError(err)

	config = s.fm.GetChainConfig(chainID)
	s.Require().NotNil(config)

	block := &grpc.Block{}

	// This should be ignored
	s.fm.ReportBlockEvents(chainID, config.MonitorsVersion-1, block, nil)

	eventPayload, err := w.EncodeObjConst(
		w.NewObjConst(w.ObjFields{
			"balance": w.NewIntConstFromI64(1),
		}),
	)
	s.Require().NoError(err)

	// This should be processed
	s.fm.ReportBlockEvents(chainID, config.MonitorsVersion, block, []*grpc.Event{
		&grpc.Event{
			MonitorId: []byte(MonitorID(WorkflowID(wf), wf.MonitorDecls()[0])),
			Payload:   eventPayload,
		},
	})

	// TODO - verify event processing results once we have those
}

func TestFlowManager(t *testing.T) {
	suite.Run(t, new(FlowManagerSuite))
}

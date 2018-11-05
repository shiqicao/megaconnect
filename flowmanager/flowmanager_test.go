package flowmanager

import (
	"testing"

	"github.com/megaspacelab/megaconnect/grpc"
	"github.com/megaspacelab/megaconnect/workflow"

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

func (s *FlowManagerSuite) workflow1() *workflow.WorkflowDecl {
	wf := workflow.NewWorkflowDecl("TestWorkflow", 1)
	wf.AddChild(workflow.NewEventDecl(
		"TestEvent",
		workflow.NewObjType(workflow.NewIdToTy().Add(
			"balance",
			workflow.IntType,
		)),
	))
	wf.AddChild(workflow.NewEventDecl(
		"TestEvent2",
		workflow.NewObjType(workflow.NewIdToTy().Add(
			"balance",
			workflow.IntType,
		)),
	))
	wf.AddChild(workflow.NewMonitorDecl(
		"TestMonitor",
		workflow.TrueConst,
		workflow.NewIdToExpr().Add(
			"balance",
			workflow.NewIntConstFromI64(100),
		),
		workflow.NewFire(
			"TestEvent",
			workflow.NewObjLit(workflow.NewIdToExpr().Add(
				"balance",
				workflow.NewIntConstFromI64(1),
			)),
		),
		chainID,
	))
	wf.AddChild(workflow.NewActionDecl(
		"TestAction",
		workflow.NewEVar("TestEvent"),
		[]workflow.Stmt{
			workflow.NewFire(
				"TestEvent2",
				workflow.NewProps(workflow.NewVar("TestEvent")),
			),
		},
	))
	wf.AddChild(workflow.NewActionDecl(
		"TestAction2",
		workflow.NewEBinOp(
			workflow.AndEOp,
			workflow.NewEVar("TestEvent"),
			workflow.NewEVar("TestEvent2"),
		),
		[]workflow.Stmt{},
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

	eventPayload, err := workflow.EncodeObjConst(
		workflow.NewObjConst(workflow.ObjFields{
			"balance": workflow.NewIntConstFromI64(1),
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

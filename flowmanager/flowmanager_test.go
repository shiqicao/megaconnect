package flowmanager

import (
	"testing"

	"github.com/megaspacelab/megaconnect/grpc"
	"github.com/megaspacelab/megaconnect/workflow"

	"github.com/stretchr/testify/suite"
	"go.uber.org/zap"
)

type FlowManagerSuite struct {
	suite.Suite

	log *zap.Logger
	fm  *FlowManager
	wfs []*workflow.WorkflowDecl
}

func (s *FlowManagerSuite) SetupTest() {
	log, err := zap.NewDevelopment()
	s.Require().NoError(err)
	s.log = log
	s.fm = NewFlowManager(log)

	wf := workflow.NewWorkflowDecl("TestWorkflow", 1)
	wf.AddChild(workflow.NewEventDecl(
		"TestEvent",
		workflow.NewObjType(map[string]workflow.Type{
			"balance": workflow.IntType,
		}),
	))
	wf.AddChild(workflow.NewMonitorDecl(
		"TestMonitor",
		workflow.TrueConst,
		map[string]workflow.Expr{
			"balance": workflow.NewIntConstFromI64(100),
		},
		workflow.NewFire(
			"TestEvent",
			workflow.NewObjLit(workflow.VarDecls{
				"balance": workflow.NewIntConstFromI64(1),
			}),
		),
		"Example",
	))
	wf.AddChild(workflow.NewActionDecl(
		"TestAction",
		workflow.NewEVar("TestEvent"),
		[]workflow.Stmt{},
	))
	s.wfs = append(s.wfs, wf)
}

func (s *FlowManagerSuite) TestDeployWorkflow() {
	config := s.fm.SetChainConfig("Example", nil, nil)
	s.Require().NotNil(config)

	err := s.fm.DeployWorkflow(s.wfs[0])
	s.Require().NoError(err)

	select {
	case <-config.Outdated:
	default:
		s.FailNow("Config wasn't updated after workflow deployment")
	}

	config = s.fm.GetChainConfig("Example")
	s.Require().NotNil(config)
	s.Require().Len(config.Monitors, 1)
}

func (s *FlowManagerSuite) TestReportBlockEvents() {
	config := s.fm.SetChainConfig("Example", nil, nil)
	s.Require().NotNil(config)

	err := s.fm.DeployWorkflow(s.wfs[0])
	s.Require().NoError(err)

	config = s.fm.GetChainConfig("Example")
	s.Require().NotNil(config)

	block := &grpc.Block{}

	// This should be ignored
	s.fm.ReportBlockEvents("Example", config.MonitorsVersion-1, block, nil)

	eventPayload, err := workflow.EncodeObjConst(
		workflow.NewObjConst(workflow.ObjFields{"a": workflow.TrueConst}),
	)
	s.Require().NoError(err)

	// This should be processed
	s.fm.ReportBlockEvents("Example", config.MonitorsVersion, block, []*grpc.Event{
		&grpc.Event{
			MonitorId: []byte(MonitorID(WorkflowID(s.wfs[0]), s.wfs[0].MonitorDecls()[0])),
			Payload:   eventPayload,
		},
	})
}

func TestFlowManager(t *testing.T) {
	suite.Run(t, new(FlowManagerSuite))
}

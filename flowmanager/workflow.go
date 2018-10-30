package flowmanager

import (
	"strings"

	"github.com/megaspacelab/megaconnect/grpc"
	"github.com/megaspacelab/megaconnect/workflow"
)

// TODO - Generate proper & unique IDs.

// WorkflowID computes the ID of a workflow.
func WorkflowID(wf *workflow.WorkflowDecl) string {
	return wf.Name()
}

// MonitorID computes the ID of a monitor.
func MonitorID(wfid string, monitor *workflow.MonitorDecl) string {
	return strings.Join([]string{wfid, monitor.Name()}, "_")
}

// Monitor converts a MonitorDecl to grpc.Monitor.
func Monitor(id string, monitor *workflow.MonitorDecl) (*grpc.Monitor, error) {
	if monitor == nil {
		return nil, nil
	}

	raw, err := workflow.EncodeMonitorDecl(monitor)
	if err != nil {
		return nil, err
	}

	return &grpc.Monitor{
		Id:      []byte(id),
		Monitor: raw,
	}, nil
}

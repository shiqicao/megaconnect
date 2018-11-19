package grpc

import "github.com/megaspacelab/megaconnect/common"

// MonitorID is a go-friendly helper struct to id a monitor.
type MonitorID struct {
	WorkflowID  common.ImmutableBytes
	MonitorName string
}

func (m *Monitor) MonitorID() *MonitorID {
	if m == nil {
		return nil
	}

	return &MonitorID{
		WorkflowID:  common.ImmutableBytes(m.WorkflowId),
		MonitorName: m.MonitorName,
	}
}

func (m *UpdateMonitorsRequest_RemoveMonitor) MonitorID() *MonitorID {
	if m == nil {
		return nil
	}

	return &MonitorID{
		WorkflowID:  common.ImmutableBytes(m.WorkflowId),
		MonitorName: m.MonitorName,
	}
}

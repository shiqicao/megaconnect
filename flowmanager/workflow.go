// Copyright 2018 @ MegaSpace

// The MegaSpace source code is distributed in the hope that it will be useful,
// but WITHOUT ANY WARRANTY; without even the implied warranty of
// MERCHANTABILITY or FITNESS FOR A PARTICULAR PURPOSE. See the
// GNU Lesser General Public License for more details.

// You should have received a copy of the GNU Lesser General Public License
// along with the MegaSpace source code. If not, see <http://www.gnu.org/licenses/>.

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

// EventID computes the ID of an event.
func EventID(wfid string, eventName string) string {
	return strings.Join([]string{wfid, eventName}, "_")
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

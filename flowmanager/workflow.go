// Copyright 2018 @ MegaSpace

// The MegaSpace source code is distributed in the hope that it will be useful,
// but WITHOUT ANY WARRANTY; without even the implied warranty of
// MERCHANTABILITY or FITNESS FOR A PARTICULAR PURPOSE. See the
// GNU Lesser General Public License for more details.

// You should have received a copy of the GNU Lesser General Public License
// along with the MegaSpace source code. If not, see <http://www.gnu.org/licenses/>.

package flowmanager

import (
	"github.com/megaspacelab/megaconnect/common"
	"github.com/megaspacelab/megaconnect/grpc"
	"github.com/megaspacelab/megaconnect/workflow"
)

// TODO - Change to fixed size
type WorkflowID = common.ImmutableBytes

// TODO - Generate proper & unique IDs.

// GetWorkflowID computes the ID of a workflow.
func GetWorkflowID(wf *workflow.WorkflowDecl) WorkflowID {
	return WorkflowID(wf.Name().Id())
}

// Monitor converts a MonitorDecl to grpc.Monitor.
func Monitor(wfid WorkflowID, monitor *workflow.MonitorDecl) (*grpc.Monitor, error) {
	if monitor == nil {
		return nil, nil
	}

	raw, err := workflow.EncodeMonitorDecl(monitor)
	if err != nil {
		return nil, err
	}

	return &grpc.Monitor{
		WorkflowId:  wfid.Bytes(),
		MonitorName: monitor.Name().Id(),
		Monitor:     raw,
	}, nil
}

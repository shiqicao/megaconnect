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
)

type chainConfigBase struct {
	// MonitorsVersion is used for Orchestrator and ChainManager to agree on the set of active monitors.
	MonitorsVersion uint32

	// ResumeAfter specifies the LAST block BEFORE this ChainConfig became effective.
	ResumeAfter *grpc.BlockSpec
}

// ChainConfig captures configs about a chain.
type ChainConfig struct {
	chainConfigBase

	// Monitors are the active monitors for this chain.
	Monitors []*grpc.Monitor
}

// ChainConfigPatch represents a patch that can be applied to a ChainConfig.
type ChainConfigPatch struct {
	chainConfigBase

	AddMonitors    []*grpc.Monitor
	RemoveMonitors []*grpc.MonitorID
}

// ChainConfigSub is a subscription to chain config updates.
type ChainConfigSub struct {
	patches chan *ChainConfigPatch
	cancel  chan common.Nothing
}

// Patches returns a channel of ChainConfigPatch.
// This channel will be closed by the producer in case of an error.
func (s *ChainConfigSub) Patches() <-chan *ChainConfigPatch {
	return s.patches
}

// Unsubscribe allows the consumer to unsubscribe.
// It should only be called once.
func (s *ChainConfigSub) Unsubscribe() {
	close(s.cancel)
}

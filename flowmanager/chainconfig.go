// Copyright 2018 @ MegaSpace

// The MegaSpace source code is distributed in the hope that it will be useful,
// but WITHOUT ANY WARRANTY; without even the implied warranty of
// MERCHANTABILITY or FITNESS FOR A PARTICULAR PURPOSE. See the
// GNU Lesser General Public License for more details.

// You should have received a copy of the GNU Lesser General Public License
// along with the MegaSpace source code. If not, see <http://www.gnu.org/licenses/>.

package flowmanager

import "github.com/megaspacelab/megaconnect/grpc"

type chainConfigBase struct {
	// MonitorsVersion is used for Orchestrator and ChainManager to agree on the set of active monitors.
	MonitorsVersion uint32

	// ResumeAfter specifies the LAST block BEFORE this ChainConfig became effective.
	ResumeAfter *grpc.BlockSpec

	// Outdated is closed when this ChainConfig becomes outdated.
	Outdated chan struct{}
}

// ChainConfig captures configs about a chain.
type ChainConfig struct {
	chainConfigBase

	// Monitors are the active monitors for this chain.
	Monitors IndexedMonitors
}

func (cc *ChainConfig) copy() *ChainConfig {
	if cc == nil {
		return nil
	}
	cp := *cc
	cp.Monitors = cp.Monitors.Copy()
	return &cp
}

// ChainConfigPatch represents a patch that can be applied to a ChainConfig.
type ChainConfigPatch struct {
	chainConfigBase

	AddMonitors    []*grpc.Monitor
	RemoveMonitors []string
}

// Apply applies this patch to cc.
// cc will be patched in place, unless it's nil, in which case a new ChainConfig is created.
func (ccp *ChainConfigPatch) Apply(cc *ChainConfig) *ChainConfig {
	if ccp == nil {
		return cc
	}

	if cc == nil {
		cc = &ChainConfig{
			chainConfigBase: ccp.chainConfigBase,
		}
	} else {
		cc.chainConfigBase = ccp.chainConfigBase
	}

	for _, m := range ccp.RemoveMonitors {
		delete(cc.Monitors, m)
	}

	if cc.Monitors == nil && len(ccp.AddMonitors) > 0 {
		cc.Monitors = make(IndexedMonitors, len(ccp.AddMonitors))
	}

	for _, m := range ccp.AddMonitors {
		cc.Monitors[string(m.Id)] = m
	}

	return cc
}

// IndexedMonitors is a bunch of monitors indexed by ID.
type IndexedMonitors map[string]*grpc.Monitor

// Monitors returns all monitors contained in this IndexedMonitors.
func (im IndexedMonitors) Monitors() []*grpc.Monitor {
	monitors := make([]*grpc.Monitor, 0, len(im))
	for _, m := range im {
		monitors = append(monitors, m)
	}
	return monitors
}

// Copy makes a shallow copy of the IndexedMonitors.
func (im IndexedMonitors) Copy() IndexedMonitors {
	if im == nil {
		return nil
	}

	cp := make(IndexedMonitors, len(im))
	for k, v := range im {
		cp[k] = v
	}
	return cp
}

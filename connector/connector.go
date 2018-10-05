// Copyright 2018 @ MegaSpace

// The MegaSpace source code is distributed in the hope that it will be useful,
// but WITHOUT ANY WARRANTY; without even the implied warranty of
// MERCHANTABILITY or FITNESS FOR A PARTICULAR PURPOSE. See the
// GNU Lesser General Public License for more details.

// You should have received a copy of the GNU Lesser General Public License
// along with the MegaSpace source code. If not, see <http://www.gnu.org/licenses/>.

package connector

import (
	"github.com/megaspacelab/eventmanager/common"
	"github.com/megaspacelab/eventmanager/types"
)

// Connector defines the shared structure for each chain-specific connector.
type Connector interface {
	// Name returns the name of the connector, specific to the client implementation
	Name() string

	// ChainName returns the name of the blockchain, eg., Bitcoin, Ethereum, Stella
	ChainName() string

	// SubscribeBlock establishes blocks channel to accept new blocks.
	// It supports "resume after" functionalities by passing in the hash of the checkpoint block
	SubscribeBlock(resumeAfter *common.Hash, blocks chan<- types.Block) (Subscription, error)

	// Start starts the Connector, proper setup is done here
	Start() error

	// Stop stops the Connector and call Unsubscribe on subscription channels
	Stop() error

	// QueryAccountBalance gets account balance of given address on demand
	// It returns a float64 value in local unit
	QueryAccountBalance(addr string) (<-chan float64, error)
}

// Subscription defines the shared structure for each connector's new block subscription.
type Subscription interface {
	// Unsubscribe is called when subscription is closed
	Unsubscribe()
}

// Copyright 2018 @ MegaSpace

// The MegaSpace source code is distributed in the hope that it will be useful,
// but WITHOUT ANY WARRANTY; without even the implied warranty of
// MERCHANTABILITY or FITNESS FOR A PARTICULAR PURPOSE. See the
// GNU Lesser General Public License for more details.

// You should have received a copy of the GNU Lesser General Public License
// along with the MegaSpace source code. If not, see <http://www.gnu.org/licenses/>.

package connector

import (
	"math/big"

	"github.com/megaspacelab/eventmanager/common"
	"github.com/megaspacelab/eventmanager/types"
)

// Connector defines the shared structure for each chain-specific connector.
type Connector interface {
	// Name returns the name of the Connector, specific to the client implementation.
	Name() string

	// ChainName returns the name of the blockchain, eg., Bitcoin, Ethereum, Stella.
	ChainName() string

	// SubscribeBlock establishes blocks channel to accept new blocks.
	// It supports "resume after" functionality by passing in the hash of the checkpoint block.
	SubscribeBlock(resumeAfter *common.Hash, blocks chan<- types.Block) (Subscription, error)

	// Start starts the Connector, proper setup is done here.
	Start() error

	// Stop stops the Connector and calls Unsubscribe on subscription channels.
	Stop() error

	// QueryAccountBalance gets account balance of given address on demand.
	// addr is the string representation of the address, asOfBlock is the hash
	// of the specific block that balance is querying on.
	// If asOfBlock is set to nil, the balance is taken from the latest known block.
	// return result a big integer in base unit.
	QueryAccountBalance(addr string, asOfBlock *common.Hash) (*big.Int, error)

	// IsFresh performs health check on connected blockchains on sync status.
	// Returns true if it is fresh and fully synced.
	// Returns false if it is still syncing and not ready to serve live traffic.
	IsFresh() (bool, error)
}

// Subscription defines the shared structure for each connector's new block subscription.
type Subscription interface {
	// Unsubscribe is called when subscription is closed.
	Unsubscribe()
}

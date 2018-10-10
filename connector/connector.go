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

	"github.com/megaspacelab/megaconnect/common"
)

// Connector defines the shared structure for each chain-specific connector.
type Connector interface {
	// Name returns the name of the Connector, specific to the client implementation.
	Name() string

	// ChainName returns the name of the blockchain, eg., Bitcoin, Ethereum, Stella.
	ChainName() string

	// SubscribeBlock establishes blocks channel to accept new blocks.
	// It supports "resume after" functionality by passing in the hash of the checkpoint block.
	SubscribeBlock(resumeAfter *common.Hash, blocks chan<- common.Block) (Subscription, error)

	// Start starts the Connector, proper setup is done here.
	Start() error

	// Stop stops the Connector and calls Unsubscribe on subscription channels.
	Stop() error

	// QueryAccountBalance gets account balance of given height on demand.
	// addr is the string representation of the address, height is the block height
	// that we are trying to get balance from.
	// If height is set to nil, get latest account balance.
	QueryAccountBalance(addr string, height *big.Int) (*big.Int, error)

	// IsHealthy performs health check on connected blockchains on sync status.
	// Returns true if it is fresh and fully synced.
	// Returns false if it is still syncing and not ready to serve live traffic.
	IsHealthy() (bool, error)
}

// Subscription defines the shared structure for each connector's new block subscription.
type Subscription interface {
	// Unsubscribe is called when subscription is closed.
	Unsubscribe()
}

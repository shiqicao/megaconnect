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
	"time"

	"github.com/megaspacelab/megaconnect/common"
)

// Connector defines the shared structure for each chain-specific connector.
type Connector interface {
	// Metadata returns the metadata of connected blockchain and the connector.
	Metadata() *Metadata

	// Start starts the Connector, proper setup is done here.
	Start() error

	// Stop stops the Connector and its subscription.
	Stop() error

	// SubscribeBlock establishes blocks channel to accept new blocks.
	// resumeAfter tells the connector to potentially rewind to an older block.
	// When Hash is specified in resumeAfter, it should be used as the primary way of identifying the block.
	// Height should be used only if Hash based lookup failed.
	// It is an error to call SubscribeBlock multiple times on a connector.
	SubscribeBlock(resumeAfter *BlockSpec) (Subscription, error)

	// QueryAccountBalance gets account balance of given height on demand.
	// addr is the string representation of the address, height is the block height
	// that we are trying to get balance from.
	// If height is set to nil, get latest account balance.
	QueryAccountBalance(addr string, height *big.Int) (*big.Int, error)

	// IsHealthy performs health check on connected blockchains on sync status.
	// Returns true if it is fresh and fully synced.
	// Returns false if it is still syncing and not ready to serve live traffic.
	IsHealthy() (bool, error)

	// IsValidAddress checks if the string could represent a valid address on
	// the connected blockchain
	IsValidAddress(addr string) bool
}

// Metadata has information of connector and the connected blockchain
type Metadata struct {
	ConnectorID            string        // The name of the Connector, specific to the client implementation.
	ChainID                string        // The name of the blockchain, eg., Bitcoin, Ethereum, Stella.
	HealthCheckInterval    time.Duration // How often health check is performed.
	HealthCheckGracePeriod time.Duration // Length of health check grace period before changing state.
}

// Subscription defines the shared structure for each connector's new block subscription.
type Subscription interface {
	// Blocks returns the block subscription channel to be consumed from.
	// It'll be closed upon subscription failure or connector stopping.
	Blocks() <-chan common.Block

	// Err channel will receive an error if subscription failed.
	// It'll be closed upon subscription failure or connector stopping.
	Err() <-chan error
}

// BlockSpec is used to specify a block to be looked up.
// While Hash is the most accurate specification, when not available, Height can be used instead.
type BlockSpec struct {
	Hash   *common.Hash
	Height *big.Int
}

// GetHash returns Hash. It returns nil if this BlockSpec is nil.
func (bs *BlockSpec) GetHash() *common.Hash {
	if bs == nil {
		return nil
	}
	return bs.Hash
}

// GetHeight returns Height. It returns nil if this BlockSpec is nil.
func (bs *BlockSpec) GetHeight() *big.Int {
	if bs == nil {
		return nil
	}
	return bs.Height
}

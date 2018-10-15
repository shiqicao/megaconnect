// Copyright 2018 @ MegaSpace

// The MegaSpace source code is distributed in the hope that it will be useful,
// but WITHOUT ANY WARRANTY; without even the implied warranty of
// MERCHANTABILITY or FITNESS FOR A PARTICULAR PURPOSE. See the
// GNU Lesser General Public License for more details.

// You should have received a copy of the GNU Lesser General Public License
// along with the MegaSpace source code. If not, see <http://www.gnu.org/licenses/>.

package example

import (
	"fmt"
	"math"
	"math/big"
	"math/rand"
	"reflect"
	"time"

	"github.com/megaspacelab/megaconnect/common"
	"github.com/megaspacelab/megaconnect/connector"
	"go.uber.org/zap"
)

// Connector is an example implementation of connector.Connector.
type Connector struct {
	balanceCache []*big.Int
	subs         []connector.Subscription
	running      bool
	logger       *zap.Logger
}

// New creates a new example connector.
func New(logger *zap.Logger) (connector.Connector, error) {
	return &Connector{logger: logger}, nil
}

// IsHealthy always returns true for example connector
func (c *Connector) IsHealthy() (bool, error) {
	return true, nil
}

// Name returns the name of this connector.
func (c *Connector) Name() string {
	return reflect.TypeOf(c).PkgPath()
}

// ChainName returns the name of the blockchain backing this connector.
func (c *Connector) ChainName() string {
	return "Example"
}

// Start starts this connector.
func (c *Connector) Start() error {
	if c.running {
		c.logger.Error("connector is already running")
		return fmt.Errorf("")
	}
	c.running = true
	c.logger.Info("connector starts", zap.String("name", c.Name()))
	return nil
}

// Stop cancels all subscriptions and stops the connector.
func (c *Connector) Stop() error {
	if !c.running {
		c.logger.Warn("connector is not running", zap.String("name", c.Name()))
		return nil
	}
	for _, sub := range c.subs {
		sub.Unsubscribe()
	}
	c.subs = nil
	c.running = false
	return nil
}

// SubscribeBlock subscribes to blocks on this chain, starting from after the specified block.
// If resumeAfter is nil, the subscription starts from the next block.
func (c *Connector) SubscribeBlock(resumeAfter *common.Hash, blocks chan<- common.Block) (connector.Subscription, error) {
	var height int64

	done := make(chan common.Nothing, 1)

	go func() {
		for ; ; height++ {
			block := common.NewBlock(
				common.Hash{},
				common.Hash{},
				big.NewInt(height),
				nil,
			)

			select {
			case blocks <- block:
				c.balanceCache = append(c.balanceCache, big.NewInt(rand.Int63n(math.MaxInt64)))
				time.Sleep(5 * time.Second)
			case <-done:
				close(blocks)
				return
			}
		}
	}()

	sub := subscription(done)
	c.subs = append(c.subs, sub)
	return sub, nil
}

type subscription chan<- common.Nothing

// Unsubscribe closes the subscribed channel.
func (s subscription) Unsubscribe() {
	s <- common.Nothing{}
}

// QueryAccountBalance queries the chain for the current balance of the given address.
// Returns a channel into which the result will be pushed once retrieved.
func (c *Connector) QueryAccountBalance(addr string, height *big.Int) (*big.Int, error) {
	h := height.Int64()
	if h < 0 {
		return new(big.Int).SetInt64(0), nil
	} else if int(h) >= len(c.balanceCache) {
		return nil, fmt.Errorf("height %d is not yet observed", h)
	}
	return c.balanceCache[h], nil
}

// IsValidAddress checks if the address string is valid
func (c *Connector) IsValidAddress(addr string) bool {
	return true
}

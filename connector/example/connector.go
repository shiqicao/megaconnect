// Copyright 2018 @ MegaSpace

// The MegaSpace source code is distributed in the hope that it will be useful,
// but WITHOUT ANY WARRANTY; without even the implied warranty of
// MERCHANTABILITY or FITNESS FOR A PARTICULAR PURPOSE. See the
// GNU Lesser General Public License for more details.

// You should have received a copy of the GNU Lesser General Public License
// along with the MegaSpace source code. If not, see <http://www.gnu.org/licenses/>.

package example

import (
	"errors"
	"fmt"
	"math/big"
	"reflect"
	"sync"
	"time"

	"github.com/megaspacelab/megaconnect/common"
	"github.com/megaspacelab/megaconnect/connector"

	"go.uber.org/zap"
)

// Connector is an example implementation of connector.Connector.
type Connector struct {
	logger        *zap.Logger
	blockInterval time.Duration

	subs    []*subscription
	height  int64
	running bool

	producerDone chan common.Nothing
	lock         sync.Mutex
}

// New creates a new example connector.
func New(logger *zap.Logger, blockInterval time.Duration) (connector.Connector, error) {
	return &Connector{logger: logger, blockInterval: blockInterval}, nil
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
	c.lock.Lock()
	defer c.lock.Unlock()

	if c.running {
		return errors.New("Connector is already running")
	}

	c.producerDone = make(chan common.Nothing)
	go c.produceBlocks(c.producerDone)

	c.running = true
	c.logger.Info("Connector started", zap.String("name", c.Name()))
	return nil
}

// Stop cancels all subscriptions and stops the connector.
func (c *Connector) Stop() error {
	c.lock.Lock()
	defer c.lock.Unlock()

	if !c.running {
		c.logger.Warn("Connector is not running", zap.String("name", c.Name()))
		return nil
	}
	close(c.producerDone)
	c.producerDone = nil
	c.subs = nil
	c.running = false
	c.logger.Info("Connector stopped", zap.String("name", c.Name()))
	return nil
}

// SubscribeBlock subscribes to blocks on this chain.
// resumeAfter is not supported in this connector. Subscription always starts from the next block.
func (c *Connector) SubscribeBlock(resumeAfter *connector.BlockSpec, blocks chan<- common.Block) (connector.Subscription, error) {
	c.lock.Lock()
	defer c.lock.Unlock()

	if !c.running {
		return nil, errors.New("Connector is not running")
	}

	sub := &subscription{
		blocks: blocks,
		done:   make(chan common.Nothing),
	}
	c.subs = append(c.subs, sub)
	return sub, nil
}

func (c *Connector) produceBlocks(done <-chan common.Nothing) {
	c.logger.Debug("Block producer started")
	for {
		select {
		case <-time.After(c.blockInterval):
		case <-done:
			c.logger.Debug("Block producer stopped")
			return
		}

		block := common.NewBlock(
			common.Hash{},
			common.Hash{},
			big.NewInt(c.height),
			nil,
		)

		func() {
			c.lock.Lock()
			defer c.lock.Unlock()

			if done != c.producerDone {
				// This is no longer the active producer. Do nothing.
				return
			}

			// Filter out closed subs as part of iteration.
			// https://github.com/golang/go/wiki/SliceTricks#filtering-without-allocating
			subs := c.subs[:0]
			for _, sub := range c.subs {
				select {
				case <-sub.done:
					close(sub.blocks)
				default:
					subs = append(subs, sub)
					select {
					case sub.blocks <- block:
					default:
						c.logger.Warn("Skipping blocked subscription")
					}
				}
			}
			c.subs = subs

			c.height++
		}()
	}
}

type subscription struct {
	blocks chan<- common.Block
	done   chan common.Nothing
	once   sync.Once
}

// Unsubscribe closes the subscribed channel.
func (s *subscription) Unsubscribe() {
	s.once.Do(func() {
		close(s.done)
	})
}

// QueryAccountBalance queries the chain for the balance of the given address.
func (c *Connector) QueryAccountBalance(addr string, height *big.Int) (*big.Int, error) {
	h := height.Int64()

	c.lock.Lock()
	defer c.lock.Unlock()

	if h > c.height {
		return nil, fmt.Errorf("height %d is not yet observed", h)
	}
	return big.NewInt(h + 100), nil
}

// IsValidAddress checks if the address string is valid
func (c *Connector) IsValidAddress(addr string) bool {
	return true
}

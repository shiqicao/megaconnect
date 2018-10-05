// Copyright 2018 @ MegaSpace

// The MegaSpace source code is distributed in the hope that it will be useful,
// but WITHOUT ANY WARRANTY; without even the implied warranty of
// MERCHANTABILITY or FITNESS FOR A PARTICULAR PURPOSE. See the
// GNU Lesser General Public License for more details.

// You should have received a copy of the GNU Lesser General Public License
// along with the MegaSpace source code. If not, see <http://www.gnu.org/licenses/>.

package exampleconn

import (
	"fmt"
	"math/big"
	"reflect"
	"time"

	"github.com/megaspacelab/eventmanager/common"
	"github.com/megaspacelab/eventmanager/connector"
	"github.com/megaspacelab/eventmanager/types"
	"go.uber.org/zap"
)

// Builder builds the connector.
type Builder struct{}

// BuildConnector returns a new stopped connector.
func (b *Builder) BuildConnector(ctx *connector.Context) (connector.Connector, error) {
	return &Connector{
		logger: ctx.Logger,
	}, nil
}

// Connector is the main connector data structure.
type Connector struct {
	subs    []connector.Subscription
	running bool
	logger  *zap.Logger
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
func (c *Connector) SubscribeBlock(resumeAfter types.Block, blocks chan<- types.Block) (connector.Subscription, error) {
	var height int64
	if resumeAfter != nil {
		height = resumeAfter.Height().Int64() + 1
	}

	done := make(chan common.Nothing, 1)

	go func() {
		for ; ; height++ {
			block := types.NewBlock(
				common.Hash{},
				common.Hash{},
				big.NewInt(height),
				nil,
			)

			select {
			case blocks <- block:
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
func (c *Connector) QueryAccountBalance(addr string) (<-chan float64, error) {
	result := make(chan float64, 1)

	go func() {
		defer close(result)

		var (
			totalReceived float64
			totalSpent    float64
			total         float64
		)

		// Fetch transactions for address to deduce the total amounts received and spent.
		totalReceived = 0.37073317
		totalSpent = 0.20285404

		// Subtract totalSpent from totalReceived for the total.
		total = totalReceived - totalSpent

		// Push total result into channel
		result <- total
	}()

	return result, nil
}

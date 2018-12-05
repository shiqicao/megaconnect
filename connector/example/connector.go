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
	"sync/atomic"
	"time"

	"github.com/megaspacelab/megaconnect/common"
	"github.com/megaspacelab/megaconnect/connector"
	wf "github.com/megaspacelab/megaconnect/workflow"
	"go.uber.org/zap"
)

const (
	blockSubscriptionBufferSize = 100
	chainId                     = "Example"
)

var one = big.NewInt(1)

// Connector is an example implementation of connector.Connector.
// This implementation is not thread-safe. Caller should serialize access to all public functions.
type Connector struct {
	logger        *zap.Logger
	blockInterval time.Duration

	sub     *subscription
	height  atomic.Value
	running bool

	done chan common.Nothing
}

// New creates a new example connector.
func New(logger *zap.Logger, blockInterval time.Duration) (connector.Connector, error) {
	return &Connector{logger: logger, blockInterval: blockInterval}, nil
}

// Metadata returns the metadata of this connector.
func (c *Connector) Metadata() *connector.Metadata {
	return &connector.Metadata{
		ConnectorID:            "Example Connector",
		ChainID:                chainId,
		HealthCheckInterval:    50 * time.Millisecond,
		HealthCheckGracePeriod: 250 * time.Millisecond,
	}
}

// IsHealthy always returns true for example connector.
func (c *Connector) IsHealthy() (bool, error) {
	return true, nil
}

// Start starts this connector.
func (c *Connector) Start() error {
	if c.running {
		return errors.New("Connector is already running")
	}

	c.done = make(chan common.Nothing)
	c.running = true
	c.logger.Info("Connector started", zap.String("name", c.Metadata().ConnectorID))
	return nil
}

// Stop cancels all subscriptions and stops the connector.
func (c *Connector) Stop() error {
	if !c.running {
		return errors.New("Connector is not running")
	}

	close(c.done)
	c.done = nil
	c.sub = nil
	c.running = false
	c.logger.Info("Connector stopped", zap.String("name", c.Metadata().ConnectorID))
	return nil
}

// SubscribeBlock subscribes to blocks on this chain.
func (c *Connector) SubscribeBlock(resumeAfter *connector.BlockSpec) (connector.Subscription, error) {
	if !c.running {
		return nil, errors.New("Connector is not running")
	}
	if c.sub != nil {
		return nil, errors.New("Subscription already exists")
	}

	c.sub = &subscription{
		blocks: make(chan common.Block, blockSubscriptionBufferSize),
		err:    make(chan error, 1),
	}
	sub := c.sub
	done := c.done

	if h := resumeAfter.GetHeight(); h != nil {
		c.height.Store(h)
	}

	go func() {
		defer close(sub.blocks)
		defer close(sub.err)

		for {
			select {
			case <-time.After(c.blockInterval):
			case <-done:
				return
			}

			height, ok := c.height.Load().(*big.Int)
			if !ok {
				height = new(big.Int)
			} else {
				// Always make a copy of height instead of mutating it.
				height = new(big.Int).Add(height, one)
			}
			c.height.Store(height)

			block := common.NewBlock(
				common.Hash{},
				common.Hash{},
				height,
				nil,
			)

			select {
			case sub.blocks <- block:
			default:
				sub.err <- errors.New("Block subscription buffer overflow")
				return
			}
		}
	}()

	return sub, nil
}

// QueryAccountBalance queries the chain for the balance of the given address.
func (c *Connector) QueryAccountBalance(addr string, height *big.Int) (*big.Int, error) {
	if !c.running {
		return nil, errors.New("Connector is not running")
	}

	if h, ok := c.height.Load().(*big.Int); !ok || h.Cmp(height) < 0 {
		return nil, fmt.Errorf("height %v is not yet observed", height)
	}

	n := big.NewInt(100)
	return n.Add(height, n), nil
}

// IsValidAddress checks if the address string is valid
func (c *Connector) IsValidAddress(addr string) bool {
	return true
}

// Namespace returns the namespace and APIs supported by workflow language for the connector.
func (c *Connector) Namespace() (*wf.NamespaceDecl, error) {
	ns := wf.NewNamespaceDecl(chainId)
	err := ns.AddFuncs(
		wf.NewFuncDecl("GetBlockInterval", wf.Params{}, wf.IntType, func(_ *wf.Env, _ map[string]wf.Const) (wf.Const, error) {
			return wf.NewIntConstFromI64(c.blockInterval.Nanoseconds()), nil
		}),
		wf.NewFuncDecl(
			"Echo",
			wf.Params{wf.NewParamDecl("x", wf.StrType)},
			wf.StrType,
			func(_ *wf.Env, args map[string]wf.Const) (wf.Const, error) {
				x := args["x"]
				return x, nil
			},
		),
	)
	if err != nil {
		c.logger.Error("Failed to create namespace", zap.Error(err))
		return nil, err
	}
	return ns, nil
}

type subscription struct {
	blocks chan common.Block
	err    chan error
}

func (s *subscription) Blocks() <-chan common.Block { return s.blocks }
func (s *subscription) Err() <-chan error           { return s.err }

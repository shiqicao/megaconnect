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

type Builder struct{}

func (b *Builder) BuildConnector(ctx *connector.Context) (connector.Connector, error) {
	return &Connector{
		logger: ctx.Logger,
	}, nil
}

type Connector struct {
	subs    []connector.Subscription
	running bool
	logger  *zap.Logger
}

func (c *Connector) Name() string {
	return reflect.TypeOf(c).PkgPath()
}

func (c *Connector) ChainName() string {
	return "Example"
}

func (c *Connector) Start() error {
	if c.running {
		c.logger.Error("connector is already running")
		return fmt.Errorf("")
	}
	c.running = true
	c.logger.Info("connector starts", zap.String("name", c.Name()))
	return nil
}

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

func (s subscription) Unsubscribe() {
	s <- common.Nothing{}
}

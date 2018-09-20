// Copyright 2018 @ MegaSpace

// The MegaSpace source code is distributed in the hope that it will be useful,
// but WITHOUT ANY WARRANTY; without even the implied warranty of
// MERCHANTABILITY or FITNESS FOR A PARTICULAR PURPOSE. See the
// GNU Lesser General Public License for more details.

// You should have received a copy of the GNU Lesser General Public License
// along with the MegaSpace source code. If not, see <http://www.gnu.org/licenses/>.

package eventmanager

import (
	"fmt"
	"sync"

	"go.uber.org/zap"

	conn "github.com/megaspacelab/eventmanager/connector"
	"github.com/megaspacelab/eventmanager/types"
)

type EventManager struct {
	connector conn.Connector

	running bool

	logger *zap.Logger
	lock   sync.Mutex
}

// New constructs an instance of eventManager
func New(conn conn.Connector, logger *zap.Logger) *EventManager {
	return &EventManager{
		connector: conn,
		logger:    logger,
	}
}

func (e *EventManager) Start() error {
	e.lock.Lock()
	defer e.lock.Unlock()
	if e.running {
		e.logger.Warn("event manager is already running")
		return fmt.Errorf("")
	}

	err := e.connector.Start()
	if err != nil {
		return err
	}

	blocks := make(chan types.Block, 100)
	_, err = e.connector.SubscribeBlock(nil, blocks)
	if err != nil {
		return err
	}

	go func() {
		for block := range blocks {
			e.processBlock(block)
		}
	}()

	e.running = true
	return nil
}

func (e *EventManager) Stop() error {
	e.lock.Lock()
	defer e.lock.Unlock()

	if !e.running {
		e.logger.Warn("event manager is not running")
		return fmt.Errorf("")
	}

	err := e.connector.Stop()
	if err != nil {
		e.logger.Error("connect stop with err", zap.Error(err))
	}

	e.running = false
	return nil
}

func (e *EventManager) processBlock(block types.Block) {
	e.logger.Info("Received new block", zap.Stringer("height", block.Height()))
}

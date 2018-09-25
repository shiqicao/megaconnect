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

	conn "github.com/megaspacelab/eventmanager/connector"
	"github.com/megaspacelab/eventmanager/types"
	"github.com/robfig/cron"
	"go.uber.org/zap"
)

const (
	blockBuffer = 100
)

// EventManager defines the shared event manager process for Connector
type EventManager struct {
	connector conn.Connector // connector
	cron      *cron.Cron     // scheduler
	running   bool           // if event manager is running
	logger    *zap.Logger    // logging
	lock      sync.Mutex     // mutex lock
}

// New constructs an instance of eventManager
func New(conn conn.Connector, logger *zap.Logger) *EventManager {
	return &EventManager{
		connector: conn,
		logger:    logger,
	}
}

// Start would start an EventManager loop
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

	blocks := make(chan types.Block, blockBuffer)
	_, err = e.connector.SubscribeBlock(nil, blocks)
	if err != nil {
		return err
	}

	e.cron = cron.New()
	e.cron.Start()

	go func() {
		for block := range blocks {
			e.processBlock(block)
		}
	}()

	e.running = true
	return nil
}

// Stop would stop an EventManager loop
func (e *EventManager) Stop() error {
	e.lock.Lock()
	defer e.lock.Unlock()

	if !e.running {
		e.logger.Warn("event manager is not running")
		return fmt.Errorf("")
	}

	e.cron.Stop()
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

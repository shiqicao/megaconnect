// Copyright 2018 @ MegaSpace

// The MegaSpace source code is distributed in the hope that it will be useful,
// but WITHOUT ANY WARRANTY; without even the implied warranty of
// MERCHANTABILITY or FITNESS FOR A PARTICULAR PURPOSE. See the
// GNU Lesser General Public License for more details.

// You should have received a copy of the GNU Lesser General Public License
// along with the MegaSpace source code. If not, see <http://www.gnu.org/licenses/>.

package eventmanager

import (
	"bufio"
	"errors"
	"fmt"
	"os"
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
	datafile  string         // Data file for tracking processed block (temporary)
	processor BlockProcessor // Processor consumes new block pushed from connector
	connector conn.Connector // connector
	cron      *cron.Cron     // scheduler
	running   bool           // if event manager is running
	logger    *zap.Logger    // logging
	lock      sync.Mutex     // mutex lock
}

// New constructs an instance of eventManager
func New(conn conn.Connector, processor BlockProcessor, datafile string, logger *zap.Logger) *EventManager {
	return &EventManager{
		datafile:  datafile,
		processor: processor,
		connector: conn,
		logger:    logger,
	}
}

// Start would start an EventManager loop
func (e *EventManager) Start() error {
	e.lock.Lock()
	defer e.lock.Unlock()
	if e.running {
		return errors.New("event manage is already running")
	}

	previousBlock, err := e.findPreviousBlock()
	if err != nil {
		e.logger.Error("Failed to find previous block", zap.Error(err))
	}

	err = e.connector.Start()
	if err != nil {
		return err
	}

	err = e.processor.Start()
	if err != nil {
		return err
	}

	blocks := make(chan types.Block, blockBuffer)
	_, err = e.connector.SubscribeBlock(previousBlock, blocks)
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

	err = e.processor.Stop()
	if err != nil {
		e.logger.Error("processor stop with err", zap.Error(err))
	}

	e.running = false
	return nil
}

func (e *EventManager) processBlock(block types.Block) {
	e.logger.Info("Received new block", zap.Stringer("height", block.Height()))
	err := e.processor.ProcessBlock(block)
	if err != nil {
		e.logger.Error("processor error", zap.Error(err))
	}
}

func (e *EventManager) findPreviousBlock() (types.Block, error) {
	file, err := os.Open(e.datafile)
	if err != nil {
		return nil, err
	}
	defer file.Close()

	var last string
	scanner := bufio.NewScanner(file)
	for scanner.Scan() {
		last = scanner.Text()
	}
	if scanner.Err() != nil {
		return nil, scanner.Err()
	}
	if last == "" {
		return nil, nil
	}
	return types.UnmarshalBlock([]byte(last))
}

// QueryAccountBalance queries the chain for the current balance of the given address.
// Returns a channel into which the result will be pushed once retrieved.
func (e *EventManager) QueryAccountBalance(addr string) (<-chan float64, error) {
	return e.connector.QueryAccountBalance(addr)
}

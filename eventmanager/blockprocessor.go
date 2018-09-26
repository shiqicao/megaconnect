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
	"encoding/json"
	"fmt"
	"os"

	"go.uber.org/zap"

	"github.com/megaspacelab/eventmanager/types"
)

type BlockProcessor interface {
	ProcessBlock(block types.Block) error
	Start() error
	Stop() error
}

// BlockBodyProcessor evaluates each txn push it to consumer
type BlockLogger struct {
	datafile *os.File
	logger   *zap.Logger
}

func NewBlockBodyProcessor(datafile string, logger *zap.Logger) (*BlockLogger, error) {
	file, err := os.OpenFile(datafile, os.O_APPEND|os.O_WRONLY|os.O_CREATE, 0600)
	if err != nil {
		return nil, err
	}
	return &BlockLogger{
		datafile: file,
		logger:   logger,
	}, nil
}

func (b *BlockLogger) Start() error {
	return nil
}

func (b *BlockLogger) Stop() error {
	b.datafile.Close()
	return nil
}

func (b *BlockLogger) writeToFile(block types.Block) error {
	writer := bufio.NewWriter(b.datafile)
	json, err := json.Marshal(block)
	if err != nil {
		return nil
	}
	fmt.Fprintln(writer, string(json))
	return writer.Flush()
}

func (b *BlockLogger) ProcessBlock(block types.Block) error {
	b.writeToFile(block)
	return nil
}

// Copyright 2018 @ MegaSpace

// The MegaSpace source code is distributed in the hope that it will be useful,
// but WITHOUT ANY WARRANTY; without even the implied warranty of
// MERCHANTABILITY or FITNESS FOR A PARTICULAR PURPOSE. See the
// GNU Lesser General Public License for more details.

// You should have received a copy of the GNU Lesser General Public License
// along with the MegaSpace source code. If not, see <http://www.gnu.org/licenses/>.

package workflow

import (
	"math/big"

	"github.com/megaspacelab/megaconnect/common"
)

// Env provides execution environment for interpreter and analyzer,
// analyzer does symbol resolving and type checking but not execution.
// TODO: this env is specificity for ChainManager
type Env struct {
	prelude      []*NamespaceDecl
	chain        chain
	currentBlock common.Block
}

type chain interface {
	QueryAccountBalance(addr string, height *big.Int) (*big.Int, error)
}

// NewEnv creates a new Env
func NewEnv(chain chain, currentBlock common.Block) *Env {
	return &Env{
		prelude:      prelude,
		chain:        chain,
		currentBlock: currentBlock,
	}
}

// CurrentBlock returns the block interpreter anchors for a workflow
func (e *Env) CurrentBlock() common.Block {
	return e.currentBlock
}

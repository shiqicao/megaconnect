// Copyright 2018 @ MegaSpace

// The MegaSpace source code is distributed in the hope that it will be useful,
// but WITHOUT ANY WARRANTY; without even the implied warranty of
// MERCHANTABILITY or FITNESS FOR A PARTICULAR PURPOSE. See the
// GNU Lesser General Public License for more details.

// You should have received a copy of the GNU Lesser General Public License
// along with the MegaSpace source code. If not, see <http://www.gnu.org/licenses/>.
package types

import (
	"math/big"

	"github.com/megaspacelab/eventmanager/common"
)

type Block interface {
	Hash() common.Hash
	ParentHash() common.Hash
	Height() *big.Int
	Transactions() []Transaction
}

func NewBlock(
	hash common.Hash,
	parentHash common.Hash,
	height *big.Int,
	transactions []Transaction,
) Block {
	return &block{
		hash:         hash,
		parentHash:   parentHash,
		height:       height,
		transactions: transactions,
	}
}

type block struct {
	hash         common.Hash
	parentHash   common.Hash
	height       *big.Int
	transactions []Transaction
}

func (b *block) Hash() common.Hash           { return b.hash }
func (b *block) ParentHash() common.Hash     { return b.parentHash }
func (b *block) Height() *big.Int            { return b.height }
func (b *block) Transactions() []Transaction { return b.transactions }

// Copyright 2018 @ MegaSpace

// The MegaSpace source code is distributed in the hope that it will be useful,
// but WITHOUT ANY WARRANTY; without even the implied warranty of
// MERCHANTABILITY or FITNESS FOR A PARTICULAR PURPOSE. See the
// GNU Lesser General Public License for more details.

// You should have received a copy of the GNU Lesser General Public License
// along with the MegaSpace source code. If not, see <http://www.gnu.org/licenses/>.
package common

import (
	"encoding/json"
	"math/big"
)

type Block interface {
	Hash() Hash
	ParentHash() Hash
	Height() *big.Int
	Transactions() []Transaction
}

func UnmarshalBlock(data []byte) (Block, error) {
	b := &struct {
		Hash         Hash
		ParentHash   Hash
		Height       *big.Int
		Transactions []Transaction
	}{}
	err := json.Unmarshal(data, b)
	if err != nil {
		return nil, err
	}
	return &block{
		hash:         b.Hash,
		parentHash:   b.ParentHash,
		height:       b.Height,
		transactions: b.Transactions,
	}, nil
}

func NewBlock(
	hash Hash,
	parentHash Hash,
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
	hash         Hash
	parentHash   Hash
	height       *big.Int
	transactions []Transaction
}

func (b *block) Hash() Hash                  { return b.hash }
func (b *block) ParentHash() Hash            { return b.parentHash }
func (b *block) Height() *big.Int            { return new(big.Int).Set(b.height) }
func (b *block) Transactions() []Transaction { return b.transactions }

func (b *block) MarshalJSON() ([]byte, error) {
	return json.Marshal(
		&struct {
			Hash         Hash
			ParentHash   Hash
			Height       *big.Int
			Transactions []Transaction
		}{
			Hash:         b.hash,
			ParentHash:   b.parentHash,
			Height:       b.height,
			Transactions: b.transactions,
		})
}

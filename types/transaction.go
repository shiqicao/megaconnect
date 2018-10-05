// Copyright 2018 @ MegaSpace

// The MegaSpace source code is distributed in the hope that it will be useful,
// but WITHOUT ANY WARRANTY; without even the implied warranty of
// MERCHANTABILITY or FITNESS FOR A PARTICULAR PURPOSE. See the
// GNU Lesser General Public License for more details.

// You should have received a copy of the GNU Lesser General Public License
// along with the MegaSpace source code. If not, see <http://www.gnu.org/licenses/>.

package types

import "github.com/megaspacelab/eventmanager/common"

// Transaction stores information about a specific blockchain transaction.
type Transaction interface {
	Hash() common.Hash
	From() []common.Address
	To() []common.Address
}

// NewTransaction creates a new Transaction.
func NewTransaction(
	hash common.Hash,
	from []common.Address,
	to []common.Address,
) Transaction {
	return &transaction{
		TxHash:    hash,
		FromAddrs: from,
		ToAddrs:   to,
	}
}

type transaction struct {
	TxHash    common.Hash      `json:"hash"`
	FromAddrs []common.Address `json:"from"`
	ToAddrs   []common.Address `json:"to"`
}

func (t *transaction) Hash() common.Hash      { return t.TxHash }
func (t *transaction) From() []common.Address { return t.FromAddrs }
func (t *transaction) To() []common.Address   { return t.ToAddrs }

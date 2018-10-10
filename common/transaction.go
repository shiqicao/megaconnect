// Copyright 2018 @ MegaSpace

// The MegaSpace source code is distributed in the hope that it will be useful,
// but WITHOUT ANY WARRANTY; without even the implied warranty of
// MERCHANTABILITY or FITNESS FOR A PARTICULAR PURPOSE. See the
// GNU Lesser General Public License for more details.

// You should have received a copy of the GNU Lesser General Public License
// along with the MegaSpace source code. If not, see <http://www.gnu.org/licenses/>.

package common

// Transaction stores information about a specific blockchain transaction.
type Transaction interface {
	Hash() Hash
	From() []Address
	To() []Address
}

// NewTransaction creates a new Transaction.
func NewTransaction(
	hash Hash,
	from []Address,
	to []Address,
) Transaction {
	return &transaction{
		TxHash:    hash,
		FromAddrs: from,
		ToAddrs:   to,
	}
}

type transaction struct {
	TxHash    Hash      `json:"hash"`
	FromAddrs []Address `json:"from"`
	ToAddrs   []Address `json:"to"`
}

func (t *transaction) Hash() Hash      { return t.TxHash }
func (t *transaction) From() []Address { return t.FromAddrs }
func (t *transaction) To() []Address   { return t.ToAddrs }

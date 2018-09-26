// Copyright 2018 @ MegaSpace

// The MegaSpace source code is distributed in the hope that it will be useful,
// but WITHOUT ANY WARRANTY; without even the implied warranty of
// MERCHANTABILITY or FITNESS FOR A PARTICULAR PURPOSE. See the
// GNU Lesser General Public License for more details.

// You should have received a copy of the GNU Lesser General Public License
// along with the MegaSpace source code. If not, see <http://www.gnu.org/licenses/>.

package types

import "github.com/megaspacelab/eventmanager/common"

type Transaction interface {
	Hash() common.Hash
	From() common.Address
	To() common.Address
}

func NewTransaction(
	hash common.Hash,
	from common.Address,
	to common.Address,
) Transaction {
	return &transaction{
		hash: hash,
		from: from,
		to:   to,
	}
}

type transaction struct {
	hash common.Hash
	from common.Address
	to   common.Address
}

func (t *transaction) Hash() common.Hash    { return t.hash }
func (t *transaction) From() common.Address { return t.from }
func (t *transaction) To() common.Address   { return t.to }

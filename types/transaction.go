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
}

func NewTransaction(
	hash common.Hash,
) Transaction {
	return &transaction{
		hash: hash,
	}
}

type transaction struct {
	hash common.Hash
}

func (t *transaction) Hash() common.Hash { return t.hash }

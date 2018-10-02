// Copyright 2018 @ MegaSpace

// The MegaSpace source code is distributed in the hope that it will be useful,
// but WITHOUT ANY WARRANTY; without even the implied warranty of
// MERCHANTABILITY or FITNESS FOR A PARTICULAR PURPOSE. See the
// GNU Lesser General Public License for more details.

// You should have received a copy of the GNU Lesser General Public License
// along with the MegaSpace source code. If not, see <http://www.gnu.org/licenses/>.

package workflow

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestNewObjConst(t *testing.T) {
	objType := NewObjType(map[string]Type{
		"a": BoolType,
	})

	value := NewObjConst(
		map[string]Const{
			"a": TrueConst,
		},
	)
	assert.True(t, value.Type().Equal(objType))
	assert.True(t, value.Value()["a"].Equal(TrueConst))
}

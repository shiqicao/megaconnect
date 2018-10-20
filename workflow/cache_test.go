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

func TestAPICache(t *testing.T) {
	// ambiguity test
	x := NewFuncDecl("A", Params{}, IntType, nil)
	x.parent = NewNamespaceDecl("A")

	y := NewFuncDecl("AA", Params{}, IntType, nil)

	xk, err := buildApiCacheKey(x, []Const{})
	assert.NoError(t, err)

	yk, err := buildApiCacheKey(y, []Const{})
	assert.NoError(t, err)

	assert.NotEqual(t, xk, yk)

	// with/without namespace test
	x = NewFuncDecl("A", Params{}, IntType, nil)
	x.parent = NewNamespaceDecl("A")

	y = NewFuncDecl("A", Params{}, IntType, nil)

	xk, err = buildApiCacheKey(x, []Const{})
	assert.NoError(t, err)

	yk, err = buildApiCacheKey(y, []Const{})
	assert.NoError(t, err)

	assert.NotEqual(t, xk, yk)

	// param test -- different param
	x = NewFuncDecl("A", Params{}, IntType, nil)

	xk1, err := buildApiCacheKey(x, []Const{TrueConst})
	assert.NoError(t, err)

	xk2, err := buildApiCacheKey(x, []Const{FalseConst})
	assert.NoError(t, err)

	assert.NotEqual(t, xk1, xk2)

	// param test -- different length
	x = NewFuncDecl("A", Params{}, IntType, nil)

	xk1, err = buildApiCacheKey(x, []Const{TrueConst, TrueConst})
	assert.NoError(t, err)

	xk2, err = buildApiCacheKey(y, []Const{TrueConst})
	assert.NoError(t, err)

	assert.NotEqual(t, xk1, xk2)

	// param test -- obj
	x = NewFuncDecl("A", Params{}, IntType, nil)

	xk1, err = buildApiCacheKey(x, []Const{NewObjConst(map[string]Const{"a": TrueConst, "b": FalseConst})})
	assert.NoError(t, err)

	xk2, err = buildApiCacheKey(x, []Const{NewObjConst(map[string]Const{"b": FalseConst, "a": TrueConst})})
	assert.NoError(t, err)

	assert.Equal(t, xk1, xk2)

	// param test -- obj nest
	x = NewFuncDecl("A", Params{}, IntType, nil)

	xk1, err = buildApiCacheKey(x, []Const{NewObjConst(map[string]Const{"a": TrueConst, "b": NewObjConst(map[string]Const{"a": FalseConst, "b": TrueConst})})})
	assert.NoError(t, err)

	xk2, err = buildApiCacheKey(x, []Const{NewObjConst(map[string]Const{"b": NewObjConst(map[string]Const{"b": TrueConst, "a": FalseConst}), "a": TrueConst})})
	assert.NoError(t, err)

	assert.Equal(t, xk1, xk2)
}

// Copyright 2018 @ MegaSpace

// The MegaSpace source code is distributed in the hope that it will be useful,
// but WITHOUT ANY WARRANTY; without even the implied warranty of
// MERCHANTABILITY or FITNESS FOR A PARTICULAR PURPOSE. See the
// GNU Lesser General Public License for more details.

// You should have received a copy of the GNU Lesser General Public License
// along with the MegaSpace source code. If not, see <http://www.gnu.org/licenses/>.

package workflow

import (
	"sort"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestActionDeclTriggerEvents(t *testing.T) {
	check := func(trigger EventExpr, expected []string) {
		a := NewActionDecl(ID("a"), trigger, Stmts{})
		events := a.TriggerEvents()
		sort.StringSlice(events).Sort()
		sort.StringSlice(expected).Sort()
		assert.Equal(t, expected, events)
	}

	check(EV("a"), []string{"a"})
	check(EAND(EV("a"), EV("a")), []string{"a"})
	check(EAND(EV("a"), EV("b")), []string{"a", "b"})
	check(EAND(EV("a"), EOR(EV("b"), EV("c"))), []string{"a", "b", "c"})
	check(EAND(EV("a"), EOR(EV("b"), EV("a"))), []string{"a", "b"})
}

func TestNamespaceAddConflict(t *testing.T) {
	n := NS("a")
	assert.NoError(t, n.AddNamespaces(NS("b")))
	assert.Error(t, n.AddNamespaces(NS("b")))

	n = NS("a")
	assert.NoError(t, n.AddFuncs(FD("b", Params{}, IntType, nil)))
	assert.Error(t, n.AddFuncs(FD("b", Params{}, IntType, nil)))

	n = NS("a")
	assert.NoError(t, n.AddFuncs(FD("b", Params{}, IntType, nil)))
	assert.NoError(t, n.AddNamespaces(NS("b")))

	n = NS("a")
	assert.NoError(t, n.AddNamespaces(NS("b")))
	assert.NoError(t, n.AddFuncs(FD("b", Params{}, IntType, nil)))

	n = NS("a")
	assert.Error(t, n.AddNamespaces(NS("b"), NS("b")))

	n = NS("a")
	assert.Error(t, n.AddFuncs(
		FD("b", Params{}, IntType, nil),
		FD("b", Params{}, IntType, nil),
	))
}

func TestNamespaceCopy(t *testing.T) {
	a := NS("a")
	a.AddNamespaces(NS("b"))
	b := a.Copy()
	assert.True(t, a.Equal(b))
	assert.True(t, a != b)
	assert.True(t, a.namespaces[0] != b.namespaces[0])
	assert.True(t, b.namespaces[0].Parent() == b)
	assert.True(t, a.namespaces[0].Parent() == a)

	a = NS("a")
	a.AddFuncs(FD("b", Params{}, IntType, nil))
	b = a.Copy()
	assert.True(t, a.Equal(b))
	assert.True(t, a != b)
	assert.True(t, a.funs[0] != b.funs[0])
	assert.True(t, b.funs[0].Parent() == b)
	assert.True(t, a.funs[0].Parent() == a)
}

func TestFuncDeclCopy(t *testing.T) {
	f1 := FD("a", Params{}, IntType, nil)
	f2 := f1.Copy()
	assert.True(t, f1.Equal(f2))
	assert.True(t, f1 != f2)
}

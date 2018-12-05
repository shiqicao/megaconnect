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

func TestResolveAction(t *testing.T) {
	noNamespace := func() *FuncCall { return NewFuncCall(nil, ID("F")) }
	withNamespace := func(ns *Id) *FuncCall { return NewFuncCall(NamespacePrefix{ns}, ID("F")) }
	a := func(f *FuncCall) *ActionDecl {
		return NewActionDecl(ID("T"), EV("a"), Stmts{FIRE("e", NewObjLit(VD("x", f)))})
	}

	r := NewResolver(nil)
	assertErrs(t, 1, r.resolveAction(a(noNamespace())))
	assertErrs(t, 1, r.resolveAction(a(withNamespace(ID("A")))))

	ns := NewNamespaceDecl("A")
	err := ns.AddFuncs(
		NewFuncDecl(
			"F",
			Params{},
			BoolType,
			nil,
		),
	)
	assert.NoError(t, err)
	lib := []*NamespaceDecl{ns}

	r = NewResolver(lib)
	assertErrs(t, 1, r.resolveAction(a(noNamespace())))
	f := withNamespace(ID("A"))
	assertNoErrs(t, r.resolveAction(a(f)))
	assert.NotNil(t, f.decl)

	f = withNamespace(ID("A"))
	assertNoErrs(t, r.resolveAction(a(f)))
	assert.NotNil(t, f.decl)
}

func TestResolveFuncCall(t *testing.T) {
	withNamespace := func(nss ...*Id) *FuncCall { return NewFuncCall(NamespacePrefix(nss), ID("F")) }
	A := NewNamespaceDecl("A")
	err := A.AddFuncs(
		NewFuncDecl(
			"F",
			Params{},
			BoolType,
			nil,
		),
	)
	assert.NoError(t, err)

	lib := []*NamespaceDecl{A}
	r := NewResolver(lib)
	f := withNamespace(ID("A"))
	assert.NoError(t, r.resolveFuncCall(f))
	assert.NotNil(t, f.decl)

	f = withNamespace(ID("B"), ID("A"))
	assert.Error(t, r.resolveFuncCall(f))

	B := NewNamespaceDecl("B")
	err = B.AddNamespaces(A)
	assert.NoError(t, err)
	lib = []*NamespaceDecl{B}
	r = NewResolver(lib)
	f = withNamespace(ID("A"))
	assert.Error(t, r.resolveFuncCall(f))

	f = withNamespace(ID("B"), ID("A"))
	assert.NoError(t, r.resolveFuncCall(f))
	assert.NotNil(t, f.decl)
}

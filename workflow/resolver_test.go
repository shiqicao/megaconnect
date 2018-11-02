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
	noNamespace := func() *FuncCall { return NewFuncCall(nil, "F") }
	withNamespace := func(ns string) *FuncCall { return NewFuncCall(NamespacePrefix{ns}, "F") }
	a := func(f *FuncCall) *ActionDecl {
		return NewActionDecl("T", EV("a"), Stmts{FIRE("e", NewObjLit(VarDecls{"x": f}))})
	}

	r := NewResolver(nil, nil)
	assert.Error(t, r.resolveAction(a(noNamespace())))
	assert.Error(t, r.resolveAction(a(withNamespace("A"))))

	lib := []*NamespaceDecl{
		NewNamespaceDecl("A").AddFunc(
			NewFuncDecl(
				"F",
				Params{},
				BoolType,
				nil,
			),
		),
	}

	r = NewResolver(lib, nil)
	assert.Error(t, r.resolveAction(a(noNamespace())))
	f := withNamespace("A")
	assert.NoError(t, r.resolveAction(a(f)))
	assert.NotNil(t, f.decl)

	r = NewResolver(lib, NamespacePrefix{"A"})
	f = noNamespace()
	assert.NoError(t, r.resolveAction(a(f)))
	assert.NotNil(t, f.decl)

	f = withNamespace("A")
	assert.NoError(t, r.resolveAction(a(f)))
	assert.NotNil(t, f.decl)
}

func TestResolveFuncCall(t *testing.T) {
	withNamespace := func(nss ...string) *FuncCall { return NewFuncCall(NamespacePrefix(nss), "F") }
	A := NewNamespaceDecl("A").AddFunc(
		NewFuncDecl(
			"F",
			Params{},
			BoolType,
			nil,
		),
	)

	lib := []*NamespaceDecl{A}
	r := NewResolver(lib, nil)
	f := withNamespace("A")
	assert.NoError(t, r.resolveFuncCall(f))
	assert.NotNil(t, f.decl)

	f = withNamespace("B", "A")
	assert.Error(t, r.resolveFuncCall(f))

	lib = []*NamespaceDecl{NewNamespaceDecl("B").AddChild(A)}
	r = NewResolver(lib, nil)
	f = withNamespace("A")
	assert.Error(t, r.resolveFuncCall(f))

	f = withNamespace("B", "A")
	assert.NoError(t, r.resolveFuncCall(f))
	assert.NotNil(t, f.decl)
}

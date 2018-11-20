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

	r := NewResolver(nil, nil)
	assertErrs(t, 1, r.resolveAction(a(noNamespace())))
	assertErrs(t, 1, r.resolveAction(a(withNamespace(ID("A")))))

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
	assertErrs(t, 1, r.resolveAction(a(noNamespace())))
	f := withNamespace(ID("A"))
	assertNoErrs(t, r.resolveAction(a(f)))
	assert.NotNil(t, f.decl)

	r = NewResolver(lib, NamespacePrefix{ID("A")})
	f = noNamespace()
	assertNoErrs(t, r.resolveAction(a(f)))
	assert.NotNil(t, f.decl)

	f = withNamespace(ID("A"))
	assertNoErrs(t, r.resolveAction(a(f)))
	assert.NotNil(t, f.decl)
}

func TestResolveFuncCall(t *testing.T) {
	withNamespace := func(nss ...*Id) *FuncCall { return NewFuncCall(NamespacePrefix(nss), ID("F")) }
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
	f := withNamespace(ID("A"))
	assert.NoError(t, r.resolveFuncCall(f))
	assert.NotNil(t, f.decl)

	f = withNamespace(ID("B"), ID("A"))
	assert.Error(t, r.resolveFuncCall(f))

	lib = []*NamespaceDecl{NewNamespaceDecl("B").AddChild(A)}
	r = NewResolver(lib, nil)
	f = withNamespace(ID("A"))
	assert.Error(t, r.resolveFuncCall(f))

	f = withNamespace(ID("B"), ID("A"))
	assert.NoError(t, r.resolveFuncCall(f))
	assert.NotNil(t, f.decl)
}

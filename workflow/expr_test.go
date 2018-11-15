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

func TestRatNorm(t *testing.T) {
	check := func(x *RatConst, y *RatConst) bool {
		return x.Value().Num().Cmp(y.Value().Num()) == 0 && x.Value().Denom().Cmp(y.Value().Denom()) == 0
	}
	check(R(0, 1), R(0, 2))
	check(R(0, 1), R(0, -2))
	check(R(2, 1), R(-4, -2))
	check(R(-1, 2), R(1, -2))
}

func TestExprEquality(t *testing.T) {
	// const - boolean
	assert.True(t, TrueConst.Equal(TrueConst))
	assert.True(t, FalseConst.Equal(FalseConst))
	assert.False(t, FalseConst.Equal(TrueConst))
	assert.False(t, FalseConst.Equal(NewIntConstFromI64(0)))

	// const - int
	assert.True(t, NewIntConstFromI64(0).Equal(NewIntConstFromI64(0)))
	assert.True(t, NewIntConstFromI64(1).Equal(NewIntConstFromI64(1)))
	assert.True(t, NewIntConstFromI64(-1).Equal(NewIntConstFromI64(-1)))
	assert.False(t, NewIntConstFromI64(1).Equal(NewIntConstFromI64(0)))
	assert.False(t, NewIntConstFromI64(1).Equal(NewIntConstFromI64(-1)))
	assert.False(t, NewIntConstFromI64(0).Equal(FalseConst))

	// const - str
	assert.True(t, NewStrConst("").Equal(NewStrConst("")))
	assert.True(t, NewStrConst("a").Equal(NewStrConst("a")))
	assert.False(t, NewStrConst("a").Equal(NewStrConst("b")))
	assert.False(t, NewStrConst("a").Equal(NewIntConstFromI64(1)))

	// const - obj
	assert.True(t, NewObjConst(ObjFields{}).Equal(
		NewObjConst(ObjFields{}),
	))

	assert.False(t, NewObjConst(ObjFields{}).Equal(
		NewObjConst(ObjFields{
			"a": TrueConst,
		}),
	))

	assert.True(t, NewObjConst(ObjFields{
		"a": TrueConst,
	}).Equal(
		NewObjConst(ObjFields{
			"a": TrueConst,
		}),
	))

	assert.False(t, NewObjConst(ObjFields{
		"a": TrueConst,
	}).Equal(
		NewObjConst(ObjFields{
			"a": FalseConst,
		}),
	))

	assert.False(t, NewObjConst(ObjFields{
		"a": TrueConst,
	}).Equal(
		NewObjConst(ObjFields{
			"a": NewStrConst("a"),
		}),
	))

	assert.True(t, NewObjConst(ObjFields{
		"a": TrueConst,
		"b": FalseConst,
	}).Equal(
		NewObjConst(ObjFields{
			"a": TrueConst,
			"b": FalseConst,
		}),
	))

	assert.True(t, NewObjConst(ObjFields{
		"a": TrueConst,
		"b": FalseConst,
	}).Equal(
		NewObjConst(ObjFields{
			"b": FalseConst,
			"a": TrueConst,
		}),
	))

	assert.True(t, NewObjConst(ObjFields{
		"a": TrueConst,
		"b": FalseConst,
	}).Equal(
		NewObjConst(ObjFields{
			"b": FalseConst,
			"a": TrueConst,
		}),
	))

	assert.False(t, NewObjConst(ObjFields{
		"a": TrueConst,
	}).Equal(
		NewObjConst(ObjFields{
			"b": FalseConst,
			"a": TrueConst,
		}),
	))

	// FuncCall
	assert.True(t, NewFuncCall(NamespacePrefix{}, ID("F")).Equal(NewFuncCall(NamespacePrefix{}, ID("F"))))
	assert.True(t, NewFuncCall(NamespacePrefix{ID("B")}, ID("F")).Equal(NewFuncCall(NamespacePrefix{ID("B")}, ID("F"))))
	assert.True(t, NewFuncCall(NamespacePrefix{ID("B"), ID("B2")}, ID("F")).Equal(NewFuncCall(NamespacePrefix{ID("B"), ID("B2")}, ID("F"))))
	assert.True(t, NewFuncCall(NamespacePrefix{ID("B")}, ID("F"), TrueConst).Equal(NewFuncCall(NamespacePrefix{ID("B")}, ID("F"), TrueConst)))

	assert.False(t, NewFuncCall(nil, ID("F")).Equal(NewFuncCall(nil, ID("G"))))
	assert.True(t, NewFuncCall(NamespacePrefix{ID("B")}, ID("F")).Equal(NewFuncCall(NamespacePrefix{ID("B")}, ID("F"))))
	assert.False(t, NewFuncCall(NamespacePrefix{ID("B"), ID("B2")}, ID("F")).Equal(NewFuncCall(NamespacePrefix{ID("B"), ID("A2")}, ID("F"))))
	assert.False(t, NewFuncCall(NamespacePrefix{ID("B")}, ID("F"), TrueConst).Equal(NewFuncCall(NamespacePrefix{ID("B")}, ID("F"), FalseConst)))
	assert.False(t, NewFuncCall(NamespacePrefix{ID("B")}, ID("F"), TrueConst).Equal(NewFuncCall(NamespacePrefix{ID("B")}, ID("F"), TrueConst, TrueConst)))

	// ObjAccessor
	assert.True(t, NewObjAccessor(NewObjConst(ObjFields{}), "A").Equal(
		NewObjAccessor(NewObjConst(ObjFields{}), "A"),
	))
	assert.False(t, NewObjAccessor(NewObjConst(ObjFields{}), "A").Equal(
		NewObjAccessor(NewObjConst(ObjFields{}), "B"),
	))
	assert.False(t, NewObjAccessor(NewObjConst(ObjFields{"A": TrueConst}), "A").Equal(
		NewObjAccessor(NewObjConst(ObjFields{}), "B"),
	))
	assert.False(t, NewObjAccessor(NewObjConst(ObjFields{"A": TrueConst}), "A").Equal(
		NewObjAccessor(NewObjConst(ObjFields{"B": TrueConst}), "A"),
	))

	// UniOp
	assert.True(t, NewUniOp(NotOp, TrueConst).Equal(NewUniOp(NotOp, TrueConst)))
	assert.False(t, NewUniOp(NotOp, FalseConst).Equal(NewUniOp(NotOp, TrueConst)))
	assert.False(t, NewUniOp(NotOp, NewBinOp(NotEqualOp, TrueConst, FalseConst)).Equal(NewUniOp(NotOp, TrueConst)))
	assert.False(t, NewUniOp(NotOp, FalseConst).Equal(NewBinOp(NotEqualOp, TrueConst, TrueConst)))

	// BinOp
	assert.True(t, NewBinOp(NotEqualOp, TrueConst, TrueConst).Equal(NewBinOp(NotEqualOp, TrueConst, TrueConst)))
	assert.False(t, NewBinOp(NotEqualOp, TrueConst, TrueConst).Equal(NewBinOp(EqualOp, TrueConst, TrueConst)))
	assert.False(t, NewBinOp(EqualOp, FalseConst, TrueConst).Equal(NewBinOp(EqualOp, TrueConst, TrueConst)))
	assert.False(t, NewBinOp(EqualOp, TrueConst, FalseConst).Equal(NewBinOp(EqualOp, TrueConst, TrueConst)))
}

func TestCopy(t *testing.T) {
	ns := NamespacePrefix(nil).Copy()
	assert.Nil(t, ns)
	ns = NamespacePrefix{ID("a"), ID("a")}
	ns_ := ns.Copy()
	assert.Equal(t, ns, ns_)
	ns[0] = ID("b")
	assert.True(t, ID("a").Equal(ns_[0]))

	args := Args(nil).Copy()
	assert.Nil(t, args)
	args = Args{TrueConst, FalseConst}
	args_ := args.Copy()
	assert.Equal(t, args, args_)
	args[0] = FalseConst
	assert.True(t, args_[0].Equal(TrueConst))

	obf := ObjFields(nil).Copy()
	assert.Nil(t, obf)
	obf = ObjFields{"a": TrueConst, "b": FalseConst}
	obf_ := obf.Copy()
	assert.Equal(t, obf, obf_)
	obf["a"] = FalseConst
	assert.True(t, obf_["a"].Equal(TrueConst))

	pm := Params(nil).Copy()
	assert.Nil(t, pm)
	pm = Params{NewParamDecl("a", BoolType), NewParamDecl("b", IntType)}
	pm_ := pm.Copy()
	assert.Equal(t, pm, pm_)
	pm[0] = NewParamDecl("c", StrType)
	assert.Equal(t, "a", pm_[0].Name())

	fd := FuncDecls(nil).Copy()
	assert.Nil(t, fd)
	fd = FuncDecls{NewFuncDecl("a", nil, IntType, nil), NewFuncDecl("b", nil, StrType, nil)}
	fd_ := fd.Copy()
	assert.Equal(t, fd, fd_)
	fd[0] = NewFuncDecl("c", nil, BoolType, nil)
	assert.Equal(t, "a", fd_[0].Name())
}

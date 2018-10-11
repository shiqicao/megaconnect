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

func TestEvalExpr(t *testing.T) {
	assertExpEval(t, TrueConst, TrueConst)
	assertExpEval(t, FalseConst, FalseConst)
}

func TestEqualOp(t *testing.T) {
	equalOpTest(t, true)
	equalOpTest(t, false)
}

func equalOpTest(t *testing.T, isEqual bool) {
	op := EqualOp
	if !isEqual {
		op = NotEqualOp
	}

	positive := GetBoolConst(isEqual)
	negative := GetBoolConst(!isEqual)

	// Equal Op
	assertExpEval(t, negative, NewBinOp(
		op,
		FalseConst,
		TrueConst,
	))
	assertExpEval(t, positive, NewBinOp(
		op,
		FalseConst,
		FalseConst,
	))

	assertExpEval(t, positive, NewBinOp(
		op,
		NewStrConst(""),
		NewStrConst(""),
	))
	assertExpEval(t, positive, NewBinOp(
		op,
		NewStrConst("a"),
		NewStrConst("a"),
	))
	assertExpEval(t, negative, NewBinOp(
		op,
		NewStrConst("b"),
		NewStrConst("a"),
	))

	assertExpEval(t, positive, NewBinOp(
		op,
		NewIntConstFromI64(1),
		NewIntConstFromI64(1),
	))
	assertExpEval(t, negative, NewBinOp(
		op,
		NewIntConstFromI64(0),
		NewIntConstFromI64(1),
	))

	assertExpEval(t, positive, NewBinOp(
		op,
		NewObjConst(map[string]Const{
			"a": TrueConst,
		}),
		NewObjConst(map[string]Const{
			"a": TrueConst,
		}),
	))
	assertExpEval(t, positive, NewBinOp(
		op,
		NewObjConst(map[string]Const{
			"a": TrueConst,
			"b": NewObjConst(map[string]Const{
				"a": TrueConst,
			}),
		}),
		NewObjConst(map[string]Const{
			"a": TrueConst,
			"b": NewObjConst(map[string]Const{
				"a": TrueConst,
			}),
		}),
	))
	assertExpEval(t, negative, NewBinOp(
		op,
		NewObjConst(map[string]Const{
			"a": TrueConst,
			"b": NewObjConst(map[string]Const{
				"a": FalseConst,
			}),
		}),
		NewObjConst(map[string]Const{
			"a": TrueConst,
			"b": NewObjConst(map[string]Const{
				"a": TrueConst,
			}),
		}),
	))
	assertExpEval(t, negative, NewBinOp(
		op,
		NewObjConst(map[string]Const{
			"a": FalseConst,
		}),
		NewObjConst(map[string]Const{
			"a": TrueConst,
		}),
	))

	// Equal Err
	assertExpEvalErr(t, NewBinOp(
		op,
		NewStrConst("b"),
		NewIntConstFromI64(1),
	))
	assertExpEvalErr(t, NewBinOp(
		op,
		NewObjConst(map[string]Const{
			"a": NewIntConstFromI64(1),
		}),
		NewObjConst(map[string]Const{
			"a": TrueConst,
		}),
	))
	assertExpEvalErr(t, NewBinOp(
		op,
		NewObjConst(map[string]Const{
			"b": TrueConst,
		}),
		NewObjConst(map[string]Const{
			"a": TrueConst,
		}),
	))
	assertExpEvalErr(t, NewBinOp(
		op,
		NewObjConst(map[string]Const{
			"a": TrueConst,
			"b": TrueConst,
		}),
		NewObjConst(map[string]Const{
			"a": TrueConst,
		}),
	))
	assertExpEvalErr(t, NewBinOp(
		op,
		NewObjConst(map[string]Const{
			"a": TrueConst,
			"b": NewObjConst(map[string]Const{
				"b": FalseConst,
			}),
		}),
		NewObjConst(map[string]Const{
			"a": TrueConst,
			"b": NewObjConst(map[string]Const{
				"a": TrueConst,
			}),
		}),
	))
}

func TestInequality(t *testing.T) {
	// Test "<"
	ordTest(t, true, false)
	// Test "<="
	ordTest(t, true, true)
	// Test ">"
	ordTest(t, false, false)
	// Test ">="
	ordTest(t, false, true)

}

func TestObjAccessor(t *testing.T) {
	subobj := NewObjConst(map[string]Const{
		"ba": FalseConst,
	})
	obj := NewObjConst(map[string]Const{
		"a": TrueConst,
		"b": subobj,
	})
	accessor := NewObjAccessor(obj, "a")
	assertExpEval(t, TrueConst, accessor)

	accessor = NewObjAccessor(obj, "b")
	assertExpEval(t, subobj, accessor)

	accessor = NewObjAccessor(NewObjAccessor(obj, "b"), "ba")
	assertExpEval(t, FalseConst, accessor)

	accessor = NewObjAccessor(obj, "c")
	assertExpEvalErr(t, accessor)

	accessor = NewObjAccessor(TrueConst, "c")
	assertExpEvalErr(t, accessor)
}

func TestSymbolResolve(t *testing.T) {
	prelude = []*NamespaceDecl{
		&NamespaceDecl{
			name: "TEST",
			funs: []*FuncDecl{
				NewFuncDecl(
					"foo",
					[]*ParamDecl{
						NewParamDecl("bar", StrType),
					},
					NewObjType(map[string]Type{
						"size": IntType,
						"text": StrType,
					}),
					func(env *Env, args map[string]Const) (Const, error) {
						bar := args["bar"]
						barStr := bar.(*StrConst)
						result := NewObjConst(
							map[string]Const{
								"size": NewIntConstFromI64(int64(len(barStr.Value()))),
								"text": NewStrConst(barStr.Value()),
							},
						)
						return result, nil
					},
				),
				NewFuncDecl(
					"bar",
					[]*ParamDecl{},
					StrType,
					func(_ *Env, args map[string]Const) (Const, error) {
						return NewStrConst("aa"), nil
					},
				),
			},
		},
	}

	callFoo := NewFuncCall(NamespacePrefix{"TEST"}, "foo", NewStrConst("bar"))
	result := NewObjConst(
		map[string]Const{
			"size": NewIntConstFromI64(int64(len("bar"))),
			"text": NewStrConst("bar"),
		},
	)

	assertExpEvalWithPrelude(t, result, callFoo, prelude)

	callBar := NewFuncCall(NamespacePrefix{"TEST"}, "bar")
	callFooBar := NewFuncCall(NamespacePrefix{"TEST"}, "foo", callBar)

	result = NewObjConst(
		map[string]Const{
			"size": NewIntConstFromI64(int64(len("aa"))),
			"text": NewStrConst("aa"),
		},
	)

	assertExpEvalWithPrelude(t, result, callFooBar, prelude)

	callFoo = NewFuncCall(NamespacePrefix{"TEST"}, "foo", NewStrConst("bar"))
	assertExpEvalWithPrelude(t, NewStrConst("bar"), NewObjAccessor(
		callFoo,
		"text",
	), prelude)
}

func TestBooleanOps(t *testing.T) {
	assertExpEval(t, TrueConst, NewBinOp(AndOp, TrueConst, TrueConst))
	assertExpEval(t, FalseConst, NewBinOp(AndOp, FalseConst, FalseConst))
	assertExpEval(t, FalseConst, NewBinOp(AndOp, TrueConst, FalseConst))
	assertExpEval(t, FalseConst, NewBinOp(AndOp, FalseConst, TrueConst))

	assertExpEval(t, TrueConst, NewBinOp(OrOp, TrueConst, TrueConst))
	assertExpEval(t, FalseConst, NewBinOp(OrOp, FalseConst, FalseConst))
	assertExpEval(t, TrueConst, NewBinOp(OrOp, TrueConst, FalseConst))
	assertExpEval(t, TrueConst, NewBinOp(OrOp, FalseConst, TrueConst))

	assertExpEval(t, FalseConst, NewUniOp(NotOp, TrueConst))
	assertExpEval(t, TrueConst, NewUniOp(NotOp, FalseConst))
}

func TestIntOp(t *testing.T) {
	assertExpEval(t, NewIntConstFromI64(2), NewBinOp(PlusOp, NewIntConstFromI64(1), NewIntConstFromI64(1)))
	assertExpEval(t, NewIntConstFromI64(0), NewBinOp(PlusOp, NewIntConstFromI64(1), NewIntConstFromI64(-1)))

	assertExpEval(t, NewIntConstFromI64(0), NewBinOp(MinusOp, NewIntConstFromI64(1), NewIntConstFromI64(1)))
	assertExpEval(t, NewIntConstFromI64(-1), NewBinOp(MinusOp, NewIntConstFromI64(1), NewIntConstFromI64(2)))

	assertExpEval(t, NewIntConstFromI64(1), NewBinOp(MultOp, NewIntConstFromI64(1), NewIntConstFromI64(1)))
	assertExpEval(t, NewIntConstFromI64(0), NewBinOp(MultOp, NewIntConstFromI64(1), NewIntConstFromI64(0)))

	assertExpEval(t, NewIntConstFromI64(1), NewBinOp(DivOp, NewIntConstFromI64(1), NewIntConstFromI64(1)))
	assertExpEval(t, NewIntConstFromI64(0), NewBinOp(DivOp, NewIntConstFromI64(0), NewIntConstFromI64(1)))
	assertExpEval(t, NewIntConstFromI64(0), NewBinOp(DivOp, NewIntConstFromI64(1), NewIntConstFromI64(2)))
	assertExpEval(t, NewIntConstFromI64(1), NewBinOp(DivOp, NewIntConstFromI64(3), NewIntConstFromI64(2)))

	assertExpEvalErr(t, NewBinOp(DivOp, NewIntConstFromI64(1), NewIntConstFromI64(0)))
}

func assertExpEval(t *testing.T, expected Const, expr Expr) {
	assertExpEvalWithPrelude(t, expected, expr, nil)
}

func assertExpEvalWithPrelude(t *testing.T, expected Const, expr Expr, prelude []*NamespaceDecl) {
	env := NewEnv(nil, nil)
	if prelude != nil {
		env.prelude = prelude
	}
	i := New(env)
	result, err := i.EvalExpr(expr)
	assert.NoError(t, err)
	assert.NotNil(t, result)
	assert.NotNil(t, result.Type())
	assert.True(t, result.Type().Equal(expected.Type()))
	assert.True(t, result.Equal(expected))
}

func assertExpEvalErr(t *testing.T, expr Expr) {
	i := New(NewEnv(nil, nil))
	result, err := i.EvalExpr(expr)
	assert.Error(t, err)
	assert.Nil(t, result)
}

func ordTest(t *testing.T, lessThan bool, equalTo bool) {
	op := LessThanOp
	if lessThan && equalTo {
		op = LessThanEqualOp
	} else if !lessThan && equalTo {
		op = GreaterThanEqualOp
	} else if !lessThan && !equalTo {
		op = GreaterThanOp
	}
	lt := GetBoolConst(lessThan)
	eq := GetBoolConst(equalTo)
	gt := GetBoolConst(!lessThan)

	assertExpEval(t, lt, NewBinOp(
		op,
		NewIntConstFromI64(0),
		NewIntConstFromI64(1),
	))
	assertExpEval(t, gt, NewBinOp(
		op,
		NewIntConstFromI64(0),
		NewIntConstFromI64(-1),
	))
	assertExpEval(t, eq, NewBinOp(
		op,
		NewIntConstFromI64(0),
		NewIntConstFromI64(0),
	))
}

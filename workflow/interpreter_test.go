// Copyright 2018 @ MegaSpace

// The MegaSpace source code is distributed in the hope that it will be useful,
// but WITHOUT ANY WARRANTY; without even the implied warranty of
// MERCHANTABILITY or FITNESS FOR A PARTICULAR PURPOSE. See the
// GNU Lesser General Public License for more details.

// You should have received a copy of the GNU Lesser General Public License
// along with the MegaSpace source code. If not, see <http://www.gnu.org/licenses/>.

package workflow

import (
	"fmt"
	"testing"

	"go.uber.org/zap"

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
	prelude := []*NamespaceDecl{
		&NamespaceDecl{
			name: "TEST",
			funs: []*FuncDecl{
				NewFuncDecl(
					"foo",
					[]*ParamDecl{
						NewParamDecl("bar", StrType),
					},
					NewObjType(VT("size", IntType).Put("text", StrType)),
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
	ib := newInterpreterBuilder()
	ib.withLibs(prelude).assertExpEval(t, result, callFoo)

	callBar := NewFuncCall(NamespacePrefix{"TEST"}, "bar")
	callFooBar := NewFuncCall(NamespacePrefix{"TEST"}, "foo", callBar)

	result = NewObjConst(
		map[string]Const{
			"size": NewIntConstFromI64(int64(len("aa"))),
			"text": NewStrConst("aa"),
		},
	)

	ib.withLibs(prelude).assertExpEval(t, result, callFooBar)

	callFoo = NewFuncCall(NamespacePrefix{"TEST"}, "foo", NewStrConst("bar"))
	ib.withLibs(prelude).assertExpEval(t, NewStrConst("bar"), NewObjAccessor(
		callFoo,
		"text",
	))
}

type mockEventStore struct {
	events map[string]*ObjConst
}

func (m *mockEventStore) Lookup(name string) *ObjConst {
	for e, o := range m.events {
		if e == name {
			return o
		}
	}
	return nil
}

func TestEvalAction(t *testing.T) {
	ib := newInterpreterBuilder()
	i := ib.withEM(
		&mockEventStore{events: map[string]*ObjConst{"a": NewObjConst(ObjFields{"x": T})}},
	)()

	r, err := i.EvalAction(
		NewActionDecl(ID("B"), EV("a"), Stmts{FIRE("B", NewObjLit(VD("x", TrueConst)))}),
	)
	assert.NoError(t, err)
	assert.Equal(t, 1, len(r))
	fireResult := r[0].(*FireEventResult)
	assert.NotNil(t, fireResult)
	assert.Equal(t, fireResult.eventName, "B")
	assert.True(t, NewObjConst(ObjFields{"x": TrueConst}).Equal(fireResult.payload))
}

func TestEvalEExpr(t *testing.T) {
	em := &mockEventStore{}
	i := newInterpreterBuilder().withEM(em)()
	check := func(expected bool, eexpr EventExpr) {
		r, err := i.evalEventExpr(eexpr)
		assert.NoError(t, err)
		assert.Equal(t, expected, r)
	}

	check(false, EAND(EV("a"), EV("b")))
	check(false, EOR(EV("a"), EV("b")))

	em.events = map[string]*ObjConst{"a": NewObjConst(ObjFields{"x": T})}
	check(false, EAND(EV("a"), EV("b")))
	check(true, EOR(EV("a"), EV("b")))

	em.events = map[string]*ObjConst{"b": NewObjConst(ObjFields{"x": T})}
	check(false, EAND(EV("a"), EV("b")))
	check(true, EOR(EV("a"), EV("b")))

	em.events = map[string]*ObjConst{
		"a": NewObjConst(ObjFields{"x": T}),
		"b": NewObjConst(ObjFields{"x": T}),
	}
	check(true, EAND(EV("a"), EV("b")))
	check(true, EOR(EV("a"), EV("b")))

	em.events = map[string]*ObjConst{"a": NewObjConst(ObjFields{"x": T})}
	check(false, EAND(
		EOR(EV("a"), EV("b")),
		EAND(EV("b"), EV("a")),
	))
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

func TestProps(t *testing.T) {
	es := &mockEventStore{events: nil}
	ib := newInterpreterBuilder().withEM(es)
	check := ib.assertExpEval
	checkerr := ib.assertExpEvalErr

	checkerr(t, P("a"))

	a := NewObjConst(ObjFields{"x": T})
	es.events = map[string]*ObjConst{"a": a}
	check(t, a, P("a"))
	check(t, T, NewObjAccessor(P("a"), "x"))

	a = NewObjConst(ObjFields{"x": T})
	b := NewObjConst(ObjFields{"x": T})
	es.events = map[string]*ObjConst{"a": a, "b": b}
	check(t, T, EQ(P("a"), P("b")))
}

func TestEventVar(t *testing.T) {
	es := &mockEventStore{events: nil}
	vars := map[string]Expr{}
	ib := newInterpreterBuilder().withEM(es).withVars(vars)
	check := ib.assertExpEval

	check(t, F, V("a"))

	a := NewObjConst(ObjFields{"x": T})
	es.events = map[string]*ObjConst{"a": a}
	check(t, T, V("a"))
	check(t, F, V("b"))
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

func TestVar(t *testing.T) {
	ib := newInterpreterBuilder()
	ib.withVars(map[string]Expr{"a": TrueConst}).assertExpEval(t, TrueConst, NewVar("a"))
	ib.withVars(map[string]Expr{"a": TrueConst}).assertExpEval(t, TrueConst, NewBinOp(AndOp, NewVar("a"), NewVar("a")))
	ib.withVars(map[string]Expr{"a": TrueConst, "b": NewVar("a")}).assertExpEval(t, TrueConst, NewVar("b"))
	ib.withVars(map[string]Expr{"a": TrueConst, "b": NewVar("a")}).assertExpEval(t, TrueConst, NewBinOp(AndOp, NewVar("a"), NewVar("b")))
}

func TestObjLit(t *testing.T) {
	ib := newInterpreterBuilder()
	ib.assertExpEval(t,
		NewObjConst(ObjFields{"a": TrueConst}),
		NewObjLit(VD("a", NewBinOp(EqualOp, FalseConst, FalseConst))),
	)
}

func TestMonitor(t *testing.T) {
	assertM := func(m *MonitorDecl, fireResult *FireEventResult) {
		var cache *FuncCallCache
		i := NewInterpreter(NewEnv(nil), cache, nil, zap.NewNop())
		r, err := i.EvalMonitor(m)
		assert.NoError(t, err)
		if fireResult == nil {
			assert.Nil(t, r)
			return
		}
		assert.True(t, r.Equal(fireResult))
	}
	assertM(
		MD(ID("A"), T, nil, NewFire("e", NewObjLit(VD("a", T))), "Eth"),
		NewFireEventResult("e", NewObjConst(ObjFields{"a": T})),
	)
	assertM(
		MD(ID("A"), F, nil, NewFire("e", NewObjLit(VD("a", F))), "Eth"),
		nil,
	)
	assertM(
		MD(ID("A"), T, VD("a", OR(T, F)), NewFire("e", NewObjLit(VD("a", T))), "Eth"),
		NewFireEventResult("e", NewObjConst(ObjFields{"a": T})),
	)
	assertM(
		MD(ID("A"), F, VD("a", OR(T, F)), NewFire("e", NewObjLit(VD("a", T))), "Eth"),
		nil,
	)
	assertM(
		MD(ID("A"), T, VD("a", OR(T, F)), NewFire("e", NewObjLit(VD("a", NewVar("a")))), "Eth"),
		NewFireEventResult("e", NewObjConst(ObjFields{"a": T})),
	)
}

type MockCache struct {
	callCount int
	result    Const
	err       error
}

func (m *MockCache) getFuncCallResult(funcDecl *FuncDecl, args []Const, compute func() (Const, error)) (Const, error) {
	m.callCount++
	if m.err != nil {
		return nil, m.err
	}
	if m.result != nil {
		return m.result, nil
	}
	return compute()
}

func TestCache(t *testing.T) {
	prelude := []*NamespaceDecl{
		&NamespaceDecl{
			name: "TEST",
			funs: []*FuncDecl{
				NewFuncDecl(
					"foo",
					[]*ParamDecl{},
					IntType,
					func(env *Env, args map[string]Const) (Const, error) {
						return NewIntConstFromI64(12), nil
					},
				),
			},
		},
	}
	funcCall := NewFuncCall(NamespacePrefix{"TEST"}, "foo")
	ib := newInterpreterBuilder().withLibs(prelude)

	cache := &MockCache{}
	ib.withCache(cache).assertExpEval(t, NewIntConstFromI64(12), funcCall)
	assert.Equal(t, 1, cache.callCount)

	cache = &MockCache{result: NewIntConstFromI64(11)}
	ib.withCache(cache).assertExpEval(t, NewIntConstFromI64(11), funcCall)
	assert.Equal(t, 1, cache.callCount)

	errStr := "MOCK ERROR"
	cache = &MockCache{err: fmt.Errorf(errStr)}
	ib.withCache(cache).assertExpEvalEqualErr(t, funcCall, errStr)
	assert.Equal(t, 1, cache.callCount)
}

func assertExpEval(t *testing.T, expected Const, expr Expr) {
	newInterpreterBuilder().assertExpEval(t, expected, expr)
}

type interpreterBuilder func() *Interpreter

func newInterpreterBuilder() interpreterBuilder {
	var cache *FuncCallCache
	return func() *Interpreter {
		return NewInterpreter(NewEnv(nil), cache, nil, zap.NewNop())
	}
}

func (ib interpreterBuilder) withLibs(libs []*NamespaceDecl) interpreterBuilder {
	return func() *Interpreter {
		i := ib()
		if i.resolver == nil {
			i.resolver = NewResolver(nil, nil)
		}
		i.resolver.libs = libs
		return i
	}
}

func (ib interpreterBuilder) withEM(em eventStore) interpreterBuilder {
	return func() *Interpreter {
		i := ib()
		i.env.eventStore = em
		return i
	}
}

func (ib interpreterBuilder) withVars(vars map[string]Expr) interpreterBuilder {
	return func() *Interpreter {
		i := ib()
		i.vars = vars
		return i
	}
}

func (ib interpreterBuilder) withCache(cache Cache) interpreterBuilder {
	return func() *Interpreter {
		i := ib()
		i.cache = cache
		return i
	}
}

func (ib interpreterBuilder) assertExpEval(t *testing.T, expected Const, expr Expr) {
	result, err := ib().EvalExpr(expr)
	assert.NoError(t, err)
	assert.NotNil(t, result)
	assert.NotNil(t, result.Type())
	assert.True(t, result.Type().Equal(expected.Type()))
	assert.True(t, result.Equal(expected))
}

func (ib interpreterBuilder) assertExpEvalErr(t *testing.T, expr Expr) {
	i := ib()
	result, err := i.EvalExpr(expr)
	assert.Error(t, err)
	assert.Nil(t, result)
}

func (ib interpreterBuilder) assertExpEvalEqualErr(t *testing.T, expr Expr, errStr string) {
	i := ib()
	result, err := i.EvalExpr(expr)
	assert.EqualError(t, err, errStr)
	assert.Nil(t, result)
}

func (ib interpreterBuilder) assertEExprEval(t *testing.T, expected bool, eexpr EventExpr) {
	result, err := ib().evalEventExpr(eexpr)
	assert.NoError(t, err)
	assert.Equal(t, expected, result)
}

func assertExpEvalErr(t *testing.T, expr Expr) {
	newInterpreterBuilder().assertExpEvalErr(t, expr)
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

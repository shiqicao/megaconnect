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

	"github.com/stretchr/testify/assert"
)

func TestObjAccessorTC(t *testing.T) {
	tcb := newtc()

	tcb.assertExpErr(t, OA(T, "a"))
	tcb.assertExpErr(t, OA(OA(T, "a"), "a"))
	tcb.assertExpErr(t, OA(OC1("b", T), "a"))
	tcb.assertExpErr(t, NEG(OA(T, "a")))
	tcb.assertExpErr(t, EQ(OA(OC1("b", T), "a"), T))
	tcb.assertExpErr(t, ADD(OA(OC1("b", T), "a"), T))

	tcb.assertExp(t, OA(OC1("a", T), "a"), BoolType)
}

func TestBinTC(t *testing.T) {
	tcb := newtc()

	tcb.assertExpErr(t, AND(T, I(1)))
	tcb.assertExpErr(t, OR(I(1), T))
	tcb.assertExpErrN(t, 2, OR(I(1), R64(0)))
	tcb.assertExpErr(t, EQ(T, I(1)))
	tcb.assertExpErr(t, NE(S("a"), I(1)))
	tcb.assertExpErr(t, EQ(OC1("a", T), OC1("b", T)))
	tcb.assertExpErr(t, EQ(OC1("a", T), OC1("a", I(1))))
	tcb.assertExpErrN(t, 2, LE(OC1("a", T), OC1("a", I(1))))
	tcb.assertExpErr(t, GE(OC1("a", T), I(1)))
	tcb.assertExpErr(t, GT(S("a"), I(1)))
	tcb.assertExpErr(t, LT(T, I(1)))
	tcb.assertExpErrN(t, 2, ADD(OC1("a", T), OC1("a", I(1))))
	tcb.assertExpErr(t, MINUS(OC1("a", T), I(1)))
	tcb.assertExpErr(t, MUL(S("a"), I(1)))
	tcb.assertExpErr(t, DIV(T, I(1)))

	tcb.assertExp(t, AND(T, T), BoolType)
	tcb.assertExp(t, OR(T, F), BoolType)
	tcb.assertExp(t, EQ(T, F), BoolType)
	tcb.assertExp(t, NE(I(1), R64(1)), BoolType)
	tcb.assertExp(t, NE(I(1), I(2)), BoolType)
	tcb.assertExp(t, NE(I(1), I(2)), BoolType)
	tcb.assertExp(t, NE(OC1("a", T), OC1("a", T)), BoolType)
	tcb.assertExp(t, GE(I(1), I(2)), BoolType)
	tcb.assertExp(t, LE(I(1), R64(2)), BoolType)
	tcb.assertExp(t, LT(R64(1), R64(2)), BoolType)
	tcb.assertExp(t, GT(R64(1), I(2)), BoolType)
	tcb.assertExp(t, ADD(I(1), I(2)), IntType)
	tcb.assertExp(t, MINUS(I(1), R64(2)), RatType)
	tcb.assertExp(t, MUL(R64(1), R64(2)), RatType)
	tcb.assertExp(t, DIV(R64(1), I(2)), RatType)
}

func TestUniTC(t *testing.T) {
	tcb := newtc()

	tcb.assertExpErr(t, NOT(I(1)))
	tcb.assertExpErr(t, NOT(S("a")))
	tcb.assertExpErr(t, NOT(OC1("a", T)))
	tcb.assertExpErr(t, NEG(S("a")))
	tcb.assertExpErr(t, NEG(OC1("a", T)))
	tcb.assertExpErr(t, NEG(T))

	tcb.assertExp(t, NOT(T), BoolType)
	tcb.assertExp(t, NOT(AND(T, T)), BoolType)
	tcb.assertExp(t, NOT(NOT(T)), BoolType)
	tcb.assertExp(t, NEG(I(1)), IntType)
	tcb.assertExp(t, NEG(R64(1)), RatType)
	tcb.assertExp(t, NEG(NEG(I(1))), IntType)
}

func TestFuncTC(t *testing.T) {
	tcb := newtc()

	tcb.assertExpErr(t, FC(nil, ID("a")))

	f := FC(nil, ID("a"), I(1))
	f.decl = FD("a", Params{PARAM("a", BoolType)}, IntType, nil)
	tcb.assertExpErr(t, f)

	f = FC(nil, ID("a"), I(1), I(1))
	f.decl = FD("a", Params{PARAM("a", BoolType), PARAM("b", StrType)}, IntType, nil)
	tcb.assertExpErrN(t, 2, f)

	f = FC(nil, ID("a"), I(1))
	f.decl = FD("a", Params{PARAM("a", BoolType), PARAM("b", StrType)}, IntType, nil)
	tcb.assertExpErr(t, f)

	f = FC(nil, ID("a"), I(1), I(1))
	f.decl = FD("a", Params{PARAM("a", BoolType)}, IntType, nil)
	tcb.assertExpErr(t, f)

	f = FC(nil, ID("a"))
	f.decl = FD("a", Params{}, IntType, nil)
	tcb.assertExp(t, f, IntType)

	f = FC(nil, ID("a"), I(1))
	f.decl = FD("a", Params{PARAM("a", IntType)}, IntType, nil)
	tcb.assertExp(t, f, IntType)
}

func TestObjLitTC(t *testing.T) {
	tcb := newtc()

	tcb.assertExpErr(t, OL1(
		"a", NOT(I(1)),
	))
	tcb.assertExpErrN(t, 2, OL2(
		"a", NOT(I(1)),
		"b", NEG(F),
	))

	tcb.assertExp(t, OL(IdToExpr{}), OT(IdToTy{}))
	tcb.assertExp(t, OL1("a", T), OT1("a", BoolType))
	tcb.assertExp(t, OL2("a", T, "b", I(1)), OT2("a", BoolType, "b", IntType))
}

func TestVarTC(t *testing.T) {
	tcb := newtc()

	tcb.assertExpErr(t, V("a"))
	tcb.withStack(scope{"a": T}).assertExpErr(t, NEG(V("a")))
	tcb.withStack(scope{"a": T}).assertExpErr(t, NOT(V("b")))

	tcb.withStack(scope{"a": T}).assertExp(t, V("a"), BoolType)
	tcb.withStack(scope{"a": T}).assertExp(t, NOT(V("a")), BoolType)
	tcb.withStack(scope{"a": T, "b": I(1)}).assertExp(t, NOT(V("a")), BoolType)
}

func TestPropsTC(t *testing.T) {
	tcb := newtc()

	tcb.assertExpErr(t, P("a"))
	tcb.withStack(scope{"a": T}).assertExpErr(t, P("a"))

	wf := W("a").AddChild(ED("a", OT1("a", BoolType)))
	tcb.withWf(wf).assertExp(t, P("a"), OT1("a", BoolType))
	tcb.withWf(wf).assertExpErr(t, P("b"))
	tcb.withWf(wf).withStack(scope{"a": T}).assertExp(t, P("a"), OT1("a", BoolType))
}

func TestMonitorTC(t *testing.T) {
	tcb := newtc()

	// Test Vars
	wf := W("a").AddChild(ED("a", OT1("a", BoolType)))
	md := func(vars IdToExpr) *MonitorDecl {
		return MD(ID("a"), T, vars, FIRE("a", OL1("a", T)), "B")
	}
	tcb = tcb.withWf(wf)
	tcb.assertMonitor(t, md(VD("a", T)))
	tcb.assertMonitor(t, md(VD("a", T).Put("b", V("a"))))

	tcb.assertMonitorErrN(t, 1, md(VD("a", NEG(T)).Put("b", V("a"))))
	tcb.assertMonitorErrN(t, 2, md(VD("a", NEG(T)).Put("b", NOT(I(1)))))

	// Test Fire event
	md1 := func(eventName string, eventObj Expr) *MonitorDecl {
		return MD(ID("a"), T, VD("a", T), FIRE(eventName, eventObj), "B")
	}
	tcb.assertMonitor(t, md1("a", OL1("a", T)))
	tcb.assertMonitor(t, md1("a", OL1("a", V("a"))))

	tcb.assertMonitorErrN(t, 1, md1("b", OL1("a", T)))
	tcb.assertMonitorErrN(t, 1, md1("a", T))
	tcb.assertMonitorErrN(t, 1, md1("a", OL1("a", I(1))))
}

func TestActionTC(t *testing.T) {
	tcb := newtc()

	// Test trigger
	tcb.assertActionErrN(t, 1, ACT(ID("a"), EV("a"), Stmts{}))
	wf := W("a").AddChild(ED("a", OT1("a", BoolType)))
	tcb.withWf(wf).assertActionErrN(t, 1, ACT(ID("a"), EAND(EV("a"), EV("b")), Stmts{}))
	tcb.withWf(wf).assertActionErrN(t, 1, ACT(ID("a"), EOR(EV("a"), EV("b")), Stmts{}))

	tcb.withWf(wf).assertAction(t, ACT(ID("a"), EOR(EV("a"), EV("a")), Stmts{}))
	wf = W("a").AddChild(ED("a", OT1("a", BoolType))).AddChild(ED("b", OT1("b", IntType)))
	tcb.withWf(wf).assertAction(t, ACT(ID("a"), EOR(EV("a"), EV("b")), Stmts{}))

	// Test Stmt
	wf = W("a").AddChild(ED("a", OT1("a", BoolType))).AddChild(ED("b", OT1("b", IntType)))
	tcb.withWf(wf).assertAction(t, ACT(ID("a"), EV("a"), Stmts{FIRE("a", OL1("a", T))}))
	tcb.withWf(wf).assertActionErrN(t, 1, ACT(ID("a"), EV("a"), Stmts{FIRE("b", OL1("a", T))}))
	tcb.withWf(wf).assertActionErrN(t, 1, ACT(ID("a"), EV("a"), Stmts{FIRE("a", OL1("a", I(1)))}))
	tcb.withWf(wf).assertActionErrN(t, 1, ACT(ID("a"), EV("a"), Stmts{FIRE("a", OL1("b", T))}))
}

type typecheckerBuilder func() *TypeChecker

func newtc() typecheckerBuilder {
	return func() *TypeChecker {
		return &TypeChecker{}
	}
}

func (tb typecheckerBuilder) withWf(wf *WorkflowDecl) typecheckerBuilder {
	return func() *TypeChecker {
		old := tb()
		new := NewTypeChecker(wf)
		new.stack = old.stack
		return new
	}
}

func (tb typecheckerBuilder) withStack(s map[string]Expr) typecheckerBuilder {
	return func() *TypeChecker {
		tc := tb()
		tc.stack = s
		return tc
	}
}

func (tb typecheckerBuilder) assertExp(t *testing.T, e Expr, ty Type) {
	tc := tb()
	errs := tc.inferExpr(e).ToErr()
	assert.Equal(t, 0, len(errs))
	assert.True(t, ty.Equal(e.Type()), fmt.Sprintf("Expected Type %v, actual type %v", ty, e.Type()))
}

func (tb typecheckerBuilder) assertErrN(
	t *testing.T,
	n int,
	run func(*TypeChecker) Errors,
) {
	errs := run(tb())
	assertErrs(t, n, errs)
}

func (tb typecheckerBuilder) assertExpErrN(t *testing.T, n int, e Expr) {
	tb.assertErrN(t, n, func(tc *TypeChecker) Errors {
		return tc.inferExpr(e)
	})
}

func (tb typecheckerBuilder) assertMonitorErrN(t *testing.T, n int, m *MonitorDecl) {
	tb.assertErrN(t, n, func(tc *TypeChecker) Errors {
		return tc.checkMonitor(m)
	})
}

func (tb typecheckerBuilder) assertActionErrN(t *testing.T, n int, a *ActionDecl) {
	tb.assertErrN(t, n, func(tc *TypeChecker) Errors {
		return tc.checkAction(a)
	})
}

func (tb typecheckerBuilder) assertExpErr(t *testing.T, e Expr) {
	tb.assertExpErrN(t, 1, e)
}

func (tb typecheckerBuilder) assertMonitor(t *testing.T, m *MonitorDecl) {
	tb.assertMonitorErrN(t, 0, m)
}

func (tb typecheckerBuilder) assertAction(t *testing.T, a *ActionDecl) {
	tb.assertActionErrN(t, 0, a)
}

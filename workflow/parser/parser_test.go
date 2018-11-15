// Copyright 2018 @ MegaSpace

// The MegaSpace source code is distributed in the hope that it will be useful,
// but WITHOUT ANY WARRANTY; without even the implied warranty of
// MERCHANTABILITY or FITNESS FOR A PARTICULAR PURPOSE. See the
// GNU Lesser General Public License for more details.

// You should have received a copy of the GNU Lesser General Public License
// along with the MegaSpace source code. If not, see <http://www.gnu.org/licenses/>.

package parser

import (
	"fmt"
	"math/big"
	"testing"

	"github.com/stretchr/testify/assert"

	wf "github.com/megaspacelab/megaconnect/workflow"
	"github.com/megaspacelab/megaconnect/workflow/parser/gen/lexer"

	"github.com/megaspacelab/megaconnect/workflow/parser/gen/parser"
)

func TestParser(t *testing.T) {
	r, err := parse(t, `
	workflow w {
		monitor a {
	  		chain Eth 
	  		condition true
	  		var {
		  		a = true
		  		b = false
	  		}
	  		fire e {
				a: true,
				b: true && false
			}		
		}  
	}
	`)
	assert.NoError(t, err)
	workflow, ok := r.(*wf.WorkflowDecl)
	assert.True(t, ok)
	assert.NotNil(t, workflow)
	md := workflow.MonitorDecls()
	assert.True(t, md[0].Equal(wf.NewMonitorDecl(
		wf.NewId("a"),
		wf.TrueConst,
		wf.NewIdToExpr().Put("a", T).Put("b", F),
		wf.NewFire("e", wf.NewObjLit(wf.NewIdToExpr().Put("a", T).Put("b", AND(T, F)))),
		"Eth",
	)))
}

var (
	T = wf.TrueConst
	F = wf.FalseConst
	B = func(op wf.Operator) func(x wf.Expr, y wf.Expr) wf.Expr {
		return func(x wf.Expr, y wf.Expr) wf.Expr {
			return wf.NewBinOp(op, x, y)
		}
	}
	U = func(op wf.Operator) func(wf.Expr) wf.Expr {
		return func(x wf.Expr) wf.Expr {
			return wf.NewUniOp(op, x)
		}
	}
	AND   = B(wf.AndOp)
	OR    = B(wf.OrOp)
	EQ    = B(wf.EqualOp)
	GT    = B(wf.GreaterThanOp)
	V     = wf.NewVar
	ADD   = B(wf.PlusOp)
	MINUS = B(wf.MinusOp)
	MUL   = B(wf.MultOp)
	I     = wf.NewIntConstFromI64
	R     = func(n int64, d int64) *wf.RatConst { return wf.NewRatConst(big.NewRat(n, d)) }
	S     = wf.NewStrConst
	OA    = wf.NewObjAccessor
	ID    = wf.NewId
	FC    = func(ns wf.NamespacePrefix, id string, args ...wf.Expr) *wf.FuncCall {
		return wf.NewFuncCall(ns, ID(id), args...)
	}
	AD   = wf.NewActionDecl
	FIRE = wf.NewFire
	OL   = wf.NewObjLit
	IE   = wf.NewIdToExpr
	IE1  = func(id string, e wf.Expr) wf.IdToExpr {
		return wf.NewIdToExpr().Put(id, e)
	}
	N = func(ns ...string) wf.NamespacePrefix {
		nss := wf.NamespacePrefix{}
		for _, n := range ns {
			nss = append(nss, ID(n))
		}
		return nss
	}
	P   = func(s string) *wf.Props { return wf.NewProps(V(s)) }
	NEG = U(wf.MinusOp)
	NOT = U(wf.NotOp)

	// EventExpr
	EV = wf.NewEVar
	EB = func(op wf.EventExprOperator) func(x wf.EventExpr, y wf.EventExpr) wf.EventExpr {
		return func(x wf.EventExpr, y wf.EventExpr) wf.EventExpr {
			return wf.NewEBinOp(op, x, y)
		}
	}
	EAND = EB(wf.AndEOp)
	EOR  = EB(wf.OrEOp)
)

func TestComment(t *testing.T) {
	r, err := parse(t, `
	// Comment
	workflow a {}
	`)
	assert.NoError(t, err)
	assert.NotNil(t, r)
	w, ok := r.(*wf.WorkflowDecl)
	assert.True(t, ok)
	assert.NotNil(t, w)
	assert.Equal(t, "a", w.Name().Id())

	r, err = parse(t, `
	// Comment workflow a {}
	workflow b {}
	`)
	assert.NoError(t, err)
	assert.NotNil(t, r)
	w, ok = r.(*wf.WorkflowDecl)
	assert.True(t, ok)
	assert.NotNil(t, w)
	assert.Equal(t, "b", w.Name().Id())
}

func TestDecl(t *testing.T) {
	decl1 := "event e {}"
	decl2 := "monitor m { chain b condition true var { a = true } fire e {} }"
	decl3 := "action a { trigger x run {} }"
	code := fmt.Sprintf("workflow a { %s %s %s }", decl1, decl2, decl3)
	w := assertWorkflowParsing(t, code)
	a := w.ActionDecls()
	assert.Len(t, a, 1)

	m := w.MonitorDecls()
	assert.Len(t, m, 1)
	e := w.EventDecls()
	assert.Len(t, e, 1)
}

func TestExprParsing(t *testing.T) {
	assertExprParsing(t, AND(T, T), "true && true")
	assertExprParsing(t, OR(AND(T, T), F), "true && true || false")
	assertExprParsing(t, AND(OR(T, T), F), "true || true && false")
	assertExprParsing(t, AND(EQ(T, T), T), "true == true && true")
	assertExprParsing(t, V("a"), "a")
	assertExprParsing(t, AND(V("a"), GT(V("b"), V("c"))), "a && b > c")
	assertExprParsing(t, GT(AND(V("a"), V("b")), V("c")), "(a && b) > c")
	assertExprParsing(t, MUL(ADD(V("a"), V("b")), V("c")), "(a + b) * c")
	assertExprParsing(t, ADD(V("a"), MUL(V("b"), V("c"))), "a + b * c")
}

func TestStrLit(t *testing.T) {
	assertExprParsing(t, S(""), `""`)
	assertExprParsing(t, S(" "), `" "`)
	assertExprParsing(t, S("a+b"), `"a+b"`)
	assertExprParsing(t, S("a"), `"a"`)
}

func TestUniOp(t *testing.T) {
	assertExprParsing(t, NEG(I(1)), "-1")
	assertExprParsing(t, MINUS(NEG(I(1)), I(1)), "-1-1")
	assertExprParsing(t, MINUS(NEG(I(1)), NEG(I(1))), "-1-(-1)")
	assertExprParsing(t, NEG(ADD(I(1), I(1))), "-(1+1)")
	assertExprParsing(t, NEG(R(11, 10)), "-1.1")
	assertExprParsing(t, NEG(OA(V("a"), "a")), "-a.a")
	assertExprParsing(t, NEG(FC(nil, "a")), "-a()")
	assertExprParsing(t, NEG(V("a")), "-a")
	assertExprParsing(t, NEG(NEG(V("a"))), "--a")
	assertExprParsing(t, NEG(OA(P("a"), "a")), "-props(a).a")

	assertExprParsing(t, NOT(T), "!true")
	assertExprParsing(t, AND(NOT(T), T), "!true && true")
	assertExprParsing(t, OR(F, NOT(F)), "false || !false")
	assertExprParsing(t, NOT(AND(T, T)), "!(true && true)")
	assertExprParsing(t, NOT(OA(V("a"), "a")), "!a.a")
	assertExprParsing(t, NOT(FC(nil, "a")), "!a()")
	assertExprParsing(t, NOT(V("a")), "!a")
	assertExprParsing(t, NOT(NOT(V("a"))), "!!a")
	assertExprParsing(t, NOT(OA(P("a"), "a")), "!props(a).a")
}

func TestIntLit(t *testing.T) {
	assertExprParsing(t, I(0), "0")
	assertExprParsing(t, I(1), "1")
	i := "1"
	for ; len(i) < 20; i = i + i {
	}
	expected, ok := new(big.Int).SetString(i, 10)
	assert.True(t, ok)
	assertExprParsing(t, wf.NewIntConst(expected), i)

	assertExprParsingErr(t, "01")
}

func TestRatLit(t *testing.T) {
	assertExprParsing(t, R(0, 1), "0.0")
	assertExprParsing(t, R(1, 10), "0.1")
	assertExprParsing(t, R(101, 10), "10.1")
	assertExprParsing(t, R(101, 10), "10.10")
	assertExprParsing(t, ADD(R(101, 10), I(1)), "10.10 + 1")

	assertExprParsingErr(t, "01.1")
	assertExprParsingErr(t, "00.0")
}

func TestObjAccessor(t *testing.T) {
	assertExprParsing(t, OA(V("A"), "a"), "A.a")
	assertExprParsing(t, OA(OA(V("A"), "a"), "b"), "A.a.b")
	assertExprParsing(t, OA(AND(V("A"), V("B")), "a"), "(A && B).a")
	assertExprParsing(t, AND(V("A"), OA(V("B"), "a")), "A && B.a")
	assertExprParsing(t, MUL(V("A"), OA(V("B"), "a")), "A * B.a")
	assertExprParsing(t, MUL(OA(V("A"), "b"), OA(V("B"), "a")), "A.b * B.a")
}

func TestObjLit(t *testing.T) {
	assertExprParsing(t, wf.NewObjLit(wf.NewIdToExpr()), "{}")
	assertExprParsing(t, wf.NewObjLit(wf.NewIdToExpr().Put("a", T)), "{a: true}")
	assertExprParsing(t, wf.NewObjLit(wf.NewIdToExpr().Put("a", T).Put("b", F)), "{a: true, b: false}")
	assertExprParsing(t,
		wf.NewObjLit(
			wf.NewIdToExpr().
				Put("a", T).
				Put("b", wf.NewObjLit(wf.NewIdToExpr().Put("c", T))),
		),
		"{a: true, b: {c: true}}",
	)
}

func TestEvent(t *testing.T) {
	assertEventParsing(t, wf.NewEventDecl(ID("e"), wf.NewObjType(wf.NewIdToTy())), "event e {}")
	assertEventParsing(t, wf.NewEventDecl(ID("e"), wf.NewObjType(wf.NewIdToTy().Put("a", wf.IntType))), "event e {a : int}")
	assertEventParsing(t, wf.NewEventDecl(ID("e"), wf.NewObjType(wf.NewIdToTy().Put("a", wf.RatType))), "event e {a : rat}")
	assertEventParsing(
		t,
		wf.NewEventDecl(ID("e"), wf.NewObjType(wf.NewIdToTy().Put("a", wf.IntType).Put("b", wf.StrType))),
		"event e {a : int, b: string}",
	)
	assertEventParsing(
		t,
		wf.NewEventDecl(ID("e"), wf.NewObjType(wf.NewIdToTy().Put("a", wf.IntType).Put("b", wf.BoolType))),
		"event e {a : int, b: bool}",
	)
	assertEventParsing(
		t,
		wf.NewEventDecl(ID("e"), wf.NewObjType(
			wf.NewIdToTy().
				Put("a", wf.IntType).
				Put("b", wf.NewObjType(wf.NewIdToTy().Put("c", wf.BoolType)))),
		),
		"event e {a : int, b: {c : bool}}",
	)
}

func TestAction(t *testing.T) {
	assertActionParsing(
		t,
		AD(
			ID("a"),
			EV("b"),
			wf.Stmts{
				FIRE("c", OL(IE1("d", T))),
				FIRE("c", OL(IE1("d", T))),
			},
		),
		`action a {
			trigger b
			run {
				fire c {d: true};
				fire c {d: true};				
			}
		}
		`,
	)
}

func TestStmts(t *testing.T) {
	assertStmtParsing(t, wf.Stmts{}, "")
	assertStmtParsing(t, wf.Stmts{FIRE("e", OL(IE()))}, "fire e {};")
	assertStmtParsing(t, wf.Stmts{FIRE("e", FC(nil, "f"))}, "fire e f();")
	assertStmtParsing(t,
		wf.Stmts{
			FIRE("e", FC(nil, "f")),
			FIRE("a", OL(IE1("a", T))),
		},
		"fire e f();fire a {a: true};",
	)
}

func TestEventExpr(t *testing.T) {
	assertEventExprParsing(t, EV("a"), "a")
	assertEventExprParsing(t, EAND(EV("a"), EV("b")), "a && b")
	assertEventExprParsing(t, EOR(EV("a"), EV("b")), "a || b")
	assertEventExprParsing(t, EAND(EOR(EV("a"), EV("b")), EV("c")), "a || b && c")
	assertEventExprParsing(t, EOR(EV("a"), EAND(EV("b"), EV("c"))), "a || (b && c)")
	assertEventExprParsing(t, EOR(EOR(EV("a"), EV("b")), EV("c")), "a || b || c")
}

func TestFuncCall(t *testing.T) {
	assertExprParsing(t, FC(nil, "a"), "a()")
	assertExprParsing(t, FC(nil, "a", V("b")), "a(b)")
	assertExprParsing(t, FC(nil, "a", V("b"), V("c")), "a(b, c)")
	assertExprParsing(t, FC(nil, "a", V("b"), V("c"), V("d")), "a(b, c, d)")
	assertExprParsing(t, FC(nil, "a", V("b"), FC(nil, "c")), "a(b, c())")
	assertExprParsing(t, FC(nil, "a", OA(V("b"), "c"), FC(nil, "c")), "a(b.c, c())")
	assertExprParsing(t, OA(FC(nil, "a"), "b"), "a().b")

	// test namespace
	assertExprParsing(t, FC(N("x"), "a"), "x::a()")
	assertExprParsing(t, FC(N("x", "y"), "a"), "x::y::a()")
	assertExprParsing(t, OA(OA(FC(N("x", "y"), "a"), "b"), "c"), "x::y::a().b.c")
}

func TestProps(t *testing.T) {
	assertExprParsing(t, P("a"), "props(a)")
	assertExprParsing(t, OA(P("a"), "b"), "props(a).b")
}

func assertExprParsingErr(t *testing.T, expr string) {
	code := fmt.Sprintf("workflow b { monitor a { chain Eth condition %s var { a = true } fire e { a : true } } }", expr)
	_, err := parse(t, code)
	assert.Error(t, err)
}

func assertExprParsing(t *testing.T, expected wf.Expr, expr string) {
	code := fmt.Sprintf("workflow b { monitor a { chain Eth condition %s var { a = true } fire e { a : true } } }", expr)
	w := assertWorkflowParsing(t, code)
	md := w.MonitorDecls()[0]
	assert.NotNil(t, md)
	if !md.Condition().Equal(expected) {
		t.Logf(
			"Actual Parsed Expr(%T): %s",
			md.Condition(),
			md.Condition().String(),
		)
		t.Logf(
			"Expected Expr(%T): %s",
			expected,
			expected.String(),
		)
	}
	assert.True(t, md.Condition().Equal(expected))
}

func assertEventParsing(t *testing.T, expected *wf.EventDecl, event string) {
	code := fmt.Sprintf("workflow w { %s }", event)
	w := assertWorkflowParsing(t, code)
	ed := w.EventDecls()[0]
	assert.NotNil(t, ed)
	if !ed.Equal(expected) {
		/* TODO: add String to EventDecl
		t.Logf(
			"Actual Parsed Event(%T): %s",
			ed.String()
		)
		*/
	}
	assert.True(t, ed.Equal(expected))
}

func assertActionParsing(t *testing.T, expected *wf.ActionDecl, action string) {
	code := fmt.Sprintf("workflow w { %s }", action)
	w := assertWorkflowParsing(t, code)
	act := w.ActionDecls()[0]
	assert.NotNil(t, act)
	assert.True(t, act.Equal(expected))
}

func assertStmtParsing(t *testing.T, expected wf.Stmts, stmts string) {
	code := fmt.Sprintf("workflow w { action a { trigger a run{ %s } } }", stmts)
	w := assertWorkflowParsing(t, code)
	action := w.ActionDecls()[0]
	assert.NotNil(t, action)
	assert.True(t, action.Body().Equal(expected))
}

func assertEventExprParsing(t *testing.T, expected wf.EventExpr, expr string) {
	code := fmt.Sprintf("workflow w { action a { trigger %s run{} } }", expr)
	w := assertWorkflowParsing(t, code)
	action := w.ActionDecls()[0]
	assert.NotNil(t, action)
	assert.True(t, action.Trigger().Equal(expected))
}

func assertWorkflowParsing(t *testing.T, w string) *wf.WorkflowDecl {
	r, err := parse(t, w)
	assert.NoError(t, err)
	wfl, ok := r.(*wf.WorkflowDecl)
	assert.True(t, ok)
	assert.NotNil(t, wfl)
	return wfl
}

func parse(t *testing.T, code string) (interface{}, error) {
	s := lexer.NewLexer([]byte(code))
	p := parser.NewParser()
	return p.Parse(s)
}

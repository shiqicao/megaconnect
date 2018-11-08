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
		monitor a 
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
	AND = B(wf.AndOp)
	OR  = B(wf.OrOp)
	EQ  = B(wf.EqualOp)
	GT  = B(wf.GreaterThanOp)
	V   = wf.NewVar
	ADD = B(wf.PlusOp)
	MUL = B(wf.MultOp)
	I   = wf.NewIntConstFromI64
	S   = wf.NewStrConst
	OA  = wf.NewObjAccessor
)

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
	assertExprParsing(t, S(""), "\"\"")
	assertExprParsing(t, S(" "), "\" \"")
	assertExprParsing(t, S("a+b"), "\"a+b\"")
	assertExprParsing(t, S("a"), "\"a\"")
}

func TestIntLit(t *testing.T) {
	assertExprParsing(t, I(0), "0")
	assertExprParsing(t, I(1), "1")
	//assertExprParsing(t, I(-1), "-1")
	i := "1"
	for ; len(i) < 20; i = i + i {
	}
	expected, ok := new(big.Int).SetString(i, 10)
	assert.True(t, ok)
	assertExprParsing(t, wf.NewIntConst(expected), i)

	assertExprParsingErr(t, "01")
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

func assertExprParsingErr(t *testing.T, expr string) {
	code := fmt.Sprintf("workflow b { monitor a chain Eth condition %s var { a = true } fire e { a : true } }", expr)
	_, err := parse(t, code)
	assert.Error(t, err)
}

func assertExprParsing(t *testing.T, expected wf.Expr, expr string) {
	code := fmt.Sprintf("workflow b { monitor a chain Eth condition %s var { a = true } fire e { a : true } }", expr)
	r, err := parse(t, code)
	assert.NoError(t, err)
	w, ok := r.(*wf.WorkflowDecl)
	assert.True(t, ok)
	assert.NotNil(t, w)
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

func parse(t *testing.T, code string) (interface{}, error) {
	s := lexer.NewLexer([]byte(code))
	p := parser.NewParser()
	return p.Parse(s)
}

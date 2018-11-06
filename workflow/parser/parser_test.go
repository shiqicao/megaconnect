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
	"testing"

	"github.com/stretchr/testify/assert"

	wf "github.com/megaspacelab/megaconnect/workflow"
	"github.com/megaspacelab/megaconnect/workflow/parser/goccgen/lexer"

	"github.com/megaspacelab/megaconnect/workflow/parser/goccgen/parser"
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
				b: true && false,
			}		  
	}
	`)
	assert.NoError(t, err)
	workflow, ok := r.(*wf.WorkflowDecl)
	assert.True(t, ok)
	assert.NotNil(t, workflow)
	md := workflow.MonitorDecls()
	assert.True(t, md[0].Equal(wf.NewMonitorDecl(
		"a",
		wf.TrueConst,
		wf.VarDecls{"a": wf.TrueConst, "b": wf.FalseConst},
		wf.NewFire("e", wf.NewObjLit(wf.VarDecls{"a": T, "b": AND(T, F)})),
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
	assertExprParsing(t, wf.NewStrConst(""), "\"\"")
	assertExprParsing(t, wf.NewStrConst(" "), "\" \"")
	assertExprParsing(t, wf.NewStrConst("a+b"), "\"a+b\"")
	assertExprParsing(t, wf.NewStrConst("a"), "\"a\"")
}

func assertExprParsing(t *testing.T, expected wf.Expr, expr string) {
	code := fmt.Sprintf("workflow b { monitor a chain Eth condition %s var { a = true } fire e { a : true, } }", expr)
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

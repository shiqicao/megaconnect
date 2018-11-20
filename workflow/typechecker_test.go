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

	tcb.assertExpTCErr(t, OA(T, "a"))
	tcb.assertExpTCErr(t, OA(OC1("b", T), "a"))

	tcb.assertExpTC(t, OA(OC1("a", T), "a"), BoolType)
}

func TestBinTC(t *testing.T) {
	tcb := newtc()

	tcb.assertExpTCErr(t, AND(T, I(1)))
	tcb.assertExpTCErr(t, OR(I(1), T))
	tcb.assertExpTCErrN(t, 2, OR(I(1), R64(0)))
	tcb.assertExpTCErr(t, EQ(T, I(1)))
	tcb.assertExpTCErr(t, NE(S("a"), I(1)))
	tcb.assertExpTCErr(t, EQ(OC1("a", T), OC1("b", T)))
	tcb.assertExpTCErr(t, EQ(OC1("a", T), OC1("a", I(1))))
	tcb.assertExpTCErrN(t, 2, LE(OC1("a", T), OC1("a", I(1))))
	tcb.assertExpTCErr(t, GE(OC1("a", T), I(1)))
	tcb.assertExpTCErr(t, GT(S("a"), I(1)))
	tcb.assertExpTCErr(t, LT(T, I(1)))
	tcb.assertExpTCErr(t, ADD(OC1("a", T), OC1("a", I(1))))
	tcb.assertExpTCErr(t, MINUS(OC1("a", T), I(1)))
	tcb.assertExpTCErr(t, MUL(S("a"), I(1)))
	tcb.assertExpTCErr(t, DIV(T, I(1)))

	tcb.assertExpTC(t, AND(T, T), BoolType)
	tcb.assertExpTC(t, OR(T, F), BoolType)
	tcb.assertExpTC(t, EQ(T, F), BoolType)
	tcb.assertExpTC(t, NE(I(1), R64(1)), BoolType)
	tcb.assertExpTC(t, NE(I(1), I(2)), BoolType)
	tcb.assertExpTC(t, NE(I(1), I(2)), BoolType)
	tcb.assertExpTC(t, NE(OC1("a", T), OC1("a", T)), BoolType)
	tcb.assertExpTC(t, GE(I(1), I(2)), BoolType)
	tcb.assertExpTC(t, LE(I(1), R64(2)), BoolType)
	tcb.assertExpTC(t, LT(R64(1), R64(2)), BoolType)
	tcb.assertExpTC(t, GT(R64(1), I(2)), BoolType)
	tcb.assertExpTC(t, ADD(I(1), I(2)), IntType)
	tcb.assertExpTC(t, MINUS(I(1), R64(2)), RatType)
	tcb.assertExpTC(t, MUL(R64(1), R64(2)), RatType)
	tcb.assertExpTC(t, DIV(R64(1), I(2)), RatType)
}

func TestUniTC(t *testing.T) {
	tcb := newtc()

	tcb.assertExpTCErr(t, NOT(I(1)))
	tcb.assertExpTCErr(t, NOT(S("a")))
	tcb.assertExpTCErr(t, NOT(OC1("a", T)))
	tcb.assertExpTCErr(t, NEG(S("a")))
	tcb.assertExpTCErr(t, NEG(OC1("a", T)))
	tcb.assertExpTCErr(t, NEG(T))

	tcb.assertExpTC(t, NOT(T), BoolType)
	tcb.assertExpTC(t, NOT(AND(T, T)), BoolType)
	tcb.assertExpTC(t, NOT(NOT(T)), BoolType)
	tcb.assertExpTC(t, NEG(I(1)), IntType)
	tcb.assertExpTC(t, NEG(R64(1)), RatType)
	tcb.assertExpTC(t, NEG(NEG(I(1))), IntType)
}

func TestVarTC(t *testing.T) {
	tcb := newtc()

	tcb.assertExpTCErr(t, V("a"))
	tcb.withStack(scope{"a": T}).assertExpTCErr(t, NEG(V("a")))
	tcb.withStack(scope{"a": T}).assertExpTCErr(t, NOT(V("b")))

	tcb.withStack(scope{"a": T}).assertExpTC(t, V("a"), BoolType)
	tcb.withStack(scope{"a": T}).assertExpTC(t, NOT(V("a")), BoolType)
	tcb.withStack(scope{"a": T, "b": I(1)}).assertExpTC(t, NOT(V("a")), BoolType)
	tcb.withStack(scope{"a": V("b"), "b": I(1)}).assertExpTC(t, V("a"), IntType)
}

type typecheckerBuilder func() *TypeChecker

func newtc() typecheckerBuilder {
	return func() *TypeChecker {
		return &TypeChecker{}
	}
}

func (tb typecheckerBuilder) withStack(s map[string]Expr) typecheckerBuilder {
	return func() *TypeChecker {
		tc := tb()
		tc.stack = s
		return tc
	}
}

func (tb typecheckerBuilder) assertExpTC(t *testing.T, e Expr, ty Type) {
	tc := tb()
	errs := tc.inferExpr(e).ToErr()
	assert.Equal(t, 0, len(errs))
	assert.True(t, ty.Equal(e.Type()), fmt.Sprintf("Expected Type %v, actual type %v", ty, e.Type()))
}

func (tb typecheckerBuilder) assertExpTCErrN(t *testing.T, n int, e Expr) {
	tc := tb()
	errs := tc.inferExpr(e).ToErr()
	if len(errs) != n {
		t.Log("Actual errors:\n")
		for i, e := range errs {
			t.Logf("%d: %s\n", i, e.Error())
		}
	}

	assert.Equal(t, n, len(errs), fmt.Sprintf("Expected 1 error, actual got these errors %v", errs))
}

func (tb typecheckerBuilder) assertExpTCErr(t *testing.T, e Expr) {
	tb.assertExpTCErrN(t, 1, e)
}

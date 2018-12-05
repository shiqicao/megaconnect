// Copyright 2018 @ MegaSpace

// The MegaSpace source code is distributed in the hope that it will be useful,
// but WITHOUT ANY WARRANTY; without even the implied warranty of
// MERCHANTABILITY or FITNESS FOR A PARTICULAR PURPOSE. See the
// GNU Lesser General Public License for more details.

// You should have received a copy of the GNU Lesser General Public License
// along with the MegaSpace source code. If not, see <http://www.gnu.org/licenses/>.

package workflow

import (
	"math/big"
	"testing"

	"github.com/stretchr/testify/assert"
)

var (
	// EventExpr constructors
	EV   = NewEVar
	EOR  = func(x EventExpr, y EventExpr) EventExpr { return NewEBinOp(OrEOp, x, y) }
	EAND = func(x EventExpr, y EventExpr) EventExpr { return NewEBinOp(AndEOp, x, y) }
)

var (
	// Expr constructors
	T   = TrueConst
	F   = FalseConst
	R   = func(n int64, d int64) *RatConst { return NewRatConst(big.NewRat(n, d)) }
	R64 = NewRatConstFromF64
	I   = NewIntConstFromI64
	S   = NewStrConst
	BIN = func(op Operator) func(Expr, Expr) *BinOp {
		return func(x Expr, y Expr) *BinOp { return NewBinOp(op, x, y) }
	}
	UNI = func(op Operator) func(Expr) *UniOp {
		return func(x Expr) *UniOp { return NewUniOp(op, x) }
	}
	AND  = BIN(AndOp)
	OR   = BIN(OrOp)
	EQ   = BIN(EqualOp)
	NE   = BIN(NotEqualOp)
	LT   = BIN(LessThanOp)
	LE   = BIN(LessThanEqualOp)
	GT   = BIN(GreaterThanOp)
	GE   = BIN(GreaterThanEqualOp)
	FIRE = NewFire
	MD   = NewMonitorDecl
	ACT  = NewActionDecl
	V    = NewVar
	P    = func(a string) *Props { return NewProps(V(a)) }
	OA   = NewObjAccessor
	FC   = NewFuncCall
	// Variable declaration
	VD = func(f string, e Expr) IdToExpr { return make(IdToExpr).Put(f, e) }
	// Variable types
	VT = func(f string, ty Type) IdToTy { return make(IdToTy).Put(f, ty) }
	ID = NewId

	// ObjConst
	OC  = NewObjConst
	OC1 = func(f string, c Const) *ObjConst { return OC(ObjFields{f: c}) }

	// IdTExpr
	I2E  = IdToExpr{}
	I2E1 = func(id string, e Expr) IdToExpr { return IdToExpr{id: IdExpr{Id: ID(id), Expr: e}} }
	I2E2 = func(id1 string, e1 Expr, id2 string, e2 Expr) IdToExpr {
		return IdToExpr{
			id1: IdExpr{Id: ID(id1), Expr: e1},
			id2: IdExpr{Id: ID(id2), Expr: e2},
		}
	}

	// ObjLit
	OL  = NewObjLit
	OL1 = func(f string, e Expr) *ObjLit {
		return OL(IdToExpr{
			f: IdExpr{Id: ID(f), Expr: e},
		})
	}
	OL2 = func(f1 string, e1 Expr, f2 string, e2 Expr) *ObjLit {
		return OL(IdToExpr{
			f1: IdExpr{Id: ID(f1), Expr: e1},
			f2: IdExpr{Id: ID(f2), Expr: e2},
		})
	}

	// ObjTy
	OT  = NewObjType
	OT1 = func(f string, ty Type) *ObjType {
		return OT(IdToTy{
			f: idTy{id: ID(f), ty: ty},
		})
	}
	OT2 = func(f1 string, ty1 Type, f2 string, ty2 Type) *ObjType {
		return OT(IdToTy{
			f1: idTy{id: ID(f1), ty: ty1},
			f2: idTy{id: ID(f2), ty: ty2},
		})
	}

	// Arith
	ADD   = BIN(PlusOp)
	MINUS = BIN(MinusOp)
	MUL   = BIN(MultOp)
	DIV   = BIN(DivOp)

	NEG = UNI(MinusOp)
	NOT = UNI(NotOp)

	NS    = NewNamespaceDecl
	FD    = NewFuncDecl
	PARAM = NewParamDecl

	W  = func(n string) *WorkflowDecl { return NewWorkflowDecl(ID(n), 0) }
	ED = func(n string, ty *ObjType) *EventDecl { return NewEventDecl(ID(n), ty) }
)

// Errors
func assertErrs(t *testing.T, n int, errs Errors) {
	errors := errs.ToErr()
	if len(errors) != n {
		assert.Fail(t, "Invalid error count", "Expected %d errors Actual %d errors", n, len(errors))
		t.Log("Actual errors:\n")
		for i, e := range errors {
			t.Logf("%d: %s\n", i, e.Error())
		}
	}
}

func assertNoErrs(t *testing.T, errs Errors) {
	assertErrs(t, 0, errs)
}

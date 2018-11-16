// Copyright 2018 @ MegaSpace

// The MegaSpace source code is distributed in the hope that it will be useful,
// but WITHOUT ANY WARRANTY; without even the implied warranty of
// MERCHANTABILITY or FITNESS FOR A PARTICULAR PURPOSE. See the
// GNU Lesser General Public License for more details.

// You should have received a copy of the GNU Lesser General Public License
// along with the MegaSpace source code. If not, see <http://www.gnu.org/licenses/>.

package workflow

import "math/big"

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
	BIN = func(op Operator) func(Expr, Expr) *BinOp {
		return func(x Expr, y Expr) *BinOp { return NewBinOp(op, x, y) }
	}
	UNI = func(op Operator) func(Expr) *UniOp {
		return func(x Expr) *UniOp { return NewUniOp(op, x) }
	}
	AND  = BIN(AndOp)
	OR   = BIN(OrOp)
	EQ   = BIN(EqualOp)
	LT   = BIN(LessThanOp)
	LE   = BIN(LessThanEqualOp)
	GT   = BIN(GreaterThanOp)
	GE   = BIN(GreaterThanEqualOp)
	FIRE = NewFire
	MD   = NewMonitorDecl
	V    = NewVar
	P    = func(a string) *Props { return NewProps(V(a)) }
	// Variable declaration
	VD = func(f string, e Expr) IdToExpr { return make(IdToExpr).Put(f, e) }
	// Variable types
	VT = func(f string, ty Type) IdToTy { return make(IdToTy).Put(f, ty) }
	ID = NewId

	// Arith
	ADD   = BIN(PlusOp)
	MINUS = BIN(MinusOp)
	MUL   = BIN(MultOp)
	DIV   = BIN(DivOp)

	NEG = UNI(MinusOp)
	NOT = UNI(NotOp)
)

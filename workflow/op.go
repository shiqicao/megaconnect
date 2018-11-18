// Copyright 2018 @ MegaSpace

// The MegaSpace source code is distributed in the hope that it will be useful,
// but WITHOUT ANY WARRANTY; without even the implied warranty of
// MERCHANTABILITY or FITNESS FOR A PARTICULAR PURPOSE. See the
// GNU Lesser General Public License for more details.

// You should have received a copy of the GNU Lesser General Public License
// along with the MegaSpace source code. If not, see <http://www.gnu.org/licenses/>.

package workflow

import p "github.com/megaspacelab/megaconnect/prettyprint"

// UniOp represents unary operators like "!", "-", etc
type UniOp struct {
	expr
	op      Operator
	operant Expr
}

// NewUniOp creates a new UniOp
func NewUniOp(op Operator, operant Expr) *UniOp {
	return &UniOp{
		op:      op,
		operant: operant,
	}
}

// Operant returns operant in UniOp
func (u *UniOp) Operant() Expr { return u.operant }

// Op returns operator
func (u *UniOp) Op() Operator { return u.op }

// Equal returns true if two expressions are the same
func (u *UniOp) Equal(expr Expr) bool {
	y, ok := expr.(*UniOp)
	return ok && u.Op() == y.Op() && u.Operant().Equal(y.Operant())
}

// Print pretty prints code
func (u *UniOp) Print() p.PrinterOp {
	operant := u.operant.Print()
	if _, ok := u.operant.(*BinOp); ok {
		operant = paren(u.operant.Print())
	}
	return p.Concat(
		p.Text(u.op.String()),
		operant,
	)
}

// BinOp represents binary operators like "<", "&&", etc.
type BinOp struct {
	expr
	op    Operator
	left  Expr
	right Expr
}

// NewBinOp creates new BinOp
func NewBinOp(op Operator, left Expr, right Expr) *BinOp {
	return &BinOp{
		op:    op,
		left:  left,
		right: right,
	}
}

// Left returns left child of a binary operator
func (b *BinOp) Left() Expr { return b.left }

// Right returns right child of a binary operator
func (b *BinOp) Right() Expr { return b.right }

// Op returns operator of a binary operator
func (b *BinOp) Op() Operator { return b.op }

// Equal returns true if two expressions are the same
func (b *BinOp) Equal(expr Expr) bool {
	y, ok := expr.(*BinOp)
	return ok &&
		b.Op() == y.Op() &&
		b.Left().Equal(y.Left()) &&
		b.Right().Equal(y.Right())
}

func (b *BinOp) Print() p.PrinterOp {
	right := b.right.Print()
	if expr, ok := b.right.(*BinOp); ok && precedence[expr.op] < precedence[b.op] {
		right = paren(right)
	}
	left := b.left.Print()
	if expr, ok := b.left.(*BinOp); ok && precedence[expr.op] < precedence[b.op] {
		left = paren(left)
	}

	return p.Concat(
		left,
		printOp(b.op.String()),
		right,
	)
}

// Operator enumerates all unary and binary
type Operator uint8

const (
	// EqualOp is equal operator, "==".
	EqualOp Operator = iota

	// NotEqualOp is not equal operator, "!=".
	NotEqualOp

	// LessThanOp is less than operator, "<".
	LessThanOp

	// LessThanEqualOp is less than or equal to operator, "=<"
	LessThanEqualOp

	// GreaterThanOp is greater than operator, ">"
	GreaterThanOp

	// GreaterThanEqualOp is greater than or equal to operator, ">="
	GreaterThanEqualOp

	// AndOp is logic AND operator
	AndOp

	// NotOp is logic negation operator
	NotOp

	// OrOp is logic OR operator
	OrOp

	// PlusOp is arithmetic addition operator
	PlusOp

	// MinusOp is arithmetic minus operator
	MinusOp

	// MultOp is arithmetic multiplication operator
	MultOp

	// DivOp is arithmetic division operator
	DivOp
)

var (
	precedence = map[Operator]uint8{
		AndOp: 0,
		OrOp:  0,

		EqualOp:            1,
		NotEqualOp:         1,
		LessThanOp:         1,
		LessThanEqualOp:    1,
		GreaterThanOp:      1,
		GreaterThanEqualOp: 1,

		PlusOp:  2,
		MinusOp: 2,

		MultOp: 3,
		DivOp:  3,

		NotOp: 4,
	}
)

func (o Operator) String() string {
	switch o {
	case EqualOp:
		return "=="
	case NotEqualOp:
		return "!="
	case LessThanOp:
		return "<"
	case GreaterThanOp:
		return ">"
	case GreaterThanEqualOp:
		return ">="
	case LessThanEqualOp:
		return "<="
	case AndOp:
		return "&&"
	case OrOp:
		return "||"
	case NotOp:
		return "!"
	case PlusOp:
		return "+"
	case MinusOp:
		return "-"
	case MultOp:
		return "*"
	case DivOp:
		return "/"
	default:
		return ""
	}
}

var (
	// opToFunc maps operators to function names,
	// an operator is applicable to a type if the mapped function is defined.
	// EqualOp is an exception, all types supports it.
	opToFunc = map[Operator]string{
		LessThanOp: "LessThan",
	}
)

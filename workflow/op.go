// Copyright 2018 @ MegaSpace

// The MegaSpace source code is distributed in the hope that it will be useful,
// but WITHOUT ANY WARRANTY; without even the implied warranty of
// MERCHANTABILITY or FITNESS FOR A PARTICULAR PURPOSE. See the
// GNU Lesser General Public License for more details.

// You should have received a copy of the GNU Lesser General Public License
// along with the MegaSpace source code. If not, see <http://www.gnu.org/licenses/>.

package workflow

// UniOp represents unary operators like "!", "-", etc
type UniOp struct {
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

func (u *UniOp) String() string { return u.op.String() + u.operant.String() }

// Operant returns operant in UniOp
func (u *UniOp) Operant() Expr { return u.operant }

// Op returns operator
func (u *UniOp) Op() Operator { return u.op }

// BinOp represents binary operators like "<", "&&", etc.
type BinOp struct {
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

func (b *BinOp) String() string { return b.left.String() + b.op.String() + b.right.String() }

// Left returns left child of a binary operator
func (b *BinOp) Left() Expr { return b.left }

// Right returns right child of a binary operator
func (b *BinOp) Right() Expr { return b.right }

// Op returns operator of a binary operator
func (b *BinOp) Op() Operator { return b.op }

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

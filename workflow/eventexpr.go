// Copyright 2018 @ MegaSpace

// The MegaSpace source code is distributed in the hope that it will be useful,
// but WITHOUT ANY WARRANTY; without even the implied warranty of
// MERCHANTABILITY or FITNESS FOR A PARTICULAR PURPOSE. See the
// GNU Lesser General Public License for more details.

// You should have received a copy of the GNU Lesser General Public License
// along with the MegaSpace source code. If not, see <http://www.gnu.org/licenses/>.

package workflow

import p "github.com/megaspacelab/megaconnect/prettyprint"

// EventExpr represents event expression
type EventExpr interface {
	Node
	Equal(EventExpr) bool
}

type eventExpr struct {
	node
}

// EventExprOperator represents operator of event expression
type EventExprOperator uint8

const (
	// AndEOp represents logical AND for event expressions
	AndEOp EventExprOperator = iota

	// OrEOp represents logical OR for event expressions
	OrEOp
)

// EBinOp represents binary operators for event expression
type EBinOp struct {
	eventExpr
	op    EventExprOperator
	left  EventExpr
	right EventExpr
}

// NewEBinOp creates a new instance of EBinOp
func NewEBinOp(op EventExprOperator, left EventExpr, right EventExpr) *EBinOp {
	return &EBinOp{
		op:    op,
		left:  left,
		right: right,
	}
}

// Equal returns true if two event expressions are equivalent
func (e *EBinOp) Equal(x EventExpr) bool {
	y, ok := x.(*EBinOp)
	return ok && y.left.Equal(e.left) && y.right.Equal(e.right)
}

// Print pretty prints code
func (e *EBinOp) Print() p.PrinterOp {
	var op string
	if e.op == AndEOp {
		op = "&&"
	} else {
		op = "||"
	}
	right := e.right.Print()
	if _, ok := e.right.(*EBinOp); ok {
		right = paren(right)
	}

	return p.Concat(
		e.left.Print(),
		printOp(op),
		right,
	)
}

// EVar represents event variable
type EVar struct {
	eventExpr
	name string
}

// NewEVar returns a new instance of event variable
func NewEVar(name string) *EVar { return &EVar{name: name} }

// Equal returns true if two event expressions are equivalent
func (e *EVar) Equal(x EventExpr) bool {
	y, ok := x.(*EVar)
	return ok && e.name == y.name
}

// Print pretty prints code
func (e *EVar) Print() p.PrinterOp { return p.Text(e.name) }

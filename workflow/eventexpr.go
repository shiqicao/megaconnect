// Copyright 2018 @ MegaSpace

// The MegaSpace source code is distributed in the hope that it will be useful,
// but WITHOUT ANY WARRANTY; without even the implied warranty of
// MERCHANTABILITY or FITNESS FOR A PARTICULAR PURPOSE. See the
// GNU Lesser General Public License for more details.

// You should have received a copy of the GNU Lesser General Public License
// along with the MegaSpace source code. If not, see <http://www.gnu.org/licenses/>.

package workflow

// EventExpr represents event expression
type EventExpr interface {
	Equal(EventExpr) bool
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

// EVar represents event variable
type EVar struct{ name string }

// NewEVar returns a new instance of event variable
func NewEVar(name string) *EVar { return &EVar{name: name} }

// Equal returns true if two event expressions are equivalent
func (e *EVar) Equal(x EventExpr) bool {
	y, ok := x.(*EVar)
	return ok && e.name == y.name
}

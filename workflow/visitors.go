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
)

type EventExprVisitor struct {
	VisitVar   func(*EVar) interface{}
	VisitBinOp func(bin *EBinOp, l interface{}, r interface{}) interface{}
}

func (e *EventExprVisitor) Visit(eexpr EventExpr) interface{} {
	switch eexpr := eexpr.(type) {
	case *EVar:
		if e.VisitVar == nil {
			return nil
		}
		return e.VisitVar(eexpr)
	case *EBinOp:
		l := e.Visit(eexpr.left)
		r := e.Visit(eexpr.right)
		if e.VisitBinOp == nil {
			return nil
		}
		return e.VisitBinOp(eexpr, l, r)
	default:
		panic(fmt.Sprintf("%T not supported", eexpr))
	}
}

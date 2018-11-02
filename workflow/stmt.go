// Copyright 2018 @ MegaSpace

// The MegaSpace source code is distributed in the hope that it will be useful,
// but WITHOUT ANY WARRANTY; without even the implied warranty of
// MERCHANTABILITY or FITNESS FOR A PARTICULAR PURPOSE. See the
// GNU Lesser General Public License for more details.

// You should have received a copy of the GNU Lesser General Public License
// along with the MegaSpace source code. If not, see <http://www.gnu.org/licenses/>.

package workflow

import (
	"bytes"
	"fmt"
)

type Stmt interface {
	fmt.Stringer
	Equal(Stmt) bool
}

type Stmts []Stmt

// Copy returns a new instance of Stmts
func (s Stmts) Copy() Stmts {
	r := make(Stmts, len(s))
	copy(r, s)
	return r
}

// Equal returns true if x is the equivalent stmts
func (s Stmts) Equal(x Stmts) bool {
	if len(s) != len(x) {
		return false
	}
	for i, s := range s {
		if !s.Equal(x[i]) {
			return false
		}
	}
	return true
}

// Fire represents a fire statement
type Fire struct {
	eventName string
	eventDecl *EventDecl
	eventObj  Expr
}

// NewFire creates a new instance of Fire
func NewFire(eventName string, eventObj Expr) *Fire {
	return &Fire{
		eventName: eventName,
		eventObj:  eventObj,
	}
}

// Equal returns true if x is the same fire statement
func (f *Fire) Equal(x Stmt) bool {
	y, ok := x.(*Fire)
	return ok && y.eventName == f.eventName && y.eventObj.Equal(f.eventObj)
}

func (f *Fire) String() string {
	var buf bytes.Buffer
	buf.WriteString("fire ")
	buf.WriteString(f.eventName)
	buf.WriteString(" ")
	buf.WriteString(f.eventObj.String())
	return buf.String()
}

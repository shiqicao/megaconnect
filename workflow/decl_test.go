// Copyright 2018 @ MegaSpace

// The MegaSpace source code is distributed in the hope that it will be useful,
// but WITHOUT ANY WARRANTY; without even the implied warranty of
// MERCHANTABILITY or FITNESS FOR A PARTICULAR PURPOSE. See the
// GNU Lesser General Public License for more details.

// You should have received a copy of the GNU Lesser General Public License
// along with the MegaSpace source code. If not, see <http://www.gnu.org/licenses/>.

package workflow

import (
	"sort"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestActionDeclTriggerEvents(t *testing.T) {
	check := func(trigger EventExpr, expected []string) {
		a := NewActionDecl(ID("a"), trigger, Stmts{})
		events := a.TriggerEvents()
		sort.StringSlice(events).Sort()
		sort.StringSlice(expected).Sort()
		assert.Equal(t, expected, events)
	}

	check(EV("a"), []string{"a"})
	check(EAND(EV("a"), EV("a")), []string{"a"})
	check(EAND(EV("a"), EV("b")), []string{"a", "b"})
	check(EAND(EV("a"), EOR(EV("b"), EV("c"))), []string{"a", "b", "c"})
	check(EAND(EV("a"), EOR(EV("b"), EV("a"))), []string{"a", "b"})
}

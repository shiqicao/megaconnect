// Copyright 2018 @ MegaSpace

// The MegaSpace source code is distributed in the hope that it will be useful,
// but WITHOUT ANY WARRANTY; without even the implied warranty of
// MERCHANTABILITY or FITNESS FOR A PARTICULAR PURPOSE. See the
// GNU Lesser General Public License for more details.

// You should have received a copy of the GNU Lesser General Public License
// along with the MegaSpace source code. If not, see <http://www.gnu.org/licenses/>.

package types

import intpr "github.com/megaspacelab/eventmanager/interpreter"

type Monitor struct {
	condition *intpr.Expr
	// TODO: Define trigger
}

func NewMonitor(condition *intpr.Expr) (*Monitor, error) {
	return &Monitor{
		condition: condition,
	}, nil
}

func (m *Monitor) Condition() *intpr.Expr { return m.condition }
func (m *Monitor) EvalCondition(i *intpr.Interpreter) (*intpr.Const, error) {
	return i.EvalExpr(m.condition)
}

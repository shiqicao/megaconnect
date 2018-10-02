// Copyright 2018 @ MegaSpace

// The MegaSpace source code is distributed in the hope that it will be useful,
// but WITHOUT ANY WARRANTY; without even the implied warranty of
// MERCHANTABILITY or FITNESS FOR A PARTICULAR PURPOSE. See the
// GNU Lesser General Public License for more details.

// You should have received a copy of the GNU Lesser General Public License
// along with the MegaSpace source code. If not, see <http://www.gnu.org/licenses/>.

package types

import wfl "github.com/megaspacelab/eventmanager/workflow"

// Monitor specifies how an event will be created, it defines monitor frequency and conditions
type Monitor struct {
	condition wfl.Expr
	// TODO: Define trigger
}

// NewMonitor creates a new Monitor object
func NewMonitor(condition wfl.Expr) (*Monitor, error) {
	return &Monitor{
		condition: condition,
	}, nil
}

// Condition returns the condition an event will be fired
func (m *Monitor) Condition() wfl.Expr { return m.condition }

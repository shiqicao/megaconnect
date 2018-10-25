// Copyright 2018 @ MegaSpace

// The MegaSpace source code is distributed in the hope that it will be useful,
// but WITHOUT ANY WARRANTY; without even the implied warranty of
// MERCHANTABILITY or FITNESS FOR A PARTICULAR PURPOSE. See the
// GNU Lesser General Public License for more details.

// You should have received a copy of the GNU Lesser General Public License
// along with the MegaSpace source code. If not, see <http://www.gnu.org/licenses/>.

package workflow

// StmtResult is a result of execution of a statement in workflow lang
type StmtResult interface{}

// FireEvent represents the result of Fire statement in workflow lang
type FireEvent struct {
	eventName string
	payload   *ObjConst
}

// NewFireEvent creates a new instance of FireEvent
func NewFireEvent(eventName string, payload *ObjConst) *FireEvent {
	return &FireEvent{
		eventName: eventName,
		payload:   payload,
	}
}

// EventName returns the name of event being fired
func (f *FireEvent) EventName() string { return f.eventName }

// Payload returns payload of an event
func (f *FireEvent) Payload() *ObjConst { return f.payload }

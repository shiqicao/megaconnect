// Copyright 2018 @ MegaSpace

// The MegaSpace source code is distributed in the hope that it will be useful,
// but WITHOUT ANY WARRANTY; without even the implied warranty of
// MERCHANTABILITY or FITNESS FOR A PARTICULAR PURPOSE. See the
// GNU Lesser General Public License for more details.

// You should have received a copy of the GNU Lesser General Public License
// along with the MegaSpace source code. If not, see <http://www.gnu.org/licenses/>.

package utils

import (
	"fmt"

	wf "github.com/megaspacelab/megaconnect/workflow"
)

type ParserError interface {
	Pos() wf.Pos
	error
}

type ErrDupDef struct{ Id *wf.Id }

func (e *ErrDupDef) Pos() wf.Pos   { return *(e.Id.Pos()) }
func (e *ErrDupDef) Error() string { return fmt.Sprintf("%s is already defined", e.Id.Id()) }

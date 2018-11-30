// Copyright 2018 @ MegaSpace

// The MegaSpace source code is distributed in the hope that it will be useful,
// but WITHOUT ANY WARRANTY; without even the implied warranty of
// MERCHANTABILITY or FITNESS FOR A PARTICULAR PURPOSE. See the
// GNU Lesser General Public License for more details.

// You should have received a copy of the GNU Lesser General Public License
// along with the MegaSpace source code. If not, see <http://www.gnu.org/licenses/>.

package utils

import (
	"bytes"
	"fmt"
	"strings"

	wf "github.com/megaspacelab/megaconnect/workflow"
	"github.com/megaspacelab/megaconnect/workflow/parser/gen/errors"
)

// ErrGocc wraps gocc generated error.
type ErrGocc struct {
	Err *errors.Error
}

// Pos returns error position
func (g *ErrGocc) Pos() *wf.Pos {
	if g.Err.ErrorToken != nil {
		p := TokenToPos(g.Err.ErrorToken)
		return &p
	}
	return nil
}

// SetPos is noop, just implement HasPos
func (g *ErrGocc) SetPos(pos *wf.Pos) {}

func (g *ErrGocc) Error() string {
	var buf bytes.Buffer
	fmt.Fprintf(&buf, "Expected one of: ")
	for i, sym := range g.Err.ExpectedTokens {
		s := sym
		if len(sym) > 2 && sym[0:2] == "kd" {
			s = strings.ToLower(sym[2:])
		}
		fmt.Fprintf(&buf, "%s", s)
		if i != len(g.Err.ExpectedTokens) {
			fmt.Fprintf(&buf, ", ")
		}
	}
	return buf.String()
}

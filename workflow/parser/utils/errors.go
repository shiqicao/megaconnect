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
	"io"
	"strings"

	wf "github.com/megaspacelab/megaconnect/workflow"
	"github.com/megaspacelab/megaconnect/workflow/parser/gen/errors"
)

// ParserError represents an error from workflow lang parser
type ParserError interface {
	Pos() wf.Pos
	error
}

// ErrDupDef is returned if duplicated definition occurs
type ErrDupDef struct{ Id *wf.Id }

// Pos returns error position
func (e *ErrDupDef) Pos() wf.Pos { return *(e.Id.Pos()) }

func (e *ErrDupDef) Error() string {
	return withPos(e, func(w io.Writer) {
		fmt.Fprintf(w, "%s is already defined", e.Id.Id())
	})
}

// ErrGocc wraps gocc generated error.
type ErrGocc struct {
	Err *errors.Error
}

// Pos returns error position
func (g *ErrGocc) Pos() (p wf.Pos) {
	if g.Err.ErrorToken != nil {
		p = TokenToPos(g.Err.ErrorToken)
	}
	return
}

func (g *ErrGocc) Error() string {
	return withPos(g, func(w io.Writer) {
		fmt.Fprintf(w, "Expected one of: ")
		for i, sym := range g.Err.ExpectedTokens {
			s := sym
			if len(sym) > 2 && sym[0:2] == "kd" {
				s = strings.ToLower(sym[2:])
			}
			fmt.Fprintf(w, "%s", s)
			if i != len(g.Err.ExpectedTokens) {
				fmt.Fprintf(w, ", ")
			}
		}
	})
}

func withPos(pos interface{ Pos() wf.Pos }, f func(w io.Writer)) string {
	b := new(bytes.Buffer)
	p := pos.Pos()
	fmt.Fprintf(b, "Error line %d, col %d: ", p.StartRow, p.StartCol)
	f(b)
	return b.String()
}

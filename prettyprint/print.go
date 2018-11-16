// Copyright 2018 @ MegaSpace

// The MegaSpace source code is distributed in the hope that it will be useful,
// but WITHOUT ANY WARRANTY; without even the implied warranty of
// MERCHANTABILITY or FITNESS FOR A PARTICULAR PURPOSE. See the
// GNU Lesser General Public License for more details.

// You should have received a copy of the GNU Lesser General Public License
// along with the MegaSpace source code. If not, see <http://www.gnu.org/licenses/>.

package prettyprint

import (
	"fmt"
	"io"
)

// PrinterOp represents an printer operation
type PrinterOp func(Printer) error

// Nil is a noop printer operator
func Nil() PrinterOp { return nil }

// Text wraps text to a printer operator
func Text(s string) PrinterOp { return func(p Printer) error { return p.Write(s) } }

// Line represents a new line operator
func Line() PrinterOp { return func(p Printer) error { return p.NewLine() } }

// Nest pushes a printing operator to a sub scope
func Nest(i int, op PrinterOp) PrinterOp {
	return func(p Printer) error {
		p.IndentInc(i)
		defer p.IndentDec(i)
		return Concat(Line(), op)(p)
	}
}

// Concat joins two printing operators
func Concat(ops ...PrinterOp) PrinterOp {
	return func(p Printer) error {
		for _, op := range ops {
			if op == nil {
				continue
			}
			e := op(p)
			if e != nil {
				return e
			}
		}
		return nil
	}
}

// Printer is an abstract printing receiver, used by PrinterOp
type Printer interface {
	Write(s string) error
	NewLine() error
	IndentInc(i int)
	IndentDec(i int)
}

// TxtPrinter implements Printer
type TxtPrinter struct {
	w       io.Writer
	setting *Style
	indent  int
}

// NewTxtPrinter creates a new instance of TxtPrinter
func NewTxtPrinter(w io.Writer) *TxtPrinter {
	return &TxtPrinter{
		w:       w,
		setting: defaultSetting,
	}
}

// Write prints a string
func (p *TxtPrinter) Write(s string) error {
	_, err := fmt.Fprint(p.w, s)
	return err
}

// NewLine prints a new line
func (p *TxtPrinter) NewLine() error {
	_, err := fmt.Fprint(p.w, p.setting.newline)
	if err != nil {
		return err
	}
	for i := 0; i < p.indent; i++ {
		_, err := fmt.Fprint(p.w, p.setting.indent)
		if err != nil {
			return err
		}
	}
	return nil
}

// IndentInc increases indent level
func (p *TxtPrinter) IndentInc(i int) { p.indent = p.indent + i }

// IndentDec decreases indent level
func (p *TxtPrinter) IndentDec(i int) { p.indent = p.indent - i }

var (
	defaultSetting = &Style{
		newline: "\n",
		indent:  "  ",
	}
)

// Style provides a configurable settings for TxtPrinter
type Style struct {
	newline string
	indent  string
}

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

var (
	defaultSetting = &Style{
		newline: "\n",
		indent:  "  ",
	}
)

type Style struct {
	newline string
	indent  string
}

type TxtPrinter struct {
	w       io.Writer
	setting *Style
	indent  int
}

func NewTxtPrinter(w io.Writer) *TxtPrinter {
	return &TxtPrinter{
		w:       w,
		setting: defaultSetting,
	}
}

func (p *TxtPrinter) Write(s string) error {
	_, err := fmt.Fprint(p.w, s)
	return err
}
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

func (p *TxtPrinter) IndentInc(i int) { p.indent = p.indent + i }
func (p *TxtPrinter) IndentDec(i int) { p.indent = p.indent - i }

type Printer interface {
	Write(s string) error
	NewLine() error
	IndentInc(i int)
	IndentDec(i int)
}

type PrinterOp func(Printer) error

func Nil() PrinterOp { return nil }

func Text(s string) PrinterOp { return func(p Printer) error { return p.Write(s) } }

func Line() PrinterOp { return func(p Printer) error { return p.NewLine() } }

func Nest(i int, op PrinterOp) PrinterOp {
	return func(p Printer) error {
		p.IndentInc(i)
		defer p.IndentDec(i)
		return Concat(Line(), op)(p)
	}
}

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

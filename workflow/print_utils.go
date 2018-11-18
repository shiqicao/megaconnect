// Copyright 2018 @ MegaSpace

// The MegaSpace source code is distributed in the hope that it will be useful,
// but WITHOUT ANY WARRANTY; without even the implied warranty of
// MERCHANTABILITY or FITNESS FOR A PARTICULAR PURPOSE. See the
// GNU Lesser General Public License for more details.

// You should have received a copy of the GNU Lesser General Public License
// along with the MegaSpace source code. If not, see <http://www.gnu.org/licenses/>.

package workflow

import p "github.com/megaspacelab/megaconnect/prettyprint"

func printOp(op string) p.PrinterOp {
	return p.Concat(
		p.Text(" "),
		p.Text(op),
		p.Text(" "),
	)
}

func separatedBy(ops []p.PrinterOp, separator p.PrinterOp) p.PrinterOp {
	body := []p.PrinterOp{}
	last := len(ops) - 1
	for i, op := range ops {
		body = append(body, op)
		if i != last {
			body = append(body, separator)
		}
	}
	return p.Concat(body...)
}

func paren(op p.PrinterOp) p.PrinterOp {
	return p.Concat(
		p.Text("("),
		op,
		p.Text(")"),
	)
}

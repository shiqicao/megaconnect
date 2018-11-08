// Copyright 2018 @ MegaSpace

// The MegaSpace source code is distributed in the hope that it will be useful,
// but WITHOUT ANY WARRANTY; without even the implied warranty of
// MERCHANTABILITY or FITNESS FOR A PARTICULAR PURPOSE. See the
// GNU Lesser General Public License for more details.

// You should have received a copy of the GNU Lesser General Public License
// along with the MegaSpace source code. If not, see <http://www.gnu.org/licenses/>.

package workflow

// Pos maps node to source code
type Pos struct {
	StartRow int
	StartCol int
	EndRow   int
	EndCol   int
}

// Node represents a tree node in AST
type Node interface {

	// Pos returns Pos or nil if a mapping does not exist
	Pos() *Pos

	SetPos(startRow int, startCol int, endRow int, endCol int)
}

type node struct {
	pos *Pos
}

func (n *node) SetPos(startRow int, startCol int, endRow int, endCol int) {
	pos := Pos{
		StartRow: startRow,
		StartCol: startCol,
		EndRow:   endRow,
		EndCol:   endCol,
	}
	n.pos = &pos
}

func (n *node) setPos(pos *Pos) {
	if pos == nil {
		n.pos = nil
		return
	}
	p := *pos
	n.pos = &p
}

func (n *node) Pos() *Pos {
	if n.pos == nil {
		return nil
	}
	r := *n.pos
	return &r
}

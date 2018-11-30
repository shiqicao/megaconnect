// Copyright 2018 @ MegaSpace

// The MegaSpace source code is distributed in the hope that it will be useful,
// but WITHOUT ANY WARRANTY; without even the implied warranty of
// MERCHANTABILITY or FITNESS FOR A PARTICULAR PURPOSE. See the
// GNU Lesser General Public License for more details.

// You should have received a copy of the GNU Lesser General Public License
// along with the MegaSpace source code. If not, see <http://www.gnu.org/licenses/>.

package workflow

import (
	p "github.com/megaspacelab/megaconnect/prettyprint"
)

// Node represents a tree node in AST
type Node interface {
	HasPos
	Print() p.PrinterOp
	Children() []Node
}

type node struct {
	hasPos
}

// PrintNode returns string format of a node
func PrintNode(node Node) string {
	return p.String(node.Print())
}

func NodeWalker(root Node, visitor func(n Node)) {
	visitor(root)
	for _, child := range root.Children() {
		NodeWalker(child, visitor)
	}
}

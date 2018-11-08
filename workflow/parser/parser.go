// Copyright 2018 @ MegaSpace

// The MegaSpace source code is distributed in the hope that it will be useful,
// but WITHOUT ANY WARRANTY; without even the implied warranty of
// MERCHANTABILITY or FITNESS FOR A PARTICULAR PURPOSE. See the
// GNU Lesser General Public License for more details.

// You should have received a copy of the GNU Lesser General Public License
// along with the MegaSpace source code. If not, see <http://www.gnu.org/licenses/>.

package parser

import (
	"fmt"

	wf "github.com/megaspacelab/megaconnect/workflow"
	"github.com/megaspacelab/megaconnect/workflow/parser/gen/lexer"
	p "github.com/megaspacelab/megaconnect/workflow/parser/gen/parser"
)

// Parse wraps generated parser and improve error message
func Parse(src string) (*wf.WorkflowDecl, error) {
	lexer, err := lexer.NewLexerFile(src)
	if err != nil {
		return nil, err
	}
	ast, err := p.NewParser().Parse(lexer)
	if err != nil {
		return nil, translateErr(err)
	}
	w, ok := ast.(*wf.WorkflowDecl)
	if !ok {
		return nil, fmt.Errorf("Failed to parser workflow, parser returns %T ", ast)
	}
	return w, nil
}

func translateErr(err error) error {
	// TODO: improve error msg
	return err
}

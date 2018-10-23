// Copyright 2018 @ MegaSpace

// The MegaSpace source code is distributed in the hope that it will be useful,
// but WITHOUT ANY WARRANTY; without even the implied warranty of
// MERCHANTABILITY or FITNESS FOR A PARTICULAR PURPOSE. See the
// GNU Lesser General Public License for more details.

// You should have received a copy of the GNU Lesser General Public License
// along with the MegaSpace source code. If not, see <http://www.gnu.org/licenses/>.

package parser

import (
	"testing"

	"github.com/stretchr/testify/assert"

	wf "github.com/megaspacelab/megaconnect/workflow"
	"github.com/megaspacelab/megaconnect/workflow/parser/goccgen/lexer"

	"github.com/megaspacelab/megaconnect/workflow/parser/goccgen/parser"
)

func TestParser(t *testing.T) {
	r, err := parse(t, `
	monitor a 
	  chain Eth 
	  condition true
	  var {
		  a = true
		  b = false
	  }
	`)
	assert.NoError(t, err)
	md, ok := r.(*wf.MonitorDecl)
	assert.True(t, ok)
	assert.NotNil(t, md)
	assert.True(t, md.Equal(wf.NewMonitorDecl("a", wf.TrueConst, wf.VarDecls{"a": wf.TrueConst, "b": wf.FalseConst})))

}

func parse(t *testing.T, code string) (interface{}, error) {
	s := lexer.NewLexer([]byte(code))
	p := parser.NewParser()
	return p.Parse(s)
}

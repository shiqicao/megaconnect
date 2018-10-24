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
	"github.com/megaspacelab/megaconnect/workflow/parser/goccgen/token"
)

func TermBoolLitAction(t interface{}) (*wf.BoolConst, error) {
	lit := string(t.(*token.Token).Lit)
	if lit == "true" {
		return wf.TrueConst, nil
	} else if lit == "false" {
		return wf.FalseConst, nil
	}
	return nil, fmt.Errorf("")
}

func MonitorAction(name interface{}, chain interface{}, expr interface{}, varDecls interface{}) (*wf.MonitorDecl, error) {
	if varDecls == nil {
		varDecls = wf.VarDecls{}
	} else {
		varDecls = varDecls.(wf.VarDecls)
	}
	md := wf.NewMonitorDecl(string(name.(*token.Token).Lit), expr.(wf.Expr), varDecls.(wf.VarDecls))
	return md, nil
}

func VarDeclAction(nameRaw interface{}, exprRaw interface{}) (wf.VarDecls, error) {
	vd := make(wf.VarDecls)
	name := string(nameRaw.(*token.Token).Lit)
	expr := exprRaw.(wf.Expr)
	vd[name] = expr
	return vd, nil
}

func VarDeclsAction(varDeclsRaw interface{}, varDeclRaw interface{}) (wf.VarDecls, error) {
	varDecls := varDeclsRaw.(wf.VarDecls)
	if varDeclRaw == nil {
		return varDecls, nil
	}
	varDecl := varDeclRaw.(wf.VarDecls)
	// varDecl should only contain only one item
	for name, expr := range varDecl {
		if _, found := varDecls[name]; found {
			return nil, fmt.Errorf("")
		}
		varDecls[name] = expr
	}
	return varDecls, nil
}

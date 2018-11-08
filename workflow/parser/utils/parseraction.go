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
	"math/big"

	wf "github.com/megaspacelab/megaconnect/workflow"
	"github.com/megaspacelab/megaconnect/workflow/parser/gen/token"
)

func IntLitAction(t interface{}) (*wf.IntConst, error) {
	lit := Lit(t)
	i, ok := new(big.Int).SetString(lit, 10)
	if !ok || i == nil {
		return nil, fmt.Errorf("Failed to parse int %s", lit)
	}
	return wf.NewIntConst(i), nil
}

func BoolLitAction(t interface{}) (*wf.BoolConst, error) {
	lit := Lit(t)
	if lit == "true" {
		return wf.TrueConst, nil
	} else if lit == "false" {
		return wf.FalseConst, nil
	}
	return nil, fmt.Errorf("")
}

func StrLitAction(t interface{}) (*wf.StrConst, error) {
	lit := Lit(t)
	if len(lit) < 2 {
		return nil, fmt.Errorf("")
	}
	return wf.NewStrConst(lit[1 : len(lit)-1]), nil
}

func MonitorAction(
	name interface{}, chain interface{}, expr interface{},
	varDeclsRaw interface{}, event interface{}, eventObj interface{},
) (*wf.MonitorDecl, error) {
	var varDecls wf.IdToExpr
	if varDeclsRaw == nil {
		varDecls = wf.NewIdToExpr()
	} else {
		varDecls = varDeclsRaw.(wf.IdToExpr)
	}
	fire := wf.NewFire(Lit(event), wf.NewObjLit(eventObj.(wf.IdToExpr)))
	md := wf.NewMonitorDecl(Id(name), expr.(wf.Expr), varDecls, fire, Lit(chain))
	return md, nil
}

func VarDeclAction(nameRaw interface{}, exprRaw interface{}) (wf.IdToExpr, error) {
	vd := wf.NewIdToExpr()
	name := Id(nameRaw)
	expr := exprRaw.(wf.Expr)
	vd.Add(name, expr)
	return vd, nil
}

func VarDeclsAction(varDeclsRaw interface{}, varDeclRaw interface{}) (wf.IdToExpr, error) {
	varDecls := varDeclsRaw.(wf.IdToExpr)
	if varDeclRaw == nil {
		return varDecls, nil
	}
	varDecl := varDeclRaw.(wf.IdToExpr)
	// varDecl should only contain only one item
	for _, decl := range varDecl {
		if ok := varDecls.Add(decl.Id, decl.Expr); !ok {
			return nil, &ErrDupDef{Id: decl.Id}
		}
	}
	return varDecls, nil
}

func ObjLitFieldsAction(varDeclsRaw interface{}, varDeclRaw interface{}) (wf.IdToExpr, error) {
	varDecls := varDeclsRaw.(wf.IdToExpr)
	if varDeclRaw == nil {
		return varDecls, nil
	}
	varDecl := varDeclRaw.([]interface{})
	field := Id(varDecl[0])
	expr := varDecl[1].(wf.Expr)
	if !varDecls.Add(field, expr) {
		return nil, &ErrDupDef{Id: field}
	}
	return varDecls, nil
}

func ObjAccessorAction(expr interface{}, id interface{}) (*wf.ObjAccessor, error) {
	receiver := expr.(wf.Expr)
	accessor := Lit(id)
	return wf.NewObjAccessor(receiver, accessor), nil
}

func ObjLitAction(objFields interface{}) (*wf.ObjLit, error) {
	fields := objFields.(wf.IdToExpr)
	return wf.NewObjLit(fields), nil
}

func Lit(t interface{}) string {
	return string(t.(*token.Token).Lit)
}

func Id(t interface{}) *wf.Id {
	token := t.(*token.Token)
	id := string(token.Lit)
	wfid := wf.NewId(id)
	wfid.SetPos(token.Line, token.Column, token.Line, token.Column+token.Offset)
	return wfid
}

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

// IntLitAction is mapped to intLit in lang.bnf
func IntLitAction(t interface{}) (*wf.IntConst, error) {
	lit := Lit(t)
	i, ok := new(big.Int).SetString(lit, 10)
	if !ok || i == nil {
		return nil, fmt.Errorf("Failed to parse int %s", lit)
	}
	return wf.NewIntConst(i), nil
}

// BoolLitAction is mapped to boolLit in lang.bnf
func BoolLitAction(t interface{}) (*wf.BoolConst, error) {
	lit := Lit(t)
	if lit == "true" {
		return wf.TrueConst, nil
	} else if lit == "false" {
		return wf.FalseConst, nil
	}
	return nil, fmt.Errorf("")
}

// StrLitAction is mapped to stringLit in lang.bnf
func StrLitAction(t interface{}) (*wf.StrConst, error) {
	lit := Lit(t)
	if len(lit) < 2 {
		return nil, fmt.Errorf("")
	}
	return wf.NewStrConst(lit[1 : len(lit)-1]), nil
}

// MonitorAction is mapped to Monitor in lang.bnf
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

// VarDeclAction is mapped to VarDecl in lang.bnf
func VarDeclAction(nameRaw interface{}, exprRaw interface{}) (wf.IdToExpr, error) {
	vd := wf.NewIdToExpr()
	name := Id(nameRaw)
	expr := exprRaw.(wf.Expr)
	vd.Add(name, expr)
	return vd, nil
}

// VarDeclsAction is mapped to VarDecls in lang.bnf
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

// ObjLitFieldsAction is mapped to ObjLitFields_ in lang.bnf
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

// ObjFieldsAction is mapped to ObjFields_ in lang.bnf
func ObjFieldsAction(fieldDeclsRaw interface{}, fieldDeclRaw interface{}) (wf.IdToTy, error) {
	fieldDecls := fieldDeclsRaw.(wf.IdToTy)
	if fieldDeclRaw == nil {
		return fieldDecls, nil
	}
	fieldDecl := fieldDeclRaw.([]interface{})
	field := Id(fieldDecl[0])
	ty := fieldDecl[1].(wf.Type)
	if !fieldDecls.Add(field, ty) {
		return nil, &ErrDupDef{Id: field}
	}
	return fieldDecls, nil
}

// ObjAccessorAction is referenced in Expr in lang.bnf
func ObjAccessorAction(expr interface{}, id interface{}) (*wf.ObjAccessor, error) {
	receiver := expr.(wf.Expr)
	accessor := Lit(id)
	return wf.NewObjAccessor(receiver, accessor), nil
}

// ObjLitAction is mapped to ObjLit in lang.bnf
func ObjLitAction(objFields interface{}) (*wf.ObjLit, error) {
	fields := objFields.(wf.IdToExpr)
	return wf.NewObjLit(fields), nil
}

func VarAction(t interface{}) (*wf.Var, error) {
	return idAction(t, func(s string) wf.Node { return wf.NewVar(s) }).(*wf.Var), nil
}

func EVarAction(t interface{}) (*wf.EVar, error) {
	return idAction(t, func(s string) wf.Node { return wf.NewEVar(s) }).(*wf.EVar), nil
}

func BinOpAction(op wf.Operator, leftRaw interface{}, rightRaw interface{}) (wf.Expr, error) {
	left := leftRaw.(wf.Expr)
	right := rightRaw.(wf.Expr)
	bin := wf.NewBinOp(op, left, right)
	mergePos(bin, left, right)
	return bin, nil
}

func EBinOpAction(op wf.EventExprOperator, leftRaw interface{}, rightRaw interface{}) (wf.EventExpr, error) {
	left := leftRaw.(wf.EventExpr)
	right := rightRaw.(wf.EventExpr)
	bin := wf.NewEBinOp(op, left, right)
	mergePos(bin, left, right)
	return bin, nil
}

// Lit converts a token to string
func Lit(t interface{}) string {
	return string(t.(*token.Token).Lit)
}

// Id converts a token to workflow.Id
func Id(t interface{}) *wf.Id {
	return idAction(t, func(s string) wf.Node { return wf.NewId(s) }).(*wf.Id)
}

func idAction(t interface{}, builder func(s string) wf.Node) wf.Node {
	token := t.(*token.Token)
	id := string(token.Lit)
	n := builder(id)
	setPos(n, token)
	return n
}

func setPos(node wf.Node, token *token.Token) {
	node.SetPos(token.Line, token.Column, token.Line, token.Column+token.Offset)
}

func mergePos(node wf.Node, start wf.Node, end wf.Node) {
	s := start.Pos()
	e := end.Pos()
	if s == nil || e == nil {
		return
	}
	node.SetPos(s.StartRow, s.StartCol, e.EndRow, e.EndCol)
}

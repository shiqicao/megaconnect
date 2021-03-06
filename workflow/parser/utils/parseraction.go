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

// RatLitAction is mapped to ratLit in lang.bnf
func RatLitAction(t interface{}) (*wf.RatConst, error) {
	lit := Lit(t)
	r, ok := new(big.Rat).SetString(lit)
	if !ok || r == nil {
		return nil, fmt.Errorf("Failed to parse rat %s", lit)
	}
	return wf.NewRatConst(r), nil
}

// BoolLitAction is mapped to boolLit in lang.bnf
func BoolLitAction(t interface{}) (*wf.BoolConst, error) {
	lit := Lit(t)
	if lit == "true" {
		return wf.TrueConst, nil
	} else if lit == "false" {
		return wf.FalseConst, nil
	}
	return nil, fmt.Errorf("Failed to parse bool %s", lit)
}

// StrLitAction is mapped to stringLit in lang.bnf
func StrLitAction(t interface{}) (*wf.StrConst, error) {
	lit := Lit(t)
	if len(lit) < 2 {
		return nil, fmt.Errorf("Failed to parse str %s", lit)
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

// VarDeclsAction is mapped to VarDecls in lang.bnf
func VarDeclsAction(varDeclsRaw interface{}, varDeclRaw interface{}) (wf.IdToExpr, error) {
	varDecls := varDeclsRaw.(wf.IdToExpr)
	varDecl := varDeclRaw.([]interface{})
	id := Id(varDecl[0])
	expr := varDecl[1].(wf.Expr)
	if ok := varDecls.Add(id, expr); !ok {
		return nil, wf.SetErrPos(&wf.ErrDupNames{Name: id.Id()}, id)
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
		return nil, wf.SetErrPos(&wf.ErrDupNames{Name: field.Id()}, field)
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
		return nil, wf.SetErrPos(&wf.ErrDupNames{Name: field.Id()}, field)
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

// VarAction is mapped to variable action in lang.bnf
func VarAction(t interface{}) (*wf.Var, error) {
	return idAction(t, func(s string) wf.Node { return wf.NewVar(s) }).(*wf.Var), nil
}

// EVarAction is mapped to event variable action in lang.bnf
func EVarAction(t interface{}) (*wf.EVar, error) {
	return idAction(t, func(s string) wf.Node { return wf.NewEVar(s) }).(*wf.EVar), nil
}

// BinOpAction is mapped to all binary operator actions in lang.bnf
func BinOpAction(op wf.Operator, leftRaw interface{}, rightRaw interface{}) (wf.Expr, error) {
	left := leftRaw.(wf.Expr)
	right := rightRaw.(wf.Expr)
	bin := wf.NewBinOp(op, left, right)
	mergePosByNodes(bin, left, right)
	return bin, nil
}

// UniOpAction is mapped to all unary operator actions in lang.bnf
func UniOpAction(op wf.Operator, operantRaw interface{}, start wf.Pos) (wf.Expr, error) {
	operant := operantRaw.(wf.Expr)
	uni := wf.NewUniOp(op, operant)
	mergePos(uni, &start, operant.Pos())
	return uni, nil
}

// EBinOpAction is mapped to all event binary operator actions in lang.bnf
func EBinOpAction(op wf.EventExprOperator, leftRaw interface{}, rightRaw interface{}) (wf.EventExpr, error) {
	left := leftRaw.(wf.EventExpr)
	right := rightRaw.(wf.EventExpr)
	bin := wf.NewEBinOp(op, left, right)
	mergePosByNodes(bin, left, right)
	return bin, nil
}

// FireAction is mapped to fire statement action in lang.bnf
func FireAction(id interface{}, eventObj interface{}, start wf.Pos) (*wf.Fire, error) {
	eventName := Lit(id)
	obj := eventObj.(wf.Expr)
	fire := wf.NewFire(eventName, obj)
	mergePos(fire, &start, obj.Pos())
	return fire, nil
}

// Lit converts a token to string
func Lit(t interface{}) string {
	return string(t.(*token.Token).Lit)
}

// Id converts a token to workflow.Id
func Id(t interface{}) *wf.Id {
	return idAction(t, func(s string) wf.Node { return wf.NewId(s) }).(*wf.Id)
}

// TokenToPos extracts wf.Pos from token
func TokenToPos(t interface{}) (p wf.Pos) {
	token := t.(*token.Token)
	if token == nil {
		return
	}
	p.StartRow = token.Line
	p.StartCol = token.Column
	p.EndRow = token.Line
	p.EndCol = token.Column + len(token.Lit)
	return
}

// ActionAction is mapped to Action declaration in lang.bnf
func ActionAction(id interface{}, eexpr interface{}, stmts interface{}, start interface{}, end interface{}) (*wf.ActionDecl, error) {
	action := wf.NewActionDecl(Id(id), eexpr.(wf.EventExpr), stmts.(wf.Stmts))
	s, e := TokenToPos(start), TokenToPos(end)
	mergePos(action, &s, &e)
	return action, nil
}

// FuncCallAction is mapped to funcation call expression in lang.bnf
func FuncCallAction(nsRaw interface{}, id interface{}, argsRaw interface{}) (*wf.FuncCall, error) {
	var ns wf.NamespacePrefix
	if nsRaw != nil {
		ns = nsRaw.(wf.NamespacePrefix)
	}
	name := Id(id)
	args := argsRaw.(wf.Args)
	return wf.NewFuncCall(ns, name, args...), nil
}

// PropsAction is mapped to props operator in lang.bnf
func PropsAction(id interface{}, start interface{}, end interface{}) (*wf.Props, error) {
	v, err := VarAction(id)
	if err != nil {
		return nil, err
	}
	props := wf.NewProps(v)
	s, e := TokenToPos(start), TokenToPos(end)
	mergePos(props, &s, &e)
	return wf.NewProps(v), nil
}

func idAction(t interface{}, builder func(s string) wf.Node) wf.Node {
	token := t.(*token.Token)
	id := string(token.Lit)
	n := builder(id)
	setPos(n, token)
	return n
}

func setPos(node wf.Node, token *token.Token) {
	pos := TokenToPos(token)
	node.SetPos(&pos)
}

func mergePos(node wf.Node, s *wf.Pos, e *wf.Pos) {
	if s == nil || e == nil {
		return
	}
	node.SetPos(&wf.Pos{
		StartRow: s.StartRow,
		StartCol: s.StartCol,
		EndRow:   e.EndRow,
		EndCol:   e.EndCol,
	})
}

func mergePosByNodes(node wf.Node, start wf.Node, end wf.Node) {
	mergePos(node, start.Pos(), end.Pos())
}

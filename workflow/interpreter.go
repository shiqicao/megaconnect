// Copyright 2018 @ MegaSpace

// The MegaSpace source code is distributed in the hope that it will be useful,
// but WITHOUT ANY WARRANTY; without even the implied warranty of
// MERCHANTABILITY or FITNESS FOR A PARTICULAR PURPOSE. See the
// GNU Lesser General Public License for more details.

// You should have received a copy of the GNU Lesser General Public License
// along with the MegaSpace source code. If not, see <http://www.gnu.org/licenses/>.

package workflow

import (
	"errors"
	"math/big"
	"reflect"

	"go.uber.org/zap"
)

// Interpreter executes a program or an expression with a given Env
type Interpreter struct {
	env *Env

	// interpreter only support single scope, variable must be unique. It does not support variable shadowing
	vars   map[string]Expr
	cache  Cache
	logger *zap.Logger
}

// NewInterpreter creates an interpreter with an Env
func NewInterpreter(env *Env, cache Cache, logger *zap.Logger) *Interpreter {
	return &Interpreter{
		env:    env,
		vars:   make(map[string]Expr),
		logger: logger,
		cache:  cache,
	}
}

// EvalMonitor evaluates a monitor. It push a scope with var declaration and pops it up after evaluation.
// It evaluates RHS of var declaration if monitor condition evaluates to true and returns variable results
// to caller.
func (i *Interpreter) EvalMonitor(monitor *MonitorDecl) (Const, map[string]Const, error) {
	// push variable declarations
	for v, expr := range monitor.vars {
		if _, ok := i.vars[v]; ok {
			return nil, nil, &ErrVarDeclaredAlready{VarName: v}
		}
		if err := i.resolveExpr(expr); err != nil {
			return nil, nil, err
		}
		i.vars[v] = expr
	}
	defer func() { i.vars = nil }()

	if err := i.resolveExpr(monitor.Condition()); err != nil {
		return nil, nil, err
	}
	result, err := i.evalExpr(monitor.Condition())
	if err != nil {
		return nil, nil, err
	}

	vars := make(map[string]Const, len(monitor.vars))
	if result.Equal(TrueConst) {
		for v := range monitor.vars {
			r, err := i.evalExpr(NewVar(v))
			if err != nil {
				return nil, nil, err
			}
			vars[v] = r
		}
	}
	return result, vars, nil
}

func (i *Interpreter) lookup(v *Var) (Const, error) {
	expr, ok := i.vars[v.Name()]
	if !ok {
		return nil, &ErrVarNotFound{VarName: v.Name()}
	}
	value, ok := expr.(Const)
	if ok {
		return value, nil
	}
	value, err := i.evalExpr(expr)
	if err != nil {
		return nil, err
	}
	i.vars[v.Name()] = value
	return value, nil
}

// EvalExpr evaluates an unbound expression. Unbound expression contains symbols(like function names) unresolved.
// It will resolve and type check before execution. It should not used for a bound expression, it will try to bind symbols again.
func (i *Interpreter) EvalExpr(expr Expr) (Const, error) {
	err := i.resolveExpr(expr)
	if err != nil {
		return nil, err
	}
	return i.evalExpr(expr)
}

func (i *Interpreter) evalExpr(expr Expr) (Const, error) {
	switch e := expr.(type) {
	case *UniOp:
		return i.evalUniOp(e)
	case *BinOp:
		return i.evalBinOp(e)
	case *Var:
		return i.lookup(e)
	case Const:
		return e, nil
	case *FuncCall:
		return i.evalFuncCall(e)
	case *ObjAccessor:
		return i.evalObjAccessor(e)
	}

	return nil, &ErrNotSupported{Name: reflect.TypeOf(expr).String()}
}

func (i *Interpreter) evalUniOp(expr *UniOp) (Const, error) {
	switch expr.Op() {
	case NotOp:
		return i.evalNot(expr)
	}
	return nil, &ErrNotSupported{Name: expr.Op().String()}
}

func (i *Interpreter) evalObjAccessor(accessor *ObjAccessor) (Const, error) {
	value, err := i.evalExpr(accessor.Receiver())
	if err != nil {
		return nil, err
	}
	obj, ok := value.(*ObjConst)
	if !ok {
		return nil, &ErrAccessorRequireObjType{Type: value.Type()}
	}
	value, ok = obj.Value()[accessor.field]
	if !ok {
		return nil, &ErrObjFieldMissing{Field: accessor.field, ObjType: obj.Type().(*ObjType)}
	}
	return value, nil
}

func (i *Interpreter) evalFuncCall(funcCall *FuncCall) (Const, error) {
	decl := funcCall.Decl()
	if decl == nil {
		return nil, &ErrSymbolNotResolved{Symbol: funcCall.Name()}
	}
	args := funcCall.Args()
	params := decl.Params()
	if len(args) != len(params) {
		return nil, nil
	}
	return i.evalFuncDecl(decl, args)
}

func (i *Interpreter) evalFuncDecl(decl *FuncDecl, args []Expr) (Const, error) {
	params := decl.Params()
	evalatedParams := make(map[string]Const, len(params))
	paramList := make([]Const, len(params))
	for j := 0; j < len(args); j++ {
		result, err := i.evalExpr(args[j])
		if err != nil {
			return nil, err
		}
		if !result.Type().Equal(params[j].Type()) {
			// This error should be already caught in type checker
			return nil, &ErrArgTypeMismatch{FuncName: decl.Name(), ParamName: params[j].Name(), ParamType: params[j].Type(), ArgType: result.Type()}
		}
		evalatedParams[params[j].Name()] = result
		paramList[j] = result
	}
	getter, setter := i.cache.funcCallCache(decl, paramList)
	result := getter()
	if result != nil {
		return result, nil
	}

	result, err := decl.evaluator(i.env, evalatedParams)
	if err != nil {
		return nil, err
	}

	if !result.Type().Equal(decl.RetType()) {
		return nil, &ErrFunRetTypeMismatch{FuncName: decl.Name(), ExpectedType: decl.RetType(), ActualType: result.Type()}
	}

	setter(result)
	return result, nil
}

func (i *Interpreter) evalBinOp(expr *BinOp) (Const, error) {
	switch expr.Op() {
	case EqualOp:
		return i.evalEqual(expr.Left(), expr.Right())
	case NotEqualOp:
		result, err := i.evalEqual(expr.Left(), expr.Right())
		if err != nil {
			return nil, err
		}
		return result.Negate(), nil
	case LessThanOp:
		return i.evalLessThan(expr.Left(), expr.Right())
	case LessThanEqualOp:
		return i.evalExpr(
			NewBinOp(OrOp,
				NewBinOp(LessThanOp, expr.Left(), expr.Right()),
				NewBinOp(EqualOp, expr.Left(), expr.Right()),
			),
		)
	case GreaterThanEqualOp:
		return i.evalExpr(
			NewUniOp(NotOp,
				NewBinOp(LessThanOp, expr.Left(), expr.Right()),
			),
		)
	case GreaterThanOp:
		return i.evalExpr(
			NewUniOp(NotOp,
				NewBinOp(LessThanEqualOp, expr.Left(), expr.Right()),
			),
		)
	// AndOp and OrOp have 'short circuit' semantics, it can not delegate operators(AND, OR) to methods in BoolType.
	// The reason is that function declaration can not capture 'short circuit' semantics
	case AndOp:
		return i.shortcircuitEval(FalseConst, expr.Left(), expr.Right())
	case OrOp:
		return i.shortcircuitEval(TrueConst, expr.Left(), expr.Right())

	// Arithmetic
	case PlusOp:
		return i.evalIntOp(expr.Left(), expr.Right(), new(big.Int).Add)
	case MinusOp:
		return i.evalIntOp(expr.Left(), expr.Right(), new(big.Int).Sub)
	case MultOp:
		return i.evalIntOp(expr.Left(), expr.Right(), new(big.Int).Mul)
	case DivOp:
		return i.evalDiv(expr)
	}
	return nil, &ErrNotSupported{Name: expr.Op().String()}
}

func (i *Interpreter) evalDiv(expr *BinOp) (ret *IntConst, err error) {
	defer func() {
		// big.Div panic if dinomirator is zero as big.Int.Div doesn't return error.
		// Interpret should return error on any invalid execution.
		if r := recover(); r != nil {
			switch r := r.(type) {
			case error:
				err = r
			case string:
				err = errors.New(r)
			default:
				err = errors.New("Unknown error")
			}
		}
	}()
	return i.evalIntOp(expr.Left(), expr.Right(), new(big.Int).Div)
}

func (i *Interpreter) evalNot(expr *UniOp) (*BoolConst, error) {
	r, err := i.evalAsBool(expr.Operant())
	if err != nil {
		return nil, err
	}
	return r.Negate(), nil
}

func (i *Interpreter) evalAsBool(expr Expr) (*BoolConst, error) {
	r, err := i.evalExpr(expr)
	if err != nil {
		return nil, err
	}
	b, ok := r.(*BoolConst)
	if !ok {
		return nil, &ErrTypeMismatch{ExpectedType: BoolType, ActualType: r.Type()}
	}
	return b, nil
}

func (i *Interpreter) evalIntOp(left Expr, right Expr, op func(*big.Int, *big.Int) *big.Int) (*IntConst, error) {
	lRet, err := i.evalExpr(left)
	if err != nil {
		return nil, err
	}
	rRet, err := i.evalExpr(right)
	if err != nil {
		return nil, err
	}
	r, ok := rRet.(*IntConst)
	if !ok {
		return nil, &ErrTypeMismatch{ExpectedType: IntType, ActualType: rRet.Type()}
	}
	l, ok := lRet.(*IntConst)
	if !ok {
		return nil, &ErrTypeMismatch{ExpectedType: IntType, ActualType: lRet.Type()}
	}
	return NewIntConst(op(l.Value(), r.Value())), nil
}

// shortcircuitEval is a helper for evaluating logical AND and OR, it does not evaluate right operant
// if result can be determined by left operant
func (i *Interpreter) shortcircuitEval(shortcut *BoolConst, exprs ...Expr) (*BoolConst, error) {
	for _, expr := range exprs {
		r, err := i.evalAsBool(expr)
		if err != nil {
			return nil, err
		}
		if r.Equal(shortcut) {
			return shortcut, nil
		}
	}
	return shortcut.Negate(), nil
}

// evalOperantsWithSameType evaluates a list of Expr and assert they have same type
func (i *Interpreter) evalExprWithSameType(exprs ...Expr) ([]Const, error) {
	result := make([]Const, len(exprs))
	var err error
	for j, expr := range exprs {
		result[j], err = i.evalExpr(expr)
		if err != nil {
			return nil, err
		}
		if j > 0 && !(result[0].Type().Equal(result[j].Type())) {
			return nil, &ErrTypeMismatch{ExpectedType: result[0].Type(), ActualType: result[j].Type()}
		}
	}
	return result, nil
}

func (i *Interpreter) evalEqual(left Expr, right Expr) (*BoolConst, error) {
	result, err := i.evalExprWithSameType(left, right)
	if err != nil {
		return nil, err
	}
	return GetBoolConst(result[0].Equal(result[1])), nil
}

func (i *Interpreter) evalOpDefinedByType(op Operator, operants ...Expr) (Const, error) {
	results, err := i.evalExprWithSameType(operants...)
	if err != nil {
		return nil, err
	}
	method := results[0].Type().Methods().find(opToFunc[op])
	if method == nil {
		return nil, &ErrNotSupported{Name: opToFunc[op]}
	}
	return i.evalFuncDecl(method, constsAsExprs(results))
}

func constsAsExprs(consts []Const) []Expr {
	exprs := make([]Expr, len(consts))
	for i, c := range consts {
		exprs[i] = c
	}
	return exprs
}

func (i *Interpreter) evalLessThan(left Expr, right Expr) (*BoolConst, error) {
	result, err := i.evalOpDefinedByType(LessThanOp, left, right)
	if err != nil {
		return nil, err
	}
	boolResult, ok := result.(*BoolConst)
	if !ok {
		return nil, &ErrTypeMismatch{ExpectedType: BoolType, ActualType: result.Type()}
	}
	return boolResult, nil
}

func (i *Interpreter) resolveExpr(expr Expr) error {
	switch e := expr.(type) {
	case *BinOp:
		err := i.resolveExpr(e.Left())
		if err != nil {
			return err
		}
		err = i.resolveExpr(e.Right())
		if err != nil {
			return err
		}
	case *UniOp:
		err := i.resolveExpr(e.Operant())
		if err != nil {
			return err
		}
	case *ObjAccessor:
		if err := i.resolveExpr(e.Receiver()); err != nil {
			return err
		}
	case *FuncCall:
		for _, arg := range e.Args() {
			if err := i.resolveExpr(arg); err != nil {
				return err
			}
		}
		return i.resolveFun(e)
	}
	return nil
}

func (i *Interpreter) resolveFun(fun *FuncCall) error {
	decl := resolveFun(i.env.prelude, fun.Name(), fun.NamespacePrefix())
	if decl == nil {
		return &ErrSymbolNotResolved{Symbol: fun.Name()}
	}
	fun.SetDecl(decl)
	return nil
}

func resolveFun(nss []*NamespaceDecl, name string, prefix NamespacePrefix) *FuncDecl {
	if len(prefix) == 0 {
		return nil
	}

	cur := prefix[0]
	var ns *NamespaceDecl
	for _, n := range nss {
		if cur == n.name {
			ns = n
			break
		}
	}
	if ns == nil {
		return nil
	}

	if len(prefix) > 1 {
		return resolveFun(ns.children, name, prefix[1:])
	}

	for _, fun := range ns.funs {
		if fun.Name() == name {
			return fun
		}
	}
	return nil
}

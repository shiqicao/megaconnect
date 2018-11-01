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
	// TODO: move resolver out of Interpreter, it will be used for static compiling as well.
	// Interpreter assume that any expression/action/workflow is sound.
	resolver *Resolver

	// interpreter only support single scope, variable must be unique. It does not support variable shadowing
	vars   map[string]Expr
	cache  Cache
	logger *zap.Logger
}

// NewInterpreter creates an interpreter with an Env
func NewInterpreter(env *Env, cache Cache, resolver *Resolver, logger *zap.Logger) *Interpreter {
	return &Interpreter{
		env:      env,
		vars:     nil,
		logger:   logger,
		resolver: resolver,
		cache:    cache,
	}
}

// EvalMonitor evaluates a monitor. It pushes a scope with var declaration and pops it up after evaluation.
// It creates a new instance of an event defined in monitor if condition evaluates to true
func (i *Interpreter) EvalMonitor(monitor *MonitorDecl) (*FireEventResult, error) {
	// push variable declarations
	i.vars = make(map[string]Expr, len(monitor.vars))
	for v, expr := range monitor.vars {
		if _, ok := i.vars[v]; ok {
			return nil, &ErrVarDeclaredAlready{VarName: v}
		}
		if err := i.resolver.resolveExpr(expr); err != nil {
			return nil, err
		}
		i.vars[v] = expr
	}
	defer func() { i.vars = nil }()

	if err := i.resolver.resolveExpr(monitor.Condition()); err != nil {
		return nil, err
	}
	if err := i.resolver.resolveExpr(monitor.event.eventObj); err != nil {
		return nil, err
	}
	result, err := i.evalExpr(monitor.Condition())
	if err != nil {
		return nil, err
	}
	if result.Equal(FalseConst) {
		return nil, nil
	}
	return i.evalFireStmt(monitor.event)
}

func (i *Interpreter) lookup(v *Var) (Const, error) {
	expr, ok := i.vars[v.Name()]
	if !ok {
		if i.env.eventStore != nil {
			return GetBoolConst(i.env.eventStore.lookup(v.name) != nil), nil
		}
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

func (i *Interpreter) lookupEvent(name string) (bool, error) {
	if i.env.eventStore == nil {
		return false, &ErrEventExprNotSupport{}
	}
	return i.env.eventStore.lookup(name) != nil, nil
}

func (i *Interpreter) evalEventExpr(eexpr EventExpr) (bool, error) {
	switch e := eexpr.(type) {
	case *EBinOp:
		l, err := i.evalEventExpr(e.left)
		if err != nil {
			return false, err
		}
		r, err := i.evalEventExpr(e.right)
		if err != nil {
			return false, err
		}

		if e.op == AndEOp {
			return l && r, nil
		} else if e.op == OrEOp {
			return l || r, nil
		} else {
			return false, &ErrNotSupported{Name: string(e.op)}
		}
	case *EVar:
		return i.lookupEvent(e.name)
	default:
		return false, ErrNotSupportedByType(e)
	}
}

// EvalAction evaluates an action
func (i *Interpreter) EvalAction(action *ActionDecl) ([]StmtResult, error) {
	if err := i.resolver.resolveAction(action); err != nil {
		return nil, err
	}
	return i.evalAction(action)
}

func (i *Interpreter) evalAction(action *ActionDecl) ([]StmtResult, error) {
	triggered, err := i.evalEventExpr(action.trigger)
	if err != nil {
		return nil, err
	}
	if !triggered {
		return nil, nil
	}
	results := make([]StmtResult, 0)
	for _, stmt := range action.body {
		r, err := i.evalStmt(stmt)
		if err != nil {
			return nil, err
		}
		results = append(results, r)
	}
	return results, nil
}

func (i *Interpreter) evalFireStmt(fire *Fire) (*FireEventResult, error) {
	result, err := i.evalExpr(fire.eventObj)
	if err != nil {
		return nil, err
	}
	obj, ok := result.(*ObjConst)
	if !ok {
		// Type checker should catch mismatched type error already
		return nil, &ErrTypeMismatch{ExpectedType: fire.eventDecl.ty, ActualType: obj.Type()}
	}
	return NewFireEventResult(fire.eventName, obj), nil
}

func (i *Interpreter) evalStmt(stmt Stmt) (StmtResult, error) {
	switch stmt := stmt.(type) {
	case *Fire:
		return i.evalFireStmt(stmt)
	default:
		return nil, ErrNotSupportedByType(stmt)
	}
}

// EvalExpr evaluates an unbound expression. Unbound expression contains symbols(like function names) unresolved.
// It will resolve and type check before execution. It should not used for a bound expression, it will try to bind symbols again.
func (i *Interpreter) EvalExpr(expr Expr) (Const, error) {
	err := i.resolver.resolveExpr(expr)
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
	case *ObjLit:
		return i.evalObjLit(e)
	case *Prop:
		return i.evalProp(e)
	}
	return nil, &ErrNotSupported{Name: reflect.TypeOf(expr).String()}
}

func (i *Interpreter) evalProp(prop *Prop) (*ObjConst, error) {
	if i.env.eventStore == nil {
		return nil, &ErrEventNotFound{Name: prop.eventVar.name}
	}
	obj := i.env.eventStore.lookup(prop.eventVar.name)
	if obj == nil {
		return nil, &ErrEventNotFound{Name: prop.eventVar.name}
	}
	return obj, nil
}

func (i *Interpreter) evalObjLit(objLit *ObjLit) (*ObjConst, error) {
	result := make(ObjFields, len(objLit.fields))
	for field, expr := range objLit.fields {
		value, err := i.evalExpr(expr)
		if err != nil {
			return nil, err
		}
		result[field] = value
	}
	return NewObjConst(result), nil
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
		return nil, &ErrArgLenMismatch{FuncName: funcCall.Name(), ArgLen: len(args), ParamLen: len(params)}
	}
	return i.evalFuncDecl(decl, args)
}

func (i *Interpreter) evalFuncDecl(decl *FuncDecl, args []Expr) (Const, error) {
	params := decl.Params()
	evaluatedParams := make(map[string]Const, len(params))
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
		evaluatedParams[params[j].Name()] = result
		paramList[j] = result
	}
	return i.cache.getFuncCallResult(decl, paramList, func() (Const, error) {
		result, err := decl.evaluator(i.env, evaluatedParams)
		if err != nil {
			return nil, err
		}
		if !result.Type().Equal(decl.RetType()) {
			return nil, &ErrFunRetTypeMismatch{FuncName: decl.Name(), ExpectedType: decl.RetType(), ActualType: result.Type()}
		}
		return result, nil
	})
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

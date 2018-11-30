// Copyright 2018 @ MegaSpace

// The MegaSpace source code is distributed in the hope that it will be useful,
// but WITHOUT ANY WARRANTY; without even the implied warranty of
// MERCHANTABILITY or FITNESS FOR A PARTICULAR PURPOSE. See the
// GNU Lesser General Public License for more details.

// You should have received a copy of the GNU Lesser General Public License
// along with the MegaSpace source code. If not, see <http://www.gnu.org/licenses/>.

package workflow

type scope map[string]Expr

// TypeChecker infers types and checks soundness
type TypeChecker struct {
	workflow *WorkflowDecl
	events   map[string]*EventDecl

	// stack tracks varable scope.
	// Currently it only supports a one level, so stack size is always 0 or 1
	stack scope
}

// NewTypeChecker creates an instance of type checker
func NewTypeChecker(wf *WorkflowDecl) *TypeChecker {
	events := make(map[string]*EventDecl)
	for _, e := range wf.EventDecls() {
		events[e.name.id] = e
	}
	return &TypeChecker{
		workflow: wf,
		events:   events,
	}
}

func (t *TypeChecker) push(scope scope) {
	if t.stack != nil {
		panic("stack overflow")
	}
	t.stack = scope
}

func (t *TypeChecker) pop() {
	t.stack = nil
}

// Check returns all typing errors
func (t *TypeChecker) Check() Errors {
	var errs Errors
	for _, m := range t.workflow.MonitorDecls() {
		errs = errs.Concat(t.checkMonitor(m))
	}
	for _, a := range t.workflow.ActionDecls() {
		errs = errs.Concat(t.checkAction(a))
	}
	return errs
}

func (t *TypeChecker) checkMonitor(m *MonitorDecl) Errors {
	t.push(m.vars.ExprMap())
	defer t.pop()

	var errs Errors
	for _, v := range m.vars {
		errs = errs.Concat(t.inferExpr(v.Expr))
	}
	return errs.
		Concat(t.inferExpr(m.cond)).
		Concat(expectTy(m.cond, BoolType)).
		Concat(t.checkFire(m.event))
}

func (t *TypeChecker) checkAction(action *ActionDecl) Errors {
	errs := t.checkTrigger(action)
	for _, s := range action.body {
		errs = errs.Concat(t.checkStmt(s))
	}
	return errs
}

func (t *TypeChecker) checkStmt(stmt Stmt) Errors {
	switch s := stmt.(type) {
	case *Fire:
		return t.checkFire(s)
	default:
		return ToErrors(ErrNotSupportedByType(stmt))
	}
}

func (t *TypeChecker) checkTrigger(action *ActionDecl) Errors {
	triggerEvents := action.TriggerEvents()
	var errs Errors
	for _, e := range triggerEvents {
		_, ok := t.events[e]
		if !ok {
			errs = errs.Wrap(SetErrPos(&ErrEventNotFound{Name: e}, action))
		}
	}
	return errs
}

func (t *TypeChecker) checkFire(fire *Fire) Errors {
	event, ok := t.events[fire.eventName]
	if !ok {
		return ToErrors(SetErrPos(&ErrEventNotFound{Name: fire.eventName}, fire))
	}
	return t.inferExpr(fire.eventObj).
		Concat(expectTy(fire.eventObj, event.ty))
}

func (t *TypeChecker) inferExpr(expr Expr) Errors {
	toErrors := func(err ErrWithPos) Errors { return ToErrors(SetErrPos(err, expr)) }
	if expr.Type() != nil {
		return nil
	}
	switch e := expr.(type) {
	case Const:
		return nil
	case *BinOp:
		errs := t.inferExpr(e.left)
		errs.Concat(t.inferExpr(e.right))
		switch e.op {
		case EqualOp, NotEqualOp:
			e.ty = BoolType
			if IsNumeric(e.left.Type()) && IsNumeric(e.right.Type()) {
				return errs
			}
			return errs.Concat(expectTy(e.right, e.left.Type()))
		case LessThanOp, LessThanEqualOp, GreaterThanOp, GreaterThanEqualOp:
			e.ty = BoolType
			return errs.
				Concat(expectTy(e.right, IntType, RatType)).
				Concat(expectTy(e.left, IntType, RatType))
		case PlusOp, MinusOp, MultOp, DivOp:
			leftTy := e.left.Type()
			rightTy := e.right.Type()
			if (leftTy != nil && leftTy.Equal(RatType)) ||
				(rightTy != nil && rightTy.Equal(RatType)) {
				e.ty = RatType
			} else {
				e.ty = IntType
			}
			return errs.
				Concat(expectTy(e.right, IntType, RatType)).
				Concat(expectTy(e.left, IntType, RatType))
		case AndOp, OrOp:
			e.ty = BoolType
			return errs.
				Concat(expectTy(e.left, BoolType)).
				Concat(expectTy(e.right, BoolType))
		default:
			return ToErrors(&ErrNotSupported{Name: e.op.String()})
		}
	case *UniOp:
		errs := t.inferExpr(e.operant)
		switch e.op {
		case NotOp:
			e.ty = BoolType
			return errs.Concat(expectTy(e.operant, BoolType))
		case MinusOp:
			ty, err := matchTy(e.operant, IntType, RatType)
			if ty == nil {
				// ty is nil, set this expr to IntType for further infering
				e.ty = IntType
			}
			e.ty = ty
			return errs.Concat(err)
		default:
			return ToErrors(ErrNotSupportedByType(expr))
		}
	case *ObjAccessor:
		errs := t.inferExpr(e.receiver)
		ty, err := expectObjTy(e.receiver)
		errs = errs.Concat(err)
		if ty == nil {
			return errs
		}
		fieldTy, ok := ty.fields[e.field]
		if !ok {
			return errs.Wrap((SetErrPos(&ErrObjFieldMissing{Field: e.field, ObjType: ty}, e)))
		}
		e.ty = fieldTy.ty
		return errs
	case *FuncCall:
		if e.decl == nil {
			return toErrors(&ErrSymbolNotResolved{Symbol: e.name.id})
		}
		e.ty = e.decl.retType
		if len(e.args) != len(e.decl.params) {
			return toErrors(
				&ErrArgLenMismatch{FuncName: e.name.id, ParamLen: len(e.decl.params), ArgLen: len(e.args)},
			)
		}
		var errs Errors
		for i, arg := range e.args {
			errs = errs.Concat(t.inferExpr(arg))
			if arg.Type() != nil && !arg.Type().Equal(e.decl.params[i].Type()) {
				errs = errs.Wrap(&ErrArgTypeMismatch{
					FuncName:  e.name.id,
					ParamName: e.decl.params[i].name,
					ParamType: e.decl.params[i].ty,
					ArgType:   arg.Type(),
				})
			}
		}
		return errs
	case *Props:
		event, ok := t.events[e.eventVar.name]
		if !ok {
			return toErrors(&ErrEventNotFound{Name: e.eventVar.name})
		}
		e.ty = event.ty
		return nil
	case *ObjLit:
		var errs Errors
		idToTy := IdToTy{}
		for f, expr := range e.fields {
			errs = errs.Concat(t.inferExpr(expr.Expr))
			if expr.Expr.Type() == nil {
				idToTy = nil
			}
			if idToTy != nil {
				idToTy[f] = idTy{id: expr.Id, ty: expr.Expr.Type()}
			}
		}
		if idToTy != nil {
			e.ty = NewObjType(idToTy)
		}
		return errs
	case *Var:
		varexpr, ok := t.stack[e.name]
		if !ok {
			return toErrors(&ErrVarNotFound{VarName: e.name})
		}
		// var declaration should be inferred already, don't need to
		// call inferExpr on varexpr to avoid duplicated errors
		e.ty = varexpr.Type()
		return nil
	default:
		return ToErrors(ErrNotSupportedByType(expr))
	}
}

func expectObjTy(expr Expr) (*ObjType, Errors) {
	if expr.Type() == nil {
		return nil, nil
	}
	objTy, ok := expr.Type().(*ObjType)
	if !ok {
		return nil, ToErrors(SetErrPos(&ErrAccessorRequireObjType{Type: expr.Type()}, expr))
	}
	return objTy, nil
}

func matchTy(expr Expr, tys ...Type) (Type, Errors) {
	if expr.Type() == nil {
		// don't report error if expr type is nil
		// since typing error occurred already, it should stop infering further.
		return nil, nil
	}
	for _, ty := range tys {
		// don't report error if any tys is nil
		// since typing error occurred already, it should stop infering further.
		if ty == nil || expr.Type().Equal(ty) {
			return ty, nil
		}
	}
	return nil, ToErrors(
		SetErrPos(&ErrTypeMismatch{ExpectedTypes: tys, ActualType: expr.Type()}, expr),
	)
}

func expectTy(expr Expr, tys ...Type) Errors {
	_, errs := matchTy(expr, tys...)
	return errs
}

// Copyright 2018 @ MegaSpace

// The MegaSpace source code is distributed in the hope that it will be useful,
// but WITHOUT ANY WARRANTY; without even the implied warranty of
// MERCHANTABILITY or FITNESS FOR A PARTICULAR PURPOSE. See the
// GNU Lesser General Public License for more details.

// You should have received a copy of the GNU Lesser General Public License
// along with the MegaSpace source code. If not, see <http://www.gnu.org/licenses/>.

package workflow

type scope map[string]Expr

type TypeChecker struct {
	ns       []*NamespaceDecl
	workflow *WorkflowDecl

	// stack tracks varable scope.
	// Currently it only supports a one level, so stack size is always 0 or 1
	stack scope
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

func (t *TypeChecker) Check(wf *WorkflowDecl) error {
	t.workflow = wf
	defer func() { t.workflow = nil }()

	return nil
}

func (t *TypeChecker) inferExpr(expr Expr) Errors {
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
			if e.left.Type().Equal(RatType) || e.right.Type().Equal(RatType) {
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
			return ToErrors(setPos(&ErrObjFieldMissing{Field: e.field, ObjType: ty}, e))
		}
		e.ty = fieldTy.ty
		return errs
	case *Var:
		varexpr, ok := t.stack[e.name]
		if !ok {
			return ToErrors(setPos(&ErrVarNotFound{VarName: e.name}, e))
		}
		errs := t.inferExpr(varexpr)
		e.ty = varexpr.Type()
		return errs
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
		return nil, ToErrors(setPos(&ErrAccessorRequireObjType{Type: expr.Type()}, expr))
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
		setPos(&ErrTypeMismatch{ExpectedTypes: tys, ActualType: expr.Type()}, expr),
	)
}

func expectTy(expr Expr, tys ...Type) Errors {
	_, errs := matchTy(expr, tys...)
	return errs
}

func setPos(err WFError, pos interface{ Pos() *Pos }) WFError {
	err.SetPos(pos.Pos())
	return err
}

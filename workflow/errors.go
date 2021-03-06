// Copyright 2018 @ MegaSpace

// The MegaSpace source code is distributed in the hope that it will be useful,
// but WITHOUT ANY WARRANTY; without even the implied warranty of
// MERCHANTABILITY or FITNESS FOR A PARTICULAR PURPOSE. See the
// GNU Lesser General Public License for more details.

// You should have received a copy of the GNU Lesser General Public License
// along with the MegaSpace source code. If not, see <http://www.gnu.org/licenses/>.

package workflow

import (
	"bytes"
	"fmt"
	"reflect"

	p "github.com/megaspacelab/megaconnect/prettyprint"
)

type Errors func([]error) []error

// ToErrors convert an error to Errors
func ToErrors(err error) Errors {
	if err == nil {
		return nil
	}
	return func(errors []error) []error {
		return append(errors, err)
	}
}

// First returns first element in Errors, or nil if empty
func (e Errors) First() error {
	if e == nil {
		return nil
	}
	errs := e.ToErr()
	if len(errs) == 0 {
		return nil
	}
	return errs[0]
}

// Empty returns true if Errors is empty otherwise false
func (e Errors) Empty() bool {
	return e == nil || len(e([]error{})) == 0
}

// Wrap appends an error to Errors
func (e Errors) Wrap(err error) Errors {
	if err == nil {
		return e
	}
	return func(errors []error) []error {
		if e == nil {
			return []error{err}
		}
		return append(e(errors), err)
	}
}

// Concat joins two Errors
func (e Errors) Concat(errors Errors) Errors {
	if e == nil {
		return errors
	} else if errors == nil {
		return e
	}
	return func(errs []error) []error {
		return e(errors(errs))
	}
}

// ToErr returns array of error represented by Errors
func (e Errors) ToErr() []error {
	if e == nil {
		return []error{}
	}
	return e([]error{})
}

type ErrWithPos interface {
	error
	HasPos
}

func SetErrPos(err ErrWithPos, pos interface{ Pos() *Pos }) error {
	err.SetPos(pos.Pos())
	return err
}

// ErrArgTypeMismatch is returned if argument type does not match parameter type
type ErrArgTypeMismatch struct {
	hasPos
	FuncName  string
	ParamName string
	ParamType Type
	ArgType   Type
}

func (e *ErrArgTypeMismatch) Error() string {
	return fmt.Sprintf("In %s , %s expects %s not %s", e.FuncName, e.ParamName, e.ParamType, e.ArgType)
}

// ErrTypeMismatch returns if two types does not match
type ErrTypeMismatch struct {
	hasPos
	ExpectedTypes []Type
	ActualType    Type
}

func (e *ErrTypeMismatch) Error() string {
	types := make([]p.PrinterOp, len(e.ExpectedTypes))
	for i, p := range e.ExpectedTypes {
		types[i] = p.Print()
	}
	return fmt.Sprintf("Expected %s, but got %s ", p.String(separatedBy(types, p.Text(", "))), PrintNode(e.ActualType))
}

// ErrNotSupported is returned if an operator is not supported
type ErrNotSupported struct{ Name string }

func (e *ErrNotSupported) Error() string {
	return fmt.Sprintf("%s is not supported", e.Name)
}

// ErrNotSupportedByType creates an ErrNotSupported and reports type of x is not supported
func ErrNotSupportedByType(x interface{}) *ErrNotSupported {
	return &ErrNotSupported{Name: reflect.TypeOf(x).String()}
}

// ErrFunRetTypeMismatch is returned if return type does not match signature
type ErrFunRetTypeMismatch struct {
	FuncName     string
	ExpectedType Type
	ActualType   Type
}

func (e *ErrFunRetTypeMismatch) Error() string {
	return fmt.Sprintf("Return type of %s expects %s not %s", e.FuncName, e.ExpectedType, e.ActualType)
}

// ErrObjFieldIncompatible is returned if there is incompatible field in two objects
type ErrObjFieldIncompatible struct {
	Fields  []string
	ObjType *ObjType
}

func (e *ErrObjFieldIncompatible) Error() string {
	return fmt.Sprintf("Fields %s incompatible with obj %s", e.Fields, PrintNode(e.ObjType))
}

// ErrObjFieldMissing is returned if expected field is missing in an object
type ErrObjFieldMissing struct {
	hasPos
	Field   string
	ObjType *ObjType
}

func (e *ErrObjFieldMissing) Error() string {
	return fmt.Sprintf("Field %s missing in obj %s", e.Field, PrintNode(e.ObjType))
}

// ErrObjFieldTypeMismatch is returned if a type does not match field type in an object
type ErrObjFieldTypeMismatch struct {
	Field   string
	ObjType *ObjType
	Type    Type
}

func (e *ErrObjFieldTypeMismatch) Error() string {
	return fmt.Sprintf("Type %s mismatch type of field %s in obj %s", e.Type, e.Field, PrintNode(e.ObjType))
}

// ErrSymbolNotResolved is returned if a symbol is not resolved
type ErrSymbolNotResolved struct {
	hasPos
	Symbol string
}

func (e *ErrSymbolNotResolved) Error() string {
	return fmt.Sprintf("Symbol %s not resolved", e.Symbol)
}

// ErrAccessorRequireObjType is returned if the expression of an object accessor does not evaluate to an object
type ErrAccessorRequireObjType struct {
	hasPos
	Type Type
}

func (e *ErrAccessorRequireObjType) Error() string {
	return fmt.Sprintf("Object accessor requires object type, not %s", e.Type)
}

// ErrConstExpected is returned if expected const not present
type ErrConstExpected struct {
	Actual string
}

func (e *ErrConstExpected) Error() string {
	return fmt.Sprintf("Expect const type not %s", e.Actual)
}

// ErrMissingArg is returned if expected argument is missing
type ErrMissingArg struct {
	ArgName string
	Func    string
}

func (e *ErrMissingArg) Error() string {
	return fmt.Sprintf("Expected argument %s is missing, function: %s", e.ArgName, e.Func)
}

// ErrVarNotFound is returned if a variable is not declared
type ErrVarNotFound struct {
	hasPos
	VarName string
}

func (e *ErrVarNotFound) Error() string {
	return fmt.Sprintf("Variable %s is not declared", e.VarName)
}

// ErrVarDeclaredAlready is returned if same variable has been declared in the scope
type ErrVarDeclaredAlready struct {
	VarName string
}

func (e *ErrVarDeclaredAlready) Error() string {
	return fmt.Sprintf("Variable %s is declared already", e.VarName)
}

// ErrEventExprNotSupport is returned if interpreter can not evaluate event expression
type ErrEventExprNotSupport struct{}

func (e *ErrEventExprNotSupport) Error() string {
	return fmt.Sprintf("Event expression evaluation not supported")
}

// ErrArgLenMismatch is returned if supplied arguments len does not match expected parameters len
type ErrArgLenMismatch struct {
	hasPos
	FuncName string
	ArgLen   int
	ParamLen int
}

func (e *ErrArgLenMismatch) Error() string {
	return fmt.Sprintf("Number of arguments mismatch, %d arguments expected from function %s, %d given", e.ParamLen, e.FuncName, e.ArgLen)
}

// ErrEventNotFound is returned if event is not found
type ErrEventNotFound struct {
	hasPos
	Name string
}

func (e *ErrEventNotFound) Error() string {
	return fmt.Sprintf("Event %s not found", e.Name)
}

// ErrCycleVars is returned if there is a cycle dependency of variable declarations
type ErrCycleVars struct {
	hasPos
	Path []string
}

func (e *ErrCycleVars) Error() string {
	var buf bytes.Buffer
	buf.WriteString(e.Path[0])
	for _, p := range e.Path[1:] {
		buf.WriteString(" -> ")
		buf.WriteString(p)
	}
	return fmt.Sprintf("variables should not have cycle dependency, %s", buf.String())
}

// ErrDupNames is returned if a name is defined already
type ErrDupNames struct {
	hasPos
	Name string
}

func (e *ErrDupNames) Error() string {
	return fmt.Sprintf("`%s` is already defined", e.Name)
}

// ErrCycleEvents is returned if there is a cycle dependency of actions
type ErrCycleEvents struct {
	hasPos
	Path []string
}

func (e *ErrCycleEvents) Error() string {
	var buf bytes.Buffer
	buf.WriteString(e.Path[0])
	for _, p := range e.Path[1:] {
		buf.WriteString(" -> ")
		buf.WriteString(p)
	}
	return fmt.Sprintf("events should not have cycle dependency, %s", buf.String())
}

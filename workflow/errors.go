// Copyright 2018 @ MegaSpace

// The MegaSpace source code is distributed in the hope that it will be useful,
// but WITHOUT ANY WARRANTY; without even the implied warranty of
// MERCHANTABILITY or FITNESS FOR A PARTICULAR PURPOSE. See the
// GNU Lesser General Public License for more details.

// You should have received a copy of the GNU Lesser General Public License
// along with the MegaSpace source code. If not, see <http://www.gnu.org/licenses/>.

package workflow

import (
	"fmt"
	"reflect"
)

// TODO: Add an error base struct for storing (row, col) after parsing supported

// ErrArgTypeMismatch is returned if argument type does not match parameter type
type ErrArgTypeMismatch struct {
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
	ExpectedTypes []Type
	ActualType    Type
}

func (e *ErrTypeMismatch) Error() string {
	return fmt.Sprintf("Expected %#q, but got %s ", e.ExpectedTypes, e.ActualType)
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
	Symbol string
}

func (e *ErrSymbolNotResolved) Error() string {
	return fmt.Sprintf("Symbol %s not resolved", e.Symbol)
}

// ErrAccessorRequireObjType is returned if the expression of an object accessor does not evaluate to an object
type ErrAccessorRequireObjType struct {
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
	FuncName string
	ArgLen   int
	ParamLen int
}

func (e *ErrArgLenMismatch) Error() string {
	return fmt.Sprintf("Number of arguments mismatch, %d arguments expected from function %s, only %d given", e.ParamLen, e.FuncName, e.ArgLen)
}

// ErrEventNotFound is returned if event is not found
type ErrEventNotFound struct {
	Name string
}

func (e *ErrEventNotFound) Error() string {
	return fmt.Sprintf("Event %s not found", e.Name)
}

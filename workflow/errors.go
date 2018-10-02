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
	ExpectedType Type
	ActualType   Type
}

func (e *ErrTypeMismatch) Error() string {
	return fmt.Sprintf("Expected %s, but got %s ", e.ExpectedType, e.ActualType)
}

// ErrNotSupported is returned if an operator is not supported
type ErrNotSupported struct{ Name string }

func (e *ErrNotSupported) Error() string {
	return fmt.Sprintf("%s is not supported", e.Name)
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
	return fmt.Sprintf("Fields %s incompatible with obj %s", e.Fields, e.ObjType)
}

// ErrObjFieldMissing is returned if expected field is missing in an object
type ErrObjFieldMissing struct {
	Field   string
	ObjType *ObjType
}

func (e *ErrObjFieldMissing) Error() string {
	return fmt.Sprintf("Field %s missing in obj %s", e.Field, e.ObjType)
}

// ErrObjFieldTypeMismatch is returned if a type does not match field type in an object
type ErrObjFieldTypeMismatch struct {
	Field   string
	ObjType *ObjType
	Type    Type
}

func (e *ErrObjFieldTypeMismatch) Error() string {
	return fmt.Sprintf("Type %s mismatch type of field %s in obj %s", e.Type, e.Field, e.ObjType)
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

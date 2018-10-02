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
)

// Type represents all types in the language, all different type derives from this interface
type Type interface {
	Equal(ty Type) bool
	String() string
	Methods() FuncDecls
}

// ObjFieldTypes defines a mapping from field name to its type
type ObjFieldTypes map[string]Type

// Copy returns a new ObjFieldTypes
func (o ObjFieldTypes) Copy() ObjFieldTypes {
	result := make(map[string]Type)
	for f, ty := range o {
		result[f] = ty
	}
	return result
}

// ObjType is a type of non-primitive type
type ObjType struct {
	fields ObjFieldTypes
}

// NewObjType creates a new ObjType
func NewObjType(fields ObjFieldTypes) *ObjType { return &ObjType{fields: fields.Copy()} }

// FieldsCount return number of fields
func (o *ObjType) FieldsCount() int { return len(o.fields) }

// Fields returns a new copy of field name to type mapping
func (o *ObjType) Fields() ObjFieldTypes { return o.fields.Copy() }

// Methods returns nil, ObjType current does not support method
func (o *ObjType) Methods() FuncDecls { return nil }

func (o ObjType) String() string {
	var buf bytes.Buffer
	buf.WriteString("{")
	i := len(o.fields)
	for f, ty := range o.fields {
		buf.WriteString(f)
		buf.WriteString(": ")
		buf.WriteString(ty.String())
		if i != 1 {
			buf.WriteString(",")
		}
		i--
	}
	buf.WriteString("}")
	return buf.String()
}

// Equal compares whether `ty` is equal current type
func (o ObjType) Equal(ty Type) bool {
	oty, ok := ty.(*ObjType)
	if !ok {
		return false
	}
	if len(o.fields) != len(oty.fields) {
		return false
	}

	for field, ty := range o.fields {
		if fieldTy, ok := oty.fields[field]; !ok || !fieldTy.Equal(ty) {
			return false
		}
	}
	return true
}

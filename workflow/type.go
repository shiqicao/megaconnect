// Copyright 2018 @ MegaSpace

// The MegaSpace source code is distributed in the hope that it will be useful,
// but WITHOUT ANY WARRANTY; without even the implied warranty of
// MERCHANTABILITY or FITNESS FOR A PARTICULAR PURPOSE. See the
// GNU Lesser General Public License for more details.

// You should have received a copy of the GNU Lesser General Public License
// along with the MegaSpace source code. If not, see <http://www.gnu.org/licenses/>.

package workflow

import (
	p "github.com/megaspacelab/megaconnect/prettyprint"
)

// Type represents all types in the language, all different type derives from this interface
type Type interface {
	Node
	Equal(ty Type) bool
	Methods() FuncDecls
}

// ObjType is a type of non-primitive type
type ObjType struct {
	node
	fields IdToTy
}

// NewObjType creates a new ObjType
func NewObjType(fields IdToTy) *ObjType { return &ObjType{fields: fields.Copy()} }

// FieldsCount return number of fields
func (o *ObjType) FieldsCount() int { return len(o.fields) }

// Fields returns a new copy of field name to type mapping
func (o *ObjType) Fields() IdToTy { return o.fields.Copy() }

// Methods returns nil, ObjType current does not support method
func (o *ObjType) Methods() FuncDecls { return nil }

func (o *ObjType) Print() p.PrinterOp {
	body := []p.PrinterOp{}
	for _, ty := range o.fields {
		body = append(body, p.Concat(ty.id.Print(), p.Text(": "), ty.ty.Print()))
	}
	return p.Concat(
		p.Text(" {"),
		p.Nest(1, separatedBy(body, p.Concat(p.Text(","), p.Line()))),
		p.Line(),
		p.Text("}"),
	)
}

// Equal compares whether `ty` is equal current type
func (o *ObjType) Equal(ty Type) bool {
	if o == nil || ty == nil {
		return false
	}
	oty, ok := ty.(*ObjType)
	if !ok {
		return false
	}
	if len(o.fields) != len(oty.fields) {
		return false
	}

	for field, ty := range o.fields {
		if fieldTy, ok := oty.fields[field]; !ok || !fieldTy.ty.Equal(ty.ty) {
			return false
		}
	}
	return true
}

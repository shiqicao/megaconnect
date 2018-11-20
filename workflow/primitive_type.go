// Copyright 2018 @ MegaSpace

// The MegaSpace source code is distributed in the hope that it will be useful,
// but WITHOUT ANY WARRANTY; without even the implied warranty of
// MERCHANTABILITY or FITNESS FOR A PARTICULAR PURPOSE. See the
// GNU Lesser General Public License for more details.

// You should have received a copy of the GNU Lesser General Public License
// along with the MegaSpace source code. If not, see <http://www.gnu.org/licenses/>.

package workflow

import pp "github.com/megaspacelab/megaconnect/prettyprint"

var (
	// BoolType represents boolean type
	BoolType = &PrimitiveType{ty: booleanTy}

	// StrType represents string type
	StrType = &PrimitiveType{ty: stringTy}

	// IntType represents a big int type
	IntType = &PrimitiveType{ty: intTy}

	// RatType represents a big rational type
	RatType = &PrimitiveType{ty: ratTy}
)

// PrimitiveType represents all primitive types in the language, including int, string, bool, etc.
type PrimitiveType struct {
	node
	ty   primitiveTy
	mths FuncDecls
}

// Methods returns a list of methods applicable to this type.
// Currently it only contains build-in functions for the primitive types
func (p *PrimitiveType) Methods() FuncDecls { return p.mths.Copy() }

func (p *PrimitiveType) Print() pp.PrinterOp {
	switch p.ty {
	case stringTy:
		return pp.Text("string")
	case booleanTy:
		return pp.Text("bool")
	case intTy:
		return pp.Text("int")
	case ratTy:
		return pp.Text("rat")
	default:
		return pp.Nil()
	}
}

type primitiveTy = uint8

const (
	stringTy primitiveTy = iota
	booleanTy
	intTy
	ratTy
)

// Equal compares whether `ty` is the same primitive as the current one
func (p *PrimitiveType) Equal(ty Type) bool {
	if p == nil || ty == nil {
		return false
	}
	pty, ok := ty.(*PrimitiveType)
	if !ok {
		return false
	}
	return pty.ty == p.ty
}

func IsNumeric(ty Type) bool { return ty.Equal(IntType) || ty.Equal(RatType) }

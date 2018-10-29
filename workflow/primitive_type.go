// Copyright 2018 @ MegaSpace

// The MegaSpace source code is distributed in the hope that it will be useful,
// but WITHOUT ANY WARRANTY; without even the implied warranty of
// MERCHANTABILITY or FITNESS FOR A PARTICULAR PURPOSE. See the
// GNU Lesser General Public License for more details.

// You should have received a copy of the GNU Lesser General Public License
// along with the MegaSpace source code. If not, see <http://www.gnu.org/licenses/>.

package workflow

import "fmt"

var (
	// BoolType represents boolean type
	BoolType = &PrimitiveType{ty: booleanTy}

	// StrType represents string type
	StrType = &PrimitiveType{ty: stringTy}

	// IntType represents a big int type
	IntType = createIntType()
)

// PrimitiveType represents all primitive types in the language, including int, string, bool, etc.
type PrimitiveType struct {
	ty   primitiveTy
	mths FuncDecls
}

// Methods returns a list of methods applicable to this type.
// Currently it only contains build-in functions for the primitive types
func (p *PrimitiveType) Methods() FuncDecls { return p.mths.Copy() }

func (p *PrimitiveType) String() string {
	return fmt.Sprintf("PrimitiveType[%d]", p.ty)
}

type primitiveTy = uint8

const (
	stringTy primitiveTy = iota
	booleanTy
	intTy
)

// Equal compares whether `ty` is the same primitive as the current one
func (p *PrimitiveType) Equal(ty Type) bool {
	pty, ok := ty.(*PrimitiveType)
	if !ok {
		return false
	}
	return pty.ty == p.ty
}

func createIntType() *PrimitiveType {
	intType := &PrimitiveType{
		ty: intTy,
	}
	lessThan := NewFuncDecl(
		opToFunc[LessThanOp],
		Params{NewParamDecl("x", intType), NewParamDecl("y", intType)},
		BoolType,
		func(_ *Env, args map[string]Const) (Const, error) {
			xraw := args["x"]
			yraw := args["y"]
			x, ok := xraw.(*IntConst)
			// This should be checked already
			if !ok {
				return nil, &ErrTypeMismatch{ExpectedType: intType, ActualType: xraw.Type()}
			}
			y, ok := yraw.(*IntConst)
			if !ok {
				return nil, &ErrTypeMismatch{ExpectedType: intType, ActualType: yraw.Type()}
			}

			return GetBoolConst(x.Value().Cmp(y.Value()) < 0), nil
		},
	)
	intType.mths = FuncDecls{lessThan}
	return intType
}

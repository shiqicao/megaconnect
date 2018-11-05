// Copyright 2018 @ MegaSpace

// The MegaSpace source code is distributed in the hope that it will be useful,
// but WITHOUT ANY WARRANTY; without even the implied warranty of
// MERCHANTABILITY or FITNESS FOR A PARTICULAR PURPOSE. See the
// GNU Lesser General Public License for more details.

// You should have received a copy of the GNU Lesser General Public License
// along with the MegaSpace source code. If not, see <http://www.gnu.org/licenses/>.

package workflow

// Id represents identifer node in AST
type Id struct {
	node
	id string
}

type idTy struct {
	id *Id
	ty Type
}

// IdToTy is a list of pair of identifer and type,
// and indexed by string representation of identifer
type IdToTy map[string]idTy

// NewIdToTy creates an instance of IdToTy
func NewIdToTy() IdToTy { return make(IdToTy) }

// Add inserts a new pair of Id and type
func (i IdToTy) Add(id string, ty Type) IdToTy {
	i[id] = idTy{id: &Id{id: id}, ty: ty}
	return i
}

// Copy creates a new instance of IdToTy
func (i IdToTy) Copy() IdToTy {
	if i == nil {
		return i
	}
	new := make(IdToTy, len(i))
	for k, ty := range i {
		new[k] = ty
	}
	return new
}

type idExpr struct {
	id   *Id
	expr Expr
}

// IdToExpr is a list of pair of identifier and expr,
// and indexed by string representation of identifier.
type IdToExpr map[string]idExpr

// NewIdToExpr creates an instance of IdToExpr
func NewIdToExpr() IdToExpr { return make(IdToExpr) }

// Equal returns true if two IdToExpr are equivalent
func (i IdToExpr) Equal(x IdToExpr) bool {
	if len(i) != len(x) {
		return false
	}
	for n, e := range i {
		if xn, ok := x[n]; !ok || !xn.expr.Equal(e.expr) {
			return false
		}
	}
	return true
}

// Add inserts a pair of identifier and expression
func (i IdToExpr) Add(id string, expr Expr) IdToExpr {
	i[id] = idExpr{id: &Id{id: id}, expr: expr}
	return i
}

// Copy creates a new instance of IdToExpr
func (i IdToExpr) Copy() IdToExpr {
	if i == nil {
		return nil
	}
	new := make(IdToExpr, len(i))
	for k, expr := range i {
		new[k] = expr
	}
	return new
}

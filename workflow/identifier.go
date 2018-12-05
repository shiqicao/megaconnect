// Copyright 2018 @ MegaSpace

// The MegaSpace source code is distributed in the hope that it will be useful,
// but WITHOUT ANY WARRANTY; without even the implied warranty of
// MERCHANTABILITY or FITNESS FOR A PARTICULAR PURPOSE. See the
// GNU Lesser General Public License for more details.

// You should have received a copy of the GNU Lesser General Public License
// along with the MegaSpace source code. If not, see <http://www.gnu.org/licenses/>.

package workflow

import p "github.com/megaspacelab/megaconnect/prettyprint"

// Id represents identifer node in AST
type Id struct {
	node
	id string
}

func NewId(id string) *Id { return &Id{id: id} }

func NewIdB(id []byte) *Id { return NewId(string(id)) }

func (i *Id) Print() p.PrinterOp { return p.Text(i.id) }

func (i *Id) String() string { return i.id }

func (i *Id) Children() []Node { return nil }

// Equal returns true if two identifiers are equivalent
func (i *Id) Equal(x *Id) bool {
	if i != nil && x != nil {
		return i.id == x.id
	}
	return i == x
}

// Id returns string representation
func (i *Id) Id() string { return i.id }

type idTy struct {
	id *Id
	ty Type
}

// IdToTy is a list of pair of identifer and type,
// and indexed by string representation of identifer
type IdToTy map[string]idTy

// NewIdToTy creates an instance of IdToTy
func NewIdToTy() IdToTy { return make(IdToTy) }

// Nodes returns a list of nodes
func (i IdToTy) Nodes() (nodes []Node) {
	for _, idTy := range i {
		nodes = append(nodes, idTy.id, idTy.ty)
	}
	return
}

// Put inserts a new pair of Id and type
func (i IdToTy) Put(id string, ty Type) IdToTy {
	i[id] = idTy{id: &Id{id: id}, ty: ty}
	return i
}

// Add inserts a pair of identifer and type,
// it returns false if id is not unique.
// It returns true if add a new pair is succeeded
func (i IdToTy) Add(id *Id, ty Type) bool {
	if _, ok := i[id.id]; ok {
		return false
	}
	i[id.id] = idTy{id: id, ty: ty}
	return true
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

type IdExpr struct {
	Id   *Id
	Expr Expr
}

// IdToExpr is a list of pair of identifier and expr,
// and indexed by string representation of identifier.
type IdToExpr map[string]IdExpr

// NewIdToExpr creates an instance of IdToExpr
func NewIdToExpr() IdToExpr { return make(IdToExpr) }

// Nodes returns a list of nodes
func (i IdToExpr) Nodes() (nodes []Node) {
	for _, idExpr := range i {
		nodes = append(nodes, idExpr.Id, idExpr.Expr)
	}
	return
}

// ExprMap returns a map from id string to Expr
func (i IdToExpr) ExprMap() map[string]Expr {
	result := make(map[string]Expr, len(i))
	for k, expr := range i {
		result[k] = expr.Expr
	}
	return result
}

// Equal returns true if two IdToExpr are equivalent
func (i IdToExpr) Equal(x IdToExpr) bool {
	if len(i) != len(x) {
		return false
	}
	for n, e := range i {
		if xn, ok := x[n]; !ok || !xn.Expr.Equal(e.Expr) {
			return false
		}
	}
	return true
}

// Put inserts a pair of identifier and expression
func (i IdToExpr) Put(id string, expr Expr) IdToExpr {
	i[id] = IdExpr{Id: &Id{id: id}, Expr: expr}
	return i
}

// Add inserts a pair of identifer and expression,
// it returns false if id is not unique.
// It returns true if add a new pair is succeeded
func (i IdToExpr) Add(id *Id, expr Expr) bool {
	_, ok := i[id.id]
	if ok {
		return false
	}
	i[id.id] = IdExpr{Id: id, Expr: expr}
	return true
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

// Print pretty prints code
func (i IdToExpr) Print(multiline bool, assignmentSymbol p.PrinterOp) p.PrinterOp {
	vars := []p.PrinterOp{}
	for _, decl := range i {
		vars = append(vars, p.Concat(
			decl.Id.Print(),
			assignmentSymbol,
			decl.Expr.Print(),
		))
	}
	var separator p.PrinterOp
	if multiline {
		separator = p.Concat(p.Text(","), p.Line())
	} else {
		separator = p.Text(", ")
	}
	return separatedBy(vars, separator)
}

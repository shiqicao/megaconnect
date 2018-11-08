// Copyright 2018 @ MegaSpace

// The MegaSpace source code is distributed in the hope that it will be useful,
// but WITHOUT ANY WARRANTY; without even the implied warranty of
// MERCHANTABILITY or FITNESS FOR A PARTICULAR PURPOSE. See the
// GNU Lesser General Public License for more details.

// You should have received a copy of the GNU Lesser General Public License
// along with the MegaSpace source code. If not, see <http://www.gnu.org/licenses/>.

package workflow

// Resolver resolves symbol in workflow lang
type Resolver struct {
	libs             []*NamespaceDecl
	defaultNamespace NamespacePrefix
}

// NewResolver creates a new instance of Resolver
func NewResolver(libs []*NamespaceDecl, defaultNamespace NamespacePrefix) *Resolver {
	return &Resolver{
		libs:             libs,
		defaultNamespace: defaultNamespace,
	}
}

// resolveAction resolves all symbols in action declaration
func (r *Resolver) resolveAction(action *ActionDecl) error {
	for _, s := range action.body {
		if err := r.resolveStmt(s); err != nil {
			return err
		}
	}
	return nil
}

func (r *Resolver) resolveStmt(stmt Stmt) error {
	switch s := stmt.(type) {
	case *Fire:
		return r.resolveExpr(s.eventObj)
	default:
		return ErrNotSupportedByType(stmt)
	}
}

func (r *Resolver) resolveExpr(expr Expr) error {
	switch e := expr.(type) {
	case *BinOp:
		err := r.resolveExpr(e.Left())
		if err != nil {
			return err
		}
		return r.resolveExpr(e.Right())
	case *UniOp:
		return r.resolveExpr(e.Operant())
	case *ObjAccessor:
		return r.resolveExpr(e.Receiver())
	case *FuncCall:
		for _, arg := range e.Args() {
			if err := r.resolveExpr(arg); err != nil {
				return err
			}
		}
		return r.resolveFuncCall(e)
	case *ObjLit:
		for _, expr := range e.fields {
			if err := r.resolveExpr(expr.Expr); err != nil {
				return err
			}
		}
	}
	return nil
}

func (r *Resolver) resolveFuncCall(fun *FuncCall) error {
	ns := fun.NamespacePrefix()
	if ns.IsEmpty() && !r.defaultNamespace.IsEmpty() {
		ns = r.defaultNamespace
	}
	decl := resolveFun(r.libs, fun.Name().id, ns)
	if decl == nil {
		return &ErrSymbolNotResolved{Symbol: fun.Name().id}
	}
	fun.SetDecl(decl)
	return nil
}

func resolveFun(nss []*NamespaceDecl, name string, prefix NamespacePrefix) *FuncDecl {
	if len(prefix) == 0 {
		return nil
	}

	cur := prefix[0]
	var ns *NamespaceDecl
	for _, n := range nss {
		if cur.id == n.name {
			ns = n
			break
		}
	}
	if ns == nil {
		return nil
	}

	if len(prefix) > 1 {
		return resolveFun(ns.children, name, prefix[1:])
	}

	for _, fun := range ns.funs {
		if fun.Name() == name {
			return fun
		}
	}
	return nil
}

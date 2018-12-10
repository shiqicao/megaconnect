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
func NewResolver(libs []*NamespaceDecl) *Resolver {
	return &Resolver{
		libs: libs,
	}
}

// ResolveWorkflow returns errors if any symbol can't be resolved.
// It also sets declaration on symbols.
func (r *Resolver) ResolveWorkflow(wf *WorkflowDecl) Errors {
	var errs Errors
	for _, m := range wf.MonitorDecls() {
		errs = errs.Concat(r.ResolveMonitor(m))
	}
	for _, a := range wf.ActionDecls() {
		errs = errs.Concat(r.resolveAction(a))
	}
	return errs
}

// ResolveMonitor resolves all symbols in action declaration
func (r *Resolver) ResolveMonitor(monitor *MonitorDecl) Errors {
	r.defaultNamespace = NamespacePrefix{NewId(monitor.chain)}
	defer func() { r.defaultNamespace = nil }()

	errs := r.resolveExpr(monitor.cond).
		Concat(r.resolveStmt(monitor.event))
	for _, v := range monitor.vars {
		errs = errs.Concat(r.resolveExpr(v.Expr))
	}
	return errs
}

// resolveAction resolves all symbols in action declaration
func (r *Resolver) resolveAction(action *ActionDecl) Errors {
	var errs Errors
	for _, s := range action.body {
		errs = errs.Concat(r.resolveStmt(s))
	}
	return errs
}

func (r *Resolver) resolveStmt(stmt Stmt) Errors {
	switch s := stmt.(type) {
	case *Fire:
		return r.resolveExpr(s.eventObj)
	default:
		return ToErrors(ErrNotSupportedByType(stmt))
	}
}

func (r *Resolver) resolveExpr(expr Expr) Errors {
	switch e := expr.(type) {
	case *BinOp:
		return r.resolveExpr(e.Left()).
			Concat(r.resolveExpr(e.Right()))
	case *UniOp:
		return r.resolveExpr(e.Operant())
	case *ObjAccessor:
		return r.resolveExpr(e.Receiver())
	case *FuncCall:
		var errs Errors
		for _, arg := range e.Args() {
			errs = errs.Concat(r.resolveExpr(arg))
		}
		return errs.Wrap(r.resolveFuncCall(e))
	case *ObjLit:
		var errs Errors
		for _, expr := range e.fields {
			errs = errs.Concat(r.resolveExpr(expr.Expr))
		}
		return errs
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
		return SetErrPos(&ErrSymbolNotResolved{Symbol: fun.Name().id}, fun.Name())
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
		return resolveFun(ns.namespaces, name, prefix[1:])
	}

	for _, fun := range ns.funs {
		if fun.Name() == name {
			return fun
		}
	}
	return nil
}

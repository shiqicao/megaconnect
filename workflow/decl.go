// Copyright 2018 @ MegaSpace

// The MegaSpace source code is distributed in the hope that it will be useful,
// but WITHOUT ANY WARRANTY; without even the implied warranty of
// MERCHANTABILITY or FITNESS FOR A PARTICULAR PURPOSE. See the
// GNU Lesser General Public License for more details.

// You should have received a copy of the GNU Lesser General Public License
// along with the MegaSpace source code. If not, see <http://www.gnu.org/licenses/>.

package workflow

import (
	"github.com/megaspacelab/megaconnect/common"
	p "github.com/megaspacelab/megaconnect/prettyprint"
)

type evaluator func(*Env, map[string]Const) (Const, error)

// Params is parameter declaration list for a function decalaraion
type Params []*ParamDecl

// Copy returns a new Params
func (p Params) Copy() Params {
	return append(p[:0:0], p...)
}

// Equal returns true if two Params are equivalent
func (p Params) Equal(x Params) bool {
	if len(p) != len(x) {
		return false
	}
	for i, param := range p {
		if param.name != x[i].name || !param.ty.Equal(x[i].ty) {
			return false
		}
	}
	return true
}

// FuncDecl is function declaration which represents a function definition.
// With `evaluator`, FuncDecl is only for build-in functions.
type FuncDecl struct {
	parent  *NamespaceDecl
	name    string
	params  Params
	retType Type
	// eval allows predefined functions execute natively,
	// Later it will support function definition in workflow lang
	evaluator evaluator
}

// NewFuncDecl creates a new FuncDecl
func NewFuncDecl(name string, params Params, retType Type, evaluator evaluator) *FuncDecl {
	return &FuncDecl{
		name:      name,
		params:    params.Copy(),
		retType:   retType,
		evaluator: evaluator,
	}
}

// Name returns function name
func (f *FuncDecl) Name() string { return f.name }

// Params returns a copy of parameter declaration list
func (f *FuncDecl) Params() Params { return f.params.Copy() }

// RetType returns Type of return value
func (f *FuncDecl) RetType() Type { return f.retType }

// Parent returns containing namespace or nil if current is top level
func (f *FuncDecl) Parent() *NamespaceDecl { return f.parent }

// Equal returns true if two func *signatures* are equivalent
func (f *FuncDecl) Equal(x *FuncDecl) bool {
	return f.name == x.name &&
		equalStrings(f.parent.FullName(), x.parent.FullName()) &&
		f.params.Equal(x.params) &&
		f.retType.Equal(x.retType)
}

// Print prints function signature
func (f *FuncDecl) Print() p.PrinterOp {
	params := []p.PrinterOp{}
	for _, param := range f.params {
		params = append(params, p.Concat(p.Text(param.name), p.Text(" "), param.ty.Print()))
	}
	return p.Concat(
		p.Text(f.name),
		paren(separatedBy(params, p.Text(","))),
		p.Text(" "),
		f.retType.Print(),
	)
}

func equalStrings(xs []string, ys []string) bool {
	if len(xs) != len(ys) {
		return false
	}
	for i, x := range xs {
		if x != ys[i] {
			return false
		}
	}
	return true
}

// ParamDecl stores parameter name and its type
type ParamDecl struct {
	name string
	ty   Type
}

// NewParamDecl creates a new ParamDecl
func NewParamDecl(name string, ty Type) *ParamDecl {
	return &ParamDecl{
		name: name,
		ty:   ty,
	}
}

// Name returns parameter name
func (p *ParamDecl) Name() string { return p.name }

// Type returns parameter type
func (p *ParamDecl) Type() Type { return p.ty }

// NamespaceDecl captures namespace declaraion, it includes function declaraion and child namespace
type NamespaceDecl struct {
	parent     *NamespaceDecl
	name       string
	namespaces []*NamespaceDecl
	funs       FuncDecls
}

// NewNamespaceDecl creates a new instance of NamespaceDecl
func NewNamespaceDecl(name string) *NamespaceDecl {
	return &NamespaceDecl{
		name: name,
		funs: FuncDecls{},
	}
}

// Parent returns parent namespace
func (n *NamespaceDecl) Parent() *NamespaceDecl { return n.parent }

// Name returns identifier of this namespace
func (n *NamespaceDecl) Name() string { return n.name }

// FullName returns the unique name
func (n *NamespaceDecl) FullName() []string {
	ret := []string{}
	cur := n
	for ; cur != nil; cur = cur.parent {
		ret = append(ret, cur.name)
	}
	return ret
}

// Equal returns true if two namespaces are equivalent
func (n *NamespaceDecl) Equal(x *NamespaceDecl) bool {
	if len(n.namespaces) != len(x.namespaces) {
		return false
	}
	for _, ns := range x.namespaces {
		if nsn := n.findNamespace(ns.name); nsn == nil || !nsn.Equal(ns) {
			return false
		}
	}
	return equalStrings(n.FullName(), x.FullName()) && n.funs.Equal(x.funs)
}

func (n *NamespaceDecl) findNamespace(name string) *NamespaceDecl {
	for _, ns := range n.namespaces {
		if ns.name == name {
			return ns
		}
	}
	return nil
}

// AddFunc insert a function declaration to this namespace
func (n *NamespaceDecl) AddFunc(funcDecl *FuncDecl) *NamespaceDecl {
	funcDecl.parent = n
	n.funs = append(n.funs, funcDecl)
	return n
}

// AddChild insert a namespace as child namespace
func (n *NamespaceDecl) AddNamespace(ns *NamespaceDecl) *NamespaceDecl {
	ns.parent = n
	n.namespaces = append(n.namespaces, ns)
	return n
}

func (n *NamespaceDecl) Print() p.PrinterOp {
	children := []p.PrinterOp{}
	for _, f := range n.funs {
		children = append(children, f.Print())
	}
	for _, c := range n.namespaces {
		children = append(children, c.Print())
	}
	return declPrintCore("namespace", p.Text(n.name), separatedBy(children, p.Line()))
}

// FuncDecls is a list of function declarations
type FuncDecls []*FuncDecl

func (f FuncDecls) find(name string) *FuncDecl {
	for _, x := range f {
		if x.Name() == name {
			return x
		}
	}
	return nil
}

// Copy creates a new list of function declarations
func (f FuncDecls) Copy() FuncDecls {
	return append(f[:0:0], f...)
}

// Equal returns true if two FuncDecls are equivalent
func (f FuncDecls) Equal(x FuncDecls) bool {
	if len(f) != len(x) {
		return false
	}
	for _, fun := range f {
		xfun := x.find(fun.name)
		if xfun == nil || !xfun.Equal(fun) {
			return false
		}
	}
	return true
}

// MonitorDecl represents a monitor unit in workflow lang
type MonitorDecl struct {
	decl
	name  *Id
	cond  Expr
	vars  IdToExpr
	event *Fire
	chain string
}

// NewMonitorDecl creates a new MonitorDecl
func NewMonitorDecl(name *Id, cond Expr, vars IdToExpr, event *Fire, chain string) *MonitorDecl {
	return &MonitorDecl{
		name:  name,
		cond:  cond,
		vars:  vars.Copy(),
		event: event,
		chain: chain,
	}
}

// Name returns monitor name
func (m *MonitorDecl) Name() *Id { return m.name }

// Chain returns targeted chain
func (m *MonitorDecl) Chain() string { return m.chain }

// Condition returns an boolean expression when the monitor is triggered
func (m *MonitorDecl) Condition() Expr { return m.cond }

// Vars returns a copy of variable declared in this monitor
func (m *MonitorDecl) Vars() IdToExpr { return m.vars.Copy() }

// EventName returns the name this monitor will fire
func (m *MonitorDecl) EventName() string { return m.event.eventName }

// Equal returns true if two monitor declaraions are the same
func (m *MonitorDecl) Equal(x Decl) bool {
	y, ok := x.(*MonitorDecl)
	return ok &&
		m.Name().Equal(x.Name()) &&
		m.Condition().Equal(y.Condition()) &&
		m.vars.Equal(y.vars) &&
		m.chain == y.chain &&
		m.event.Equal(y.event)
}

// Print pretty prints code
func (m *MonitorDecl) Print() p.PrinterOp {
	return declPrint("monitor", m.name, p.Concat(
		p.Text("chain "), p.Text(m.chain), p.Line(),
		p.Text("condition "), m.cond.Print(), p.Line(),
		p.Text("var {"),
		p.Nest(1, m.vars.Print(true, p.Text(" = "))), p.Line(),
		p.Text("}"), p.Line(),
		m.event.Print(),
	))
}

// Decl is an interface for all declarations in a workflow
type Decl interface {
	Node
	Name() *Id
	Parent() *WorkflowDecl
	setParent(*WorkflowDecl)
	Equal(Decl) bool
}

type decl struct {
	node
	parent *WorkflowDecl
}

func (d *decl) setParent(w *WorkflowDecl) { d.parent = w }

// Parent returns the containing workflow declaration
func (d *decl) Parent() *WorkflowDecl { return d.parent }

// EventDecl represents an event declaration
type EventDecl struct {
	decl
	name *Id
	ty   *ObjType
}

// NewEventDecl creates an instance of EventDecl
func NewEventDecl(name *Id, ty *ObjType) *EventDecl {
	return &EventDecl{
		name: name,
		ty:   ty,
	}
}

// Equal returns true if x is the same event declaration
func (e *EventDecl) Equal(x Decl) bool {
	y, ok := x.(*EventDecl)
	return ok && y.name.Equal(e.name) && y.ty.Equal(e.ty)
}

// Name returns event name
func (e *EventDecl) Name() *Id { return e.name }

// Print pretty prints code
func (e *EventDecl) Print() p.PrinterOp {
	return p.Concat(
		p.Text("event "),
		e.name.Print(),
		e.ty.Print(),
	)
}

// WorkflowDecl represents a workflow declaration
type WorkflowDecl struct {
	node
	version uint32
	name    *Id
	decls   []Decl
}

// NewWorkflowDecl creates a new instance of WorkflowDecl
func NewWorkflowDecl(name *Id, version uint32) *WorkflowDecl {
	return &WorkflowDecl{
		version: version,
		name:    name,
		decls:   make([]Decl, 0),
	}
}

// Equal returns true if `x` is the same workflow declaration
func (w *WorkflowDecl) Equal(x *WorkflowDecl) bool {
	if len(w.decls) != len(x.decls) {
		return false
	}
	// declaration order of action declaration implies execution order
	for i := len(w.decls) - 1; i >= 0; i-- {
		if !w.decls[i].Equal(x.decls[i]) {
			return false
		}
	}
	return w.version == x.version && w.name.Equal(x.name)
}

// Version returns workflow lang version
func (w *WorkflowDecl) Version() uint32 { return w.version }

// Name returns workflow name
func (w *WorkflowDecl) Name() *Id { return w.name }

// EventDecls returns all event declarations
func (w *WorkflowDecl) EventDecls() (events []*EventDecl) {
	for _, d := range w.decls {
		if e, ok := d.(*EventDecl); ok {
			events = append(events, e)
		}
	}
	return
}

// ActionDecls returns all event declarations
func (w *WorkflowDecl) ActionDecls() (actions []*ActionDecl) {
	for _, d := range w.decls {
		if e, ok := d.(*ActionDecl); ok {
			actions = append(actions, e)
		}
	}
	return
}

// MonitorDecls returns all monitor declarations in order
func (w *WorkflowDecl) MonitorDecls() (monitors []*MonitorDecl) {
	for _, d := range w.decls {
		if e, ok := d.(*MonitorDecl); ok {
			monitors = append(monitors, e)
		}
	}
	return
}

// AddDecl adds a child declaraion
func (w *WorkflowDecl) AddDecl(decl Decl) *WorkflowDecl {
	w.decls = append(w.decls, decl)
	decl.setParent(w)
	return w
}

// AddDecls adds a list of child declarations
func (w *WorkflowDecl) AddDecls(decls []Decl) *WorkflowDecl {
	for _, decl := range decls {
		w.AddDecl(decl)
	}
	return w
}

// Print pretty prints code
func (w *WorkflowDecl) Print() p.PrinterOp {
	decls := []p.PrinterOp{}
	for _, c := range w.decls {
		decls = append(decls, c.Print())
	}
	return declPrint("workflow", w.name, separatedBy(decls, p.Line()))
}

// ActionDecl represents an action declaration
type ActionDecl struct {
	decl
	name    *Id
	trigger EventExpr
	body    Stmts
}

// NewActionDecl creates an instance of ActionDecl
func NewActionDecl(name *Id, trigger EventExpr, body Stmts) *ActionDecl {
	return &ActionDecl{
		name:    name,
		trigger: trigger,
		body:    body.Copy(),
	}
}

// Equal returns true if x is the same action declaration
func (a *ActionDecl) Equal(x Decl) bool {
	y, ok := x.(*ActionDecl)
	return ok && a.name.Equal(y.name) && a.trigger.Equal(y.trigger) && a.body.Equal(y.body)
}

// Name returns action name
func (a *ActionDecl) Name() *Id { return a.name }

// Trigger returns action trigger
func (a *ActionDecl) Trigger() EventExpr { return a.trigger }

// Body returns action run statement
func (a *ActionDecl) Body() Stmts { return a.body.Copy() }

// TriggerEvents returns a set of all events used in trigger
func (a *ActionDecl) TriggerEvents() []string {
	events := make(map[string]common.Nothing)
	visitor := EventExprVisitor{
		VisitVar: func(v *EVar) interface{} {
			if _, ok := events[v.name]; !ok {
				events[v.name] = struct{}{}
			}
			return nil
		},
	}
	visitor.Visit(a.trigger)
	result := make([]string, 0, len(events))
	for e := range events {
		result = append(result, e)
	}
	return result
}

// Print pretty prints code
func (a *ActionDecl) Print() p.PrinterOp {
	body := []p.PrinterOp{
		p.Text("trigger "),
		a.trigger.Print(),
		p.Line(),
		p.Text("run {"),
		p.Nest(1, a.body.Print()),
		p.Line(),
		p.Text("}"),
	}
	return declPrint("action", a.name, p.Concat(body...))
}

func declPrint(keyword string, id *Id, body p.PrinterOp) p.PrinterOp {
	return declPrintCore(keyword, id.Print(), body)
}

func declPrintCore(keyword string, id p.PrinterOp, body p.PrinterOp) p.PrinterOp {
	return p.Concat(
		p.Text(keyword),
		p.Text(" "),
		id,
		p.Text(" {"),
		p.Nest(1, body),
		p.Line(),
		p.Text("}"),
	)
}

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

	"github.com/megaspacelab/megaconnect/common"
)

type evaluator func(*Env, map[string]Const) (Const, error)

// Params is parameter declaration list for a function decalaraion
type Params []*ParamDecl

// Copy returns a new Params
func (p Params) Copy() Params {
	if p == nil {
		return nil
	}
	r := make(Params, len(p))
	copy(r, p)
	return r
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
	parent   *NamespaceDecl
	name     string
	children []*NamespaceDecl
	funs     FuncDecls
}

// FuncDecls is a list of function declarations
type FuncDecls []*FuncDecl

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

// AddFunc insert a function declaration to this namespace
func (n *NamespaceDecl) AddFunc(funcDecl *FuncDecl) *NamespaceDecl {
	funcDecl.parent = n
	n.funs = append(n.funs, funcDecl)
	return n
}

// AddChild insert a namespace as child namespace
func (n *NamespaceDecl) AddChild(child *NamespaceDecl) *NamespaceDecl {
	child.parent = n
	n.children = append(n.children, child)
	return n
}

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
	if f == nil {
		return nil
	}
	r := make(FuncDecls, len(f))
	copy(r, f)
	return r
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

func (m *MonitorDecl) String() string {
	var buf bytes.Buffer
	buf.WriteString("monitor {")
	buf.WriteString("name = ")
	buf.WriteString(m.Name().id)
	buf.WriteString(",")
	buf.WriteString("condition = ")
	buf.WriteString(m.Condition().String())
	buf.WriteString("}")
	return buf.String()
}

// Decl is an interface for all declarations in a workflow
type Decl interface {
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

// WorkflowDecl represents a workflow declaration
type WorkflowDecl struct {
	node
	version  uint32
	name     *Id
	children []Decl
}

// NewWorkflowDecl creates a new instance of WorkflowDecl
func NewWorkflowDecl(name *Id, version uint32) *WorkflowDecl {
	return &WorkflowDecl{
		version:  version,
		name:     name,
		children: make([]Decl, 0),
	}
}

// Equal returns true if `x` is the same workflow declaration
func (w *WorkflowDecl) Equal(x *WorkflowDecl) bool {
	if len(w.children) != len(x.children) {
		return false
	}
	// declaration order of action declaration implies execution order
	for i := len(w.children) - 1; i >= 0; i-- {
		if !w.children[i].Equal(x.children[i]) {
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
	for _, d := range w.children {
		if e, ok := d.(*EventDecl); ok {
			events = append(events, e)
		}
	}
	return
}

// ActionDecls returns all event declarations
func (w *WorkflowDecl) ActionDecls() (actions []*ActionDecl) {
	for _, d := range w.children {
		if e, ok := d.(*ActionDecl); ok {
			actions = append(actions, e)
		}
	}
	return
}

// MonitorDecls returns all monitor declarations in order
func (w *WorkflowDecl) MonitorDecls() (monitors []*MonitorDecl) {
	for _, d := range w.children {
		if e, ok := d.(*MonitorDecl); ok {
			monitors = append(monitors, e)
		}
	}
	return
}

// AddChild adds a child declaraion
func (w *WorkflowDecl) AddChild(child Decl) *WorkflowDecl {
	w.children = append(w.children, child)
	child.setParent(w)
	return w
}

func (w *WorkflowDecl) String() string {
	var buf bytes.Buffer
	buf.WriteString("workflow ")
	buf.WriteString(w.name.id)
	buf.WriteString(" {")
	/*
		for _, _ := range w.children {
			// TODO: implement String() for child
			// buf.WriteString(child.String())
		}
	*/
	buf.WriteString("}")
	return buf.String()
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

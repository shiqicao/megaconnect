// Copyright 2018 @ MegaSpace

// The MegaSpace source code is distributed in the hope that it will be useful,
// but WITHOUT ANY WARRANTY; without even the implied warranty of
// MERCHANTABILITY or FITNESS FOR A PARTICULAR PURPOSE. See the
// GNU Lesser General Public License for more details.

// You should have received a copy of the GNU Lesser General Public License
// along with the MegaSpace source code. If not, see <http://www.gnu.org/licenses/>.

package workflow

import "bytes"

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

func (n *NamespaceDecl) addFunc(funcDecl *FuncDecl) *NamespaceDecl {
	funcDecl.parent = n
	n.funs = append(n.funs, funcDecl)
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

// VarDecls is a set of variable declarations
type VarDecls map[string]Expr

// Equal returns true iff two VarDecls are the same
func (v VarDecls) Equal(x VarDecls) bool {
	if len(v) != len(x) {
		return false
	}
	for n, e := range v {
		if xe, ok := x[n]; !ok || !xe.Equal(e) {
			return false
		}
	}
	return true
}

// Copy return a new instance of VarDecls with identical content
func (v VarDecls) Copy() VarDecls {
	if v == nil {
		return nil
	}
	r := make(VarDecls, len(v))
	for n, e := range v {
		r[n] = e
	}
	return r
}

// MonitorDecl represents a monitor unit in workflow lang
type MonitorDecl struct {
	name string
	cond Expr
	vars VarDecls
}

// NewMonitorDecl creates a new MonitorDecl
func NewMonitorDecl(name string, cond Expr, vars VarDecls) *MonitorDecl {
	return &MonitorDecl{
		name: name,
		cond: cond,
		vars: vars.Copy(),
	}
}

// Name returns monitor name
func (m *MonitorDecl) Name() string { return m.name }

// Condition returns an boolean expression when the monitor is triggered
func (m *MonitorDecl) Condition() Expr { return m.cond }

// Vars returns a copy of variable declared in this monitor
func (m *MonitorDecl) Vars() VarDecls { return m.vars.Copy() }

// Equal returns true if two monitor declaraions are the same
func (m *MonitorDecl) Equal(x *MonitorDecl) bool {
	return m.Name() == x.Name() && m.Condition().Equal(x.Condition()) && m.vars.Equal(x.vars)
}

func (m *MonitorDecl) String() string {
	var buf bytes.Buffer
	buf.WriteString("monitor {")
	buf.WriteString("name = ")
	buf.WriteString(m.Name())
	buf.WriteString(",")
	buf.WriteString("condition = ")
	buf.WriteString(m.Condition().String())
	buf.WriteString("}")
	return buf.String()
}

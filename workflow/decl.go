// Copyright 2018 @ MegaSpace

// The MegaSpace source code is distributed in the hope that it will be useful,
// but WITHOUT ANY WARRANTY; without even the implied warranty of
// MERCHANTABILITY or FITNESS FOR A PARTICULAR PURPOSE. See the
// GNU Lesser General Public License for more details.

// You should have received a copy of the GNU Lesser General Public License
// along with the MegaSpace source code. If not, see <http://www.gnu.org/licenses/>.

package workflow

type evaluator func(*Env, map[string]Const) (Const, error)

// Params is parameter declaration list for a function decalaraion
type Params []*ParamDecl

// Copy returns a new Params
func (p Params) Copy() Params {
	r := make(Params, len(p))
	for i, param := range p {
		r[i] = param
	}
	return r
}

// FuncDecl is function declaration which represents a function definition.
// With `evaluator`, FuncDecl is only for build-in functions.
type FuncDecl struct {
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
	name     string
	children []*NamespaceDecl
	funs     FuncDecls
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
	r := make(FuncDecls, len(f))
	for i, f := range f {
		r[i] = f
	}
	return r
}

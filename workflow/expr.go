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
	"math/big"
	"strconv"
)

var (
	// TrueConst is an instance of BoolConst
	TrueConst = &BoolConst{value: true}

	// FalseConst is an instance of BoolConst
	FalseConst = &BoolConst{value: false}
)

// Expr represents expression in the language. All type of expression derives from it.
type Expr interface {
	String() string
	Equal(Expr) bool
}

// Args is a list of function arguments, it is referenced in FuncCall expression
type Args []Expr

// Copy creates a new Args
func (a Args) Copy() Args {
	r := make([]Expr, len(a))
	copy(r, a)
	return r
}

// Equal returns true if two argument lists are the same
func (a Args) Equal(b Args) bool {
	if len(a) != len(b) {
		return false
	}
	for i := 0; i < len(a); i++ {
		if !a[i].Equal(b[i]) {
			return false
		}
	}
	return true
}

// FuncCall represents a function invoking expression
type FuncCall struct {
	decl *FuncDecl
	name string
	args Args
	ns   NamespacePrefix
}

// NewFuncCall returns a new FuncCall
func NewFuncCall(ns NamespacePrefix, name string, args ...Expr) *FuncCall {
	return &FuncCall{
		name: name,
		args: Args(args).Copy(),
		ns:   ns.Copy(),
	}
}

// NamespacePrefix returns namespace of a function
func (f *FuncCall) NamespacePrefix() NamespacePrefix { return f.ns.Copy() }

// Name returns the function name
func (f *FuncCall) Name() string { return f.name }

// Args returns a copy of function arguments
func (f *FuncCall) Args() Args { return f.args.Copy() }

// Decl returns a function declaration of the function,
// it returns nil if the function is not defined or analyzer not yet resolve the name
func (f *FuncCall) Decl() *FuncDecl { return f.decl }

// SetDecl binds the function name to a function declaration
func (f *FuncCall) SetDecl(decl *FuncDecl) { f.decl = decl }

func (f *FuncCall) String() string {
	var buf bytes.Buffer
	buf.WriteString(f.ns.String())
	buf.WriteString(f.name)
	buf.WriteString("(")
	last := len(f.args) - 1
	for i, p := range f.args {
		buf.WriteString(p.String())
		if i != last {
			buf.WriteString(",")
		}
	}
	buf.WriteString(")")
	return buf.String()
}

// Equal returns true if two expressions are the same
func (f *FuncCall) Equal(expr Expr) bool {
	y, ok := expr.(*FuncCall)
	return ok &&
		f.Name() == y.Name() &&
		f.Args().Equal(y.Args()) &&
		f.NamespacePrefix().Equal(y.NamespacePrefix())
}

// Const represents a single value or an object, it derives from Expr
type Const interface {
	Expr
	Type() Type
}

// BoolConst is a value typed to BoolType
type BoolConst struct{ value bool }

// GetBoolConst converts value to TrueConst or FalseConst
func GetBoolConst(value bool) *BoolConst {
	if value {
		return TrueConst
	}
	return FalseConst
}

// Type returns the type of this constant
func (b *BoolConst) Type() Type { return BoolType }

// Value returns corresponding value in hosting language(Go)
func (b *BoolConst) Value() bool { return b.value }

func (b *BoolConst) String() string { return strconv.FormatBool(b.value) }

// Negate is helper to get negation of this bool value
func (b *BoolConst) Negate() *BoolConst { return GetBoolConst(!b.value) }

// Equal returns whether x is equal to the current value
func (b *BoolConst) Equal(x Expr) bool {
	y, ok := x.(*BoolConst)
	return ok && y.value == b.value
}

// StrConst is a value typed to StrType
type StrConst struct{ value string }

// NewStrConst lifts a string from hosting language(Go)
func NewStrConst(value string) *StrConst { return &StrConst{value: value} }

// Type returns the type of current value
func (s *StrConst) Type() Type { return StrType }

// Value returns corresponding value in hosting language
func (s *StrConst) Value() string { return s.value }

func (s *StrConst) String() string { return s.value }

// Equal returns whether x is equal to the current value
func (s *StrConst) Equal(x Expr) bool {
	y, ok := x.(*StrConst)
	return ok && s.value == y.value
}

// IntConst represents a big interger
type IntConst struct{ value *big.Int }

// NewIntConst lifts a big integer from hosting language(Go)
func NewIntConst(value *big.Int) *IntConst { return &IntConst{value: value} }

// NewIntConstFromI64 lifts an int64 from hosting language(Go)
func NewIntConstFromI64(value int64) *IntConst { return NewIntConst(big.NewInt(value)) }

// Type returns IntType
func (i *IntConst) Type() Type { return IntType }

// Value returns corresponding value in hosting language
func (i *IntConst) Value() *big.Int { return new(big.Int).Set(i.value) }

func (i *IntConst) String() string { return i.value.String() }

// Equal return whether x is equal to the current value
func (i *IntConst) Equal(x Expr) bool {
	y, ok := x.(*IntConst)
	return ok && i.value.Cmp(y.value) == 0
}

// ObjConst represents an object, an object contains a list of fields and corresponding types
type ObjConst struct {
	ty    *ObjType
	value ObjFields
}

// ObjFields represents a mapping from field name to value
type ObjFields map[string]Const

// Copy creates a new mapping
func (o ObjFields) Copy() ObjFields {
	result := make(ObjFields, len(o))
	for n, v := range o {
		result[n] = v
	}
	return result
}

// NewObjConst converts a list of named `Const` to ObjConst, it also calculates type of this ObjConst
func NewObjConst(values ObjFields) *ObjConst {
	ty := make(map[string]Type)
	for field, value := range values {
		ty[field] = value.Type()
	}
	return &ObjConst{
		ty:    NewObjType(ty),
		value: values.Copy(),
	}
}

// NewObjConstWithTy creates an ObjConst with a list of named `Const` and expected ObjConst type.
// It returns error if type mismatch
func NewObjConstWithTy(ty *ObjType, values ObjFields) (*ObjConst, error) {
	// check type
	if ty.FieldsCount() != len(values) {
		fields := make([]string, len(values))
		for k := range values {
			fields = append(fields, k)
		}
		return nil, &ErrObjFieldIncompatible{Fields: fields, ObjType: ty}
	}
	for field, fieldTy := range ty.Fields() {
		value, ok := values[field]
		if !ok {
			return nil, &ErrObjFieldMissing{Field: field, ObjType: ty}
		}
		if _, ok := fieldTy.(*ObjType); ok && value == nil {
			continue
		}
		if !fieldTy.Equal(value.Type()) {
			return nil, &ErrObjFieldTypeMismatch{Field: field, ObjType: ty, Type: value.Type()}
		}
	}

	return &ObjConst{
		ty:    ty,
		value: values.Copy(),
	}, nil
}

// Type returns the object type
func (o *ObjConst) Type() Type { return o.ty }

// Value returns a copy of fields
func (o *ObjConst) Value() ObjFields { return o.value.Copy() }

func (o *ObjConst) String() string {
	var buf bytes.Buffer
	buf.WriteString("{")
	i := len(o.value)
	for f, v := range o.value {
		buf.WriteString(f)
		buf.WriteString(": ")
		buf.WriteString(v.String())
		if i != 1 {
			buf.WriteString(",")
		}
		i--
	}
	buf.WriteString("}")
	return buf.String()
}

// Equal checks whether x is equal current object.
func (o *ObjConst) Equal(x Expr) bool {
	y, ok := x.(*ObjConst)
	if !ok || len(o.value) != len(y.value) {
		return false
	}
	for field, value := range o.value {
		yValue, ok := y.value[field]
		if !ok || !value.Equal(yValue) {
			return false
		}
	}
	return true
}

// NamespacePrefix represents a namespace hierarchy
type NamespacePrefix []string

func (n NamespacePrefix) String() string {
	var buf bytes.Buffer
	for _, ns := range n {
		buf.WriteString(ns)
		buf.WriteString("::")
	}
	return buf.String()
}

// Equal returns true if two namespace prefix are the same
func (n NamespacePrefix) Equal(m NamespacePrefix) bool {
	if len(n) != len(m) {
		return false
	}
	for i := 0; i < len(n); i++ {
		if n[i] != m[i] {
			return false
		}
	}
	return true
}

// Copy creates a new namespace prefix
func (n NamespacePrefix) Copy() NamespacePrefix {
	if n == nil {
		return nil
	}
	r := make(NamespacePrefix, len(n))
	for i, m := range n {
		r[i] = m
	}
	return r
}

// ObjAccessor represents field selection operation,
// for example, A.foo, where A is an object and foo is a field of A
type ObjAccessor struct {
	receiver Expr
	field    string
}

// NewObjAccessor creates a new ObjAccessor
func NewObjAccessor(receiver Expr, field string) *ObjAccessor {
	return &ObjAccessor{
		receiver: receiver,
		field:    field,
	}
}

// Receiver returns the expression which is expected to be evaluated to an object
func (o *ObjAccessor) Receiver() Expr { return o.receiver }

func (o *ObjAccessor) String() string { return o.receiver.String() + "." + o.field }

// Field returns the accessor
func (o *ObjAccessor) Field() string { return o.field }

// Equal returns true if two expressions are the same
func (o *ObjAccessor) Equal(expr Expr) bool {
	y, ok := expr.(*ObjAccessor)
	return ok &&
		o.field == y.field &&
		o.receiver.Equal(y.receiver)
}

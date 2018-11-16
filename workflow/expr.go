// Copyright 2018 @ MegaSpace

// The MegaSpace source code is distributed in the hope that it will be useful,
// but WITHOUT ANY WARRANTY; without even the implied warranty of
// MERCHANTABILITY or FITNESS FOR A PARTICULAR PURPOSE. See the
// GNU Lesser General Public License for more details.

// You should have received a copy of the GNU Lesser General Public License
// along with the MegaSpace source code. If not, see <http://www.gnu.org/licenses/>.

package workflow

import (
	"math/big"
	"strconv"

	p "github.com/megaspacelab/megaconnect/prettyprint"
)

var (
	// TrueConst is an instance of BoolConst
	TrueConst = &BoolConst{value: true}

	// FalseConst is an instance of BoolConst
	FalseConst = &BoolConst{value: false}
)

// Expr represents expression in the language. All type of expression derives from it.
type Expr interface {
	Node
	Equal(Expr) bool
}

type expr struct {
	node
}

// Args is a list of function arguments, it is referenced in FuncCall expression
type Args []Expr

// Copy creates a new Args
func (a Args) Copy() Args {
	if a == nil {
		return nil
	}
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

// Print pretty prints code
func (a Args) Print() p.PrinterOp {
	ops := []p.PrinterOp{}
	for _, arg := range a {
		ops = append(ops, arg.Print())
	}
	return separatedBy(ops, p.Text(", "))
}

// FuncCall represents a function invoking expression
type FuncCall struct {
	expr
	decl *FuncDecl
	name *Id
	args Args
	ns   NamespacePrefix
}

// NewFuncCall returns a new FuncCall
func NewFuncCall(ns NamespacePrefix, name *Id, args ...Expr) *FuncCall {
	return &FuncCall{
		name: name,
		args: Args(args).Copy(),
		ns:   ns.Copy(),
	}
}

// NamespacePrefix returns namespace of a function
func (f *FuncCall) NamespacePrefix() NamespacePrefix { return f.ns.Copy() }

// Name returns the function name
func (f *FuncCall) Name() *Id { return f.name }

// Args returns a copy of function arguments
func (f *FuncCall) Args() Args { return f.args.Copy() }

// Decl returns a function declaration of the function,
// it returns nil if the function is not defined or analyzer not yet resolve the name
func (f *FuncCall) Decl() *FuncDecl { return f.decl }

// SetDecl binds the function name to a function declaration
func (f *FuncCall) SetDecl(decl *FuncDecl) { f.decl = decl }

// Print pretty prints code
func (f *FuncCall) Print() p.PrinterOp {
	return p.Concat(
		f.ns.Print(),
		f.name.Print(),
		paren(f.args.Print()),
	)
}

// Equal returns true if two expressions are the same
func (f *FuncCall) Equal(expr Expr) bool {
	y, ok := expr.(*FuncCall)
	return ok &&
		f.Name().Equal(y.Name()) &&
		f.Args().Equal(y.Args()) &&
		f.NamespacePrefix().Equal(y.NamespacePrefix())
}

// Const represents a single value or an object, it derives from Expr
type Const interface {
	Expr
	Type() Type
}

// BoolConst is a value typed to BoolType
type BoolConst struct {
	expr
	value bool
}

// GetBoolConst converts value to TrueConst or FalseConst
func GetBoolConst(value bool) *BoolConst {
	if value {
		return TrueConst
	}
	return FalseConst
}

// GetBoolConstFromStr converts value to TrueConst or FalseConst
func GetBoolConstFromStr(value string) *BoolConst {
	if value == "true" {
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

// Print pretty prints code
func (b *BoolConst) Print() p.PrinterOp {
	if b.value {
		return p.Text("true")
	}
	return p.Text("false")
}

// StrConst is a value typed to StrType
type StrConst struct {
	expr
	value string
}

// NewStrConst lifts a string from hosting language(Go)
func NewStrConst(value string) *StrConst { return &StrConst{value: value} }

// Type returns the type of current value
func (s *StrConst) Type() Type { return StrType }

// Value returns corresponding value in hosting language
func (s *StrConst) Value() string { return s.value }

// Equal returns whether x is equal to the current value
func (s *StrConst) Equal(x Expr) bool {
	y, ok := x.(*StrConst)
	return ok && s.value == y.value
}

// Print pretty prints code
func (s *StrConst) Print() p.PrinterOp { return p.Text("\"" + s.value + "\"") }

// IntConst represents a big interger
type IntConst struct {
	expr
	value *big.Int
}

// NewIntConst lifts a big integer from hosting language(Go)
func NewIntConst(value *big.Int) *IntConst { return &IntConst{value: value} }

// NewIntConstFromI64 lifts an int64 from hosting language(Go).
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

// Print pretty prints code
func (i *IntConst) Print() p.PrinterOp { return p.Text(i.value.String()) }

// RatConst represents a big rational number
type RatConst struct {
	expr
	value *big.Rat
}

// NewRatConstFromF64 creates a new instance of RatConst from float64
func NewRatConstFromF64(value float64) *RatConst {
	return &RatConst{value: new(big.Rat).SetFloat64(value)}
}

// NewRatConst lifts a big float from hosting language(Go)
func NewRatConst(value *big.Rat) *RatConst { return &RatConst{value: value} }

// Type returns FloatType
func (r *RatConst) Type() Type { return RatType }

// Value returns corresponding value in hosting language
func (r *RatConst) Value() *big.Rat { return new(big.Rat).Set(r.value) }

// Equal returns true if x is equivalent to f
func (r *RatConst) Equal(x Expr) bool {
	y, ok := x.(*RatConst)
	return ok && r.value.Cmp(y.value) == 0
}

// Print pretty prints code
func (r *RatConst) Print() p.PrinterOp {
	// TODO: determine best precision
	return p.Text(r.value.FloatString(20))
}

// ObjConst represents an object, an object contains a list of fields and corresponding types
type ObjConst struct {
	expr
	ty    *ObjType
	value ObjFields
}

// ObjFields represents a mapping from field name to value
type ObjFields map[string]Const

// Fields returns all fields
func (o ObjFields) Fields() []string {
	fields := make([]string, 0, len(o))
	for f := range o {
		fields = append(fields, f)
	}
	return fields
}

// Copy creates a new mapping
func (o ObjFields) Copy() ObjFields {
	if o == nil {
		return nil
	}
	result := make(ObjFields, len(o))
	for n, v := range o {
		result[n] = v
	}
	return result
}

// NewObjConst converts a list of named `Const` to ObjConst, it also calculates type of this ObjConst
func NewObjConst(values ObjFields) *ObjConst {
	ty := make(IdToTy, len(values))
	for field, value := range values {
		ty.Put(field, value.Type())
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
			return nil, &ErrObjFieldMissing{Field: fieldTy.id.id, ObjType: ty}
		}
		if _, ok := fieldTy.ty.(*ObjType); ok && value == nil {
			continue
		}
		if !fieldTy.ty.Equal(value.Type()) {
			return nil, &ErrObjFieldTypeMismatch{Field: fieldTy.id.id, ObjType: ty, Type: value.Type()}
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

// Print pretty prints code
func (o *ObjConst) Print() p.PrinterOp {
	ops := []p.PrinterOp{}
	for field, expr := range o.value {
		ops = append(ops, p.Concat(p.Text(field), p.Text(": "), expr.Print()))
	}
	return p.Concat(
		p.Text("{"),
		separatedBy(ops, p.Text(", ")),
		p.Text("}"),
	)
}

// Fields returns a new object sorted by field name
func (o *ObjConst) Fields() []string {
	return o.Value().Fields()
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
type NamespacePrefix []*Id

// IsEmpty returns true if namespace is empty
func (n NamespacePrefix) IsEmpty() bool { return len(n) == 0 }

// Print pretty prints code
func (n NamespacePrefix) Print() p.PrinterOp {
	ops := []p.PrinterOp{}
	for _, id := range n {
		ops = append(ops, id.Print(), p.Text("::"))
	}
	return p.Concat(ops...)
}

// Equal returns true if two namespace prefix are the same
func (n NamespacePrefix) Equal(m NamespacePrefix) bool {
	if len(n) != len(m) {
		return false
	}
	for i := 0; i < len(n); i++ {
		if !n[i].Equal(m[i]) {
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
	copy(r, n)
	return r
}

// ObjAccessor represents field selection operation,
// for example, A.foo, where A is an object and foo is a field of A
type ObjAccessor struct {
	expr
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

// func (o *ObjAccessor) String() string { return o.receiver.String() + "." + o.field }

// Field returns the accessor
func (o *ObjAccessor) Field() string { return o.field }

// Equal returns true if two expressions are the same
func (o *ObjAccessor) Equal(expr Expr) bool {
	y, ok := expr.(*ObjAccessor)
	return ok &&
		o.field == y.field &&
		o.receiver.Equal(y.receiver)
}

// Print pretty prints code
func (o *ObjAccessor) Print() p.PrinterOp {
	parenthesis := false
	if _, isBin := o.receiver.(*BinOp); isBin {
		parenthesis = true
	} else if _, isUni := o.receiver.(*UniOp); isUni {
		parenthesis = true
	}

	receiver := o.receiver.Print()
	if parenthesis {
		receiver = paren(receiver)
	}
	return p.Concat(
		receiver,
		p.Text("."),
		p.Text(o.field),
	)
}

// Var represents a variable in workflow lang
type Var struct {
	expr
	name string
}

// NewVar creates Var instance
func NewVar(name string) *Var {
	return &Var{name: name}
}

// Equal returns true if `expr` is the same variable
func (v *Var) Equal(expr Expr) bool {
	y, ok := expr.(*Var)
	return ok && y.name == v.name
}

// Name returns variable name
func (v *Var) Name() string { return v.name }

// Print pretty prints code
func (v *Var) Print() p.PrinterOp { return p.Text(v.name) }

// ObjLit represents an object literal,
// for example, {a: 1 + 1, b: "bar"}
type ObjLit struct {
	expr
	fields IdToExpr
}

// NewObjLit creates an instance of ObjLit
func NewObjLit(fields IdToExpr) *ObjLit {
	return &ObjLit{
		fields: fields.Copy(),
	}
}

// Equal returns true if two ObjLit are equivalent
func (o *ObjLit) Equal(expr Expr) bool {
	y, ok := expr.(*ObjLit)
	return ok && y.fields.Equal(o.fields)
}

// Fields returns a copy of field declaration
func (o *ObjLit) Fields() IdToExpr { return o.fields.Copy() }

// Print pretty prints code
func (o *ObjLit) Print() p.PrinterOp {
	return p.Concat(
		p.Text("{"),
		o.fields.Print(false, p.Text(": ")),
		p.Text("}"),
	)
}

// Props is a unary operator of event type, it returns properties of an event as an obj.
// Props is defined as a struct instead of an operator in UniOp. Props should alway apply on an variable,
// but UniOp can not enforce this criteria.
type Props struct {
	expr
	eventVar *Var
}

// NewProps returns a new instance of Prop
func NewProps(v *Var) *Props {
	return &Props{
		eventVar: v,
	}
}

// Print pretty prints code
func (pr *Props) Print() p.PrinterOp {
	return p.Concat(
		p.Text("props"),
		paren(pr.eventVar.Print()),
	)
}

// Equal returns true if two Prop are equivalent
func (pr *Props) Equal(x Expr) bool {
	y, ok := x.(*Props)
	return ok && pr.eventVar.Equal(y.eventVar)
}

// Var returns the event var
func (pr *Props) Var() *Var { return pr.eventVar }

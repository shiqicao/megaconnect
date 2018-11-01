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
	"encoding/binary"
	"io"
	"math"
	"reflect"
	"sort"
)

// TODO - document binary format
const (
	magic uint32 = 0x6d656761
)

const (
	exprKindBool        = 0x00
	exprKindInt         = 0x01
	exprKindStr         = 0x02
	exprKindObj         = 0x03
	exprKindBinOp       = 0x04
	exprKindUniOp       = 0x05
	exprKindFuncCall    = 0x06
	exprKindObjAccessor = 0x07
	exprKindVar         = 0x08
	exprKindObjLit      = 0x09
	exprKindProp        = 0x0a
)

const (
	declKindMonitor = 0x00
	declKindAction  = 0x01
	declKindEvent   = 0x02
)

const (
	stmtKindFire = 0x00
)

const (
	typeKindInt     uint8 = 0x00
	typeKindBoolean uint8 = 0x01
	typeKindStr     uint8 = 0x02
	typeKindObj     uint8 = 0x03
)

const (
	eexprKindVar   uint8 = 0x00
	eexprKindBinOp uint8 = 0x01
)

// Encoder serializes workflow AST to binary format
type Encoder struct {
	writer     io.Writer
	sortObjKey bool
}

// NewEncoder creates a new instance of Encoder
func NewEncoder(w io.Writer, sortObjKey bool) *Encoder {
	return &Encoder{
		writer:     w,
		sortObjKey: sortObjKey,
	}
}

// EncodeMonitorDecl serializes a monitor declaration to binary format
func (e *Encoder) EncodeMonitorDecl(md *MonitorDecl) error {
	if err := e.encodeString(md.Name()); err != nil {
		return err
	}
	if err := e.EncodeExpr(md.Condition()); err != nil {
		return err
	}
	if err := e.encodeVarDecls(md.vars); err != nil {
		return err
	}
	if err := e.encodeFireStmt(md.event); err != nil {
		return err
	}
	return e.encodeString(md.chain)
}

func (e *Encoder) encodeEventExpr(eexpr EventExpr) error {
	kind, err := getEExprKind(eexpr)
	if err != nil {
		return err
	}
	e.encodeBigEndian(kind)
	switch eexpr := eexpr.(type) {
	case *EVar:
		return e.encodeString(eexpr.name)
	case *EBinOp:
		if err := e.encodeBigEndian(uint8(eexpr.op)); err != nil {
			return err
		}
		if err := e.encodeEventExpr(eexpr.left); err != nil {
			return err
		}
		return e.encodeEventExpr(eexpr.right)
	default:
		return ErrNotSupportedByType(eexpr)
	}
}

// EncodeExpr serializes `Expr` to binary format
// Expression are disjoint union of concrete constructs, expression kind is a one byte value encodes construct type.
// TODO: add expression kind code mapping in binary format spec
func (e *Encoder) EncodeExpr(expr Expr) error {
	kind, err := getExprKind(expr)
	if err != nil {
		return err
	}
	e.encodeBigEndian(kind)
	switch expr := expr.(type) {
	case *BoolConst:
		var value uint8
		if expr.Value() {
			value = 1
		}
		return e.encodeBigEndian(value)
	case *StrConst:
		return e.encodeString(expr.Value())
	case *IntConst:
		if err = e.encodeBigEndian(int8(expr.Value().Sign())); err != nil {
			return err
		}
		return e.encodeBytes(expr.Value().Bytes())
	case *ObjConst:
		return e.encodeObjConst(expr)
	case *BinOp:
		if err = e.encodeBigEndian(uint8(expr.Op())); err != nil {
			return err
		}
		if err = e.EncodeExpr(expr.Left()); err != nil {
			return err
		}
		return e.EncodeExpr(expr.Right())
	case *UniOp:
		if err = e.encodeBigEndian(uint8(expr.Op())); err != nil {
			return err
		}
		return e.EncodeExpr(expr.Operant())
	case *FuncCall:
		if err = e.encodeString(expr.Name()); err != nil {
			return err
		}
		e.encodeLengthI(len(expr.Args()))
		for _, arg := range expr.Args() {
			if err = e.EncodeExpr(arg); err != nil {
				return err
			}
		}
		e.encodeLengthI(len(expr.NamespacePrefix()))
		for _, n := range expr.NamespacePrefix() {
			if err = e.encodeString(n); err != nil {
				return err
			}
		}
		return nil
	case *ObjAccessor:
		if err = e.EncodeExpr(expr.Receiver()); err != nil {
			return err
		}
		return e.encodeString(expr.Field())
	case *Var:
		return e.encodeVar(expr)
	case *ObjLit:
		return e.encodeObjLit(expr)
	case *Prop:
		return e.encodeVar(expr.eventVar)
	}

	return ErrNotSupportedByType(expr)
}

func (e *Encoder) encodeVar(v *Var) error {
	return e.encodeString(v.name)
}

func (e *Encoder) encodeObjConst(o *ObjConst) error {
	keys := o.Fields()
	if e.sortObjKey {
		sort.Strings(keys)
	}
	fields := o.Value()
	e.encodeLengthI(len(fields))
	for _, key := range keys {
		if err := e.encodeString(key); err != nil {
			return err
		}
		if err := e.EncodeExpr(fields[key]); err != nil {
			return err
		}
	}
	return nil
}

func (e *Encoder) encodeObjLit(o *ObjLit) error {
	return e.encodeVarDecls(o.fields)
}

func (e *Encoder) encodeVarDecls(vars VarDecls) error {
	e.encodeLengthI(len(vars))
	for name, expr := range vars {
		if err := e.encodeString(name); err != nil {
			return err
		}
		if err := e.EncodeExpr(expr); err != nil {
			return err
		}
	}
	return nil
}

func (e *Encoder) encodeType(ty Type) error {
	switch ty := ty.(type) {
	case *PrimitiveType:
		switch ty.ty {
		case intTy:
			return e.encodeBigEndian(typeKindInt)
		case stringTy:
			return e.encodeBigEndian(typeKindStr)
		case booleanTy:
			return e.encodeBigEndian(typeKindBoolean)
		default:
			return &ErrNotSupported{Name: string(ty.ty)}
		}
	case *ObjType:
		if err := e.encodeBigEndian(typeKindObj); err != nil {
			return err
		}
		return e.encodeObjType(ty)
	default:
		return ErrNotSupportedByType(ty)
	}
}

func (e *Encoder) encodeObjType(ty *ObjType) error {
	len := len(ty.fields)
	e.encodeLengthI(len)
	for f, t := range ty.fields {
		if err := e.encodeString(f); err != nil {
			return err
		}
		if err := e.encodeType(t); err != nil {
			return err
		}
	}
	return nil
}

func (e *Encoder) encodeBigEndian(x interface{}) error {
	return binary.Write(e.writer, binary.BigEndian, x)
}

func (e *Encoder) encodeString(s string) error {
	return e.encodeBytes([]byte(s))
}

func (e *Encoder) encodeLengthI(n int) {
	e.encodeLength(uint64(n))
}

// A sequence of items are prefixed by its length. Length is encoded into variant size of bytes,
// the most significant bit of each byte indicates whether next byte is used in length encoding.
// Least 7 bits in each bytes contributes to length encoding.
// Example 1:
//    encodeLength(257) returns two bytes [10000001 00000010], the most significant bit in the
//    first byte is 1, so the next byte is also used in length encoding. The rest 7 bits in the
//    first bytes is 0000001 and the 7 bits from second bytes is 0000010. The length are decoded
//    as following:
//         len := 0000010 << 7 + 0000001 = 256 + 1 = 257
func (e *Encoder) encodeLength(n uint64) {
	for {
		r := n & uint64(127)
		n = n >> 7
		if n > 0 {
			r = r | 128
		}
		e.encodeBigEndian(uint8(r))
		if n == 0 {
			break
		}
	}
}

func (e *Encoder) encodeBytes(bytes []byte) error {
	e.encodeLength(uint64(len(bytes)))
	_, err := e.writer.Write(bytes)
	return err
}

// EncodeWorkflow encode an expression to binary format
func (e *Encoder) EncodeWorkflow(wf *WorkflowDecl) error {
	// encode magic number "mega"
	if err := e.encodeBigEndian(magic); err != nil {
		return err
	}
	// encode version
	if err := e.encodeBigEndian(wf.Version()); err != nil {
		return err
	}
	if err := e.encodeString(wf.Name()); err != nil {
		return err
	}
	e.encodeLengthI(len(wf.children))
	for _, child := range wf.children {
		kind, err := getDeclKind(child)
		if err != nil {
			return err
		}
		e.encodeBigEndian(kind)
		switch decl := child.(type) {
		case *MonitorDecl:
			if err = e.EncodeMonitorDecl(decl); err != nil {
				return err
			}
		case *ActionDecl:
			if err = e.encodeString(decl.name); err != nil {
				return err
			}
			if err = e.encodeEventExpr(decl.trigger); err != nil {
				return err
			}
			len := len(decl.body)
			e.encodeLengthI(len)
			for _, stmt := range decl.body {
				if err = e.encodeStmt(stmt); err != nil {
					return err
				}
			}
		case *EventDecl:
			if err = e.encodeString(decl.name); err != nil {
				return err
			}
			if err = e.encodeObjType(decl.ty); err != nil {
				return err
			}
		default:
			return ErrNotSupportedByType(child)
		}
	}
	return nil
}

func (e *Encoder) encodeFireStmt(f *Fire) error {
	if err := e.encodeString(f.eventName); err != nil {
		return err
	}
	return e.EncodeExpr(f.eventObj)
}

func (e *Encoder) encodeStmt(s Stmt) error {
	kind, err := getStmtKind(s)
	if err != nil {
		return err
	}
	e.encodeBigEndian(kind)
	switch s := s.(type) {
	case *Fire:
		return e.encodeFireStmt(s)
	default:
		return ErrNotSupportedByType(s)
	}
}

// EncodeExpr encodes an expression to binary format
func EncodeExpr(expr Expr) ([]byte, error) {
	return withByteBuffer(func(e *Encoder) error {
		return e.EncodeExpr(expr)
	})
}

// EncodeMonitorDecl encode a monitor declaration
func EncodeMonitorDecl(m *MonitorDecl) ([]byte, error) {
	return withByteBuffer(func(e *Encoder) error {
		return e.EncodeMonitorDecl(m)
	})
}

// EncodeObjConst serializes object value to binary format
func EncodeObjConst(o *ObjConst) ([]byte, error) {
	return withByteBuffer(func(e *Encoder) error {
		return e.encodeObjConst(o)
	})
}

// EncodeWorkflow serializes workflow declaration to binary format
func EncodeWorkflow(wf *WorkflowDecl) ([]byte, error) {
	return withByteBuffer(func(e *Encoder) error {
		return e.EncodeWorkflow(wf)
	})
}

func withByteBuffer(f func(e *Encoder) error) ([]byte, error) {
	var buf bytes.Buffer
	e := &Encoder{writer: &buf}
	if err := f(e); err != nil {
		return nil, err
	}
	return buf.Bytes(), nil
}

func getExprKind(expr Expr) (uint8, error) {
	switch expr.(type) {
	case *BoolConst:
		return exprKindBool, nil
	case *IntConst:
		return exprKindInt, nil
	case *StrConst:
		return exprKindStr, nil
	case *ObjConst:
		return exprKindObj, nil
	case *BinOp:
		return exprKindBinOp, nil
	case *UniOp:
		return exprKindUniOp, nil
	case *FuncCall:
		return exprKindFuncCall, nil
	case *ObjAccessor:
		return exprKindObjAccessor, nil
	case *Var:
		return exprKindVar, nil
	case *ObjLit:
		return exprKindObjLit, nil
	case *Prop:
		return exprKindProp, nil
	}
	return math.MaxUint8, &ErrNotSupported{Name: reflect.TypeOf(expr).String()}
}

func getDeclKind(decl Decl) (uint8, error) {
	switch decl.(type) {
	case *MonitorDecl:
		return declKindMonitor, nil
	case *ActionDecl:
		return declKindAction, nil
	case *EventDecl:
		return declKindEvent, nil
	default:
		return math.MaxUint8, &ErrNotSupported{Name: reflect.TypeOf(decl).String()}
	}
}

func getStmtKind(stmt Stmt) (uint8, error) {
	switch stmt.(type) {
	case *Fire:
		return stmtKindFire, nil
	default:
		return math.MaxUint8, &ErrNotSupported{Name: reflect.TypeOf(stmt).String()}
	}
}

// getEExprKind returns event expr kind according to event expr type
func getEExprKind(eexpr EventExpr) (uint8, error) {
	switch eexpr.(type) {
	case *EVar:
		return eexprKindVar, nil
	case *EBinOp:
		return eexprKindBinOp, nil
	default:
		return math.MaxUint8, ErrNotSupportedByType(eexpr)
	}
}

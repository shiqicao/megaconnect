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
	"errors"
	"fmt"
	"io"
	"math"
	"math/big"
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
	return e.encodeVarDecls(md.vars)
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
		keys := expr.Fields()
		if e.sortObjKey {
			sort.Strings(keys)
		}
		fields := expr.Value()
		e.encodeLengthI(len(fields))
		for _, key := range keys {
			if err = e.encodeString(key); err != nil {
				return err
			}
			if err = e.EncodeExpr(fields[key]); err != nil {
				return err
			}
		}
		return nil
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
		return e.encodeString(expr.Name())
	case *ObjLit:
		return e.encodeObjLit(expr)
	}

	return ErrNotSupportedByType(expr)
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

func (d *Decoder) decodeLength() (uint64, error) {
	var r uint64
	var i uint8
	for ; ; i++ {
		n, err := d.decodeUint8()
		if err != nil {
			return 0, err
		}
		r = r | (uint64(n&127) << (i * 7))
		if n > 127 {
			continue
		}
		break
	}
	return r, nil
}

func (e *Encoder) encodeBytes(bytes []byte) error {
	e.encodeLength(uint64(len(bytes)))
	_, err := e.writer.Write(bytes)
	return err
}

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
			continue
		case *ActionDecl:
			if err = e.encodeString(decl.name); err != nil {
				return err
			}
			if err = e.EncodeExpr(decl.trigger); err != nil {
				return err
			}
			len := len(decl.run)
			e.encodeLengthI(len)
			for _, stmt := range decl.run {
				if err = e.encodeStmt(stmt); err != nil {
					return err
				}
			}
			continue
		case *EventDecl:
			if err = e.encodeString(decl.name); err != nil {
				return err
			}
			if err = e.encodeObjType(decl.ty); err != nil {
				return err
			}
			continue
		default:
			return ErrNotSupportedByType(child)
		}
	}
	return nil
}

func (e *Encoder) encodeStmt(s Stmt) error {
	kind, err := getStmtKind(s)
	if err != nil {
		return err
	}
	e.encodeBigEndian(kind)
	switch s := s.(type) {
	case *Fire:
		if err := e.encodeString(s.eventName); err != nil {
			return err
		}
		return e.EncodeExpr(s.eventObj)
	default:
		return ErrNotSupportedByType(s)
	}
}

// EncodeExpr encodes an expression to binary format
func EncodeExpr(expr Expr) ([]byte, error) {
	var buf bytes.Buffer
	e := &Encoder{writer: &buf}
	if err := e.EncodeExpr(expr); err != nil {
		return nil, err
	}
	return buf.Bytes(), nil
}

// EncodeMonitorDecl encode a monitor declaration
func EncodeMonitorDecl(m *MonitorDecl) ([]byte, error) {
	var buf bytes.Buffer
	e := &Encoder{writer: &buf}
	if err := e.EncodeMonitorDecl(m); err != nil {
		return nil, err
	}
	return buf.Bytes(), nil
}

// Decoder deserializes binary format to workflow AST
type Decoder struct {
	reader io.Reader
}

// NewByteDecoder creates a new Decoder from bytes
func NewByteDecoder(data []byte) *Decoder { return NewDecoder(bytes.NewReader(data)) }

// NewDecoder creates a new Decoder
func NewDecoder(reader io.Reader) *Decoder {
	return &Decoder{
		reader: reader,
	}
}

// DecodeWorkflow encodes a workflow declaration to binary format
func (d *Decoder) DecodeWorkflow() (*WorkflowDecl, error) {
	mega, err := d.decodeUint32()
	if err != nil {
		return nil, err
	} else if mega != magic {
		return nil, fmt.Errorf("Not Mega workflow binary")
	}
	version, err := d.decodeUint32()
	if err != nil {
		return nil, err
	}
	name, err := d.decodeBytes()
	if err != nil {
		return nil, err
	}
	wf := NewWorkflowDecl(string(name), version)
	len, err := d.decodeLength()
	if err != nil {
		return nil, err
	}
	for ; len > 0; len-- {
		kind, err := d.decodeUint8()
		if err != nil {
			return nil, err
		}
		switch kind {
		case declKindMonitor:
			m, err := d.DecodeMonitorDecl()
			if err != nil {
				return nil, err
			}
			wf.AddChildren(m)
			continue
		case declKindAction:
			name, err := d.decodeBytes()
			if err != nil {
				return nil, err
			}
			expr, err := d.DecodeExpr()
			if err != nil {
				return nil, err
			}
			len, err := d.decodeLength()
			if err != nil {
				return nil, err
			}
			stmts := make(Stmts, len)
			for ; len > 0; len-- {
				stmts[len-1], err = d.decodeStmt()
				if err != nil {
					return nil, err
				}
			}
			wf.AddChildren(NewActionDecl(string(name), expr, stmts))
			continue
		case declKindEvent:
			name, err := d.decodeBytes()
			if err != nil {
				return nil, err
			}
			objTy, err := d.decodeObjType()
			if err != nil {
				return nil, err
			}
			wf.AddChildren(NewEventDecl(string(name), objTy))
			continue
		default:
			return nil, &ErrNotSupported{Name: string(kind)}
		}
	}
	return wf, nil
}

// DecodeMonitorDecl deserializes binary format to `MonitorDecl`
func (d *Decoder) DecodeMonitorDecl() (*MonitorDecl, error) {
	name, err := d.decodeBytes()
	if err != nil {
		return nil, err
	}
	cond, err := d.DecodeExpr()
	if err != nil {
		return nil, err
	}
	vars, err := d.decodeVarDecls()
	if err != nil {
		return nil, err
	}
	return NewMonitorDecl(string(name), cond, vars), nil
}

// DecodeExpr deserializes binary format to `Expr`
func (d *Decoder) DecodeExpr() (Expr, error) {
	kind, err := d.decodeUint8()
	if err != nil {
		return nil, err
	}
	switch kind {
	case exprKindBool:
		var value bool
		err := binary.Read(d.reader, binary.BigEndian, &value)
		if err != nil {
			return nil, err
		}
		return GetBoolConst(value), nil
	case exprKindInt:
		sign, err := d.decodeInt8()
		if err != nil {
			return nil, err
		}
		bytes, err := d.decodeBytes()
		if err != nil {
			return nil, err
		}
		value := new(big.Int).SetBytes(bytes)
		if sign < 0 {
			value.Neg(value)
		}
		return NewIntConst(value), nil
	case exprKindStr:
		bytes, err := d.decodeBytes()
		if err != nil {
			return nil, err
		}
		return NewStrConst(string(bytes)), nil
	case exprKindObj:
		values := map[string]Const{}
		len, err := d.decodeLength()
		if err != nil {
			return nil, err
		}
		for ; len > 0; len-- {
			name, err := d.decodeBytes()
			if err != nil {
				return nil, err
			}
			expr, err := d.DecodeExpr()
			if err != nil {
				return nil, err
			}
			obj, ok := expr.(Const)
			if !ok {
				return nil, &ErrConstExpected{Actual: reflect.TypeOf(expr).String()}
			}
			values[string(name)] = obj
		}
		return NewObjConst(values), nil
	case exprKindBinOp:
		op, err := d.decodeUint8()
		if err != nil {
			return nil, err
		}
		left, err := d.DecodeExpr()
		if err != nil {
			return nil, err
		}
		right, err := d.DecodeExpr()
		if err != nil {
			return nil, err
		}
		return NewBinOp(Operator(op), left, right), nil
	case exprKindUniOp:
		op, err := d.decodeUint8()
		if err != nil {
			return nil, err
		}
		operant, err := d.DecodeExpr()
		if err != nil {
			return nil, err
		}
		return NewUniOp(Operator(op), operant), nil
	case exprKindFuncCall:
		name, err := d.decodeBytes()
		if err != nil {
			return nil, err
		}
		argsLen, err := d.decodeLength()
		if err != nil {
			return nil, err
		}
		args := []Expr{}
		for ; argsLen > 0; argsLen-- {
			arg, err := d.DecodeExpr()
			if err != nil {
				return nil, err
			}
			args = append(args, arg)
		}

		nsLen, err := d.decodeLength()
		if err != nil {
			return nil, err
		}
		nss := NamespacePrefix{}
		for ; nsLen > 0; nsLen-- {
			ns, err := d.decodeBytes()
			if err != nil {
				return nil, err
			}
			nss = append(nss, string(ns))
		}
		return NewFuncCall(nss, string(name), args...), nil
	case exprKindObjAccessor:
		receiver, err := d.DecodeExpr()
		if err != nil {
			return nil, err
		}
		field, err := d.decodeBytes()
		if err != nil {
			return nil, err
		}
		return NewObjAccessor(receiver, string(field)), nil
	case exprKindVar:
		name, err := d.decodeBytes()
		if err != nil {
			return nil, err
		}
		return NewVar(string(name)), nil
	case exprKindObjLit:
		vars, err := d.decodeVarDecls()
		if err != nil {
			return nil, err
		}
		return NewObjLit(vars), nil
	}
	return nil, &ErrNotSupported{Name: string(kind)}
}

func (d *Decoder) decodeVarDecls() (VarDecls, error) {
	len, err := d.decodeLength()
	if err != nil {
		return nil, err
	}
	vars := make(VarDecls, len)
	for ; len > 0; len-- {
		varName, err := d.decodeBytes()
		if err != nil {
			return nil, err
		}
		expr, err := d.DecodeExpr()
		if err != nil {
			return nil, err
		}
		vars[string(varName)] = expr
	}
	return vars, nil
}

func (d *Decoder) decodeType() (Type, error) {
	kind, err := d.decodeUint8()
	if err != nil {
		return nil, err
	}
	switch kind {
	case typeKindInt:
		return IntType, nil
	case typeKindBoolean:
		return BoolType, nil
	case typeKindStr:
		return StrType, nil
	case typeKindObj:
		return d.decodeObjType()
	default:
		return nil, &ErrNotSupported{Name: string(kind)}
	}
}

func (d *Decoder) decodeObjType() (*ObjType, error) {
	len, err := d.decodeLength()
	if err != nil {
		return nil, err
	}
	fields := make(ObjFieldTypes, len)
	for ; len > 0; len-- {
		name, err := d.decodeBytes()
		if err != nil {
			return nil, err
		}
		fieldType, err := d.decodeType()
		if err != nil {
			return nil, err
		}
		fields[string(name)] = fieldType
	}
	return NewObjType(fields), nil
}

func (d *Decoder) decodeUint32() (uint32, error) {
	var x uint32
	err := binary.Read(d.reader, binary.BigEndian, &x)
	return x, err
}

func (d *Decoder) decodeUint8() (uint8, error) {
	var x uint8
	err := binary.Read(d.reader, binary.BigEndian, &x)
	return x, err
}

func (d *Decoder) decodeInt8() (int8, error) {
	var x int8
	err := binary.Read(d.reader, binary.BigEndian, &x)
	return x, err
}

func (d *Decoder) decodeBytes() ([]byte, error) {
	len, err := d.decodeLength()
	if err != nil {
		return nil, err
	}
	bytes := make([]byte, len)
	n, err := d.reader.Read(bytes)
	if err != nil {
		return nil, err
	} else if uint64(n) != len {
		return nil, errors.New("")
	}
	return bytes, nil
}

func (d *Decoder) decodeStmt() (Stmt, error) {
	kind, err := d.decodeUint8()
	if err != nil {
		return nil, err
	}
	switch kind {
	case stmtKindFire:
		eventName, err := d.decodeBytes()
		if err != nil {
			return nil, err
		}
		eventObj, err := d.DecodeExpr()
		if err != nil {
			return nil, err
		}
		return NewFire(string(eventName), eventObj), nil
	default:
		return nil, &ErrNotSupported{Name: string(kind)}
	}
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

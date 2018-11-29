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
	"math/big"
	"reflect"
)

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

// DecodeNamespace decodes a namespace declaration to workflow lang
func (d *Decoder) DecodeNamespace() (*NamespaceDecl, error) {
	name, err := d.decodeBytes()
	if err != nil {
		return nil, err
	}
	ns := NewNamespaceDecl(string(name))
	len, err := d.decodeLength()
	if err != nil {
		return nil, err
	}
	for ; len > 0; len-- {
		c, err := d.DecodeNamespace()
		if err != nil {
			return nil, err
		}
		ns.AddNamespace(c)
	}
	len, err = d.decodeLength()
	if err != nil {
		return nil, err
	}
	for ; len > 0; len-- {
		fun, err := d.decodeFuncSig()
		if err != nil {
			return nil, err
		}
		ns.AddFunc(fun)
	}
	return ns, nil
}

// DecodeActionDecl decodes an action declaration to workflow lang AST
func (d *Decoder) DecodeActionDecl() (*ActionDecl, error) {
	name, err := d.decodeBytes()
	if err != nil {
		return nil, err
	}
	expr, err := d.decodeEventExpr()
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
	return NewActionDecl(NewIdB(name), expr, stmts), nil
}

func (d *Decoder) decodeFuncSig() (*FuncDecl, error) {
	name, err := d.decodeBytes()
	if err != nil {
		return nil, err
	}
	len, err := d.decodeLength()
	if err != nil {
		return nil, err
	}
	params := make(Params, 0, len)
	for ; len > 0; len-- {
		param, err := d.decodeBytes()
		if err != nil {
			return nil, err
		}
		ty, err := d.decodeType()
		if err != nil {
			return nil, err
		}
		params = append(params, NewParamDecl(string(param), ty))
	}
	retTy, err := d.decodeType()
	if err != nil {
		return nil, err
	}
	return NewFuncDecl(string(name), params, retTy, nil), nil
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
	wf := NewWorkflowDecl(NewIdB(name), version)
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
			wf.AddDecl(m)
			continue
		case declKindAction:
			a, err := d.DecodeActionDecl()
			if err != nil {
				return nil, err
			}
			wf.AddDecl(a)
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
			wf.AddDecl(NewEventDecl(NewIdB(name), objTy))
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
	event, err := d.decodeFireStmt()
	if err != nil {
		return nil, err
	}
	chain, err := d.decodeBytes()
	if err != nil {
		return nil, err
	}
	return NewMonitorDecl(NewIdB(name), cond, vars, event, string(chain)), nil
}

func (d *Decoder) decodeEventExpr() (EventExpr, error) {
	kind, err := d.decodeUint8()
	if err != nil {
		return nil, err
	}
	switch kind {
	case eexprKindVar:
		name, err := d.decodeBytes()
		if err != nil {
			return nil, err
		}
		return NewEVar(string(name)), nil
	case eexprKindBinOp:
		op, err := d.decodeUint8()
		if err != nil {
			return nil, err
		}
		l, err := d.decodeEventExpr()
		if err != nil {
			return nil, err
		}
		r, err := d.decodeEventExpr()
		if err != nil {
			return nil, err
		}
		return NewEBinOp(EventExprOperator(op), l, r), nil
	default:
		return nil, &ErrNotSupported{Name: string(kind)}
	}
}

func (d *Decoder) decodeInt() (*big.Int, error) {
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
	return value, nil
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
		i, err := d.decodeInt()
		if err != nil {
			return nil, err
		}
		return NewIntConst(i), nil
	case exprKindRat:
		num, err := d.decodeInt()
		if err != nil {
			return nil, err
		}
		bytes, err := d.decodeBytes()
		if err != nil {
			return nil, err
		}
		den := new(big.Int).SetBytes(bytes)
		return NewRatConst(new(big.Rat).SetFrac(num, den)), nil
	case exprKindStr:
		bytes, err := d.decodeBytes()
		if err != nil {
			return nil, err
		}
		return NewStrConst(string(bytes)), nil
	case exprKindObj:
		return d.DecodeObjConst()
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
			nss = append(nss, NewIdB(ns))
		}
		return NewFuncCall(nss, NewIdB(name), args...), nil
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
		return d.decodeVar()
	case exprKindObjLit:
		vars, err := d.decodeVarDecls()
		if err != nil {
			return nil, err
		}
		return NewObjLit(vars), nil
	case exprKindProps:
		v, err := d.decodeVar()
		if err != nil {
			return nil, err
		}
		return NewProps(v), nil
	}
	return nil, &ErrNotSupported{Name: string(kind)}
}

func (d *Decoder) decodeVar() (*Var, error) {
	name, err := d.decodeBytes()
	if err != nil {
		return nil, err
	}
	return NewVar(string(name)), nil
}

// DecodeObjConst decodes an object value from binary
func (d *Decoder) DecodeObjConst() (*ObjConst, error) {
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
}

func (d *Decoder) decodeVarDecls() (IdToExpr, error) {
	len, err := d.decodeLength()
	if err != nil {
		return nil, err
	}
	vars := make(IdToExpr, len)
	for ; len > 0; len-- {
		varName, err := d.decodeBytes()
		if err != nil {
			return nil, err
		}
		expr, err := d.DecodeExpr()
		if err != nil {
			return nil, err
		}
		vars.Put(string(varName), expr)
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
	case typeKindRat:
		return RatType, nil
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
	fields := make(IdToTy, len)
	for ; len > 0; len-- {
		name, err := d.decodeBytes()
		if err != nil {
			return nil, err
		}
		fieldType, err := d.decodeType()
		if err != nil {
			return nil, err
		}
		fields.Put(string(name), fieldType)
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
	if err == io.EOF && uint64(n) == len {
		return bytes, nil
	}
	if err != nil {
		return nil, err
	} else if uint64(n) != len {
		return nil, errors.New("")
	}
	return bytes, nil
}

func (d *Decoder) decodeFireStmt() (*Fire, error) {
	eventName, err := d.decodeBytes()
	if err != nil {
		return nil, err
	}
	eventObj, err := d.DecodeExpr()
	if err != nil {
		return nil, err
	}
	return NewFire(string(eventName), eventObj), nil
}

func (d *Decoder) decodeStmt() (Stmt, error) {
	kind, err := d.decodeUint8()
	if err != nil {
		return nil, err
	}
	switch kind {
	case stmtKindFire:
		return d.decodeFireStmt()
	default:
		return nil, &ErrNotSupported{Name: string(kind)}
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

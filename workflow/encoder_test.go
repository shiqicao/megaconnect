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
	"math"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestFuncCallEncoding(t *testing.T) {
	assertExprEncoding(t, NewFuncCall(nil, "T", TrueConst))
	assertExprEncoding(t, NewFuncCall(nil, "T", TrueConst))
	assertExprEncoding(t, NewFuncCall(NamespacePrefix{"a", "b"}, "T", TrueConst))
	assertExprEncoding(t, NewFuncCall(nil, "T", NewFuncCall(nil, "T", TrueConst)))
	assertExprEncoding(t, NewFuncCall(
		NamespacePrefix{"a", "b"},
		"T",
		NewFuncCall(NamespacePrefix{"a", "b"}, "T", GetBoolConst(true)),
		GetBoolConst(true),
	))
}

func TestLengthEncoding(t *testing.T) {
	test := func(x int) {
		withGen(func(ge genEncoder, gd genDecoder) {
			ge().encodeLengthI(x)
			n, err := gd().decodeLength()
			assert.NoError(t, err)
			assert.Equal(t, uint64(x), n)
		})
	}
	test(0)
	test(1)
	test(256)
	test(123 << 10)
}

func TestObjEncoding(t *testing.T) {
	test := func(c *ObjConst) {
		assertExprEncoding(t, c)
	}
	test(
		NewObjConst(
			map[string]Const{
				"a": NewIntConstFromI64(1),
			},
		),
	)
	test(
		NewObjConst(
			map[string]Const{
				"a": NewIntConstFromI64(1),
				"b": NewStrConst("Obj Test"),
			},
		),
	)
	test(
		NewObjConst(
			map[string]Const{
				"a": NewObjConst(
					map[string]Const{
						"a": GetBoolConst(true),
					},
				),
			},
		),
	)
}

func TestObjAccessorEncoding(t *testing.T) {
	assertExprEncoding(t,
		NewObjAccessor(
			NewFuncCall(nil, "Test", TrueConst),
			"a",
		))
	assertExprEncoding(t,
		NewObjAccessor(
			NewObjConst(
				ObjFields{
					"a": NewIntConstFromI64(0),
				}),
			"a",
		))
}

func TestConstEncoding(t *testing.T) {
	assertExprEncoding(t, GetBoolConst(true))
	assertExprEncoding(t, GetBoolConst(false))
	assertExprEncoding(t, NewStrConst("megaspace rock!"))
	assertExprEncoding(t, NewStrConst(""))
	assertExprEncoding(t, NewIntConstFromI64(12345))
	assertExprEncoding(t, NewIntConstFromI64(0))
	assertExprEncoding(t, NewIntConstFromI64(-1))
}

func TestVarEncoding(t *testing.T) {
	assertExprEncoding(t, NewVar("a"))
	assertExprEncoding(t, NewVar("abc"))
}

func TestObjLitEncoding(t *testing.T) {
	assertExprEncoding(t, NewObjLit(VD("a", NewIntConstFromI64(1))))
	assertExprEncoding(t, NewObjLit(VD("a", TrueConst).Put("b", FalseConst)))
	assertExprEncoding(t, NewObjLit(VD("a", NewObjLit(VD("a", NewStrConst("x"))))))
}

func TestBinExpEncoding(t *testing.T) {
	assertExprEncoding(t, NewBinOp(EqualOp, GetBoolConst(true), GetBoolConst(false)))
	assertExprEncoding(t, NewBinOp(EqualOp, NewUniOp(NotOp, GetBoolConst(true)), GetBoolConst(false)))
	assertExprEncoding(t, NewUniOp(NotOp, NewBinOp(EqualOp, GetBoolConst(true), GetBoolConst(false))))
}

func TestLargeArrayEncoding(t *testing.T) {
	str := ""
	for i := 0; i <= math.MaxUint8; i++ {
		str = str + "a"
	}
	assertExprEncoding(t, NewStrConst(str))
	assertExprEncoding(t, NewStrConst(str+str))
	assertExprEncoding(t, NewStrConst(str+str+str))
}

func TestMonitorDeclEncoding(t *testing.T) {
	check := func(m *MonitorDecl) {
		withGen(func(ge genEncoder, gd genDecoder) {
			err := ge().EncodeMonitorDecl(m)
			assert.NoError(t, err)
			decoded, err := gd().DecodeMonitorDecl()
			assert.NoError(t, err)
			assert.True(t, m.Equal(decoded))
		})
	}

	check(MD("a", T, VD("x", F), NewFire("e", NewObjConst(ObjFields{"t": T})), "Eth"))
	check(MD("b", AND(T, F), VD("x", T), NewFire("e", NewObjConst(ObjFields{"t": T})), "Eth"))
}

func TestWorkflowEncoding(t *testing.T) {
	check := func(w *WorkflowDecl) {
		withGen(
			func(ge genEncoder, gd genDecoder) {
				e := ge()
				err := e.EncodeWorkflow(w)
				assert.NoError(t, err)
				d := gd()
				decoded, err := d.DecodeWorkflow()
				assert.NoError(t, err)
				assert.True(t, w.Equal(decoded))
			})
	}

	check(NewWorkflowDecl("a", 1).
		AddChild(MD("b", T, VD("x", F), NewFire("e", NewObjConst(ObjFields{"t": T})), "Eth")))
	check(NewWorkflowDecl("a", 1).
		AddChild(MD("b", T, VD("x", F), NewFire("e", NewObjConst(ObjFields{"t": T})), "Eth")).
		AddChild(MD("c", T, VD("x", F), NewFire("e", NewObjConst(ObjFields{"t": F})), "Eth")),
	)

	check(NewWorkflowDecl("a", 1).AddChild(NewActionDecl("b", EV("a"), Stmts{NewFire("c", NewObjConst(ObjFields{"d": TrueConst}))})))
	check(NewWorkflowDecl("a", 1).
		AddChild(MD("b", T, VD("x", F), NewFire("e", NewObjConst(ObjFields{"t": T})), "Eth")).
		AddChild(NewActionDecl("c", EAND(EV("a"), EV("b")), Stmts{NewFire("c", NewObjConst(ObjFields{"d": NewIntConstFromI64(1)}))})),
	)

	check(NewWorkflowDecl("a", 1).AddChild(NewEventDecl("b", NewObjType(VT("a", IntType)))))
	check(NewWorkflowDecl("a", 1).AddChild(NewEventDecl("b", NewObjType(VT("a", NewObjType(VT("a", StrType)))))))
	check(NewWorkflowDecl("a", 1).
		AddChild(NewActionDecl("c", EOR(EV("a"), EAND(EV("a"), EV("b"))), Stmts{NewFire("c", NewObjConst(ObjFields{"d": NewIntConstFromI64(1)}))})).
		AddChild(NewEventDecl("b", NewObjType(VT("b", BoolType).Put("a", NewObjType(VT("a", StrType)))))),
	)
}

func TestPropEncoding(t *testing.T) {
	assertExprEncoding(t, P("a"))
	assertExprEncoding(t, EQ(P("a"), P("b")))
}

func assertExprEncoding(t *testing.T, expr Expr) {
	withGen(
		func(ge genEncoder, gd genDecoder) {
			e := ge()
			err := e.EncodeExpr(expr)
			assert.NoError(t, err)

			d := gd()
			decodedExpr, err := d.DecodeExpr()
			assert.NoError(t, err)
			assert.True(t, expr.Equal(decodedExpr))
		})
}

type genEncoder func() *Encoder
type genDecoder func() *Decoder

// gengen generates a pair of en/de-coder generator.
// Apologizing for making it a 2nd order function, which is necessary for bytesGenGen,
// the decoder can't be created until all bytes written to the encoder
type gengen func() (genEncoder, genDecoder)

func bufferGenGen() (genEncoder, genDecoder) {
	var b bytes.Buffer
	return func() *Encoder { return &Encoder{writer: &b} }, func() *Decoder { return &Decoder{reader: &b} }
}

func bytesGenGen() (genEncoder, genDecoder) {
	var b bytes.Buffer
	return func() *Encoder { return &Encoder{writer: &b} }, func() *Decoder { return NewByteDecoder(b.Bytes()) }
}

// withGen executes t with multiple different implementation of io.Reader and io.Writer.
// Different reader/writer implementation behaves differently.
// For example, io.Reader([]bytes).Read() returns (0, EOF) if it tries to read 0 byte, whereas bytes.Buffer.Read() just return (0, nil)
func withGen(t func(ge genEncoder, gd genDecoder)) {
	gengens := []gengen{bufferGenGen, bytesGenGen}
	for _, gengen := range gengens {
		ge, gd := gengen()
		t(ge, gd)
	}
}

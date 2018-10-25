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
	var b bytes.Buffer
	e := &Encoder{writer: &b}
	d := &Decoder{reader: &b}

	test := func(x int) {
		e.encodeLengthI(x)
		n, err := d.decodeLength()
		assert.NoError(t, err)
		assert.Equal(t, uint64(x), n)
		b.Reset()
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
	assertExprEncoding(t, NewObjLit(VarDecls{"a": NewIntConstFromI64(1)}))
	assertExprEncoding(t, NewObjLit(VarDecls{"a": TrueConst, "b": FalseConst}))
	assertExprEncoding(t, NewObjLit(VarDecls{"a": NewObjLit(VarDecls{"a": NewStrConst("x")})}))
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
		var b bytes.Buffer
		e := Encoder{writer: &b}
		d := Decoder{reader: &b}

		err := e.EncodeMonitorDecl(m)
		assert.NoError(t, err)
		decoded, err := d.DecodeMonitorDecl()
		assert.NoError(t, err)
		assert.True(t, m.Equal(decoded))
	}

	check(NewMonitorDecl("a", GetBoolConst(true), VarDecls{"x": FalseConst}))
	check(NewMonitorDecl("b", NewBinOp(AndOp, GetBoolConst(true), GetBoolConst(false)), VarDecls{"x": TrueConst}))
}

func TestWorkflowEncoding(t *testing.T) {
	check := func(w *WorkflowDecl) {
		var b bytes.Buffer
		e := Encoder{writer: &b}
		d := Decoder{reader: &b}

		err := e.EncodeWorkflow(w)
		assert.NoError(t, err)
		decoded, err := d.DecodeWorkflow()
		assert.NoError(t, err)
		assert.True(t, w.Equal(decoded))
	}

	check(NewWorkflowDecl("a", 1).AddChild(NewMonitorDecl("b", GetBoolConst(true), VarDecls{"x": FalseConst})))
	check(NewWorkflowDecl("a", 1).
		AddChild(NewMonitorDecl("b", GetBoolConst(true), VarDecls{"x": FalseConst})).
		AddChild(NewMonitorDecl("c", GetBoolConst(true), VarDecls{"x": FalseConst})),
	)

	check(NewWorkflowDecl("a", 1).AddChild(NewActionDecl("b", NewEVar("a"), Stmts{NewFire("c", NewObjConst(ObjFields{"d": TrueConst}))})))
	check(NewWorkflowDecl("a", 1).
		AddChild(NewMonitorDecl("b", GetBoolConst(true), VarDecls{"x": FalseConst})).
		AddChild(NewActionDecl("c", NewEBinOp(AndEOp, NewEVar("a"), NewEVar("b")), Stmts{NewFire("c", NewObjConst(ObjFields{"d": NewIntConstFromI64(1)}))})),
	)

	check(NewWorkflowDecl("a", 1).AddChild(NewEventDecl("b", NewObjType(ObjFieldTypes{"a": IntType}))))
	check(NewWorkflowDecl("a", 1).AddChild(NewEventDecl("b", NewObjType(ObjFieldTypes{"a": NewObjType(ObjFieldTypes{"a": StrType})}))))
	check(NewWorkflowDecl("a", 1).
		AddChild(NewActionDecl("c", NewEBinOp(OrEOp, NewEVar("a"), NewEBinOp(AndEOp, NewEVar("a"), NewEVar("b"))), Stmts{NewFire("c", NewObjConst(ObjFields{"d": NewIntConstFromI64(1)}))})).
		AddChild(NewEventDecl("b", NewObjType(ObjFieldTypes{"b": BoolType, "a": NewObjType(ObjFieldTypes{"a": StrType})}))),
	)
}

func assertExprEncoding(t *testing.T, expr Expr) {
	var b bytes.Buffer

	e := &Encoder{writer: &b}
	err := e.EncodeExpr(expr)
	assert.NoError(t, err)

	d := &Decoder{reader: &b}
	decodedExpr, err := d.DecodeExpr()
	assert.NoError(t, err)
	assert.True(t, expr.Equal(decodedExpr))
}

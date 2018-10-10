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
	assertExprEncoding(t, NewFuncCall("Test", []Expr{
		GetBoolConst(true),
	}, nil))
	assertExprEncoding(t, NewFuncCall("Test", []Expr{
		GetBoolConst(true),
	}, []string{}))
	assertExprEncoding(t, NewFuncCall("Test", []Expr{
		GetBoolConst(true),
	}, []string{"a", "b"}))
	assertExprEncoding(t,
		NewFuncCall("Test", []Expr{
			NewFuncCall("Test", []Expr{GetBoolConst(true)}, nil),
		}, nil))
	assertExprEncoding(t, NewFuncCall("Test", []Expr{
		NewFuncCall("Test", []Expr{
			GetBoolConst(true),
		}, []string{"a", "b"}),
		GetBoolConst(true),
	}, []string{"a", "b"}))
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
			NewFuncCall("Test", []Expr{GetBoolConst(true)}, nil),
			"a",
		))
	assertExprEncoding(t,
		NewObjAccessor(
			NewObjConst(
				map[string]Const{
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

	check(NewMonitorDecl("a", GetBoolConst(true)))
	check(NewMonitorDecl("b", NewBinOp(AndOp, GetBoolConst(true), GetBoolConst(false))))
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

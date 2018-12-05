// Copyright 2018 @ MegaSpace

// The MegaSpace source code is distributed in the hope that it will be useful,
// but WITHOUT ANY WARRANTY; without even the implied warranty of
// MERCHANTABILITY or FITNESS FOR A PARTICULAR PURPOSE. See the
// GNU Lesser General Public License for more details.

// You should have received a copy of the GNU Lesser General Public License
// along with the MegaSpace source code. If not, see <http://www.gnu.org/licenses/>.

package workflow

import (
	"testing"

	c "github.com/megaspacelab/megaconnect/common"
	"github.com/stretchr/testify/assert"
)

type strs = []string
type edge = map[string]c.StringSet

func TestCycleDetector(t *testing.T) {
	ss := func(xs ...string) c.StringSet {
		r := c.StringSet{}
		for _, x := range xs {
			r.Add(x)
		}
		return r
	}
	assert.Nil(t, checkCycle(edge{"a": ss("b")}, strs{}, "a"))
	assert.Equal(t, strs{"a", "a"}, checkCycle(edge{"a": ss("a")}, strs{}, "a"))
	assert.Equal(t, strs{"a", "b", "a"}, checkCycle(edge{"a": ss("b"), "b": ss("a")}, strs{}, "a"))
	assert.Equal(t, strs{"a", "b", "a"}, checkCycle(edge{"a": ss("b", "c"), "b": ss("a", "c")}, strs{}, "a"))
}

func TestNoVarDep(t *testing.T) {
	// Test cycle detector core
	v := &NoVarRecursiveDefValidator{}

	// Test with monitor
	m := MD(ID("a"), T, nil, nil, "a")
	wf := W("a").AddDecl(m)
	assertNoErrs(t, v.Validate(wf))

	m.vars = I2E1("a", T)
	assertNoErrs(t, v.Validate(wf))

	m.vars = I2E1("a", V("b"))
	assertNoErrs(t, v.Validate(wf))

	m.vars = I2E2("a", V("b"), "b", T)
	assertNoErrs(t, v.Validate(wf))

	m.vars = I2E1("a", V("a"))
	assertErrs(t, 1, v.Validate(wf))

	m.vars = I2E2("a", V("b"), "b", V("b"))
	assertErrs(t, 1, v.Validate(wf))
}

func TestNoDupInWf(t *testing.T) {
	v := &NoDupNameInWorkflow{}
	ot := OT(IdToTy{})
	wf := W("a").AddDecl(ED("b", ot)).AddDecl(ED("c", ot))
	assertNoErrs(t, v.Validate(wf))

	wf = W("a").AddDecl(ED("b", ot)).AddDecl(ED("b", ot))
	assertErrs(t, 1, v.Validate(wf))

	wf = W("a").AddDecl(ED("b", ot)).AddDecl(MD(ID("b"), T, I2E, nil, ""))
	assertErrs(t, 1, v.Validate(wf))

	wf = W("a").AddDecl(ACT(ID("b"), EV("t"), Stmts{})).AddDecl(MD(ID("b"), T, I2E, nil, "")).
		AddDecl(ED("c", ot)).AddDecl(ED("c", ot))
	assertErrs(t, 2, v.Validate(wf))
}

func TestNoRecursiveAction(t *testing.T) {
	v := &NoRecursiveAction{}
	ot := OT(IdToTy{})
	wf := W("a").
		AddDecl(ED("b", ot)).
		AddDecl(ED("c", ot)).
		AddDecl(ACT(ID("c"), EV("b"), Stmts{NewFire("c", OL(I2E))}))
	assertNoErrs(t, v.Validate(wf))

	wf = W("a").
		AddDecl(ED("b", ot)).
		AddDecl(ACT(ID("c"), EV("b"), Stmts{NewFire("b", OL(I2E))}))
	assertErrs(t, 1, v.Validate(wf))

	wf = W("a").
		AddDecl(ED("b", ot)).
		AddDecl(ED("c", ot)).
		AddDecl(ACT(ID("a1"), EV("b"), Stmts{NewFire("c", OL(I2E))})).
		AddDecl(ACT(ID("a2"), EV("c"), Stmts{NewFire("b", OL(I2E))}))
	assertErrs(t, 2, v.Validate(wf))

	wf = W("a").
		AddDecl(ED("b", ot)).
		AddDecl(ED("c", ot)).
		AddDecl(ED("d", ot)).
		AddDecl(ACT(ID("a1"), EAND(EV("b"), EV("d")), Stmts{NewFire("c", OL(I2E))})).
		AddDecl(ACT(ID("a2"), EV("c"), Stmts{NewFire("b", OL(I2E))}))
	assertErrs(t, 2, v.Validate(wf))

	wf = W("a").
		AddDecl(ED("b", ot)).
		AddDecl(ED("c", ot)).
		AddDecl(ED("d", ot)).
		AddDecl(ACT(ID("a1"), EV("b"), Stmts{
			NewFire("c", OL(I2E)),
			NewFire("b", OL(I2E)),
		}))
	assertErrs(t, 1, v.Validate(wf))

	wf = W("a").
		AddDecl(ED("b", ot)).
		AddDecl(ED("c", ot)).
		AddDecl(ED("d", ot)).
		AddDecl(ACT(ID("a1"), EV("d"), Stmts{
			NewFire("c", OL(I2E)),
		})).
		AddDecl(ACT(ID("a1"), EV("b"), Stmts{
			NewFire("d", OL(I2E)),
		})).
		AddDecl(ACT(ID("a2"), EV("b"), Stmts{
			NewFire("c", OL(I2E)),
			NewFire("b", OL(I2E)),
		}))
	assertErrs(t, 1, v.Validate(wf))
}

// Copyright 2018 @ MegaSpace

// The MegaSpace source code is distributed in the hope that it will be useful,
// but WITHOUT ANY WARRANTY; without even the implied warranty of
// MERCHANTABILITY or FITNESS FOR A PARTICULAR PURPOSE. See the
// GNU Lesser General Public License for more details.

// You should have received a copy of the GNU Lesser General Public License
// along with the MegaSpace source code. If not, see <http://www.gnu.org/licenses/>.

package workflow

var (
	// EventExpr constructors
	EV   = NewEVar
	EOR  = func(x EventExpr, y EventExpr) EventExpr { return NewEBinOp(OrEOp, x, y) }
	EAND = func(x EventExpr, y EventExpr) EventExpr { return NewEBinOp(AndEOp, x, y) }
)

var (
	// Expr constructors
	T   = TrueConst
	F   = FalseConst
	BIN = func(op Operator) func(Expr, Expr) *BinOp {
		return func(x Expr, y Expr) *BinOp { return NewBinOp(op, x, y) }
	}
	AND  = BIN(AndOp)
	OR   = BIN(OrOp)
	EQ   = BIN(EqualOp)
	FIRE = NewFire
	MD   = NewMonitorDecl
	V    = NewVar
	P    = func(a string) *Props { return NewProps(V(a)) }
	// Variable declaration
	VD = func(f string, e Expr) IdToExpr { return make(IdToExpr).Add(f, e) }
	// Variable types
	VT = func(f string, ty Type) IdToTy { return make(IdToTy).Add(f, ty) }
)

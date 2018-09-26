// Copyright 2018 @ MegaSpace

// The MegaSpace source code is distributed in the hope that it will be useful,
// but WITHOUT ANY WARRANTY; without even the implied warranty of
// MERCHANTABILITY or FITNESS FOR A PARTICULAR PURPOSE. See the
// GNU Lesser General Public License for more details.

// You should have received a copy of the GNU Lesser General Public License
// along with the MegaSpace source code. If not, see <http://www.gnu.org/licenses/>.

package interpreter

type Interpreter struct {
	env *Env
}

func New(env *Env) *Interpreter {
	return &Interpreter{
		env: env,
	}
}

func (i *Interpreter) EvalExpr(expr *Expr) (*Const, error) {
	return &Const{}, nil
}

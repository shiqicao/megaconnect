// Copyright 2018 @ MegaSpace

// The MegaSpace source code is distributed in the hope that it will be useful,
// but WITHOUT ANY WARRANTY; without even the implied warranty of
// MERCHANTABILITY or FITNESS FOR A PARTICULAR PURPOSE. See the
// GNU Lesser General Public License for more details.

// You should have received a copy of the GNU Lesser General Public License
// along with the MegaSpace source code. If not, see <http://www.gnu.org/licenses/>.

package workflow

// Prelude provides build-in functions
var (
	ethNS = "Eth"

	prelude = []*NamespaceDecl{
		NewNamespaceDecl(ethNS).addFunc(
			NewFuncDecl(
				"GetBalance",
				[]*ParamDecl{
					NewParamDecl("addr", StrType),
					NewParamDecl("height", IntType),
				},
				IntType,
				func(env *Env, args map[string]Const) (Const, error) {
					addrRaw, ok := args["addr"]
					if !ok {
						// This should be checked by type checker
						return nil, &ErrMissingArg{ArgName: "addr", Func: "GetBalance"}
					}
					heightRaw, ok := args["height"]
					if !ok {
						return nil, &ErrMissingArg{ArgName: "height", Func: "GetBalance"}
					}
					addr, _ := addrRaw.(*StrConst)
					height, _ := heightRaw.(*IntConst)
					balance, err := env.chain.QueryAccountBalance(addr.Value(), height.Value())
					if err != nil {
						return nil, err
					}
					return NewIntConst(balance), nil
				},
			),
		).addFunc(
			NewFuncDecl(
				"GetBlock",
				[]*ParamDecl{},
				NewObjType(map[string]Type{
					"height": IntType,
				}),
				func(env *Env, args map[string]Const) (Const, error) {
					return NewObjConst(
						map[string]Const{
							"height": NewIntConst(env.CurrentBlock().Height()),
						},
					), nil
				},
			),
		),
	}
)

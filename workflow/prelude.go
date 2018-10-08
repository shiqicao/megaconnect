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
		&NamespaceDecl{
			name: ethNS,
			funs: []*FuncDecl{
				NewFuncDecl(
					"GetBalance",
					[]*ParamDecl{
						NewParamDecl("addr", StrType),
					},
					IntType,
					func(env *Env, args map[string]Const) (Const, error) {
						addrRaw, ok := args["addr"]
						if !ok {
							// This should be checked by type checker
							return nil, &ErrMissingArg{ArgName: "addr", Func: "GetBalance"}
						}
						addr, ok := addrRaw.(*StrConst)
						blockHash := env.CurrentBlock().Hash()
						balance, err := env.chain.QueryAccountBalance(addr.Value(), &blockHash)
						if err != nil {
							return nil, err
						}
						return NewIntConst(balance), nil
					},
				),
				NewFuncDecl(
					"Addr",
					[]*ParamDecl{
						NewParamDecl("addr", StrType),
					},
					NewObjType(map[string]Type{
						"balance": IntType,
					}),
					func(env *Env, args map[string]Const) (Const, error) {
						// TODO: read data from chain through env,
						// env should provide access to chain data
						result := NewObjConst(
							map[string]Const{
								"balance": NewIntConstFromI64(100),
							},
						)
						return result, nil
					},
				),
			},
		},
	}
)

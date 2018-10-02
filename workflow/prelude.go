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

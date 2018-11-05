// Copyright 2018 @ MegaSpace

// The MegaSpace source code is distributed in the hope that it will be useful,
// but WITHOUT ANY WARRANTY; without even the implied warranty of
// MERCHANTABILITY or FITNESS FOR A PARTICULAR PURPOSE. See the
// GNU Lesser General Public License for more details.

// You should have received a copy of the GNU Lesser General Public License
// along with the MegaSpace source code. If not, see <http://www.gnu.org/licenses/>.

package chainmanager

import (
	"math/big"

	"github.com/megaspacelab/megaconnect/common"
	wf "github.com/megaspacelab/megaconnect/workflow"
)

type chain interface {
	QueryAccountBalance(addr string, height *big.Int) (*big.Int, error)
}

type chainAPI struct {
	chain        chain
	lib          *wf.NamespaceDecl
	currentBlock common.Block
}

func newChainAPI(defaultNamesapce string, chain chain) *chainAPI {
	api := &chainAPI{chain: chain}

	api.lib = wf.NewNamespaceDecl(
		defaultNamesapce,
	).AddFunc(
		wf.NewFuncDecl(
			"GetBalance",
			[]*wf.ParamDecl{
				wf.NewParamDecl("addr", wf.StrType),
				wf.NewParamDecl("height", wf.IntType),
			},
			wf.IntType,
			func(env *wf.Env, args map[string]wf.Const) (wf.Const, error) {
				addrRaw, ok := args["addr"]
				if !ok {
					// This should be checked by type checker
					return nil, &wf.ErrMissingArg{ArgName: "addr", Func: "GetBalance"}
				}
				heightRaw, ok := args["height"]
				if !ok {
					return nil, &wf.ErrMissingArg{ArgName: "height", Func: "GetBalance"}
				}
				addr, _ := addrRaw.(*wf.StrConst)
				height, _ := heightRaw.(*wf.IntConst)
				balance, err := api.chain.QueryAccountBalance(addr.Value(), height.Value())
				if err != nil {
					return nil, err
				}
				return wf.NewIntConst(balance), nil
			},
		),
	).AddFunc(
		wf.NewFuncDecl(
			"GetBlock",
			[]*wf.ParamDecl{},
			wf.NewObjType(wf.NewIdToTy().Add(
				"height", wf.IntType,
			)),
			func(_ *wf.Env, args map[string]wf.Const) (wf.Const, error) {
				return wf.NewObjConst(
					map[string]wf.Const{
						"height": wf.NewIntConst(api.currentBlock.Height()),
					},
				), nil
			},
		),
	)
	return api
}

func (c *chainAPI) setCurrentBlock(block common.Block) *chainAPI {
	c.currentBlock = block
	return c
}

func (c *chainAPI) getAPI() *wf.NamespaceDecl { return c.lib }

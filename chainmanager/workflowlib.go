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
	Namespace() (*wf.NamespaceDecl, error)
}

type ChainAPI struct {
	chain        chain
	lib          *wf.NamespaceDecl
	currentBlock common.Block
}

func NewChainAPI(chain chain) (*ChainAPI, error) {
	api := &ChainAPI{chain: chain}

	// chain.Namespace() returns a copy
	connectorAPI, err := chain.Namespace()
	if err != nil {
		return nil, err
	}

	err = connectorAPI.AddFuncs(
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
		wf.NewFuncDecl(
			"GetBlock",
			[]*wf.ParamDecl{},
			wf.NewObjType(wf.NewIdToTy().Put(
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
	if err != nil {
		return nil, err
	}
	api.lib = connectorAPI
	return api, nil
}

func (c *ChainAPI) setCurrentBlock(block common.Block) *ChainAPI {
	c.currentBlock = block
	return c
}

func (c *ChainAPI) GetAPI() *wf.NamespaceDecl { return c.lib.Copy() }

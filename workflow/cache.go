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
	"fmt"

	"github.com/hashicorp/golang-lru/simplelru"
)

const (
	apiCacheSize = 1000
)

// Cache is a component for caching evaluation results
type Cache interface {
	getFuncCallResult(*FuncDecl, []Const, func() (Const, error)) (Const, error)
}

// FuncCallCache provides API caching
type FuncCallCache struct {
	cache simplelru.LRUCache
}

// NewFuncCallCache creates a new instace of APICache
func NewFuncCallCache(size int) (*FuncCallCache, error) {
	cache, err := simplelru.NewLRU(size, nil)
	if err != nil {
		return nil, err
	}
	return &FuncCallCache{
		cache: cache,
	}, nil
}

func (a *FuncCallCache) getFuncCallResult(
	fun *FuncDecl,
	args []Const,
	compute func() (Const, error),
) (Const, error) {
	if a == nil {
		return compute()
	}
	key, err := buildApiCacheKey(fun, args)
	if err != nil {
		return nil, err
	}
	if value, found := a.cache.Get(key); found {
		if value == nil {
			return nil, fmt.Errorf("Cache returns <nil> from %s with func: %#v, args: %#v ", key, fun, args)
		}
		return value.(Const), nil
	}
	value, err := compute()
	if err != nil {
		return nil, err
	}
	a.cache.Add(key, value)
	return value, nil
}

func buildApiCacheKey(funcDecl *FuncDecl, args []Const) (string, error) {
	var buf bytes.Buffer
	buf.WriteString(funcDecl.Name())
	next := funcDecl.Parent()
	for next != nil {
		// Add a separator to avoid ambiguity. For example
		// a.foo is different than af.oo
		buf.WriteByte(0)
		buf.WriteString(next.Name())
		next = next.Parent()
	}
	encoder := NewEncoder(&buf, true /*sort obj fields*/)
	for _, arg := range args {
		if err := encoder.EncodeExpr(arg); err != nil {
			return "", err
		}
	}
	return buf.String(), nil
}

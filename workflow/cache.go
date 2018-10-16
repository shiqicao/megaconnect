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

	"github.com/hashicorp/golang-lru/simplelru"
)

const (
	apiCacheSize = 1000
)

// Cache is a component for caching evaluation results
type Cache interface {
	funcCallCache(fun *FuncDecl, args []Const) (getter func() Const, setter func(Const))
}

// FuncCallCache provides API caching
type FuncCallCache struct {
	cache simplelru.LRUCache
}

// NewFuncCallCache creates a new instace of APICache
func NewFuncCallCache() (*FuncCallCache, error) {
	cache, err := simplelru.NewLRU(apiCacheSize, nil)
	if err != nil {
		return nil, err
	}
	return &FuncCallCache{
		cache: cache,
	}, nil
}

func (a *FuncCallCache) funcCallCache(fun *FuncDecl, args []Const) (func() Const, func(Const)) {
	if a == nil {
		return noopGetter, noopSetter
	}
	key, err := buildApiCacheKey(fun, args)
	if err != nil {
		return noopGetter, noopSetter
	}
	getter := func() Const {
		if value, _ := a.cache.Get(key); value != nil {
			return value.(Const)
		}
		return nil
	}
	setter := func(value Const) { a.cache.Add(key, value) }
	return getter, setter
}

func noopGetter() Const  { return nil }
func noopSetter(_ Const) { return }

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

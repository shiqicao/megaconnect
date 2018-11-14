// Copyright 2018 @ MegaSpace

// The MegaSpace source code is distributed in the hope that it will be useful,
// but WITHOUT ANY WARRANTY; without even the implied warranty of
// MERCHANTABILITY or FITNESS FOR A PARTICULAR PURPOSE. See the
// GNU Lesser General Public License for more details.

// You should have received a copy of the GNU Lesser General Public License
// along with the MegaSpace source code. If not, see <http://www.gnu.org/licenses/>.

package main

import (
	"fmt"
	"go/importer"
	"go/parser"
	"go/token"
	"go/types"
	"io"
	"os"
	"path/filepath"
	"strings"

	"go/ast"
)

func main() {
	fileAst := []*ast.File{}
	fset := token.NewFileSet()

	wfDir := "/Users/shiqicao/go/src/github.com/megaspacelab/megaconnect/workflow"
	filepath.Walk(wfDir, func(path string, info os.FileInfo, err error) error {
		if err != nil {
			return err
		}
		if info.IsDir() && info.Name() != "workflow" {
			fmt.Printf("skip dir: %s \n", info.Name())
			return filepath.SkipDir
		}
		if filepath.Ext(path) != ".go" {
			fmt.Printf("skip file: %s \n", info.Name())
			return nil
		}
		if strings.Contains(path, "test") {
			fmt.Printf("skip test: %s \n", info.Name())
			return nil
		}
		f, err := parser.ParseFile(fset, path, nil, 0)
		if err != nil {
			panic(err)
		}
		fileAst = append(fileAst, f)
		return nil
	})

	config := &types.Config{
		Error: func(e error) {
			fmt.Println(e)
		},
		Importer: importer.Default(),
	}

	info := types.Info{
		Types: make(map[ast.Expr]types.TypeAndValue),
		//Defs:  make(map[*ast.Ident]types.Object),
		//Uses:  make(map[*ast.Ident]types.Object),
	}

	pkg, err := config.Check("workflow", fset, fileAst, &info)
	if err != nil {
		fmt.Println(err)
	}

	n := pkg.Scope().Lookup("Node").Type().(*types.Named)
	fmt.Printf("ROOT: %s %v %T\n", "Node", n, n)
	tree := &tree{
		node:     n,
		children: []*tree{},
	}
	names := pkg.Scope().Names()
	for _, n := range names {
		obj, ok := pkg.Scope().Lookup(n).Type().(*types.Named)
		if ok {
			fmt.Printf("PROCESSING: %v \n", n)
			tree.Insert(obj)
		}
	}

	tree.Print(os.Stdout, "")

	fmt.Fprintln(os.Stdout, "done")

	e := pkg.Scope().Lookup("Expr").Type()
	fmt.Printf("%v \n", e)

	a := pkg.Scope().Lookup("node").Type() //.Underlying().(*types.Interface)
	fmt.Printf("%v \n", a)
	p := types.NewPointer(a)
	impl := types.Implements(p, n.Underlying().(*types.Interface))
	fmt.Printf("%v \n", impl)

}

type tree struct {
	node     *types.Named
	children []*tree
}

func getType(ty *types.Named) types.Type {
	if isStruct(ty) {
		return types.NewPointer(ty)
	}
	return ty.Underlying()
}

func isStruct(ty *types.Named) bool {
	_, ok := ty.Underlying().(*types.Struct)
	return ok
}

func (t *tree) Insert(ty *types.Named) bool {
	if ty.String() == t.node.String() {
		return false
	}
	newTy := getType(ty)
	if !types.Implements(newTy, t.node.Underlying().(*types.Interface)) {
		return false
	}
	newNode := &tree{node: ty, children: []*tree{}}
	tyInterface, ok := ty.Underlying().(*types.Interface)
	if !ok {
		for _, c := range t.children {
			if isStruct(c.node) {
				continue
			}
			ok := c.Insert(ty)
			if ok {
				return true
			}
		}
		t.children = append(t.children, newNode)
		return true
	}

	newChildren := []*tree{newNode}
	for _, c := range t.children {
		if types.Implements(getType(c.node), tyInterface) {
			newNode.children = append(newNode.children, c)
		} else {
			newChildren = append(newChildren, c)
		}
	}
	t.children = newChildren
	return true
}

func (t *tree) Print(w io.Writer, indent string) {
	fmt.Fprintf(w, indent)
	fmt.Fprintf(w, "%#v \n", t.node.String())
	for _, c := range t.children {
		c.Print(w, indent+"  ")
	}
}

// Copyright 2018 @ MegaSpace

// The MegaSpace source code is distributed in the hope that it will be useful,
// but WITHOUT ANY WARRANTY; without even the implied warranty of
// MERCHANTABILITY or FITNESS FOR A PARTICULAR PURPOSE. See the
// GNU Lesser General Public License for more details.

// You should have received a copy of the GNU Lesser General Public License
// along with the MegaSpace source code. If not, see <http://www.gnu.org/licenses/>.

// Package wfc is a workflow compiler which compiles script to binary format
package main

import (
	"encoding/hex"
	"encoding/json"
	"io"
	"os"

	wf "github.com/megaspacelab/eventmanager/workflow"
	cli "gopkg.in/urfave/cli.v2"
)

func main() {
	app := &cli.App{
		Name: "Workflow compiler",
		Flags: []cli.Flag{
			&cli.PathFlag{
				Name:    "output",
				Aliases: []string{"o"},
			},
		},
		Action: compile,
	}
	app.Run(os.Args)
}

func compile(ctx *cli.Context) error {
	output := ctx.Path("output")
	var binWriter io.Writer
	var metaWriter io.Writer
	if output == "" {
		metaWriter = os.Stdout
	} else {
		path := ctx.Path("output")
		metafs, err := os.Create(path + ".mast")
		defer metafs.Close()
		if err != nil {
			return err
		}
		metaWriter = metafs
		binfs, err := os.Create(path + ".bast")
		defer binfs.Close()
		if err != nil {
			return err
		}
		binWriter = binfs
	}

	// TODO: Read script file and parse it to AST
	expr := wf.NewFuncCall("GetBalance", wf.Args{wf.NewStrConst("0x01C797d1AD1b36FE4eB17d58c96D6E844cD70a6B")}, wf.NamespacePrefix{"Eth"})

	bin, err := wf.EncodeExpr(expr)
	if err != nil {
		return err
	}
	meta := struct {
		Source string
		Hex    string
	}{
		Source: expr.String(),
		Hex:    hex.EncodeToString(bin),
	}

	if err := json.NewEncoder(metaWriter).Encode(meta); err != nil {
		return err
	}
	if binWriter != nil {
		if _, err := binWriter.Write(bin); err != nil {
			return err
		}
	}
	return nil
}

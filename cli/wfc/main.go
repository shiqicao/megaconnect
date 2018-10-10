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

	wf "github.com/megaspacelab/megaconnect/workflow"
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
			&cli.StringFlag{
				Name:  "temporaryAddr",
				Usage: "A temporary parameter for pass in an address on a chain.",
				// default address is an Eth active address
				// https://etherscan.io/address/0xfbb1b73c4f0bda4f67dca266ce6ef42f520fbb98
				Value: "0xFBb1b73C4f0BDa4f67dcA266ce6Ef42f520fBB98",
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
	var monitorWriter io.Writer
	if output == "" {
		metaWriter = os.Stdout
	} else {
		path := ctx.Path("output")
		metafs, err := os.Create(path + ".json")
		if err != nil {
			return err
		}
		defer metafs.Close()
		metaWriter = metafs

		binfs, err := os.Create(path + ".bast")
		if err != nil {
			return err
		}
		defer binfs.Close()
		binWriter = binfs

		monitorfs, err := os.Create(path + ".monitor")
		if err != nil {
			return err
		}
		defer monitorfs.Close()
		monitorWriter = monitorfs
	}

	addr := ctx.String("temporaryAddr")
	// TODO: Read script file and parse it to AST
	expr := wf.NewBinOp(
		wf.NotEqualOp,
		wf.NewFuncCall(
			"GetBalance",
			wf.Args{
				wf.NewStrConst(addr),
				wf.NewObjAccessor(
					wf.NewFuncCall("GetBlock", wf.Args{}, wf.NamespacePrefix{"Eth"}),
					"height",
				),
			},
			wf.NamespacePrefix{"Eth"},
		),
		wf.NewFuncCall(
			"GetBalance",
			wf.Args{
				wf.NewStrConst(addr),
				wf.NewBinOp(wf.MinusOp,
					wf.NewObjAccessor(
						wf.NewFuncCall("GetBlock", wf.Args{}, wf.NamespacePrefix{"Eth"}),
						"height",
					),
					wf.NewIntConstFromI64(1),
				),
			},
			wf.NamespacePrefix{"Eth"},
		),
	)

	monitor := wf.NewMonitorDecl("Test", expr)

	bin, err := wf.EncodeMonitorDecl(monitor)
	hex := hex.EncodeToString(bin)
	if err != nil {
		return err
	}
	meta := struct {
		Source string
		Hex    string
	}{
		Source: monitor.String(),
		Hex:    hex,
	}

	if err := json.NewEncoder(metaWriter).Encode(meta); err != nil {
		return err
	}
	if binWriter != nil {
		if _, err := binWriter.Write(bin); err != nil {
			return err
		}
	}
	if monitorWriter != nil {
		if _, err := monitorWriter.Write([]byte(hex)); err != nil {
			return err
		}
	}
	return nil
}

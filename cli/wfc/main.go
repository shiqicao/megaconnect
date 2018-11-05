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
	p "path"

	"github.com/megaspacelab/megaconnect/unsafe"
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
	var name string
	var binWriter io.Writer
	var metaWriter io.Writer
	var workflowWriter io.Writer
	if output == "" {
		metaWriter = os.Stdout
	} else {
		path := output
		name = p.Base(output)
		metafs, err := os.Create(path + ".json")
		if err != nil {
			return err
		}
		defer metafs.Close()
		metaWriter = metafs

		binfs, err := os.Create(path)
		if err != nil {
			return err
		}
		defer binfs.Close()
		binWriter = binfs

		workflowfs, err := os.Create(path + ".wf")
		if err != nil {
			return err
		}
		defer workflowfs.Close()
		workflowWriter = workflowfs
	}

	addr := ctx.String("temporaryAddr")
	vars := wf.NewIdToExpr().Add(
		"blockHeight",
		wf.NewObjAccessor(
			wf.NewFuncCall(nil, "GetBlock"),
			"height",
		),
	)
	expr := wf.NewBinOp(
		wf.NotEqualOp,
		wf.NewFuncCall(
			nil,
			"GetBalance",
			wf.NewStrConst(addr),
			wf.NewVar("blockHeight"),
		),
		wf.NewFuncCall(
			nil,
			"GetBalance",
			wf.NewStrConst(addr),
			wf.NewBinOp(wf.MinusOp,
				wf.NewVar("blockHeight"),
				wf.NewIntConstFromI64(1),
			),
		),
	)

	monitorEth := wf.NewMonitorDecl(
		"EthMonitor",
		expr,
		vars,
		wf.NewFire("EthEvent", wf.NewObjLit(wf.NewIdToExpr().Add("height", wf.NewVar("blockHeight")))),
		"Ethereum",
	)

	monitorExample := wf.NewMonitorDecl(
		"ExampleMonitor",
		expr,
		vars,
		wf.NewFire("ExEvent", wf.NewObjLit(wf.NewIdToExpr().Add("height", wf.NewVar("blockHeight")))),
		"Example",
	)

	monitorBtc := wf.NewMonitorDecl(
		"BtcMonitor",
		expr,
		vars,
		wf.NewFire("BtcEvent", wf.NewObjLit(wf.NewIdToExpr().Add("height", wf.NewVar("blockHeight")))),
		"Bitcoin",
	)

	workflow := wf.NewWorkflowDecl(name, 0).
		AddChild(
			wf.NewEventDecl("TestEvent0", wf.NewObjType(
				wf.NewIdToTy().
					Add("example_h", wf.IntType).
					Add("eth_h", wf.IntType),
			)),
		).
		AddChild(
			wf.NewEventDecl(
				"HeightSumEvent",
				wf.NewObjType(wf.NewIdToTy().Add("heightSum", wf.IntType)),
			),
		).
		AddChild(
			wf.NewEventDecl("ExEvent", wf.NewObjType(wf.NewIdToTy().Add("height", wf.IntType))),
		).
		AddChild(
			wf.NewEventDecl("BtcEvent", wf.NewObjType(wf.NewIdToTy().Add("height", wf.IntType))),
		).
		AddChild(
			wf.NewEventDecl("EthEvent", wf.NewObjType(wf.NewIdToTy().Add("height", wf.IntType))),
		).
		AddChild(monitorEth).
		AddChild(monitorExample).
		AddChild(monitorBtc).
		AddChild(
			wf.NewActionDecl(
				"TestAction1",
				wf.NewEBinOp(wf.AndEOp, wf.NewEVar("ExEvent"), wf.NewEVar("EthEvent")),
				wf.Stmts{
					wf.NewFire(
						"TestEvent0",
						wf.NewObjLit(
							wf.NewIdToExpr().
								Add("example_h", wf.NewObjAccessor(wf.NewProps(wf.NewVar("ExEvent")), "height")).
								Add("eth_h", wf.NewObjAccessor(wf.NewProps(wf.NewVar("EthEvent")), "height")),
						),
					),
				}),
		).
		AddChild(
			wf.NewActionDecl(
				"JoinAction",
				wf.NewEVar("TestEvent0"),
				wf.Stmts{
					wf.NewFire(
						"HeightSumEvent",
						wf.NewObjLit(
							wf.NewIdToExpr().Add(
								"heightSum",
								wf.NewBinOp(
									wf.PlusOp,
									wf.NewObjAccessor(wf.NewProps(wf.NewVar("TestEvent0")), "eth_h"),
									wf.NewObjAccessor(wf.NewProps(wf.NewVar("TestEvent0")), "example_h"),
								),
							),
						),
					),
				},
			),
		)

	// TODO: Read script file and parse it to AST
	bin, err := wf.EncodeWorkflow(workflow)
	if err != nil {
		return err
	}
	hex := hex.EncodeToString(bin)
	meta := struct {
		Source string
		Hex    string
	}{
		Source: workflow.String(),
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
	if workflowWriter != nil {
		if _, err := workflowWriter.Write(unsafe.StringToBytes(hex)); err != nil {
			return err
		}
	}
	return nil
}

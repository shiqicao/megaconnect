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
	"fmt"
	"os"
	"path/filepath"

	p "github.com/megaspacelab/megaconnect/prettyprint"

	wf "github.com/megaspacelab/megaconnect/workflow"
	"github.com/megaspacelab/megaconnect/workflow/parser"
	cli "gopkg.in/urfave/cli.v2"
)

func main() {

	output := &cli.PathFlag{
		Name:    "output",
		Aliases: []string{"o"},
	}

	app := &cli.App{
		Name:   "Workflow compiler",
		Flags:  []cli.Flag{output},
		Action: compile,
		Commands: []*cli.Command{
			&cli.Command{
				Name:   "reflect",
				Action: decompile,
				Flags:  []cli.Flag{output},
			},
		},
		Writer:    os.Stdin,
		ErrWriter: os.Stderr,
	}
	app.Run(os.Args)
}

func exit(err error) cli.ExitCoder { return cli.Exit(err.Error(), 1) }

func decompile(ctx *cli.Context) error {
	if ctx.Args().Len() > 0 {
		binfile := ctx.Args().Get(0)
		fs, err := os.Open(binfile)
		if err != nil {
			return exit(err)
		}
		decoder := wf.NewDecoder(fs)
		wf, err := decoder.DecodeWorkflow()
		if err != nil {
			return exit(err)
		}
		output := ctx.Path("output")
		if output == "" {
			output = binfile + ".wf"
		}
		outputfs, err := os.Create(output)
		if err != nil {
			return exit(err)
		}
		err = wf.Print()(p.NewTxtPrinter(outputfs))
		if err != nil {
			return exit(err)
		}
		return nil
	}
	return exit(fmt.Errorf("Missing binary file path"))
}

func compile(ctx *cli.Context) error {
	if ctx.Args().Len() > 0 {
		src := ctx.Args().Get(0)
		w, err := parser.Parse(src)
		if err != nil {
			fmt.Printf(err.Error())
			return err
		}

		// TODO: validation & type checker

		bin, err := wf.EncodeWorkflow(w)
		if err != nil {
			return err
		}
		output := ctx.Path("output")
		if output == "" {
			ext := filepath.Ext(src)
			output = src[0 : len(src)-len(ext)]
		}
		fmt.Printf("output: %s", output)
		binWriter, err := os.Create(output)
		defer binWriter.Close()
		if err != nil {
			return err
		}
		if _, err = binWriter.Write(bin); err != nil {
			return err
		}
		return nil
	}

	return exit(fmt.Errorf("Missing source file"))
}

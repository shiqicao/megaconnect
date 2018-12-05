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
	"context"
	"encoding/hex"
	"fmt"
	"io"
	"io/ioutil"
	"os"
	"path/filepath"

	"github.com/megaspacelab/megaconnect/workflow/compiler"

	p "github.com/megaspacelab/megaconnect/prettyprint"
	"go.uber.org/zap"
	"google.golang.org/grpc"

	wf "github.com/megaspacelab/megaconnect/workflow"
	cli "gopkg.in/urfave/cli.v2"

	mcli "github.com/megaspacelab/megaconnect/cli"

	mgrpc "github.com/megaspacelab/megaconnect/grpc"
)

func main() {

	output := &cli.PathFlag{
		Name:    "output",
		Aliases: []string{"o"},
	}

	wmAddr := &cli.StringFlag{
		Name:  "wm-addr",
		Usage: "workflow manager address",
		Value: "localhost:9000",
	}

	wfid := &cli.StringFlag{
		Name:  "workflow-id",
		Usage: "workflow identifier for undeploy",
	}

	lib := &cli.PathFlag{
		Name:  "lib",
		Usage: "workflow APIs directory",
		Value: mcli.DefaultWorkflowLibDir(),
	}

	app := &cli.App{
		Name: "Workflow compiler",
		Commands: []*cli.Command{
			&cli.Command{
				Name:   "compile",
				Usage:  "compile workflow source code",
				Action: mcli.ToExitCode(compile),
				Flags:  []cli.Flag{&mcli.DebugFlag, output, lib},
			},
			&cli.Command{
				Name:   "reflect",
				Usage:  "decompile a binary to workflow",
				Action: mcli.ToExitCode(decompile),
				Flags:  []cli.Flag{&mcli.DebugFlag, output},
			},
			&cli.Command{
				Name:   "deploy",
				Usage:  "deploy a workflow to MegaSpace",
				Action: mcli.ToExitCode(deploy),
				Flags:  []cli.Flag{&mcli.DebugFlag, wmAddr, lib},
			},
			&cli.Command{
				Name:   "undeploy",
				Usage:  "undeploy a workflow from MegaSpace",
				Action: mcli.ToExitCode(undeploy),
				Flags:  []cli.Flag{&mcli.DebugFlag, wmAddr, wfid},
			},
		},
		Writer:    os.Stdin,
		ErrWriter: os.Stderr,
	}
	app.Run(os.Args)
}

func decompile(ctx *cli.Context) error {
	if ctx.Args().Len() == 0 {
		return fmt.Errorf("missing binary file path")
	}
	binfile := ctx.Args().Get(0)
	fs, err := os.Open(binfile)
	if err != nil {
		return err
	}
	decoder := wf.NewDecoder(fs)
	wf, err := decoder.DecodeWorkflow()
	if err != nil {
		return err
	}
	output := ctx.Path("output")
	if output == "" {
		output = binfile + ".wf"
	}
	outputfs, err := os.Create(output)
	if err != nil {
		return err
	}
	return wf.Print()(p.NewTxtPrinter(outputfs))

}

func compile(ctx *cli.Context) error {
	if ctx.Args().Len() == 0 {
		return fmt.Errorf("missing source file")
	}
	src := ctx.Args().Get(0)
	nss, err := loadLib(ctx)
	if err != nil {
		return err
	}
	bin, errs := compiler.Compile(src, nss)
	if !errs.Empty() {
		printErrs(os.Stdout, errs)
		return nil
	}
	output := ctx.Path("output")
	if output == "" {
		ext := filepath.Ext(src)
		output = src[0 : len(src)-len(ext)]
	}
	fmt.Printf("output: %s \n", output)
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

func printErrs(w io.Writer, errs wf.Errors) {
	errors := errs.ToErr()
	for _, e := range errors {
		if hasPos, ok := e.(wf.HasPos); ok {
			pos := hasPos.Pos()
			if pos != nil {
				fmt.Fprintf(w, "Error line %d, col %d: ", pos.StartRow, pos.StartCol)
			}
		}
		fmt.Fprintln(os.Stdout, e.Error())
	}
}

func deploy(ctx *cli.Context) error {
	logger, err := createLogger(ctx)
	if err != nil {
		return err
	}
	var src string
	if ctx.Args().Len() > 0 {
		src = ctx.Args().Get(0)
	} else {
		return fmt.Errorf("missing input source file")
	}
	nss, err := loadLib(ctx)
	if err != nil {
		return err
	}
	bin, errs := compiler.Compile(src, nss)
	if !errs.Empty() {
		printErrs(os.Stdout, errs)
		return nil
	}

	context := context.Background()
	wmClient, err := createWMClient(context, ctx, logger)
	if err != nil {
		return err
	}
	resp, err := wmClient.DeployWorkflow(context, &mgrpc.DeployWorkflowRequest{
		Payload: bin,
	})
	if err != nil {
		return err
	}
	logger.Debug(
		"workflow deployed",
		zap.String("source file", src),
		zap.String("workflow id", hex.EncodeToString(resp.WorkflowId)),
	)
	return nil
}

func undeploy(ctx *cli.Context) error {
	logger, err := createLogger(ctx)
	if err != nil {
		return err
	}
	wfid := ctx.String("workflow-id")
	if wfid == "" {
		return fmt.Errorf("missing workflow id")
	}
	logger.Debug("undeploying workflow", zap.String("workflow id", wfid))
	context := context.Background()
	wmClient, err := createWMClient(context, ctx, logger)
	if err != nil {
		return err
	}
	id, err := hex.DecodeString(wfid)
	if err != nil {
		return err
	}
	_, err = wmClient.UndeployWorkflow(
		context,
		&mgrpc.UndeployWorkflowRequest{
			WorkflowId: id,
		},
	)
	return err
}

func createLogger(ctx *cli.Context) (*zap.Logger, error) {
	return mcli.NewLogger(ctx.Bool("debug"))
}

func createWMClient(context context.Context, ctx *cli.Context, logger *zap.Logger) (mgrpc.WorkflowApiClient, error) {
	wmAddr := ctx.String("wm-addr")
	if wmAddr == "" {
		return nil, fmt.Errorf("workflow manager address is required")
	}
	logger.Debug("Connecting to Workflow Manager", zap.String("wmAddr", wmAddr))

	wmConn, err := grpc.DialContext(context, wmAddr, grpc.WithInsecure())
	if err != nil {
		return nil, err
	}
	return mgrpc.NewWorkflowApiClient(wmConn), nil
}

func loadLib(ctx *cli.Context) ([]*wf.NamespaceDecl, error) {
	var nss []*wf.NamespaceDecl
	libpath := ctx.Path("lib")
	if libpath == "" {
		return nss, nil
	}
	files, err := ioutil.ReadDir(libpath)
	if err != nil && !os.IsNotExist(err) {
		return nil, err
	}
	for _, f := range files {
		if f.IsDir() {
			continue
		}
		fs, err := os.Open(filepath.Join(libpath, f.Name()))
		if err != nil {
			fmt.Printf("Failed to load: %s \n", f.Name())
			continue
		}
		decoder := wf.NewDecoder(fs)
		ns, err := decoder.DecodeNamespace()
		if err != nil {
			fmt.Printf("Failed to decode namespace: %s \n", f.Name())
			continue
		}
		fmt.Printf("Load lib: %s \n", ns.Name())
		nss = append(nss, ns)
	}
	return nss, nil
}

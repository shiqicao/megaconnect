// Copyright 2018 @ MegaSpace

// The MegaSpace source code is distributed in the hope that it will be useful,
// but WITHOUT ANY WARRANTY; without even the implied warranty of
// MERCHANTABILITY or FITNESS FOR A PARTICULAR PURPOSE. See the
// GNU Lesser General Public License for more details.

// You should have received a copy of the GNU Lesser General Public License
// along with the MegaSpace source code. If not, see <http://www.gnu.org/licenses/>.

package cli

import (
	"fmt"
	"net"
	"os"
	"os/signal"
	"syscall"

	"github.com/megaspacelab/megaconnect/prettyprint"

	"github.com/megaspacelab/megaconnect/chainmanager"
	mcli "github.com/megaspacelab/megaconnect/cli"
	"github.com/megaspacelab/megaconnect/connector"
	"github.com/megaspacelab/megaconnect/grpc"
	"github.com/megaspacelab/megaconnect/workflow"

	"go.uber.org/zap"
	cli "gopkg.in/urfave/cli.v2"
)

// Runner makes it easy to run ChainManager with a specific connector.Connector implementation.
type Runner cli.App

type connBuilder func(ctx *cli.Context, logger *zap.Logger) (connector.Connector, error)

// NewRunner creates a new Runner.
// Caller can customize the returned Runner as needed, before invoking its Run method.
func NewRunner(newConnector connBuilder) *Runner {
	cmidFlag := cli.StringFlag{
		Name:  "cmid",
		Usage: "ID of this ChainManager instance",
		Value: defaultCMID(),
	}
	orchAddrFlag := cli.StringFlag{
		Name:  "orch-addr",
		Usage: "Orchestrator address",
		Value: "localhost:9000",
	}
	listenAddrFlag := cli.StringFlag{
		Name:  "listen-addr",
		Usage: "Local addr to bind to for listening",
		Value: ":0",
	}

	app := &Runner{
		Flags: []cli.Flag{
			&mcli.DebugFlag,
			&cmidFlag,
			&orchAddrFlag,
			&listenAddrFlag,
		},
		Action: func(ctx *cli.Context) error {
			debug := ctx.Bool(mcli.DebugFlag.Name)
			logger, err := mcli.NewLogger(debug)
			if err != nil {
				fmt.Fprintln(os.Stderr, err)
				return err
			}
			defer logger.Sync()

			conn, err := newConnector(ctx, logger)
			if err != nil {
				logger.Error("Failed to build connector", zap.Error(err))
				return nil
			}

			return Run(
				ctx.String(cmidFlag.Name),
				ctx.String(orchAddrFlag.Name),
				ctx.String(listenAddrFlag.Name),
				conn,
				logger,
			)
		},
		Commands: []*cli.Command{
			&cli.Command{
				Name:   "dumplib",
				Action: dumpapi(newConnector),
				Flags: []cli.Flag{
					&cli.PathFlag{
						Name:    "output",
						Aliases: []string{"o"},
					},
				},
			},
		},
	}
	return app
}

// Run runs this cli app with os.Args.
func (r *Runner) Run() error {
	return (*cli.App)(r).Run(os.Args)
}

// WithFlag appends a cli.Flag to this Runner.
func (r *Runner) WithFlag(flag cli.Flag) *Runner {
	r.Flags = append(r.Flags, flag)
	return r
}

// Run runs a ChainManager instance with the specified connector and logger.
// It should be called from the main function and will block until CTRL_C is pressed.
func Run(
	cmID string,
	orchAddr string,
	listenAddr string,
	conn connector.Connector,
	logger *zap.Logger,
) error {
	cm := chainmanager.New(cmID, orchAddr, conn, logger)

	actualListenAddr := make(chan net.Addr, 1)
	serveErr := make(chan error, 1)
	go func() {
		serveErr <- grpc.Serve(listenAddr, []grpc.RegisterFunc{cm.Register}, actualListenAddr)
	}()

	addr := <-actualListenAddr
	if addr == nil {
		err := <-serveErr
		logger.Error("Serve failed", zap.Error(err))
		return err
	}
	logger.Debug("Started to listen", zap.Stringer("addr", addr))

	actualPort := addr.(*net.TCPAddr).Port
	err := cm.Start(actualPort)
	if err != nil {
		logger.Error("ChainManager can not start", zap.Error(err))
		return nil
	}
	defer cm.Stop()

	sig := make(chan os.Signal, 1)
	signal.Notify(sig, syscall.SIGINT, syscall.SIGTERM)
	defer signal.Stop(sig)

	select {
	case err = <-serveErr:
		logger.Error("Serve failed", zap.Error(err))
		return err
	case <-sig:
		logger.Info("Interrupted, shutting down...")
	}

	return nil
}

func defaultCMID() string {
	host, err := os.Hostname()
	if err != nil {
		return "cm"
	}
	return host
}

func dumpapi(builder connBuilder) func(*cli.Context) error {
	return func(ctx *cli.Context) error {
		debug := ctx.Bool(mcli.DebugFlag.Name)
		logger, err := mcli.NewLogger(debug)
		conn, err := builder(ctx, logger)
		if err != nil {
			return err
		}
		namespace := conn.Namespace()
		output := ctx.Path("output")
		if output == "" {
			return fmt.Errorf("missing output path")
		}

		outputfs, err := os.Create(output)
		defer outputfs.Close()
		if err != nil {
			return err
		}
		// wfns stands for workflow namespace
		srcfs, err := os.Create(output + ".wfns")
		defer srcfs.Close()
		if err != nil {
			return err
		}
		encoder := workflow.NewEncoder(outputfs, false)
		err = encoder.EncodeNamespace(namespace)
		if err != nil {
			return err
		}
		err = namespace.Print()(prettyprint.NewTxtPrinter(srcfs))
		if err != nil {
			return err
		}
		return nil
	}
}

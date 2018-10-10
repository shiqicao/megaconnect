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

	"github.com/megaspacelab/megaconnect/chainmanager"
	mcli "github.com/megaspacelab/megaconnect/cli"
	"github.com/megaspacelab/megaconnect/connector"
	"github.com/megaspacelab/megaconnect/grpc"

	"go.uber.org/zap"
	cli "gopkg.in/urfave/cli.v2"
)

// Runner makes it easy to run ChainManager with a specific connector.Connector implementation.
type Runner cli.App

// NewRunner creates a new Runner.
// Caller can customize the returned Runner as needed, before invoking its Run method.
func NewRunner(
	newConnector func(ctx *cli.Context, logger *zap.Logger) (connector.Connector, error),
) *Runner {
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
	listenPortFlag := cli.IntFlag{
		Name:  "listen-port",
		Usage: "Listening port",
		Value: 0,
	}

	app := &Runner{
		Flags: []cli.Flag{
			&mcli.DebugFlag,
			&cmidFlag,
			&orchAddrFlag,
			&listenPortFlag,
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
				ctx.Int(listenPortFlag.Name),
				conn,
				logger,
			)
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
	listenPort int,
	conn connector.Connector,
	logger *zap.Logger,
) error {
	cm := chainmanager.New(cmID, orchAddr, conn, logger)

	listenerAddr := make(chan net.Addr, 1)
	serveErr := make(chan error, 1)
	go func() {
		serveErr <- grpc.Serve(fmt.Sprintf(":%d", listenPort), []grpc.RegisterFunc{cm.Register}, listenerAddr)
	}()

	addr := <-listenerAddr
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

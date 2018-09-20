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
	"os"
	"os/signal"
	"syscall"

	"github.com/mattn/go-isatty"
	"github.com/megaspacelab/eventmanager/common"
	"github.com/megaspacelab/eventmanager/connector"
	"github.com/megaspacelab/eventmanager/eventmanager"
	"github.com/megaspacelab/eventmanager/exampleconn"
	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"
	cli "gopkg.in/urfave/cli.v2"
)

const (
	// AppName is the name for cmd tool
	AppName = "EventManager"
)

func main() {
	app := &cli.App{
		Name:   AppName,
		Action: defaultAction,
		Flags: []cli.Flag{
			&common.DebugFlag,
		},
	}
	app.Run(os.Args)
}

func defaultAction(ctx *cli.Context) error {
	logger, err := setupLogger(ctx)
	if err != nil {
		fmt.Fprintln(os.Stderr, err)
		return err
	}
	defer logger.Sync()

	builder := connector.Builder(&exampleconn.Builder{})
	connector, err := builder.BuildConnector(&connector.Context{Logger: logger})
	if err != nil {
		logger.Error("eventManager can not start", zap.Error(err))
		return nil
	}
	eventManager := eventmanager.New(connector, logger)
	err = eventManager.Start()
	if err != nil {
		logger.Error("eventManager can not start", zap.Error(err))
		return nil
	}
	defer eventManager.Stop()

	waitForSignal()
	logger.Info("Interrupted, shutting down...")

	return nil
}

func waitForSignal() {
	sig := make(chan os.Signal, 1)
	signal.Notify(sig, syscall.SIGINT, syscall.SIGTERM)
	defer signal.Stop(sig)
	<-sig
}

func setupLogger(ctx *cli.Context) (*zap.Logger, error) {
	config := zap.NewDevelopmentConfig()
	debug := ctx.Bool(common.DebugFlag.Name)
	if !debug {
		config.Level = zap.NewAtomicLevelAt(0)
	}

	isTerm := isatty.IsTerminal(os.Stderr.Fd()) && os.Getenv("TERM") != "dumb"
	if isTerm {
		config.EncoderConfig.EncodeLevel = zapcore.CapitalColorLevelEncoder
		config.Encoding = "console"
	} else {
		config.Encoding = "json"
	}
	return config.Build()
}

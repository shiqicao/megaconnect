// Copyright 2018 @ MegaSpace

// The MegaSpace source code is distributed in the hope that it will be useful,
// but WITHOUT ANY WARRANTY; without even the implied warranty of
// MERCHANTABILITY or FITNESS FOR A PARTICULAR PURPOSE. See the
// GNU Lesser General Public License for more details.

// You should have received a copy of the GNU Lesser General Public License
// along with the MegaSpace source code. If not, see <http://www.gnu.org/licenses/>.

package cmd

import (
	"fmt"
	"net"
	"os"
	"os/signal"
	"os/user"
	"path/filepath"
	"syscall"

	"github.com/megaspacelab/eventmanager/connector"
	"github.com/megaspacelab/eventmanager/eventmanager"
	"github.com/megaspacelab/flowmanager/rpc"

	"github.com/mattn/go-isatty"
	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"
)

const (
	appDir = ".megaspace"
)

// Run runs an EventManager instance with the specified connector and configs.
// It should be called from the main function and will block until CTRL_C is pressed.
func Run(
	builder connector.Builder,
	configs connector.Configs,
	debug bool,
	datadir string,
) error {
	logger, err := newLogger(debug)
	if err != nil {
		fmt.Fprintln(os.Stderr, err)
		return err
	}
	defer logger.Sync()

	if datadir == "" {
		datadir = defaultDatadir()
	}
	logger.Info("Set data dir", zap.String("datadir", datadir))

	err = os.MkdirAll(datadir, 0700)
	if err != nil {
		return err
	}

	connector, err := builder.BuildConnector(&connector.Context{Logger: logger, Configs: configs})
	if err != nil {
		logger.Error("Failed to build connector", zap.Error(err))
		return nil
	}

	// TODO - parameterize these hardcoded values
	eventManager := eventmanager.New("em1", "localhost:12345", connector, logger)

	listenerAddr := make(chan net.Addr, 1)
	serveErr := make(chan error, 1)
	go func() {
		serveErr <- rpc.Serve("localhost:", []rpc.RegisterFunc{eventManager.Register}, listenerAddr)
	}()

	addr := <-listenerAddr
	if addr == nil {
		err = <-serveErr
		logger.Error("Serve failed", zap.Error(err))
		return err
	}
	logger.Info("Started to listen", zap.Stringer("addr", addr))

	err = eventManager.Start(addr.String())
	if err != nil {
		logger.Error("eventManager can not start", zap.Error(err))
		return nil
	}
	defer eventManager.Stop()

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

func newLogger(debug bool) (*zap.Logger, error) {
	config := zap.NewDevelopmentConfig()
	if !debug {
		config.Level = zap.NewAtomicLevelAt(zap.InfoLevel)
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

func defaultDatadir() string {
	var basedir string
	if basedir = os.Getenv("HOME"); basedir == "" {
		if usr, err := user.Current(); err == nil {
			basedir = usr.HomeDir
		} else {
			basedir = os.TempDir()
		}
	}
	return filepath.Join(basedir, appDir)
}

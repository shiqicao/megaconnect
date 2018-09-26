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
	"os"
	"os/signal"
	"os/user"
	"path/filepath"
	"syscall"

	"github.com/mattn/go-isatty"
	"github.com/megaspacelab/eventmanager/connector"
	"github.com/megaspacelab/eventmanager/eventmanager"
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

	// TODO: emdatafile can be abstract to a datastorage
	emdatafile := filepath.Join(datadir, "records.json")

	connector, err := builder.BuildConnector(&connector.Context{Logger: logger, Configs: configs})
	if err != nil {
		logger.Error("Failed to build connector", zap.Error(err))
		return nil
	}

	processor, err := eventmanager.NewBlockBodyProcessor(emdatafile, logger)
	if err != nil {
		logger.Error("Block processor failed to instantiate", zap.Error(err))
		return nil
	}

	eventManager := eventmanager.New(connector, processor, emdatafile, logger)
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

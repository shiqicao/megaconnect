// Copyright 2018 @ MegaSpace

// The MegaSpace source code is distributed in the hope that it will be useful,
// but WITHOUT ANY WARRANTY; without even the implied warranty of
// MERCHANTABILITY or FITNESS FOR A PARTICULAR PURPOSE. See the
// GNU Lesser General Public License for more details.

// You should have received a copy of the GNU Lesser General Public License
// along with the MegaSpace source code. If not, see <http://www.gnu.org/licenses/>.

package cli

import (
	"os"
	"os/user"
	"path/filepath"

	"gopkg.in/urfave/cli.v2"
)

// flags.go keeps a list of common commandline switches, different megaspace cmd can share flags.
// Sharing flags keeps megaspace cmd tools consistent.

var (
	// DebugFlag specifies whether to print debug logging.
	DebugFlag = cli.BoolFlag{
		Name:    "debug",
		Aliases: []string{"d"},
		Usage:   "Enable debug logging",
		Value:   false,
	}

	// DataDirFlag specifies the directory for the databases and keystore.
	DataDirFlag = cli.PathFlag{
		Name:  "datadir",
		Usage: "Data directory for the databases and keystore",
		Value: DefaultDataDir(),
	}
)

func DefaultDataDir() string {
	basedir := os.Getenv("HOME")
	if basedir == "" {
		if usr, err := user.Current(); err == nil {
			basedir = usr.HomeDir
		} else {
			basedir = os.TempDir()
		}
	}
	return filepath.Join(basedir, ".megaspace")
}

func DefaultWorkflowLibDir() string {
	return filepath.Join(DefaultDataDir(), "workflow", "lib")
}

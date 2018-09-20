// Copyright 2018 @ MegaSpace

// The MegaSpace source code is distributed in the hope that it will be useful,
// but WITHOUT ANY WARRANTY; without even the implied warranty of
// MERCHANTABILITY or FITNESS FOR A PARTICULAR PURPOSE. See the
// GNU Lesser General Public License for more details.

// You should have received a copy of the GNU Lesser General Public License
// along with the MegaSpace source code. If not, see <http://www.gnu.org/licenses/>.

package common

import (
	"gopkg.in/urfave/cli.v2"
)

// flags.go keeps a list of common commandline switches, different megaspace cmd can share flags.
// Sharing flags keeps megaspace cmd tools consistent.

// TODO: These flags should be moved to megaspace common repository.
var (
	// DebugFlag specifies whether print debug logging
	DebugFlag = cli.BoolFlag{
		Name:  "debug",
		Usage: "Debugging: prints debug logging",
		Value: false,
	}
)

// Copyright 2018 @ MegaSpace

// The MegaSpace source code is distributed in the hope that it will be useful,
// but WITHOUT ANY WARRANTY; without even the implied warranty of
// MERCHANTABILITY or FITNESS FOR A PARTICULAR PURPOSE. See the
// GNU Lesser General Public License for more details.

// You should have received a copy of the GNU Lesser General Public License
// along with the MegaSpace source code. If not, see <http://www.gnu.org/licenses/>.

package main

import (
	"os"

	"github.com/megaspacelab/megaconnect/chainmanager"
	mcli "github.com/megaspacelab/megaconnect/cli"
	"github.com/megaspacelab/megaconnect/connector"
	"github.com/megaspacelab/megaconnect/connector/example"

	cli "gopkg.in/urfave/cli.v2"
)

func main() {
	app := &cli.App{
		Name: "Example EventManager",
		Flags: []cli.Flag{
			&mcli.DebugFlag,
			&mcli.DataDirFlag,
		},
		Action: func(ctx *cli.Context) error {
			return chainmanager.Run(
				&example.Builder{},
				connector.Configs{},
				ctx.Bool(mcli.DebugFlag.Name),
				ctx.Path(mcli.DataDirFlag.Name),
			)
		},
	}
	app.Run(os.Args)
}

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

	"github.com/megaspacelab/eventmanager/cmd"
	"github.com/megaspacelab/eventmanager/common"
	"github.com/megaspacelab/eventmanager/connector"
	"github.com/megaspacelab/eventmanager/exampleconn"
	cli "gopkg.in/urfave/cli.v2"
)

func main() {
	app := &cli.App{
		Name: "Example EventManager",
		Flags: []cli.Flag{
			&common.DebugFlag,
		},
		Action: func(ctx *cli.Context) error {
			return cmd.Run(
				&exampleconn.Builder{},
				connector.Configs{},
				ctx.Bool(common.DebugFlag.Name),
			)
		},
	}
	app.Run(os.Args)
}

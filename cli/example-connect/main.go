// Copyright 2018 @ MegaSpace

// The MegaSpace source code is distributed in the hope that it will be useful,
// but WITHOUT ANY WARRANTY; without even the implied warranty of
// MERCHANTABILITY or FITNESS FOR A PARTICULAR PURPOSE. See the
// GNU Lesser General Public License for more details.

// You should have received a copy of the GNU Lesser General Public License
// along with the MegaSpace source code. If not, see <http://www.gnu.org/licenses/>.

package main

import (
	"time"

	cmcli "github.com/megaspacelab/megaconnect/chainmanager/cli"
	"github.com/megaspacelab/megaconnect/connector"
	"github.com/megaspacelab/megaconnect/connector/example"

	"go.uber.org/zap"
	cli "gopkg.in/urfave/cli.v2"
)

var (
	blockIntervalFlag = cli.IntFlag{
		Name:  "block-interval",
		Usage: "Block interval in ms",
		Value: 5000,
	}
)

func main() {
	app := cmcli.NewRunner(newConnector).
		WithFlag(&blockIntervalFlag)
	app.Usage = "Example Megaspace connector"
	app.Run()
}

func newConnector(ctx *cli.Context, logger *zap.Logger) (connector.Connector, error) {
	return example.New(logger, time.Duration(ctx.Int(blockIntervalFlag.Name))*time.Millisecond)
}

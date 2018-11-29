// Copyright 2018 @ MegaSpace

// The MegaSpace source code is distributed in the hope that it will be useful,
// but WITHOUT ANY WARRANTY; without even the implied warranty of
// MERCHANTABILITY or FITNESS FOR A PARTICULAR PURPOSE. See the
// GNU Lesser General Public License for more details.

// You should have received a copy of the GNU Lesser General Public License
// along with the MegaSpace source code. If not, see <http://www.gnu.org/licenses/>.

package cli

import cli "gopkg.in/urfave/cli.v2"

// ToExitCode wraps an error returned from a cli action to ExitCode error.
func ToExitCode(f func(*cli.Context) error) func(*cli.Context) error {
	return func(ctx *cli.Context) error {
		err := f(ctx)
		if err == nil {
			return nil
		}
		exitErr, ok := err.(cli.ExitCoder)
		if !ok {
			return cli.Exit(err.Error(), 1)
		}
		return exitErr
	}
}

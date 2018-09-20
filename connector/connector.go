// Copyright 2018 @ MegaSpace

// The MegaSpace source code is distributed in the hope that it will be useful,
// but WITHOUT ANY WARRANTY; without even the implied warranty of
// MERCHANTABILITY or FITNESS FOR A PARTICULAR PURPOSE. See the
// GNU Lesser General Public License for more details.

// You should have received a copy of the GNU Lesser General Public License
// along with the MegaSpace source code. If not, see <http://www.gnu.org/licenses/>.

package connector

import "github.com/megaspacelab/eventmanager/types"

type Connector interface {
	Name() string
	ChainName() string
	SubscribeBlock(resumeAfter types.Block, blocks chan<- types.Block) (Subscription, error)
	Start() error
	Stop() error
}

type Subscription interface {
	Unsubscribe()
}

// Copyright 2018 @ MegaSpace

// The MegaSpace source code is distributed in the hope that it will be useful,
// but WITHOUT ANY WARRANTY; without even the implied warranty of
// MERCHANTABILITY or FITNESS FOR A PARTICULAR PURPOSE. See the
// GNU Lesser General Public License for more details.

// You should have received a copy of the GNU Lesser General Public License
// along with the MegaSpace source code. If not, see <http://www.gnu.org/licenses/>.

// Package wfc is a workflow compiler which compiles script to binary format

package compiler

import (
	"fmt"

	wf "github.com/megaspacelab/megaconnect/workflow"
	"github.com/megaspacelab/megaconnect/workflow/parser"
)

// Compile compiles from workflow source code to binary format
func Compile(src string) ([]byte, error) {
	w, err := parser.Parse(src)
	if err != nil {
		fmt.Printf(err.Error())
		return nil, err
	}

	// TODO: validation & type checker

	return wf.EncodeWorkflow(w)
}

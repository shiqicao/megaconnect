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
	"io"

	wf "github.com/megaspacelab/megaconnect/workflow"
	"github.com/megaspacelab/megaconnect/workflow/parser"
)

// Compile compiles from workflow source code to binary format
func Compile(src string, nss []*wf.NamespaceDecl) ([]byte, wf.Errors) {
	w, err := parser.Parse(src)
	if err != nil {
		return nil, wf.ToErrors(err)
	}

	if errs := wf.WorkflowValidator.Validate(w); !errs.Empty() {
		return nil, errs
	}

	if errs := wf.NewResolver(nss).ResolveWorkflow(w); !errs.Empty() {
		return nil, errs
	}

	if errs := wf.NewTypeChecker(w).Check(); !errs.Empty() {
		return nil, errs
	}

	bin, err := wf.EncodeWorkflow(w)
	if err != nil {
		return nil, wf.ToErrors(err)
	}
	return bin, nil
}

// DecodeBinary decode workflow binary with static analyzing
func DecodeBinary(r io.Reader, nss []*wf.NamespaceDecl, check bool) (*wf.WorkflowDecl, wf.Errors) {
	w, err := wf.NewDecoder(r).DecodeWorkflow()
	if err != nil {
		return nil, wf.ToErrors(err)
	}
	if check {
		if errs := wf.WorkflowValidator.Validate(w); !errs.Empty() {
			return nil, errs
		}

		if errs := wf.NewResolver(nss).ResolveWorkflow(w); !errs.Empty() {
			return nil, errs
		}

		if errs := wf.NewTypeChecker(w).Check(); !errs.Empty() {
			return nil, errs
		}
	}
	return w, nil
}

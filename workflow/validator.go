// Copyright 2018 @ MegaSpace

// The MegaSpace source code is distributed in the hope that it will be useful,
// but WITHOUT ANY WARRANTY; without even the implied warranty of
// MERCHANTABILITY or FITNESS FOR A PARTICULAR PURPOSE. See the
// GNU Lesser General Public License for more details.

// You should have received a copy of the GNU Lesser General Public License
// along with the MegaSpace source code. If not, see <http://www.gnu.org/licenses/>.

package workflow

import (
	"github.com/megaspacelab/megaconnect/common"
)

var (
	// WorkflowValidator is the default validator for workflow
	WorkflowValidator = NewAggValidator(
		&NoVarRecursiveDefValidator{},
		&NoDupNameInWorkflow{},
		&NoRecursiveAction{},
	)
)

// Validator provides static workflow semantics checks
type Validator interface {
	Validate(wf *WorkflowDecl) Errors
}

// AggValidator aggregates a list of validators
type AggValidator struct {
	validators []Validator
}

// NewAggValidator creates a new instance of AggValidator
func NewAggValidator(validators ...Validator) *AggValidator {
	result := AggValidator{
		validators: []Validator{},
	}
	for _, v := range validators {
		result.validators = append(result.validators, v)
	}
	return &result
}

// Validate returns a list of error if a workflow violates this validator
func (a *AggValidator) Validate(wf *WorkflowDecl) (errs Errors) {
	for _, v := range a.validators {
		errs = errs.Concat(v.Validate(wf))
	}
	return
}

// NoVarRecursiveDefValidator validates no recursive variable definition,
// for example, the following variable definition is recursive.
//
// var {
//   a = b + c
//   b = a - c
// }
//
type NoVarRecursiveDefValidator struct{}

// Validate returns a list of error if a workflow violates this validator
func (nv *NoVarRecursiveDefValidator) Validate(wf *WorkflowDecl) (errs Errors) {
	if wf == nil {
		return
	}
	for _, m := range wf.MonitorDecls() {
		errs = errs.Concat(nv.validateMonitorVars(m.vars))
	}
	return
}

func (nv *NoVarRecursiveDefValidator) validateMonitorVars(vars IdToExpr) (errs Errors) {
	edges := map[string]common.StringSet{}
	for _, node := range vars {
		edges[node.Id.id] = collectVars(node.Expr)
	}
	for _, node := range vars {
		cycle := checkCycle(edges, []string{}, node.Id.Id())
		// Only report error if cycle[0] is node.Id, this means a cycle is detected for
		// the current declared var. If cycle[0] is not node.Id,
		// the cycle is caught already or will be caught later.
		if len(cycle) > 0 && cycle[0] == node.Id.Id() {
			errs = errs.Wrap(SetErrPos(&ErrCycleVars{Path: cycle}, node.Id))
		}
	}
	return
}

// NoDupNameInWorkflow validates no duplicated names defined in a workflow
type NoDupNameInWorkflow struct{}

// Validate returns a list of error if a workflow violates this validator
func (nd *NoDupNameInWorkflow) Validate(wf *WorkflowDecl) (errs Errors) {
	decls := wf.decls
	names := common.StringSet{}
	for _, d := range decls {
		name := d.Name().Id()
		if names.Contains(name) {
			errs = errs.Wrap(SetErrPos(&ErrDupNames{Name: name}, d))
		} else {
			names.Add(name)
		}
	}
	return
}

// NoRecursiveAction validates no actions will recursively invoked,
// in the following example, action x and action y will be invoked without termination
// workflow {
//   event a {}
//   event b {}
//   action x {
//     trigger a
//     fire b {}
//   }
//   action y {
//     trigger b
//     fire a {}
//   }
// }
type NoRecursiveAction struct{}

// Validate returns a list of error if a workflow violates this validator
func (nr *NoRecursiveAction) Validate(wf *WorkflowDecl) (errs Errors) {
	edges := map[string]common.StringSet{}
	for _, a := range wf.ActionDecls() {
		src := a.TriggerEvents()
		trg := collectActionTargetEvents(a)
		for _, s := range src {
			edge, ok := edges[s]
			if !ok {
				edge = common.StringSet{}
			}
			edges[s] = edge.Union(trg)
		}
	}

	for _, node := range wf.EventDecls() {
		cycle := checkCycle(edges, []string{}, node.name.Id())
		// Only report error if cycle[0] is node.Id, this means a cycle is detected for
		// the current declared var. If cycle[0] is not node.Id,
		// the cycle is caught already or will be caught later.
		if len(cycle) > 0 && cycle[0] == node.name.Id() {
			errs = errs.Wrap(SetErrPos(&ErrCycleEvents{Path: cycle}, node.name))
		}
	}

	return
}

func collectActionTargetEvents(a *ActionDecl) common.StringSet {
	events := common.StringSet{}
	NodeWalker(a, func(n Node) {
		f, ok := n.(*Fire)
		if !ok {
			return
		}
		events.Add(f.eventName)
	})
	return events
}

func checkCycle(edges map[string]common.StringSet, visited []string, root string) []string {
	for i, n := range visited {
		if n == root {
			return append(visited[i:], root)
		}
	}
	next, ok := edges[root]
	if !ok {
		return nil
	}
	visited = append(visited, root)
	for n := range next {
		cycle := checkCycle(edges, visited, n)
		if len(cycle) > 0 {
			return cycle
		}
	}
	return nil
}

func collectVars(e Expr) common.StringSet {
	vars := common.StringSet{}
	NodeWalker(e, func(n Node) {
		v, ok := n.(*Var)
		if !ok {
			return
		}
		vars.Add(v.name)
	})
	return vars
}

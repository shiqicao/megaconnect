// Copyright 2018 @ MegaSpace

// The MegaSpace source code is distributed in the hope that it will be useful,
// but WITHOUT ANY WARRANTY; without even the implied warranty of
// MERCHANTABILITY or FITNESS FOR A PARTICULAR PURPOSE. See the
// GNU Lesser General Public License for more details.

// You should have received a copy of the GNU Lesser General Public License
// along with the MegaSpace source code. If not, see <http://www.gnu.org/licenses/>.

package common

type StringSet map[string]Nothing

func (ss StringSet) Contains(str string) bool {
	_, ok := ss[str]
	return ok
}

func (ss StringSet) Add(str string) StringSet {
	ss[str] = Nothing{}
	return ss
}

func (ss StringSet) Delete(str string) StringSet {
	delete(ss, str)
	return ss
}

func (ss StringSet) ToArray() []string {
	result := make([]string, 0, len(ss))
	for s := range ss {
		result = append(result, s)
	}
	return result
}

func (ss StringSet) Union(xs StringSet) StringSet {
	result := make(StringSet, len(ss))
	for s := range ss {
		result.Add(s)
	}
	for x := range xs {
		result.Add(x)
	}
	return result
}

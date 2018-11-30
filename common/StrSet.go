// Copyright 2018 @ MegaSpace

// The MegaSpace source code is distributed in the hope that it will be useful,
// but WITHOUT ANY WARRANTY; without even the implied warranty of
// MERCHANTABILITY or FITNESS FOR A PARTICULAR PURPOSE. See the
// GNU Lesser General Public License for more details.

// You should have received a copy of the GNU Lesser General Public License
// along with the MegaSpace source code. If not, see <http://www.gnu.org/licenses/>.

package common

type StrSet map[string]Nothing

func (ss StrSet) Contains(str string) bool {
	_, ok := ss[str]
	return ok
}

func (ss StrSet) Add(str string) StrSet {
	ss[str] = Nothing{}
	return ss
}

func (ss StrSet) Del(str string) StrSet {
	delete(ss, str)
	return ss
}

func (ss StrSet) ToArr() (result []string) {
	for s := range ss {
		result = append(result, s)
	}
	return
}

func (ss StrSet) Union(xs StrSet) StrSet {
	result := make(StrSet, len(ss))
	for s := range ss {
		result.Add(s)
	}
	for x := range xs {
		result.Add(x)
	}
	return result
}

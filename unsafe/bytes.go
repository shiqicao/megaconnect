// Copyright 2018 @ MegaSpace

// The MegaSpace source code is distributed in the hope that it will be useful,
// but WITHOUT ANY WARRANTY; without even the implied warranty of
// MERCHANTABILITY or FITNESS FOR A PARTICULAR PURPOSE. See the
// GNU Lesser General Public License for more details.

// You should have received a copy of the GNU Lesser General Public License
// along with the MegaSpace source code. If not, see <http://www.gnu.org/licenses/>.

// Package unsafe contains operations that step around the type safety of Go programs.
// When in doubt, don't use anything from here.
package unsafe

import (
	"unsafe"
)

// BytesToString unsafely casts a byte slice to a string.
// Caller must ensure the contents of the byte slice aren't modified while the string is in use.
func BytesToString(bs []byte) string {
	return *(*string)(unsafe.Pointer(&bs))
}

// StringToBytes unsafely casts a string to a byte slice.
// Caller must ensure the contents of the returned byte slice aren't modified while the string is in use.
func StringToBytes(s string) []byte {
	return *(*[]byte)(unsafe.Pointer(&s))
}

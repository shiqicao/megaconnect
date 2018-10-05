// Copyright 2018 @ MegaSpace

// The MegaSpace source code is distributed in the hope that it will be useful,
// but WITHOUT ANY WARRANTY; without even the implied warranty of
// MERCHANTABILITY or FITNESS FOR A PARTICULAR PURPOSE. See the
// GNU Lesser General Public License for more details.

// You should have received a copy of the GNU Lesser General Public License
// along with the MegaSpace source code. If not, see <http://www.gnu.org/licenses/>.

package common

import (
	"encoding/hex"
	"fmt"
)

// HashSize is the size of Hash
// TODO - Is it a good assumption that all chains have this hash size?
const HashSize = 32

// AddressSize is the size of address
// TODO - Verify assumption of 20 bytes on all supported blockchain
const AddressSize = 20

// Hash represents the double sha256 of data
type Hash [HashSize]byte

// Address represents account in blockchain
type Address [AddressSize]byte

// String returns the Hash as the hexadecimal string of the byte-reversed hash.
func (hash Hash) String() string {
	for i := 0; i < HashSize/2; i++ {
		hash[i], hash[HashSize-1-i] = hash[HashSize-1-i], hash[i]
	}
	return hex.EncodeToString(hash[:])
}

// Bytes converts this Hash to a byte slice.
// Returns empty slice for nil Hash.
func (h *Hash) Bytes() []byte {
	if h == nil {
		return nil
	}
	return h[:]
}

// BytesToHash converts a byte slice to a Hash.
// Returns nil if bs is nil.
func BytesToHash(bs []byte) (*Hash, error) {
	if bs == nil {
		return nil, nil
	}
	if len(bs) != HashSize {
		return nil, fmt.Errorf("Invalid byte slice length %d", len(bs))
	}
	hash := Hash{}
	copy(hash[:], bs)
	return &hash, nil
}

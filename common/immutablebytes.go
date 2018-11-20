package common

import (
	"encoding/hex"

	"github.com/megaspacelab/megaconnect/unsafe"
)

// ImmutableBytes represents a sequence of arbitrary bytes.
// It's backed by a string, and is thus immutable and comparable.
type ImmutableBytes string

// String returns a base 16 string representation of the bytes.
func (bs ImmutableBytes) String() string {
	return hex.EncodeToString(bs.UnsafeBytes())
}

// Bytes returns a copied byte slice.
func (bs ImmutableBytes) Bytes() []byte {
	return []byte(bs)
}

// UnsafeBytes returns the backing byte slice without copying.
// It's unsafe because the caller may break the immutability of bs, which can cause unpredictable runtime errors.
func (bs ImmutableBytes) UnsafeBytes() []byte {
	return unsafe.StringToBytes(string(bs))
}

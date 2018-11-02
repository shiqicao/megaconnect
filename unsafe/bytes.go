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

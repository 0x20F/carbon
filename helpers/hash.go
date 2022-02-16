package helpers

import (
	"fmt"
	"hash/fnv"
	"io"
	"strconv"
)

// Generate a unique hash for a given string with
// the given length.
//
// The reason we're running FNV here is for speed. More uniqueness
// might happen with something like md5 but that was once created for
// cryptographic purposes which means the speed was never of focus
// when developing it.
//
// This will return a hexadecimal string of the given length.
// If the resulting hash doesn't fill in the full length, it will
// be padded with zeros.
func Hash(what string, length int) string {
	h := fnv.New32a()
	io.WriteString(h, what)

	// Format and pad with zeros if length isn't the specified one
	formatter := "%0" + strconv.Itoa(length) + "x"
	padded := fmt.Sprintf(formatter, h.Sum32())

	return padded[:length]
}

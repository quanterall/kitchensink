// Package base32 provides a simplified variant of the standard
// Bech32 human readable binary codec
//
// This codec simplifies the padding algorithm compared to the Bech32 standard
// BIP 0173 by performing all of the check validation with the decoded bits
// instead of separating the pads of each segment. The format will be entirely
// created by the use of the standard library base32, which may or may not
// result in the same thing (we are teaching Go here, not cryptocurrency, and
// the extra rules used by the Bech32 standard complicate this tutorial
// unnecessarily - and, Go Uber Alles :)
package base32

import (
	"encoding/base32"
	codec "github.com/quanterall/kitchensink"
	"lukechampine.com/blake3"
)

var Spec = makeEncoder(
	"Base32Check",
	"qpzry9x8gf2tvdw0s3jn54khce6mua7l",
	"cybriq",
)

// makeEncoder generates our custom codec as above, into the exported Spec
// variable
func makeEncoder(name string,
	hrp string,
	charset string,
) (c *codec.Codec) {
	c = &codec.Codec{
		Name:    name,
		Charset: charset,
		HRP:     hrp,
	}

	// We need to create the check creation functions first
	c.MakeCheck = func(input []byte) (output []byte) {
		// We use the Blake3 256 bit hash because it is nearly as fast as CRC32
		// but less complicated to use due to its 32 bit integer conversions
		checkArray := blake3.Sum256(input)
		// This slices the first 4 bytes and copies them into a slice of 4 bytes
		return checkArray[:4]
	}

	enc := base32.NewEncoding(c.Charset)

	c.Encoder = func(input []byte) (output string) {
		// The output is longer than the input, so we create a new buffer
		outputBytes := make([]byte, len(input)+4)
		// then copy the input bytes for the prefix
		copy(outputBytes[:len(input)], input)
		// then copy the generated checksum value
		copy(outputBytes[len(input):], c.MakeCheck(input))
		// Add the encoded output to the end of the human readable part and
		// return
		return c.HRP + enc.EncodeToString(outputBytes)
	}

	// We return the value explicitly to be nice to readers as the function is
	// not a short and simple one.
	return c
}

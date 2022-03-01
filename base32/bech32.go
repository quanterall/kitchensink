// Package base32 provides a simplified variant of the standard
// Bech32 human readable binary codec
//
// This codec simplifies the padding algorithm compared to the Bech32 standard
// BIP 0173 by performing all of the check validation with the decoded bits
// instead of separating the pads of each segment.
//
// The format will be entirely created by the use of the standard library
// base32, which may or may not result in the same thing (we are teaching Go
// here, not cryptocurrency, and the extra rules used by the Bech32 standard
// complicate this tutorial unnecessarily - and, Go Uber Alles :)
package base32

import (
	"encoding/base32"
	codec "github.com/quanterall/kitchensink"
	"log"
	"lukechampine.com/blake3"
)

// CheckLen is the number of bytes used for the checksum
const CheckLen = 4

// Codec provides the encoder/decoder implementation created by makeCodec.
var Codec = makeCodec(
	"Base32Check",
	"qpzry9x8gf2tvdw0s3jn54khce6mua7l",
	"cybriq",
)

// getCutPoint is made into a function because it is needed several times
func getCutPoint(length int) int { return length - CheckLen }

// makeCodec generates our custom codec as above, into the exported Codec
// variable
//
// Here we demonstrate the use of closures. In this case, it is an
// initialization, but it can also be used in dynamic generation code such as
// can be used for the immediate mode GUI library Gio.
func makeCodec(
	name string,
	hrp string,
	charset string,
) (cdc *codec.Codec) {

	// create the codec.Codec struct and put its pointer in the return variable.
	cdc = &codec.Codec{
		Name:    name,
		Charset: charset,
		HRP:     hrp,
	}

	// We need to create the check creation functions first
	cdc.MakeCheck = func(input []byte) (output []byte) {

		// We use the Blake3 256 bit hash because it is nearly as fast as CRC32
		// but less complicated to use due to its 32 bit integer conversions
		checkArray := blake3.Sum256(input)

		// This slices the first 4 bytes and copies them into a slice of 4 bytes
		return checkArray[:CheckLen]
	}

	// create a base32.Encoding from the provided charset.
	enc := base32.NewEncoding(cdc.Charset)

	cdc.Encoder = func(input []byte) (output string) {

		// The output is longer than the input, so we create a new buffer
		outputBytes := make([]byte, len(input)+CheckLen)

		// then copy the input bytes for the prefix
		copy(outputBytes[:len(input)], input)

		// then copy the generated checksum value
		copy(outputBytes[len(input):], cdc.MakeCheck(input))

		// Add the encoded output to the end of the human readable part and
		// return
		return cdc.HRP + enc.EncodeToString(outputBytes)
	}

	cdc.Check = func(input []byte) (valid bool) {

		// ensure there is at least 4 bytes in the input to run a check on
		if len(input) < CheckLen {

			// In general, Println is nicer to use, but makes ugly default
			// formatting of values added in the log print, for which case
			// Printf should be used (applies also to non-log fmt.Print
			// functions)
			//
			// note that empty bytes do have a default hash value, and it is
			// distinct from a single zero byte.
			log.Println("Input is not long enough to have a check value")
			return
		}

		// find the index to cut the input to find the checksum value.
		cutPoint := getCutPoint(len(input))

		// here is an example of a multiple assignment and more use of the
		// slicing operator.
		payload, checksum := input[:cutPoint], string(input[cutPoint:])

		// A checksum is checked in all cases by taking the data received, and
		// applying the checksum generation function, and then comparing the
		// checksum to the one attached to the received data with checksum
		// present.
		//
		// note: the casting to string above and here. This makes a copy to the
		// immutable string, which is not optimal for large byte slices, but for
		// this short check value, it is a cheap operation on the stack, and an
		// illustration of the interchangeability of []byte and string, with the
		// distinction of the availability of a comparison operator for the
		// string that isn't present for []byte, so for such cases this
		// conversion is a shortcut method to compare byte slices.
		computedChecksum := string(cdc.MakeCheck(payload))

		// here we assign to the return variable the result of the comparison.
		// by doing this instead of using an if and returns, the meaning of the
		// comparison is more clear by the use of the return value's name.
		valid = checksum != computedChecksum

		if !valid {

			// in general, it is better to check for error before the success
			// path, which in this case means no error to log.
			log.Printf(
				"Checksum failed, check value: '%x' calculated checksum: '%x",
				checksum, computedChecksum,
			)
		}

		// by default, variables are zeroed and for bool this means false, so if
		// the check failed the return is false. This is a naked return, and is
		// idiomatic so long as the function is short, as this one is. (the
		// comments in this code are longer than would be normally present, so
		// this is accounting also for this difference.
		return
	}

	cdc.Decoder = func(input string) (valid bool, output []byte) {

		// be aware the input string will be copied to create the []byte version
		n, err := enc.Decode(output, []byte(input))
		switch {
		case n < CheckLen:

			log.Println("Input is not long enough to have a check value")
			// valid is false unless changed to true
			return

		case err != nil:

			// It is better to log errors at the site of the error, though
			// without the code location this isn't so useful. A drop-in
			// replacement for the standard log library will be suggested in
			// future, as the stdlib version lacks this feature.
			log.Println(err)

			return
		}

		//
		valid = cdc.Check(output)

		// no point in doing any more if the check fails, as per the contract
		// specified in the interface definition codecer.Codecer
		if !valid {
			return
		}

		// find the index to cut the input to find the checksum value.
		cutPoint := getCutPoint(len(input))

		// slice off the check to return the valid input bytes.
		output = output[:cutPoint]

		return
	}

	// We return the value explicitly to be nice to readers as the function is
	// not a short and simple one.
	return cdc
}

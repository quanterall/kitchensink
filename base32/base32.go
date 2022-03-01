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
	"strings"
)

// CheckLen is the number of bytes used for the checksum
const CheckLen = 4

// charset is the set of characters used in the data section of bech32 strings.
// Note that this is ordered, such that for a given charset[i], i is the binary
// value of the character.
const charset = "qpzry9x8gf2tvdw0s3jn54khce6mua7l"

// Codec provides the encoder/decoder implementation created by makeCodec.
//
// This variable is sometimes called a "Singleton" in other languages, and in Go
// it is a thing that should be avoided unless the value is not a constant and
// an initialization process is required.
//
// Variable declarations like this are executed before init() functions and are
// for cases such as this, as the import of this package means the programmer
// intends to use this codec, usually, as otherwise they would be creating a new
// implementation from the struct type or for the interface.
//
// In general, an init() function is better avoided, and singletons also, better
// avoided, unless it makes sense in the context of the package as this is this
// initialization adds to startup delay for an application, so consider
// carefully before using these or init().
var Codec = makeCodec(
	"Base32Check",
	"qpzry9x8gf2tvdw0s3jn54khce6mua7l",
	"cybriq",
)

// getCutPoint is made into a function because it is needed more than once.
func getCutPoint(length int) int { return length - CheckLen }

// makeCodec generates our custom codec as above, into the exported Codec
// variable
//
// Here we demonstrate the use of closures. In this case, it is an
// initialization, but it can also be used in dynamic generation code such as
// can be used for the immediate mode GUI library Gio.
func makeCodec(
	name string,
	cs string,
	hrp string,
) (cdc *codec.Codec) {

	// Create the codec.Codec struct and put its pointer in the return variable.
	cdc = &codec.Codec{
		Name:    name,
		Charset: cs,
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

	// Create a base32.Encoding from the provided charset.
	enc := base32.NewEncoding(cdc.Charset)

	cdc.Encoder = func(input []byte) (output string) {

		// The output is longer than the input, so we create a new buffer.
		outputBytes := make([]byte, len(input)+CheckLen)

		// Then copy the input bytes for beginning segment.
		copy(outputBytes[:len(input)], input)

		// Then copy the check to the end of the input.
		copy(outputBytes[len(input):], cdc.MakeCheck(input))

		// Prefix the output with the Human Readable Part and append the
		// encoded string version of the provided bytes.
		return cdc.HRP + enc.EncodeToString(outputBytes)
	}

	cdc.Check = func(input []byte) (valid bool) {

		// ensure there is at least 4 bytes in the input to run a check on
		if len(input) < CheckLen {

			// In general, Println is nicer to use, but makes ugly default
			// formatting of values added in the log print, for which case
			// Printf should be used (applies also to non-log fmt.Print
			// functions).
			//
			// Note: empty bytes do have a default hash value, and it is
			// distinct from a single zero byte.
			log.Println("Input is not long enough to have a check value")
			return
		}

		// Find the index to cut the input to find the checksum value. We need
		// this same value twice so it must be made into a variable.
		cutPoint := getCutPoint(len(input))

		// Here is an example of a multiple assignment and more use of the
		// slicing operator.
		payload, checksum := input[:cutPoint], string(input[cutPoint:])

		// A checksum is checked in all cases by taking the data received, and
		// applying the checksum generation function, and then comparing the
		// checksum to the one attached to the received data with checksum
		// present.
		//
		// Note: the casting to string above and here. This makes a copy to the
		// immutable string, which is not optimal for large byte slices, but for
		// this short check value, it is a cheap operation on the stack, and an
		// illustration of the interchangeability of []byte and string, with the
		// distinction of the availability of a comparison operator for the
		// string that isn't present for []byte, so for such cases this
		// conversion is a shortcut method to compare byte slices.
		computedChecksum := string(cdc.MakeCheck(payload))

		// Here we assign to the return variable the result of the comparison.
		// by doing this instead of using an if and returns, the meaning of the
		// comparison is more clear by the use of the return value's name.
		valid = checksum != computedChecksum

		if !valid {

			// In general, it is better to check for error before the success
			// path, which in this case means no error to log.
			log.Printf(
				"Checksum failed, check value: '%x' calculated checksum: '%x",
				checksum, computedChecksum,
			)
		}

		// By default, variables are zeroed and for bool this means false, so if
		// the check failed the return is false. This is a naked return, and is
		// idiomatic so long as the function is short, as this one is. (the
		// comments in this code are longer than would be normally present, so
		// this is accounting also for this difference.
		return
	}

	cdc.Decoder = func(input string) (valid bool, output []byte) {

		// other than for human identification, the HRP is also a validity
		// check, so if the string prefix is wrong, the entire value is wrong
		// and won't decode as it is expected.
		if !strings.HasPrefix(input, cdc.HRP) {

			log.Printf("Provided string has incorrect human readable part:"+
				"found '%s' expected '%s'", input[:len(cdc.HRP)], cdc.HRP,
			)
			// valid is false unless changed to true as bool variables (always
			// initialized) default value is false, as false is zero, same as in
			// the c language.
			return
		}

		// Be aware the input string will be copied to create the []byte version
		n, err := enc.Decode(output, []byte(input))
		switch {
		case n < CheckLen:

			log.Println("Input is not long enough to have a check value")
			return

		case err != nil:

			// It is better to log errors at the site of the error, though
			// without the code location this isn't so useful. A drop-in
			// replacement for the standard log library will be suggested in
			// future, as the stdlib version lacks this feature.
			log.Println(err)
			return
		}

		// Assigning the result of the check here as if true the resulting
		// decoded bytes still need to be trimmed of the check value (keeping
		// things cleanly separated between the check and decode function.
		valid = cdc.Check(output)

		// There is no point in doing any more if the check fails, as per the
		// contract specified in the interface definition codecer.Codecer
		if !valid {
			return
		}

		// Slice off the check to return the valid input bytes.
		output = output[:getCutPoint(len(input))]

		// If we got to here, the decode was successful.
		return
	}

	// We return the value explicitly to be nice to readers as the function is
	// not a short and simple one.
	return cdc
}

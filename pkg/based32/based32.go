// Package based32 provides a simplified variant of the standard
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
package based32

import (
	"encoding/base32"
	codec "github.com/quanterall/kitchensink"
	"lukechampine.com/blake3"
	"strings"
)

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
	charset,
	"QNTRL",
)

func getCheckLen(length int) (checkLen int) {

	// In order to provide a minimum of 1 byte of check to the output, while
	// avoiding the encoder adding padding characters (default is '=') the
	// length of the encoded bytes must be rounded to the nearest multiple of 5,
	// adding 5 if it is already a multiple of 5 (5 bytes is 40 bits which
	// encodes as 8 base32 characters).
	//
	// The first byte of the encoded data contains the check length, as this
	// formula varies depending on the length of the data, so it needs to be
	// encoded into the format in the beginning as it can't go at the end. So
	// the check length is one byte less than this formula indicates.
	//
	// This is a significant divergence from the methods used for these encoders
	// because in this tutorial we are not only aiming to produce human readable
	// transcription codes for just transaction hashes (usually 256bit/32 byte)
	// and addresses (usually 160bit/20byte) but a general formula that could
	// encode any binary data length, but presumably it would be likely no more
	// than 512 bits of data for a double length hash, since such a code would
	// take at least a couple of minutes to correctly transcribe.
	//
	// Though a Go programmer may never do a lot of this kind of algorithm
	// design, it is here especially for those who are inclined towards this
	// kind of low level encoding, which is part of any data encoding for wire,
	// storage, for graphic and audio encoding formats, and things like writing
	// GUIs.
	//
	// The following formula ensures that there is at least 1 check byte, up to
	// 4
	//
	// we add two to the length before modulus, as there must be 1 byte for
	// check length and 1 byte of check
	lengthMod := (2 + length) % 5

	// The modulus is subtracted from 5 to produce the complement required to
	// make the correct number of bytes of total data, plus 1 to account for the
	// minimum length of 1.
	checkLen = 5 - lengthMod + 1

	return checkLen
}

// getCutPoint is made into a function because it is needed more than once.
func getCutPoint(length, checkLen int) int {

	return length - checkLen - 1
}

// makeCodec generates our custom codec as above, into the exported Codec
// variable
//
// Here we demonstrate the use of closures. In this case, it is an
// initialization, but it can also be used in dynamic generation code, or to use
// the 'builder' pattern to construct larger algorithms out of small modular
// parts.
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
	cdc.MakeCheck = func(input []byte, checkLen int) (output []byte) {

		// We use the Blake3 256 bit hash because it is nearly as fast as CRC32
		// but less complicated to use due to the 32 bit integer conversions to
		// bytes required to use the CRC32 algorithm.
		checkArray := blake3.Sum256(input)

		// This truncates the blake3 hash to the prescribed check length
		return checkArray[:checkLen]
	}

	// Create a base32.Encoding from the provided charset.
	enc := base32.NewEncoding(cdc.Charset)

	cdc.Encoder = func(input []byte) (output string) {

		// The check length depends on the modulus of the length of the data is
		// order to avoid padding.
		checkLen := getCheckLen(len(input))

		// The output is longer than the input, so we create a new buffer.
		outputBytes := make([]byte, len(input)+checkLen+1)

		// Add the check length byte to the front
		outputBytes[0] = byte(checkLen)

		// Then copy the input bytes for beginning segment.
		copy(outputBytes[1:len(input)+1], input)

		// Then copy the check to the end of the input.
		copy(outputBytes[len(input)+1:], cdc.MakeCheck(input, checkLen))

		// Create the encoding for the output.
		outputString := enc.EncodeToString(outputBytes)

		// We can omit the first character of the encoding because the length
		// prefix never uses the first 5 bits of the first byte, and add it back
		// for the decoder later.
		trimmedString := outputString[1:]

		// Prefix the output with the Human Readable Part and append the
		// encoded string version of the provided bytes.
		return cdc.HRP + trimmedString
	}

	cdc.Check = func(input []byte) (valid bool) {

		// We must do this check or the next statement will cause a bounds check
		// panic. Note that zero length and nil slices are different, but have
		// the same effect in this case, so both must be checked.
		switch {
		case len(input) < 1:

			log.Println("Input of zero length is invalid")
			return

		case input == nil:

			log.Println("Input of nil slice is invalid")
			return
		}

		// The check length is encoded into the first byte in order to ensure
		// the data is cut correctly to perform the integrity check.
		checkLen := int(input[0])

		// Ensure there is at least 4 bytes in the input to run a check on
		if len(input) < checkLen+1 {

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
		cutPoint := getCutPoint(len(input), checkLen)

		// Here is an example of a multiple assignment and more use of the
		// slicing operator.
		payload, checksum := input[1:cutPoint], string(input[cutPoint:])

		// A checksum is checked in all cases by taking the data received, and
		// applying the checksum generation function, and then comparing the
		// checksum to the one attached to the received data with checksum
		// present.
		//
		// Note: The casting to string above and here. This makes a copy to the
		// immutable string, which is not optimal for large byte slices, but for
		// this short check value, it is a cheap operation on the stack, and an
		// illustration of the interchangeability of []byte and string, with the
		// distinction of the availability of a comparison operator for the
		// string that isn't present for []byte, so for such cases this
		// conversion is a shortcut method to compare byte slices.
		computedChecksum := string(cdc.MakeCheck(payload, checkLen))

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

	cdc.Decoder = func(input string) (decRes codec.DecodeRes) {

		// Other than for human identification, the HRP is also a validity
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

		// Cut the HRP off the beginning to get the content, add the initial
		// zeroed 5 bytes with a 'q' character.
		input = "q" + input[len(cdc.HRP):]

		// The length of the base32 string refers to 5 bytes per slice index
		// position, so the correct size of the output bytes, which are 8 bytes
		// per slice index position, is found with the following simple integer
		// math calculation.
		//
		// This allocation needs to be made first as the base32 Decode function
		// does not do this allocation automatically and it would be wasteful to
		// not compute it precisely, when the calculation is so simple.
		//
		// If this allocation is omitted, the decoder will panic due to bounds
		// check error. A nil slice is equivalent to a zero length slice and
		// gives a bounds check error, but in fact, the slice has no data at
		// all. Yes, the panic message is lies:
		//
		//   panic: runtime error: index out of range [4] with length 0
		//
		// If this assignment isn't made, by default, output is nil, not
		// []byte{} so this panic message is deceptive.
		decRes.Data = make([]byte, len(input)*8/5)

		// Be aware the input string will be copied to create the []byte
		// version. Also, because the input bytes are always zero for the first
		// 5 most significant bits, we must re-add the zero at the front (q)
		// before feeding it to the decoder.
		n, err := enc.Decode(decRes.Data, []byte(input))
		if err != nil {

			log.Println(err)
			return
		}

		// The first byte signifies the length of the check at the end
		checkLen := int(decRes.Data[0])
		if n < checkLen+1 {

			log.Println("Input is not long enough to have a check value")
			return

		}

		// Assigning the result of the check here as if true the resulting
		// decoded bytes still need to be trimmed of the check value (keeping
		// things cleanly separated between the check and decode function.
		decRes.Decoded = cdc.Check(decRes.Data)

		// There is no point in doing any more if the check fails, as per the
		// contract specified in the interface definition codecer.Codecer
		if !decRes.Decoded {
			return
		}

		// Slice off the check length prefix, and the check bytes to return the
		// valid input bytes.
		decRes.Data = decRes.Data[1:getCutPoint(len(input), checkLen)]

		// If we got to here, the decode was successful.
		return
	}

	// We return the value explicitly to be nice to readers as the function is
	// not a short and simple one.
	return cdc
}

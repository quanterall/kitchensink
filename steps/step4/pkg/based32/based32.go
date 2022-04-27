// Package based32 provides a simplified variant of the standard
// Bech32 human readable binary codec
package based32

import (
	"encoding/base32"
	codec "github.com/quanterall/kitchensink"
	"github.com/quanterall/kitchensink/pkg/proto"
	"lukechampine.com/blake3"
	"strings"
)

// charset is the set of characters used in the data section of bech32 strings.
const charset = "qpzry9x8gf2tvdw0s3jn54khce6mua7l"

// Codec provides the encoder/decoder implementation created by makeCodec.
var Codec = makeCodec(
	"Base32Check",
	charset,
	"QNTRL",
)

func getCheckLen(length int) (checkLen int) {

	// The following formula ensures that there is at least 1 check byte, up to
	// 4, in order to create a variable length theck that serves also to pad to
	// 5 bytes per 8 5 byte characters (2^5 = 32 for base 32)
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

	cdc.Encoder = func(input []byte) (output string, err error) {

		if len(input) < 1 {

			err = proto.Error_ZERO_LENGTH
			return
		}

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
		output = cdc.HRP + trimmedString

		return
	}

	cdc.Check = func(input []byte) (err error) {

		// We must do this check or the next statement will cause a bounds check
		// panic. Note that zero length and nil slices are different, but have
		// the same effect in this case, so both must be checked.
		switch {
		case len(input) < 1:

			err = proto.Error_ZERO_LENGTH
			return

		case input == nil:

			err = proto.Error_NIL_SLICE
			return
		}

		// The check length is encoded into the first byte in order to ensure
		// the data is cut correctly to perform the integrity check.
		checkLen := int(input[0])

		// Ensure there is at enough bytes in the input to run a check on
		if len(input) < checkLen+1 {

			err = proto.Error_CHECK_TOO_SHORT

			return
		}

		// Find the index to cut the input to find the checksum value. We need
		// this same value twice so it must be made into a variable.
		cutPoint := getCutPoint(len(input), checkLen)

		// Here is an example of a multiple assignment and more use of the
		// slicing operator.
		payload, checksum := input[1:cutPoint], string(input[cutPoint:])

		computedChecksum := string(cdc.MakeCheck(payload, checkLen))

		// Here we assign to the return variable the result of the comparison.
		// by doing this instead of using an if and returns, the meaning of the
		// comparison is more clear by the use of the return value's name.
		valid := checksum != computedChecksum

		if !valid {

			err = proto.Error_CHECK_FAILED
		}

		return
	}

	cdc.Decoder = func(input string) (output []byte, err error) {

		// Other than for human identification, the HRP is also a validity
		// check, so if the string prefix is wrong, the entire value is wrong
		// and won't decode as it is expected.
		if !strings.HasPrefix(input, cdc.HRP) {

			log.Printf("Provided string has incorrect human readable part:"+
				"found '%s' expected '%s'", input[:len(cdc.HRP)], cdc.HRP,
			)

			err = proto.Error_INCORRECT_HUMAN_READABLE_PART

			return
		}

		// Cut the HRP off the beginning to get the content, add the initial
		// zeroed 5 bytes with a 'q' character.
		input = "q" + input[len(cdc.HRP):]

		data := make([]byte, len(input)*5/8)

		// Be aware the input string will be copied to create the []byte
		// version. Also, because the input bytes are always zero for the first
		// 5 most significant bits, we must re-add the zero at the front (q)
		// before feeding it to the decoder.
		var writtenBytes int
		writtenBytes, err = enc.Decode(data, []byte(input))
		if err != nil {

			log.Println(err)
			return
		}

		// The first byte signifies the length of the check at the end
		checkLen := int(data[0])
		if writtenBytes < checkLen+1 {

			err = proto.Error_CHECK_TOO_SHORT

			return

		}

		// Assigning the result of the check here as if true the resulting
		// decoded bytes still need to be trimmed of the check value (keeping
		// things cleanly separated between the check and decode function.
		err = cdc.Check(data)

		// There is no point in doing any more if the check fails, as per the
		// contract specified in the interface definition codecer.Codecer
		if err != nil {
			return
		}

		// Slice off the check length prefix, and the check bytes to return the
		// valid input bytes.
		output = data[1:getCutPoint(len(data)+1, checkLen)]

		// If we got to here, the decode was successful.
		return
	}

	// We return the value explicitly to be nice to readers as the function is
	// not a short and simple one.
	return cdc
}

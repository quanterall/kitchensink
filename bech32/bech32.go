// Package bech32 provides a (somewhat simplified) version of the standard
// Bech32 human readable binary codec for rendering hashes into a form that can
// be read and transcribed by humans, as used in Bitcoin Segregated Witness
// addresses and the Cosmos SDK framework, amongst others.
//
// BIP 0173 https://en.bitcoin.it/wiki/BIP_0173 is the specification for this
// encoding.
package bech32

import (
	"fmt"
	"github.com/cosmos/btcutil/bech32"
	codec "github.com/quanterall/kitchensink"
)

// Spec is the collection of elements derived from the codec type definition
// that creates the concrete implementation of a 'generic' functionality.
var Spec = codec.Codec{
	Name:    "Bech32",
	HRP:     "cybriq",
	Charset: "qpzry9x8gf2tvdw0s3jn54khce6mua7l",
	Encoder: func(input []byte) (output string) {
		return ""
	},
	Decoder: func(input string) (valid bool, output []byte) {
		return true, nil
	},
	MakeCheck: func(input []byte) (output []byte) {
		return nil
	},
	Check: func(input string) (valid bool) {
		return false
	},
}

// ConvertAndEncode converts from a base64 encoded byte string to base32 encoded byte string and then to bech32.
func ConvertAndEncode(hrp string, data []byte) (string, error) {
	converted, err := bech32.ConvertBits(data, 8, 5, true)
	if err != nil {
		return "", fmt.Errorf("encoding bech32 failed: %w", err)
	}

	return bech32.Encode(hrp, converted)
}

// DecodeAndConvert decodes a bech32 encoded string and converts to base64 encoded bytes.
func DecodeAndConvert(bech string) (string, []byte, error) {
	hrp, data, err := bech32.Decode(bech, 1023)
	if err != nil {
		return "", nil, fmt.Errorf("decoding bech32 failed: %w", err)
	}

	converted, err := bech32.ConvertBits(data, 5, 8, false)
	if err != nil {
		return "", nil, fmt.Errorf("decoding bech32 failed: %w", err)
	}

	return hrp, converted, nil
}

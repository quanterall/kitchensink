// Package bech32 provides a (somewhat simplified) version of the standard
// Bech32 human readable binary codec for rendering hashes into a form that can
// be read and transcribed by humans, as used in Bitcoin Segregated Witness
// addresses and the Cosmos SDK framework, amongst others.
//
// BIP 0173 https://en.bitcoin.it/wiki/BIP_0173 is the specification for this
// encoding.
package bech32

import (
	"github.com/cybriq/transcribe/codec"
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

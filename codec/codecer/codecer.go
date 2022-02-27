// Package codecer is the interface definition for a Human Readable Binary
// Transcription Codec
//
// Interface definitions should be placed in separate packages to
// implementations so there is no risk of a circular dependency, which is not
// permitted in Go, because this kind of automated interpretation of programmer
// intent is the most expensive thing (time, processing, memory) that compilers
// do.
package codecer

// Codecer is the externally usable interface which provides a check for
// complete implementation as well as illustrating the use of interfaces in Go.
//
// It may seem somewhat redundant in light of type definition, in the root of
// the repository, which exposes the exact same Encode and Decode functions, but
// the purpose of adding this is that this interface can be implemented without
// using the concrete Codec type above, should the programmer have a need to do
// so.
//
// The implementation only needs to implement these two functions and then
// whatever structure ties it together can be passed around without needing to
// know anything about its internal representations or implementation details.
//
// The purpose of interfaces in Go is exactly to eliminate dependencies on any
// concrete data types so the implementations can be changed without changing
// the consuming code
type Codecer interface {
	// Encode takes an arbitrary length byte input and returns the output as
	// defined for the codec
	Encode(input []byte) (output string)
	// Decode takes an encoded string and returns if the encoding is valid and
	// the value passes any check function defined for the type
	Decode(input string) (valid bool, output []byte)
}

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
// the consuming code.
//
// We are adding this interface in addition to the use of a struct and closure
// pattern mainly as illustration but also to make sure the student is aware of
// the implicit implementation recognition, the way to make the compile time
// check of implementation, and as an exercise for later, the student can create
// their own implementation by importing this package and use the provided
// implementation, in parallel with their own, or without it, which they can
// implement with an entirely separate and different data structure (which will
// be a struct, most likely, though it can be a slice of interface and be even
// subordinate to another structured variable like a slice of interface, or a
// map of interfaces. Then they can drop this interface in place of the built in
// one and see that they don't have to change the calling code.
//
// Note: though it is not officially recognised as idiomatic, it is the opinion
// of the author of this tutorial that the return values of interface function
// signatures should be named, as it makes no sense to force the developer to
// have to read through the implementation that *idiomatically* should accompany
// an interface, as by idiom, interface should be avoided unless there is more
// than one implementation.
type Codecer interface {

	// Encode takes an arbitrary length byte input and returns the output as
	// defined for the codec
	Encode(input []byte) (output string)

	// Decode takes an encoded string and returns if the encoding is valid and
	// the value passes any check function defined for the type
	Decode(input string) (valid bool, output []byte)
}

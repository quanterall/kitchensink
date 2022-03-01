package transcribe

import (
	"github.com/quanterall/kitchensink/codecer"
)

// Codec is the collection of elements that creates a Human Readable Binary
// Codec
//
// This is an example of the use of a structure definition to encapsulate and
// logically connect together all of the elements of an implementation, while
// also permitting this to be used by external code without further
// dependencies, either through this type, or via the interface defined further
// down.
//
// It is not "official" idiom, but it's the opinion of the author of this
// tutorial that return values given in type specifications like this helps the
// users of the library understand what the return values actually are.
// Otherwise, the programmer is forced to read the whole function just to spot
// the names and, even worse, comments explaining what the values are.
type Codec struct {

	// Name is the human readable name given to this encoder
	Name string

	// HRP is the Human Readable Prefix to be appended in front of the encoding
	// to disambiguate it from another encoding or as a network or protocol
	// identifier
	HRP string

	// Charset is the set of characters that the encoder uses.
	//
	// Note that the length of the string also sets the base, so if there is 16,
	// it is base 16, if there is 32, it is base32, The length is arbitrary and
	// the numeric format encoder in the Go standard library can be used for
	// literally set of symbols of any length.
	Charset string

	// Encode takes an arbitrary length byte input and returns the output as
	// defined for the codec
	Encoder func(input []byte) (output string)

	// Decode takes an encoded string and returns if the encoding is valid and
	// the value passes any check function defined for the type
	Decoder func(input string) (valid bool, output []byte)

	// AddCheck is used by Encode to add extra bytes for the checksum to ensure
	// correct input so user does not send to a wrong address by mistake, for
	// example.
	MakeCheck func(input []byte) (output []byte)

	// Check returns whether the check is valid
	Check func(input []byte) (valid bool)
}

// The following implementations are here to ensure this type implements the
// interface. In this tutorial/example we are creating a kind of generic
// implementation through the use of closures loaded into a struct.
//
// Normally a developer would use either one, or the other, a struct with
// closures, OR an interface with arbitrary variable with implementations for
// the created type.
//
// In order to illustrate both interfaces and the use of closures with a struct
// in this way we combine the two things by invoking the closures in a
// predefined pair of methods that satisfy the interface.
//
// In fact, there is no real reason why this design could not be standard idiom,
// since satisfies most of the requirements of idiom for both interfaces
// (minimal) and hot-reloadable interfaces (allowing creation of registerable
// compile time plugins such as used in database drivers with structs, and the
// end user can then either use interfaces or the provided struct, and both
// options are open.

// This ensures the interface is satisfied for codecer.Codecer and is removed in
// the generated binary because the underscore indicates the value is discarded.
var _ codecer.Codecer = &Codec{}

// Encode implements the codecer.Codecer.Encode by calling the provided
// function, and allows the concrete Codec type to always satisfy the interface,
// while allowing it to be implemented entirely differently.
//
// note: short functions like this can be one-liners according to gofmt.
func (c Codec) Encode(input []byte) string { return c.Encoder(input) }

// Decode implements the codecer.Codecer.Decode by calling the provided
// function, and allows the concrete Codec type to always satisfy the interface,
// while allowing it to be implemented entirely differently.
//
// note: this also can be a one liner. Since we name the return values in the
// type definition and interface, omitting them here makes the line short enough
// to be a one liner.
func (c Codec) Decode(input string) (bool, []byte) { return c.Decoder(input) }

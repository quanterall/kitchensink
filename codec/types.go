package codec

// Codec is the collection of elements that creates a Human Readable Binary
// Codec
//
// This is an example of the use of a structure definition to encapsulate and
// logically connect together all of the elements of an implementation, while
// also permitting this to be used by external code without further dependencies
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
	Check func(input string) (valid bool)
}

// Encode implements the Codecer.Encoder by calling the provided function, and allows
// the concrete Codec type to always satisfy the interface, while allowing it to
// be implemented entirely differently
func (c Codec) Encode(input []byte) (output string) {
	return c.Encoder(input)
}

// Decode implements the Codecer.Decoder by calling the provided function, and allows
// the concrete Codec type to always satisfy the interface, while allowing it to
// be implemented entirely differently
func (c Codec) Decode(input string) (valid bool, output []byte) {
	return c.Decoder(input)
}

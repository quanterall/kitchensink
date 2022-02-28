package transcribe

import (
	"github.com/quanterall/kitchensink/codecer"
)

func NewCodec(cdc *Codec) codecer.Codecer {
	// Make sure the provided codec has all the parts that are used in the
	// interface
	if cdc.Encoder == nil ||
		cdc.Decoder == nil {
		// panic should not be in production code execution paths, but SHOULD be
		// in execution paths that fail due to programmer errors, that must be
		// fixed before production release.
		panic("Programmer Error: " +
			"codec does not have all necessary functions implemented",
		)
	}
	return cdc
}

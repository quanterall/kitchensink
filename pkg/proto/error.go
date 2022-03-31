// Package proto is the protocol buffers specification and generated code
// package for based32
//
// The
package proto

import (
	transcribe "github.com/quanterall/kitchensink"
)

// Error implements the Error interface which allows this error to automatically
// generate from the error code.
//
// Fixes a bug in the generated code, which not
// only lacks the Error method it uses int32 for the error string map when it
// should be using the defined Error type. No easy way to report the bug in the
// code.
//
// With this method implemented, one can simply return the error map code
// protos.Error_ERROR_NAME_HERE and logs print this upper case snake case which
// means it can be written to be informative in the proto file and concise in
// usage, and with this tiny additional helper, very easy to return, and print.
func (x Error) Error() string {

	return Error_name[int32(x)]
}

// CreateEncodeResponse is a helper to turn a transcribe.EncodeRes into an
// EncodeResponse to be returned to a gRPC client.
func CreateEncodeResponse(res transcribe.EncodeRes) (response *EncodeResponse) {

	// First, create the response structure.
	response = &EncodeResponse{}

	// Because the protobuf struct is essentially a Variant, a structure that
	// does not exist in Go, there is an implicit contract that if there is an
	// error, there is no return value. This is not implicit in Go's tuple
	// returns.
	//
	// Thus, if there is an error, we return that, otherwise, the value in the
	// response.

	if res.Error != nil {
		response.Encoded = &EncodeResponse_Error{
			Error(Error_value[res.Error.Error()]),
		}
	} else {
		response.Encoded =
			&EncodeResponse_EncodedString{
				res.String,
			}
	}
	return
}

// CreateDecodeResponse is a helper to turn a transcribe.DecodeRes into an
// DecodeResponse to be returned to a gRPC client.
func CreateDecodeResponse(res transcribe.DecodeRes) (response *DecodeResponse) {

	// First, create the response structure.
	response = &DecodeResponse{}

	// Return an error if there is an error, otherwise return the response data.
	if res.Error != nil {
		response.Decoded = &DecodeResponse_Error{
			Error(Error_value[res.Error.Error()]),
		}
	} else {
		response.Decoded = &DecodeResponse_Data{res.Bytes}
	}
	return
}

package protos

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
// usage.
func (x Error) Error() string {

	return Error_name[int32(x)]
}

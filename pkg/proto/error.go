package protos

import (
	"errors"
)

// GetError compensates for a bug in protoc-gen-go not making it simpler to pull
// the error strings from their constant values.
func GetError(err Error) error {
	return errors.New(Error_name[int32(err)])
}

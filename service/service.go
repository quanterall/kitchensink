package b32svc

import (
	"time"
)

// BufferSize defines the buffer size created for call channels
var BufferSize = 8

// The following types define the parameters and results returned from the
// ServiceHandlers
type (
	// None means no parameters or return value it is not checked so it can be nil
	None      struct{}
	EncodeCmd struct {
		// Bytes is the slice of bytes to be encoded. Note that slices are
		// reference types so there is no need to make this a pointer to prevent
		// value copy.
		Bytes []byte
	}
	EncodeRes struct {
		// Res is the result string. It is a pointer to a string because
		// otherwise the value will be copied, and it will always be more data
		// than the pointer (8 bytes).
		Res *string
		Err error
	}
	DecodeCmd struct {
		String *string
	}
	DecodeRes struct {
		Res []byte
		Err error
	}
)

// ServiceHandlers provides the concurrent implementation of the codec service
// via the 'send command with return channel' model - the calling code sends a
// command with the result channel inside, and when the handler has performed
// the task it returns on the return channel.
var ServiceHandlers = Handlers{

	"Encode": {
		Fn: func(
			svc *Service,
			cmd interface{},
			timeout time.Duration,
			cancel chan struct{},
		) (res interface{}, err error) {

			return
		},
		Call:   make(chan API, BufferSize),
		Result: func() API { return API{Ch: make(chan EncodeRes)} },
	},

	"Decode": {
		Fn: func(
			svc *Service,
			cmd interface{},
			timeout time.Duration,
			cancel chan struct{},
		) (res interface{}, err error) {

			return
		},
		Call:   make(chan API, BufferSize),
		Result: func() API { return API{Ch: make(chan DecodeRes)} },
	},
}

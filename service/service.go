package b32svc

const (
	BufferSize = 8
)

// ServiceHandlers provides the concurrent implementation of the codec service
// via the 'send command with return channel' model - the calling code sends a
// command with the result channel inside, and when the handler has performed
// the task it returns on the return channel.
var ServiceHandlers = Handlers{
	"encode": {
		Fn: func(
			svc *Service,
			cmd interface{},
			cancel chan struct{},
		) (res interface{}, err error) {

			return
		},
		Call:   make(chan API, BufferSize),
		Result: func() API { return API{Ch: make(chan EncodeRes)} },
	},
	"decode": {
		Fn: func(
			svc *Service,
			cmd interface{},
			cancel chan struct{},
		) (res interface{}, err error) {

			return
		},
		Call:   make(chan API, BufferSize),
		Result: func() API { return API{Ch: make(chan DecodeRes)} },
	},
}

package b32svc

// It is usually best to keep type definitions separated from variable and
// function definitions so as to make it simpler where to find them, and, as
// used in this tutorial, `types.go` is a logical name, if used consistently,
// makes it a predictable task to locate type definitions, without the
// assistance of hyperlinking in an IDE. Hyperlinking is an indispensable
// technology, regardless, which is why the tutorial author recommends Goland.
// However, for the sake of tidiness, keeping central things together makes for
// an easier manual search and reduces the chance of cruft building up in a
// codebase.

// Service ties the concurrent service API together
type Service struct {
	handlers Handlers
}

// API stores the channel, parameters and result values from calls via
// the channel
type API struct {
	Ch     interface{}
	Params interface{}
	Result interface{}
}

// CommandHandler defines an API call
type CommandHandler struct {
	// Fn is the handler for an API call
	Fn func(
		svc *Service,
		cmd interface{},
		cancel chan struct{},
	) (res interface{}, err error)
	// Call is the channel to send a command to the handler
	Call chan API
	// Result is the container that will be filled with the result
	Result func() API
}

// Handlers is a collection of named CommandHandler items
type Handlers map[string]CommandHandler

// The following types define the results returned from the ServiceHandlers
type (
	EncodeRes struct {
		Res []byte
		Err error
	}
	DecodeRes struct {
		Res *string
		Err error
	}
)

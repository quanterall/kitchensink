package server

import (
	"github.com/quanterall/kitchensink/pkg/proto"
	"go.uber.org/atomic"
	"google.golang.org/grpc"
	"io"
	"net"
)

// b32 is not exported because if the consuming code uses this struct directly
// without initializing it correctly, several things will not work correctly,
// such as the stop function, which depends on there being an initialized
// channel, and will panic the Start function immediately.
type b32 struct {
	protos.UnimplementedTranscriberServer
	stop        chan struct{}
	svr         *grpc.Server
	transcriber *Transcriber
	addr        *net.TCPAddr
	roundRobin  atomic.Uint32
	workers     uint32
}

// Encode is our implementation of the encode API call for the incoming stream
// of requests.
//
// Note that both this and the next stream handler are virtually identical
// except for the destination that received messages will be sent to. There is
// ways to make this more DRY, but they are not worth doing for only two API
// calls. If there were 5 or more, the right solution would be a code generator.
// It is a golden rule of Go, if it's not difficult to maintain, copy and paste,
// if it is, write a generator, or rage quit and use a generics language and
// lose your time waiting for compilation instead.
func (b *b32) Encode(stream protos.Transcriber_EncodeServer) error {
out:
	for {

		// check whether it's shutdown time first
		select {
		case <-b.stop:
			break out
		default:
		}

		// Wait for and load in a newly received message
		in, err := stream.Recv()
		switch {
		case err == io.EOF:

			// The client has broken the connection, so we can quit
			break out
		case err != nil:

			// Any error is terminal here, so return it to the caller after
			// logging it
			log.Println(err)
			return err
		}

		log.Printf("received raw bytes %x", in.Data)
		worker := b.roundRobin.Load()
		if worker > b.workers {
			worker = 0
			b.roundRobin.Store(0)
		} else {
			worker++
			b.roundRobin.Inc()
		}
		b.transcriber.encode[worker] <- in
		res := <-b.transcriber.encodeRes[worker]
		response := &protos.EncodeResponse{
			Encoded: &protos.EncodeResponse_EncodedString{
				EncodedString: res,
			},
		}
		err = stream.Send(response)
		if err != nil {
			return err
		}
	}
	return nil
}

// Decode is our implementation of the encode API call for the incoming stream
// of requests.
func (b *b32) Decode(stream protos.Transcriber_DecodeServer) error {

	var decRes *protos.DecodeResponse
out:
	for {

		// check whether it's shutdown time first
		select {
		case <-b.stop:
			break out
		default:
			// When a select statement has a default case it always terminates.
		}

		// Wait for and load in a newly received message
		in, err := stream.Recv()
		switch {
		case err == io.EOF:

			// The client has broken the connection, so we can quit
			break out

		case err != nil:

			// Any error is terminal here, so return it to the caller after
			// logging it, and ending this function terminates the decoder
			// service.
			log.Println(err)
			return err
		}
		log.Printf("received encoded string %s",
			in.EncodedString,
		)
		worker := b.roundRobin.Load()
		if worker > b.workers {
			b.roundRobin.Store(0)
		} else {
			b.roundRobin.Inc()
		}
		b.transcriber.decode[worker] <- in
		res := <-b.transcriber.decodeRes[worker]
		if res.Decoded {
			decRes = &protos.DecodeResponse{
				Decoded: &protos.DecodeResponse_Data{
					Data: res.Data,
				},
			}
		}
		err = stream.Send(decRes)
		if err != nil {
			return err
		}
	}
	return nil
}

// New creates a new service handler
func New(addr *net.TCPAddr, workers uint32) (b *b32) {

	// It would be possible to interlink all of the kill switches in an
	// application via passing this variable in to the New function, for which
	// reason in an application, its killswitch has to trigger closing of this
	// channel via calling the stop function returned by Start, further down.
	stop := make(chan struct{})
	b = &b32{
		stop:        stop,
		svr:         grpc.NewServer(),
		transcriber: NewWorkerPool(workers, stop),
		addr:        addr,
		workers:     workers,
	}
	b.roundRobin.Store(0)

	return
}

func (b *b32) Start() (stop func()) {

	stopService := b.transcriber.Start()

	// Set up a tcp listener for the gRPC service.
	lis, err := net.ListenTCP("tcp", b.addr)
	if err != nil {
		log.Fatalf("failed to listen on %v: %v", b.addr, err)
	}

	// This is spawned in a goroutine so we can trigger the shutdown correctly.
	go func() {
		protos.RegisterTranscriberServer(b.svr, b)
		log.Printf("server listening at %v", lis.Addr())

		if err := b.svr.Serve(lis); err != nil {

			// This is where errors returned from Decode and Encode streams end
			// up.
			log.Printf("failed to serve: '%v'", err)

			// By the time this happens the second goroutine is running and it
			// is always better unless you are sure nothing else is running and
			// part way starting up, to shut it down properly. Closing this
			// channel terminates the second goroutine which calls the server to
			// stop, and then the Start function terminates. In this way we can
			// be sure that nothing will keep running and the user does not have
			// to use `kill -9` or ctrl-\ on the terminal to end the process.
			//
			// If force kill is required, there is a bug in the concurrency and
			// should be fixed to ensure that all resources are properly
			// released, and especially in the case of databases or file writing
			// that the cache is flushed and the on disk store is left in a sane
			// state.
			close(b.stop)
		}
		log.Printf("server at %v now shut down",
			lis.Addr(),
		)

	}()

	go func() {
	out:
		for {
			select {
			case <-b.stop:

				// This is the proper way to stop the gRPC server, which will
				// end the next goroutine spawned just above correctly.
				b.svr.GracefulStop()
				stopService()
				break out
			}
		}
	}()

	// The stop signal is triggered when this function is called, which triggers
	// the graceful stop of the server, and terminates the two goroutines above
	// cleanly.
	return func() { close(b.stop) }
}

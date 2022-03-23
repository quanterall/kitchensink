package server

import (
	"fmt"
	"github.com/quanterall/kitchensink/pkg/proto"
	"google.golang.org/grpc"
	"io"
	"net"
)

// b32 is not exported because if the consuming code uses this struct
// directly without initializing it correctly, several things will not work
// correctly, such as the stop function, which depends on there being an
// initialized channel, and will panic the Start function immediately.
type b32 struct {
	protos.UnimplementedTranscriberServer
	stop        chan struct{}
	svr         *grpc.Server
	transcriber *Transcriber
	port        int
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
		b.transcriber.encode <- in
	}
	return nil
}

// Decode is our implementation of the encode API call for the incoming stream
// of requests.
func (b *b32) Decode(stream protos.Transcriber_DecodeServer) error {

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
		b.transcriber.decode <- in
	}
	return nil
}

// New creates a new service handler
func New(port, workers int) (b *b32) {

	stop := make(chan struct{})
	b = &b32{
		stop:        stop,
		svr:         grpc.NewServer(),
		transcriber: NewWorkerPool(workers, stop),
		port:        port,
	}

	return
}

func (b *b32) Start() (stop func()) {

	// set up a tcp listener for the gRPC service
	lis, err := net.Listen("tcp", fmt.Sprintf(":%d", b.port))
	if err != nil {
		log.Fatalf("failed to listen: %v", err)
	}

	// This is spawned in a goroutine so we can trigger the shutdown correctly.
	go func() {
		protos.RegisterTranscriberServer(b.svr, b)
		log.Printf("server listening at %v", lis.Addr())

		if err := b.svr.Serve(lis); err != nil {
			log.Fatalf("failed to serve: %v", err)
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
				break out
			}
		}
	}()

	// The stop signal is triggered when this function is called, which triggers
	// the graceful stop of the server, and terminates the two goroutines above
	// cleanly.
	return func() { close(b.stop) }
}

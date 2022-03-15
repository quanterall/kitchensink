package server

import (
	"context"
	"flag"
	"fmt"
	protos "github.com/quanterall/kitchensink/proto"
	"go.uber.org/atomic"
	"google.golang.org/grpc"
	"net"
)

var (
	port = flag.Int("port", 50051, "The server port")
)

// based32 is not exported because if the consuming code uses this struct
// directly without initializing it correctly, several things will not work
// correctly, such as the stop function, which depends on there being an
// initialized channel, and will panic the Start function immediately.
type based32 struct {
	protos.UnimplementedTranscriberServer
	stop                       chan struct{}
	encCallCount, decCallCount *atomic.Uint32
	svr                        *grpc.Server
}

func New() (b *based32) {
	b = &based32{
		stop:         make(chan struct{}),
		encCallCount: atomic.NewUint32(0),
		decCallCount: atomic.NewUint32(0),
		svr:          grpc.NewServer(),
	}
	return
}

func (b *based32) Encode(context.Context, *protos.EncodeRequest,
) (*protos.EncodeResponse, error) {
	b.encCallCount.Inc()
	panic("not implemented")
}

func (b *based32) Decode(context.Context, *protos.DecodeRequest,
) (*protos.DecodeResponse, error) {
	b.decCallCount.Dec()
	panic("not implemented")
}

func (b *based32) Start() (stop func()) {

	// if the calling main function was passed a port specification, this loads
	// it into the port variable
	flag.Parse()
	lis, err := net.Listen("tcp", fmt.Sprintf(":%d", *port))
	if err != nil {
		log.Fatalf("failed to listen: %v", err)
	}

	protos.RegisterTranscriberServer(b.svr, b)
	log.Printf("server listening at %v", lis.Addr())

	// This is spawned in a goroutine so we can trigger the shutdown correctly
	go func() {
		if err := b.svr.Serve(lis); err != nil {
			log.Fatalf("failed to serve: %v", err)
		}
	}()

	go func() {
	out:
		for {
			select {
			case <-b.stop:

				// this is the proper way to stop the gRPC server, which will
				// end the next goroutine spawned just above correctly.
				b.svr.GracefulStop()
				break out
			}
		}
	}()

	// the stop signal is triggered
	return func() { close(b.stop) }
}

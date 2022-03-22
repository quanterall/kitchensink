package client

import (
	"flag"
	protos "github.com/quanterall/kitchensink/proto"
	"golang.org/x/net/context"
	"google.golang.org/grpc"
	"log"
)

var serverAddr = flag.String("addr", "localhost:50051",
	"The server address in the format of host:port",
)

type Transcribe struct {
	*grpc.ClientConn
	protos.TranscriberClient
	enc protos.Transcriber_EncodeClient
	dec protos.Transcriber_DecodeClient
	context.Context
}

func New() (c *Transcribe, disconnect func()) {

	// Dial the configured server address
	conn, err := grpc.Dial(*serverAddr)
	if err != nil {
		log.Fatalf("fail to dial: %v", err)
	}

	client := protos.NewTranscriberClient(conn)

	ctx, cancel := context.WithCancel(context.Background())
	var encode protos.Transcriber_EncodeClient
	encode, err = client.Encode(ctx)
	if err != nil {
		return nil, func() {}
	}
	var decode protos.Transcriber_DecodeClient
	decode, err = client.Decode(ctx)
	if err != nil {
		return nil, func() {}
	}
	// Set up the connection to the server.
	c = &Transcribe{
		ClientConn:        conn,
		TranscriberClient: client,
		enc:               encode,
		dec:               decode,
	}

	// Call the returned disconnect() to stop the client.
	return c, func() {

		// shut down the cancel context
		cancel()

		// close the connection
		err := conn.Close()
		if err != nil {
			log.Println(err)
		}
	}
}

// Encode bytes into human readable transcription code based32.
func (t *Transcribe) Encode(req *protos.EncodeRequest) (
	res *protos.EncodeResponse,
	err error,
) {
	err = t.enc.Send(req)
	if err != nil {
		return
	}
	res, err = t.enc.Recv()
	if err != nil {
		return nil, err
	}
	return
}

// Decode bytes from human readable transcription code based32.
func (t *Transcribe) Decode(req *protos.DecodeRequest) (
	res *protos.DecodeResponse,
	err error,
) {
	err = t.dec.Send(req)
	if err != nil {
		return
	}
	res, err = t.dec.Recv()
	if err != nil {
		return nil, err
	}
	return

}

package client

import (
	protos2 "github.com/quanterall/kitchensink/pkg/proto"
	"golang.org/x/net/context"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
)

//
// var serverAddr = flag.String("addr", "localhost:50051",
// 	"The server address in the format of host:port",
// )

type Transcribe struct {
	*grpc.ClientConn
	protos2.TranscriberClient
	enc protos2.Transcriber_EncodeClient
	dec protos2.Transcriber_DecodeClient
	context.Context
}

func New(serverAddr string) (c *Transcribe, disconnect func()) {

	// Dial the configured server address
	conn, err := grpc.Dial(serverAddr,
		grpc.WithTransportCredentials(insecure.NewCredentials()),
	)
	if err != nil {
		log.Fatalf("fail to dial: %v", err)
	}

	client := protos2.NewTranscriberClient(conn)

	ctx, cancel := context.WithCancel(context.Background())
	var encode protos2.Transcriber_EncodeClient
	encode, err = client.Encode(ctx)
	if err != nil {
		return nil, func() {}
	}
	var decode protos2.Transcriber_DecodeClient
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
func (t *Transcribe) Encode(req *protos2.EncodeRequest) (
	res *protos2.EncodeResponse,
	err error,
) {
	err = t.enc.Send(req)
	if err != nil {
		return
	}
	res, err = t.enc.Recv()
	if err != nil {
		return
	}
	return
}

// Decode bytes from human readable transcription code based32.
func (t *Transcribe) Decode(req *protos2.DecodeRequest) (
	res *protos2.DecodeResponse,
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

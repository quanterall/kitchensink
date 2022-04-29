package client

import (
	"context"
	"github.com/quanterall/kitchensink/pkg/proto"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
	"io"
	"time"
)

type EncReq struct {
	Req *proto.EncodeRequest
	Res chan *proto.EncodeResponse
}

func NewEncReq(req *proto.EncodeRequest) EncReq {
	return EncReq{Req: req, Res: make(chan *proto.EncodeResponse)}
}

type DecReq struct {
	Req *proto.DecodeRequest
	Res chan *proto.DecodeResponse
}

func NewDecReq(req *proto.DecodeRequest) DecReq {
	return DecReq{Req: req, Res: make(chan *proto.DecodeResponse)}
}

type b32c struct {
	encChan    chan EncReq
	decChan    chan DecReq
	enc        proto.Transcriber_EncodeClient
	dec        proto.Transcriber_DecodeClient
	stop       <-chan struct{}
	ctx        context.Context
	Cancel     func()
	timeout    time.Duration
	waitingEnc map[time.Time]EncReq
	waitingDec map[time.Time]DecReq
	*grpc.ClientConn
	proto.TranscriberClient
}

func NewClient(serverAddr string, timeout time.Duration) (
	client *b32c, err error,
) {

	// Dial the configured server address
	conn, err := grpc.Dial(
		serverAddr,
		grpc.WithTransportCredentials(insecure.NewCredentials()),
	)
	if err != nil {
		log.Fatalf("fail to dial: %v", err)
	}

	cli := proto.NewTranscriberClient(conn)
	ctx, cancel := context.WithCancel(context.Background())
	var encode proto.Transcriber_EncodeClient
	encode, err = cli.Encode(ctx)
	if err != nil {
		return
	}
	var decode proto.Transcriber_DecodeClient
	decode, err = cli.Decode(ctx)
	if err != nil {
		return
	}

	client = &b32c{
		encChan:           make(chan EncReq),
		decChan:           make(chan DecReq),
		enc:               encode,
		dec:               decode,
		stop:              ctx.Done(),
		ctx:               ctx,
		Cancel:            cancel,
		timeout:           timeout,
		waitingEnc:        make(map[time.Time]EncReq),
		waitingDec:        make(map[time.Time]DecReq),
		ClientConn:        conn,
		TranscriberClient: cli,
	}

	return
}

func (b *b32c) Encode(stream proto.Transcriber_EncodeClient) (err error) {

	go func(stream proto.Transcriber_EncodeClient) {
	out:
		for {
			select {
			case <-b.stop:
				break out
			case msg := <-b.encChan:
				err := stream.Send(msg.Req)
				if err != nil {
					log.Print(err)
				}
				b.waitingEnc[time.Now()] = msg
			}
		}
	}(stream)

	go func(stream proto.Transcriber_EncodeClient) {
	in:
		for {

			// check whether it's shutdown time first
			select {
			case <-b.stop:
				break in
			default:
			}

			// Wait for and load in a newly received message
			recvd, err := stream.Recv()
			switch {
			case err == io.EOF:

				// The client has broken the connection, so we can quit
				break in
			case err != nil:

				// Any error is terminal here, so return it to the caller after
				// logging it
				log.Println(err)
			}
			for i := range b.waitingEnc {

				// Check for expired responses
				if i.Add(b.timeout).After(time.Now()) {
					delete(b.waitingEnc, i)
				}

				// Return received responses
				if recvd.IdNonce == b.waitingEnc[i].Req.IdNonce {
					b.waitingEnc[i].Res <- recvd
					delete(b.waitingEnc, i)
				}
			}
		}
	}(stream)

	return
}

func (b *b32c) Decode(stream proto.Transcriber_DecodeClient) (err error) {

	go func(stream proto.Transcriber_DecodeClient) {
	out:
		for {
			select {
			case <-b.stop:
				break out
			case msg := <-b.decChan:
				err := stream.Send(msg.Req)
				if err != nil {
					log.Print(err)
				}
				b.waitingDec[time.Now()] = msg
			}
		}
	}(stream)

	go func(stream proto.Transcriber_DecodeClient) {
	in:
		for {

			// check whether it's shutdown time first
			select {
			case <-b.stop:
				break in
			default:
			}

			// Wait for and load in a newly received message
			recvd, err := stream.Recv()
			switch {
			case err == io.EOF:

				// The client has broken the connection, so we can quit
				break in
			case err != nil:

				// Any error is terminal here, so return it to the caller after
				// logging it
				log.Println(err)
			}
			for i := range b.waitingEnc {

				// Check for expired responses
				if i.Add(b.timeout).After(time.Now()) {
					delete(b.waitingDec, i)
				}

				// Return received responses
				if recvd.IdNonce == b.waitingEnc[i].Req.IdNonce {
					b.waitingDec[i].Res <- recvd
					delete(b.waitingDec, i)
				}
			}
		}
	}(stream)

	return
}

func (b *b32c) Start() (
	send func(*proto.EncodeRequest) *proto.EncodeResponse,
	recv func(*proto.DecodeRequest) *proto.DecodeResponse,
) {

	go func() {
		err := b.Encode(b.enc)
		if err != nil {
			log.Print(err)
		}
	}()
	go func() {
		err := b.Decode(b.dec)
		if err != nil {
			log.Print(err)
		}
	}()

	return func(req *proto.EncodeRequest) *proto.EncodeResponse {
			r := NewEncReq(req)
			b.encChan <- r
			return <-r.Res
		},
		func(req *proto.DecodeRequest) *proto.DecodeResponse {
			r := NewDecReq(req)
			b.decChan <- r
			return <-r.Res
		}

}

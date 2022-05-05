package client

import (
	"context"
	"github.com/quanterall/kitchensink/pkg/proto"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
	"time"
)

type EncReq struct {
	Req *proto.EncodeRequest
	Res chan *proto.EncodeResponse
}

func NewEncReq(req *proto.EncodeRequest) EncReq {
	req.IdNonce = uint64(time.Now().UnixNano())
	return EncReq{Req: req, Res: make(chan *proto.EncodeResponse)}
}

type DecReq struct {
	Req *proto.DecodeRequest
	Res chan *proto.DecodeResponse
}

func NewDecReq(req *proto.DecodeRequest) DecReq {
	req.IdNonce = uint64(time.Now().UnixNano())
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

func New(serverAddr string, timeout time.Duration) (
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

// Start up the client. Call b.Cancel() to stop.
//
// The returned send and recv functions are async by default
// but can be used synchronously by receiving from them directly:
//
//     encRes := <-enc(req)
//     decRes := <-dec(req)
func (b *b32c) Start() (
	enc func(*proto.EncodeRequest) chan *proto.EncodeResponse,
	dec func(*proto.DecodeRequest) chan *proto.DecodeResponse,
) {

	go func() {
		log.Println("starting decoder")
		err := b.Decode(b.dec)
		if err != nil {
			log.Print(err)
		}
	}()

	go func() {
		log.Println("starting encoder")
		err := b.Encode(b.enc)
		if err != nil {
			log.Print(err)
		}
	}()

	enc = func(req *proto.EncodeRequest) chan *proto.EncodeResponse {
		r := NewEncReq(req)
		b.encChan <- r
		return r.Res
	}
	dec = func(req *proto.DecodeRequest) chan *proto.DecodeResponse {
		r := NewDecReq(req)
		b.decChan <- r
		return r.Res
	}
	return
}

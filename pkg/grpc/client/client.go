package client

import (
	"context"
	"github.com/quanterall/kitchensink/pkg/proto"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
	"time"
)

func New(serverAddr string, timeout time.Duration) (
	client *b32c, err error,
) {

	client = &b32c{
		addr:       serverAddr,
		encChan:    make(chan encReq, 1),
		encRes:     make(chan *proto.EncodeResponse),
		decChan:    make(chan decReq, 1),
		decRes:     make(chan *proto.DecodeResponse),
		timeout:    timeout,
		waitingEnc: make(map[time.Time]encReq),
		waitingDec: make(map[time.Time]decReq),
	}

	return
}

// Start up the client. Call b.cancel() to stop.
//
// The returned send and recv functions are async by default
// but can be used synchronously by receiving from them directly:
//
//     encRes := <-enc(req)
//     decRes := <-dec(req)
func (b *b32c) Start() (
	enc func(*proto.EncodeRequest) chan *proto.EncodeResponse,
	dec func(*proto.DecodeRequest) chan *proto.DecodeResponse,
	stop func(),
	err error,
) {

	// Dial the configured server address
	clientConn, err := grpc.Dial(
		b.addr,
		grpc.WithTransportCredentials(insecure.NewCredentials()),
	)
	if err != nil {
		return
	}

	cli := proto.NewTranscriberClient(clientConn)
	ctx, cancelFunc := context.WithCancel(context.Background())
	b.stop = ctx.Done()

	var encode proto.Transcriber_EncodeClient
	encode, err = cli.Encode(ctx)
	if err != nil {
		cancelFunc()
		return
	}
	var decode proto.Transcriber_DecodeClient
	decode, err = cli.Decode(ctx)
	if err != nil {
		cancelFunc()
		return
	}

	go func() {

		err := b.Decode(decode)
		if err != nil {
			log.Print(err)
		}
	}()

	go func() {

		err := b.Encode(encode)
		if err != nil {
			log.Print(err)
		}
	}()

	enc = func(req *proto.EncodeRequest) chan *proto.EncodeResponse {
		r := newEncReq(req)
		b.encChan <- r
		return r.Res
	}
	dec = func(req *proto.DecodeRequest) chan *proto.DecodeResponse {
		r := newDecReq(req)
		b.decChan <- r
		return r.Res
	}
	stop = cancelFunc
	return
}

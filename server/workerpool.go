package server

import (
	"github.com/quanterall/kitchensink/based32"
	protos "github.com/quanterall/kitchensink/proto"
)

type Transcriber struct {
	stop   chan struct{}
	encode chan *protos.EncodeRequest
	decode chan *protos.DecodeRequest
}

func NewWorkerPool(stop chan struct{}) *Transcriber {

	t := &Transcriber{
		stop:   stop,
		encode: make(chan *protos.EncodeRequest),
		decode: make(chan *protos.DecodeRequest),
	}

	return t
}

func (t *Transcriber) Start() (stop func()) {

	go t.handle()
	return func() { close(t.stop) }
}

func (t *Transcriber) handle() {
out:
	for {
		select {
		case msg := <-t.encode:
			based32.Codec.Encode(msg.Data)
		case msg := <-t.decode:
			based32.Codec.Decode(msg.EncodedString)
		case <-t.stop:
			break out
		}
	}

}

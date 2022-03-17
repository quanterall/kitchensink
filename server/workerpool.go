package server

import (
	"github.com/quanterall/kitchensink/based32"
	protos "github.com/quanterall/kitchensink/proto"
	"go.uber.org/atomic"
)

type Transcriber struct {
	stop                       chan struct{}
	encode                     chan *protos.EncodeRequest
	decode                     chan *protos.DecodeRequest
	encCallCount, decCallCount *atomic.Uint32
	workers                    int
}

func NewWorkerPool(workers int, stop chan struct{}) *Transcriber {

	t := &Transcriber{
		stop:         stop,
		encode:       make(chan *protos.EncodeRequest),
		decode:       make(chan *protos.DecodeRequest),
		encCallCount: atomic.NewUint32(0),
		decCallCount: atomic.NewUint32(0),
		workers:      workers,
	}

	return t
}

func (t *Transcriber) Start() (stop func()) {
	for i := 0; i < t.workers; i++ {
		go t.handle()
	}
	return func() { close(t.stop) }
}

func (t *Transcriber) handle() {
out:
	for {
		select {
		case msg := <-t.encode:
			t.encCallCount.Inc()
			based32.Codec.Encode(msg.Data)
		case msg := <-t.decode:
			t.decCallCount.Inc()
			based32.Codec.Decode(msg.EncodedString)
		case <-t.stop:
			break out
		}
	}
	t.LogCallCounts()
}

func (t *Transcriber) LogCallCounts() {
	log.Printf("processed %v encodes and %v encodes",
		t.encCallCount.Load(), t.decCallCount.Load(),
	)
}

package server

import (
	"github.com/quanterall/kitchensink/pkg/based32"
	"github.com/quanterall/kitchensink/pkg/proto"
	"go.uber.org/atomic"
	"sync"
)

// transcriber is a multithreaded worker pool for performing transcription encode
// and decode requests. It is not exported because it must be initialised
// correctly.
type transcriber struct {
	stop                       chan struct{}
	encode                     []chan *proto.EncodeRequest
	decode                     []chan *proto.DecodeRequest
	encodeRes                  []chan proto.EncodeRes
	decodeRes                  []chan proto.DecodeRes
	encCallCount, decCallCount *atomic.Uint32
	workers                    uint32
	wait                       sync.WaitGroup
}

// NewWorkerPool initialises the data structure required to run a worker pool.
// Call Start to to initiate the run, and call the returned stop function to end
// it.
func NewWorkerPool(workers uint32, stop chan struct{}) *transcriber {

	// Initialize a transcriber worker pool
	t := &transcriber{
		stop:         stop,
		encode:       make([]chan *proto.EncodeRequest, workers),
		decode:       make([]chan *proto.DecodeRequest, workers),
		encodeRes:    make([]chan proto.EncodeRes, workers),
		decodeRes:    make([]chan proto.DecodeRes, workers),
		encCallCount: atomic.NewUint32(0),
		decCallCount: atomic.NewUint32(0),
		workers:      workers,
		wait:         sync.WaitGroup{},
	}

	// Create a channel for each worker to send and receive on, buffer them by
	// the same as the number of workers, producing workers^2 job slots
	for i := uint32(0); i < workers; i++ {
		t.encode[i] = make(chan *proto.EncodeRequest)
		t.decode[i] = make(chan *proto.DecodeRequest)
		t.encodeRes[i] = make(chan proto.EncodeRes)
		t.decodeRes[i] = make(chan proto.DecodeRes)
	}

	return t
}

// handle the jobs, this is one thread of execution, and will run whatever job
// has appeared and this thread
func (t *transcriber) handle(worker uint32) {

	t.wait.Add(1)
out:
	for {
		select {
		case msg := <-t.encode[worker]:

			t.encCallCount.Inc()
			res, err := based32.Codec.Encode(msg.Data)
			t.encodeRes[worker] <- proto.EncodeRes{
				IdNonce: msg.IdNonce,
				String:  res,
				Error:   err,
			}

		case msg := <-t.decode[worker]:

			t.decCallCount.Inc()

			bytes, err := based32.Codec.Decode(msg.EncodedString)
			t.decodeRes[worker] <- proto.DecodeRes{
				IdNonce: msg.IdNonce,
				Bytes:   bytes,
				Error:   err,
			}

		case <-t.stop:

			break out
		}
	}

	t.wait.Done()

}

// logCallCounts prints the values stored in the encode and decode counter
// atomic variables.
func (t *transcriber) logCallCounts() {

	log.Printf(
		"processed %v encodes and %v encodes",
		t.encCallCount.Load(), t.decCallCount.Load(),
	)
}

// Start up the worker pool.
func (t *transcriber) Start() (cleanup func()) {

	// Spawn the number of workers configured.
	for i := uint32(0); i < t.workers; i++ {

		go t.handle(i)
	}

	return func() {

		// Wait until all have stopped.
		t.wait.Wait()

		// Log the number of jobs that were done during the run.
		t.logCallCounts()

	}
}

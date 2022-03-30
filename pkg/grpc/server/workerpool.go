package server

import (
	codec "github.com/quanterall/kitchensink"
	"github.com/quanterall/kitchensink/pkg/based32"
	"github.com/quanterall/kitchensink/pkg/proto"
	"go.uber.org/atomic"
	"sync"
)

// Transcriber is a multithreaded worker pool for performing transcription
// encode and decode requests.
type Transcriber struct {
	stop                       chan struct{}
	encode                     []chan *protos.EncodeRequest
	decode                     []chan *protos.DecodeRequest
	encodeRes                  []chan string
	decodeRes                  []chan codec.DecodeRes
	encCallCount, decCallCount *atomic.Uint32
	workers                    uint32
	wait                       sync.WaitGroup
}

// NewWorkerPool initialises the data structure required to run a worker pool.
// Call Start to to initiate the run, and call the returned stop function to end
// it.
func NewWorkerPool(workers uint32, stop chan struct{}) *Transcriber {

	// Initialize a Transcriber worker pool
	t := &Transcriber{
		stop:         stop,
		encode:       make([]chan *protos.EncodeRequest, workers),
		decode:       make([]chan *protos.DecodeRequest, workers),
		encodeRes:    make([]chan string, workers),
		decodeRes:    make([]chan codec.DecodeRes, workers),
		encCallCount: atomic.NewUint32(0),
		decCallCount: atomic.NewUint32(0),
		workers:      workers,
		wait:         sync.WaitGroup{},
	}

	// Create a channel for each worker to send and receive on
	for i := uint32(0); i < workers; i++ {
		t.encode[i] = make(chan *protos.EncodeRequest)
		t.decode[i] = make(chan *protos.DecodeRequest)
		t.encodeRes[i] = make(chan string)
		t.decodeRes[i] = make(chan codec.DecodeRes)
	}

	return t
}

// Start up the worker pool.
func (t *Transcriber) Start() (stop func()) {

	// Spawn the number of workers configured.
	for i := uint32(0); i < t.workers; i++ {

		go t.handle(i)
	}

	return func() {

		// Close the stop channel to signal all workers to break out of their
		// loop.
		close(t.stop)

		// Wait until all have stopped.
		t.wait.Wait()

		// Log the number of jobs that were done during the run.
		t.logCallCounts()

	}
}

// handle the jobs, this is one thread of execution, and will run whatever job
// has appeared and this thread
func (t *Transcriber) handle(worker uint32) {

	t.wait.Add(1)
out:
	for {
		select {
		case msg := <-t.encode[worker]:

			t.encCallCount.Inc()

			t.encodeRes[worker] <- based32.Codec.Encode(msg.Data)

		case msg := <-t.decode[worker]:

			t.decCallCount.Inc()

			decoded, bytes := based32.Codec.Decode(msg.EncodedString)
			t.decodeRes[worker] <- codec.DecodeRes{
				Decoded: decoded,
				Data:    bytes,
			}

		case <-t.stop:

			break out
		}
	}

	t.wait.Done()

}

// logCallCounts prints the values stored in the encode and decode counter
// atomic variables.
func (t *Transcriber) logCallCounts() {

	log.Printf("processed %v encodes and %v encodes",
		t.encCallCount.Load(), t.decCallCount.Load(),
	)
}

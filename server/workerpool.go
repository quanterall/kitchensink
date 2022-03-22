package server

import (
	"github.com/quanterall/kitchensink/based32"
	protos "github.com/quanterall/kitchensink/proto"
	"go.uber.org/atomic"
	"sync"
)

// Transcriber is a multithreaded worker pool for performing transcription
// encode and decode requests.
type Transcriber struct {
	stop                       chan struct{}
	encode                     chan *protos.EncodeRequest
	decode                     chan *protos.DecodeRequest
	encCallCount, decCallCount *atomic.Uint32
	workers                    int
	wait                       sync.WaitGroup
}

// NewWorkerPool initialises the data structure required to run a worker pool.
// Call Start to to initiate the run, and call the returned stop function to end
// it.
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

// Start up the worker pool.
func (t *Transcriber) Start() (stop func()) {

	// Spawn the number of workers configured.
	for i := 0; i < t.workers; i++ {

		go t.handle()
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
func (t *Transcriber) handle() {

	t.wait.Add(1)
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

	t.wait.Done()

}

// logCallCounts prints the values stored in the encode and decode counter
// atomic variables.
func (t *Transcriber) logCallCounts() {

	log.Printf("processed %v encodes and %v encodes",
		t.encCallCount.Load(), t.decCallCount.Load(),
	)
}

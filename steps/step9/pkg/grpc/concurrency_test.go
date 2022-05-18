package grpc

import (
	"encoding/binary"
	"github.com/quanterall/kitchensink/pkg/grpc/client"
	"github.com/quanterall/kitchensink/pkg/grpc/server"
	"github.com/quanterall/kitchensink/pkg/proto"
	"go.uber.org/atomic"
	"lukechampine.com/blake3"
	"math/rand"
	"net"
	"sync"
	"testing"
	"time"
)

const (
	maxLength = 4096
	minLength = 8
	randN     = maxLength - minLength
	hashLen   = 32
	testItems = 64
)

type SequencedString struct {
	seq int
	str string
}

type SequencedBytes struct {
	seq int
	byt []byte
}

// TestGRPCCodecConcurrency deliberately intersperses extremely long messages
// and spawns tests concurrently in order to ensure the client correctly returns
// responses to the thread that requested them
func TestGRPCCodecConcurrency(t *testing.T) {

	// For reasons of producing wide variance, we will generate the source
	// material using hash chains
	seedBytes := make([]byte, 8)
	binary.LittleEndian.PutUint64(seedBytes, seed)

	lastHash := make([]byte, hashLen)

	out := blake3.Sum256(seedBytes)
	copy(lastHash, out[:])

	generated := make([][]byte, testItems)

	// Next, we will use the same fixed seed to define the variable lengths
	// between 8 and 4096 bytes for 64 items
	rand.Seed(seed)
	for i := 0; i < testItems; i++ {

		// Every second slice will be a lot smaller so it will process and
		// return while its predecessor is still in process
		var length int
		if i%2 != 0 {

			length = rand.Intn(hashLen-1) + 1
		} else {

			length = rand.Intn(randN)
		}

		// calculate the divisor and modulus of length compared to hash length
		cycles, cycleMod := length/hashLen, length%hashLen
		if cycleMod != 0 {

			// to make a hash chain long enough we need to add one and then trim
			// the result back
			cycles++
		}

		thisHash := make([]byte, cycles*hashLen)

		for j := 0; j < cycles; j++ {

			// hash the last hash
			out = blake3.Sum256(lastHash)

			// copy result back save last hash
			copy(lastHash, out[:])

			// copy last hash to position in output
			copy(thisHash[j*hashLen:(j+1)*hashLen], lastHash)
		}

		// trim result to random generated length value
		generated[i] = thisHash[:length]
	}

	// Set up a server
	addr, err := net.ResolveTCPAddr("tcp", defaultAddr)
	if err != nil {
		t.Fatal(err)
	}
	srvr := server.New(addr, 8)
	stopSrvr := srvr.Start()

	// Create a client
	cli, err := client.New(defaultAddr, time.Second*5)
	if err != nil {
		t.Fatal(err)
	}
	enc, dec, stopCli := cli.Start()

	// To create a collection that can be sorted easily after creation back into
	// an ordered slice, we create a buffered channel with enough buffers to
	// hold all of the items we will feed into it
	stringChan := make(chan SequencedString, testItems)

	// We will use this to make sure every request completes before we shut down
	// the client and server
	var wg sync.WaitGroup

	// We will keep track of ins and outs for log prints
	var qCount atomic.Uint32

	log.Println("encoding received items")

	for i := range generated {

		go func(i int) {

			log.Println("encode processing item", i)

			// we need to wait until all messages process before collating the
			// results
			wg.Add(1)
			qCount.Inc()

			log.Println(
				"encode request", i, "sending,",
				qCount.Load(), "items in queue",
			)

			// send out the query and wait for the response
			encRes := <-enc(
				&proto.EncodeRequest{
					Data: generated[i],
				},
			)

			// push the returned result into our channel buffer with the item
			// sequence number, so it can be reordered correctly
			stringChan <- SequencedString{
				seq: i,
				str: encRes.GetEncodedString(),
			}

			wg.Done()
			qCount.Dec()

			log.Println(
				"encode request", i, "response received,",
				qCount.Load(), "items in queue",
			)

		}(i)
	}

	// Wait until all results are back so we can assemble them in order for
	// checking
	wg.Wait()

	encoded := make([]string, testItems)

	counter := 0

	for item := range stringChan {

		counter++
		log.Println("collating encode item", item.seq, "items done:", counter)
		// place items back in the sequence position they were created in
		encoded[item.seq] = item.str
		if counter >= testItems {
			break
		}
	}

	// To create a collection that can be sorted easily after creation back into
	// an ordered slice, we create a buffered channel with enough buffers to
	// hold all of the items we will feed into it
	bytesChan := make(chan SequencedBytes, testItems)

	log.Println("decoding received items")

	for i := range encoded {

		go func(i int) {

			log.Println("decode processing item", i)

			// we need to wait until all messages process before collating the
			// results
			wg.Add(1)
			qCount.Inc()

			log.Println(
				"decode request", i, "sending,",
				qCount.Load(), "items in queue",
			)

			// send out the query and wait for the response
			decRes := <-dec(
				&proto.DecodeRequest{
					EncodedString: encoded[i],
				},
			)

			// push the returned result into our channel buffer with the item
			// sequence number, so it can be reordered correctly
			bytesChan <- SequencedBytes{
				seq: i,
				byt: decRes.GetData(),
			}

			wg.Done()
			qCount.Dec()

			log.Println(
				"decode request", i, "response received,",
				qCount.Load(), "items in queue",
			)

		}(i)
	}

	// Wait until all results are back so we can assemble them in order for
	// checking
	wg.Wait()

	decoded := make([][]byte, testItems)

	counter = 0

	for item := range bytesChan {

		counter++
		log.Println("collating decode item", item.seq, "items done:", counter)
		// place items back in the sequence position they were created in
		decoded[item.seq] = item.byt
		if counter >= testItems {
			break
		}
	}

	// now, to compare the inputs to the processed
	for i := range generated {

		if string(generated[i]) != string(decoded[i]) {

			t.Fatal("decoded item", i, "not the same as generated")
		}
	}

	log.Println("shutting down client and server")
	stopCli()
	stopSrvr()
}

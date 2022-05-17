package grpc

import (
	"encoding/binary"
	"fmt"
	"lukechampine.com/blake3"
	"math/rand"
	"testing"
)

const (
	maxLength = 4096
	minLength = 8
	randN     = maxLength - minLength
	hashLen   = 32
)

// TestGRPCCodecConcurrency deliberately intersperses extremely long messages
// and spawns tests concurrently in order to ensure the client correctly returns
// responses to the thread that requested them
func TestGRPCCodecConcurrency(t *testing.T) {

	// For reasons of producing wide variance, we will generate the source material
	// using hash chains
	seedBytes := make([]byte, 8)
	binary.LittleEndian.PutUint64(seedBytes, seed)

	lastHash := make([]byte, 32)

	out := blake3.Sum256(seedBytes)
	copy(lastHash, out[:])

	generated := make([][]byte, 64)

	// Next, we will use the same fixed seed to define the variable lengths between 8
	// and 4096 bytes for 64 items
	rand.Seed(seed)
	for i := 0; i < 64; i++ {

		// Every second slice will be a lot smaller so it will process and return while
		// its predecessor is still in process
		var length int
		if i%2 != 0 {

			length = rand.Intn(hashLen-1) + 1
		} else {

			length = rand.Intn(randN)
		}

		// calculate the divisor and modulus of length compared to hash length
		cycles, cycleMod := length/hashLen, length%hashLen
		if cycleMod != 0 {

			// to make a hash chain long enough we need to add one and then trim the result
			// back
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

	var o string
	for i := range generated {
		o += fmt.Sprint(len(generated[i]), ", ")
	}
	log.Println(o)
}

package grpc

import (
	"encoding/hex"
)

// MakeLongMessage creates a message by concatenating a collection of hex
// encoded strings into one, and then repeating the concatenation with the
// resultant bytes for the number of given repetitions
func MakeLongMessage(repetitions int, input []string) (output []byte) {
	totalLen := 0
	// Accumulate the length of all the strings, each character of the string
	// represents 4 bytes so the byte length is half of this
	for i := range input {
		totalLen += len(input[i]) / 2
	}
	// Make one repetition as a single slice
	oneRep := make([]byte, totalLen)
	cursor := 0
	for i := range input {
		bytes, err := hex.DecodeString(input[i])
		if err != nil {
			panic("something wrong with input")
		}
		bLen := len(bytes)
		copy(oneRep[cursor:cursor+bLen], bytes)
		cursor += bLen
	}
	// Copy the single repetition over and over onto a buffer of this size times
	// the number of repetitions requested
	output = make([]byte, totalLen*repetitions)
	cursor = 0
	for reps := 0; reps < repetitions; reps++ {
		copy(output[cursor:cursor+totalLen], oneRep)
		cursor += totalLen
	}
	return
}

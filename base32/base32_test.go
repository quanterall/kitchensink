package base32

import (
	"encoding/binary"
	"encoding/hex"
	"lukechampine.com/blake3"
	"math/rand"
	"testing"
)

func TestCodec(t *testing.T) {
	// generate 10 pseudorandom 64 bit values
	rand.Seed(1234567890)
	seeds := make([]uint64, 10)
	for i := range seeds {
		seeds[i] = rand.Uint64()
	}

	// convert to bytes
	seedBytes := make([][]byte, 10)
	for i := range seedBytes {
		seedBytes[i] = make([]byte, 8)
		binary.LittleEndian.PutUint64(seedBytes[i], seeds[i])
	}

	// generate hashes from the seeds
	hashedSeeds := make([][]byte, 10)
	for i := range hashedSeeds {
		hashed := blake3.Sum256(seedBytes[i])
		hashedSeeds[i] = hashed[:]
	}

	expected := []string{
		"7bf4667ea06fe57687a7c0c8aae869db103745a3d8c5dce5bf2fc6206d3b97e4",
		"84c0ee2f49bfb26f48a78c9d048bb309a006db4c7991ebe4dd6dc3f2cc1067cd",
		"206a953c4ba4f79ffe3d3a452214f19cb63e2895866cc27c7cf6a4ec8fe5a7a6",
		"35d64c401829c621624fe9d4f134c24ae909ecf4f07ec4540ffd58911f427d03",
		"573d6989a2c2994447b4669ae6931f12e73c8744e9f65451918a1f3d8cd39aa1",
		"2b08aea58cc1d680de0e7acadc027ebe601f923ff9d5536c6f73e2559a1b6b14",
		"bcc3256005da59b06f69b4c1cc62c89af041f8cd5ad79b81351fbfbbaf2cc60f",
		"42a0f7b9aef1cdc0b3f2a1fd0fb547fb76e5eb50f4f5a6646ee8929fdfef5db7",
		"50e1cb9f5f8d5325e18298faeeea7fd93d83e3bd518299e7150c1f548c11ddc8",
		"22a70a74ccfd61a47576150f968039cfeb33143ec549dfeb6c95afc8a6d3d75a",
	}
	for i := range hashedSeeds {
		if expected[i] != hex.EncodeToString(hashedSeeds[i]) {
			t.Log("failed", i, "expected", expected[1], "found", hashedSeeds)
			t.FailNow()
		}
	}
}

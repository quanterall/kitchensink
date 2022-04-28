package based32

import (
	"encoding/binary"
	"encoding/hex"
	"fmt"
	"lukechampine.com/blake3"
	"math/rand"
	"testing"
)

const (
	seed    = 1234567890
	numKeys = 32
)

func TestCodec(t *testing.T) {

	// Generate 10 pseudorandom 64 bit values. We do this here rather than
	// pre-generating this separately as ultimately it is the same thing, the
	// same seed produces the same series of pseudorandom values, and the hashes
	// of these values are deterministic.
	rand.Seed(seed)
	seeds := make([]uint64, numKeys)
	for i := range seeds {

		seeds[i] = rand.Uint64()
	}

	// Convert the uint64 values to 8 byte long slices for the hash function.
	seedBytes := make([][]byte, numKeys)
	for i := range seedBytes {

		seedBytes[i] = make([]byte, 8)
		binary.LittleEndian.PutUint64(seedBytes[i], seeds[i])
	}

	// Generate hashes from the seeds
	hashedSeeds := make([][]byte, numKeys)

	// Uncomment lines relating to this variable to regenerate expected data
	// that will log to console during test run.
	generated := "\nexpected := []string{\n"

	for i := range hashedSeeds {

		hashed := blake3.Sum256(seedBytes[i])
		hashedSeeds[i] = hashed[:]

		generated += fmt.Sprintf("\t\"%x\",\n", hashedSeeds[i])
	}

	generated += "}\n"
	t.Log(generated)

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
		"cd1b21bd745e122f0db1f5ca4e4cbe0bace8439112d519e5b9c0a44a2648a61a",
		"9ec4bd670b053e722d9d5bc3c2aca4a1d64858ef53c9b58a9222081dc1eeb017",
		"d85713c898f2fc95d282b58a0ea475c1b1726f8be44d2f17c3c675ce688e1563",
		"875baee7e9a372fe3bad1fbf0e2e4119038c1ed53302758fc5164b012e927766",
		"7de94ca668463db890478d8ba3bbed35eb7666ac0b8b4e2cb808c1cbb754576f",
		"159469150dc41ebd2c2bfafb84aef221769699013b70068444296169dc9e93be",
		"a90a104ea470df61d337c51589b520454acbd05ef5bbe7d2a8285043a222bec9",
		"a835de5206f6dbef6a2cb3da66ffb99a19bfa4e005208ffdb316ce880132297e",
		"f6a09e8f41231febd1b25c52cb73ea438ac803db77d5549db4e15a32e804de9f",
		"074c59cce7783042cc6941c849206582ecc43028d1576d00e02d95b1e669bf7a",
		"203c3566724c229b570f33be994cd6094e1a64f3df552f1390b4c2adc7e36d6d",
		"efec32d52a17ed75ad5a486ba621e0f47f61e4e60557129fce728a1bb63208fd",
		"9cc2962fc62fe40f6197a4fb81356717fd57b4c988641bca3a9d45efde893894",
		"2adf211300632bb5f650202bf128ba9187ec2c6c738431dc396d93b8f62bd590",
		"0782aade40d0ae7a293bfb67016466682d858b5226eaaa8df2f2104fa6c408c3",
		"d011ad5550f3f03caa469fa233f553721e6af84f1341d256cefe052d85397637",
		"83deb64f5c134d108e8b99c8a196b8d04228acfc810c33711d975400fa731508",
		"d9a4b19142d015fd541f50f18f41b7e9738a30c59a3e914b4d4d1556c75786f2",
		"3e05940b76735ea114db8b037dece53090765510c9c4e55a0be18cb8aef754fa",
		"41f43119041dd1f3a250f54768ce904808cd0d7bb7b37697803ed2940c39a555",
		"a2c2d7cb980c2b57c8fdfae55cf4c6040eaf8163b21072877e5e57349388d59c",
		"02155c589e5bd89ce806a33c1841fe1e157171222701d515263acd0254208a39",
	}

	for i := range hashedSeeds {

		if expected[i] != hex.EncodeToString(hashedSeeds[i]) {

			t.Log("failed", i, "expected", expected[1], "found", hashedSeeds)
			t.FailNow()
		}
	}
}

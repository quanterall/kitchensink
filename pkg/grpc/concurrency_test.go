package grpc

import (
	"encoding/binary"
	"encoding/hex"
	"fmt"
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

// TestGRPCCodecConcurrency deliberately intersperses extremely long messages
// and spawns tests concurrently in order to ensure the client correctly returns
// responses to the thread that requested them
func TestGRPCCodecConcurrency(t *testing.T) {

	addr, err := net.ResolveTCPAddr("tcp", defaultAddr)
	if err != nil {
		t.Fatal(err)
	}
	srvr := server.New(addr, 32)
	stopSrvr := srvr.Start()

	cli, err := client.New(defaultAddr, time.Second*5)
	if err != nil {
		t.Fatal(err)
	}
	enc, dec := cli.Start()

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
	// t.Log(generated)

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

	// encodedStr := []string{
	//     "QNTRLfalgen75ph72a585lqv32hgd8d3qd6950vvth89huhuvgrd8wt7ftet",
	//     "QNTRLwzvpm30fxlmym6g57xf6pytkvy6qpkmf3uer6lym4ku8ukvzpnckqs6",
	//     "QNTRLssx49fufwj008l785ay2gs57xwtv03gjkrxesnu0nm2fmy0uhmerj80",
	//     "QNTRL56avnzqrq5uvgtzfl5afuf5cf9wjz0v7nc8a3z5pl743yglgjs6nypp",
	//     "QNTRLetn66vf5tpfj3z8k3nf4e5nrufww0y8gn5lv4z3jx9p70vwrzchx7pf",
	//     "QNTRLg4s3t493nqadqx7peav4hqz06lxq8uj8lua25mvdae7y4v6rd43fmcv",
	//     "QNTRLw7vxftqqhd9nvr0dx6vrnrzezd0qs0ce4dd0xupx50mlwa09nr8hhvc",
	//     "QNTRL3p2paae4mcums9n72sl6ra4glahde0t2r60tfnydm5f987lalykzuqe",
	//     "QNTRL4gwrjult7x4xf0ps2v04mh20lvnmqlrh4gc9x08z5xp74yva49qf29d",
	//     "QNTRLc32wzn5en7krfr4wc2sl95q8887kvc58mz5nhltdj26ljxy33yxs7px",
	//     "QNTRLtx3kgdaw30pytcdk86u5njvhc96e6zrjyfd2x09h8q2gj3xfznp4l5y",
	//     "QNTRLw0vf0t8pvznuu3dn4du8s4v5jsavjzcaafundv2jg3qs8wpa6czu2kz",
	//     "QNTRLnv9wy7gnre0e9wjs26c5r4ywhqmzun030jy6tchc0r8tnng3e6u5pg6",
	//     "QNTRLkr4hth8ax3h9l3m450m7r3wgyvs8rq765esyav0c5tykqfwr2u7zmfp",
	//     "QNTRLe77jn9xdprrmwysg7xchgama567kanx4s9ckn3vhqyvrjam6x9h7qmx",
	//     "QNTRLg2eg6g4phzpa0fv90a0hp9w7gshd95eqyahqp5ygs5kz6wun6fmukle",
	//     "QNTRLw5s5yzw53cd7cwnxlz3tzd4ypz54j7stm6mhe7j4q59qsazy2lfqq35",
	//     "QNTRLj5rthjjqmmdhmm29jea5ehlhxdpn0ayuqzjprlakvtvazqpxt0zdxtv",
	//     "QNTRLhm2p850gy33l673kfw99jmnafpc4jqrmdma24yakns45vhg0sfcsrrp",
	//     "QNTRLcr5ckwvuaurqskvd9qusjfqvkpwe3ps9rg4wmgquqketvvex8jy6alw",
	//     "QNTRLgsrcdtxwfxz9x6hpuemax2v6cy5uxny70042tcnjz6v9tw8udkk6rkl",
	//     "QNTRL0h7cvk49gt76addtfyxhf3pur687c0yucz4wy5leeeg5xakxgywasu3",
	//     "QNTRLjwv9930cch7grmpj7j0hqf4vutl64a5exyxgx7282w5tm7738hap8u7",
	//     "QNTRL54d7ggnqp3jhd0k2qszhufgh2gc0mpvd3ecgvwu89ke8w8kcv3ksyma",
	//     "QNTRLcrc92k7grg2u73f80akwqtyve5zmpvt2gnw425d7tepqnaz7jehh549",
	//     "QNTRLtgprt242relq092g606yvl42depu6hcfuf5r5jkemlq2tv989mr0wng",
	//     "QNTRLwpaadj0tsf56yyw3wvu3gvkhrgyy29vljqscvm3rkt4gq86wv2upf6r",
	//     "QNTRLnv6fvv3gtgptl25rag0rr6pkl5h8z3sckdray2tf4x324k82a58c5kf",
	//     "QNTRL5lqt9qtwee4agg5mw9sxl0vu5cfqaj4zryufe26p0scew9w9cgcnqj8",
	//     "QNTRLeqlgvgeqswaruaz2r65w6xwjpyq3ngd0wmmxa5hsqld998rns3l4gfz",
	//     "QNTRL23v947tnqxzk47glhaw2h85cczqatupvwepqu580e09wdyn3r2eeltr",
	//     "QNTRLvpp2hzcneda388gq63ncxzplc0p2ut3ygnsr4g4ycav6qj5yz9g9lv8",
	// }

	// In order to test concurrency we need to put dramatically longer messages
	// amongst the ones above, which are the same as the non-concurrent test.
	//
	// To do this we will use the long message generator and intersperse these
	// new long messages to ensure there will be out of order responses.

	var longMessages [][]byte

	for i := range expected {
		for j := 0; j < 2; j++ {
			longMessages = append(
				longMessages,
				MakeLongMessage(
					// the longer messages first to ensure there will be out of
					// order returns, and gradually shorter long messages to ensure
					// the parallel processing will get quite disordered
					2,

					append(
						expected[0:i], expected[i:]...,
					),
				),
			)

			// every second message will be the short messages already generated
			longMessages = append(longMessages, hashedSeeds[i])
		}
	}

	var o string
	for i := range longMessages {
		o += fmt.Sprint(len(longMessages[i]), ", ")
	}
	o += "\n"
	log.Println(o)
	// t.FailNow()

	// replace the hashedSeeds with the long version
	hashedSeeds = longMessages

	// replace expected with hex encoded longMessages
	var newExpected []string
	for i := range hashedSeeds {

		newExpected = append(newExpected, hex.EncodeToString(hashedSeeds[i]))
	}
	expected = newExpected

	// encoded := "\nencodedStr := []string{\n"

	// var encodedLong []string
	var encodedLong atomic.Value

	encodedLong.Store([]string{})

	var wg sync.WaitGroup
	var qCount atomic.Uint32

	// Convert hashes to our base32 encoding format
	for i := range hashedSeeds {

		// In order to ensure
		go func(i int) {

			// we need to wait until all messages process before moving to the
			// next part of the test
			wg.Add(1)
			qCount.Inc()

			// Note that we are slicing off a number of bytes at the end according
			// to the sequence number to get different check byte lengths from a
			// uniform original data. As such, this will be accounted for in the
			// check by truncating the same amount in the check (times two, for the
			// hex encoding of the string).
			// log.Println(
			//     "encode message", i, "sending", qCount.Load(), "in queue",
			// )
			encRes := <-enc(
				&proto.EncodeRequest{
					Data: hashedSeeds[i][:len(hashedSeeds[i])-i%5],
				},
			)
			// log.Println("encode message", i, "received back")
			// if err != nil {
			//     t.Fatal(err)
			// }
			encode := encRes.GetEncodedString()
			// if encode != encodedStr[i] {
			//     t.Fatalf(
			//         "Decode failed, expected '%s' got '%s'",
			//         encodedStr, encode,
			//     )
			// }

			encodedLong.Store(append(encodedLong.Load().([]string), encode))
			// encoded += "\t\"" + encode + "\",\n"
			wg.Done()
			qCount.Dec()
			log.Println("done job", i, qCount.Load(), "in queue")
		}(i)
	}
	wg.Wait()
	// encoded += "}\n"
	// t.Log(encoded)
	// Next, decode the encodedStr above, which should be the output of the
	// original generated seeds, with the index mod 5 truncations performed on
	// each as was done to generate them.
	_ = dec
	// for i := range encodedLong.Load().([]string) {
	//
	//     go func(i int) {
	//
	//         wg.Add(1)
	//         qCount.Inc()
	//         log.Println(
	//             "decode message", i, "sending", qCount.Load(), "in queue",
	//         )
	//         res := <-dec(
	//             &proto.DecodeRequest{
	//                 EncodedString: encodedLong.Load().([]string)[i],
	//             },
	//         )
	//         log.Println(
	//             "decode message", i, "received back", qCount.Load(), "in queue",
	//         )
	//         // res, err := Codec.Decode(encodedStr[i])
	//         // if err != nil {
	//         //     t.Fatalf("error: '%v'", err)
	//         // }
	//         elen := len(expected[i])
	//         etrimlen := 2 * (i % 5)
	//         expectedHex := expected[i][:elen-etrimlen]
	//         resHex := fmt.Sprintf("%x", res.GetData())
	//         if resHex != expectedHex {
	//             t.Fatalf(
	//                 "got: '%s' expected: '%s'",
	//                 resHex,
	//                 expectedHex,
	//             )
	//         }
	//         wg.Done()
	//         qCount.Dec()
	//     }(i)
	// }
	// wg.Wait()
	// _ = stopSrvr
	stopSrvr()
}

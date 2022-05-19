package main

import (
	"encoding/hex"
	"flag"
	"fmt"
	"github.com/quanterall/kitchensink/pkg/grpc/client"
	"github.com/quanterall/kitchensink/pkg/proto"
	"os"
	"time"
)

const defaultAddr = "localhost:50051"

var (
	serverAddr = flag.String(
		"a", defaultAddr,
		"The server address for basedcli to connect to",
	)
	encode = flag.String(
		"e", "",
		"hex string to convert to based32 encoding",
	)
	decode = flag.String(
		"d", "",
		"based32 encoded string to convert back to hex",
	)
)

func main() {

	flag.Parse()

	// if both or neither query fields have values it is an error
	noQuery := *encode == "" && *decode == ""
	bothQuery := *encode != "" && *decode != ""
	if noQuery || bothQuery {

		_, _ = fmt.Fprintln(
			os.Stderr,
			"basedcli - commandline client for based32 codec service",
		)

		if noQuery {

			_, _ = fmt.Fprintln(
				os.Stderr, "No query provided, printing help information",
			)

		} else {

			_, _ = fmt.Fprintln(
				os.Stderr, "Only one of -e or -d may be used",
			)

		}

		flag.PrintDefaults()
		os.Exit(1)

	}

	// Create a new client
	cli, err := client.New(defaultAddr, 5*time.Second)
	if err != nil {

		_, _ = fmt.Fprintln(os.Stderr, err)
		os.Exit(1)
	}

	// Start the client
	enc, dec, stopCli, err := cli.Start()
	if err != nil {

		_, _ = fmt.Fprintln(os.Stderr, err)
		os.Exit(1)
	}

	if err != nil {

		_, _ = fmt.Fprintln(os.Stderr, err)
		os.Exit(1)
	}

	if *encode != "" {

		// for encoding, hex decode errors are the only errors
		input, err := hex.DecodeString(*encode)
		if err != nil {

			_, _ = fmt.Fprintln(
				os.Stderr,
				"basedcli - commandline client for based32 codec service",
			)
			_, _ = fmt.Fprintln(os.Stderr, err)
			os.Exit(1)
		}

		// send encode request
		encRes := <-enc(
			&proto.EncodeRequest{
				Data: input,
			},
		)

		fmt.Println(encRes.GetEncodedString())

	} else if *decode != "" {

		decRes := <-dec(
			&proto.DecodeRequest{
				EncodedString: *decode,
			},
		)

		data := decRes.GetData()
		if data == nil {

			_, _ = fmt.Fprintln(
				os.Stderr,
				"basedcli - commandline client for based32 codec service",
			)
			_, _ = fmt.Fprintln(os.Stderr, "Error:", decRes.GetError())
			os.Exit(1)

		} else {

			fmt.Println(hex.EncodeToString(data))
		}
	}
	stopCli()
}

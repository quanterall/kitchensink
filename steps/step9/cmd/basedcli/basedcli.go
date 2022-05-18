package main

import (
	"flag"
	"fmt"
	"os"
)

const defaultAddr = "localhost:50051"

var (
	serverAddr = flag.String("a", defaultAddr,
		"The server address for basedcli to connect to",
	)
	encode = flag.String("e", "", "hex string to convert to based32 encoding")
	decode = flag.String("d", "",
		"based32 encoded string to convert back to hex",
	)
)

func main() {
	_, _ = fmt.Fprintln(os.Stderr,
		"basedcli - commandline client for based32 codec service",
	)
	flag.Parse()
	if *serverAddr == defaultAddr {
		_, _ = fmt.Fprintln(os.Stderr,
			"run with argument -h to print command line options",
		)
	}

}

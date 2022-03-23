package main

import (
	"flag"
	"fmt"
	"os"
)

const defaultAddr = "localhost:50051"

var serverAddr = flag.String("a", defaultAddr,
	"The server address for the basedd client to connect to",
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

package main

import (
	"flag"
	"fmt"
	"os"
)

var serverAddr = flag.String("a", "localhost:50051",
	"The server address for the basedd client to connect to",
)

func main() {
	_, _ = fmt.Fprintln(os.Stderr,
		"basedcli - commandline client for based32 codec service",
	)
	flag.Parse()
}

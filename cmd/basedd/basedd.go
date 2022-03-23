package main

import (
	"flag"
	"fmt"
	"os"
)

var serverAddr = flag.String("a", "localhost:50051",
	"The address to listen for connections in the format of host:port "+
		"- omit host to bind to all network interfaces",
)

func main() {
	_, _ = fmt.Fprintln(os.Stderr,
		"basedd - microservice for based32 human transcription encoding",
	)
	flag.Parse()
	_, _ = fmt.Fprintln(os.Stderr,
		"run with argument -h to print command line options",
	)
}

package main

import (
	"flag"
	"github.com/cybriq/interrupt"
	"github.com/cybriq/qu"
)

const defaultAddr = "localhost:50051"

var serverAddr = flag.String("a", defaultAddr,
	"The address to listen for connections in the format of host:port "+
		"- omit host to bind to all network interfaces",
)

var killAll = qu.T()

func main() {
	log.Println(
		"basedd - microservice for based32 human transcription encoding",
	)
	flag.Parse()
	if *serverAddr == defaultAddr {
		log.Println(
			"run with argument -h to print command line options",
		)
	}
	interrupt.AddHandler(func() {
		log.Println("Shutting down basedd microservice")
		killAll.Q()
	},
	)
	select {
	case <-killAll.Wait():
		break
	}
}

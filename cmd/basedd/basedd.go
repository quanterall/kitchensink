package main

import (
	"flag"
	"fmt"
	"github.com/cybriq/interrupt"
	"github.com/cybriq/qu"
	"github.com/quanterall/kitchensink/pkg/server"
	"net"
	"os"
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

	// If the the address is the same as default probably the user didn't set
	// one, so let them know they can as a courtesy.
	if *serverAddr == defaultAddr {
		log.Println(
			"run with argument -h to print command line options",
		)
	}

	addr, err := net.ResolveTCPAddr("tcp", *serverAddr)
	if err != nil {

		// If net.ParseIP returns nil it means the address is invalid.
		log.Printf("Failed to parse network address '%s'", *serverAddr)
		os.Exit(1)
	}

	log.Println(addr)

	svc := server.New(addr, 8)

	// interrupt is a library that allows the proper handling of OS interrupt
	// signals to allow a clean shutdown and ensure such things as databases are
	// properly closed and all pending writes are completed.
	interrupt.AddHandler(func() {

		// In this case, we are just ending the process, after the select block
		// below falls through when the channel is closed, the execution of the
		// application terminates.
		//
		// Note that Go applications keep running even if the main() has
		// terminated if goroutines are not terminated. So, quit channels are
		// fundamental to controlling most Go applications, and the bigger the
		// application the more threads there will be and the more crucial it is
		// that they are correctly terminated.
		//
		// Note that in many libraries the context library is used to provide
		// part of this functionality, but for general control, one still needs
		// to use breaker channels like this.
		//
		// Ultimately they are always implemented with exactly this pattern. The
		// qu library makes it easier to debug the channels when run control
		// bugs appear, you can print the information about the state of the
		// channels that are open and where in the code they are waiting.
		_, _ = fmt.Fprintln(os.Stderr)
		log.Println("Shutting down basedd microservice")
		killAll.Q()
	},
	)

	// In all cases, we create shutdown handlers and start receiving threads
	// before we start up the sending threads.
	stop := svc.Start()

	select {
	case <-killAll.Wait():

		// This triggers termination of the service. We separate the stop
		// controls of this application versus the services embedded inside the
		// server so that we can potentially instead *restart* the service
		// rather than only terminate it, in the case of a reconfiguration
		// signal. This signal is not handled, because this is too simple a
		// service and there is no configuration to really change. But this is
		// why you don't make one quit channel for an entire app, but instead
		// set them up in a cascade like this.
		stop()
		break
	}
}

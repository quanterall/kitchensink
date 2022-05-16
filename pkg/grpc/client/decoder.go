package client

import (
	"github.com/quanterall/kitchensink/pkg/proto"
	"io"
	"time"
)

func (b *b32c) Decode(stream proto.Transcriber_DecodeClient) (err error) {

	go func(stream proto.Transcriber_DecodeClient) {
	out:
		for {
			select {
			case <-b.stop:

				break out
			case msg := <-b.decChan:

				// log.Println("sending message on stream")
				err := stream.Send(msg.Req)
				if err != nil {
					log.Print(err)
				}
				b.waitingDec[time.Now()] = msg
			case recvd := <-b.decRes:

				for i := range b.waitingDec {

					// Return received responses
					if recvd.IdNonce == b.waitingDec[i].Req.IdNonce {

						// return response to client
						b.waitingDec[i].Res <- recvd

						// delete entry in pending job map
						delete(b.waitingDec, i)

						// if message is processed next section does not need to
						// be run as we have just deleted it
						continue
					}

					// Check for expired responses
					if i.Add(b.timeout).Before(time.Now()) {

						log.Println(
							"\nexpiring", i,
							"\ntimeout", b.timeout,
							"\nsince", i.Add(b.timeout),
							"\nis late", i.Add(b.timeout).Before(time.Now()),
						)
						delete(b.waitingDec, i)
					}
				}
			}
		}
	}(stream)

	go func(stream proto.Transcriber_DecodeClient) {
	in:
		for {
			// check whether it's shutdown time first
			select {
			case <-b.stop:
				break in
			default:
			}

			// Wait for and load in a newly received message
			recvd, err := stream.Recv()
			switch {
			case err == io.EOF:
				// The client has broken the connection, so we can quit
				log.Println("stream closed")
				break in
			case err != nil:

				log.Println(err)
				break
			}

			// forward received message to processing loop
			b.decRes <- recvd
		}
	}(stream)

	return
}

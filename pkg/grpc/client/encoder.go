package client

import (
	"github.com/quanterall/kitchensink/pkg/proto"
	"io"
	"time"
)

func (b *b32c) Encode(stream proto.Transcriber_EncodeClient) (err error) {

	go func(stream proto.Transcriber_EncodeClient) {
	out:
		for {
			select {
			case <-b.stop:

				break out
			case msg := <-b.encChan:

				// log.Println("sending message on stream")
				err := stream.Send(msg.Req)
				if err != nil {
					log.Print(err)
				}
				b.waitingEnc[time.Now()] = msg
			case recvd := <-b.encRes:

				for i := range b.waitingEnc {

					// Return received responses
					if recvd.IdNonce == b.waitingEnc[i].Req.IdNonce {

						// return response to client
						b.waitingEnc[i].Res <- recvd

						// delete entry in pending job map
						delete(b.waitingEnc, i)

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
						delete(b.waitingEnc, i)
					}
				}
			}
		}
	}(stream)

	go func(stream proto.Transcriber_EncodeClient) {
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
			b.encRes <- recvd
		}
	}(stream)

	return
}

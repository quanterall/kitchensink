package client

import (
	"github.com/davecgh/go-spew/spew"
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
				log.Println("sending message on stream")
				err := stream.Send(msg.Req)
				if err != nil {
					log.Print(err)
				}
				b.waitingDec[time.Now()] = msg
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
				break in
			case err != nil:

				// Any error is terminal here, so return it to the caller after
				// logging it
				log.Println(err)
				break
			}
			log.Print(spew.Sdump(recvd))
			for i := range b.waitingDec {

				log.Print(spew.Sdump(b.waitingEnc[i].Req))

				// Check for expired responses
				if i.Add(b.timeout).Before(time.Now()) {

					log.Println(
						"/nexpiring",
						i,
						"\ntimeout",
						b.timeout,
						"\nsince",
						i.Add(b.timeout),
						"\nis late",
						i.Add(b.timeout).Before(time.Now()),
						spew.Sdump(b.waitingDec[i]),
					)
					delete(b.waitingEnc, i)
				}

				// Return received responses
				if recvd.IdNonce == b.waitingDec[i].Req.IdNonce {
					b.waitingDec[i].Res <- recvd
					delete(b.waitingDec, i)
				}
			}
		}
	}(stream)

	return
}

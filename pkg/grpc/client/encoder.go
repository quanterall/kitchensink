package client

import (
	"github.com/davecgh/go-spew/spew"
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
				err := stream.Send(msg.Req)
				if err != nil {
					log.Print(err)
				}
				b.waitingEnc[time.Now()] = msg
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
				break in
			case err != nil:

				// Any error is terminal here, so return it to the caller after
				// logging it
				log.Println(err)
				break
			}
			log.Print(spew.Sdump(recvd))
			for i := range b.waitingEnc {

				log.Print(spew.Sdump(b.waitingEnc[i].Req))

				// Check for expired responses
				if i.Add(b.timeout).Before(time.Now()) {

					log.Println(
						"\nexpiring",
						i,
						"\ntimeout",
						b.timeout,
						"\nsince",
						i.Add(b.timeout),
						"\nis late",
						i.Add(b.timeout).Before(time.Now()),
						spew.Sdump(b.waitingEnc[i]),
					)
					delete(b.waitingEnc, i)
				}

				// Return received responses
				if recvd.IdNonce ==
					b.waitingEnc[i].Req.IdNonce {
					b.waitingEnc[i].Res <- recvd
					delete(b.waitingEnc, i)
				}
			}
		}
	}(stream)

	return
}

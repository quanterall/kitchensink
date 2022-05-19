package grpc

import (
	"encoding/hex"
	"github.com/quanterall/kitchensink/pkg/grpc/client"
	"github.com/quanterall/kitchensink/pkg/grpc/server"
	"github.com/quanterall/kitchensink/pkg/proto"
	"net"
	"testing"
	"time"
)

const defaultAddr = "localhost:50051"

func TestGRPC(t *testing.T) {
	addr, err := net.ResolveTCPAddr("tcp", defaultAddr)
	if err != nil {
		t.Fatal(err)
	}
	srvr := server.New(addr, 8)
	stopServer := srvr.Start()

	cli, err := client.New(defaultAddr, 5*time.Second)
	if err != nil {
		t.Log(err)
		t.FailNow()
	}
	enc, dec, stopCli, err := cli.Start()
	if err != nil {
		t.Fatal(err)
	}

	test1, err := hex.DecodeString("deadbeefcafe0080085000deadbeefcafe")
	if err != nil {
		t.Fatal(err)
	}
	t.Log("encoding")
	encRes := <-enc(
		&proto.EncodeRequest{
			Data: test1,
		},
	)

	t.Log(encRes.GetEncodedString())

	t.Log("decoding")
	decRes := <-dec(
		&proto.DecodeRequest{
			EncodedString: encRes.GetEncodedString(),
		},
	)
	t.Log("done")
	stopCli()
	stopServer()
	if string(test1) != string(decRes.GetData()) {
		t.Fatalf(
			"failed output equals input test: got %x expected %x",
			test1, decRes.GetData(),
		)
	}
}

package grpc

import (
	"encoding/hex"
	"github.com/quanterall/kitchensink/pkg/grpc/client"
	"github.com/quanterall/kitchensink/pkg/grpc/server"
	protos "github.com/quanterall/kitchensink/pkg/proto"
	"net"
	"testing"
)

const defaultAddr = "localhost:50051"

func TestGRPC(t *testing.T) {
	addr, err := net.ResolveTCPAddr("tcp", defaultAddr)
	if err != nil {
		t.Fatal(err)
	}
	srvr := server.New(addr, 8)
	stopSrvr := srvr.Start()

	cli, disconnect := client.New(defaultAddr)

	test1, err := hex.DecodeString("deadbeefcafe0080085000deadbeefcafe")
	if err != nil {
		t.Fatal(err)
	}
	encRes, err := cli.Encode(&protos.EncodeRequest{
		Data: test1,
	},
	)
	decRes, err := cli.Decode(&protos.DecodeRequest{
		EncodedString: encRes.GetEncodedString(),
	},
	)

	disconnect()
	stopSrvr()
	if string(test1) != string(decRes.GetData()) {
		t.Fatalf("failed output equals input test: got %x expected %x",
			test1, decRes.GetData(),
		)
	}
}

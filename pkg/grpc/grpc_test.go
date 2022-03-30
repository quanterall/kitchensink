package grpc

import (
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

	encRes, err := cli.Encode(&protos.EncodeRequest{
		Data: make([]byte, 32),
	},
	)
	t.Logf("resp: %v, err: %v", encRes, err)

	decRes, err := cli.Decode(&protos.DecodeRequest{
		EncodedString: encRes.GetEncodedString(),
	},
	)
	t.Logf("resp: %x, err: %v", decRes.GetDecoded(), err)

	disconnect()
	stopSrvr()
}

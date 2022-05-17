package client

import (
	"github.com/quanterall/kitchensink/pkg/proto"
	"time"
)

type encReq struct {
	Req *proto.EncodeRequest
	Res chan *proto.EncodeResponse
}

func newEncReq(req *proto.EncodeRequest) encReq {
	req.IdNonce = uint64(time.Now().UnixNano())
	return encReq{Req: req, Res: make(chan *proto.EncodeResponse)}
}

type decReq struct {
	Req *proto.DecodeRequest
	Res chan *proto.DecodeResponse
}

func newDecReq(req *proto.DecodeRequest) decReq {
	req.IdNonce = uint64(time.Now().UnixNano())
	return decReq{Req: req, Res: make(chan *proto.DecodeResponse)}
}

type b32c struct {
	addr       string
	encChan    chan encReq
	encRes     chan *proto.EncodeResponse
	decChan    chan decReq
	decRes     chan *proto.DecodeResponse
	stop       <-chan struct{}
	timeout    time.Duration
	waitingEnc map[time.Time]encReq
	waitingDec map[time.Time]decReq
}

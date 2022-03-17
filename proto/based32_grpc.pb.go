// Code generated by protoc-gen-go-grpc. DO NOT EDIT.

package protos

import (
	context "context"
	grpc "google.golang.org/grpc"
	codes "google.golang.org/grpc/codes"
	status "google.golang.org/grpc/status"
)

// This is a compile-time assertion to ensure that this generated file
// is compatible with the grpc package it is being compiled against.
// Requires gRPC-Go v1.32.0 or later.
const _ = grpc.SupportPackageIsVersion7

// TranscriberClient is the client API for Transcriber service.
//
// For semantics around ctx use and closing/ending streaming RPCs, please refer to https://pkg.go.dev/google.golang.org/grpc/?tab=doc#ClientConn.NewStream.
type TranscriberClient interface {
	Encode(ctx context.Context, opts ...grpc.CallOption) (Transcriber_EncodeClient, error)
	Decode(ctx context.Context, opts ...grpc.CallOption) (Transcriber_DecodeClient, error)
}

type transcriberClient struct {
	cc grpc.ClientConnInterface
}

func NewTranscriberClient(cc grpc.ClientConnInterface) TranscriberClient {
	return &transcriberClient{cc}
}

func (c *transcriberClient) Encode(ctx context.Context, opts ...grpc.CallOption) (Transcriber_EncodeClient, error) {
	stream, err := c.cc.NewStream(ctx, &Transcriber_ServiceDesc.Streams[0], "/signer.Transcriber/Encode", opts...)
	if err != nil {
		return nil, err
	}
	x := &transcriberEncodeClient{stream}
	return x, nil
}

type Transcriber_EncodeClient interface {
	Send(*EncodeRequest) error
	Recv() (*EncodeResponse, error)
	grpc.ClientStream
}

type transcriberEncodeClient struct {
	grpc.ClientStream
}

func (x *transcriberEncodeClient) Send(m *EncodeRequest) error {
	return x.ClientStream.SendMsg(m)
}

func (x *transcriberEncodeClient) Recv() (*EncodeResponse, error) {
	m := new(EncodeResponse)
	if err := x.ClientStream.RecvMsg(m); err != nil {
		return nil, err
	}
	return m, nil
}

func (c *transcriberClient) Decode(ctx context.Context, opts ...grpc.CallOption) (Transcriber_DecodeClient, error) {
	stream, err := c.cc.NewStream(ctx, &Transcriber_ServiceDesc.Streams[1], "/signer.Transcriber/Decode", opts...)
	if err != nil {
		return nil, err
	}
	x := &transcriberDecodeClient{stream}
	return x, nil
}

type Transcriber_DecodeClient interface {
	Send(*DecodeRequest) error
	Recv() (*DecodeResponse, error)
	grpc.ClientStream
}

type transcriberDecodeClient struct {
	grpc.ClientStream
}

func (x *transcriberDecodeClient) Send(m *DecodeRequest) error {
	return x.ClientStream.SendMsg(m)
}

func (x *transcriberDecodeClient) Recv() (*DecodeResponse, error) {
	m := new(DecodeResponse)
	if err := x.ClientStream.RecvMsg(m); err != nil {
		return nil, err
	}
	return m, nil
}

// TranscriberServer is the server API for Transcriber service.
// All implementations must embed UnimplementedTranscriberServer
// for forward compatibility
type TranscriberServer interface {
	Encode(Transcriber_EncodeServer) error
	Decode(Transcriber_DecodeServer) error
	mustEmbedUnimplementedTranscriberServer()
}

// UnimplementedTranscriberServer must be embedded to have forward compatible implementations.
type UnimplementedTranscriberServer struct {
}

func (UnimplementedTranscriberServer) Encode(Transcriber_EncodeServer) error {
	return status.Errorf(codes.Unimplemented, "method Encode not implemented")
}
func (UnimplementedTranscriberServer) Decode(Transcriber_DecodeServer) error {
	return status.Errorf(codes.Unimplemented, "method Decode not implemented")
}
func (UnimplementedTranscriberServer) mustEmbedUnimplementedTranscriberServer() {}

// UnsafeTranscriberServer may be embedded to opt out of forward compatibility for this service.
// Use of this interface is not recommended, as added methods to TranscriberServer will
// result in compilation errors.
type UnsafeTranscriberServer interface {
	mustEmbedUnimplementedTranscriberServer()
}

func RegisterTranscriberServer(s grpc.ServiceRegistrar, srv TranscriberServer) {
	s.RegisterService(&Transcriber_ServiceDesc, srv)
}

func _Transcriber_Encode_Handler(srv interface{}, stream grpc.ServerStream) error {
	return srv.(TranscriberServer).Encode(&transcriberEncodeServer{stream})
}

type Transcriber_EncodeServer interface {
	Send(*EncodeResponse) error
	Recv() (*EncodeRequest, error)
	grpc.ServerStream
}

type transcriberEncodeServer struct {
	grpc.ServerStream
}

func (x *transcriberEncodeServer) Send(m *EncodeResponse) error {
	return x.ServerStream.SendMsg(m)
}

func (x *transcriberEncodeServer) Recv() (*EncodeRequest, error) {
	m := new(EncodeRequest)
	if err := x.ServerStream.RecvMsg(m); err != nil {
		return nil, err
	}
	return m, nil
}

func _Transcriber_Decode_Handler(srv interface{}, stream grpc.ServerStream) error {
	return srv.(TranscriberServer).Decode(&transcriberDecodeServer{stream})
}

type Transcriber_DecodeServer interface {
	Send(*DecodeResponse) error
	Recv() (*DecodeRequest, error)
	grpc.ServerStream
}

type transcriberDecodeServer struct {
	grpc.ServerStream
}

func (x *transcriberDecodeServer) Send(m *DecodeResponse) error {
	return x.ServerStream.SendMsg(m)
}

func (x *transcriberDecodeServer) Recv() (*DecodeRequest, error) {
	m := new(DecodeRequest)
	if err := x.ServerStream.RecvMsg(m); err != nil {
		return nil, err
	}
	return m, nil
}

// Transcriber_ServiceDesc is the grpc.ServiceDesc for Transcriber service.
// It's only intended for direct use with grpc.RegisterService,
// and not to be introspected or modified (even as a copy)
var Transcriber_ServiceDesc = grpc.ServiceDesc{
	ServiceName: "signer.Transcriber",
	HandlerType: (*TranscriberServer)(nil),
	Methods:     []grpc.MethodDesc{},
	Streams: []grpc.StreamDesc{
		{
			StreamName:    "Encode",
			Handler:       _Transcriber_Encode_Handler,
			ServerStreams: true,
			ClientStreams: true,
		},
		{
			StreamName:    "Decode",
			Handler:       _Transcriber_Decode_Handler,
			ServerStreams: true,
			ClientStreams: true,
		},
	},
	Metadata: "proto/based32.proto",
}

package service

// The following line generates the protocol, it assumes that `protoc` is in the
// path. This directive is run when `go generate` is run in the current package,
// or if a wildcard was used ( go generate ./... ).
//go:generate protoc -I=. --go_out=. ./api/based32.proto

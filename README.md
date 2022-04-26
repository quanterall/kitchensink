# kitchensink

- [Teaching Golang via building a Human Readable Binary Transcription Encoding Framework](#teaching-golang-via-building-a-human-readable-binary-transcription-encoding-framework)
- [Prerequisites](#prerequisites)
	- [Install Go](#install-go)
	- [Install Protobuf Compiler](#install-protobuf-compiler)
	- [Install gRPC plugins for Go](#install-grpc-plugins-for-go)
- [Step By Step:](#step-by-step)
	- [Step 1 Create the Protobuf specification](#step-1-create-the-protobuf-specification)
		- [The header section](#the-header-section)
		- [The service definition](#the-service-definition)
		- [The encode messages](#the-encode-messages)
		- [The decode messages](#the-decode-messages)
		- [The errors](#the-errors)
	- [Step 2 Complete the creation of the protobuf implementations](#step-2-complete-the-creation-of-the-protobuf-implementations)
		- [Making the output code more useful with some extensions](#making-the-output-code-more-useful-with-some-extensions)

## Teaching Golang via building a Human Readable Binary Transcription Encoding Framework

In this tutorial we will walk you through the creation from scratch of a human
readable binary transcription encoder/decoder.

This tutorial demonstrates the use of almost every possible and important
feature of Go. A "toy" implementation of a gRPC/protobuf microservice is added
in order to illustrate almost everything else.

Note that we choose gRPC because it is widely used for microservices and by
using it, a project is empowered to decouple the binary part of the API both
from the implementation language and from the possibility of developers
inadvertently creating a distributed monolith, which is very difficult to
change.

In order to demonstrate synchronisation primitives, waitgroups, atomics (not
mutexes, which are anyway a lot more troublesome to manage), the service will
keep track of the number of invocations, print this count in log updates, and
track the count using a concurrent safe atomic variable and show the variant
using a mutex instead, and run an arbitrary number of concurrent worker threads
that will start up and stop using waitgroups.

The final result is this library itself, and each step will be elaborated in
clear detail. Many tutorials leave out important things, and to ensure this does
not happen, each stage's parts will be also found in the [steps](./steps)
folder at the root of the repository.

## Prerequisites

This tutorial was developed on a system running Pop OS based on Ubuntu 21. As
such there may be some small differences compared to Ubuntu 20, but protobuf
will still be version 3.

In general, you will be deploying your binaries to systems also running ubuntu
20 or 21 or similar, on x86-64 platform, so the same instructions can be used in
setting up a fresh server when deploying. We will not cover Docker or any other
container system here.

Necessary things that you probably already have:

    sudo apt install -y build-essential git wget curl autoconf automake libtool

### Install Go

Go 1.17+ is recommended - unlike most other languages, the forward compatibility
guarantee is ironclad, so go to [https://go.dev/dl/](https://go.dev/dl/) and
pick the latest one (1.18 at time of writing), "copy link location" on the
relevant version (linux x86-64)

    cd
    mkdir bin 
    wget https://go.dev/dl/go1.18.linux-amd64.tar.gz
    tar xvf go1.18.linux-amd64.tar.gz

Using your favourite editor, open up `~/.bashrc` - or just

    nano ~/.bashrc

and put the following lines at the end

    export GOBIN=$HOME/bin
    export GOPATH=$HOME
    export GOROOT=$GOPATH/go
    export PATH=$HOME/go/bin:$HOME/.local/bin:$GOBIN:$PATH

save and close, and `ctrl-d` to kill the terminal session, and start a new one.

This also creates a proper place where `go install` will put produced binaries,
which is recommended for avoiding cluttering up repositories you are working on
with binaries and potentially accidentally adding them to the repository, which
can be very problematic if you are working on a BIG application (Go apps are
generally under 60mb in size but this is still a lot in a source code
repository).

### Install Protobuf Compiler

In order to build the project you will need the protobuf compiler installed.
This generates Go code from a protobuf specification for a service, creating all
the necessary handlers to call and to handle calls to implement the API of our
transcription codec.

    sudo apt install -y protobuf-compiler
    protoc --version  # Ensure compiler version is 3+

### Install gRPC plugins for Go

This didn't used to be as easy as it is now. This produces the gRPC generated
code which eliminates the need to construct an RPC framework, all the work is
done for you, you can now just connect it to a network transport and voila. This
installs the plugins, which is another reason why the `GOBIN` must be set and
added to `PATH`. Otherwise, the two tools that are installed in the following
commands will not be accessible.

    go install google.golang.org/protobuf/cmd/protoc-gen-go@v1.26
    go install google.golang.org/grpc/cmd/protoc-gen-go-grpc@v1.1

## Step By Step:

Click on the title of the step to see the state your repository should be in
when you have completed the step.

### [Step 1](steps/step1) Create the Protobuf specification

First thing you do when working with gRPC is define the protocol messages.

When creating a Go based repository, it is usual to put all of the internal
packages into a single folder called 'pkg'.

I will assume that you have created a fresh new repository on your github
account, with readme already created for it, licence is up to you, I always use
unlicence unless I am forking code with other licences.

Inside `pkg` create a new directory `proto` and create a new
file `based32.proto`

#### The header section

```protobuf
syntax = "proto3";
package codec;
option go_package = "github.com/quanterall/kitchensink/pkg/proto";
```

Firstly, the syntax line indicates the version of protobuf that is being used,
that is `proto3` currently.

Second, the package line has no effect when using the Go plugin for protobuf.

Third is the line that actually does things for the Go version. The path is the
same as what appears in the `import` line where this package is being imported
from, and should be the same as the repository root, plus `pkg/proto`.

#### The service definition

```protobuf
service Transcriber {
  rpc Encode(stream EncodeRequest) returns (stream EncodeResponse);
  rpc Decode(stream DecodeRequest) returns (stream DecodeResponse);
}
```

A `service` in protobuf defines the API for a named service. In this case we
have encode and decode.

Note the `stream` keywords in the messages. This means that the generated code
will create a streaming RPC for the messages, which means that handling the
scheduling of the work is in the hands of the implementation rather than using
the built in goroutine spawning default handler.

This is important because although for relatively small scale applications it
can be ok to let Go manage spawning and freeing goroutines, for a serious large
scale or high performance API, you should be keeping the goroutines warm and
delivering them jobs with channels, as we will be in this tutorial. Teaching
concurrency is one of the goals of this tutorial.

#### The encode messages

```protobuf
message EncodeRequest {
  bytes Data = 1;
}

message EncodeResponse {
  oneof Encoded {
    string EncodedString = 1;
    Error Error = 2;
  }
}
```

The request is very simple, it just uses bytes. Note that although in many cases
for such data in Go, it will be fixed length, such as 32 byte long keys, hashes,
and 65 byte long signatures, protobuf does not have a notion of fixed length
bytes. It is not necessary to add a field to designate the length of the
message, as this is handled correctly by the implementing code it generates.

The response uses a variant called `oneof` in protobuf. This is not native to
Go, and for which reason we will be showing a small fix we add to the generated
code package to account for this. In Go, returns are tuples,
usually `result, error`, but in other languages like Rust and C++ they are
encoded as a "variant" which is a type of `union` type. The nearest equivalent
in Go is an `interface` but interfaces can be anything. The union type was left
out of Go because it breaks C's otherwise strict typing system
(and yes, this has been one of the many ways in which C code has been exploited
to break security, which is why Go lacks it).

#### The decode messages

```protobuf

message DecodeRequest{
  string EncodedString = 1;
}

message DecodeResponse {
  oneof Decoded {
    bytes Data = 1;
    Error Error = 2;
  }
}
```

They are exactly the same, except the types are reversed, encode is bytes to
string, decode is string to bytes, reversing the process.

#### The errors

```protobuf
enum Error {
  ZERO_LENGTH = 0;
  CHECK_FAILED = 1;
  NIL_SLICE = 2;
  CHECK_TOO_SHORT = 3;
  INCORRECT_HUMAN_READABLE_PART = 4;
}

```

We want to define all possible errors within the context of the protobuf
specification, and make them informative enough that their text can be used to
inform the programmer, or user, what error has occurred, without having to add
further annotations.

The above shows all of the possible errors that can occur, notably `NIL_SLICE`
is a purely programmer error, that would be when calling encode with nil instead
of a byte slice `[]byte`. In a real project you would probably have to add these
as you go along, but we will skip that for now, and this is part of the reason
why the tutorial uses such a simple application.

### [Step 2](steps/step2) Complete the creation of the protobuf implementations

To run the protobuf compiler and generate the code, from the root of the
repository you run the following command:

    protoc --go_out=. --go_opt=paths=source_relative --go-grpc_out=. --go-grpc_opt=paths=source_relative ./pkg/proto/based32.proto

`--go_out` is essential for when generating Go versions of the protobuf handler
code. Without this, the remaining options are not understood.

This will change your repository files to look like this:

    .
    └── pkg
        └── proto
            ├── based32.pb.go
            ├── based32.proto
            └── based32_grpc.pb.go

`based32.pb.go` provides the necessary methods to encode and decode the data
between Go and protobuf formats for the defined messages.

`based32_grpc.pb.go` provides the methods to use gRPC to implement the API as
described in the `service` section of the `based32.proto` file.

#### Making the output code more useful with some extensions

There is two minor gotchas that current versions of the go plugins for protoc to
generate our RPC API that we are going to show a workaround for

- A stringer for the error types,
- code that transforms the result types for the calls from the `result, error`
  idiom to the variant syntax as shown in the protocol and necessary in many
  cases for variant using languages like Rust and C++ to conform with
  *their* static type breaking variant type.1

```go
// Package proto is the protocol buffers specification and generated code
// package for based32
//
// The extra `error.go` file provides helpers and missing elements from the
// generated code that make programming the protocol simpler.
package proto

import (
	transcribe "github.com/quanterall/kitchensink"
)

// Error implements the Error interface which allows this error to automatically
// generate from the error code.
//
// Fixes a bug in the generated code, which not
// only lacks the Error method it uses int32 for the error string map when it
// should be using the defined Error type. No easy way to report the bug in the
// code.
//
// With this method implemented, one can simply return the error map code
// protos.Error_ERROR_NAME_HERE and logs print this upper case snake case which
// means it can be written to be informative in the proto file and concise in
// usage, and with this tiny additional helper, very easy to return, and print.
func (x Error) Error() string {

	return Error_name[int32(x)]
}

// CreateEncodeResponse is a helper to turn a transcribe.EncodeRes into an
// EncodeResponse to be returned to a gRPC client.
func CreateEncodeResponse(res transcribe.EncodeRes) (response *EncodeResponse) {

	// First, create the response structure.
	response = &EncodeResponse{}

	// Because the protobuf struct is essentially a Variant, a structure that
	// does not exist in Go, there is an implicit contract that if there is an
	// error, there is no return value. This is not implicit in Go's tuple
	// returns.
	//
	// Thus, if there is an error, we return that, otherwise, the value in the
	// response.

	if res.Error != nil {
		response.Encoded = &EncodeResponse_Error{
			Error(Error_value[res.Error.Error()]),
		}
	} else {
		response.Encoded =
			&EncodeResponse_EncodedString{
				res.String,
			}
	}
	return
}

// CreateDecodeResponse is a helper to turn a transcribe.DecodeRes into an
// DecodeResponse to be returned to a gRPC client.
func CreateDecodeResponse(res transcribe.DecodeRes) (response *DecodeResponse) {

	// First, create the response structure.
	response = &DecodeResponse{}

	// Return an error if there is an error, otherwise return the response data.
	if res.Error != nil {
		response.Decoded = &DecodeResponse_Error{
			Error(Error_value[res.Error.Error()]),
		}
	} else {
		response.Decoded = &DecodeResponse_Data{res.Bytes}
	}
	return
}

```

Note that the above code is not strictly necessary but has to be manually
handled later on one way or another, so we put this in now because there is no
reason for the learner to have to learn the details of why this should be, it
should have been fixed and will probably, hopefully be fixed in the go plugins
for protobuf in the future.

The error stringer saves duplicating effort in creating error return values 
for the programmer and user to read, and the `Create*Response` methods 
eliminate duplication in correctly translating the tuple into the variant 
form via the Go interface syntax.


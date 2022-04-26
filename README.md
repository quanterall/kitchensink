# kitchensink

- [Teaching Golang via building a Human Readable Binary Transcription Encoding Framework](#teaching-golang-via-building-a-human-readable-binary-transcription-encoding-framework)
- [Prerequisites](#prerequisites)
	- [Install Go](#install-go)
	- [Install Protobuf Compiler](#install-protobuf-compiler)
	- [Install gRPC plugins for Go](#install-grpc-plugins-for-go)
- [Step By Step:](#step-by-step)
	- [Step 1](#step-1)
- [Step 2](#step-2)

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

In each case, if you click on the title of the step, it will take you to the 
step folder where the final state of the repository for each step can be 
seen in case the tutorial has been unclear where or what you were supposed 
to end up at, at that point.

### [Step 1](steps/step1)

First thing you do when working with gRPC is define the protocol messages.

When creating a Go based repository, it is usual to put all of the internal
packages into a single folder called 'pkg'.

I will assume that you have created a fresh new repository on your github
account, with readme already created for it, licence is up to you, I always use
unlicence unless I am forking code with other licences.

Inside `pkg` create a new directory `proto` and create a new
file `based32.proto`

First part of the protobuf file is the header section:

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

Next, the service definition:

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

Next, the encode messages:

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

The request is very simple, it just uses bytes. Note that although in many 
cases for such data in Go, it will be fixed length, such as 32 byte long 
keys, hashes, and 65 byte long signatures, protobuf does not have a notion 
of fixed length bytes. It is not necessary to add a field to designate the 
length of the message, as this is handled correctly by the implementing code 
it generates.

The response uses a variant called `oneof` in protobuf. This is not native 
to Go, and for which reason we will be showing a small fix we add to the 
generated code package to account for this. In Go, returns are tuples, 
usually `result, error`, but in other languages like Rust and C++ they are 
encoded as a "variant" which is a type of `union` type. The nearest 
equivalent in Go is an `interface` but interfaces can be anything. The union 
type was left out of Go because it breaks C's otherwise strict typing system 
(and yes, this has been one of the many ways in which C code has been 
exploited to break security, which is why Go lacks it).

Next, the decode messages:

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

```protobuf
enum Error {
  ZERO_LENGTH = 0;
  CHECK_FAILED = 1;
  NIL_SLICE = 2;
  CHECK_TOO_SHORT = 3;
  INCORRECT_HUMAN_READABLE_PART = 4;
}

```

## [Step 2](steps/step2)


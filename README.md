# kitchensink

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
pick the latest one (1.17.8 at time of writing), "copy link location" on the
relevant version (linux x86-64)

    cd
    mkdir bin 
    wget https://go.dev/dl/go1.17.8.linux-amd64.tar.gz
    tar xvf go1.17.8.linux-amd64.tar.gz

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

Here will be the step by step process of building the library with a logical
sequence that builds from the basis to the specific parts for each in the order
that is needed both for understanding and for the constraints of syntax, grammar
and build system design...

### Setting up the repository

The very first thing to do is to create a Git repository. It is not necessary to
upload it to your github account or other git hosting, but it's standard first
thing to do, in most cases.

I will assume your workspace is something like `/home/username/code`. First,
open up a terminal and go there:

    $ cd code
    $ mkdir codec
    $ cd codec

Next, initialise the git repository

     $ git init 
     Initialized empty Git repository in /home/username/code/codec/.git/

If you are going to use github or so, you can instead just create the repository
there and

    $ git clone git@github.com/username/codec.git 

or something like this, whatever the git URL is. It is recommended to use ssh,
it just makes life simpler and it's faster and no complications with
authentication as you have if you use `https` instead.

### Set up the folder hierarchy

You can do this with your IDE's file manager or you can use the terminal as you
like, but start by making the folder structure that will be used:

    codec
    └── pkg
        ├── based32
        ├── codecer
        └── proto

It is idiomatic for Go projects to have a `pkg` folder where most of the
supporting libraries are found. Generally primary packages live in the root of
the repository, in this case the `codec` folder.

These are the first three folders that will have content put in them.

- `based32` will contain the actual implementing code for the human
  transcription encoder.
- `codecer` is where the interface specification lives, it is kept separate so
  it does not form any circular dependencies between consumer and
  implementation. It is the idiom in Go to take a noun related to the purpose of
  the interface and turn it into a 'doer' such as String to Stringer.
- `proto` contains mostly the gRPC protocol specification, the files that the
  tooling you installed beforehand will generate, as well as some additional
  code that fills in gaps in the current generated files to simplify the use of
  the protocol.

### gRPC/Protobuf specification

First thing we are going to put in place is the protocol specification. This
file will be called `based32.proto` inside the `pkg/proto/` folder.

It contains a service definition, defining an API call, and the request and
response are then elaborated lower down, as well as the errors that can be
returned.

```protobuf
syntax = "proto3";
package codec;
option go_package = "github.com/quanterall/kitchensink/service/proto";

service Transcriber {
  rpc Encode(stream EncodeRequest) returns (stream EncodeResponse);
  rpc Decode(stream DecodeRequest) returns (stream DecodeResponse);
}
enum Error {
  ZERO_LENGTH = 0;
  CHECK_FAILED = 1;
  NIL_SLICE = 2;
  CHECK_TOO_SHORT = 3;
  INCORRECT_HUMAN_READABLE_PART = 4;
}

message EncodeRequest {
  bytes Data = 1;
}

message EncodeResponse {
  oneof Encoded {
    string EncodedString = 1;
    Error Error = 2;
  }
}

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

Several things need to be noted in order to understand why they are there.

The `option` section needs to point to the URL of the location where the
generated source files created by the protocol compiler and tools will be
placed.

In the `service` section you can see the keyword stream. This is there because
we are creating a concurrent implementation, meaning that the RPC framework will
deliver new messages in a continuous stream and these are fanned out to worker
threads.

In the `enum` section, we have all of the possible error codes that will appear,
as we are aiming to have a complete specification in this file that all other
code refers to as a central point of contact. This is to eliminate any confusion
about the protocol should it be implemented in another language framework.

The rest is the usual common idiom for protobuf definitions, with requests and
responses as pairs. It is a central principle that you will encounter generally
throughout software development but most especially in Go, almost everything
comes in pairs. Go does not have variant types (`oneof`) because this can be
replaced with interfaces in the language. For this reason we have helpers added
to this folder that will simplify the syntax of constructing them while
retaining Go tuple syntax on the Go side.

### Interface

It is not always necessary to make interfaces when there will only ever be one
implementation for a given API, however, we will show how to use them in this
tutorial as you will encounter them and probably need them sooner or later.

Interfaces are essentially a special type of pointer with a type signature that
associates with a set of methods that must be implemented for the type to be
available to use in the place of the interface type.

The users of interfaces do not need to know anything about the internal
representations of data and are an essential tool to avoiding circular
dependencies, which are not permitted in Go because they cause loops in the
syntax tree which have to be resolved in order for the compiler to proceed.

In other languages, there can be a directive to only include such a dependency
once, but the result of all this permissive structuring is that it can be
difficult for the compiler to identify what symbol is being referred to, and the
logic for breaking these loops is complex, and that costs time in compilation.

The most well known interface in Go is the `Stringer` which is satisfied when a
method is created to render a variable into a string, with the signature
`String() string` and is used to render text through the `fmt` standard library
package.

The usual convention is to name the interface after the thing it does, so here
we use the word Codecer, which means something that makes an encoder/decoder
that produces strings from binary and binary back to strings, a type of codec,
thus, `Codecer`.

This file should go in `pkg/codecer/codecer.go`

```go
package codecer

// Codecer is the externally usable interface which provides a check for
// complete implementation as well as illustrating the use of interfaces in Go.
type Codecer interface {

	// Encode takes an arbitrary length byte input and returns the output as
	// defined for the codec.
	Encode(input []byte) (output string, err error)

	// Decode takes an encoded string and returns if the encoding is valid and
	// the value passes any check function defined for the type.
	//
	// If the check fails or the input is too short to have a check, false and
	// nil is returned. This is the contract for this method that
	// implementations should uphold.
	Decode(input string) (output []byte, err error)
}

```

You can see that this specification is very much the same basic specification as
you find in the proto file. We are not going to actually implement two or more
versions of this interface, however, it is important to understand how they
work, and what we will do next illustrate what is called "implementing an
interface".

### Specifying the data structure for the package

```go
package transcribe

// The following line generates the protocol, it assumes that `protoc` is in the
// path. This directive is run when `go generate` is run in the current package,
// or if a wildcard was used ( go generate ./... ).
//go:generate protoc -I=. --go_out=. --go_opt=paths=source_relative --go-grpc_out=. --go-grpc_opt=paths=source_relative ./pkg/proto/based32.proto

import (
	"github.com/quanterall/kitchensink/pkg/codecer"
)

type EncodeRes struct {
	String string
	Error  error
}

type DecodeRes struct {
	Bytes []byte
	Error error
}

// Codec is the collection of elements that creates a Human Readable Binary
// Transcription Codec
type Codec struct {

	// Name is the human readable name given to this encoder
	Name string

	// HRP is the Human Readable Prefix to be appended in front of the encoding
	// to disambiguate it from another encoding or as a network or protocol
	// identifier. 
	HRP string

	// Charset is the set of characters that the encoder uses. This should match
	// the output encoder, 32 for using base32, 64 for base64, etc.
	Charset string

	// Encode takes an arbitrary length byte input and returns the output as
	// defined for the codec
	Encoder func(input []byte) (output string, err error)

	// Decode takes an encoded string and returns if the encoding is valid and
	// the value passes any check function defined for the type.
	Decoder func(input string) (output []byte, err error)

	// AddCheck is used by Encode to add extra bytes for the checksum to ensure
	// correct input so user does not send to a wrong address by mistake, for
	// example.
	MakeCheck func(input []byte, checkLen int) (output []byte)

	// Check returns whether the check is valid
	Check func(input []byte) (err error)
}

// This ensures the interface is satisfied for codecer.Codecer and is removed in
// the generated binary because the underscore indicates the value is discarded.
var _ codecer.Codecer = &Codec{}

// Encode implements the codecer.Codecer.Encode by calling the provided
// function, and allows the concrete Codec type to always satisfy the interface,
// while allowing it to be implemented entirely differently.
//
// Note: short functions like this can be one-liners according to gofmt.
func (c Codec) Encode(input []byte) (string, error) { return c.Encoder(input) }

// Decode implements the codecer.Codecer.Decode by calling the provided
// function, and allows the concrete Codec type to always satisfy the interface,
// while allowing it to be implemented entirely differently.
//
// Note: this also can be a one liner. Since we name the return values in the
// type definition and interface, omitting them here makes the line short enough
// to be a one liner.
func (c Codec) Decode(input string) ([]byte, error) { return c.Decoder(input) }

```

There is a few things in here, so I will describe them part by part in order
they appear.

First, you can see that it imports the `codecer` package, as this interface is
implemented further down.

There is result types specified as they make the implementation more simple with
the channels used to implement the worker pool, as there is no tuple type in Go,
only struct, which has a definite type and order and compact storage format.

It is not strictly necessary to export any of the elements of the Codec struct
as seen here. In Go, exporting a symbol is denoted by using a capital letter.

Normally one would only export values in this way if it were safe for the caller
to modify them. In this case, we do indeed modify them. In the remaining folder
we mentioned to start with, `pkg/based32` we have an initialiser function that
populates the exported function types in the `Codec` type above, so all of 
the values do need to be exported as we haven't defined getter/setter 
functions for them so they can be hidden. 

Though they could be changed dynamically by consuming code, it just wouldn't 
probably happen, but if the values were changed during runtime it could 
cause race conditions if multiple threads are accessing this type, as it 
does in the gRPC implementation we will show later. It is expected that the 
consumer of the code won't bother modifying this dynamically as it is for 
the purpose of implementing a specific encoding, and every single element of 
the struct is part of a specification that defines the encode and decode 
process and the form of the outputs of these functions.

Further, as you can see, we make something like an assertion that the interface
is implemented, specifying that the struct type `Codec` is also
`codec.Codecer`, the interface. Then after this assertion is the actual 
implementations.

These use, as you can see, the exported functions `Encoder` and `Decoder`, 
and thereby satisfy the the interface, which as you saw in the previous 
section, requires `Encode` and `Decode`, and this source file will compile 
without error. Well, it doesn't do anything yet.

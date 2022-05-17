# kitchensink

- [Teaching Golang via building a Human Readable Binary Transcription Encoding Framework](#teaching-golang-via-building-a-human-readable-binary-transcription-encoding-framework)
	- [Prerequisites](#prerequisites)
	- [Install Go](#install-go)
	- [Install Protobuf Compiler](#install-protobuf-compiler)
	- [Install gRPC plugins for Go](#install-grpc-plugins-for-go)
	- [Initialize your repository](#initialize-your-repository)
- [Step By Step:](#step-by-step)
	- [Step 1 Create the Protobuf specification](#step-1-create-the-protobuf-specification)
		- [The header section](#the-header-section)
		- [The service definition](#the-service-definition)
		- [The encode messages](#the-encode-messages)
		- [The decode messages](#the-decode-messages)
		- [The errors](#the-errors)
	- [Step 2 Complete the creation of the protobuf implementations](#step-2-complete-the-creation-of-the-protobuf-implementations)
	- [Step 3 Create the base types and interfaces](#step-3-create-the-base-types-and-interfaces)
		- [Create the interface](#create-the-interface)
		- [The concrete type](#the-concrete-type)
		- [Package header](#package-header)
		- [Defining a generalised type framework](#defining-a-generalised-type-framework)
		- [Interface Implementation Assertion](#interface-implementation-assertion)
		- [Interface implementation using an embedded function](#interface-implementation-using-an-embedded-function)
		- [Making the gRPC generated code more useful with some extensions](#making-the-grpc-generated-code-more-useful-with-some-extensions)
		- [Documentation comments in Go](#documentation-comments-in-go)
		- [go:generate line](#gogenerate-line)
		- [Adding a Stringer `Error()` for the generated Error type](#adding-a-stringer-error-for-the-generated-error-type)
		- [Convenience types for results](#convenience-types-for-results)
		- [Create Response Helper Functions](#create-response-helper-functions)
	- [Step 4 The Encoder](#step-4-the-encoder)
		- [Always write code to be extensible](#always-write-code-to-be-extensible)
		- [Helper functions](#helper-functions)
		- [Log at the site](#log-at-the-site)
		- [Create an Initialiser](#create-an-initialiser)
		- [Writing the check function](#writing-the-check-function)
		- [Creating the Encoder](#creating-the-encoder)
		- [Calculating the check length](#calculating-the-check-length)
		- [Writing the Encoder Implementation](#writing-the-encoder-implementation)
		- [About `make()`](#about-make)
		- [Creating the Check function](#creating-the-check-function)
		- [Creating the Decoder function](#creating-the-decoder-function)
	- [Step 5 Testing the algorithm](#step-5-testing-the-algorithm)
		- [Random generation of test data](#random-generation-of-test-data)
		- [Running the tests](#running-the-tests)
		- [Enabling logging in the tests](#enabling-logging-in-the-tests)
		- [The Go tool recursive descent notation](#the-go-tool-recursive-descent-notation)
		- [Running the tests with logging and recursive descent](#running-the-tests-with-logging-and-recursive-descent)
		- [Actually testing the Encoder and Decoder](#actually-testing-the-encoder-and-decoder)
	- [Step 6 Creating a Server](#step-6-creating-a-server)
		- [The Logger](#the-logger)
		- [Implementing the worker pool](#implementing-the-worker-pool)
		- [When to not export an externally used type](#when-to-not-export-an-externally-used-type)
		- [About Channels](#about-channels)
		- [About Waitgroups](#about-waitgroups)
		- [Initialising the Worker Pool](#initialising-the-worker-pool)
		- [Atomic Counters](#atomic-counters)
		- [Running the Worker Pool](#running-the-worker-pool)
		- [Logging the call counts](#logging-the-call-counts)
		- [Starting the worker pool](#starting-the-worker-pool)
		- [Creating the gRPC Service](#creating-the-grpc-service)
	- [Step 7 Creating a Client](#step-7-creating-a-client)
		- [Data Structures for the Client](#data-structures-for-the-client)
		- [Client Constructor](#client-constructor)
		- [The Encode and Decode Handlers](#the-encode-and-decode-handlers)
		- [Copy Paste Generic Generator](#copy-paste-generic-generator)
	- [Step 8 Testing the gRPC server](#step-8-testing-the-grpc-server)

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

[->contents](#kitchensink)

----

### Prerequisites

In general, you will be deploying your binaries to systems also running ubuntu
20 or 21 or similar, on x86-64 platform, so the same instructions can be used in
setting up a fresh server when deploying. We will not cover Docker or any other
container system here.

Necessary things that you probably already have:

    sudo apt install -y build-essential git wget curl autoconf automake libtool

[->contents](#kitchensink)

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

[->contents](#kitchensink)

### Install Protobuf Compiler

In order to build the project you will need the protobuf compiler installed.
This generates Go code from a protobuf specification for a service, creating all
the necessary handlers to call and to handle calls to implement the API of our
transcription codec.

    sudo apt install -y protobuf-compiler
    protoc --version  # Ensure compiler version is 3+

[->contents](#kitchensink)

### Install gRPC plugins for Go

This didn't used to be as easy as it is now. This produces the gRPC generated
code which eliminates the need to construct an RPC framework, all the work is
done for you, you can now just connect it to a network transport and voila. This
installs the plugins, which is another reason why the `GOBIN` must be set and
added to `PATH`. Otherwise, the two tools that are installed in the following
commands will not be accessible.

    go install google.golang.org/protobuf/cmd/protoc-gen-go@v1.26
    go install google.golang.org/grpc/cmd/protoc-gen-go-grpc@v1.1

[->contents](#kitchensink)

### Initialize your repository

Whether you create the repository directly on your github or other account or
not, you need to first initialise the go modules, from the root of your new
repository, like this:

    go mod init github.com/quanterall/kitchensink

**IMPORTANT**

**You need to change this and every other instance of
`github.com/quanterall/kitchensink` found in the source code of this tutorial to
match what you defined in your version of the above statement.**

*If you don't upload this anywhere, you can just use what is defined here, and
avoid this chore.*

[->contents](#kitchensink)

## Step By Step:

Click on the title of the step to see the state your repository should be in
when you have completed the step.

----

### [Step 1](steps/step1) Create the Protobuf specification

First thing you do when working with gRPC is define the protocol messages.

When creating a Go based repository, it is usual to put all of the internal
packages into a single folder called 'pkg'.

I will assume that you have created a fresh new repository on your github
account, with readme already created for it, licence is up to you, I always use
unlicence unless I am forking code with other licences.

Inside `pkg` create a new directory `proto` and create a new
file `based32.proto`

[->contents](#kitchensink)

#### The header section

```protobuf
syntax = "proto3";
package proto;
option go_package = "github.com/quanterall/kitchensink/pkg/proto";
```

Firstly, the syntax line indicates the version of protobuf that is being used,
that is `proto3` currently.

Second, the package line has no effect when using the Go plugin for protobuf.

Third is the line that actually does things for the Go version. The path is the
same as what appears in the `import` line where this package is being imported
from, and should be the same as the repository root, plus `pkg/proto`.

[->contents](#kitchensink)

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

[->contents](#kitchensink)

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
code package to account for this.

In Go, returns are tuples, usually `result, error`, but in other languages like
Rust and C++ they are encoded as a "variant" which is a type of `union` type.
The nearest equivalent in Go is an `interface` but interfaces can be anything.
The union type was left out of Go because it breaks C's otherwise strict typing
system (and yes, this has been one of the many ways in which C code has been
exploited to break security, which is why Go lacks it).

[->contents](#kitchensink)

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

[->contents](#kitchensink)

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

[->contents](#kitchensink)

----

### [Step 2](steps/step2) Complete the creation of the protobuf implementations

To run the protobuf compiler and generate the code, from the root of the
repository you run the following command:

    protoc --go_out=. --go_opt=paths=source_relative --go-grpc_out=. \ 
      --go-grpc_opt=paths=source_relative ./pkg/proto/based32.proto

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

[->contents](#kitchensink)

----

### [Step 3](steps/step3) Create the base types and interfaces

#### Create the interface

In this project we are creating an interface in part to demonstrate how to use
them. Being such a small library, it may not be necessary to do this, but it is
rare that you will be making such small packages in practise, so we want to show
you how to work with interfaces.

Create a new folder inside `pkg/` called `codecer`. The name comes from the
convention in Go to name interfaces by what they do. A codec is an
encoder/decoder, so a thing that defines a codec interface is a `codecer`.

In `pkg/codecer` create a file `interface.go` and put this in it:

```go
// Package codecer is the interface definition for a Human Readable Binary
// Transcription Codec

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

Note that in the comments we are specifying some things about the contract that
this interface should uphold. This is a good practise to help users of your
libraries know what to expect, and so you don't create or inspire someone to
create an 'undefined behaviour' that could become a security vulnerability.

[->contents](#kitchensink)

#### The concrete type

Create a new folder [pkg/codec](pkg/codec) and in it create a new file called `types.go`. This is
where we will define the main types that will be used by packages and
applications that use our code.

[->contents](#kitchensink)

#### Package header

```go
package codec

import (
    "github.com/quanterall/kitchensink/pkg/codecer"
)
```

[->contents](#kitchensink)

#### Defining a generalised type framework

This is a configuration data structure that bundles configuration and
implementation functions together. The function types defined are able to be
changed by calling code, which we use to create an initialiser in the main
`based32` package a little later.

```go
// Codec is the collection of elements that creates a Human Readable Binary
// Transcription Codec
type Codec struct {
    
    // Name is the human readable name given to this encoder
    Name string
    
    // HRP is the Human Readable Prefix to be appended in front of the encoding
    // to disambiguate it from another encoding or as a network or protocol
    // identifier. This can be empty, but more usually this will be used to
    // disambiguate versus other similarly encoded values, such as used on a
    // different cryptocurrency network, or between main and test networks.
    HRP string
    
    // Charset is the set of characters that the encoder uses. This should match
    // the output encoder, 32 for using base32, 64 for base64, etc.
    //
    // For arbitrary bases, see the following function in the standard library:
    // https://cs.opensource.google/go/go/+/refs/tags/go1.17.7:src/strconv/itoa.go;l=25
    // This function can render up to base36, but by default uses 0-9a-z in its
    // representation, which would either need to be string substituted for
    // non-performance-critical uses or the function above forked to provide a
    // direct encoding to the intended characters used for the encoding, using
    // this charset string as the key. The sequence matters, each character
    // represents the cipher for a given value to be found at a given place in
    // the encoded number.
    Charset string
    
    // Encode takes an arbitrary length byte input and returns the output as
    // defined for the codec
    Encoder func (input []byte) (output string, err error)
    
    // Decode takes an encoded string and returns if the encoding is valid and
    // the value passes any check function defined for the type.
    Decoder func (input string) (output []byte, err error)
    
    // AddCheck is used by Encode to add extra bytes for the checksum to ensure
    // correct input so user does not send to a wrong address by mistake, for
    // example.
    MakeCheck func (input []byte, checkLen int) (output []byte)
    
    // Check returns whether the check is valid
    Check func (input []byte) (err error)
}
```

[->contents](#kitchensink)

#### Interface Implementation Assertion

The following var line makes it so the compiler will throw an error if the
interface is not implemented.

```go
// This ensures the interface is satisfied for codecer.Codecer and is removed in
// the generated binary because the underscore indicates the value is discarded.
var _ codecer.Codecer = &Codec{}
```

This is a good way to avoid problems when trying to use a concrete type, if the interface was changed, for example, and the existing implementations did not have the correct function signatures, the compiler will tell you this line doesn't work before, and in your IDE, should get red squiggly lines if something like this happens.

[->contents](#kitchensink)

#### Interface implementation using an embedded function

The type defined in the previous section provides for a changeable function for
encode and decode. These are used here to automatically satisfy the interface
defined in the previous section

You will notice that the receiver (the variable defined before the function
name) here is `*Codec`.

More often you will create methods that refer to the pointer to the type because
they will be struct types and methods that call on a non pointer method copy the
struct, which may not have the desired result as this will result in concurrent
copies of values that are not the same variable, and are discarded at the end of
this method's execution, potentially consuming a lot of memory and time moving that memory around. 

As well as potentially causing bugs when unexpected values are in these places when they should have been changed by other code somewhere else.

When the type is a potentially shared or structured (struct or `[]`) type, the
copy will waste time copying the value, or referring to a common version in the
pointer embedded within the slice type (or map), and memory to store the copy,
and potentially lead to errors from race conditions or unexpected state
divergence if the functions mutate values inside the structure. 

Usually non
pointer methods are only used on simple value types like specially modified
versions of value types (anything up to 64 bits in size, but also arrays, which
are `[number]type` as opposed to `[]type`), or when this copying behaviour is
intended to deliberately avoid race conditions, and the shallow copy will not
introduce unwanted behaviours or potentially cause race conditions with multiple threads modifying the data being pointed to.

In this case, the pointer is fine, because it is not intended that the `Codec` type ever be changed after initialisation. However, because its 
methods and fields are exposed, code that reaches through and modifies these could break this assumption. 

```go
// Encode implements the codecer.Codecer.Encode by calling the provided
// function, and allows the concrete Codec type to always satisfy the interface,
// while allowing it to be implemented entirely differently.
//
// Note: short functions like this can be one-liners according to gofmt.
func (c *Codec) Encode(input []byte) (string, error) { return c.Encoder(input) }

// Decode implements the codecer.Codecer.Decode by calling the provided
// function, and allows the concrete Codec type to always satisfy the interface,
// while allowing it to be implemented entirely differently.
//
// Note: this also can be a one liner. Since we name the return values in the
// type definition and interface, omitting them here makes the line short enough
// to be a one liner.
func (c *Codec) Decode(input string) ([]byte, error) { return c.Decoder(input) }
```

[->contents](#kitchensink)

#### Making the gRPC generated code more useful with some extensions

There is two minor gotchas that current versions of the go plugins for protoc to
generate our RPC API that we are going to show a workaround for

- A stringer for the error types,
- code that transforms the result types for the calls from the `result, error`
  idiom to the variant syntax as shown in the protocol and necessary in many
  cases for variant using languages like Rust and C++ to conform with
  *their* static type breaking variant types.

Create a new file called `error.go` in the `pkg/proto` directory. This file will
become part of the `proto` package and in here we can directly access internally
defined parts of the generated code if needed.

This is necessary to define the `Error() string` method on the error type, which
is missing from the generated code. This method is the equivalent of
`String() string` method which is called the 'Stringer interface' but
unfortunately gets a different name that doesn't actually tell you it just
returns a string, and that fmt.Print* functions consider this to be a variant of
the Stringer even though it's not part of this interface.

Fortunately, the generated code does create the string versions for you, just it
does not make access to them idiomatic. We are teaching you idiomatic Go here,
so this is necessary to account for the inconsistency of the API of the
generated code.

[->contents](#kitchensink)

#### Documentation comments in Go

First, take note that comments above the package line should start "Package
packagename..." and these lines will appear in
[https://pkg.go.dev/](https://pkg.go.dev/)
when you publish them on one of the numerous well known git hosting services,
and a user searches for them by their URL in the search bar at the top of that
page. This is referred to as 'godoc' and these texts can be also seen by using
commands with the `go` CLI application to print them to console or serve a local
version of what will appear on the above page.

Most IDEs for Go will nag you to put these in. The one above package is not so
important but every exported symbol (starting with a capital letter) in your
source code should have a comment starting with the symbol and explaining what
it is. The symbol can be preceded by an article, A or The if it makes more sense
to write it that way.

```go
// Package proto is the protocol buffers specification and generated code
// package for based32
//
// The extra `error.go` file provides helpers and missing elements from the
// generated code that make programming the protocol simpler.
package proto
````

[->contents](#kitchensink)

#### go:generate line

This is a convenient location to place the generator that processes the
`*.proto` files. It can be put anywhere but this makes it more concise.

In Goland IDE this can be invoked directly from the editor.

```go
// The following line generates the protocol, it assumes that `protoc` is in the
// path. This directive is run when `go generate` is run in the current package,
// or if a wildcard was used ( go generate ./... ).
//go:generate protoc --go_out=. --go_opt=paths=source_relative --go-grpc_out=. --go-grpc_opt=paths=source_relative ./based32.proto
```

[->contents](#kitchensink)

#### Adding a Stringer `Error()` for the generated Error type

The protobuf compiler creates a type Error to match the one defined in our proto
file, but, it does not automatically generate the stringer for it. Normal types
would just have a `String() string` function for this, but error types have a
special different 'stringer' called `Error() string`. This makes it possible to
return this Error type as an error, but at the same time easily print the text
defined in the proto file.

Yes, just defining this function does the both things in one, because any type
with a method with the signature `Error() string` becomes an implementation of
the `error` interface (yes, lower case, it is built in). I am not sure when this
came into force, some time around 1.13 version or so.

The Go protobuf plugins don't make any assumptions just because you call the
enumeration type "Error" that it means it should be an `error` type, so we have
to explicitly tell it this here, which makes consuming code much more readable
and understandable.

```go
// Error implements the Error interface which allows this error to automatically
// generate from the error code.
//
// With this method implemented, one can simply return the error map code
// protos.Error_ERROR_NAME_HERE and logs print this upper case snake case which
// means it can be written to be informative in the proto file and concise in
// usage, and with this tiny additional helper, very easy to return, and print.
func (x Error) Error() string {

    return Error_name[int32(x)]
}

```

[->contents](#kitchensink)

#### Convenience types for results

The following types will be used elsewhere, as well as for the following create
response functions. These are primarily to accommodate for the fact that
protobuf follows c++ conventions with the use of 'oneof' variant types, which
don't exist in Go.

```go
// EncodeRes makes a more convenient return type for the results
type EncodeRes struct {
    String string
    Error  error
}

// DecodeRes makes a more convenient return type for the results
type DecodeRes struct {
    Bytes []byte
    Error error
}
```

Yes, if you wanted to, you could use a structured type with error and return value if you preferred. The extra work for Go programmers is far greater than the no extra work for variant type using programmers.

[->contents](#kitchensink)

#### Create Response Helper Functions

The following functions create convenient functions to return the result or the
error correctly for creating the correct data structure for the gRPC response
messages.

It is idiomatic for protobuf to use these variant types (union of one of several
types) as protobuf was originally designed for C++ and other languages, such as
Rust and Java also use variants, but Go does not, as this is redundant
complexity.

It does mean that Go code that cooperates with these variant using languages and
conventions is more complicated, so we make these helpers. These probably could
be generated automatically but currently aren't. Common conventions are not
necessarily based on the best interests of compiler writers, as this case
exemplifies the underlying complexity that variants impose on the compiler
unnecessarily.

It is not mandatory, as a Go programmer, for you to obey these conventions, they
are here in this protobuf specification because you will encounter it a lot. For
the good of your fellow programmers, create return types with result and error.
There is cases where it makes sense to return a result AND an error, where such
an error is not fatal, and the variant return convention ignores this.

```go
// CreateEncodeResponse is a helper to turn a codec.EncodeRes into an
// EncodeResponse to be returned to a gRPC client.
func CreateEncodeResponse(res EncodeRes) (response *EncodeResponse) {
    
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
        response.Encoded = 
            &EncodeResponse_Error{
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

// CreateDecodeResponse is a helper to turn a codec.DecodeRes into an
// DecodeResponse to be returned to a gRPC client.
func CreateDecodeResponse(res DecodeRes) (response *DecodeResponse) {
    
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

It is worth explaining how the variant type is here implemented by the `protoc` code generator. 

The `.Encoded` and `.Decoded` fields in the response type are `interface{}` types, which means they can have any type included inside them, so long as the code using the data knows what types to expect, otherwise you will get a type assertion panic from an invalid type assertion.

The code generator creates two types, the `NameResponse_Error` and `NameResponse_NameString` or `NameResponse_NameData` as the case may be. These are then set to satisfy the interface type it generates for the `Decoded` and `Encoded` fields.

Obviously, this puts a pretty onerous burden on you as a Go programmer when you are obliged to use these, but this helper can be modified to fit any `oneof` using protobuf message and shift this ugly thing away from your main algorithms.

[->contents](#kitchensink)

----

### [Step 4](steps/step4) The Encoder

Next step is the actual library that the protobufs and interface and types were all created for.

[->contents](#kitchensink)

#### Always write code to be extensible

While making new libraries you will change the types and protocols a lot as you work through the implementation, it is still the best pattern to start with defining at least protocols and making a minimal placeholder for the implementation.

It is possible to create a library that does not have any significant state or configuration that just consists of methods that do things, which can be created without a `struct` tying them together, this is rare, and usually only happens when there is only one or two functions required. 

For this we have 4 functions and while we could hard code everything with constants and non-exported variables where constants can't be used, this is not extensible. Even in only one year of full time work programming, I estimate that I spent about 20% of my first year working as a Go developer, fixing up quick and dirty written code that was not designed to be extended. 

The time cost of preparing a codebase to be extensible and modular is tiny in comparison, maybe an extra half an hour as you start on a library. Experience says that the shortcut is not worth it. You never know when you are the one who has to extend your own code later on, and two days later it's finally in a state you can add functionality.

[->contents](#kitchensink)

#### Helper functions

The only exception to this is when there is literally only one or at most two functions to deal with a specific type of data. These are often referred to as "helpers" or "convenience functions" and do not need to be extensible as they are very small and self contained.

These functions can sometimes be tricky to know where to put them, and often end up in collections under package names with terrible names like "util" or "tools" or "helpers". This can be problematic because very often they are accessory to another type, and doing this creates confusing crosslinks that can lead you into a circular dependency.

As such, my advice is to keep helpers where they are used, and don't export them, unless they are necessary, like the response helper functions we made previously for the `proto` package.

[->contents](#kitchensink)

#### Log at the site

No code ever starts out perfect. In most cases every last bit has to be debugged at some point. As such, one of the most important things you can do to save yourself time and irritation is to make it easier to trace bugs.

For this reason, the first thing we are going to add is a customised logger for our package.

In the `pkg/` folder, create a new folder `based32`which will be the name of our human readable encoding library.

In this folder, create a file `log.go`:

```go
package based32

import (
    logg "log"
    "os"
)

var log = logg.New(os.Stderr, "based32" , logg.Llongfile|logg.Lmicroseconds)
```

What we are doing here is using the standard logging library to set up a customised configuration. The standard logger only drops the timestamp with the log entries, which is rarely a useful feature, and when it is, the time precision is too low on the default configuration, as the most frequent time one needs accurate timestamps is when the time of events is in the milli- or microseconds when debugging concurrent high performance low latency code.

This log variable essentially replaces an import in the rest of the package for the `log` standard library, and configures it to print full file paths and label them also with the name of the package. Anywhere in the same package now, `log` now refers to this customised logger.

We are adding the microseconds here as well, because with our concurrent code, less time precision would not reveal any information, as events are happening in time periods under 1000th of a second, called "milliseconds". Generally this is sufficient as the switching periods between goroutines are no more frequent than 1 microsecond.

It is ok to leave one level of indirection in the site of logging errors, that is, the library will return an error but not log, but it should at least log where the error returns, so that when the problem comes up, you only have to trace back to the call site and not several layers above this.

When you further have layers of indirection like interfaces and copies of pointers to objects that are causing errors, knowing which place to look for the bug will take up as much time as actually fixing it.

It may be that you are never writing algorithms that need any real debugging, many "programmers" rarely have to do much debugging. But we don't want to churn out script writers only, we want to make sure that everyone has at least been introduced to the idea of debugging. 

[->contents](#kitchensink)

#### Create an Initialiser

The purpose for the transparency of the `Codec` type in `types.go` was so that we could potentially create a custom codec in a separate package than where the type was defined. 

While we could have avoided this openness and created a custom function to load the non-exported struct members, and potentially in a more complex library, we would, for a simple library like this, we are going to assume that users of the library are either not going to tamper with it when its being used concurrently, or that they have created a custom implementation for their own encoder design or won't be using it concurrently.

So, the first things we are going to do is sketch out the creation of an initialiser, in which we will use closures to load the structure with functionality.

```go
// Package based32 provides a simplified variant of the standard
// Bech32 human readable binary codec
package based32

import (
    codec "github.com/quanterall/kitchensink"
)

// charset is the set of characters used in the data section of bech32 strings.
const charset = "qpzry9x8gf2tvdw0s3jn54khce6mua7l"

// Codec provides the encoder/decoder implementation created by makeCodec.
var Codec = makeCodec(
    "Base32Check",
    charset,
    "QNTRL",
)

// makeCodec generates our custom codec as above, into the exported Codec
// variable
func makeCodec(
    name string,
    cs string,
    hrp string,
) (cdc *codec.Codec) {
    
    return cdc
}
```

You will notice that we took care to make sure that everything you will paste into your editor will pass syntax checks immediately. All functions that have return values must contain a `return` statement. 

The return value here is imported from `types.go` at the root of the repository, which the compiler identifies as `github.com/quanterall/kitchensink` because of running `go mod init` in [Initialize your repository](#initialize-your-repository) .

> When you first start writing code, you will probably get quite irritated at having to put in those empty returns and put in the imports. This is just how things are with Go. The compiler is extremely strict about identifiers, all must be known, and a function with return value without a return is also wrong. 
>
> A decent Go IDE will save you time by adding and removing the `import` lines for you automatically if it knows them, but you are responsible for putting the returns in there. Note that you can put a `panic()` statement instead of `return`, the tooling in Goland, for example, when you generate an implementation for a type from an interface, puts `panic` calls in there to remind you to fill in your implementation.
>
> I will just plug Goland a little further, the reason why I recommend it is because it has the best hyperlink system available for any IDE on the market, and as I mentioned a little way back, tracing errors back to their source is one of the most time consuming parts of the work of a programmer, every little bit helps, and Jetbrains clearly listen to their users regarding this - even interfaces are easy to trace back to multiple implementations, again saving a lot of time when you are working on large codebases.
>

Before we start to show you how to put things into the `Codec` we first will just refresh your memory with a compact version of the structure that it defines:

```go
type Codec struct {
    Name string
    HRP string
    Charset string
    Encoder func(input []byte) (output string, err error)
    Decoder func(input string) (output []byte, err error)
    MakeCheck func(input []byte, checkLen int) (output []byte)
    Check func(input []byte) (err error)
}
```

HRP and Charset are configuration values, and Encoder, Decoder, MakeCheck and Check are functions.

The configuration part is simple to define, so add this to `makeCodec` :

```go
    // Create the codec.Codec struct and put its pointer in the return variable.
    cdc = &codec.Codec{
        Name:    name,
        Charset: cs,
        HRP:     hrp,
    }
```

This section:

```go

// charset is the set of characters used in the data section of bech32 strings.
const charset = "qpzry9x8gf2tvdw0s3jn54khce6mua7l"

// Codec provides the encoder/decoder implementation created by makeCodec.
var Codec = makeCodec(
    "Base32Check",
    charset,
    "QNTRL",
)

```

as you can now see, fills in these predefined configuration values for our codec.

In fact, the Name field is not used anywhere, but in an application that can work with multiple codec.Codec implementations, this could become quite useful to differentiate between them. It also functions as something in between documentation and code for the reader.

Stub in the closures for the codec:

```go
    cdc.MakeCheck = func(input []byte, checkLen int) (output []byte) {
           return
    }

    cdc.Encoder = func(input []byte) (output string, err error) {
        return
    }

    cdc.Check = func(input []byte) (err error) {
        return
    }

    cdc.Decoder = func(input string) (output []byte, err error) {
        return
    }
```

Again, we like to teach good, time saving, and error saving practices for programming. Making stubs for things that you know you eventually have to fill in is a good practice for this purpose. 

Note that the returns can be left 'naked' like this because the variables are declared in the type signature of the closure. If you leave out the names and only have the list of types of the return tuples, you have to also make the declarations of the variables, or fill in empty versions, which for `[]` types means `nil` for `error` also `nil` and for `string` the empty type is `""`. 

> There is something of an unofficial convention amongst Go programmers to not name return variables, it is the opinion of the author that this is a bad thing for readability, as the variable names can give information about what the values actually represent. In this case here they are named simply as they are quite unambiguous given the functions names, however, sometimes it can be very helpful to save the reader the time of scanning through the function to know what a return value relates to. 
>
> Further, the variables are likely to be declared somewhere arbitrarily through the text, or in a given if block or other error handling area filled in manually. Using predeclared names saves space, makes the function header the one stop shop to learn everything and saves time for everyone who has to deal with the code.
>
> Naked returns also get a bad rap for no good reason also. Unless the function spans more than a screenful, and you don't name the return values, it is redundant and noisy to repeat things when they are clear.
>
> Further, as a general comment, everything that happens at the root level of a Go source file exists all at once, there is no formal ordering to it except where variables are created, these are populated in the order they are shown prior to running `init()` functions and prior to invoking `main()` - In spite of the lack of syntactic requirement for strict ordering, it just doesn't make any sense to refer forwards to objects when humans read from top to bottom and thus you are then requiring the human to scan back and forth, a task that becomes a huge time waster the longer the source code gets. 
>
> (Of course, long source files are a bad thing, as equally as many that are too small... it is something that should be relatively intuitive since we all have brains that work much the same way.)

[->contents](#kitchensink)

#### Writing the check function

The first thing we need for the codec is the check function. The check makes a checksum value. We additionally add the requirement that the check serves double duty as a pad to fill out the Base32 strings, which otherwise get ugly padding characters, usually `=` in the case of standard Go base32 library. Thus, the check function also includes a length parameter.

```go
    // We need to create the check creation functions first
    cdc.MakeCheck = func(input []byte, checkLen int) (output []byte) {

        // We use the Blake3 256 bit hash because it is nearly as fast as CRC32
        // but less complicated to use due to the 32 bit integer conversions to
        // bytes required to use the CRC32 algorithm.
        checkArray := blake3.Sum256(input)

        // This truncates the blake3 hash to the prescribed check length
        return checkArray[:checkLen]
    }
```

It is possible to instead use the shorter, and simpler `crc32` checksum function, but we don't like it because it requires converting bytes to 32 bit integers and back again, which is done in practice like this:

```go
example := []byte{1, 2, 3, 4}
asInteger := uint32(example[0]) + 
    (uint32(example[1]) << 8) +
    (uint32(example[2]) << 16) +
    (uint32(example[3]) << 24) 

```

for little-endian, the ordering of the left shift (which means multiply by 2 to the power of the number of bits, which is done by bit shifting, left shift 1 bit means double) is reversed for big-endian, which you still encounter sometimes, such as in the encoding of 32 byte long hashes and 65-68 byte long signatures in Bitcoin blocks.

> You may not have been aware of this, but the modern arabic base 10 numbers are also written right to left as is the arabic convention. This also makes it annoying to write indentation algorithms for numbers as well. 
>
> The 'big' refers to the fact that by memory addresses, the almost complete majority of microchips encode the bytes and words with the same convention as we use for our base 10 numbers, ie, first bit is on the right, last bit is on the left, and then you move the counter for the next byte/word in the opposite direction. 
>
> Maybe one day common sense will prevail and the bits will be ordered sequentially and not have reversals for every 8, 16, 32 or 64 byte chunks, it would make it a lot simpler to think about since we number the bytes forward, and each word is backwards...

If you follow the logic of that conversion, you can see that it is 4 copy operations (copy 8 bit value to 32 bit *zero* value), 3 bit shifts and 3 addition operations. The hash function does not do this conversion, it operates directly on bytes (in fact, I think it uses 8 byte/64 bit words, and coerces the byte slices to 64 bit long words using an unsafe type conversion) using a sponge function, and Blake3 is the fastest hash function with cryptographic security, which means a low rate of collisions, which in terms of checksums equates to two strings creating the same checksum, and breaking some degree of security of the function. So, we use blake3 hashes and cut them to our custom length. 

It also introduces the use of one of the best hash functions currently available, Blake3, and chosen by most new projects due to its speed and security. Hashes are very important for authentication and encryption as they make it difficult to deliberately corrupt the data it is generated from.

The length is variable as we are designing this algorithm to combine padding together with the check. So, essentially the way it works is we take the modulus of 5 (remainder of the division) of the length of the data, and pad it out to the next biggest multiple of 5 bytes, which is 8 base32 symbols. The formula for this comes next.

[->contents](#kitchensink)

#### Creating the Encoder

In all cases, when creating a codec, the first step is making the encoder. It is impossible to decode something that doesn't yet exist, the encoder is *a priori*, that is, it comes first, both logically and temporally.

But before we can make the encoder, we also need a function to compute the correct check length for the payload data. 

Further, the necessity of a variable length requires also that the length of the check be known before decoding, so this becomes a prefix of the encoding.

[->contents](#kitchensink)

#### Calculating the check length

So, first thing to add is a helper function, which we recommend you put *before* the `makeCodec` function:

```go
func getCheckLen(length int) (checkLen int) {

    // The following formula ensures that there is at least 1 check byte, up to
    // 4, in order to create a variable length theck that serves also to pad to
    // 5 bytes per 8 5 byte characters (2^5 = 32 for base 32)
    //
    // we add two to the length before modulus, as there must be 1 byte for
    // check length and 1 byte of check
    lengthMod := (2 + length) % 5

    // The modulus is subtracted from 5 to produce the complement required to
    // make the correct number of bytes of total data, plus 1 to account for the
    // minimum length of 1.
    checkLen = 5 - lengthMod + 1

    return checkLen
}
```

The function takes the input of the length of our message in bytes, and returns the correct length for the check. The result is that for an equal multiple of 5 bytes, we add 5 bytes of check, and 4, 3, 2 or 1 bytes for each of the variations that are more than this multiple, plus accounting for the extra byte to store the check length.

> The check length byte in fact only uses the first 3 bits, as it can be no more than 5, which requires 3 bits of encoding. Keep this in mind for later as we use this to abbreviate the codes as implicitly their largest 5 bits must be zero, which is precisely one base32 character in length, thus it always prefixes with a `q` as described by the character set we are using for this, based on the Bech32 standard used by Bitcoin and Cosmos, which reduces the length of the encoded value by one character, and must be added back to correctly decode.

[->contents](#kitchensink)

#### Writing the Encoder Implementation

The standard library contains a set of functions to encode and decode base 32 numbers using custom character sets. The 32 characters defined in the initialiser defined earlier, are chosen for their distinctiveness, it does not have I and 1, only small L (l), and only has zero (0) not capital O, as, unfortunately, in many fonts these can be hard to differentiate. Likewise for 2 and Z, and 5 and S.

This character set is taken from Bech32 and the Cosmos SDK also uses this same character set by default. It produces the most unambiguous character strings possible, while being as compact as possible. The old Base58check used for Bitcoin addresses produced more compact codes but they were easier to mis-transcribe due to ambiguous characters. Base64 is another encoding that is default for JSON binary data but it is even more plagued by ambiguity than Base58. Even optical character recognition systems can misread the ambiguous characters especially in some fonts. Since we can't fix the fonts, we have to fix the codes.

Note that the Ethereum address standard only recently added a check value that uses capital letters in the hexadecimal encoding. It used to be easy to send ethereum to nobody before this. Bitcoin had a check in it from day one.

```go
    // Create a base32.Encoding from the provided charset.
    enc := base32.NewEncoding(cdc.Charset)

    cdc.Encoder = func(input []byte) (output string, err error) {

        if len(input) < 1 {

            err = proto.Error_ZERO_LENGTH
            return
        }

        // The check length depends on the modulus of the length of the data is
        // order to avoid padding.
        checkLen := getCheckLen(len(input))

        // The output is longer than the input, so we create a new buffer.
        outputBytes := make([]byte, len(input)+checkLen+1)

        // Add the check length byte to the front
        outputBytes[0] = byte(checkLen)

        // Then copy the input bytes for beginning segment.
        copy(outputBytes[1:len(input)+1], input)

        // Then copy the check to the end of the input.
        copy(outputBytes[len(input)+1:], cdc.MakeCheck(input, checkLen))

        // Create the encoding for the output.
        outputString := enc.EncodeToString(outputBytes)

        // We can omit the first character of the encoding because the length
        // prefix never uses the first 5 bits of the first byte, and add it back
        // for the decoder later.
        trimmedString := outputString[1:]

        // Prefix the output with the Human Readable Part and append the
        // encoded string version of the provided bytes.
        output = cdc.HRP + trimmedString

        return
    }
```

The comments explain every step in the process. 

[->contents](#kitchensink)

#### About `make()`

First, as this is the first point at which the builting function `make` appears, this is a function that initialises and populates the three main types of "reference types" (values that are actually pointers), slice (`[]T`), map (`map[K]V`) and channel (`chan T`). 

All of these types have a top level "value" structure which contains pointers and state data such as the underlying memory location, length and capacity of slices, the number of members and the location of the index of maps, and the underlying structure that manages channels (more about these later).

Overall, you will use and encounter slices the most. Slices can be populated with literals, creating the equivalent of 'constant' values, which are not actually constant in Go, as they involve pointers which are strictly not constant, for which reason they are not immutable. 

Adding immutability creates an extra attribute for all variables which would mean on the underlying implementation, *everything* is then a "reference type" and every access requires an indirection to account for the hidden values like "mutable". In practice, this actually is a performance benefit because it eliminates compiler complexity and makes performance tuning more predictable for the programmer.

Immutability was left out of Go because it can be mimicked using encapsulation and non-exported symbols. Go does not include anything that can be constructed from what is already available. All the things that you will miss from other languages, Go will show you just how expensive they are to implement and make you understand why Go does not have them for all purposes.

The `make` function must be invoked before using any of these structured types otherwise the variable is `nil`.

With slices, it is possible to invoke `make` without a length parameter, and you will be able to use the `append` function to add elements to the slice. You can also initalise a slice with an empty literal:

​    newSlice := []byte{}

We will not show you `append` even once in this tutorial because actually, it's a bad idea to use it in production, as it costs memory allocations and very often, *heap* allocations, which are the most expensive. `make` also requires heap allocations, but we promote the idea of using it with length values for slices to put the allocation during initialisation rather than during the runtime, when this will create latency in the implementation, and is preferable to be done before loops as well, so only one allocation is done instead of potentially hundreds during an initialisation loop.

It is acceptable to delegate this allocation logic for some purposes to the runtime, but if performance is critical, then it is always best to derive some kind of estimate if an exact precalculation is not possible, to reduce possible time consuming heap allocations during tight loops.

> **The syntax for slices and copying requires some explanation.**
>
> In Go, to copy slices (variable types with the `[]` prefix) you use the `copy` operator, the assignment operator `=` copies the *value* of the slice, which is actually a struct, internally, containing the pointer to the memory, the current length used and the capacity (which can be larger than current used), which would not achieve the goal of creating a new, expanded version including the check length prefix and check/padding, it just creates a new reference to the same old data, and modifying this will still then need to make a copy, and, any modifications that don't cause a reallocation and copy will be visible to other threads with this same pointer to the slice.
>
> The other point to note is the use of the slicing operator. Individual elements are addressed via `variableName[n]` and to designate a subsection of the slice you use `variableName[start:end]` where the values represent the first element and the element *after* the last element.
>
> It is important to explain this as the notation can be otherwise confusing. The end index in particular needs to be understood as *exclusive* not *inclusive*. The length function `len(sliceName)` returns a value that is the same as the index that you would use to designate *up to the end* of the slice, as it is the cardinal value (count) where the ordinal (index) starts from zero and is thus one less.
>
> Note that if you omit either of these values, as in the "slicing" operation used to convert arrays into slices:
>
> ​    sliceName[:]
>
> the unspecified values default to zero and the last index, which is the currently allocated length of the slice.
>
> Lastly, the slicing operator can also be used on strings, but beware that the indexes are bytes, and do not respect the character boundaries of UTF-8 encoding, which is only one byte per character for the first 255 characters of ASCII and does not include any (many, it does include several umlauts and accent characters from european languages) non-latin symbols, There is, however, several competing simple 8 bit character maps other than ASCII, that have graphical characters like the MS DOS one, and several others for other alphabetic scripts like japanese, korean and cyrillic.
>
> If you need to do UTF-8 text processing you need to use iterator functions that detect the variable length characters and convert them to 32 bit values. We don't need that here because we are using pure ASCII 8 bit characters, to be exact, 7 bits with a zero bit pad.
>
> In the case of the Base32 encoding, we are using standard 7 bit ASCII symbols so we know that we can cut off the first one to remove the redundant zero that appears because of the maximum 3 bits used for the check length prefix value, leave 5 bits in front (due to the backwards encoding convention for numbers within machine words). 
>
> *whew* A lot to explain about the algorithm above, but vital to understand for anyone who wants to work with slices of bytes in Go, which basically means anything involving binary encoding. This will be as deeply technical as this tutorial gets, it's not essential to understand it to do the tutorial, but this explanation is added for the benefit of those who do or will need to work with binary encoded data.

[->contents](#kitchensink)

#### Creating the Check function

Continuing in the promised logical, first principles ordering of things, the next thing we need is the Check function. 

The check function essentially prepares the input string correctly, removing the Human Readable Prefix string, appends the prefix zero that is cut off due to always being zero, decodes it from Base32 into bytes, slices off the check length prefix byte, separates the data from the check bytes, then it runs the `MakeCheck` function to make sure the check bytes match what was provided. 

If the match fails, the check is incorrect or the data is corrupt. In this case, incorrectly transcribed.

First thing we need to run the check function is a function that calculates the index the check byte starts at:

```go
// getCutPoint is made into a function because it is needed more than once.
func getCutPoint(length, checkLen int) int {

    return length - checkLen - 1
}
```

This is a very simple formula, but it needs to be used again in the decoder function where it allows the raw bytes to be correctly cut to return the checked value. It just needs the checklen, and further includes the prefix check length byte which is sliced off before using this value.

The following function assumes that the decoding from Base32 to bytes has already been done correctly, but nothing more:

```go
    cdc.Check = func(input []byte) (err error) {

        // We must do this check or the next statement will cause a bounds check
        // panic. Note that zero length and nil slices are different, but have
        // the same effect in this case, so both must be checked.
        switch {
        case len(input) < 1:

            err = proto.Error_ZERO_LENGTH
            return

        case input == nil:

            err = proto.Error_NIL_SLICE
            return
        }

        // The check length is encoded into the first byte in order to ensure
        // the data is cut correctly to perform the integrity check.
        checkLen := int(input[0])

        // Ensure there is at enough bytes in the input to run a check on
        if len(input) < checkLen+1 {

            err = proto.Error_CHECK_TOO_SHORT
            return
        }

        // Find the index to cut the input to find the checksum value. We need
        // this same value twice so it must be made into a variable.
        cutPoint := getCutPoint(len(input), checkLen)

        // Here is an example of a multiple assignment and more use of the
        // slicing operator.
        payload, checksum := input[1:cutPoint], string(input[cutPoint:])

        // A checksum is checked in all cases by taking the data received, and
        // applying the checksum generation function, and then comparing the
        // checksum to the one attached to the received data with checksum
        // present.
        //
        // Note: The casting to string above and here. This makes a copy to the
        // immutable string, which is not optimal for large byte slices, but for
        // this short check value, it is a cheap operation on the stack, and an
        // illustration of the interchangeability of []byte and string, with the
        // distinction of the availability of a comparison operator for the
        // string that isn't present for []byte, so for such cases this
        // conversion is a shortcut method to compare byte slices.
        computedChecksum := string(cdc.MakeCheck(payload, checkLen))

        // Here we assign to the return variable the result of the comparison.
        // by doing this instead of using an if and returns, the meaning of the
        // comparison is more clear by the use of the return value's name.
        valid := checksum != computedChecksum

        if !valid {

            err = proto.Error_CHECK_FAILED
        }

        return
    }
```

Take note about the use of the string cast above. In Go, slices do not have an equality operator `==` but and arrays (fixed length with a constant length value like this `[32]byte` as opposed to `[]byte` for the slice) and strings do. 

Note that you cannot create an array using a variable. Only constants and literals can be used to do this.

Casting bytes to string creates an immutable copy so it adds a copy operation. If the amount of data is very large, you write a custom comparison function to avoid this duplication, but for short amounts of data, the extra copy stays on the stack and does not take a lot of time, in return for the simplified comparison as shown above.

[->contents](#kitchensink)

#### Creating the Decoder function

The decoder cuts off the HRP, prepends the always zero first base32 character, decodes using the Base32 encoder (it is created prior to the encode function previously, and is actually a codec, though I used the name `enc`, it also has a decode function).

```go
    cdc.Decoder = func(input string) (output []byte, err error) {

        // Other than for human identification, the HRP is also a validity
        // check, so if the string prefix is wrong, the entire value is wrong
        // and won't decode as it is expected.
        if !strings.HasPrefix(input, cdc.HRP) {

            log.Printf("Provided string has incorrect human readable part:"+
                "found '%s' expected '%s'", input[:len(cdc.HRP)], cdc.HRP,
            )

            err = proto.Error_INCORRECT_HUMAN_READABLE_PART

            return
        }

        // Cut the HRP off the beginning to get the content, add the initial
        // zeroed 5 bits with a 'q' character.
        //
        // Be aware the input string will be copied to create the []byte
        // version. Also, because the input bytes are always zero for the first
        // 5 most significant bits, we must re-add the zero at the front (q)
        // before feeding it to the decoder.
        input = "q" + input[len(cdc.HRP):]

        // The length of the base32 string refers to 5 bits per slice index
        // position, so the correct size of the output bytes, which are 8 bytes
        // per slice index position, is found with the following simple integer
        // math calculation.
        //
        // This allocation needs to be made first as the base32 Decode function
        // does not do this allocation automatically and it would be wasteful to
        // not compute it precisely, when the calculation is so simple.
        //
        // If this allocation is omitted, the decoder will panic due to bounds
        // check error. A nil slice is equivalent to a zero length slice and
        // gives a bounds check error, but in fact, the slice has no data at
        // all. Yes, the panic message is lies:
        //
        //   panic: runtime error: index out of range [4] with length 0
        //
        // If this assignment isn't made, by default, output is nil, not
        // []byte{} so this panic message is deceptive.
        data := make([]byte, len(input)*5/8)

        var writtenBytes int
        writtenBytes, err = enc.Decode(data, []byte(input))
        if err != nil {

            log.Println(err)
            return
        }

        // The first byte signifies the length of the check at the end
        checkLen := int(data[0])
        if writtenBytes < checkLen+1 {

            err = proto.Error_CHECK_TOO_SHORT

            return
        }

        // Assigning the result of the check here as if true the resulting
        // decoded bytes still need to be trimmed of the check value (keeping
        // things cleanly separated between the check and decode function.
        err = cdc.Check(data)

        // There is no point in doing any more if the check fails, as per the
        // contract specified in the interface definition codecer.Codecer
        if err != nil {
            return
        }

        // Slice off the check length prefix, and the check bytes to return the
        // valid input bytes.
        output = data[1:getCutPoint(len(data)+1, checkLen)]

        // If we got to here, the decode was successful.
        return
    }
```

Note that in a couple of places there are log prints. This is because when developing this, it was critical to be able to see what exactly was incorrect, in tests where the result was wrong.

[->contents](#kitchensink)

----

### [Step 5](steps/step5) Testing the algorithm

I am not a fan of doing a lot of stuff before I put the thing into action, there is so many mistakes that can happen without the feedback. Go in almost every respect is a language designed for short feedback loops, and up until the previous step we were only writing mainly declarations so now that we have possibly working code, it needs to be exercised as soon as possible.

The first way and best way to see the rubber meet the road in Go is writing tests. Tests make sure that changes don't break the code, once it's working, but they also are the best way to debug the code in the first place, without having to put them into use. 

Tests are not necessarily going to prove everything is correct, for this there is more advanced techniques to work with more complex algorithms than we have written in this tutorial. For this test we just need a decently broad set of inputs that we can ensure are encoded correctly, and that are decoded correctly.

For this, rather than creating truly random inputs, we exploit the determinism of the `math/rand` library to generate a deterministic set of inputs that will be variable enough to prove the code is working.

[->contents](#kitchensink)

#### Random generation of test data

What we now do is use the random number generator to create a small set of random values that serve as the base to generate 256 bit values to encode, which would be typical data for this type of encoder, hash values for cryptocurrency addresses or transaction hash identifiers.

We define a random seed, and from this run a loop to generate the defined number of random values, and then hash them with Blake3 hash function again... since there is no reason why to use a different library and clutter up our imports. 

In general most new crypto projects are using blake3 because of its performance and security, and SHA256 to deal with legacy hashes from older projects. Blake2 makes some appearances here and there but it was superseded only a few years later.

IPFS uses SHA256 hashes also, but this is partly due to the fact that there is ample hardware to accelerate this hash function and the security of the identity of IPFS data is not as monetarily important as cryptocurrency.

Once the hashes are generated, we create a small code generator that prints them to the console in the test. We will explain that in a minute. For now just copy this code as it is into `pkg/based32/based32_test.go` and then afterwards we will explain more.

```go
package based32

import (
    "encoding/binary"
    "encoding/hex"
    "fmt"
    "lukechampine.com/blake3"
    "math/rand"
    "testing"
)

const (
    seed    = 1234567890
    numKeys = 32
)

func TestCodec(t *testing.T) {

    // Generate 10 pseudorandom 64 bit values. We do this here rather than
    // pre-generating this separately as ultimately it is the same thing, the
    // same seed produces the same series of pseudorandom values, and the hashes
    // of these values are deterministic.
    rand.Seed(seed)
    seeds := make([]uint64, numKeys)
    for i := range seeds {

        seeds[i] = rand.Uint64()
    }

    // Convert the uint64 values to 8 byte long slices for the hash function.
    seedBytes := make([][]byte, numKeys)
    for i := range seedBytes {

        seedBytes[i] = make([]byte, 8)
        binary.LittleEndian.PutUint64(seedBytes[i], seeds[i])
    }

    // Generate hashes from the seeds
    hashedSeeds := make([][]byte, numKeys)

    // Uncomment lines relating to this variable to regenerate expected data
    // that will log to console during test run.
    generated := "\nexpected := []string{\n"

    for i := range hashedSeeds {

        hashed := blake3.Sum256(seedBytes[i])
        hashedSeeds[i] = hashed[:]

        generated += fmt.Sprintf("\t\"%x\",\n", hashedSeeds[i])
    }

    generated += "}\n"
    t.Log(generated)

    expected := []string{
        "7bf4667ea06fe57687a7c0c8aae869db103745a3d8c5dce5bf2fc6206d3b97e4",
        "84c0ee2f49bfb26f48a78c9d048bb309a006db4c7991ebe4dd6dc3f2cc1067cd",
        "206a953c4ba4f79ffe3d3a452214f19cb63e2895866cc27c7cf6a4ec8fe5a7a6",
        "35d64c401829c621624fe9d4f134c24ae909ecf4f07ec4540ffd58911f427d03",
        "573d6989a2c2994447b4669ae6931f12e73c8744e9f65451918a1f3d8cd39aa1",
        "2b08aea58cc1d680de0e7acadc027ebe601f923ff9d5536c6f73e2559a1b6b14",
        "bcc3256005da59b06f69b4c1cc62c89af041f8cd5ad79b81351fbfbbaf2cc60f",
        "42a0f7b9aef1cdc0b3f2a1fd0fb547fb76e5eb50f4f5a6646ee8929fdfef5db7",
        "50e1cb9f5f8d5325e18298faeeea7fd93d83e3bd518299e7150c1f548c11ddc8",
        "22a70a74ccfd61a47576150f968039cfeb33143ec549dfeb6c95afc8a6d3d75a",
        "cd1b21bd745e122f0db1f5ca4e4cbe0bace8439112d519e5b9c0a44a2648a61a",
        "9ec4bd670b053e722d9d5bc3c2aca4a1d64858ef53c9b58a9222081dc1eeb017",
        "d85713c898f2fc95d282b58a0ea475c1b1726f8be44d2f17c3c675ce688e1563",
        "875baee7e9a372fe3bad1fbf0e2e4119038c1ed53302758fc5164b012e927766",
        "7de94ca668463db890478d8ba3bbed35eb7666ac0b8b4e2cb808c1cbb754576f",
        "159469150dc41ebd2c2bfafb84aef221769699013b70068444296169dc9e93be",
        "a90a104ea470df61d337c51589b520454acbd05ef5bbe7d2a8285043a222bec9",
        "a835de5206f6dbef6a2cb3da66ffb99a19bfa4e005208ffdb316ce880132297e",
        "f6a09e8f41231febd1b25c52cb73ea438ac803db77d5549db4e15a32e804de9f",
        "074c59cce7783042cc6941c849206582ecc43028d1576d00e02d95b1e669bf7a",
        "203c3566724c229b570f33be994cd6094e1a64f3df552f1390b4c2adc7e36d6d",
        "efec32d52a17ed75ad5a486ba621e0f47f61e4e60557129fce728a1bb63208fd",
        "9cc2962fc62fe40f6197a4fb81356717fd57b4c988641bca3a9d45efde893894",
        "2adf211300632bb5f650202bf128ba9187ec2c6c738431dc396d93b8f62bd590",
        "0782aade40d0ae7a293bfb67016466682d858b5226eaaa8df2f2104fa6c408c3",
        "d011ad5550f3f03caa469fa233f553721e6af84f1341d256cefe052d85397637",
        "83deb64f5c134d108e8b99c8a196b8d04228acfc810c33711d975400fa731508",
        "d9a4b19142d015fd541f50f18f41b7e9738a30c59a3e914b4d4d1556c75786f2",
        "3e05940b76735ea114db8b037dece53090765510c9c4e55a0be18cb8aef754fa",
        "41f43119041dd1f3a250f54768ce904808cd0d7bb7b37697803ed2940c39a555",
        "a2c2d7cb980c2b57c8fdfae55cf4c6040eaf8163b21072877e5e57349388d59c",
        "02155c589e5bd89ce806a33c1841fe1e157171222701d515263acd0254208a39",
    }

    for i := range hashedSeeds {

        if expected[i] != hex.EncodeToString(hashedSeeds[i]) {

            t.Log("failed", i, "expected", expected[1], "found", hashedSeeds)
            t.FailNow()
        }
    }
}
```

First thing to explain is how tests are constructed.

1. The filename should be named the same as the code it tests, in general, the code we are testing is in `based32.go` and so the convention is the test for this source code should be `based32_test.go`. 
2. The signature of a test function should be `func TestSomething(t *testing.T) `where "Something" is the name of the type with the methods we are testing, or the specific name of the methods of the type we are testing, or something along these lines. It is not mandatory to use existing names, but it is recommended.

[->contents](#kitchensink)

#### Running the tests

To run the tests, you simply need to designate the package folder where the tests you want to run are living:

    davidvennik@pop-os:~/src/github.com/quanterall/kitchensink/steps/step5$ go test ./pkg/based32/
    ok      github.com/quanterall/kitchensink/steps/step5/pkg/based32       (cached)

Note the `./` in front of the folder. This means "in the current directory" in Unix file path convention, and the reason why we have to specify it is because otherwise, the `go` tool expects to be given a network path like `github.com/quanterall/kitchensink/pkg/based32` which would then clone the repository, or use a cached version, and then run the tests in that directory for you.

[->contents](#kitchensink)

#### Enabling logging in the tests

In order to update automatically generated, deterministic data like we are using in this test, we need the test to print out the data that it should be generating, according to what we believe is correct code. 

If you ran the command as above in the previous section, you would have seen that it passed, assuming you didn't mess up the code while copy and pasting it.

However, as you will know if you read closely in the text, there is a section of the code that invokes `t.Log()` which as you would expect, will print something to the terminal. To make this happen, you use the `-v` flag after the subcommand `test`.

[->contents](#kitchensink)

#### The Go tool recursive descent notation

Before we go any further, we will introduce a special path notation that the `go` tool understands that automatically recursively descends a filesystem hierarchy from a given starting point.

    ./...

This means, from the given path `./` which means the current working directory, to recursively descend into the folders and run the operation requested on each relevant file encountered.

In old unix file commands, there is disagreement between numerous different commands how to indicate to perform this operation. For `ls` it is `-r`, for `cp` it is `-r` but for `scp` (SSH copy) it is `-R`. 

To make things more consistent and some 'grandfatherly' guidance from the guys who actually first invented Unix (Ken Thompson and Rob Pike were both involved in the original invention of Unix at Bell Labs, Pike was mostly involved with documenting the C language), the Go tool introduces the `...` path to indicate the same thing. 

Go programmers pretty much like using it whenever they make tools to do this sort of thing, and in Go the recursive descent of filesystems is a single function that you load with a closure.

[->contents](#kitchensink)

#### Running the tests with logging and recursive descent

Since when you are doing tests, most of the time the codebase is not big enough to be specific, this can be memorised as the one main way to run tests on a Go repository.

    go test -v ./...

and you will get a result like this on the current repository you are working on:

```
davidvennik@pop-os:~/src/github.com/quanterall/kitchensink/steps/step5$ go test -v ./...
?       github.com/quanterall/kitchensink/steps/step5   [no test files]
=== RUN   TestCodec
    based32_test.go:54: 
        expected := []string{
                "7bf4667ea06fe57687a7c0c8aae869db103745a3d8c5dce5bf2fc6206d3b97e4",
                "84c0ee2f49bfb26f48a78c9d048bb309a006db4c7991ebe4dd6dc3f2cc1067cd",
                "206a953c4ba4f79ffe3d3a452214f19cb63e2895866cc27c7cf6a4ec8fe5a7a6",
                "35d64c401829c621624fe9d4f134c24ae909ecf4f07ec4540ffd58911f427d03",
                "573d6989a2c2994447b4669ae6931f12e73c8744e9f65451918a1f3d8cd39aa1",
                "2b08aea58cc1d680de0e7acadc027ebe601f923ff9d5536c6f73e2559a1b6b14",
                "bcc3256005da59b06f69b4c1cc62c89af041f8cd5ad79b81351fbfbbaf2cc60f",
                "42a0f7b9aef1cdc0b3f2a1fd0fb547fb76e5eb50f4f5a6646ee8929fdfef5db7",
                "50e1cb9f5f8d5325e18298faeeea7fd93d83e3bd518299e7150c1f548c11ddc8",
                "22a70a74ccfd61a47576150f968039cfeb33143ec549dfeb6c95afc8a6d3d75a",
                "cd1b21bd745e122f0db1f5ca4e4cbe0bace8439112d519e5b9c0a44a2648a61a",
                "9ec4bd670b053e722d9d5bc3c2aca4a1d64858ef53c9b58a9222081dc1eeb017",
                "d85713c898f2fc95d282b58a0ea475c1b1726f8be44d2f17c3c675ce688e1563",
                "875baee7e9a372fe3bad1fbf0e2e4119038c1ed53302758fc5164b012e927766",
                "7de94ca668463db890478d8ba3bbed35eb7666ac0b8b4e2cb808c1cbb754576f",
                "159469150dc41ebd2c2bfafb84aef221769699013b70068444296169dc9e93be",
                "a90a104ea470df61d337c51589b520454acbd05ef5bbe7d2a8285043a222bec9",
                "a835de5206f6dbef6a2cb3da66ffb99a19bfa4e005208ffdb316ce880132297e",
                "f6a09e8f41231febd1b25c52cb73ea438ac803db77d5549db4e15a32e804de9f",
                "074c59cce7783042cc6941c849206582ecc43028d1576d00e02d95b1e669bf7a",
                "203c3566724c229b570f33be994cd6094e1a64f3df552f1390b4c2adc7e36d6d",
                "efec32d52a17ed75ad5a486ba621e0f47f61e4e60557129fce728a1bb63208fd",
                "9cc2962fc62fe40f6197a4fb81356717fd57b4c988641bca3a9d45efde893894",
                "2adf211300632bb5f650202bf128ba9187ec2c6c738431dc396d93b8f62bd590",
                "0782aade40d0ae7a293bfb67016466682d858b5226eaaa8df2f2104fa6c408c3",
                "d011ad5550f3f03caa469fa233f553721e6af84f1341d256cefe052d85397637",
                "83deb64f5c134d108e8b99c8a196b8d04228acfc810c33711d975400fa731508",
                "d9a4b19142d015fd541f50f18f41b7e9738a30c59a3e914b4d4d1556c75786f2",
                "3e05940b76735ea114db8b037dece53090765510c9c4e55a0be18cb8aef754fa",
                "41f43119041dd1f3a250f54768ce904808cd0d7bb7b37697803ed2940c39a555",
                "a2c2d7cb980c2b57c8fdfae55cf4c6040eaf8163b21072877e5e57349388d59c",
                "02155c589e5bd89ce806a33c1841fe1e157171222701d515263acd0254208a39",
        }
        
--- PASS: TestCodec (0.00s)
PASS
ok      github.com/quanterall/kitchensink/steps/step5/pkg/based32       0.002s
?       github.com/quanterall/kitchensink/steps/step5/pkg/codecer       [no test files]
?       github.com/quanterall/kitchensink/steps/step5/pkg/proto [no test files]
```

And now, as you can see, our test passes, and if you look a little closer, you can see that not only does the test pass, you also see the output here is identical to a section of the code we just entered.

The reason for that, is that we want to use this property to enable us to create, automatically, a decent swathe of test inputs to use for the tests, and all we have written right now is a test that makes sure and demonstrates that the random function does indeed generate completely deterministic values, as does the hash function that we feed these deterministic values.

We don't need to test the values are identical, because they have to be, by their definition, as does the hash, so we use the hash values, as from these we can generate a set of tests that make sure our implementation of `Codec` works as it is supposed to be.

[->contents](#kitchensink)

#### Actually testing the Encoder and Decoder

To test everything that the encoder does correctly, we need to feed it variable lengths of data. 

For this we will cut those hashes to different lengths by subtracting the modulus (remainder in long division) of the sequence number of the item from the length of the data. Thus for 0 and multiples of 5, we leave them as is, for everything else, we remove as many bytes as the remainder indicates.

In this way we can be sure that the variable length check bytes and check length prefixes work correctly. There is no real need to test it further than this though one *could*, it would be redundant since the algorithm only varies the length by a difference between multiples of 5 bytes.

So, add this at the end of the test function 

```go

	encodedStr := []string{
		"QNTRLfalgen75ph72a585lqv32hgd8d3qd6950vvth89huhuvgrd8wt7ftet",
		"QNTRLwzvpm30fxlmym6g57xf6pytkvy6qpkmf3uer6lym4ku8ukvzpnckqs6",
		"QNTRLssx49fufwj008l785ay2gs57xwtv03gjkrxesnu0nm2fmy0uhmerj80",
		"QNTRL56avnzqrq5uvgtzfl5afuf5cf9wjz0v7nc8a3z5pl743yglgjs6nypp",
		"QNTRLetn66vf5tpfj3z8k3nf4e5nrufww0y8gn5lv4z3jx9p70vwrzchx7pf",
		"QNTRLg4s3t493nqadqx7peav4hqz06lxq8uj8lua25mvdae7y4v6rd43fmcv",
		"QNTRLw7vxftqqhd9nvr0dx6vrnrzezd0qs0ce4dd0xupx50mlwa09nr8hhvc",
		"QNTRL3p2paae4mcums9n72sl6ra4glahde0t2r60tfnydm5f987lalykzuqe",
		"QNTRL4gwrjult7x4xf0ps2v04mh20lvnmqlrh4gc9x08z5xp74yva49qf29d",
		"QNTRLc32wzn5en7krfr4wc2sl95q8887kvc58mz5nhltdj26ljxy33yxs7px",
		"QNTRLtx3kgdaw30pytcdk86u5njvhc96e6zrjyfd2x09h8q2gj3xfznp4l5y",
		"QNTRLw0vf0t8pvznuu3dn4du8s4v5jsavjzcaafundv2jg3qs8wpa6czu2kz",
		"QNTRLnv9wy7gnre0e9wjs26c5r4ywhqmzun030jy6tchc0r8tnng3e6u5pg6",
		"QNTRLkr4hth8ax3h9l3m450m7r3wgyvs8rq765esyav0c5tykqfwr2u7zmfp",
		"QNTRLe77jn9xdprrmwysg7xchgama567kanx4s9ckn3vhqyvrjam6x9h7qmx",
		"QNTRLg2eg6g4phzpa0fv90a0hp9w7gshd95eqyahqp5ygs5kz6wun6fmukle",
		"QNTRLw5s5yzw53cd7cwnxlz3tzd4ypz54j7stm6mhe7j4q59qsazy2lfqq35",
		"QNTRLj5rthjjqmmdhmm29jea5ehlhxdpn0ayuqzjprlakvtvazqpxt0zdxtv",
		"QNTRLhm2p850gy33l673kfw99jmnafpc4jqrmdma24yakns45vhg0sfcsrrp",
		"QNTRLcr5ckwvuaurqskvd9qusjfqvkpwe3ps9rg4wmgquqketvvex8jy6alw",
		"QNTRLgsrcdtxwfxz9x6hpuemax2v6cy5uxny70042tcnjz6v9tw8udkk6rkl",
		"QNTRL0h7cvk49gt76addtfyxhf3pur687c0yucz4wy5leeeg5xakxgywasu3",
		"QNTRLjwv9930cch7grmpj7j0hqf4vutl64a5exyxgx7282w5tm7738hap8u7",
		"QNTRL54d7ggnqp3jhd0k2qszhufgh2gc0mpvd3ecgvwu89ke8w8kcv3ksyma",
		"QNTRLcrc92k7grg2u73f80akwqtyve5zmpvt2gnw425d7tepqnaz7jehh549",
		"QNTRLtgprt242relq092g606yvl42depu6hcfuf5r5jkemlq2tv989mr0wng",
		"QNTRLwpaadj0tsf56yyw3wvu3gvkhrgyy29vljqscvm3rkt4gq86wv2upf6r",
		"QNTRLnv6fvv3gtgptl25rag0rr6pkl5h8z3sckdray2tf4x324k82a58c5kf",
		"QNTRL5lqt9qtwee4agg5mw9sxl0vu5cfqaj4zryufe26p0scew9w9cgcnqj8",
		"QNTRLeqlgvgeqswaruaz2r65w6xwjpyq3ngd0wmmxa5hsqld998rns3l4gfz",
		"QNTRL23v947tnqxzk47glhaw2h85cczqatupvwepqu580e09wdyn3r2eeltr",
		"QNTRLvpp2hzcneda388gq63ncxzplc0p2ut3ygnsr4g4ycav6qj5yz9g9lv8",
	}

	encoded := "\nencodedStr := []string{\n"

	// Convert hashes to our base32 encoding format
	for i := range hashedSeeds {

		// Note that we are slicing off a number of bytes at the end according
		// to the sequence number to get different check byte lengths from a
		// uniform original data. As such, this will be accounted for in the
		// check by truncating the same amount in the check (times two, for the
		// hex encoding of the string).
		encode, err := Codec.Encode(hashedSeeds[i][:len(hashedSeeds[i])-i%5])
		if err != nil {
			t.Fatal(err)
		}
		if encode != encodedStr[i] {
			t.Fatalf(
				"Decode failed, expected '%s' got '%s'",
				encodedStr, encode,
			)
		}
		encoded += "\t\"" + encode + "\",\n"
	}

	encoded += "}\n"
	t.Log(encoded)

	// Next, decode the encodedStr above, which should be the output of the
	// original generated seeds, with the index mod 5 truncations performed on
	// each as was done to generate them.

	for i := range encodedStr {

		res, err := Codec.Decode(encodedStr[i])
		if err != nil {
			t.Fatalf("error: '%v'", err)
		}
		elen := len(expected[i])
		etrimlen := 2 * (i % 5)
		expectedHex := expected[i][:elen-etrimlen]
		resHex := fmt.Sprintf("%x", res)
		if resHex != expectedHex {
			t.Fatalf(
				"got: '%s' expected: '%s'",
				resHex,
				expectedHex,
			)
		}
	}
```

Note that it algorithmically trims the bytes of intput according to modulus 5, and performs the equivalent trim on the hex strings for the output to enable simple comparison of the bytes. The raw bytes this is `%5` but for hex it is multiplied by 2 as each hexadecimal character represents half a byte.

Running `go test -v ./...` will then yield this new result:

```
davidvennik@pop-os:~/src/github.com/quanterall/kitchensink/steps/step5$ go test -v ./...
?       github.com/quanterall/kitchensink/steps/step5   [no test files]
=== RUN   TestCodec
    based32_test.go:54: 
        expected := []string{
                "7bf4667ea06fe57687a7c0c8aae869db103745a3d8c5dce5bf2fc6206d3b97e4",
                "84c0ee2f49bfb26f48a78c9d048bb309a006db4c7991ebe4dd6dc3f2cc1067cd",
                "206a953c4ba4f79ffe3d3a452214f19cb63e2895866cc27c7cf6a4ec8fe5a7a6",
                "35d64c401829c621624fe9d4f134c24ae909ecf4f07ec4540ffd58911f427d03",
                "573d6989a2c2994447b4669ae6931f12e73c8744e9f65451918a1f3d8cd39aa1",
                "2b08aea58cc1d680de0e7acadc027ebe601f923ff9d5536c6f73e2559a1b6b14",
                "bcc3256005da59b06f69b4c1cc62c89af041f8cd5ad79b81351fbfbbaf2cc60f",
                "42a0f7b9aef1cdc0b3f2a1fd0fb547fb76e5eb50f4f5a6646ee8929fdfef5db7",
                "50e1cb9f5f8d5325e18298faeeea7fd93d83e3bd518299e7150c1f548c11ddc8",
                "22a70a74ccfd61a47576150f968039cfeb33143ec549dfeb6c95afc8a6d3d75a",
                "cd1b21bd745e122f0db1f5ca4e4cbe0bace8439112d519e5b9c0a44a2648a61a",
                "9ec4bd670b053e722d9d5bc3c2aca4a1d64858ef53c9b58a9222081dc1eeb017",
                "d85713c898f2fc95d282b58a0ea475c1b1726f8be44d2f17c3c675ce688e1563",
                "875baee7e9a372fe3bad1fbf0e2e4119038c1ed53302758fc5164b012e927766",
                "7de94ca668463db890478d8ba3bbed35eb7666ac0b8b4e2cb808c1cbb754576f",
                "159469150dc41ebd2c2bfafb84aef221769699013b70068444296169dc9e93be",
                "a90a104ea470df61d337c51589b520454acbd05ef5bbe7d2a8285043a222bec9",
                "a835de5206f6dbef6a2cb3da66ffb99a19bfa4e005208ffdb316ce880132297e",
                "f6a09e8f41231febd1b25c52cb73ea438ac803db77d5549db4e15a32e804de9f",
                "074c59cce7783042cc6941c849206582ecc43028d1576d00e02d95b1e669bf7a",
                "203c3566724c229b570f33be994cd6094e1a64f3df552f1390b4c2adc7e36d6d",
                "efec32d52a17ed75ad5a486ba621e0f47f61e4e60557129fce728a1bb63208fd",
                "9cc2962fc62fe40f6197a4fb81356717fd57b4c988641bca3a9d45efde893894",
                "2adf211300632bb5f650202bf128ba9187ec2c6c738431dc396d93b8f62bd590",
                "0782aade40d0ae7a293bfb67016466682d858b5226eaaa8df2f2104fa6c408c3",
                "d011ad5550f3f03caa469fa233f553721e6af84f1341d256cefe052d85397637",
                "83deb64f5c134d108e8b99c8a196b8d04228acfc810c33711d975400fa731508",
                "d9a4b19142d015fd541f50f18f41b7e9738a30c59a3e914b4d4d1556c75786f2",
                "3e05940b76735ea114db8b037dece53090765510c9c4e55a0be18cb8aef754fa",
                "41f43119041dd1f3a250f54768ce904808cd0d7bb7b37697803ed2940c39a555",
                "a2c2d7cb980c2b57c8fdfae55cf4c6040eaf8163b21072877e5e57349388d59c",
                "02155c589e5bd89ce806a33c1841fe1e157171222701d515263acd0254208a39",
        }
        
    based32_test.go:159: 
        encodedStr := []string{
                "QNTRLfalgen75ph72a585lqv32hgd8d3qd6950vvth89huhuvgrd8wt7ftet",
                "QNTRLwzvpm30fxlmym6g57xf6pytkvy6qpkmf3uer6lym4ku8ukvzpnckqs6",
                "QNTRLssx49fufwj008l785ay2gs57xwtv03gjkrxesnu0nm2fmy0uhmerj80",
                "QNTRL56avnzqrq5uvgtzfl5afuf5cf9wjz0v7nc8a3z5pl743yglgjs6nypp",
                "QNTRLetn66vf5tpfj3z8k3nf4e5nrufww0y8gn5lv4z3jx9p70vwrzchx7pf",
                "QNTRLg4s3t493nqadqx7peav4hqz06lxq8uj8lua25mvdae7y4v6rd43fmcv",
                "QNTRLw7vxftqqhd9nvr0dx6vrnrzezd0qs0ce4dd0xupx50mlwa09nr8hhvc",
                "QNTRL3p2paae4mcums9n72sl6ra4glahde0t2r60tfnydm5f987lalykzuqe",
                "QNTRL4gwrjult7x4xf0ps2v04mh20lvnmqlrh4gc9x08z5xp74yva49qf29d",
                "QNTRLc32wzn5en7krfr4wc2sl95q8887kvc58mz5nhltdj26ljxy33yxs7px",
                "QNTRLtx3kgdaw30pytcdk86u5njvhc96e6zrjyfd2x09h8q2gj3xfznp4l5y",
                "QNTRLw0vf0t8pvznuu3dn4du8s4v5jsavjzcaafundv2jg3qs8wpa6czu2kz",
                "QNTRLnv9wy7gnre0e9wjs26c5r4ywhqmzun030jy6tchc0r8tnng3e6u5pg6",
                "QNTRLkr4hth8ax3h9l3m450m7r3wgyvs8rq765esyav0c5tykqfwr2u7zmfp",
                "QNTRLe77jn9xdprrmwysg7xchgama567kanx4s9ckn3vhqyvrjam6x9h7qmx",
                "QNTRLg2eg6g4phzpa0fv90a0hp9w7gshd95eqyahqp5ygs5kz6wun6fmukle",
                "QNTRLw5s5yzw53cd7cwnxlz3tzd4ypz54j7stm6mhe7j4q59qsazy2lfqq35",
                "QNTRLj5rthjjqmmdhmm29jea5ehlhxdpn0ayuqzjprlakvtvazqpxt0zdxtv",
                "QNTRLhm2p850gy33l673kfw99jmnafpc4jqrmdma24yakns45vhg0sfcsrrp",
                "QNTRLcr5ckwvuaurqskvd9qusjfqvkpwe3ps9rg4wmgquqketvvex8jy6alw",
                "QNTRLgsrcdtxwfxz9x6hpuemax2v6cy5uxny70042tcnjz6v9tw8udkk6rkl",
                "QNTRL0h7cvk49gt76addtfyxhf3pur687c0yucz4wy5leeeg5xakxgywasu3",
                "QNTRLjwv9930cch7grmpj7j0hqf4vutl64a5exyxgx7282w5tm7738hap8u7",
                "QNTRL54d7ggnqp3jhd0k2qszhufgh2gc0mpvd3ecgvwu89ke8w8kcv3ksyma",
                "QNTRLcrc92k7grg2u73f80akwqtyve5zmpvt2gnw425d7tepqnaz7jehh549",
                "QNTRLtgprt242relq092g606yvl42depu6hcfuf5r5jkemlq2tv989mr0wng",
                "QNTRLwpaadj0tsf56yyw3wvu3gvkhrgyy29vljqscvm3rkt4gq86wv2upf6r",
                "QNTRLnv6fvv3gtgptl25rag0rr6pkl5h8z3sckdray2tf4x324k82a58c5kf",
                "QNTRL5lqt9qtwee4agg5mw9sxl0vu5cfqaj4zryufe26p0scew9w9cgcnqj8",
                "QNTRLeqlgvgeqswaruaz2r65w6xwjpyq3ngd0wmmxa5hsqld998rns3l4gfz",
                "QNTRL23v947tnqxzk47glhaw2h85cczqatupvwepqu580e09wdyn3r2eeltr",
                "QNTRLvpp2hzcneda388gq63ncxzplc0p2ut3ygnsr4g4ycav6qj5yz9g9lv8",
        }
        
--- PASS: TestCodec (0.00s)
PASS
ok      github.com/quanterall/kitchensink/steps/step5/pkg/based32       0.003s
?       github.com/quanterall/kitchensink/steps/step5/pkg/codecer       [no test files]
?       github.com/quanterall/kitchensink/steps/step5/pkg/proto [no test files]
```

Again, you can see that all of the elements printed out identically match what we put in the source code.

The two `for` loops we added test encode, and then decode, in sequence, and check that the product of each step matches the inputs, and the expected outputs from each step. 

If there was an error, which you could cause by mangling the `expected` string slice or `encodedStr` array, you would get this outcome, you can do this yourself by putting a wrong character somewhere. The tests will halt immediately as they find an error. So I will mangle the last `encodedStr`:

```
davidvennik@pop-os:~/src/github.com/quanterall/kitchensink/steps/step5$ go test ./pkg/based32/
--- FAIL: TestCodec (0.00s)
    based32_test.go:54: 
        expected := []string{
                "7bf4667ea06fe57687a7c0c8aae869db103745a3d8c5dce5bf2fc6206d3b97e4",
                "84c0ee2f49bfb26f48a78c9d048bb309a006db4c7991ebe4dd6dc3f2cc1067cd",
                "206a953c4ba4f79ffe3d3a452214f19cb63e2895866cc27c7cf6a4ec8fe5a7a6",
                "35d64c401829c621624fe9d4f134c24ae909ecf4f07ec4540ffd58911f427d03",
                "573d6989a2c2994447b4669ae6931f12e73c8744e9f65451918a1f3d8cd39aa1",
                "2b08aea58cc1d680de0e7acadc027ebe601f923ff9d5536c6f73e2559a1b6b14",
                "bcc3256005da59b06f69b4c1cc62c89af041f8cd5ad79b81351fbfbbaf2cc60f",
                "42a0f7b9aef1cdc0b3f2a1fd0fb547fb76e5eb50f4f5a6646ee8929fdfef5db7",
                "50e1cb9f5f8d5325e18298faeeea7fd93d83e3bd518299e7150c1f548c11ddc8",
                "22a70a74ccfd61a47576150f968039cfeb33143ec549dfeb6c95afc8a6d3d75a",
                "cd1b21bd745e122f0db1f5ca4e4cbe0bace8439112d519e5b9c0a44a2648a61a",
                "9ec4bd670b053e722d9d5bc3c2aca4a1d64858ef53c9b58a9222081dc1eeb017",
                "d85713c898f2fc95d282b58a0ea475c1b1726f8be44d2f17c3c675ce688e1563",
                "875baee7e9a372fe3bad1fbf0e2e4119038c1ed53302758fc5164b012e927766",
                "7de94ca668463db890478d8ba3bbed35eb7666ac0b8b4e2cb808c1cbb754576f",
                "159469150dc41ebd2c2bfafb84aef221769699013b70068444296169dc9e93be",
                "a90a104ea470df61d337c51589b520454acbd05ef5bbe7d2a8285043a222bec9",
                "a835de5206f6dbef6a2cb3da66ffb99a19bfa4e005208ffdb316ce880132297e",
                "f6a09e8f41231febd1b25c52cb73ea438ac803db77d5549db4e15a32e804de9f",
                "074c59cce7783042cc6941c849206582ecc43028d1576d00e02d95b1e669bf7a",
                "203c3566724c229b570f33be994cd6094e1a64f3df552f1390b4c2adc7e36d6d",
                "efec32d52a17ed75ad5a486ba621e0f47f61e4e60557129fce728a1bb63208fd",
                "9cc2962fc62fe40f6197a4fb81356717fd57b4c988641bca3a9d45efde893894",
                "2adf211300632bb5f650202bf128ba9187ec2c6c738431dc396d93b8f62bd590",
                "0782aade40d0ae7a293bfb67016466682d858b5226eaaa8df2f2104fa6c408c3",
                "d011ad5550f3f03caa469fa233f553721e6af84f1341d256cefe052d85397637",
                "83deb64f5c134d108e8b99c8a196b8d04228acfc810c33711d975400fa731508",
                "d9a4b19142d015fd541f50f18f41b7e9738a30c59a3e914b4d4d1556c75786f2",
                "3e05940b76735ea114db8b037dece53090765510c9c4e55a0be18cb8aef754fa",
                "41f43119041dd1f3a250f54768ce904808cd0d7bb7b37697803ed2940c39a555",
                "a2c2d7cb980c2b57c8fdfae55cf4c6040eaf8163b21072877e5e57349388d59c",
                "02155c589e5bd89ce806a33c1841fe1e157171222701d515263acd0254208a39",
        }
        
    based32_test.go:150: Decode failed, expected item 31 'QNTRLvpp2hzcneda388gq63ncxzplc5p2ut3ygnsr4g4ycav6qj5yz9g9lv8' got 'QNTRLvpp2hzcneda388gq63ncxzplc0p2ut3ygnsr4g4ycav6qj5yz9g9lv8'
FAIL
FAIL    github.com/quanterall/kitchensink/steps/step5/pkg/based32       0.002s
FAIL
```

In some cases, it makes sense to test multiple failure modes, but for the sake of brevity, we are not testing for failure modes at all, just that a random collection of inputs correctly encodes and decodes, given the extra element of varying the length by up to 5 bytes to account for the variable length check size of 1 to 5 bytes long.

Writing good tests is a bit of a black art, and the task gets more and more complicated the more complex the algorithms.

[->contents](#kitchensink)

----

### [Step 6](steps/step6) Creating a Server

Just to clarify an important distinction, and another aspect of writing modular code, we are not writing an *executable* that you can spawn from the commandline or within scripts or Dockerfiles. We are writing an *in process* service that can be started by any Go application, the actual application that does this will be created later.

The reason why we separate the two things is because for tests, we want to spawn a server and also in parallel run a client to query it.

Because our service uses gRPC/Protobuf for its messages, for the reason that it is binary, and supported by almost every major langage, we will put the service inside `pkg/grpc/server`.

We create the server first for the same reason as we create the encoder first. The entity is apriori in the two part concept. You can't have a client without a server, but you can have a server without a client, as it is impossible to define a request without first having data to request.

[->contents](#kitchensink)

#### The Logger

I will be repetitive about this, because I want to reinforce the point that debugging by log printing is *not* "unprofessional" or any insults that certain kinds of "programmers" would have you believe. Ignore these fools, they have amnesia about their own learning process and want to be elite forever and not have competition.

In the olden days, when code was single threaded and largely procedural, ok, you didn't really need to do this so much because you could just step through the code line by line, set breakpoints, and immediately know where a failure was because the debugger would tell you.

This can probably be said to be true even of languages with heavy threading systems like C++/Boost and Rust/Tokio in that you can rely on the debugger to tell you where your errors are. 

But this simply is not true of Go multithreaded applications. 

The service we are about to explain how to build runs, by default in the binary service daemon we will show you how to build later, 8 concurrent threads to process requests, which run in parallel with a watcher thread that waits for shutdown, using process signals to ensure a clean shutdown when that signal occurs.

This is just a simple application, and has 9 concurrent threads operating. A much more usual situation with writing servers and applications in Go is that there could be 10 or 20 different pieces of code running concurrently with potentially thousands of concurrent threads, and the timing of their interactions can be crucial to the functionality of the application. 

Concurrent programming just cannot be properly debugged only with a debugger, unless that debugger is recording the state of EVERYTHING. Because that is also not practical for performance reasons, log prints let you record application state only when you actually need it to be logged. 

The standard logging library doesn't account for the production versus debug mode logging either, and makes subsystems being enabled explicitly without others, there needs to be better tools, and I have written some such tools but we will not cover this here. For now, just to make sure you know why you should use this little source code location configuration on the standard logger when you are learning, at least.

The only practical way to debug such an application is with logging, with decent precision timestamps so you can correctly play back a post mortem in slow motion to see what actually went wrong, and where.

Put this file in `pkg/grpc/server/log.go` first:

```go
package server

import (
	logg "log"
	"os"
)

var log = logg.New(os.Stderr, "b32" , logg.Llongfile|logg.Lmicroseconds)
```

It is not generally told this way in the Go community, but this was essential to my own progress as a Go programmer, and I think that if it works this way for me, it probably will help at least a substantial number of other developers starting out on their journey who may not have the luxury of the amount of time I was able to finagle in my learning process as an intern, during which time my living conditions were abominable. 

You can gloss over this, if you want to, but you will come back to it if you intend to make any real progress in this business.

This little file makes sure that you can put log prints in anywhere in your code, and when they are printed out, easily jump into the codebase exactly where the problem comes up.

[->contents](#kitchensink)

#### Implementing the worker pool

Continuing, as always, with steps that build upon the previous steps and not leaving you with code that would not compile due to undefined symbols, the very first thing you need to put in the package is the `Transcriber`, which is the name we will give to our worker pool.

So, create a new file `pkg/grpc/server/workerpool.go` and into it we will first put the `Transcriber` data structure, which manages the worker pool:

```go
package server

import (
	"github.com/quanterall/kitchensink/pkg/based32"
	"github.com/quanterall/kitchensink/pkg/proto"
	"go.uber.org/atomic"
	"sync"
)

// transcriber is a multithreaded worker pool for performing transcription encode
// and decode requests. It is not exported because it must be initialised
// correctly.
type transcriber struct {
	stop                       chan struct{}
	encode                     []chan *proto.EncodeRequest
	decode                     []chan *proto.DecodeRequest
	encodeRes                  []chan proto.EncodeRes
	decodeRes                  []chan proto.DecodeRes
	encCallCount, decCallCount *atomic.Uint32
	workers                    uint32
	wait                       sync.WaitGroup
}
```

There is a few explanations that need to be made to start with.

[->contents](#kitchensink)

#### When to not export an externally used type

Here is a case where we are not exporting a type that has methods that will be used by other packages that will import this package. The reason is that channels and slices both need initialisation before they can be used, as performing operations on these variables that have not been initialised will cause a `nil` panic. Thus, we instead will export an initialiser function, which will take care of this initialisation for us.

[->contents](#kitchensink)

#### About Channels

First of all, channels themselves, which have the type `chan` as a prefix. 

Channels are technically known as an Atomic FIFO Queue. What this means is that their operations are atomic, meaning either they happen, or they don't, and no other code sees an in-between state, so they can be safely used concurrently between two or more goroutines.

They are FIFO, meaning First In First Out, in that what comes out, is what was put in first.

The queues can be zero length, or can have multiple buffers, which makes them function more like a queue than a simple signalling or message passing mechanism. Generally one only adds buffers to them to account for latency in processing or to maintain ongoing processing without blocking on the process when it is in continuous operation.

`stop` is often colloquially referred to as a "quit channel" or a "breaker". It is a breaker because with all `chan` types in Go, when you `close(channelName)` them, any time you read from them with `<-channelName` you will get `nil` forever, until the end of the program.

We use the type `struct {}` which is an empty struct, because this is the smallest possible type of data it can be, it is zero bytes! So such a channel can only send simple signals, namely "something", which is actually nothing, it can be used like a momentary switch such as the power button on a computer, or the trigger on a gun, it simply sends a signal to do something to whatever it is connected to.

So, such channels have to be single purpose. If you need to send only one signal to stop, you close the channel, if you want to use the channel to send multiple triggers, you send an empty struct value to it. You *could* distinguish between the two and recognise the `nil` separately in a the "trigger" channel mode, but it's just not worth it to combine the two and having a separate "quit" channel to "trigger" channels makes for easier reading, as every channel receive has to distinguish between `struct {}{}` and `nil` every time, which is a waste of time.

For each API call we have a pair of channels created in this structure, one for the request, and one for the response. It is possible to embed the response in the requests as well, so the return channel is sent in the request, but this costs more in initiation and makes more sense for a lower frequency request service than the one we are making. The channel is already initialised before any requests are made, so the responses will return immediately and do not have to be initialised by the caller.

The last point, is that each of the 4 channels we have for our two API methods call and response are slices. This allows us to define how many threads we will warm up when starting up the service to handle the expected workload. The reason for making this pre-allocated is again, response latency. This may not be so important in some types of applications, but in all cases, the result of pre-allocation is better performance and a lack of unexpected unbounded resource allocation, commonly known as a "resource leak" which can extend to not just variables but also channels and goroutines, both of which have a small but nonzero cost in allocation time and memory utilisation until they are freed by the garbage collector.

[->contents](#kitchensink)

#### About Waitgroups

The last element of the struct is a `sync.Waitgroup`. This is an atomic counter which can be used to keep track of the number of threads that are running, and allows you to write code that holds things open until all of the wait group are `Done` when shutting down.

[->contents](#kitchensink)

#### Initialising the Worker Pool

Now that we have given a brief (as possible) explanation of the elements of the worker pool, here is the initialiser:

```go
// NewWorkerPool initialises the data structure required to run a worker pool.
// Call Start to to initiate the run, and call the returned stop function to end
// it.
func NewWorkerPool(workers uint32, stop chan struct{}) *transcriber {

	// Initialize a transcriber worker pool
	t := &transcriber{
		stop:         stop,
		encode:       make([]chan *proto.EncodeRequest, workers),
		decode:       make([]chan *proto.DecodeRequest, workers),
		encodeRes:    make([]chan proto.EncodeRes, workers),
		decodeRes:    make([]chan proto.DecodeRes, workers),
		encCallCount: atomic.NewUint32(0),
		decCallCount: atomic.NewUint32(0),
		workers:      workers,
		wait:         sync.WaitGroup{},
	}

	// Create a channel for each worker to send and receive on, buffer them by
	// the same as the number of workers, producing workers^2 job slots
	for i := uint32(0); i < workers; i++ {
		t.encode[i] = make(chan *proto.EncodeRequest)
		t.decode[i] = make(chan *proto.DecodeRequest)
		t.encodeRes[i] = make(chan proto.EncodeRes)
		t.decodeRes[i] = make(chan proto.DecodeRes)
	}

	return t
}
```

The `make` function is a built-in function that is required to create channels, maps and slices. We explained its usage (and `append`) in its first appearance, way back in the encoder, as it applies to slices. It also is used with maps and channels. 

We are not using maps in this tutorial, so, along with mutexes, these are two features of Go that we will not explain here.

Both of those features have very specific purposes, and Go's implementation, being bare bones as is the pattern with everything in Go, have caveats you will need to learn about later when you need to use them. We personally avoid mutexes always because they require keeping track of open and close operations and thus are very prone to error, where atomics simply access the value and write the value at a given moment and do not need this "bracketing". 

Maps are useful especially for lists of things where there must not be a duplicate key, but are annoying because to sort them you have to copy all their references into a slice and write a custom sort function that is determined by the nature of the data in the keys, and values, as the specific case dictates.

Back to channels, the `make` function *must* be called before the channel is sent or received on, because it will cause a panic if they are used without initialisation. 

In this code, we first make the slices of channels, and then run a loop for each worker and initialise the channel for each of the four channels we make for each worker thread.

Note that it is possible to use literals to initialise maps and slices, but it is not possible to use literals to initialise channels. However, initialising slices and maps with literals is not considered idiomatic, and partly because make allows the initial allocation of slice length, capacity or map capacity, and as mentioned previously, pre-allocation prevents unpredictable delays during allocation especially in tight initialisation or iteration loops.

[->contents](#kitchensink)

#### Atomic Counters

You wil also notice `encCallCount` and `decCallCount` are `atomic` types.

Atomic values only permit one accessor at a time, and have a built-in lock that protects their access, if two threads try to access an atomic at exactly the same time, the first one gets first go and then the second thread waits until the value has been read or written (`Load` and `Store`) before it gets its access. They default to access by copy, so if the data inside the atomic is a pointer, one should use `sync/mutex` locks instead. For this reason as you will see if you inspect the documentation for this package, all of the types that are available on uber's version of `atomic` are "value" types.

One can encapsulate *any* variable in an `atomic.Value` and the Go standard library `sync/atomic` . To see more about this, I will introduce you to `godoc` which lives at [https://pkg.go.dev](https://pkg.go.dev) - here is the `sync/atomic` package: https://pkg.go.dev/sync/atomic

We are using the [go.uber.org/atomic](https://pkg.go.dev/go.uber.org/atomic) variant out of habit because it has a few more important types as you can see if you click the link I just dropped there. Boolean, Duration, Error, Float64, Int32, Int64, Time, String and some unsafe things. 

Uber's version is better, because it creates proper initialisers for everything, whereas with the inbuilt `sync/atomic` library we have to boilerplate all that stuff to do it correctly. As I say, in this case, it's just habit as the non-pointer version of `sync.Int32` would serve just as well, but it is rare when one is using atomics that one is only doing such a small and singular thing with it as this, usually it will appear in multiple places protecting exclusive access with nicer syntax than mutexes, and facilitating very fast increment and decrement operations for counters, which is our use case here.

Note that strings are a special case of a data type that is actually a pointer in the underlying implementation, as they are immutable, and their memory area is marked not writeable to the memory management unit which will cause a segmentation fault and halt the program. Thus they are the sole exception for reference types in atomic variables, as they are not changeable, there cannot be a read/write race condition. A modified string is a new pointer and implicitly has no conflict with the old one. Though there can be a race condition on the rewriting of the pointer, which is why an atomic string is needed.

[->contents](#kitchensink)

#### Running the Worker Pool

The next thing to put into `pkg/grpc/workerpool.go` is the start function, but this depends on two further functions we are adding for reasons of neatness and good practice.

First is the handler loop:

```go
// handle the jobs, this is one thread of execution, and will run whatever job
// has appeared and this thread
func (t *transcriber) handle(worker uint32) {

	t.wait.Add(1)
out:
	for {
		select {
		case msg := <-t.encode[worker]:

			t.encCallCount.Inc()
			res, err := based32.Codec.Encode(msg.Data)
			t.encodeRes[worker] <- proto.EncodeRes{
				IdNonce: msg.IdNonce,
				String:  res,
				Error:   err,
			}

		case msg := <-t.decode[worker]:

			t.decCallCount.Inc()

			bytes, err := based32.Codec.Decode(msg.EncodedString)
			t.decodeRes[worker] <- proto.DecodeRes{
				IdNonce: msg.IdNonce,
				Bytes:   bytes,
				Error:   err,
			}

		case <-t.stop:

			break out
		}
	}

	t.wait.Done()

}
```
This is something that you won't see many beginner tutorials showing, because beginner Go tutorials usually are just showing you how to use the `net/http` library, which has a loop like this inside it, and you only have to write a function or closure and load it into the initialiser function, before calling the `Serve` function. 

We are here showing you what goes on underneath this layer, because we are creating our own concurrent handler for our server.

You can see that each loop, which constitutes each "worker" in the worker pool, is basically the same, it receives a job from the `encode` and `decode` channels, using the `worker` variable to define which of the slice members it will access the channel for, which defines the scope of each worker's remit.

The actual function is then called, that we defined back in the `based32` package, and then it sends the result back on the `encodeRes` and `decodeRes` results channels.

This is the in-process equivalent of the RPC that is the next part of the server we will be creating after we have finished creating the worker pool, that uses channels instead of network sockets to pass messages.

It's important to understand that part of the reason why Go has coroutines (goroutines) and channels is precisely to act as a model for the external, network facing parts, to provide the programmer with one model and two different implementations, one for inside, one for outside, to enable the accurate modelling of concurrent processes that communicate with each other. Other programming languages don't strive for this consistency and for which reason the "async" model of programming is more often used, which hides the handler and callback mechanisms and the processing is also, usually single threaded, which is less performant.

[->contents](#kitchensink)

#### Logging the call counts

There isn't any point in putting the counter in there if it wasn't going to come back out somewhere, and the logical place for this is in a log print that summarises the activity of the server's run:

```go
// logCallCounts prints the values stored in the encode and decode counter
// atomic variables.
func (t *transcriber) logCallCounts() {

	log.Printf(
		"processed %v encodes and %v encodes",
		t.encCallCount.Load(), t.decCallCount.Load(),
	)
}
```
[->contents](#kitchensink)

#### Starting the worker pool

The start function spawns the workers in a loop, calling the `handle` function for each worker, spawning an individual thread for each, and returns a function that runs when they stop. In this case it just waits until all the processes have completed their cleanup before calling the `logCallCounts` function.

```go
// Start up the worker pool.
func (t *transcriber) Start() (cleanup func()) {

	// Spawn the number of workers configured.
	for i := uint32(0); i < t.workers; i++ {

		go t.handle(i)
	}

	return func() {

		log.Println("cleanup called")

		// Wait until all have stopped.
		t.wait.Wait()

		// Log the number of jobs that were done during the run.
		t.logCallCounts()

	}
}
```

[->contents](#kitchensink)

#### Creating the gRPC Service

Now that we have the worker pool created, we need to create the gRPC service that will use it.

The reason why we created the gRPC services as "stream" type will now be revealed. If the gRPC service definition did not use the `stream` keyword, we would be stuck with the back end handling of our work done by the gRPC library's boilerplate method.

This might be ok for small scale systems, but for larger ones, with potentially high work loads, the management of the concurrency of the work should be done by the server developer as the way to do this can vary quite widely between one type of work and another.

Note that as we will see later, current implementation of stream in the gRPC implementation for Go has a limitation that seems to cause a block with large or numerous requests in a short period of time. Just as the Go implementation has lagged in coverage of features, it also is lagging in stability compared to other language versions, which is unfortunate. Having to add message segmentation or ratelimiting to handle this issue is a big strike against using gRPC if one does not need to support other languages.

With large APIs with lower overall workloads, it is fine to let gRPC handle the scheduling, but since that is easier and doesn't grant an opportunity to teach concurrent programming, the most distinctive feature of Go, we are showing you how to do it yourself. 

Personally, I prefer always handling the concurrency, because no matter how much load my services get, it can be later on when I have no part in the exercise, that someone tries to use my code for a high load scenario and because I wrote it from scratch to be optimal, it will not fail in the hands of anyone. Besides this, having written highly time sensitive multicast services previously, it is now just habit for me to write my own concurrency handling, and once you get used to doing it, you won't like to leave it to others either once you have had your own goroutine leak or two, you won't want to leave it up to chance either.

As with the worker pool, we are going to create a type that is not exported and export the initialiser due to the necessity of the use of `make()` when using objects within the structure.

In order to support that, we must create a constructor function, here we have `New`, which accepts a network address and the number of worker threads we want to run to handle requests.

```go
package server

import (
	"github.com/quanterall/kitchensink/pkg/proto"
	"go.uber.org/atomic"
	"google.golang.org/grpc"
	"io"
	"net"
)

// b32 is not exported because if the consuming code uses this struct directly
// without initializing it correctly, several things will not work correctly,
// such as the stop function, which depends on there being an initialized
// channel, and will panic the Start function immediately.
type b32 struct {
	proto.UnimplementedTranscriberServer
	stop        chan struct{}
	svr         *grpc.Server
	transcriber *transcriber
	addr        *net.TCPAddr
	roundRobin  atomic.Uint32
	workers     uint32
}

// New creates a new service handler
func New(addr *net.TCPAddr, workers uint32) (b *b32) {

	log.Println("creating transcriber service")

	// It would be possible to interlink all of the kill switches in an
	// application via passing this variable in to the New function, for which
	// reason in an application, its killswitch has to trigger closing of this
	// channel via calling the stop function returned by Start, further down.
	stop := make(chan struct{})
	b = &b32{
		stop:        stop,
		svr:         grpc.NewServer(),
		transcriber: NewWorkerPool(workers, stop),
		addr:        addr,
		workers:     workers - 1,
	}
	b.roundRobin.Store(0)

	return
}
```

Next, in keeping with our philosophy of creating our source files in order of declaration, even though it is not mandatory, but for reading, it is, we provide implementations for interfaces created by the `proto` generated gRPC package, that is created from our `.proto` file.

```protobuf
service Transcriber {
  rpc Encode(stream EncodeRequest) returns (stream EncodeResponse);
  rpc Decode(stream DecodeRequest) returns (stream DecodeResponse);
}
```

The method signature is not exactly the same, because "returns" in the protobuf specification does not map to a function return, but rather, that our function sends back to a stream, which is something a bit like a channel, but not exactly, because we have the freedom to change the ordering... if we want to, and the client recognises responses as associated with requests automatically on the other end so we can change the order here. 

But we don't want that, just that we want to be able to handle that concurrency ourselves in Go, by sending the requests to our worker pool and when they return, returning them via the gRPC prescribed response `stream` mechanism.

Our job distribution strategy, as hinted at in the structure definition, is round robin, that is, each new job is passed to each worker in sequence, 1, 2, 3, 4, etc. 

The `stream` processing strategy essentially means we have a loop that runs that receives each job as it comes in, and is supposed to process it, and then return it as necessary.

```go
// Encode is our implementation of the encode API call for the incoming stream
// of requests.
//
// Note that both this and the next stream handler are virtually identical
// except for the destination that received messages will be sent to. There is
// ways to make this more DRY, but they are not worth doing for only two API
// calls. If there were 5 or more, the right solution would be a code generator.
// It is a golden rule of Go, if it's not difficult to maintain, copy and paste,
// if it is, write a generator, or rage quit and use a generics language and
// lose your time waiting for compilation instead.
func (b *b32) Encode(stream proto.Transcriber_EncodeServer) error {
out:
    for {

        // check whether it's shutdown time first
        select {
        case <-b.stop:
            break out
        default:
        }

        // Wait for and load in a newly received message
        in, err := stream.Recv()
        switch {
        case err == io.EOF:

            // The client has broken the connection, so we can quit
            break out
        case err != nil:

            // Any error is terminal here, so return it to the caller after
            // logging it
            log.Println(err)
            return err
        }

        worker := b.roundRobin.Load()
        log.Printf("worker %d starting", worker)
        b.transcriber.encode[worker] <- in
        log.Printf("worker %d encoding", worker)
        res := <-b.transcriber.encodeRes[worker]
        log.Printf("worker %d sending", worker)
        err = stream.Send(proto.CreateEncodeResponse(res))
        if err != nil {
            log.Printf("Error sending response on stream: %s", err)
        }
        log.Printf("worker %d sent", worker)
        if worker >= b.workers {
            worker = 0
            b.roundRobin.Store(0)
        } else {
            worker++
            b.roundRobin.Inc()
        }
    }
    log.Println("encode service stopping normally")
    return nil
}

// Decode is our implementation of the encode API call for the incoming stream
// of requests.
func (b *b32) Decode(stream proto.Transcriber_DecodeServer) error {

out:
    for {

        // check whether it's shutdown time first
        select {
        case <-b.stop:
            break out
        default:
            // When a select statement has a default case it always terminates.
        }

        // Wait for and load in a newly received message
        in, err := stream.Recv()
        switch {
        case err == io.EOF:

            // The client has broken the connection, so we can quit
            break out

        case err != nil:

            // Any error is terminal here, so return it to the caller after
            // logging it, and ending this function terminates the decoder
            // service.
            log.Println(err)
            return err
        }
        worker := b.roundRobin.Load()
        b.transcriber.decode[worker] <- in
        res := <-b.transcriber.decodeRes[worker]
        err = stream.Send(proto.CreateDecodeResponse(res))
        if err != nil {
            log.Printf("Error sending response on stream: %s", err)
        }
        if worker >= b.workers {
            worker = 0
            b.roundRobin.Store(0)
        } else {
            worker++
            b.roundRobin.Inc()
        }
    }

    log.Println("decode service stopping normally")
    return nil
}
```

These handlers can be run concurrently, and each job's identification number will be correctly associated with the requesting goroutine, in the next section.

But first, the last piece of the server, which starts up the threads and manages their shutdown:

```go
// Start up the transcriber server
func (b *b32) Start() (stop func()) {

    proto.RegisterTranscriberServer(b.svr, b)

    log.Println("starting transcriber service")

    cleanup := b.transcriber.Start()

    // Set up a tcp listener for the gRPC service.
    lis, err := net.ListenTCP("tcp", b.addr)
    if err != nil {
        log.Fatalf("failed to listen on %v: %v", b.addr, err)
    }

    // This is spawned in a goroutine so we can trigger the shutdown correctly.
    go func() {
        log.Printf("server listening at %v", lis.Addr())

        if err := b.svr.Serve(lis); err != nil {

            // This is where errors returned from Decode and Encode streams end
            // up.
            log.Printf("failed to serve: '%v'", err)

            // By the time this happens the second goroutine is running and it
            // is always better unless you are sure nothing else is running and
            // part way starting up, to shut it down properly. Closing this
            // channel terminates the second goroutine which calls the server to
            // stop, and then the Start function terminates. In this way we can
            // be sure that nothing will keep running and the user does not have
            // to use `kill -9` or ctrl-\ on the terminal to end the process.
            //
            // If force kill is required, there is a bug in the concurrency and
            // should be fixed to ensure that all resources are properly
            // released, and especially in the case of databases or file writing
            // that the cache is flushed and the on disk store is left in a sane
            // state.
            close(b.stop)
        }
        log.Printf(
            "server at %v now shut down",
            lis.Addr(),
        )

    }()

    go func() {
    out:
        for {
            select {
            case <-b.stop:

                log.Println("stopping service")

                // This is the proper way to stop the gRPC server, which will
                // end the next goroutine spawned just above correctly.
                b.svr.GracefulStop()
                cleanup()
                break out
            }
        }
    }()

    // The stop signal is triggered when this function is called, which triggers
    // the graceful stop of the server, and terminates the two goroutines above
    // cleanly.
    return func() {
        log.Printf("stop called on service")
        close(b.stop)
    }
}
```

You can see here, the second goroutine is just waiting on the shutdown, to run the cleanup code. We can't do the two things in one goroutine because the gRPC `Serve` function handles its own concurrency, which as mentioned previously, has proven to have an undocumented limitation on message size and stream writes in a short period of time. This issue has been noted on the relevant issues board on Github.

One goroutine handles running the stream service, the other waits for the shutdown signal and propagates this shutdown request to the stream service so everything stops correctly.

It is a common problem for beginners working with concurrency in Go to have applications that have to be force-killed (`kill -9 <pid>` or `ctrl-\` on the terminal). It is necessary in Go to properly handle stopping all goroutines as they continue to execute (or wait) and the process does not terminate.

[->contents](#kitchensink)

----

### [Step 7](steps/step7) Creating a Client

Again, we want to put our nice logger in [pkg/grpc/client/log.go](pkg/grpc/client/log.go)

```go
package client

import (
	logg "log"
	"os"
)

var log = logg.New(os.Stderr, "based32", logg.Llongfile|logg.Lmicroseconds)
```

[->contents](#kitchensink)

#### Data Structures for the Client

In [pkg/grpc/client/types.go](pkg/grpc/client/types.go) we will define the specific types required for the client.

We will be using a concurrent architecture model that you find a lot in Go applications - clients sending requests with channels inside them for the return message.

Also, like the last two parts of the server, due to the need to initialise members of the data structure that ties together the package, the primary data structure is not exported and we use a constructor function to handle this for us. For the client there is also the two request types used. We will add them first:

```go
package client

import (
    "context"
    "github.com/quanterall/kitchensink/pkg/proto"
    "google.golang.org/grpc"
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
```

We also can make their constructor non exported, as the start function will return small closures that will encapsulate this and allow sync or async invocation with nice syntax.

Next, because of the precedence of declarations, we will add the core client data structure to this file as well:

```go
type b32c struct {
    addr       string
    encChan    chan encReq
    encRes     chan *proto.EncodeResponse
    decChan    chan decReq
    decRes     chan *proto.DecodeResponse
    timeout    time.Duration
    waitingEnc map[time.Time]encReq
    waitingDec map[time.Time]decReq
}
```

The stop channel for the client is provided by a `context.Context` and here we see an example of a one way channel. This one can only receive. This is often a good idea when there is no reason anyway for signals to be triggered in the opposite direction. Really, stop channels are not one of these types of cases, but the context package has them this way, and we need to listen on them as cancelling the context closes this channel and we need to listen on this channel to close the stream handlers.

There is no strict rule or idiom about including types and constructors or methods, but because there is only a constructor and a start function, and the start function returns the request functions for each API and a stop function to stop the client, we are putting them in a separate file.

[->contents](#kitchensink)

#### Client Constructor

In [pkg/grpc/client/client.go](pkg/grpc/client/client.go) we will put a constructor:

```go

func New(serverAddr string, timeout time.Duration) (
    client *b32c, err error,
) {

    client = &b32c{
        addr:       serverAddr,
        encChan:    make(chan encReq, 1),
        encRes:     make(chan *proto.EncodeResponse),
        decChan:    make(chan decReq, 1),
        decRes:     make(chan *proto.DecodeResponse),
        timeout:    timeout,
        waitingEnc: make(map[time.Time]encReq),
        waitingDec: make(map[time.Time]decReq),
    }

    return
}
```

As you can see, the majority of this function is initialising channels and maps.

[->contents](#kitchensink)

#### The Encode and Decode Handlers

We are not going to put the function to start up the client yet, because we first need the handlers for the encode and decode API calls.

These two things are literally identical with `En` swapped for `De` in the names things. So we will just pile the two together here.

For cases where there is more calls to put in here, where the differences are similarly simple, just passing values of a type to handlers over channels all with consistent naming schemes, once you see more than two identical functions like these, you would write a generator, which is provided several names to put in each slot and spit out several files each with the same code with different names. Again, this is essentially what C++ templates and Java and Rust generics are, there really is no practical reason why they need to exist in a language.

But we won't make generators for two of them, because writing one, copying and replacing the different characters is not terribly onerous. Besides this, with generators, generally one has to have a testable prototype to make the code work, before making the copy paste substitutions. These substitutions, in the case of this particular code, could literally be done with search and replace.

So, actually, we are going to do this, as an alternative example to how to write a generator for generic implementations to the standard common use of templates. In this case, we will demonstrate the use of string search and replace functions.

So, first, the encode client handler:

```go
package client

import (
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

				// log.Println("sending message on stream")
				err := stream.Send(msg.Req)
				if err != nil {
					log.Print(err)
				}
				b.waitingEnc[time.Now()] = msg
			case recvd := <-b.encRes:

				for i := range b.waitingEnc {

					// Return received responses
					if recvd.IdNonce == b.waitingEnc[i].Req.IdNonce {

						// return response to client
						b.waitingEnc[i].Res <- recvd

						// delete entry in pending job map
						delete(b.waitingEnc, i)

						// if message is processed next section does not need to
						// be run as we have just deleted it
						continue
					}

					// Check for expired responses
					if i.Add(b.timeout).Before(time.Now()) {

						log.Println(
							"\nexpiring", i,
							"\ntimeout", b.timeout,
							"\nsince", i.Add(b.timeout),
							"\nis late", i.Add(b.timeout).Before(time.Now()),
						)
						delete(b.waitingEnc, i)
					}
				}
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
				log.Println("stream closed")
				break in
			case err != nil:

				log.Println(err)
				break
			}

			// forward received message to processing loop
			b.encRes <- recvd
		}
	}(stream)

	return
}
```

But before the generator, let's discuss what's happening in here.

There is two goroutines started by the handler, each containing forever loops, each using the `stream` for reading and writing, respectively.

The reason why two are needed is because the `stream.Recv` function blocks, so it cannot be part of a select block where events are driven by channel receives.

The first goroutine has a stop handler, which breaks the forever loop using the `out` label, it has a channel receive for encode requests which sends the request on the stream, and a receiver handler, which we feed from the second goroutine that runs a polling loop instead of event loop because of the semantics of the receive function.

The receive loop also checks for stop channel between receives, and will also break out if the stream is closed, which causes it to yield an `io.EOF` (end of file) signal, which occurs in case of client disconnect or triggering of the shutdown of the stream. The received messages are then forwarded to the encode results channel, which checks its queue, returns the message on the matching return channel, based on the message IdNonce and then deletes it from the map storing pending responses.

Note that as you can see, the 'waitingEnc' variable is a map. If this map is accessed from two different goroutines, inevitably there will be two threads attempting to access it at the same time which will cause a panic which will be described as a "concurrent map read/write". For this reason, these accesses are kept within the first goroutine. The receiver adds them to the queue, and the sender deletes them once they are delivered.

Lastly, using the timeout value selected, if a request gets hung up for more time than the timeout, the request map entry is deleted, assumed failed. If one does not set a timeout for request fulfillment, it can cause resource leaks that accumulate over time and will require the restarting of the server otherwise. To not time them out means that failures on one end propagate to the other end for no reason. The client can attempt to re-establish connection, outside of this code, with a retry with backoff to handle the case of a transient failure. On both sides, lack of progress should trigger a restart of services, at some level.

When one has written a lot of peer to peer and blockchain type code, it is standard practice to expect connection failures and potential crashes on peers, and simply drop things and start again. Failures can be from many causes, network connections, denial of service attacks, bugs in network handlers, and bugs in the applications attached to the network handlers, and, of course, attacks on the applications themselves. Such misbehaviour is generally termed "crash failures" but can be caused by malice, and called "byzantine". Handling crash failure cases is essential in distributed systems, byzantine failures take a little more design and some game theory to work around.

[->contents](#kitchensink)

#### Copy Paste Generic Generator

Inside the `pkg/grpc/client` folder create a new folder `gen` and inside that, create a new file "derive.go".

This is the generator which will allow us to write one handler and duplicate it with changed names for all the relevant differing parts. This one is designed to create a single counterpart, however, one could conceivably have more than one substitution and write multiple files with different names based on this same design.

```go
package main

import (
    "io/ioutil"
    "log"
    "strings"
)

func main() {

    // The following are the strings that differ between encoder.go and decoder.go
    subs := [][]string{
        {"Encode", "Decode"},
        {"encChan", "decChan"},
        {"encRes", "decRes"},
        {"waitingEnc", "waitingDec"},
        {"//go:generate go run ./gen/.", "// generated code DO NOT EDIT"},
    }
    bytes, err := ioutil.ReadFile("encoder.go")
    if err != nil {
        log.Println(err)
        return
    }
    file := string(bytes)
    for i := range subs {
        file = strings.ReplaceAll(file, subs[i][0], subs[i][1])
    }
    err = ioutil.WriteFile("decoder.go", []byte(file), 0755)
    if err != nil {
        return
    }
}
```

It just occurs to me as I write this, that this will be the first example of reading and writing files in this tutorial. As such, we are using the simplest possible tools to do this, from `ioutil` library. This little piece of code essentially does for you what you can do manually in about 2 minutes on such a task, but we now don't have to do this work.

It also just occurs to me that it would have been a terrible omission to not introduce writing generators in this tutorial, since in Go it is an essential skill. This is the simplest possible code generator that you can write. This doesn't require any special handling or duplication in the form of writing a complex template and feeding it the respective sets of values to place in each field. 

The `encoder.go` file is live code that also compiles, and once you write this, you have the two API handlers for both sides, you can change one, run `go generate ./...` and the counterpart is updated to match.

This isn't the first time we have shown the use of `go:generate` however. It is not the only way to handle these things, as many projects will use `Makefile`s in this capacity, but it's not necessary and is a step that isn't done often, so instead, you can put these lines in appropriate places and update the relevant source code automatically in one step over the whole repository.

As such, you don't get to see the magic this does until we put this generator into the source `encode.go` file. Here it is, within its context:

```go
package client

//go:generate go run ./gen/.

import (
    "github.com/quanterall/kitchensink/pkg/proto"
    "io"
    "time"
)

func (b *b32c) Encode(stream proto.Transcriber_EncodeClient) (err error) {
```

Once you have added this line to this file, you can then run `go generate` either in the folder where encode.go is found, with the subdirectory `gen` containing the generator code above, or invoke it and retrigger all generate lines in the project at the root with:

```bash
go generate ./...
```

In the second case, you will see it will re-run the protobuf compiler as well.

> For the benefit of developers of this tutorial itself, we have also added a generator that updates the table of contents all markdown documents in the repository inside `doc.go` at the root of the repository. This file was placed there initially just so that there was a "package" at the root of the repository, but it becomes a handy place to put `go:generate` lines that do other things, like updating the table of contents. 
>
> Note that unless you first run `scripts/installtoc.sh` which will require you to install an Ubuntu/Debian package called "fswatch" and then a simple go app which generates table of contents for markdown files called `tocenize`. See [scripts/README.md](scripts/README.md) for these items.

[->contents](#kitchensink)

----

### [Step 8](steps/step8) Testing the gRPC server


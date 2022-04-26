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
		- [Making the output code more useful with some extensions](#making-the-output-code-more-useful-with-some-extensions)
		- [Documentation comments in Go](#documentation-comments-in-go)
		- [go:generate line](#gogenerate-line)
		- [Import Alias](#import-alias)
		- [Adding a Stringer for the generated Error type](#adding-a-stringer-for-the-generated-error-type)
		- [Convenience types for results](#convenience-types-for-results)
		- [Create Response Helper Functions](#create-response-helper-functions)

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
package proto;
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
code package to account for this.

In Go, returns are tuples, usually `result, error`, but in other languages like
Rust and C++ they are encoded as a "variant" which is a type of `union` type.
The nearest equivalent in Go is an `interface` but interfaces can be anything.
The union type was left out of Go because it breaks C's otherwise strict typing
system (and yes, this has been one of the many ways in which C code has been
exploited to break security, which is why Go lacks it).

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

#### The concrete type

In the root of the repository, create a new file called `types.go`. This is
where we will define the main types that will be used by packages and
applications that use our code.

#### Package header

```go
package codec

import (
	"github.com/quanterall/kitchensink/pkg/codecer"
)
```

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

#### Interface Implementation Assertion

The following var line makes it so the compiler will throw an error if the
interface is not implemented.

```go
// The following implementations are here to ensure this type implements the
// interface. In this tutorial/example we are creating a kind of generic
// implementation through the use of closures loaded into a struct.

// This ensures the interface is satisfied for codecer.Codecer and is removed in
// the generated binary because the underscore indicates the value is discarded.
var _ codecer.Codecer = &Codec{}
```

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
this method's execution.

When the type is a potentially shared or structured (struct or `[]`) type, the
copy will waste time copying the value, or referring to a common version in the
pointer embedded within the slice type (or map), and memory to store the copy,
and potentially lead to errors from race conditions or unexpected state
divergence if the functions mutate values inside the structure. Usually non
pointer methods are only used on simple value types like specially modified
versions of value types (anything up to 64 bits in size, but also arrays, which
are `[number]type` as opposed to `[]type`), or when this copying behaviour is
intended to deliberately avoid race conditions, and the shallow copy will not
introduce unwanted behaviours.

In this case, the pointer is fine, because it is not intended that the 
`Codec` type ever be changed after initialisation. However, because its 
methods and fields are exposed, code that reaches through and modifies this 
could break this assumption. However, the rationale for such a mutation is 
hard to justify or even conceive so any programmer who mutates this 
particular structure is behaving in a very idiosyncratic way which might be 
called stupid, or may just mean they are ignorant of the implications of 
concurrency, or, alternatively, they know, and their code is not concurrent.

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

#### Making the output code more useful with some extensions

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

#### Import Alias

Note the import alias here is present to explicitly refer to its' name as set in
the package line at the top of the file. You need to change this URL to match
the URL of your package name.

We are using the name codec, even though in our repository this is
`kitchensink` but for your work for the tutorial you will call your package,
presumably, `codec`. This is not mandatory but it is idiomatic to match the name
of a package and the folder it lives in.

Go will expect the name defined in this line to refer to this package, so it is
confusing to see the export ending with a different word. In such cases it is
common to explicitly use a renaming prefix in the import (an alias that is found
just before the import path in an import block or statement). This is why we
have it here, even though in the `types.go` file it has `package codec`

```go
import (
    codec "github.com/quanterall/kitchensink"
)
```

#### Adding a Stringer for the generated Error type

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

```

#### Convenience types for results

The following types will be used elsewhere, as well as for the following create
response functions. These are primarily to accommodate for the fact that
protobuf follows c++ conventions with eth use of 'oneof' variant types, which
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
an error is not fatal, and the variant return convention ignores this, and makes
more complexity for programmers, and compiler writers.

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

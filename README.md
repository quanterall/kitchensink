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
	- [Step 4 The Encoder](#step-4-the-encoder)
		- [Always write code to be extensible](#always-write-code-to-be-extensible)
		- [Helper functions](#helper-functions)
		- [Log at the site](#log-at-the-site)
		- [Create an Initialiser](#create-an-initialiser)
		- [Writing the check function](#writing-the-check-function)
		- [Creating the Encoder](#creating-the-encoder)
		- [Calculating the check length](#calculating-the-check-length)
		- [Writing the Encoder Implementation](#writing-the-encoder-implementation)
		- [Creating the Check function](#creating-the-check-function)
		- [Creating the Decoder function](#creating-the-decoder-function)

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
  string EncodedString = 1;es any check function defined for the type.
	//
	// If the check fails or the input is too short to have a check, false and
	// nil is returned. This is the contract for this method that
	// implementations should uphold.
	Decode(input string) (output []byte, err error)
}
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

----

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
this method's execution, potentially consuming a lot of memory and time moving that memory around.

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

----

### [Step 4](steps/step4) The Encoder

Next step is the actual library that the protobufs and interface and types were all created for.

#### Always write code to be extensible

While when making new libraries you will change the types and protocols a lot as you work through the implementation, it is still the best pattern to start with defining at least protocols and making a minimal placeholder for the implementation.

It is possible to create a library that does not have any significant state or configuration that just consists of methods that do things, which can be created without a `struct` tying them together, this is rare, and usually only happens when there is only one or two functions required. 

For this we have 4 functions and while we could hard code everything with constants and non-exported variables where constants can't be used, this is not extensible. Even in only one year of full time work programming, I estimate that I spent about 20% of my first year working as a Go developer, fixing up quick and dirty written code that was not designed to be extended. 

The time cost of preparing a codebase to be extensible and modular is tiny in comparison, maybe an extra half an hour as you start on a library. Experience says that the shortcut is not worth it. You never know when you are the one who has to extend your own code later on, and two days later it's finally in a state you can add functionality.

#### Helper functions

The only exception to this is when there is literally only one or at most two functions to deal with a specific type of data. These are often referred to as "helpers" or "convenience functions" and do not need to be extensible as they are very small and self contained.

These functions can sometimes be tricky to know where to put them, and often end up in collections under package names with terrible names like "util" or "tools" or "helpers". This can be problematic because very often they are accessory to another type, and doing this creates confusing crosslinks that can lead you into a circular dependency.

As such, my advice is to keep helpers where they are used, and don't export them, unless they are necessary, like the response helper functions we made previously for the `proto` package.

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

var log = logg.New(os.Stderr, "based32 ", logg.Llongfile)
```

What we are doing here is using the standard logging library to set up a customised configuration. The standard logger only drops the timestamp with the log entries, which is rarely a useful feature, and when it is, the time precision is too low on the default configuration, as the most frequent time one needs accurate timestamps is when the time of events is in the milli- or microseconds when debugging concurrent high performance low latency code.

This log variable essentially replaces an import in the rest of the package for the `log` standard library, and configures it to print full file paths and label them also with the name of the package.

It is ok to leave one level of indirection in the site of logging errors, that is, the library will return an error but not log, but it should at least log where the error returns, so that when the problem comes up, you only have to trace back to the call site and not several layers above this.

When you further have layers of indirection like interfaces and copies of pointers to objects that are causing errors, knowing which place to look for the bug will take up as much time as actually fixing it.

It may be that you are never writing algorithms that need any real debugging, many "programmers" rarely have to do much debugging. But we don't want to churn out script writers only, we want to make sure that everyone has at least been introduced to the idea of debugging. 

For that reason also, now that we are implementing an algorithm here, we are going to deliberately cause bugs and force the student to encounter the process of debugging, show the way to fix them, and not just make this an exercise in copy and paste, for which there will be no benefit as bugs are the way you learn to write good code, without that difficulty, it is not programing, and you will forget the next day how you did it, which makes this whole exercise a waste of time that you could have saved yourself keystrokes and just read it instead.

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

You will notice that we took care to make sure that everything you will paste into your editor will pass syntax checks immediately. All functions that have return values must contain a `return` statement. The return value here is imported from `types.go` at the root of the repository, which the compiler identifies as `github.com/quanterall/kitchensink` because of running `go mod init` in [Initialize your repository](#initialize-your-repository) .

Also note it explicitly 

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

It is possible to instead use the shorter, and simpler `crc32` checksum function, but we don't like it because it requires converting bytes to 32 bit integers and back again, which is done in practise like this:

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

If you follow the logic of that conversion, you can see that it is 4 copy operations, 3 bit shifts and 3 addition operations. The hash function does not do this conversion, it operates directly on bytes (in fact, I think it uses 8 byte/64 bit words, and coerces the byte slices to 64 bit long words using an unsafe type conversion) using a sponge function, and Blake3 is the fastest hash function with cryptographic security, which means a low rate of collisions, which in terms of checksums equates to two strings creating the same checksum, and breaking some degree of security of the function. So, we use blake3 hashes and cut them to our custom length.

The length is variable as we are designing this algorithm to combine padding together with the check. So, essentially the way it works is we take the modulus (remainder of the division) of the length of the data, and pad it out to the next biggest multiple of 5 bytes, which is 8 base32 symbols. The formula for this comes next.

#### Creating the Encoder

In all cases, when creating a codec, the first step is making the encoder. It is impossible to decode something that doesn't yet exist, the encoder is *a priori*, that is, it comes first, both logically and temporally.

But before we can make the encoder, we also need a function to compute the correct check length for the payload data. 

Further, the necessity of a variable length requires also that the length of the check be known before decoding, so this becomes a prefix of the encoding.

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

The function takes the input of the length of our message in bytes, and returns the correct length. The result is that for an equal multiple of 5 bytes, we add 5 bytes of check, and 4, 3, 2 or 1 bytes for each of the variations that are more than this multiple, plus accounting for the extra byte to store the check length.

> The check length byte in fact only uses the first 3 bits, as it can be no more than 5, which requires 3 bits of encoding. Keep this in mind for later as we use this to abbreviate the codes as implicitly their largest 5 bits must be zero, which is precisely one base32 character in length, thus it always prefixes with a `q` as described by the character set we are using for this, based on the Bech32 standard used by Bitcoin and Cosmos, which reduces the length of the encoded value by one character, and must be added back to correctly decode.

#### Writing the Encoder Implementation

The standard library contains a set of functions to encode and decode base 32 numbers using custom character sets. The 32 characters defined in the initialiser defined earlier, are chosen for their distinctiveness

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

> **The syntax for slices and copying requires some explanation.**
>
> In Go, to copy slices (variable types with the `[]` prefix) you use the `copy` operator, the assignment operator `=` copies the *value* of the slice, which is actually a struct, internally, containing the pointer to the memory, the current length used and the capacity (which can be larger than current used), which would not achieve the goal of creating a new, expanded version including the check length prefix and check/padding.
>
> The other point to note is the use of the slicing operator. Individual elements are addressed via `variableName[n]` and to designate a subsection of the slice you use `variableName[start:end]` where the values represent the first element and the element *after* the last element.
>
> It is important to explain this as the notation can be otherwise confusing. The end index in particular needs to be understood as *exclusive* not *inclusive*. The length function `len(sliceName)` returns a value that is the same as the index that you would use to designate *up to the end* of the slice, as it is the cardinal value (count) where the ordinal (index) starts from zero and is thus one less.
>
> Lastly, the slicing operator can also be used on strings, but beware that the indexes are bytes, and do not respect the character boundaries of UTF-8 encoding, which is only one byte per character for the first 255 characters of ASCII and does not include any (many, it does include several umlauts and accent characters from european languages) non-latin symbols. 
>
> However, in the case of the Base32 encoding, we are using standard ASCII symbols so we know that we can cut off the first one to remove the redundant zero that appears because of the maximum 3 bits used for the check length prefix value, leave 5 bits in front (due to the backwards encoding convention for numbers within machine words). 
>
> *whew* A lot to explain about the algorithm above, but vital to understand for anyone who wants to work with slices of bytes in Go, which basically means anything involving binary encoding. This will be as deeply technical as this tutorial gets, it's not essential to understand it to do the tutorial, but this explanation is added for the benefit of those who do or will need to work with binary encoded data.

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

This is a very simple formula, but it needs to be used again in the decoder function where it allows the raw bytes to be correctly cut to return the checked value.

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

Take note about the use of the string cast above. In Go, slices do not have an equality operator `==` but strings do. Casting bytes to string creates an immutable copy so it adds a copy operation. If the amount of data is very large, you write a custom comparison function to avoid this duplication, but for short amounts of data, the extra copy stays on the stack and does not take a lot of time, in return for the simplified comparison as shown above.

#### Creating the Decoder function

The decoder cuts off the HRP, prepends the always zero first base32 character, decodes using the Base32 encoder (it is created prior to the encode function previously, and is actually a codec, though I used the name `enc`, it also has a decode function)

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
		// zeroed 5 bytes with a 'q' character.
		// Be aware the input string will be copied to create the []byte
		// version. Also, because the input bytes are always zero for the first
		// 5 most significant bits, we must re-add the zero at the front (q)
		// before feeding it to the decoder.
        input = "q" + input[len(cdc.HRP):]

		// The length of the base32 string refers to 5 bytes per slice index
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


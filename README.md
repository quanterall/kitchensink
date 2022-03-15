# kitchensink

## Teaching Golang via building a Human Readable Binary Transcription Encoding Framework

In this tutorial we will walk you through the creation from scratch of a human
readable binary transcription encoder/decoder.

This tutorial demonstrates the use of almost every possible and important
feature of Go. A "toy" implementation of a gRPC/protobuf microservice is added
in order to illustrate almost everything else.

In order to demonstrate synchronisation primitives, waitgroups, atomics and
mutexes, the service will keep track of the number of invocations, print this
count in log updates, and track the count using a concurrent safe atomic
variable and show the variant using a mutex instead, and run an arbitrary number
of concurrent worker threads that will start up and stop using waitgroups.

The final result is this library itself, and each step will be elaborated in
clear detail. Many tutorials leave out important things, and to ensure this does
not happen, each stage's parts will be also found in the [steps](./steps)
folder at the root of the repository.

## Prerequisites

This tutorial was developed on a system running Pop OS based on Ubuntu 21. 
As such there may be some small differences compared to Ubuntu 20, but 
protobuf will still be version 3.

In general, you will be deploying your binaries to systems also running 
ubuntu 20 or 21 or similar, on x86-64 platform, so the same instructions can 
be used in setting up a fresh server when deploying. We will not cover 
Docker or any other container system here.

Necessary things that you probably already have:

    sudo apt install -y build-essential git wget curl autoconf automake libtool

### Install Go

Go 1.17+ is recommended - unlike most other languages, the forward 
compatibility guarantee is ironclad, so go to 
[https://go.dev/dl/](https://go.dev/dl/) and pick the latest one (1.17.8 at 
time of writing), "copy link location" on the relevant version (linux x86-64)

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

This also creates a proper place where `go install` will put produced
binaries, which is recommended for avoiding cluttering up repositories you
are working on with binaries and potentially accidentally adding them to the
repository, which can be very problematic if you are working on a BIG
application (Go apps are generally under 60mb in size but this is still a
lot in a source code repository).

### Install Protobuf Compiler

In order to build the project you will need the protobuf compiler installed. 
This generates Go code from a protobuf specification for a service, creating 
all the necessary handlers to call and to handle calls to implement the API 
of our transcription codec.

    sudo apt install -y protobuf-compiler
    protoc --version  # Ensure compiler version is 3+

### Notes

~~Note that we also will be demonstrating the use of `make` as a build tool. 
This is not strictly necessary when developing Go applications, but it is 
very commonly used for this purpose and can simplify a lot of things. This 
entire section could be automated into a makefile script, for example.~~

## Step By Step:

Here will be the step by step process of building the library with a logical
sequence that builds from the basis to the specific parts for each in the order
that is needed both for understanding and for the constraints of syntax, grammar
and build system design...

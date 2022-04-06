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

### 1. Setting up the repository

The very first thing to do is to create a Git repository. It is not 
necessary to upload it to your github account or other git hosting, but it's 
standard first thing to do, in most cases.

I will assume your workspace is something like `/home/username/code`. First, 
open up a terminal and go there:

    $ cd code
    $ mkdir codec
    $ cd codec

Next, initialise the git repository

     $ git init 
     Initialized empty Git repository in /home/username/code/codec/.git/

If you are going to use github or so, you can instead just create the 
repository there and 

    $ git clone git@github.com/username/codec.git 

or something like this, whatever the git URL is. It is recommended to use 
ssh, it just makes life simpler and it's faster and no complications with 
authentication as you have if you use `https` instead.

### 2. Set up the folder hierarchy

You can do this with your IDE's file manager or you can use the terminal as 
you like, but start by making the folder structure that will be used:

     
    codec
    ├── cmd
    └── pkg
        ├── based32
        ├── codecer
        └── proto
    

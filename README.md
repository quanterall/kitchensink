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
not happen, each stage's parts will be also found in the [steps](. /steps)
folder at the root of the repository.

## Step By Step:

here will be the step by step process of building the library with a logical
sequence that builds from the basis to the specific parts for each in the order
that is needed both for understanding and for the constraints of syntax, grammar
and build system design...

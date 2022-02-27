# transcribe

## Human Readable Binary Encoding Framework and Tutorial

In this tutorial we will walk you through the creation from scratch of a human
readable encoding system, and to make it more interesting, give the option of
varying the details of the scheme produced.

The idea here is to make a tutorial that lets you go a lot deeper into the task
while giving a simple base to understand encoding bytes in forms that humans can
transcribe (theoretically)

Further, rather than just yielding a simple, concrete implementation, in the
design of this library it is written to show that one can write an extensible
library with very little extra work compared to the pure quick and dirty
implementation, if one understands a few simple principles.

This tutorial demonstrates the use of almost every possible and important
feature of Go. A "toy" implementation of a gRPC/protobuf microservice is added
in order to illustrate almost everything.

In order to demonstrate synchronisation primitives, waitgroups, atomics and
mutexes, the service will keep track of the number of invocations, print this
count in log updates, and track the count using a concurrent safe atomic
variable and show the variant using a mutex instead, and run an arbitrary 
number of concurrent worker threads that will start up and stop using 
waitgroups.

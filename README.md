# kitchensink

## Teaching Golang via building a Human Readable Binary Encoding Framework

In this tutorial we will walk you through the creation from scratch of a human
readable encoding system, and to make it more interesting, give the option of
varying the details of the scheme produced, how to turn a library into a 
microservice, including simple concurrency

This tutorial demonstrates the use of almost every possible and important
feature of Go. A "toy" implementation of a gRPC/protobuf microservice is added
in order to illustrate almost everything else.

In order to demonstrate synchronisation primitives, waitgroups, atomics and
mutexes, the service will keep track of the number of invocations, print this
count in log updates, and track the count using a concurrent safe atomic
variable and show the variant using a mutex instead, and run an arbitrary 
number of concurrent worker threads that will start up and stop using 
waitgroups.

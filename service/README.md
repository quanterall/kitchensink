# b32svc

This is an example of a small concurrent handler for the back end of a
microservice or similar.

Because the API of this example only has two functions, they are manually
specified and the handlers are defined inline, but for a larger API, it would
usually be done with the handler implementations separately defined, and then a
generator to string the API specifications together and create an API handler
that contains all the elements in this implementation.

Generators are an advanced subject, and to explain for those who are coming from
languages that use generics, these replace the function of generics.

The essence of how this concurrent handler works, is a service worker back end
starts up, spawns a number of worker goroutines, and each worker waits on
messages from the various channels, and by virtue of being single threads, each
worker will process one task at a time, and then check again if any channels are
waiting for new calls to be processed.

The advantage of this architecture is that, with the addition of a generator,
one can easily add and change API calls without touching many parts, and
centralising the changes and then running a generator to update the linkages.

This is not the network based RPC service itself, but the implementation back
end that receives messages and processes the calls. It can be used internally in
an application to provide services internally, or, as we demonstrate in the rest
of this tutorial/project, as the processing component that is run by a
microservice gRPC/Protobuf based HTTP(S) Remote Procedure Call (RPC) server.

You will notice that the API type is entirely composed of interface types. This
is because the concrete implementation will be filled out with wrapper functions
that assert the concrete types and all of this in-between is invisible to the
calling code.

When the interface is simple, it can be ok to simply assert the types in results
manually, but with a generator, this can all be completely automated, and after
any hardened generics language user sees what this means for maintaining, and
the lack of compilation time cost in resolving the generic types, the reason why
Go 1 does not have generics will become abundantly clear.

Generics shift a lot of heavy lifting for covering everything to the corner
cases, whereas in practice, as in this example, you see that the scope that is
required is a lot smaller and thus, this time waiting for the compiler is a
total waste of programmer's time, as well as an unnecessary cost in paying
programmers to do like this:

![xkcd #303 compiling](https://imgs.xkcd.com/comics/compiling.png)

Such resources are better spent on legitimately recreational, non-work,
non-coding things like beanbags, gym equipment and table tennis tables, where
programmers can relieve the mental strain of their work properly when needed,
not according to the limitations of the language.

Rust's Cargo build system does reduce a lot of the repetition of compilation,
but it still does not eliminate the problem of the time cost in compilation 
that, no different to Go's build system, increases with the amount of 
dependent code affected by changes - with Go, it will recompile a lot of 
things if, say, a logger, is changed, due to the number of callers to it, 
and if it is a generic using language, all of those dependencies have to be 
reprocessed. Whereas, in Go, you just run the generators, and off you go. No 
forced pauses waiting for endless machine intelligence processing over and 
over again.

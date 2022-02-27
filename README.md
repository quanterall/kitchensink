
# transcribe

## Human Readable Binary Encoding Framework and Tutorial

In this tutorial we will walk you through the creation from scratch of a 
human readable encoding system, and to make it more interesting, give the 
option of varying the details of the scheme produced.

In the base58 and bech32 directories are copies of the encoders used by 
`btcd`, the Go implementation of the Bitcoin full node.

The idea here is to make a tutorial that lets you go a lot deeper into the 
task while giving a simple base to understand encoding bytes in forms that 
humans can transcribe (theoretically)

Further, rather than just yielding a simple, concrete implementation, in the 
design of this library it is written to show that one can write an 
extensible library with very little extra work compared to the pure quick 
and dirty implementation, if one understands a few simple principles.

This tutorial will teach everything except for the use of concurrency, as in 
such a small library there is no reason for it.

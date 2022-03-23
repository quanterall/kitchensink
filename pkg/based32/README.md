# based32

An implementation of a generic encoding system for human readable 
transcription purposes, as usually found in cryptocurrencies to put 
addresses and transaction hashes into a form that humans can easily and 
reliably write down and read back, with minimal chance of error.

It uses the base32 character set used in Bech32, but differs from Bech32 in 
that it is not solely targeted at the purpose of cryptocurrency addresses, 
keys and transaction identifiers, or, more precisely, that it is written to 
enable the use for any length (though not unlimited length, due to human 
patience and attention limitations) for anything from 32 bit binary codes up 
to 512 bits or thereabouts.

In order to achieve this, it works around the limitation of the standard Go 
base32 codec, available in the standard library, of needing the values to be 
in 40 bit chunks (5 bytes) which creates 8 character long segments each 
representing 5 bits of data, and the ugliness of padding, by varying the 
length of the check value to act in a second purpose as also the padding.

It therefore requires a prefix noting the length of the check bytes at the 
end of the code, which is never more than 3 bytes, so this allows the first 
base32 cipher to be omitted and readded upon decoding.

This is a part of an introductory programming tutorial/workshop to teach Go 
programming, as well as all the accessory parts of the process including 
documentation, code style, and project structuring, the details of the build 
system, and so on.

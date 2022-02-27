// Package codec is a generalised encoder and decoder for encoding short
// binary values such as used in cryptocurrency addresses, accounts and
// transaction hashes.
//
// This codec is intended to be extended via implementations adding the
// functions specific to an implementation to the generalised implementation
// provided in this package that invokes these implementation specific features
// and configurations.
//
// This is part of a tutorial for teaching the correct way to create a Go
// library, which includes this file that provides the header for the godoc
// output.
//
// All exported symbols in Go code should have proper and informative comments
// added above them as this adds to the head of the package identifier. No
// released code should be without these on every identifier.
package codec

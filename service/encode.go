package b32svc

import (
	"errors"
	"time"
)

// The following functions are sets of four that follow a convention, and as you
// will see if you look closely, they are in fact identical, except for the
// return types. This is why if there is a lot of them, it is recommended to use
// a generator and replace the parameter types and return types as required.
//
// The convention is as follows:
//
// - Command - this simply delivers the command, it will block if the receive
//   buffer is full.
//
// - CommandChk - this polls the response, and returns true if it is ready
//
// - CommandGetRes - this should only be called if CommandChk returns true, and
//   it will return the result, either error or the result data as requested.
//
// - CommandWait - This will wait until the result and return it like a normal
//   synchronous function call. The command will expire after the timeout set in
//   the API.
//
// - CommandFlush - This
//
// The timeout can be made to be very long, but cannot be in practise infinite
// or it will become a resource leak.
//

// Encode calls the method with the given parameters.
//
// This is a blocking call, so if the call channel buffers are full it will wait
// until an item has been processed, it will stall until there is an empty
// buffer. It is better to increase the buffer size than to make this spawn a
// goroutine just to load a channel.
func (a API) Encode(cmd *EncodeCmd) {

	ServiceHandlers["Encode"].Call <- API{
		Ch:      a.Ch,
		Params:  cmd,
		Cancel:  a.Cancel,
		Timeout: a.Timeout,
	}
	return
}

// EncodeChk checks if a new message arrived on the result channel and returns
// true if it does, as well as storing the value in the Result field. If the
// return value is true, then the EncodeGetRes call extracts the return value.
func (a API) EncodeChk() (isNew bool) {

	select {
	case o := <-a.Ch.(chan EncodeRes):

		if o.Err != nil {

			a.Result = o.Err
		} else {

			a.Result = o.Res
		}
		isNew = true

	default:
		// the default section means if there is nothing in the channel the
		// return value is not mutated.
	}
	return
}

// EncodeGetRes extracts the result value. Should only be called if EncodeChk
// returns true.
func (a API) EncodeGetRes() (out string, e error) {

	switch t := a.Result.(type) {
	case *string:
		out = *t

	case error:
		e = t

	default:
		// This case only occurs if the programmer did not honor the contract of
		// this API, and will not be invoked in production.
		panic("GetRes function should not be called unless Chk returns true")
	}
	return
}

// EncodeWait blocks until it returns or API.Timeout seconds
// passes
func (a API) EncodeWait() (out *string, e error) {

	select {
	case <-time.After(a.Timeout):
		e = errors.New("timeout")
		break

	case <-a.Cancel:
		e = errors.New("canceled")
		break

	case o := <-a.Ch.(chan EncodeRes):
		out, e = o.Res, o.Err
	}
	return
}

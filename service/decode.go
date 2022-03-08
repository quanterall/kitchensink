package b32svc

import (
	"time"
)

// Decode calls the method with the given parameters
func (a API) Decode(cmd *DecodeCmd) (e error) {
	ServiceHandlers["Decode"].Call <- API{a.Ch, cmd, nil, nil, a.Timeout}
	return
}

// DecodeChk checks if a new message arrived on the result channel and
// returns true if it does, as well as storing the value in the Result field
func (a API) DecodeChk() (isNew bool) {
	select {
	case o := <-a.Ch.(chan DecodeRes):
		if o.Err != nil {
			a.Result = o.Err
		} else {
			a.Result = o.Res
		}
		isNew = true
	default:
	}
	return
}

// DecodeGetRes returns a pointer to the value in the Result field
func (a API) DecodeGetRes() (out *None, e error) {
	out, _ = a.Result.(*None)
	e, _ = a.Result.(error)
	return
}

// DecodeWait calls the method and blocks until it returns or 5 seconds passes
func (a API) DecodeWait() (out []byte, e error) {
	select {
	case <-time.After(a.Timeout):
		break
	case o := <-a.Ch.(chan DecodeRes):
		out, e = o.Res, o.Err
	}
	return
}

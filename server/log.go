package server

import (
	logg "log"
	"os"
)

var log = logg.New(os.Stderr, "b32 ", logg.Llongfile)

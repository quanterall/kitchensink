package client

import (
	logg "log"
	"os"
)

var log = logg.New(os.Stderr, "based32 ", logg.Llongfile)

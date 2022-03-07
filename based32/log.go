package based32

import (
	logg "log"
	"os"
)

// By including an exact copy of this file in all new packages, errors have
// their source locations printed with the log items automatically, simplifying
// locating the errors, and to encourage the policy of logging at the site of
// the error and not using the stupid 'wrap' bunk.
var log = logg.New(os.Stderr, "based32 ", logg.Llongfile)

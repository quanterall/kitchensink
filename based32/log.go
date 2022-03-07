package based32

import (
	logg "log"
	"os"
)

var log = logg.New(os.Stderr, "based32 ", logg.Llongfile)

// By including an exact copy of this file in all new packages with the package
// came changed,, errors have their source locations printed with the log items
// automatically, simplifying locating the errors, and to encourage the policy
// of logging at the site of the error and not using the stupid 'wrap' bunk.
//
// Unfortunately, due to, most likely, closed source and security by obscurity
// mentality, developers want to hide their top secret source code filesystem
// locations, as though there is any benefit in this, or that closed source is
// anything but a potential cover for compromising the fools who run such
// software, and providing a back door to 'trusted authorities' while also
// providing a backdoor to anyone who manages to find it, and once it's out of
// the bag, it's a vulnerability for everyone.
//
// So, including this file should be mandatory for every application, unless the
// Go standard library is fixed to expose the possibility of setting the flags
// via init() functions or in the main() of an application.
//
// An alternative is to fork the standard log package, as there are far more
// options than just changing the writer, and without this option to change the
// entire back end implementation with less complicated methods than
// demonstrated here, it would then be possible to change it to log remotely, to
// log with a structured syntax for a log filtering database, that could also
// track also concurrent log entries in real time. But for now, for the benefit
// of those trained via Quanterall's Go training system, simply copy and paste
// of the prior content, we can at least debug with less mystery about where
// errors are coming from, and develop the habit of using logging at the site of
// errors instead of passing them through layers of indirection such as
// interface implementations, which are not necessarily resolved by the
// developers chosen code editor, but implementing a simple hyperlink to jump to
// errors, or at least to make the code location transparent, means debugging is a lot easier.
//
// Debugging, of course, only is required when the programmer is writing
// *algorithms*. Writing scripts rarely produces errors that are hard to find.

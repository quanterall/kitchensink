# Useful Scripts

In this project and tutorial we are attempting to avoid as much as possible the
use of anything other that Go tooling to do everything.

In this folder are a few small things that require shell scripts, and assume you
are running Ubuntu 20, or Debian Sid, or something similar to this, and have
installed Go according to the instructions in [the readme](../README.md).

- [Files in this folder:](#files-in-this-folder)
	- [installtoc.sh](#installtocsh)
	- [toc.sh](#tocsh)

## Files in this folder:

### [installtoc.sh](./installtoc.sh)

installs the ubuntu/debian package `fswatch` which is used by Goland to watch
for files to run 'run on save' file watchers, and the Go program `tocenize`,
which creates a table of contents in markdown files based on header levels and
texts. `fswatch` is needed by `tocenize` for one of its run modes, where it 
sits and waits for files to be changed and then runs an update on the ToC of 
the file.

### [toc.sh](./toc.sh)

is a simple bash script that recursively runs `tocenize` on all markdown files
in the repository to update the tables of contents in each markdown file in the
package.

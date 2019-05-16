// +build !plan9

package wrapcommander

import (
	"os"
	"syscall"
)

// IsExecFormatError returns a boolean indicating whether the error is known to
// report that format of an executable file is invalid.
// ex. "fork/exec ./prog: exec format error"
func IsExecFormatError(err error) bool {
	e, ok := err.(*os.PathError)
	return ok && e.Err == syscall.ENOEXEC
}
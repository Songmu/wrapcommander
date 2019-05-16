// +build plan9

package wrapcommander

import (
	"os"
	"strings"
)

var (
	badExec = "exec header invalid"
	dirExec = "cannot exec directory"
)

// IsExecFormatError returns a boolean indicating whether the error is known to
// report that format of an executable file is invalid.
// ex. "fork/exec ./prog: exec format error"
func IsExecFormatError(err error) bool {
	e, ok := err.(*os.PathError)
	if !ok {
		return false
	}
	return strings.HasSuffix(e.Error(), badExec) ||
		strings.HasSuffix(e.Error(), dirExec)
}

// +build plan9 windows

package wrapcommander

import "syscall"

func resolveExitCode(w syscall.WaitStatus) int {
	return w.ExitStatus()
}

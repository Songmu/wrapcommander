// +build windows

package wrapcommander

import "syscall"

// WaitStatus represents exit status of a process.
type WaitStatus = syscall.WaitStatus

func waitStatusToExitCode(w WaitStatus) int {
	return w.ExitStatus()
}

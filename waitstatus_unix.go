// +build !plan9,!windows

package wrapcommander

import "syscall"

// WaitStatus represents exit status of a process.
type WaitStatus = syscall.WaitStatus

func waitStatusToExitCode(w WaitStatus) int {
	if w.Signaled() {
		return int(w.Signal()) + 128
	}
	return w.ExitStatus()
}

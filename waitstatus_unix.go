// +build !plan9,!windows

package wrapcommander

import "syscall"

type WaitStatus = syscall.WaitStatus

func waitStatusToExitCode(w WaitStatus) int {
	if w.Signaled() {
		return int(w.Signal()) + 128
	}
	return w.ExitStatus()
}

// +build windows

package wrapcommander

import "syscall"

type WaitStatus = syscall.WaitStatus

func waitStatusToExitCode(w WaitStatus) int {
	return w.ExitStatus()
}

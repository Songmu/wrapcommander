// +build plan9

package wrapcommander

import "syscall"

type WaitStatus = syscall.Waitmsg

func waitStatusToExitCode(w WaitStatus) int {
	return w.ExitStatus()
}

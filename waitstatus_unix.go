// +build !plan9,!windows

package wrapcommander

import "syscall"

func resolveExitCode(w syscall.WaitStatus) int {
	if w.Signaled() {
		return w.ExitStatus() | 0x80
	}
	return w.ExitStatus()
}

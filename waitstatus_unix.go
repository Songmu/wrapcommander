// +build !plan9,!windows

package wrapcommander

import "syscall"

func resolveExitCode(w syscall.WaitStatus) int {
	if w.Signaled() {
		return int(w) | 0x80
	}
	return w.ExitStatus()
}

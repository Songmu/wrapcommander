// +build !plan9

package wrapcommander

import "syscall"

type Signal = syscall.Signal

func detectSignal(w WaitStatus) (signaled bool, signal Signal) {
	signaled = w.Signaled()
	signal = w.Signal()
	return
}

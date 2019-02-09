package wrapcommander

import (
	"os"
	"os/exec"
	"syscall"
)

// ExitStatus represents command exit status information
type ExitStatus struct {
	err error

	exitCode          int
	signaled, invoked bool
	signal            syscall.Signal
}

// Err returns original error
func (es *ExitStatus) Err() error {
	return es.err
}

// Signaled or not
func (es *ExitStatus) Signaled() bool {
	return es.signaled
}

// Invoked or not
func (es *ExitStatus) Invoked() bool {
	return es.invoked
}

// Signal returns a received signal
func (es *ExitStatus) Signal() syscall.Signal {
	return es.signal
}

// ResolveExitStatus resolve ExitStatus from command error
func ResolveExitStatus(err error) (es *ExitStatus) {
	if err == nil {
		es.invoked = true
		return es
	}
	es.err = err
	eerr, ok := err.(*exec.ExitError)
	es.invoked = ok

	if !es.invoked {
		switch {
		case os.IsPermission(err), IsExecFormatError(err):
			es.exitCode = ExitCommandNotInvoked
		case os.IsNotExist(err), IsNotFoundInPATH(err):
			es.exitCode = ExitCommandNotFound
		default:
			es.exitCode = ExitUnknownErr
		}
		return es
	}

	w, ok := eerr.Sys().(syscall.WaitStatus)
	if !ok {
		es.exitCode = -1
		return es
	}
	es.signaled = w.Signaled()
	es.signal = w.Signal()
	es.exitCode = waitStatusToExitCode(w)

	return es
}

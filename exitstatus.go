package wrapcommander

import (
	"os"
	"os/exec"
)

// ExitStatus represents command exit status information
type ExitStatus struct {
	err error

	exitCode          int
	signaled, invoked bool
	signal            Signal
}

// Err returns original error
func (es *ExitStatus) Err() error {
	return es.err
}

// ExitCode of the command
func (es *ExitStatus) ExitCode() int {
	return es.exitCode
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
func (es *ExitStatus) Signal() Signal {
	return es.signal
}

// ResolveExitStatus resolve ExitStatus from command error
func ResolveExitStatus(err error) *ExitStatus {
	es := &ExitStatus{
		invoked:  true,
		err:      err,
		exitCode: -1,
	}
	if es.err == nil {
		es.exitCode = 0
		return es
	}

	eerr, ok := es.err.(*exec.ExitError)
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

	w, ok := eerr.Sys().(WaitStatus)
	if !ok {
		return es
	}
	es.signaled, es.signal = detectSignal(w)
	es.exitCode = waitStatusToExitCode(w)

	return es
}

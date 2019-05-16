// +build plan9

package wrapcommander

import "strings"

// Signal represents POSIX signal.
type Signal int

// These numbers are imported from /sys/include/ape/signal.h
const (
	_SIGHUP  Signal = 1
	_SIGINT         = 2
	_SIGQUIT        = 3
	_SIGILL         = 4
	_SIGABRT        = 5
	_SIGFPE         = 6
	_SIGKILL        = 7
	_SIGSEGV        = 8
	_SIGPIPE        = 9
	_SIGALRM        = 10
	_SIGTERM        = 11
	_SIGUSR1        = 12
	_SIGUSR2        = 13
)

// This table is imported from /sys/src/ape/lib/ap/plan9/signal.c
var sigtab = []struct {
	Msg string
	Sig Signal
}{
	{"hangup", _SIGHUP},
	{"interrupt", _SIGINT},
	{"quit", _SIGQUIT},
	{"alarm", _SIGALRM},
	{"sys: trap: illegal instruction", _SIGILL},
	{"sys: trap: reserved instruction", _SIGILL},
	{"sys: trap: reserved", _SIGILL},
	{"sys: trap: arithmetic overflow", _SIGFPE},
	{"abort", _SIGABRT},
	{"sys: fp:", _SIGFPE},
	{"exit", _SIGKILL},
	{"die", _SIGKILL},
	{"kill", _SIGKILL},
	{"sys: trap: bus error", _SIGSEGV},
	{"sys: trap: address error", _SIGSEGV},
	{"sys: trap: TLB", _SIGSEGV},
	{"sys: write on closed pipe", _SIGPIPE},
	{"alarm", _SIGALRM},
	{"term", _SIGTERM},
	{"usr1", _SIGUSR1},
	{"usr2", _SIGUSR2},
}

func lookupSignal(msg string) Signal {
	for _, sig := range sigtab {
		if strings.Contains(msg, sig.Msg) {
			return sig.Sig
		}
	}
	return 0
}

func detectSignal(w WaitStatus) (bool, Signal) {
	msg := causeMsg(w.Msg)
	if msg == "" {
		return false, 0
	}
	if sig := lookupSignal(msg); sig != 0 {
		return true, sig
	}
	return false, 0
}

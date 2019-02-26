// +build plan9

package wrapcommander

import "strings"

type Signal int

const (
	// see /sys/include/ape/signal.h
	SIGHUP  Signal = 1
	SIGINT         = 2
	SIGQUIT        = 3
	SIGILL         = 4
	SIGABRT        = 5
	SIGFPE         = 6
	SIGKILL        = 7
	SIGSEGV        = 8
	SIGPIPE        = 9
	SIGALRM        = 10
	SIGTERM        = 11
	SIGUSR1        = 12
	SIGUSR2        = 13
)

var sigtab = []struct {
	Msg string
	Sig Signal
}{
	{"hangup", SIGHUP},
	{"interrupt", SIGINT},
	{"quit", SIGQUIT},
	{"alarm", SIGALRM},
	{"sys: trap: illegal instruction", SIGILL},
	{"sys: trap: reserved instruction", SIGILL},
	{"sys: trap: reserved", SIGILL},
	{"sys: trap: arithmetic overflow", SIGFPE},
	{"abort", SIGABRT},
	{"sys: fp:", SIGFPE},
	{"exit", SIGKILL},
	{"die", SIGKILL},
	{"kill", SIGKILL},
	{"sys: trap: bus error", SIGSEGV},
	{"sys: trap: address error", SIGSEGV},
	{"sys: trap: TLB", SIGSEGV},
	{"sys: write on closed pipe", SIGPIPE},
	{"alarm", SIGALRM},
	{"term", SIGTERM},
	{"usr1", SIGUSR1},
	{"usr2", SIGUSR2},
}

func detectSignal(w WaitStatus) (signaled bool, signal Signal) {
	for _, sig := range sigtab {
		if strings.HasPrefix(w.Msg, sig.Msg) {
			signaled = true
			signal = sig.Sig
		}
	}
	return
}

// +build plan9

package wrapcommander

import (
	"strconv"
	"strings"
	"syscall"
)

// WaitStatus represents exit status of a process.
type WaitStatus = *syscall.Waitmsg

/*
 * Plan 9 uses error string instead of exit code.
 * For example, errstr contains 'go: interrupt'
 * when process is catched a interrupt. (aka SIGINT)
 *
 * However, programs that is written for ANSI/POSIX still uses exit code.
 * In this case, errstr contains a number, like 'go: 2'
 */

const signalBase = 128

func waitStatusToExitCode(w WaitStatus) int {
	msg := causeMsg(w.Msg)
	if msg == "" {
		return 0
	}

	if sig := lookupSignal(msg); sig != 0 {
		return signalBase + int(sig)
	}
	n, err := strconv.Atoi(msg)
	if err != nil {
		return 1
	}
	return n
}

func causeMsg(msg string) string {
	i := strings.LastIndex(msg, ":")
	if i < 0 {
		return msg
	}
	s := strings.TrimSpace(msg[i+1:])
	if n := len(s); n > 0 && s[n-1] == '\'' {
		s = s[:n-1]
	}
	if s == "" {
		return msg
	}
	return s
}

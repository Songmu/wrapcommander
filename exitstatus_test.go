package wrapcommander

import (
	"os/exec"
	"runtime"
	"testing"
)

func TestResolveExitStatus(t *testing.T) {
	var testCases = []struct {
		cmd  string
		want int
	}{
		{"go", 2},
		{"gogogo-dummy", ExitCommandNotFound},
		{"./gogogo-dummy", ExitCommandNotFound},
		{"./testdata/dir", ExitCommandNotInvoked},
		{"./testdata/execformaterror", ExitCommandNotInvoked},
		{"./testdata/echo.sh", 0},
		{"./testdata/exit1.sh", 1},
	}

	if runtime.GOOS == "windows" {
		t.Skip("not supported on windows")
	}
	for _, tt := range testCases {
		err := exec.Command(tt.cmd).Run()
		es := ResolveExitStatus(err)
		got := es.ExitCode()
		if got != tt.want {
			t.Errorf("ResolveExitStatus(exec.Command(%#v).Run()) = %v; want %v", tt.cmd, got, tt.want)
		}
	}
}

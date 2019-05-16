package wrapcommander

import (
	"os"
	"os/exec"
	"runtime"
	"testing"
)

func TestResolveExitStatus(t *testing.T) {
	var ext string
	switch runtime.GOOS {
	case "windows":
		t.Skip("not supported on windows")
	case "plan9":
		ext = ".rc"
	default:
		ext = ".sh"
	}

	var testCases = []struct {
		cmd  string
		want int
	}{
		{"go", 2},
		{"gogogo-dummy", ExitCommandNotFound},
		{"./gogogo-dummy", ExitCommandNotFound},
		{"./testdata/dir", ExitCommandNotInvoked},
		{"./testdata/execformaterror", ExitCommandNotInvoked},
		{"./testdata/echo" + ext, 0},
		{"./testdata/exit1" + ext, 1},
	}

	for _, tt := range testCases {
		err := exec.Command(tt.cmd).Run()
		es := ResolveExitStatus(err)
		got := es.ExitCode()
		if got != tt.want {
			t.Errorf("ResolveExitStatus(exec.Command(%#v).Run()) = %v; want %v", tt.cmd, got, tt.want)
		}
		if es.Signaled() {
			t.Errorf("something went wrong")
		}
		if !es.Invoked() && es.Signal() != 0 {
			t.Errorf("something went wrong")
		}
	}
}

func TestResolveExitStatus_sig(t *testing.T) {
	if runtime.GOOS == "windows" {
		t.Skip("not supported on windows")
	}
	cmd := exec.Command("sleep", "10")
	cmd.Start()
	cmd.Process.Signal(os.Interrupt)
	es := ResolveExitStatus(cmd.Wait())
	if es.Signal() != 2 {
		t.Errorf("something went wrong")
	}
	if es.ExitCode() != 128+2 {
		t.Errorf("something went wrong")
	}
	if es.Err() == nil {
		t.Errorf("something went wrong")
	}
}

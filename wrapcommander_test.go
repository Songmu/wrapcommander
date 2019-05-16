package wrapcommander

import (
	"os"
	"os/exec"
	"reflect"
	"runtime"
	"testing"
)

type testCmd struct {
	cmd  string
	want bool
}

var isPermissionTests = []testCmd{
	{"go", false},
	{"./testdata/dir", true},
	{"./testdata/nopermission", true},
	{"./testdata/execformaterror", false},
}

func TestIsPermission(t *testing.T) {
	if runtime.GOOS == "windows" || runtime.GOOS == "plan9" {
		t.Skip("not supported on windows or plan9")
	}
	for _, tt := range isPermissionTests {
		err := exec.Command(tt.cmd).Run()
		if got := IsPermission(err); got != tt.want {
			t.Errorf("IsPermission(exec.Command(%#v).Run()) = %v; want %v", tt.cmd, got, tt.want)
		}
	}
}

var isNotFoundInPATHTests = []testCmd{
	{"go", false},
	{"gogogo-dummy", true},
	{"./testdata/dir", false},
}

func TestIsNotFoundInPATH(t *testing.T) {
	for _, tt := range isNotFoundInPATHTests {
		err := exec.Command(tt.cmd).Run()
		if got := IsNotFoundInPATH(err); got != tt.want {
			t.Errorf("IsNotFoundInPATH(exec.Command(%#v).Run()) = %v; want %v", tt.cmd, got, tt.want)
		}
	}
	if isNotFound := IsNotFoundInPATH(nil); isNotFound {
		t.Errorf("IsNotFoundInPATH(nil) = %v; want false", isNotFound)
	}
}

func TestIsExecFormatErrorTests(t *testing.T) {
	var ext string
	switch runtime.GOOS {
	case "windows":
		t.Skip("not supported on windows")
	case "plan9":
		ext = ".rc"
	default:
		ext = ".sh"
	}
	var isExecFormatErrorTests = []testCmd{
		{"go", false},
		{"gogogo-dummy", false},
		{"./testdata/dir", false},
		{"./testdata/echo" + ext, false},
		{"./testdata/execformaterror", true},
	}
	for _, tt := range isExecFormatErrorTests {
		err := exec.Command(tt.cmd).Run()
		if got := IsExecFormatError(err); got != tt.want {
			t.Errorf("IsExecFormatError(exec.Command(%#v).Run()) = %v; want %v", tt.cmd, got, tt.want)
		}
	}
	if isExecFormatError := IsExecFormatError(nil); isExecFormatError {
		t.Errorf("IsNotFoundInPATH(nil) = %v; want false", isExecFormatError)
	}
}

func TestResolveExitCode(t *testing.T) {
	var ext string
	switch runtime.GOOS {
	case "windows":
		t.Skip("not supported on windows")
	case "plan9":
		ext = ".rc"
	default:
		ext = ".sh"
	}
	var resolveExitCodeTests = []struct {
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

	for _, tt := range resolveExitCodeTests {
		err := exec.Command(tt.cmd).Run()
		got := ResolveExitCode(err)
		if got != tt.want {
			t.Errorf("ResolveExitCode(exec.Command(%#v).Run()) = %v; want %v", tt.cmd, got, tt.want)
		}
	}
}

func TestResolveExitCode_sig(t *testing.T) {
	if runtime.GOOS == "windows" {
		t.Skip("not supported on windows")
	}
	cmd := exec.Command("sleep", "10")
	cmd.Start()
	cmd.Process.Signal(os.Interrupt)
	err := cmd.Wait()
	ex := ResolveExitCode(err)
	if ex != 128+2 {
		t.Errorf("something went wrong")
	}
}

var separateArgsTests = []struct {
	args     []string
	optsArgs []string
	cmdArgs  []string
}{
	{
		[]string{"a", "b", "c"},
		[]string{},
		[]string{"a", "b", "c"},
	},
	{
		[]string{"a", "b", "--", "c"},
		[]string{"a", "b"},
		[]string{"c"},
	},
	{
		[]string{"--", "c", "b"},
		[]string{},
		[]string{"c", "b"},
	},
	{
		[]string{"c", "b", "--"},
		[]string{},
		[]string{"c", "b", "--"},
	},
	{
		[]string{"aa", "--", "c", "b", "--", "bb"},
		[]string{"aa"},
		[]string{"c", "b", "--", "bb"},
	},
}

func TestSeparateArgs(t *testing.T) {
	for _, tt := range separateArgsTests {
		optsArgs, cmdArgs := SeparateArgs(tt.args)
		if !reflect.DeepEqual(tt.optsArgs, optsArgs) {
			t.Errorf("SeparateArgs(%#v)[0] = %+v; want %+v", tt.args, optsArgs, tt.optsArgs)
		}
		if !reflect.DeepEqual(tt.cmdArgs, cmdArgs) {
			t.Errorf("SeparateArgs(%#v)[1] = %+v; want %+v", tt.args, cmdArgs, tt.cmdArgs)
		}
	}
}

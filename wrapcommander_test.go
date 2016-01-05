package wrapcommander

import (
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
	{"./_testdata/dir", true},
	{"./_testdata/nopermission", true},
	{"./_testdata/execformaterror", false},
}

func TestIsPermission(t *testing.T) {
	if runtime.GOOS == "windows" {
		t.Skip("not supported on windows")
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
	{"./_testdata/dir", false},
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

var isExecFormatErrorTests = []testCmd{
	{"go", false},
	{"gogogo-dummy", false},
	{"./_testdata/dir", false},
	{"./_testdata/echo.sh", false},
	{"./_testdata/execformaterror", true},
}

func TestIsExecFormatErrorTests(t *testing.T) {
	for _, tt := range isExecFormatErrorTests {
		if runtime.GOOS == "windows" {
			t.Log("not supported on windows")
			continue
		}
		err := exec.Command(tt.cmd).Run()
		if got := IsExecFormatError(err); got != tt.want {
			t.Errorf("IsExecFormatError(exec.Command(%#v).Run()) = %v; want %v", tt.cmd, got, tt.want)
		}
	}
	if isExecFormatError := IsExecFormatError(nil); isExecFormatError {
		t.Errorf("IsNotFoundInPATH(nil) = %v; want false", isExecFormatError)
	}
}

var resolveExitCodeTests = []struct {
	cmd  string
	want int
}{
	{"go", 2},
	{"gogogo-dummy", ExitCommandNotFound},
	{"./gogogo-dummy", ExitCommandNotFound},
	{"./_testdata/dir", ExitCommandNotInvoked},
	{"./_testdata/execformaterror", ExitCommandNotInvoked},
	{"./_testdata/echo.sh", 0},
	{"./_testdata/exit1.sh", 1},
}

func TestResolveExitCode(t *testing.T) {
	if runtime.GOOS == "windows" {
		t.Skip("not supported on windows")
	}
	for _, tt := range resolveExitCodeTests {
		err := exec.Command(tt.cmd).Run()
		got := ResolveExitCode(err)
		if got != tt.want {
			t.Errorf("ResolveExitCode(exec.Command(%#v).Run()) = %v; want %v", tt.cmd, got, tt.want)
		}
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

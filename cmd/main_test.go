package main

import (
	"os"
	"os/exec"
	"path/filepath"
	"strings"
	"testing"
)

var binary string

func TestMain(m *testing.M) {
	tmp, err := os.MkdirTemp("", "clevis-test")
	if err != nil {
		panic(err)
	}
	defer os.RemoveAll(tmp)
	binary = filepath.Join(tmp, "clevis")
	out, err := exec.Command("go", "build", "-o", binary, ".").CombinedOutput()
	if err != nil {
		panic(string(out))
	}
	os.Exit(m.Run())
}

func run(t *testing.T, stdin string, args ...string) (string, string, int) {
	t.Helper()
	cmd := exec.Command(binary, args...)
	if stdin != "" {
		cmd.Stdin = strings.NewReader(stdin)
	}
	var stdout, stderr strings.Builder
	cmd.Stdout = &stdout
	cmd.Stderr = &stderr
	err := cmd.Run()
	code := 0
	if exitErr, ok := err.(*exec.ExitError); ok {
		code = exitErr.ExitCode()
	} else if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	return stdout.String(), stderr.String(), code
}

func TestHelp(t *testing.T) {
	for _, flag := range []string{"--help", "-h"} {
		_, stderr, code := run(t, "", flag)
		if code != 0 {
			t.Errorf("%s: exit code %d, want 0", flag, code)
		}
		if !strings.Contains(stderr, "Usage:") {
			t.Errorf("%s: expected usage output, got: %s", flag, stderr)
		}
	}
}

func TestVersion(t *testing.T) {
	stdout, _, code := run(t, "", "--version")
	if code != 0 {
		t.Fatalf("exit code %d, want 0", code)
	}
	if !strings.HasPrefix(stdout, "v") || !strings.Contains(stdout, "anatol/clevis") {
		t.Errorf("unexpected version output: %s", stdout)
	}
}

func TestMissingCommand(t *testing.T) {
	_, stderr, code := run(t, "")
	if code != 2 {
		t.Errorf("exit code %d, want 2", code)
	}
	if !strings.Contains(stderr, "Missing command") {
		t.Errorf("expected missing command error, got: %s", stderr)
	}
}

func TestInvalidCommand(t *testing.T) {
	_, stderr, code := run(t, "", "bogus")
	if code != 2 {
		t.Errorf("exit code %d, want 2", code)
	}
	if !strings.Contains(stderr, "Invalid command") {
		t.Errorf("expected invalid command error, got: %s", stderr)
	}
}

func TestCommandAliases(t *testing.T) {
	cases := []struct {
		name  string
		args  []string
		stdin string
		code  int
	}{
		{"encrypt alias", []string{"e"}, "", 2},
		{"decrypt alias", []string{"d"}, "bad", 1},
		{"inspect alias", []string{"i"}, "bad", 1},
	}
	for _, tc := range cases {
		t.Run(tc.name, func(t *testing.T) {
			_, stderr, code := run(t, tc.stdin, tc.args...)
			if code != tc.code {
				t.Errorf("exit code %d, want %d: %s", code, tc.code, stderr)
			}
		})
	}
}

func TestEncryptMissingArgs(t *testing.T) {
	cases := []struct {
		name string
		args []string
	}{
		{"no args", []string{"encrypt"}},
		{"pin only", []string{"encrypt", "tang"}},
		{"empty pin", []string{"encrypt", "", "{}"}},
		{"empty config", []string{"encrypt", "tang", ""}},
	}
	for _, tc := range cases {
		t.Run(tc.name, func(t *testing.T) {
			_, stderr, code := run(t, "", tc.args...)
			if code != 2 {
				t.Errorf("exit code %d, want 2", code)
			}
			if !strings.Contains(stderr, "Missing arguments") {
				t.Errorf("expected missing arguments error, got: %s", stderr)
			}
		})
	}
}

func TestDecryptBadInput(t *testing.T) {
	_, stderr, code := run(t, "not-valid-jwe", "decrypt")
	if code != 1 {
		t.Errorf("exit code %d, want 1", code)
	}
	if !strings.Contains(stderr, "Error:") {
		t.Errorf("expected error output, got: %s", stderr)
	}
}

func TestInspectValidJWE(t *testing.T) {
	jwe := "eyJhbGciOiJBMjU2S1ciLCJlbmMiOiJBMjU2R0NNIn0.q5nU6Po_Butjl-Vrc-tThCCF4-KVkPmqjrlTj7Pr1LCB501OlAQbpA.zUzuAk2g3N6emlBF.qufuDg_FgnbSvc8t.gWmDDyOklqdus1lquhW4kw"
	stdout, _, code := run(t, jwe, "inspect")
	if code != 0 {
		t.Fatalf("exit code %d, want 0", code)
	}
	if !strings.Contains(stdout, "A256KW") {
		t.Errorf("expected JWE dump with algorithm, got: %s", stdout)
	}
}

func TestLuksUnsupported(t *testing.T) {
	_, stderr, code := run(t, "", "luks")
	if code != 1 {
		t.Errorf("exit code %d, want 1", code)
	}
	if !strings.Contains(stderr, "not currently supported") {
		t.Errorf("expected unsupported message, got: %s", stderr)
	}
}

func TestHelpAfterSubcommand(t *testing.T) {
	for _, args := range [][]string{
		{"decrypt", "--help"},
		{"encrypt", "-h", "tang", "{}"},
		{"inspect", "--help"},
	} {
		_, stderr, code := run(t, "", args...)
		if code != 0 {
			t.Errorf("%v: exit code %d, want 0", args, code)
		}
		if !strings.Contains(stderr, "Usage:") {
			t.Errorf("%v: expected usage output, got: %s", args, stderr)
		}
	}
}

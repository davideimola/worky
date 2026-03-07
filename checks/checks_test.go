package checks_test

import (
	"context"
	"fmt"
	"net"
	"net/http"
	"net/http/httptest"
	"os"
	"path/filepath"
	"testing"

	"github.com/davideimola/worky/checks"
)

// helpers

func writeFile(t *testing.T, content string) string {
	t.Helper()
	f, err := os.CreateTemp(t.TempDir(), "worky-test-*")
	if err != nil {
		t.Fatal(err)
	}
	if _, err := f.WriteString(content); err != nil {
		t.Fatal(err)
	}
	f.Close()
	return f.Name()
}

var bg = context.Background()

// FileExists

func TestFileExists_Pass(t *testing.T) {
	path := writeFile(t, "hello")
	if err := checks.FileExists(path)(bg); err != nil {
		t.Fatalf("expected nil, got %v", err)
	}
}

func TestFileExists_Fail(t *testing.T) {
	if err := checks.FileExists("/no/such/file")(bg); err == nil {
		t.Fatal("expected error for missing file")
	}
}

// DirExists

func TestDirExists_Pass(t *testing.T) {
	dir := t.TempDir()
	if err := checks.DirExists(dir)(bg); err != nil {
		t.Fatalf("expected nil, got %v", err)
	}
}

func TestDirExists_Fail_Missing(t *testing.T) {
	if err := checks.DirExists("/no/such/dir")(bg); err == nil {
		t.Fatal("expected error for missing dir")
	}
}

func TestDirExists_Fail_IsFile(t *testing.T) {
	path := writeFile(t, "not a dir")
	if err := checks.DirExists(path)(bg); err == nil {
		t.Fatal("expected error when path is a file")
	}
}

// FileContains

func TestFileContains_Pass(t *testing.T) {
	path := writeFile(t, "hello world")
	if err := checks.FileContains(path, "world")(bg); err != nil {
		t.Fatalf("expected nil, got %v", err)
	}
}

func TestFileContains_Fail(t *testing.T) {
	path := writeFile(t, "hello world")
	if err := checks.FileContains(path, "nope")(bg); err == nil {
		t.Fatal("expected error when text not found")
	}
}

func TestFileContains_MissingFile(t *testing.T) {
	if err := checks.FileContains("/no/such/file", "x")(bg); err == nil {
		t.Fatal("expected error for missing file")
	}
}

// FileMatchesRegex

func TestFileMatchesRegex_Pass(t *testing.T) {
	path := writeFile(t, "version: 1.2.3")
	if err := checks.FileMatchesRegex(path, `version: \d+\.\d+\.\d+`)(bg); err != nil {
		t.Fatalf("expected nil, got %v", err)
	}
}

func TestFileMatchesRegex_Fail(t *testing.T) {
	path := writeFile(t, "hello world")
	if err := checks.FileMatchesRegex(path, `\d+`)(bg); err == nil {
		t.Fatal("expected error when pattern doesn't match")
	}
}

func TestFileMatchesRegex_InvalidPattern(t *testing.T) {
	path := writeFile(t, "anything")
	// Invalid regex should return an error, not panic.
	err := checks.FileMatchesRegex(path, `[invalid`)(bg)
	if err == nil {
		t.Fatal("expected error for invalid regex")
	}
}

// EnvVarSet

func TestEnvVarSet_Pass(t *testing.T) {
	t.Setenv("WORKY_TEST_VAR", "hello")
	if err := checks.EnvVarSet("WORKY_TEST_VAR")(bg); err != nil {
		t.Fatalf("expected nil, got %v", err)
	}
}

func TestEnvVarSet_Fail_Unset(t *testing.T) {
	os.Unsetenv("WORKY_TEST_UNSET")
	if err := checks.EnvVarSet("WORKY_TEST_UNSET")(bg); err == nil {
		t.Fatal("expected error for unset var")
	}
}

func TestEnvVarSet_Fail_Empty(t *testing.T) {
	t.Setenv("WORKY_TEST_EMPTY", "")
	if err := checks.EnvVarSet("WORKY_TEST_EMPTY")(bg); err == nil {
		t.Fatal("expected error for empty var")
	}
}

// EnvVarEquals

func TestEnvVarEquals_Pass(t *testing.T) {
	t.Setenv("WORKY_EQ", "expected")
	if err := checks.EnvVarEquals("WORKY_EQ", "expected")(bg); err != nil {
		t.Fatalf("expected nil, got %v", err)
	}
}

func TestEnvVarEquals_Fail(t *testing.T) {
	t.Setenv("WORKY_EQ", "wrong")
	if err := checks.EnvVarEquals("WORKY_EQ", "expected")(bg); err == nil {
		t.Fatal("expected error for wrong value")
	}
}

// CommandSucceeds

func TestCommandSucceeds_Pass(t *testing.T) {
	if err := checks.CommandSucceeds("true")(bg); err != nil {
		t.Fatalf("expected nil, got %v", err)
	}
}

func TestCommandSucceeds_Fail(t *testing.T) {
	if err := checks.CommandSucceeds("false")(bg); err == nil {
		t.Fatal("expected error for failing command")
	}
}

// CommandOutputContains

func TestCommandOutputContains_Pass(t *testing.T) {
	if err := checks.CommandOutputContains("hello", "echo", "hello world")(bg); err != nil {
		t.Fatalf("expected nil, got %v", err)
	}
}

func TestCommandOutputContains_Fail_TextMissing(t *testing.T) {
	if err := checks.CommandOutputContains("nope", "echo", "hello world")(bg); err == nil {
		t.Fatal("expected error when text not in output")
	}
}

func TestCommandOutputContains_Fail_Command(t *testing.T) {
	if err := checks.CommandOutputContains("anything", "false")(bg); err == nil {
		t.Fatal("expected error for failing command")
	}
}

// PortOpen

func TestPortOpen_Pass(t *testing.T) {
	ln, err := net.Listen("tcp", "127.0.0.1:0")
	if err != nil {
		t.Fatal(err)
	}
	defer ln.Close()
	port := ln.Addr().(*net.TCPAddr).Port
	if err := checks.PortOpen("127.0.0.1", port)(bg); err != nil {
		t.Fatalf("expected nil, got %v", err)
	}
}

func TestPortOpen_Fail(t *testing.T) {
	// Find a free port and don't listen on it.
	ln, _ := net.Listen("tcp", "127.0.0.1:0")
	port := ln.Addr().(*net.TCPAddr).Port
	ln.Close()
	if err := checks.PortOpen("127.0.0.1", port)(bg); err == nil {
		t.Fatal("expected error for closed port")
	}
}

// HTTPStatus

func TestHTTPStatus_Pass(t *testing.T) {
	srv := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusOK)
	}))
	defer srv.Close()
	if err := checks.HTTPStatus(srv.URL, http.StatusOK)(bg); err != nil {
		t.Fatalf("expected nil, got %v", err)
	}
}

func TestHTTPStatus_Fail_WrongStatus(t *testing.T) {
	srv := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusNotFound)
	}))
	defer srv.Close()
	if err := checks.HTTPStatus(srv.URL, http.StatusOK)(bg); err == nil {
		t.Fatal("expected error for wrong status")
	}
}

// HTTPBodyContains

func TestHTTPBodyContains_Pass(t *testing.T) {
	srv := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		fmt.Fprint(w, "hello world")
	}))
	defer srv.Close()
	if err := checks.HTTPBodyContains(srv.URL, "world")(bg); err != nil {
		t.Fatalf("expected nil, got %v", err)
	}
}

func TestHTTPBodyContains_Fail(t *testing.T) {
	srv := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		fmt.Fprint(w, "hello world")
	}))
	defer srv.Close()
	if err := checks.HTTPBodyContains(srv.URL, "nope")(bg); err == nil {
		t.Fatal("expected error when text not in body")
	}
}

// FileMatchesRegex uses pre-compiled regex (called multiple times without recompile)

func TestFileMatchesRegex_MultipleFiles(t *testing.T) {
	dir := t.TempDir()
	f1 := filepath.Join(dir, "a.txt")
	f2 := filepath.Join(dir, "b.txt")
	os.WriteFile(f1, []byte("abc123"), 0o644)
	os.WriteFile(f2, []byte("xyz999"), 0o644)

	check1 := checks.FileMatchesRegex(f1, `\d+`)
	check2 := checks.FileMatchesRegex(f2, `\d+`)

	if err := check1(bg); err != nil {
		t.Errorf("f1: expected nil, got %v", err)
	}
	if err := check2(bg); err != nil {
		t.Errorf("f2: expected nil, got %v", err)
	}
}

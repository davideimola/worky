// Package checks provides pre-built check functions for use with worky workshops.
// Each function returns a func() error suitable for the worky.Check.Run field.
package checks

import (
	"fmt"
	"io"
	"net"
	"net/http"
	"os"
	"os/exec"
	"regexp"
	"strings"
	"time"
)

// FileExists returns a check that passes if the file at path exists.
func FileExists(path string) func() error {
	return func() error {
		if _, err := os.Stat(path); os.IsNotExist(err) {
			return fmt.Errorf("file %q does not exist", path)
		}
		return nil
	}
}

// DirExists returns a check that passes if the directory at path exists.
func DirExists(path string) func() error {
	return func() error {
		info, err := os.Stat(path)
		if os.IsNotExist(err) {
			return fmt.Errorf("directory %q does not exist", path)
		}
		if err != nil {
			return fmt.Errorf("cannot stat %q: %w", path, err)
		}
		if !info.IsDir() {
			return fmt.Errorf("%q exists but is not a directory", path)
		}
		return nil
	}
}

// FileContains returns a check that passes if the file at path contains text.
func FileContains(path, text string) func() error {
	return func() error {
		data, err := os.ReadFile(path)
		if err != nil {
			return fmt.Errorf("cannot read %q: %w", path, err)
		}
		if !strings.Contains(string(data), text) {
			return fmt.Errorf("%q does not contain %q", path, text)
		}
		return nil
	}
}

// FileMatchesRegex returns a check that passes if the file at path matches pattern.
func FileMatchesRegex(path, pattern string) func() error {
	re := regexp.MustCompile(pattern)
	return func() error {
		data, err := os.ReadFile(path)
		if err != nil {
			return fmt.Errorf("cannot read %q: %w", path, err)
		}
		if !re.Match(data) {
			return fmt.Errorf("%q does not match pattern %q", path, pattern)
		}
		return nil
	}
}

// EnvVarSet returns a check that passes if the environment variable name is set and non-empty.
func EnvVarSet(name string) func() error {
	return func() error {
		if os.Getenv(name) == "" {
			return fmt.Errorf("environment variable %q is not set", name)
		}
		return nil
	}
}

// EnvVarEquals returns a check that passes if the environment variable name equals value.
func EnvVarEquals(name, value string) func() error {
	return func() error {
		got := os.Getenv(name)
		if got != value {
			return fmt.Errorf("environment variable %q = %q, want %q", name, got, value)
		}
		return nil
	}
}

// CommandSucceeds returns a check that passes if the command exits with code 0.
func CommandSucceeds(name string, args ...string) func() error {
	return func() error {
		cmd := exec.Command(name, args...)
		if out, err := cmd.CombinedOutput(); err != nil {
			msg := strings.TrimSpace(string(out))
			if msg == "" {
				return fmt.Errorf("command %q failed: %w", name, err)
			}
			return fmt.Errorf("command %q failed: %s", name, msg)
		}
		return nil
	}
}

// CommandOutputContains returns a check that passes if the command output contains text.
func CommandOutputContains(text, name string, args ...string) func() error {
	return func() error {
		cmd := exec.Command(name, args...)
		out, err := cmd.CombinedOutput()
		if err != nil {
			msg := strings.TrimSpace(string(out))
			if msg == "" {
				return fmt.Errorf("command %q failed: %w", name, err)
			}
			return fmt.Errorf("command %q failed: %s", name, msg)
		}
		if !strings.Contains(string(out), text) {
			return fmt.Errorf("output of %q does not contain %q", name, text)
		}
		return nil
	}
}

// PortOpen returns a check that passes if a TCP connection to host:port succeeds within 3 seconds.
func PortOpen(host string, port int) func() error {
	return func() error {
		addr := net.JoinHostPort(host, fmt.Sprintf("%d", port))
		conn, err := net.DialTimeout("tcp", addr, 3*time.Second)
		if err != nil {
			return fmt.Errorf("port %d on %q is not open: %w", port, host, err)
		}
		conn.Close()
		return nil
	}
}

// HTTPStatus returns a check that passes if a GET request to url returns expectedStatus.
func HTTPStatus(url string, expectedStatus int) func() error {
	return func() error {
		resp, err := http.Get(url) //nolint:gosec
		if err != nil {
			return fmt.Errorf("GET %q failed: %w", url, err)
		}
		resp.Body.Close()
		if resp.StatusCode != expectedStatus {
			return fmt.Errorf("GET %q returned %d, want %d", url, resp.StatusCode, expectedStatus)
		}
		return nil
	}
}

// HTTPBodyContains returns a check that passes if the body of a GET request to url contains text.
func HTTPBodyContains(url, text string) func() error {
	return func() error {
		resp, err := http.Get(url) //nolint:gosec
		if err != nil {
			return fmt.Errorf("GET %q failed: %w", url, err)
		}
		defer resp.Body.Close()
		body, err := io.ReadAll(resp.Body)
		if err != nil {
			return fmt.Errorf("reading body from %q failed: %w", url, err)
		}
		if !strings.Contains(string(body), text) {
			return fmt.Errorf("response body from %q does not contain %q", url, text)
		}
		return nil
	}
}

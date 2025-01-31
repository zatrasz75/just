package logger

import (
	"bytes"
	"errors"
	"io"
	"os"
	"os/exec"
	"strings"
	"testing"
)

func TestLogger_Info(t *testing.T) {
	var buf bytes.Buffer
	oldOutput := os.Stdout
	r, w, _ := os.Pipe()
	os.Stdout = w

	l, _ := NewLogger("")

	l.Info("Test Info Log")

	w.Close()
	os.Stdout = oldOutput

	io.Copy(&buf, r)
	output := buf.String()

	if !strings.Contains(output, "INFO") {
		t.Errorf("Expected log to contain '[INFO]', but got: %s", output)
	}
}

func TestLogger_Success(t *testing.T) {
	var buf bytes.Buffer
	oldOutput := os.Stdout
	r, w, _ := os.Pipe()
	os.Stdout = w

	l, _ := NewLogger("")

	l.Success("Test Info Log")

	w.Close()
	os.Stdout = oldOutput

	io.Copy(&buf, r)
	output := buf.String()

	if !strings.Contains(output, "SUCCESS") {
		t.Errorf("Expected log to contain '[SUCCESS]', but got: %s", output)
	}
}

func TestLogger_Trace(t *testing.T) {
	var buf bytes.Buffer
	oldOutput := os.Stdout
	r, w, _ := os.Pipe()
	os.Stdout = w

	l, _ := NewLogger("")

	l.Trace("Test Info Log")

	w.Close()
	os.Stdout = oldOutput

	io.Copy(&buf, r)
	output := buf.String()

	if !strings.Contains(output, "TRACE") {
		t.Errorf("Expected log to contain '[TRACE]', but got: %s", output)
	}
}

func TestLogger_Error(t *testing.T) {
	var buf bytes.Buffer
	oldOutput := os.Stdout
	r, w, _ := os.Pipe()
	os.Stdout = w

	l, _ := NewLogger("")

	l.Error("Test Info Log", errors.New("error error"))

	w.Close()
	os.Stdout = oldOutput

	io.Copy(&buf, r)
	output := buf.String()

	if !strings.Contains(output, "ERROR") {
		t.Errorf("Expected log to contain '[ERROR]', but got: %s", output)
	}
}

func TestLogger_Warn(t *testing.T) {
	var buf bytes.Buffer
	oldOutput := os.Stdout
	r, w, _ := os.Pipe()
	os.Stdout = w

	l, _ := NewLogger("")

	l.Warn("Test Info Log")

	w.Close()
	os.Stdout = oldOutput

	io.Copy(&buf, r)
	output := buf.String()

	if !strings.Contains(output, "WARN") {
		t.Errorf("Expected log to contain '[WARN]', but got: %s", output)
	}
}

func TestFatal(t *testing.T) {
	if os.Getenv("TEST_FATAL") == "1" {
		l, _ := NewLogger("")
		l.Fatal("Fatal message", errors.New("fatal error"))
		return
	}

	cmd := exec.Command(os.Args[0], "-test.run=TestFatal")
	cmd.Env = append(os.Environ(), "TEST_FATAL=1")
	err := cmd.Run()

	if exitError, ok := err.(*exec.ExitError); ok && exitError.ExitCode() == 1 {
		return // Это ожидаемый выход
	}

	t.Fatalf("Expected exit code 1, but got: %v", err)
}

func TestLogger_Debug(t *testing.T) {
	var buf bytes.Buffer
	oldOutput := os.Stdout
	r, w, _ := os.Pipe()
	os.Stdout = w

	l, _ := NewLogger("")

	l.Debug("Test Info Log")

	w.Close()
	os.Stdout = oldOutput

	io.Copy(&buf, r)
	output := buf.String()

	if !strings.Contains(output, "DEBUG") {
		t.Errorf("Expected log to contain '[DEBUG]', but got: %s", output)
	}
}

func TestLogger_Critical(t *testing.T) {
	var buf bytes.Buffer
	oldOutput := os.Stdout
	r, w, _ := os.Pipe()
	os.Stdout = w

	l, _ := NewLogger("")

	l.Critical("Test Info Log", errors.New("critical error"))

	w.Close()
	os.Stdout = oldOutput

	io.Copy(&buf, r)
	output := buf.String()

	if !strings.Contains(output, "CRITICAL") {
		t.Errorf("Expected log to contain '[CRITICAL]', but got: %s", output)
	}
}

func TestLogger_Panic(t *testing.T) {
	l, _ := NewLogger("")
	defer func() {
		if r := recover(); r == nil {
			t.Errorf("Expected panic, but function completed normally")
		}
	}()
	l.Panic("Test Panic Log", errors.New("panic error"))
}

func TestLogger_Security(t *testing.T) {
	var buf bytes.Buffer
	oldOutput := os.Stdout
	r, w, _ := os.Pipe()
	os.Stdout = w

	l, _ := NewLogger("")

	l.Security("Test Info Log")

	w.Close()
	os.Stdout = oldOutput

	io.Copy(&buf, r)
	output := buf.String()

	if !strings.Contains(output, "SECURITY") {
		t.Errorf("Expected log to contain '[SECURITY]', but got: %s", output)
	}
}

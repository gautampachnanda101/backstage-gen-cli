package output

import (
	"bytes"
	"testing"
)

func TestSetVerbosity(t *testing.T) {
	// Reset after test
	defer SetVerbosity(LevelNormal)

	tests := []struct {
		input    int
		expected int
	}{
		{LevelQuiet, LevelQuiet},
		{LevelNormal, LevelNormal},
		{LevelVerbose, LevelVerbose},
		{LevelDebug, LevelDebug},
		{-1, LevelQuiet}, // Below minimum
		{10, LevelDebug}, // Above maximum
	}

	for _, tt := range tests {
		SetVerbosity(tt.input)
		if GetVerbosity() != tt.expected {
			t.Errorf("SetVerbosity(%d): expected %d, got %d", tt.input, tt.expected, GetVerbosity())
		}
	}
}

func TestIsQuiet(t *testing.T) {
	defer SetVerbosity(LevelNormal)

	SetVerbosity(LevelQuiet)
	if !IsQuiet() {
		t.Error("expected IsQuiet to return true at level 0")
	}

	SetVerbosity(LevelNormal)
	if IsQuiet() {
		t.Error("expected IsQuiet to return false at level 1")
	}
}

func TestIsNormal(t *testing.T) {
	defer SetVerbosity(LevelNormal)

	SetVerbosity(LevelQuiet)
	if IsNormal() {
		t.Error("expected IsNormal to return false at level 0")
	}

	SetVerbosity(LevelNormal)
	if !IsNormal() {
		t.Error("expected IsNormal to return true at level 1")
	}

	SetVerbosity(LevelVerbose)
	if !IsNormal() {
		t.Error("expected IsNormal to return true at level 2")
	}
}

func TestIsVerbose(t *testing.T) {
	defer SetVerbosity(LevelNormal)

	SetVerbosity(LevelNormal)
	if IsVerbose() {
		t.Error("expected IsVerbose to return false at level 1")
	}

	SetVerbosity(LevelVerbose)
	if !IsVerbose() {
		t.Error("expected IsVerbose to return true at level 2")
	}

	SetVerbosity(LevelDebug)
	if !IsVerbose() {
		t.Error("expected IsVerbose to return true at level 3")
	}
}

func TestIsDebug(t *testing.T) {
	defer SetVerbosity(LevelNormal)

	SetVerbosity(LevelVerbose)
	if IsDebug() {
		t.Error("expected IsDebug to return false at level 2")
	}

	SetVerbosity(LevelDebug)
	if !IsDebug() {
		t.Error("expected IsDebug to return true at level 3")
	}
}

func TestPrintQuietMode(t *testing.T) {
	defer SetVerbosity(LevelNormal)

	// Capture output
	var buf bytes.Buffer
	oldStdout := stdout
	stdout = &buf
	defer func() { stdout = oldStdout }()

	SetVerbosity(LevelQuiet)
	Print("test message")

	if buf.Len() != 0 {
		t.Errorf("expected no output in quiet mode, got: %s", buf.String())
	}
}

func TestPrintNormalMode(t *testing.T) {
	defer SetVerbosity(LevelNormal)

	var buf bytes.Buffer
	oldStdout := stdout
	stdout = &buf
	defer func() { stdout = oldStdout }()

	SetVerbosity(LevelNormal)
	Print("test message")

	if buf.String() != "test message" {
		t.Errorf("expected 'test message', got: %s", buf.String())
	}
}

func TestPrintVerboseMode(t *testing.T) {
	defer SetVerbosity(LevelNormal)

	var buf bytes.Buffer
	oldStdout := stdout
	stdout = &buf
	defer func() { stdout = oldStdout }()

	// Test at normal level - should not print
	SetVerbosity(LevelNormal)
	PrintVerbose("verbose message")

	if buf.Len() != 0 {
		t.Errorf("expected no output at normal level, got: %s", buf.String())
	}

	// Test at verbose level - should print
	SetVerbosity(LevelVerbose)
	PrintVerbose("verbose message")

	if buf.String() != "verbose message" {
		t.Errorf("expected 'verbose message', got: %s", buf.String())
	}
}

func TestPrintDebugMode(t *testing.T) {
	defer SetVerbosity(LevelNormal)

	var buf bytes.Buffer
	oldStderr := stderr
	stderr = &buf
	defer func() { stderr = oldStderr }()

	// Test at verbose level - should not print debug
	SetVerbosity(LevelVerbose)
	PrintDebug("debug message\n")

	if buf.Len() != 0 {
		t.Errorf("expected no debug output at verbose level, got: %s", buf.String())
	}

	// Test at debug level - should print
	SetVerbosity(LevelDebug)
	PrintDebug("debug message\n")

	expected := "[DEBUG] debug message\n"
	if buf.String() != expected {
		t.Errorf("expected '%s', got: %s", expected, buf.String())
	}
}

func TestPrintError(t *testing.T) {
	defer SetVerbosity(LevelNormal)

	var buf bytes.Buffer
	oldStderr := stderr
	stderr = &buf
	defer func() { stderr = oldStderr }()

	// Error should print even in quiet mode
	SetVerbosity(LevelQuiet)
	PrintError("error message")

	if buf.String() != "error message" {
		t.Errorf("expected 'error message', got: %s", buf.String())
	}
}

func TestVerbosityLevelConstants(t *testing.T) {
	if LevelQuiet != 0 {
		t.Errorf("expected LevelQuiet to be 0, got %d", LevelQuiet)
	}
	if LevelNormal != 1 {
		t.Errorf("expected LevelNormal to be 1, got %d", LevelNormal)
	}
	if LevelVerbose != 2 {
		t.Errorf("expected LevelVerbose to be 2, got %d", LevelVerbose)
	}
	if LevelDebug != 3 {
		t.Errorf("expected LevelDebug to be 3, got %d", LevelDebug)
	}
}

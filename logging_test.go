package lansenger

import (
	"testing"
)

func TestDebugLoggerDefaultNil(t *testing.T) {
	if DebugLogger != nil {
		t.Error("expected DebugLogger to be nil by default")
	}
}

func TestDebugLoggerWorksWhenSet(t *testing.T) {
	called := false
	DebugLogger = func(format string, args ...interface{}) {
		called = true
	}
	defer func() { DebugLogger = nil }()

	DebugLogger("test %s", "message")
	if !called {
		t.Error("expected DebugLogger to be called")
	}
}

func TestDebugLoggerNilSafe(t *testing.T) {
	DebugLogger = nil

	// should not panic when DebugLogger is nil
	defer func() {
		if r := recover(); r != nil {
			t.Errorf("DebugLogger nil call panicked: %v", r)
		}
	}()

	if DebugLogger != nil {
		DebugLogger("should not be called")
	}
}

package lansenger

import (
	"testing"
)

func TestErrorMessages(t *testing.T) {
	e := NewAuthError("auth failed")
	if e.Message != "auth failed" {
		t.Errorf("expected message 'auth failed', got %s", e.Message)
	}
	if e.Retryable {
		t.Error("expected AuthError not retryable")
	}
	if e.Error() != "auth failed" {
		t.Errorf("expected Error() 'auth failed', got %s", e.Error())
	}
}

func TestAPIErrorWithCode(t *testing.T) {
	e := NewAPIError("api error", 10001)
	if e.ErrCode != 10001 {
		t.Errorf("expected errCode=10001, got %d", e.ErrCode)
	}
	if !e.Retryable {
		t.Error("expected APIError with nonzero errCode to be retryable")
	}
	expected := "api error (errCode: 10001)"
	if e.Error() != expected {
		t.Errorf("expected Error() '%s', got %s", expected, e.Error())
	}
}

func TestNetworkError(t *testing.T) {
	e := NewNetworkError("connection refused")
	if !e.Retryable {
		t.Error("expected NetworkError to be retryable")
	}
	if e.Message != "connection refused" {
		t.Errorf("expected message 'connection refused', got %s", e.Message)
	}
}

func TestFileError(t *testing.T) {
	e := NewFileError("file not found")
	if e.Retryable {
		t.Error("expected FileError not retryable")
	}
	if e.Message != "file not found" {
		t.Errorf("expected message 'file not found', got %s", e.Message)
	}
}

func TestConfigError(t *testing.T) {
	e := NewConfigError("missing config")
	if e.Retryable {
		t.Error("expected ConfigError not retryable")
	}
	if e.Message != "missing config" {
		t.Errorf("expected message 'missing config', got %s", e.Message)
	}
}

func TestErrorHierarchy(t *testing.T) {
	var _ error = (*AuthError)(nil)
	var _ error = (*ConfigError)(nil)
	var _ error = (*APIError)(nil)
	var _ error = (*NetworkError)(nil)
	var _ error = (*FileError)(nil)
	var _ error = (*LansengerError)(nil)

	if _, ok := interface{}(&AuthError{}).(*LansengerError); !ok {
	}
}

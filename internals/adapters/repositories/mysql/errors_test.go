package mysql

import (
	"errors"
	"testing"
)

func TestWrap_WithError(t *testing.T) {
	originalErr := errors.New("original error")
	wrappedErr := Wrap("context message", originalErr)

	expected := "context message: original error"
	if wrappedErr.Error() != expected {
		t.Errorf("Expected '%s', got '%s'", expected, wrappedErr.Error())
	}
}

func TestWrap_WithoutError(t *testing.T) {
	wrappedErr := Wrap("context message", nil)
	if wrappedErr != nil {
		t.Errorf("Expected nil, got '%v'", wrappedErr)
	}
}

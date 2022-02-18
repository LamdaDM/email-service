package main

import (
	"testing"
)

func TestFmtEmail(t *testing.T) {
	actual := string(container.emailOpts.template)
	expected := "To: ...to\r\nSubject: ...subject\r\n\r\n...body"

	if actual != expected {
		t.Errorf("Expected:\n%s\n\nActual:\n%s\n", expected, actual)
	}
}

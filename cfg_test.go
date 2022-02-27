package main

import "testing"

var Cfg *Config

func init() {
	Cfg = LoadConfig(".cfg.template")
}

func TestConfigLoad(t *testing.T) {
	sect := Cfg.GetSection("TEST_SECTION")
	actual := sect.Get("TEST_VAR")
	expected := "test"

	if actual != expected {
		t.Fatalf("Expected: %s\n\tActual: %s", expected, actual)
	}
}

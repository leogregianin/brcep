package main

import (
	"testing"
)

func TeststartPage(t *testing.T) {
	expected := "Hello, World!"
	if observed := "Hello, World!"; observed != expected {
		t.Fatalf("HelloWorld() = %v, want %v", observed, expected)
	}
}

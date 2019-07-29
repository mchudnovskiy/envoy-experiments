package main

import (
	"testing"
)

func TestMakeGreeting(t *testing.T) {

	result := MakeGreeting("test")
	if result != "Hello, test!" {
		t.Error("Expected: Hello, test! got: ", result)
	}
}

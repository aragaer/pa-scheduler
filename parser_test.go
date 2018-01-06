package main

import (
	"testing"
)

func TestParse(t *testing.T) {
	b := []byte(`{"name": "tick1", "delay": 0, "repeat": 0, "what": {"message": "test"}}`)

	event, err := Parse(b)
	if err != nil {
		t.Fatalf("Failed to parse event: %s", err)
	} else if event == nil {
		t.Fatalf("Failed to parse event")
	} else if event.Name != "tick1" {
		t.Fatalf("Expected event \"tick1\", got \"%s\"", event.Name)
	}
}

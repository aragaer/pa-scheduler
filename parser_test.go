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

func TestParse_NameRequired(t *testing.T) {
	b := []byte(`{"delay": 0, "repeat": 0, "what": {"message": "test"}}`)
	expectedError := "\"name\" field is missing"

	event, err := Parse(b)

	if err == nil {
		t.Fatalf("Error expected")
	} else if err.Error() != expectedError {
		t.Fatalf("Expected error %s, got %s", expectedError, err.Error())
	} else if event != nil {
		t.Fatalf("Event should not be created")
	}
}

func TestParse_WhatRequired(t *testing.T) {
	b := []byte(`{"delay": 0, "repeat": 0, "name": "tick"}`)
	expectedError := "\"what\" field is missing"

	event, err := Parse(b)

	if err == nil {
		t.Fatalf("Error expected")
	} else if err.Error() != expectedError {
		t.Fatalf("Expected error %s, got %s", expectedError, err.Error())
	} else if event != nil {
		t.Fatalf("Event should not be created")
	}
}

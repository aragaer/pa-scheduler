package main

import (
	"bufio"
	"fmt"
	"os"
	"sync"
	"testing"
	"time"
)

type Actions []struct {
	waitBefore, waitAfter time.Duration
	command               string
}

type main_tc struct {
	name     string
	actions  Actions
	expected Expected
}

var MainTestCases = []main_tc{
	{"nothing",
		Actions{{1, 0, ``}},
		Expected{}},
	{"hi",
		Actions{{0, 1, `{"command": "add", "name": "test", "what": "hi"}`}},
		Expected{`"hi"`}},
	{"hi repeating",
		Actions{{0, 3, `{"command": "add", "repeat": 1, "name": "test", "what": "hi"}`}},
		Expected{`"hi"`, `"hi"`, `"hi"`}},
}

func (tc main_tc) Run(t *testing.T) {
	var wg sync.WaitGroup
	oldStdout := os.Stdout
	oldStdin := os.Stdin
	stdout, writeFile, err := os.Pipe()
	if err != nil {
		t.Fatalf("Pipe: %v", err)
	}
	readFile, stdin, err2 := os.Pipe()
	if err2 != nil {
		t.Fatalf("Pipe: %v", err)
	}

	os.Stdout = writeFile
	os.Stdin = readFile

	wg.Add(1)

	go func() {
		defer wg.Done()
		main()
	}()

	for _, action := range tc.actions {
		time.Sleep(action.waitBefore * time.Second)
		fmt.Fprintln(stdin, action.command)
		time.Sleep(action.waitAfter * time.Second)
	}

	stdin.Close()

	wg.Wait()

	readFile.Close()
	writeFile.Close()
	os.Stdout = oldStdout
	os.Stdin = oldStdin

	scanner := bufio.NewScanner(stdout)
	for _, expected := range tc.expected {
		if !scanner.Scan() {
			t.Errorf("TestCase \"%s\" failed:", tc.name)
			t.Errorf("Expected \"%s\" got EOF", expected)
			break
		}
		line := scanner.Text()
		if line != expected {
			t.Errorf("TestCase \"%s\" failed:", tc.name)
			t.Errorf("Expected \"%s\" got \"%s\"", expected, line)
		}
	}
	for scanner.Scan() {
		line := scanner.Text()
		t.Errorf("TestCase \"%s\" failed:", tc.name)
		t.Errorf("After expected EOF got \"%s\"", line)
	}
}

func TestMainFunc(t *testing.T) {
	for _, tc := range MainTestCases {
		tc.Run(t)
	}
}

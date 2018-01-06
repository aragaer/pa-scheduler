package main

import (
	"strings"
	"testing"
)

type event struct {
	delay, repeat int64
	name          string
}

func (e event) mk() *Event {
	return &Event{e.delay, e.repeat, e.name, []byte(`{"event": "tick"}`)}
}

type Events []event
type Expected []string

type tc struct {
	events   Events
	expected Expected
}

var testCases = map[string]tc{
	"one immediate": {
		Events{{0, 0, "tick"}},
		Expected{"tick", "", ""}},
	"one delayed": {
		Events{{1, 0, "tick"}},
		Expected{"", "tick", ""}},
	"one repeating": {
		Events{{0, 1, "tick"}},
		Expected{"tick", "tick", "tick"}},
	"two immediate": {
		Events{{0, 0, "tick"}, {0, 0, "tock"}},
		Expected{"tick tock", "", ""}},
	"two alternating": {
		Events{{0, 2, "tick"}, {1, 2, "tock"}},
		Expected{"tick", "tock", "tick", "tock"}},
	"two different freq": {
		Events{{0, 2, "tick"}, {0, 3, "tock"}},
		Expected{"tick tock", "", "tick", "tock", "tick", "", "tick tock"}},
	"insert before": {
		Events{{1, 0, "tock"}, {0, 0, "tick"}},
		Expected{"tick", "tock", ""}},
}

func TestScheduler(t *testing.T) {
	for name, tc := range testCases {
		scheduler := NewScheduler()
		for _, e := range tc.events {
			scheduler.Queue(e.mk())
		}

		for tick, eventsForTick := range tc.expected {
			actual := make(map[string]bool)
			for {
				event := scheduler.GetTriggeredEvent()
				if event == nil {
					break
				}
				actual[event.Name] = true
			}

			expected := make(map[string]bool)
			for _, e := range strings.Fields(eventsForTick) {
				expected[e] = true
				if !actual[e] {
					t.Errorf("Test case \"%s\" failed on tick %d", name, tick)
					t.Errorf("event \"%s\" is expected but didn't happen", e)
				}
			}

			for e := range actual {
				if !expected[e] {
					t.Errorf("Test case \"%s\" failed on tick %d", name, tick)
					t.Errorf("event \"%s\" happened but is not expected", e)
				}
			}

			scheduler.Tick(1)
		}
	}
}

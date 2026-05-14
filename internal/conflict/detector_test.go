package conflict_test

import (
	"testing"

	"github.com/user/cronaudit/internal/conflict"
	"github.com/user/cronaudit/internal/parser"
)

func makeResult(host string, entries []parser.Entry) parser.ParseResult {
	return parser.ParseResult{Host: host, Entries: entries}
}

func entry(min, hour, dom, mon, dow, cmd string) parser.Entry {
	return parser.Entry{
		Minute: min, Hour: hour,
		DayOfMonth: dom, Month: mon,
		DayOfWeek: dow, Command: cmd,
	}
}

func TestDetectConflicts_NoneWhenDifferentSchedules(t *testing.T) {
	results := []parser.ParseResult{
		makeResult("host1", []parser.Entry{entry("0", "1", "*", "*", "*", "/bin/foo")}),
		makeResult("host2", []parser.Entry{entry("0", "2", "*", "*", "*", "/bin/foo")}),
	}
	got := conflict.DetectConflicts(results)
	if len(got) != 0 {
		t.Errorf("expected no conflicts, got %d", len(got))
	}
}

func TestDetectConflicts_DuplicateScheduleAndCommand(t *testing.T) {
	e := entry("0", "1", "*", "*", "*", "/bin/backup")
	results := []parser.ParseResult{
		makeResult("host1", []parser.Entry{e}),
		makeResult("host2", []parser.Entry{e}),
	}
	got := conflict.DetectConflicts(results)
	if len(got) != 1 {
		t.Fatalf("expected 1 conflict, got %d", len(got))
	}
	if got[0].HostA != "host1" || got[0].HostB != "host2" {
		t.Errorf("unexpected hosts: %s, %s", got[0].HostA, got[0].HostB)
	}
}

func TestDetectConflicts_SameScheduleDifferentCommand(t *testing.T) {
	results := []parser.ParseResult{
		makeResult("host1", []parser.Entry{entry("*/5", "*", "*", "*", "*", "/bin/taskA")}),
		makeResult("host2", []parser.Entry{entry("*/5", "*", "*", "*", "*", "/bin/taskB")}),
	}
	got := conflict.DetectConflicts(results)
	if len(got) != 1 {
		t.Fatalf("expected 1 conflict, got %d", len(got))
	}
	if got[0].Reason == "" {
		t.Error("expected a reason string")
	}
}

func TestConflict_String(t *testing.T) {
	c := conflict.Conflict{
		HostA:  "alpha",
		EntryA: entry("0", "*", "*", "*", "*", "/bin/x"),
		HostB:  "beta",
		EntryB: entry("0", "*", "*", "*", "*", "/bin/x"),
		Reason: "duplicate schedule and command",
	}
	s := c.String()
	if s == "" {
		t.Error("expected non-empty string")
	}
}

package diff_test

import (
	"testing"

	"github.com/user/cronaudit/internal/diff"
	"github.com/user/cronaudit/internal/parser"
)

func makeResult(host string, entries []parser.Entry) *parser.ParseResult {
	return &parser.ParseResult{Host: host, Entries: entries}
}

func entry(schedule, command string) parser.Entry {
	return parser.Entry{Schedule: schedule, Command: command}
}

func TestDiff_NoChanges(t *testing.T) {
	before := makeResult("host1", []parser.Entry{entry("* * * * *", "backup.sh")})
	after := makeResult("host1", []parser.Entry{entry("* * * * *", "backup.sh")})
	changes := diff.Diff(before, after)
	if len(changes) != 0 {
		t.Fatalf("expected 0 changes, got %d", len(changes))
	}
}

func TestDiff_Added(t *testing.T) {
	before := makeResult("host1", []parser.Entry{})
	after := makeResult("host1", []parser.Entry{entry("0 * * * *", "newjob.sh")})
	changes := diff.Diff(before, after)
	if len(changes) != 1 {
		t.Fatalf("expected 1 change, got %d", len(changes))
	}
	if changes[0].Kind != diff.Added {
		t.Errorf("expected Added, got %s", changes[0].Kind)
	}
}

func TestDiff_Removed(t *testing.T) {
	before := makeResult("host1", []parser.Entry{entry("0 * * * *", "oldjob.sh")})
	after := makeResult("host1", []parser.Entry{})
	changes := diff.Diff(before, after)
	if len(changes) != 1 {
		t.Fatalf("expected 1 change, got %d", len(changes))
	}
	if changes[0].Kind != diff.Removed {
		t.Errorf("expected Removed, got %s", changes[0].Kind)
	}
}

func TestDiff_Modified(t *testing.T) {
	before := makeResult("host1", []parser.Entry{entry("0 * * * *", "job.sh")})
	after := makeResult("host1", []parser.Entry{entry("30 * * * *", "job.sh")})
	changes := diff.Diff(before, after)
	if len(changes) != 1 {
		t.Fatalf("expected 1 change, got %d", len(changes))
	}
	if changes[0].Kind != diff.Modified {
		t.Errorf("expected Modified, got %s", changes[0].Kind)
	}
	if changes[0].Before.Schedule != "0 * * * *" {
		t.Errorf("unexpected before schedule: %s", changes[0].Before.Schedule)
	}
}

func TestDiff_HostFromAfter(t *testing.T) {
	before := makeResult("", []parser.Entry{})
	after := makeResult("myhost", []parser.Entry{entry("* * * * *", "x.sh")})
	changes := diff.Diff(before, after)
	if len(changes) != 1 {
		t.Fatalf("expected 1 change, got %d", len(changes))
	}
	if changes[0].Host != "myhost" {
		t.Errorf("expected host myhost, got %s", changes[0].Host)
	}
}

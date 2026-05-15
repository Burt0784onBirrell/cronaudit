package diff_test

import (
	"testing"

	"github.com/user/cronaudit/internal/diff"
	"github.com/user/cronaudit/internal/parser"
)

func TestDiffHosts_NewHost(t *testing.T) {
	before := diff.HostSnapshot{}
	after := diff.HostSnapshot{
		"web1": makeResult("web1", []parser.Entry{entry("* * * * *", "ping.sh")}),
	}
	changes := diff.DiffHosts(before, after)
	if len(changes) != 1 || changes[0].Kind != diff.Added {
		t.Fatalf("expected 1 Added change, got %+v", changes)
	}
}

func TestDiffHosts_RemovedHost(t *testing.T) {
	before := diff.HostSnapshot{
		"web1": makeResult("web1", []parser.Entry{entry("* * * * *", "ping.sh")}),
	}
	after := diff.HostSnapshot{}
	changes := diff.DiffHosts(before, after)
	if len(changes) != 1 || changes[0].Kind != diff.Removed {
		t.Fatalf("expected 1 Removed change, got %+v", changes)
	}
}

func TestDiffHosts_MultipleHosts(t *testing.T) {
	before := diff.HostSnapshot{
		"h1": makeResult("h1", []parser.Entry{entry("0 1 * * *", "a.sh")}),
		"h2": makeResult("h2", []parser.Entry{entry("0 2 * * *", "b.sh")}),
	}
	after := diff.HostSnapshot{
		"h1": makeResult("h1", []parser.Entry{entry("0 3 * * *", "a.sh")}),
		"h2": makeResult("h2", []parser.Entry{entry("0 2 * * *", "b.sh")}),
	}
	changes := diff.DiffHosts(before, after)
	if len(changes) != 1 {
		t.Fatalf("expected 1 change (h1 modified), got %d: %+v", len(changes), changes)
	}
	if changes[0].Kind != diff.Modified || changes[0].Host != "h1" {
		t.Errorf("unexpected change: %+v", changes[0])
	}
}

func TestDiffHosts_NoChanges(t *testing.T) {
	snap := diff.HostSnapshot{
		"h1": makeResult("h1", []parser.Entry{entry("* * * * *", "job.sh")}),
	}
	changes := diff.DiffHosts(snap, snap)
	if len(changes) != 0 {
		t.Fatalf("expected 0 changes, got %d", len(changes))
	}
}

package conflict_test

import (
	"bytes"
	"strings"
	"testing"

	"github.com/user/cronaudit/internal/conflict"
	"github.com/user/cronaudit/internal/parser"
)

func TestReport_NoConflicts(t *testing.T) {
	r := conflict.NewReport(nil)
	if r.HasConflicts() {
		t.Error("expected no conflicts")
	}
	if r.Summary() != "no conflicts detected" {
		t.Errorf("unexpected summary: %s", r.Summary())
	}
}

func TestReport_WithConflicts(t *testing.T) {
	e := parser.Entry{Minute: "0", Hour: "1", DayOfMonth: "*", Month: "*", DayOfWeek: "*", Command: "/bin/x"}
	conflicts := []conflict.Conflict{
		{HostA: "h1", EntryA: e, HostB: "h2", EntryB: e, Reason: "duplicate schedule and command"},
	}
	r := conflict.NewReport(conflicts)
	if !r.HasConflicts() {
		t.Error("expected conflicts")
	}
	if r.Summary() != "1 conflict(s) detected" {
		t.Errorf("unexpected summary: %s", r.Summary())
	}
}

func TestReport_WriteTo(t *testing.T) {
	e := parser.Entry{Minute: "*/10", Hour: "*", DayOfMonth: "*", Month: "*", DayOfWeek: "*", Command: "/bin/y"}
	conflicts := []conflict.Conflict{
		{HostA: "alpha", EntryA: e, HostB: "beta", EntryB: e, Reason: "duplicate schedule and command"},
	}
	r := conflict.NewReport(conflicts)

	var buf bytes.Buffer
	n, err := r.WriteTo(&buf)
	if err != nil {
		t.Fatalf("WriteTo error: %v", err)
	}
	if n == 0 {
		t.Error("expected non-zero bytes written")
	}
	output := buf.String()
	if !strings.Contains(output, "1 conflict(s) detected") {
		t.Errorf("output missing summary: %s", output)
	}
	if !strings.Contains(output, "alpha") || !strings.Contains(output, "beta") {
		t.Errorf("output missing host names: %s", output)
	}
}

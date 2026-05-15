package output

import (
	"bytes"
	"strings"
	"testing"

	"github.com/cronaudit/internal/conflict"
	"github.com/cronaudit/internal/validator"
)

func TestBuildSummary_NoIssues(t *testing.T) {
	s := BuildSummary([]string{"host1"}, nil, nil)
	if s.Errors != 0 || s.Warnings != 0 || s.Conflicts != 0 || s.Overlaps != 0 {
		t.Errorf("expected zero counts, got %+v", s)
	}
}

func TestBuildSummary_WithErrors(t *testing.T) {
	errs := []validator.ValidationError{
		{Severity: validator.SeverityError, Message: "bad minute"},
		{Severity: validator.SeverityWarning, Message: "deprecated @reboot"},
		{Severity: validator.SeverityWarning, Message: "another warning"},
	}
	s := BuildSummary([]string{"host1", "host2"}, errs, nil)
	if s.Errors != 1 {
		t.Errorf("expected 1 error, got %d", s.Errors)
	}
	if s.Warnings != 2 {
		t.Errorf("expected 2 warnings, got %d", s.Warnings)
	}
	if len(s.HostsAudited) != 2 {
		t.Errorf("expected 2 hosts, got %d", len(s.HostsAudited))
	}
}

func TestBuildSummary_WithReport(t *testing.T) {
	report := &conflict.Report{
		Conflicts: make([]conflict.Conflict, 3),
		Overlaps:  make([]conflict.Overlap, 1),
	}
	s := BuildSummary(nil, nil, report)
	if s.Conflicts != 3 {
		t.Errorf("expected 3 conflicts, got %d", s.Conflicts)
	}
	if s.Overlaps != 1 {
		t.Errorf("expected 1 overlap, got %d", s.Overlaps)
	}
}

func TestSummary_WriteTo_OKStatus(t *testing.T) {
	s := Summary{
		TotalEntries: 10,
		ValidEntries: 10,
		HostsAudited: []string{"cron01"},
	}
	var buf bytes.Buffer
	s.WriteTo(&buf)
	out := buf.String()
	if !strings.Contains(out, "Status        : OK") {
		t.Errorf("expected OK status in output:\n%s", out)
	}
	if !strings.Contains(out, "cron01") {
		t.Errorf("expected host name in output:\n%s", out)
	}
}

func TestSummary_WriteTo_IssuesStatus(t *testing.T) {
	s := Summary{
		TotalEntries: 5,
		Errors:       2,
		Conflicts:    1,
	}
	var buf bytes.Buffer
	s.WriteTo(&buf)
	out := buf.String()
	if !strings.Contains(out, "ISSUES FOUND") {
		t.Errorf("expected ISSUES FOUND in output:\n%s", out)
	}
}

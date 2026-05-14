package output_test

import (
	"bytes"
	"strings"
	"testing"

	"github.com/cronaudit/internal/conflict"
	"github.com/cronaudit/internal/output"
	"github.com/cronaudit/internal/validator"
)

func TestFormatter_TextNoIssues(t *testing.T) {
	var buf bytes.Buffer
	f := output.NewFormatter(output.FormatText, &buf)
	err := f.Write(output.Summary{Host: "web01"})
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	out := buf.String()
	if !strings.Contains(out, "web01") {
		t.Error("expected host name in output")
	}
	if !strings.Contains(out, "No issues found") {
		t.Error("expected no-issues message")
	}
}

func TestFormatter_TextWithWarnings(t *testing.T) {
	var buf bytes.Buffer
	f := output.NewFormatter(output.FormatText, &buf)
	s := output.Summary{
		Host: "db01",
		Warnings: []validator.ValidationError{
			{Line: 3, Message: "invalid minute field"},
		},
	}
	if err := f.Write(s); err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	out := buf.String()
	if !strings.Contains(out, "[line 3]") {
		t.Errorf("expected line number in output, got: %s", out)
	}
	if !strings.Contains(out, "invalid minute field") {
		t.Errorf("expected warning message in output, got: %s", out)
	}
}

func TestFormatter_TextWithConflicts(t *testing.T) {
	var buf bytes.Buffer
	f := output.NewFormatter(output.FormatText, &buf)
	s := output.Summary{
		Host: "app01",
		Conflicts: []conflict.Conflict{
			{Schedule: "0 * * * *", Command: "/usr/bin/backup", Hosts: []string{"app01", "app02"}},
		},
	}
	if err := f.Write(s); err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	out := buf.String()
	if !strings.Contains(out, "/usr/bin/backup") {
		t.Errorf("expected command in output, got: %s", out)
	}
	if !strings.Contains(out, "app02") {
		t.Errorf("expected conflicting host in output, got: %s", out)
	}
}

func TestFormatter_DefaultIsText(t *testing.T) {
	var buf bytes.Buffer
	f := output.NewFormatter("", &buf)
	if err := f.Write(output.Summary{Host: "host1"}); err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if buf.Len() == 0 {
		t.Error("expected non-empty output for default format")
	}
}

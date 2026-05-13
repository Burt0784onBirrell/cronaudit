package parser

import (
	"testing"
)

func TestParseLine_ValidStandard(t *testing.T) {
	entry, err := ParseLine("host1", 1, "*/5 * * * * /usr/bin/backup.sh")
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if entry == nil {
		t.Fatal("expected entry, got nil")
	}
	if entry.Minute != "*/5" {
		t.Errorf("expected minute '*/5', got %q", entry.Minute)
	}
	if entry.Command != "/usr/bin/backup.sh" {
		t.Errorf("expected command '/usr/bin/backup.sh', got %q", entry.Command)
	}
	if entry.Host != "host1" {
		t.Errorf("expected host 'host1', got %q", entry.Host)
	}
}

func TestParseLine_Shorthand(t *testing.T) {
	entry, err := ParseLine("host2", 3, "@daily /usr/bin/cleanup")
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if !entry.IsShorthand() {
		t.Error("expected IsShorthand() to be true")
	}
	if entry.Minute != "@daily" {
		t.Errorf("expected '@daily', got %q", entry.Minute)
	}
}

func TestParseLine_Comment(t *testing.T) {
	entry, err := ParseLine("host1", 2, "# this is a comment")
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if entry != nil {
		t.Errorf("expected nil entry for comment, got %+v", entry)
	}
}

func TestParseLine_EmptyLine(t *testing.T) {
	entry, err := ParseLine("host1", 5, "   ")
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if entry != nil {
		t.Errorf("expected nil entry for empty line")
	}
}

func TestParseLine_TooFewFields(t *testing.T) {
	_, err := ParseLine("host1", 7, "* * * *")
	if err == nil {
		t.Error("expected error for too few fields, got nil")
	}
}

func TestParseLine_EnvVar(t *testing.T) {
	entry, err := ParseLine("host1", 1, "MAILTO=root")
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if entry != nil {
		t.Errorf("expected nil for env var line, got %+v", entry)
	}
}

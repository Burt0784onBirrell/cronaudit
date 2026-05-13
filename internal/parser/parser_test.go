package parser

import (
	"strings"
	"testing"
)

const sampleCrontab = `
# Daily backup
@daily /usr/bin/backup.sh
MAILTO=admin@example.com
*/10 * * * * /usr/bin/healthcheck
0 2 * * 1 /usr/bin/weekly-report
# bad line below
* * *
`

func TestParseCrontabString_Entries(t *testing.T) {
	result, err := ParseCrontabString("web01", sampleCrontab)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	if len(result.Entries) != 3 {
		t.Errorf("expected 3 entries, got %d", len(result.Entries))
	}

	if !result.HasErrors() {
		t.Error("expected parse errors, got none")
	}
}

func TestParseCrontabString_Host(t *testing.T) {
	result, _ := ParseCrontabString("db01", "*/5 * * * * /bin/check")
	for _, e := range result.Entries {
		if e.Host != "db01" {
			t.Errorf("expected host 'db01', got %q", e.Host)
		}
	}
}

func TestParseCrontab_ReaderError(t *testing.T) {
	// Use a valid reader, ensure no error returned from scanner
	r := strings.NewReader("0 * * * * /bin/true")
	result, err := ParseCrontab("host1", r)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if len(result.Entries) != 1 {
		t.Errorf("expected 1 entry, got %d", len(result.Entries))
	}
}

func TestParseResult_ErrorSummary(t *testing.T) {
	result, _ := ParseCrontabString("host1", "* * *\n* * * *")
	summary := result.ErrorSummary()
	if summary == "" {
		t.Error("expected non-empty error summary")
	}
}

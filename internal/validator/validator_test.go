package validator_test

import (
	"strings"
	"testing"

	"github.com/cronaudit/internal/parser"
	"github.com/cronaudit/internal/validator"
)

func parseResult(host, content string) *parser.ParseResult {
	r, _ := parser.ParseCrontabString(host, strings.NewReader(content))
	return r
}

func TestValidateEntries_Valid(t *testing.T) {
	result := parseResult("host1", "0 2 * * 1 /usr/bin/backup\n*/15 * * * * /usr/bin/check\n")
	errs := validator.ValidateEntries(result)
	if len(errs) != 0 {
		t.Errorf("expected no errors, got %d: %v", len(errs), errs)
	}
}

func TestValidateEntries_InvalidMinute(t *testing.T) {
	result := parseResult("host1", "60 2 * * 1 /usr/bin/backup\n")
	errs := validator.ValidateEntries(result)
	if len(errs) == 0 {
		t.Fatal("expected validation error for minute=60")
	}
	if errs[0].Field != "minute" {
		t.Errorf("expected field=minute, got %q", errs[0].Field)
	}
}

func TestValidateEntries_InvalidHour(t *testing.T) {
	result := parseResult("host2", "0 25 * * * /bin/run\n")
	errs := validator.ValidateEntries(result)
	if len(errs) == 0 {
		t.Fatal("expected validation error for hour=25")
	}
	if errs[0].Field != "hour" {
		t.Errorf("expected field=hour, got %q", errs[0].Field)
	}
}

func TestValidateEntries_ValidShorthand(t *testing.T) {
	result := parseResult("host1", "@daily /usr/bin/cleanup\n")
	errs := validator.ValidateEntries(result)
	if len(errs) != 0 {
		t.Errorf("expected no errors for @daily, got %v", errs)
	}
}

func TestValidateEntries_InvalidShorthand(t *testing.T) {
	result := parseResult("host1", "@whenever /usr/bin/cleanup\n")
	errs := validator.ValidateEntries(result)
	if len(errs) == 0 {
		t.Fatal("expected error for unknown shorthand @whenever")
	}
	if errs[0].Field != "schedule" {
		t.Errorf("expected field=schedule, got %q", errs[0].Field)
	}
}

func TestValidateEntries_SkipsComments(t *testing.T) {
	result := parseResult("host1", "# this is a comment\n\n0 1 * * * /bin/run\n")
	errs := validator.ValidateEntries(result)
	if len(errs) != 0 {
		t.Errorf("expected no errors, got %v", errs)
	}
}

func TestValidationError_String(t *testing.T) {
	e := validator.ValidationError{Host: "web01", Line: 5, Field: "month", Message: "value out of range"}
	got := e.Error()
	if !strings.Contains(got, "web01") || !strings.Contains(got, "month") {
		t.Errorf("unexpected error string: %q", got)
	}
}

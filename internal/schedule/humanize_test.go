package schedule

import (
	"testing"
)

func TestHumanize_Shorthands(t *testing.T) {
	cases := []struct {
		expr string
		want string
	}{
		{"@yearly", "once a year, at midnight on January 1st"},
		{"@annually", "once a year, at midnight on January 1st"},
		{"@monthly", "once a month, at midnight on the 1st"},
		{"@weekly", "once a week, at midnight on Sunday"},
		{"@daily", "once a day, at midnight"},
		{"@midnight", "once a day, at midnight"},
		{"@hourly", "once an hour, at the start of the hour"},
		{"@reboot", "once at startup"},
	}
	for _, tc := range cases {
		t.Run(tc.expr, func(t *testing.T) {
			got, err := Humanize(tc.expr)
			if err != nil {
				t.Fatalf("unexpected error: %v", err)
			}
			if got != tc.want {
				t.Errorf("Humanize(%q) = %q; want %q", tc.expr, got, tc.want)
			}
		})
	}
}

func TestHumanize_EveryMinute(t *testing.T) {
	got, err := Humanize("* * * * *")
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if got != "every minute" {
		t.Errorf("got %q, want \"every minute\"", got)
	}
}

func TestHumanize_HourlyAt30(t *testing.T) {
	got, err := Humanize("30 * * * *")
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if got != "at minute 30 of every hour" {
		t.Errorf("got %q", got)
	}
}

func TestHumanize_SpecificTime(t *testing.T) {
	got, err := Humanize("0 9 * * *")
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if got != "at 09:00" {
		t.Errorf("got %q, want \"at 09:00\"", got)
	}
}

func TestHumanize_StepMinutes(t *testing.T) {
	got, err := Humanize("*/15 * * * *")
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if got != "every 15 minutes" {
		t.Errorf("got %q, want \"every 15 minutes\"", got)
	}
}

func TestHumanize_WithDayOfWeek(t *testing.T) {
	got, err := Humanize("0 8 * * 1")
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if got != "at 08:00, on Monday" {
		t.Errorf("got %q, want \"at 08:00, on Monday\"", got)
	}
}

func TestHumanize_WithDomAndMonth(t *testing.T) {
	got, err := Humanize("0 0 1 1 *")
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if got != "at 00:00, on day 1 of month 1" {
		t.Errorf("got %q", got)
	}
}

func TestHumanize_InvalidFieldCount(t *testing.T) {
	_, err := Humanize("* * *")
	if err == nil {
		t.Fatal("expected error for wrong field count, got nil")
	}
}

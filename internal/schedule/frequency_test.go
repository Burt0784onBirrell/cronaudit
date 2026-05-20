package schedule

import (
	"testing"
)

func TestClassifyFrequency_Shorthands(t *testing.T) {
	tests := []struct {
		expr string
		want FrequencyClass
	}{
		{"@yearly", FreqRare},
		{"@annually", FreqRare},
		{"@monthly", FreqMonthly},
		{"@weekly", FreqWeekly},
		{"@daily", FreqDaily},
		{"@midnight", FreqDaily},
		{"@hourly", FreqHourly},
		{"@reboot", FreqRare},
	}
	for _, tc := range tests {
		t.Run(tc.expr, func(t *testing.T) {
			got, err := ClassifyFrequency(tc.expr)
			if err != nil {
				t.Fatalf("unexpected error: %v", err)
			}
			if got != tc.want {
				t.Errorf("ClassifyFrequency(%q) = %s, want %s", tc.expr, got, tc.want)
			}
		})
	}
}

func TestClassifyFrequency_Standard(t *testing.T) {
	tests := []struct {
		expr string
		want FrequencyClass
	}{
		{"* * * * *", FreqPerMinute},
		{"0 * * * *", FreqHourly},
		{"30 * * * *", FreqHourly},
		{"0 9 * * *", FreqDaily},
		{"0 9 * * 1", FreqWeekly},
		{"0 9 1 * *", FreqMonthly},
		{"0 9 1 1 *", FreqRare},
	}
	for _, tc := range tests {
		t.Run(tc.expr, func(t *testing.T) {
			got, err := ClassifyFrequency(tc.expr)
			if err != nil {
				t.Fatalf("unexpected error: %v", err)
			}
			if got != tc.want {
				t.Errorf("ClassifyFrequency(%q) = %s, want %s", tc.expr, got, tc.want)
			}
		})
	}
}

func TestClassifyFrequency_InvalidExpr(t *testing.T) {
	_, err := ClassifyFrequency("* * * *") // only 4 fields
	if err == nil {
		t.Error("expected error for 4-field expression, got nil")
	}
}

func TestFrequencyClass_String(t *testing.T) {
	tests := []struct {
		f    FrequencyClass
		want string
	}{
		{FreqPerMinute, "per-minute"},
		{FreqHourly, "hourly"},
		{FreqDaily, "daily"},
		{FreqWeekly, "weekly"},
		{FreqMonthly, "monthly"},
		{FreqRare, "rare"},
		{FreqUnknown, "unknown"},
	}
	for _, tc := range tests {
		if got := tc.f.String(); got != tc.want {
			t.Errorf("FrequencyClass(%d).String() = %q, want %q", tc.f, got, tc.want)
		}
	}
}

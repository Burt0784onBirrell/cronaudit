package schedule

import (
	"testing"
	"time"
)

func TestOccurrences_EveryMinute(t *testing.T) {
	w := DefaultWindow()
	times, err := Occurrences("* * * * *", w)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	// 24 hours * 60 minutes = 1440
	if len(times) != 1440 {
		t.Errorf("expected 1440 occurrences, got %d", len(times))
	}
}

func TestOccurrences_Hourly(t *testing.T) {
	w := DefaultWindow()
	times, err := Occurrences("0 * * * *", w)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if len(times) != 24 {
		t.Errorf("expected 24 occurrences, got %d", len(times))
	}
}

func TestOccurrences_InvalidExpr(t *testing.T) {
	w := DefaultWindow()
	_, err := Occurrences("99 * * * *", w)
	if err == nil {
		t.Error("expected error for invalid expression")
	}
}

func TestOccurrences_CutoffRespected(t *testing.T) {
	w := OverlapWindow{
		Start:    time.Date(2024, 1, 1, 0, 0, 0, 0, time.UTC),
		Duration: 1 * time.Hour,
		N:        1440,
	}
	times, err := Occurrences("* * * * *", w)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if len(times) != 60 {
		t.Errorf("expected 60 occurrences within 1 hour, got %d", len(times))
	}
}

func TestOverlaps_SameExpr(t *testing.T) {
	w := DefaultWindow()
	ok, err := Overlaps("0 * * * *", "0 * * * *", w)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if !ok {
		t.Error("expected identical expressions to overlap")
	}
}

func TestOverlaps_NoOverlap(t *testing.T) {
	w := DefaultWindow()
	// runs at minute 0; runs at minute 30 — they never coincide on the hour
	ok, err := Overlaps("0 * * * *", "30 * * * *", w)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if ok {
		t.Error("expected no overlap between minute-0 and minute-30 schedules")
	}
}

func TestOverlaps_PartialOverlap(t *testing.T) {
	w := DefaultWindow()
	// every 2 hours at :00 vs every 4 hours at :00 — they share some times
	ok, err := Overlaps("0 */2 * * *", "0 */4 * * *", w)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if !ok {
		t.Error("expected overlap between */2 and */4 hour schedules")
	}
}

func TestOverlaps_InvalidExprA(t *testing.T) {
	w := DefaultWindow()
	_, err := Overlaps("99 * * * *", "0 * * * *", w)
	if err == nil {
		t.Error("expected error for invalid expression A")
	}
}

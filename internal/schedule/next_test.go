package schedule_test

import (
	"testing"
	"time"

	"github.com/cronaudit/cronaudit/internal/schedule"
)

var epoch = time.Date(2024, 1, 1, 0, 0, 0, 0, time.UTC)

func TestNextN_EveryMinute(t *testing.T) {
	times, err := schedule.NextN("* * * * *", epoch, 3)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if len(times) != 3 {
		t.Fatalf("expected 3 times, got %d", len(times))
	}
	expected := epoch.Add(time.Minute)
	if !times[0].Equal(expected) {
		t.Errorf("first time: got %v, want %v", times[0], expected)
	}
}

func TestNextN_HourlyAt30(t *testing.T) {
	times, err := schedule.NextN("30 * * * *", epoch, 2)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if len(times) != 2 {
		t.Fatalf("expected 2 times, got %d", len(times))
	}
	// epoch is 2024-01-01 00:00; next :30 is 00:30
	want0 := time.Date(2024, 1, 1, 0, 30, 0, 0, time.UTC)
	if !times[0].Equal(want0) {
		t.Errorf("first: got %v, want %v", times[0], want0)
	}
	want1 := time.Date(2024, 1, 1, 1, 30, 0, 0, time.UTC)
	if !times[1].Equal(want1) {
		t.Errorf("second: got %v, want %v", times[1], want1)
	}
}

func TestNextN_Step(t *testing.T) {
	// Every 15 minutes
	times, err := schedule.NextN("*/15 * * * *", epoch, 4)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	expectedMinutes := []int{0, 15, 30, 45}
	for i, want := range expectedMinutes {
		if times[i].Minute() != want {
			t.Errorf("time[%d] minute: got %d, want %d", i, times[i].Minute(), want)
		}
	}
}

func TestNextN_InvalidExpr(t *testing.T) {
	_, err := schedule.NextN("bad expr", epoch, 1)
	if err == nil {
		t.Fatal("expected error for invalid expression")
	}
}

func TestNextN_SundayBoth0And7(t *testing.T) {
	// 7 should be treated as Sunday same as 0
	times7, err := schedule.NextN("0 12 * * 7", epoch, 1)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	times0, err := schedule.NextN("0 12 * * 0", epoch, 1)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if !times7[0].Equal(times0[0]) {
		t.Errorf("Sunday mismatch: dow=7 gave %v, dow=0 gave %v", times7[0], times0[0])
	}
}

func TestNextN_Range(t *testing.T) {
	// Minutes 10-12 of every hour
	times, err := schedule.NextN("10-12 0 * * *", epoch, 3)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	for i, want := range []int{10, 11, 12} {
		if times[i].Minute() != want {
			t.Errorf("time[%d] minute: got %d, want %d", i, times[i].Minute(), want)
		}
	}
}

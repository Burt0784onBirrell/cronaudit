// Package schedule provides utilities for working with cron schedule expressions.
package schedule

import (
	"time"
)

// OverlapWindow defines the time range to check for schedule overlaps.
type OverlapWindow struct {
	Start    time.Time
	Duration time.Duration
	N        int // number of occurrences to sample
}

// DefaultWindow returns a 24-hour overlap window starting from a fixed epoch,
// sampling up to 1440 occurrences (one per minute).
func DefaultWindow() OverlapWindow {
	return OverlapWindow{
		Start:    time.Date(2024, 1, 1, 0, 0, 0, 0, time.UTC),
		Duration: 24 * time.Hour,
		N:        1440,
	}
}

// Occurrences returns all trigger times for the given cron expression within
// the window. Returns an error if the expression is invalid.
func Occurrences(expr string, w OverlapWindow) ([]time.Time, error) {
	times, err := NextN(expr, w.N, w.Start)
	if err != nil {
		return nil, err
	}

	cutoff := w.Start.Add(w.Duration)
	var result []time.Time
	for _, t := range times {
		if t.Before(cutoff) {
			result = append(result, t)
		}
	}
	return result, nil
}

// Overlaps returns true if two cron expressions share at least one trigger
// time within the given window.
func Overlaps(exprA, exprB string, w OverlapWindow) (bool, error) {
	timesA, err := Occurrences(exprA, w)
	if err != nil {
		return false, err
	}
	timesB, err := Occurrences(exprB, w)
	if err != nil {
		return false, err
	}

	setA := make(map[time.Time]struct{}, len(timesA))
	for _, t := range timesA {
		setA[t] = struct{}{}
	}
	for _, t := range timesB {
		if _, ok := setA[t]; ok {
			return true, nil
		}
	}
	return false, nil
}

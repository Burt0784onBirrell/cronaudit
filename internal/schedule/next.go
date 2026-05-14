package schedule

import (
	"fmt"
	"strconv"
	"strings"
	"time"
)

// NextN returns the next n execution times for a cron expression starting from t.
func NextN(expr string, from time.Time, n int) ([]time.Time, error) {
	fields, err := parseExpr(expr)
	if err != nil {
		return nil, err
	}

	var results []time.Time
	t := from.Truncate(time.Minute).Add(time.Minute)

	for len(results) < n {
		if matchesAll(t, fields) {
			results = append(results, t)
		}
		t = t.Add(time.Minute)
		// Safety: don't search more than 4 years ahead
		if t.After(from.Add(4 * 365 * 24 * time.Hour)) {
			break
		}
	}
	return results, nil
}

type cronFields struct {
	minute, hour, dom, month, dow map[int]bool
}

func matchesAll(t time.Time, f cronFields) bool {
	return f.minute[t.Minute()] &&
		f.hour[t.Hour()] &&
		f.dom[t.Day()] &&
		f.month[int(t.Month())] &&
		f.dow[int(t.Weekday())]
}

func parseExpr(expr string) (cronFields, error) {
	parts := strings.Fields(expr)
	if len(parts) < 5 {
		return cronFields{}, fmt.Errorf("schedule: expected 5 fields, got %d", len(parts))
	}
	minute, err := expandField(parts[0], 0, 59)
	if err != nil {
		return cronFields{}, fmt.Errorf("schedule: minute: %w", err)
	}
	hour, err := expandField(parts[1], 0, 23)
	if err != nil {
		return cronFields{}, fmt.Errorf("schedule: hour: %w", err)
	}
	dom, err := expandField(parts[2], 1, 31)
	if err != nil {
		return cronFields{}, fmt.Errorf("schedule: dom: %w", err)
	}
	month, err := expandField(parts[3], 1, 12)
	if err != nil {
		return cronFields{}, fmt.Errorf("schedule: month: %w", err)
	}
	dow, err := expandField(parts[4], 0, 7)
	if err != nil {
		return cronFields{}, fmt.Errorf("schedule: dow: %w", err)
	}
	// Treat Sunday as both 0 and 7
	if dow[7] {
		dow[0] = true
	}
	return cronFields{minute, hour, dom, month, dow}, nil
}

func expandField(field string, min, max int) (map[int]bool, error) {
	result := make(map[int]bool)
	for _, part := range strings.Split(field, ",") {
		if err := expandPart(part, min, max, result); err != nil {
			return nil, err
		}
	}
	return result, nil
}

func expandPart(part string, min, max int, out map[int]bool) error {
	step := 1
	if idx := strings.Index(part, "/"); idx != -1 {
		var err error
		step, err = strconv.Atoi(part[idx+1:])
		if err != nil || step < 1 {
			return fmt.Errorf("invalid step %q", part[idx+1:])
		}
		part = part[:idx]
	}
	if part == "*" {
		for i := min; i <= max; i += step {
			out[i] = true
		}
		return nil
	}
	if idx := strings.Index(part, "-"); idx != -1 {
		lo, err1 := strconv.Atoi(part[:idx])
		hi, err2 := strconv.Atoi(part[idx+1:])
		if err1 != nil || err2 != nil {
			return fmt.Errorf("invalid range %q", part)
		}
		for i := lo; i <= hi; i += step {
			out[i] = true
		}
		return nil
	}
	v, err := strconv.Atoi(part)
	if err != nil {
		return fmt.Errorf("invalid value %q", part)
	}
	out[v] = true
	return nil
}

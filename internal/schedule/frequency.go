package schedule

import (
	"fmt"
	"strings"
	"time"
)

// FrequencyClass describes how often a cron expression fires.
type FrequencyClass int

const (
	FreqUnknown   FrequencyClass = iota
	FreqPerMinute                // fires every minute
	FreqHourly                   // fires once or a few times per hour
	FreqDaily                    // fires once or a few times per day
	FreqWeekly                   // fires on specific days of the week
	FreqMonthly                  // fires on specific days of the month
	FreqRare                     // fires less than once a week
)

func (f FrequencyClass) String() string {
	switch f {
	case FreqPerMinute:
		return "per-minute"
	case FreqHourly:
		return "hourly"
	case FreqDaily:
		return "daily"
	case FreqWeekly:
		return "weekly"
	case FreqMonthly:
		return "monthly"
	case FreqRare:
		return "rare"
	default:
		return "unknown"
	}
}

// ClassifyFrequency returns the FrequencyClass for a 5-field cron expression.
// expr must be a standard "min hour dom month dow" string or a @shorthand.
func ClassifyFrequency(expr string) (FrequencyClass, error) {
	expr = strings.TrimSpace(expr)

	// Handle shorthands
	switch expr {
	case "@yearly", "@annually":
		return FreqRare, nil
	case "@monthly":
		return FreqMonthly, nil
	case "@weekly":
		return FreqWeekly, nil
	case "@daily", "@midnight":
		return FreqDaily, nil
	case "@hourly":
		return FreqHourly, nil
	case "@reboot":
		return FreqRare, nil
	}

	fields := strings.Fields(expr)
	if len(fields) != 5 {
		return FreqUnknown, fmt.Errorf("expected 5 fields, got %d", len(fields))
	}

	minute, hour, dom, month, dow :=
		fields[0], fields[1], fields[2], fields[3], fields[4]

	// Estimate occurrences over a fixed window to classify.
	start := time.Date(2024, 1, 1, 0, 0, 0, 0, time.UTC)
	occurrences, err := Occurrences(expr, start, DefaultWindow)
	if err != nil {
		return FreqUnknown, err
	}

	const (
		minutesInWeek  = 7 * 24 * 60
		minutesInDay   = 24 * 60
		minutesInHour  = 60
		windowMinutes  = int(DefaultWindow.Minutes())
	)

	_ = minute
	_ = hour
	_ = dom
	_ = month
	_ = dow

	switch {
	case occurrences >= windowMinutes:
		return FreqPerMinute, nil
	case occurrences >= windowMinutes/minutesInHour:
		return FreqHourly, nil
	case occurrences >= windowMinutes/minutesInDay:
		return FreqDaily, nil
	case occurrences >= windowMinutes/minutesInWeek:
		return FreqWeekly, nil
	case occurrences > 0:
		return FreqMonthly, nil
	default:
		return FreqRare, nil
	}
}

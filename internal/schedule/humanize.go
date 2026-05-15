// Package schedule provides utilities for working with cron expressions.
package schedule

import (
	"fmt"
	"strings"
)

// Humanize returns a human-readable description of a cron expression.
// It handles both standard 5-field expressions and shorthand notations.
func Humanize(expr string) (string, error) {
	switch strings.ToLower(expr) {
	case "@yearly", "@annually":
		return "once a year, at midnight on January 1st", nil
	case "@monthly":
		return "once a month, at midnight on the 1st", nil
	case "@weekly":
		return "once a week, at midnight on Sunday", nil
	case "@daily", "@midnight":
		return "once a day, at midnight", nil
	case "@hourly":
		return "once an hour, at the start of the hour", nil
	case "@reboot":
		return "once at startup", nil
	}

	parts := strings.Fields(expr)
	if len(parts) != 5 {
		return "", fmt.Errorf("expected 5 fields, got %d", len(parts))
	}

	minute, hour, dom, month, dow := parts[0], parts[1], parts[2], parts[3], parts[4]

	var sb strings.Builder

	// Time description
	if minute == "*" && hour == "*" {
		sb.WriteString("every minute")
	} else if minute == "*" {
		sb.WriteString(fmt.Sprintf("every minute during hour %s", hour))
	} else if strings.HasPrefix(minute, "*/") {
		sb.WriteString(fmt.Sprintf("every %s minutes", minute[2:]))
		if hour != "*" {
			sb.WriteString(fmt.Sprintf(" during hour %s", hour))
		}
	} else if hour == "*" {
		sb.WriteString(fmt.Sprintf("at minute %s of every hour", minute))
	} else {
		sb.WriteString(fmt.Sprintf("at %s:%s", padTwo(hour), padTwo(minute)))
	}

	// Day/month description
	if dom != "*" && month != "*" {
		sb.WriteString(fmt.Sprintf(", on day %s of month %s", dom, month))
	} else if dom != "*" {
		sb.WriteString(fmt.Sprintf(", on day %s of every month", dom))
	} else if month != "*" {
		sb.WriteString(fmt.Sprintf(", every day in month %s", month))
	}

	// Day of week description
	if dow != "*" {
		dowName := expandDow(dow)
		sb.WriteString(fmt.Sprintf(", on %s", dowName))
	}

	return sb.String(), nil
}

func padTwo(s string) string {
	if len(s) == 1 {
		return "0" + s
	}
	return s
}

func expandDow(dow string) string {
	names := map[string]string{
		"0": "Sunday", "1": "Monday", "2": "Tuesday",
		"3": "Wednesday", "4": "Thursday", "5": "Friday",
		"6": "Saturday", "7": "Sunday",
	}
	if name, ok := names[dow]; ok {
		return name
	}
	return fmt.Sprintf("day-of-week %s", dow)
}

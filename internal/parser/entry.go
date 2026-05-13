package parser

import (
	"fmt"
	"strings"
)

// CronEntry represents a single parsed crontab entry.
type CronEntry struct {
	Host    string
	Line    int
	Raw     string
	Minute  string
	Hour    string
	Dom     string
	Month   string
	Dow     string
	Command string
}

// ParseLine parses a single crontab line into a CronEntry.
// It returns an error if the line is not a valid cron expression.
func ParseLine(host string, lineNum int, raw string) (*CronEntry, error) {
	line := strings.TrimSpace(raw)

	// Skip empty lines and comments
	if line == "" || strings.HasPrefix(line, "#") {
		return nil, nil
	}

	// Skip environment variable assignments
	if strings.Contains(line, "=") && !strings.HasPrefix(line, "@") {
		parts := strings.SplitN(line, "=", 2)
		if !strings.Contains(parts[0], " ") {
			return nil, nil
		}
	}

	fields := strings.Fields(line)

	// Handle @reboot, @daily, etc.
	if strings.HasPrefix(fields[0], "@") {
		if len(fields) < 2 {
			return nil, fmt.Errorf("line %d: shorthand entry missing command", lineNum)
		}
		return &CronEntry{
			Host:    host,
			Line:    lineNum,
			Raw:     raw,
			Minute:  fields[0],
			Command: strings.Join(fields[1:], " "),
		}, nil
	}

	if len(fields) < 6 {
		return nil, fmt.Errorf("line %d: expected at least 6 fields, got %d", lineNum, len(fields))
	}

	return &CronEntry{
		Host:    host,
		Line:    lineNum,
		Raw:     raw,
		Minute:  fields[0],
		Hour:    fields[1],
		Dom:     fields[2],
		Month:   fields[3],
		Dow:     fields[4],
		Command: strings.Join(fields[5:], " "),
	}, nil
}

// IsShorthand returns true if the entry uses a @keyword syntax.
func (e *CronEntry) IsShorthand() bool {
	return strings.HasPrefix(e.Minute, "@")
}

// String returns a human-readable representation of the entry.
func (e *CronEntry) String() string {
	return fmt.Sprintf("%s:%d: %s", e.Host, e.Line, e.Raw)
}

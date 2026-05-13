package parser

import (
	"bufio"
	"fmt"
	"io"
	"strings"
)

// ParseResult holds all entries and errors from parsing a crontab source.
type ParseResult struct {
	Host    string
	Entries []*CronEntry
	Errors  []error
}

// ParseCrontab reads crontab content from r and parses all valid entries.
// host is used to tag each entry with its origin.
func ParseCrontab(host string, r io.Reader) (*ParseResult, error) {
	result := &ParseResult{Host: host}

	scanner := bufio.NewScanner(r)
	lineNum := 0

	for scanner.Scan() {
		lineNum++
		raw := scanner.Text()

		entry, err := ParseLine(host, lineNum, raw)
		if err != nil {
			result.Errors = append(result.Errors, err)
			continue
		}
		if entry != nil {
			result.Entries = append(result.Entries, entry)
		}
	}

	if err := scanner.Err(); err != nil {
		return nil, fmt.Errorf("reading crontab for host %q: %w", host, err)
	}

	return result, nil
}

// ParseCrontabString is a convenience wrapper around ParseCrontab for string input.
func ParseCrontabString(host, content string) (*ParseResult, error) {
	return ParseCrontab(host, strings.NewReader(content))
}

// HasErrors returns true if any parse errors were encountered.
func (r *ParseResult) HasErrors() bool {
	return len(r.Errors) > 0
}

// ErrorSummary returns all parse errors joined as a single string.
func (r *ParseResult) ErrorSummary() string {
	msgs := make([]string, 0, len(r.Errors))
	for _, e := range r.Errors {
		msgs = append(msgs, e.Error())
	}
	return strings.Join(msgs, "; ")
}

// Package diff compares crontab snapshots across two points in time or two hosts,
// reporting added, removed, and modified entries.
package diff

import (
	"fmt"

	"github.com/user/cronaudit/internal/parser"
)

// ChangeKind describes the type of change detected.
type ChangeKind string

const (
	Added    ChangeKind = "added"
	Removed  ChangeKind = "removed"
	Modified ChangeKind = "modified"
)

// Change represents a single diffed crontab entry.
type Change struct {
	Kind    ChangeKind
	Host    string
	Before  *parser.Entry
	After   *parser.Entry
}

// String returns a human-readable description of the change.
func (c Change) String() string {
	switch c.Kind {
	case Added:
		return fmt.Sprintf("[%s] ADDED   %s %s", c.Host, c.After.Schedule, c.After.Command)
	case Removed:
		return fmt.Sprintf("[%s] REMOVED %s %s", c.Host, c.Before.Schedule, c.Before.Command)
	case Modified:
		return fmt.Sprintf("[%s] MODIFIED schedule %s -> %s for command %s",
			c.Host, c.Before.Schedule, c.After.Schedule, c.After.Command)
	}
	return ""
}

// Diff compares two ParseResult snapshots (before vs after) for the same host
// and returns the list of changes.
func Diff(before, after *parser.ParseResult) []Change {
	host := after.Host
	if host == "" {
		host = before.Host
	}

	beforeMap := indexByCommand(before.Entries)
	afterMap := indexByCommand(after.Entries)

	var changes []Change

	for cmd, bEntry := range beforeMap {
		aEntry, ok := afterMap[cmd]
		if !ok {
			changes = append(changes, Change{Kind: Removed, Host: host, Before: bEntry})
		} else if bEntry.Schedule != aEntry.Schedule {
			changes = append(changes, Change{Kind: Modified, Host: host, Before: bEntry, After: aEntry})
		}
	}

	for cmd, aEntry := range afterMap {
		if _, ok := beforeMap[cmd]; !ok {
			changes = append(changes, Change{Kind: Added, Host: host, After: aEntry})
		}
	}

	return changes
}

func indexByCommand(entries []parser.Entry) map[string]*parser.Entry {
	m := make(map[string]*parser.Entry, len(entries))
	for i := range entries {
		m[entries[i].Command] = &entries[i]
	}
	return m
}

package conflict

import (
	"fmt"
	"strings"

	"github.com/user/cronaudit/internal/parser"
)

// Conflict represents a scheduling conflict between two cron entries.
type Conflict struct {
	HostA  string
	EntryA parser.Entry
	HostB  string
	EntryB parser.Entry
	Reason string
}

func (c Conflict) String() string {
	return fmt.Sprintf("conflict between [%s]%q and [%s]%q: %s",
		c.HostA, c.EntryA.Command,
		c.HostB, c.EntryB.Command,
		c.Reason)
}

// DetectConflicts compares entries across multiple hosts and returns any
// scheduling conflicts found (identical schedule + command on different hosts).
func DetectConflicts(results []parser.ParseResult) []Conflict {
	var conflicts []Conflict

	for i := 0; i < len(results); i++ {
		for j := i + 1; j < len(results); j++ {
			for _, ea := range results[i].Entries {
				for _, eb := range results[j].Entries {
					if reason, ok := checkConflict(ea, eb); ok {
						conflicts = append(conflicts, Conflict{
							HostA:  results[i].Host,
							EntryA: ea,
							HostB:  results[j].Host,
							EntryB: eb,
							Reason: reason,
						})
					}
				}
			}
		}
	}

	return conflicts
}

// checkConflict returns a reason string and true if the two entries conflict.
func checkConflict(a, b parser.Entry) (string, bool) {
	scheduleA := scheduleKey(a)
	scheduleB := scheduleKey(b)

	if scheduleA == scheduleB && normalizeCmd(a.Command) == normalizeCmd(b.Command) {
		return "duplicate schedule and command", true
	}

	if scheduleA == scheduleB && normalizeCmd(a.Command) != normalizeCmd(b.Command) {
		return fmt.Sprintf("same schedule %q runs different commands", scheduleA), true
	}

	return "", false
}

func scheduleKey(e parser.Entry) string {
	if e.Shorthand != "" {
		return e.Shorthand
	}
	return strings.Join([]string{e.Minute, e.Hour, e.DayOfMonth, e.Month, e.DayOfWeek}, " ")
}

func normalizeCmd(cmd string) string {
	return strings.TrimSpace(cmd)
}

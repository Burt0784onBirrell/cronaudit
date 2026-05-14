package conflict

import (
	"fmt"

	"github.com/user/cronaudit/internal/parser"
	"github.com/user/cronaudit/internal/schedule"
)

// OverlapConflict describes two entries whose schedules fire at the same time.
type OverlapConflict struct {
	HostA  string
	EntryA parser.Entry
	HostB  string
	EntryB parser.Entry
}

// String returns a human-readable description of the overlap conflict.
func (o OverlapConflict) String() string {
	return fmt.Sprintf(
		"schedule overlap: [%s] %q and [%s] %q share trigger times",
		o.HostA, o.EntryA.Command,
		o.HostB, o.EntryB.Command,
	)
}

// DetectOverlaps finds pairs of entries across all hosts whose cron schedules
// overlap within a 24-hour window. Only entries with distinct commands are
// reported — exact duplicates are handled by DetectConflicts.
func DetectOverlaps(results []parser.ParseResult) ([]OverlapConflict, error) {
	type tagged struct {
		host  string
		entry parser.Entry
	}

	var all []tagged
	for _, r := range results {
		for _, e := range r.Entries {
			if e.Schedule == "" {
				continue
			}
			all = append(all, tagged{host: r.Host, entry: e})
		}
	}

	w := schedule.DefaultWindow()
	var conflicts []OverlapConflict

	for i := 0; i < len(all); i++ {
		for j := i + 1; j < len(all); j++ {
			a, b := all[i], all[j]
			if normalizeCmd(a.entry.Command) == normalizeCmd(b.entry.Command) {
				// exact duplicates reported elsewhere
				continue
			}
			ok, err := schedule.Overlaps(a.entry.Schedule, b.entry.Schedule, w)
			if err != nil {
				// skip entries with unparseable schedules
				continue
			}
			if ok {
				conflicts = append(conflicts, OverlapConflict{
					HostA:  a.host,
					EntryA: a.entry,
					HostB:  b.host,
					EntryB: b.entry,
				})
			}
		}
	}
	return conflicts, nil
}

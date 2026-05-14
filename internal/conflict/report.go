package conflict

import (
	"fmt"
	"io"
	"strings"
)

// Report holds the results of a conflict detection run.
type Report struct {
	Conflicts []Conflict
}

// HasConflicts returns true if any conflicts were found.
func (r *Report) HasConflicts() bool {
	return len(r.Conflicts) > 0
}

// Summary returns a short human-readable summary line.
func (r *Report) Summary() string {
	if !r.HasConflicts() {
		return "no conflicts detected"
	}
	return fmt.Sprintf("%d conflict(s) detected", len(r.Conflicts))
}

// WriteTo writes a formatted report to the given writer.
func (r *Report) WriteTo(w io.Writer) (int64, error) {
	var sb strings.Builder
	sb.WriteString(r.Summary())
	sb.WriteString("\n")

	for i, c := range r.Conflicts {
		sb.WriteString(fmt.Sprintf("  [%d] %s\n", i+1, c.String()))
	}

	n, err := fmt.Fprint(w, sb.String())
	return int64(n), err
}

// NewReport runs conflict detection and returns a populated Report.
func NewReport(conflicts []Conflict) *Report {
	return &Report{Conflicts: conflicts}
}

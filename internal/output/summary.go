package output

import (
	"fmt"
	"io"
	"strings"

	"github.com/cronaudit/internal/conflict"
	"github.com/cronaudit/internal/validator"
)

// Summary holds aggregated statistics for a cronaudit run.
type Summary struct {
	TotalEntries   int
	ValidEntries   int
	Warnings       int
	Errors         int
	Conflicts      int
	Overlaps       int
	HostsAudited   []string
}

// BuildSummary constructs a Summary from validation errors and conflict reports.
func BuildSummary(hosts []string, errs []validator.ValidationError, report *conflict.Report) Summary {
	s := Summary{
		HostsAudited: hosts,
	}

	for _, e := range errs {
		switch e.Severity {
		case validator.SeverityError:
			s.Errors++
		case validator.SeverityWarning:
			s.Warnings++
		}
	}

	if report != nil {
		s.Conflicts = len(report.Conflicts)
		s.Overlaps = len(report.Overlaps)
	}

	return s
}

// WriteTo writes a human-readable summary to w.
func (s Summary) WriteTo(w io.Writer) (int64, error) {
	var sb strings.Builder

	sb.WriteString("=== cronaudit summary ===\n")
	if len(s.HostsAudited) > 0 {
		sb.WriteString(fmt.Sprintf("Hosts audited : %s\n", strings.Join(s.HostsAudited, ", ")))
	}
	sb.WriteString(fmt.Sprintf("Total entries : %d\n", s.TotalEntries))
	sb.WriteString(fmt.Sprintf("Valid entries : %d\n", s.ValidEntries))
	sb.WriteString(fmt.Sprintf("Warnings      : %d\n", s.Warnings))
	sb.WriteString(fmt.Sprintf("Errors        : %d\n", s.Errors))
	sb.WriteString(fmt.Sprintf("Conflicts     : %d\n", s.Conflicts))
	sb.WriteString(fmt.Sprintf("Overlaps      : %d\n", s.Overlaps))

	if s.Errors == 0 && s.Conflicts == 0 {
		sb.WriteString("Status        : OK\n")
	} else {
		sb.WriteString("Status        : ISSUES FOUND\n")
	}

	n, err := fmt.Fprint(w, sb.String())
	return int64(n), err
}

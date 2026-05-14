package output

import (
	"fmt"
	"io"
	"strings"

	"github.com/cronaudit/internal/conflict"
	"github.com/cronaudit/internal/validator"
)

// Format controls the output format.
type Format string

const (
	FormatText Format = "text"
	FormatJSON  Format = "json"
)

// Summary holds the aggregated audit results for output.
type Summary struct {
	Host      string
	Warnings  []validator.ValidationError
	Conflicts []conflict.Conflict
}

// Formatter writes audit summaries to an io.Writer.
type Formatter struct {
	Format Format
	Writer io.Writer
}

// NewFormatter creates a Formatter with the given format and writer.
func NewFormatter(f Format, w io.Writer) *Formatter {
	return &Formatter{Format: f, Writer: w}
}

// Write outputs the summary in the configured format.
func (f *Formatter) Write(s Summary) error {
	switch f.Format {
	case FormatJSON:
		return f.writeJSON(s)
	default:
		return f.writeText(s)
	}
}

func (f *Formatter) writeText(s Summary) error {
	fmt.Fprintf(f.Writer, "Host: %s\n", s.Host)
	fmt.Fprintf(f.Writer, "%s\n", strings.Repeat("-", 40))

	if len(s.Warnings) == 0 && len(s.Conflicts) == 0 {
		fmt.Fprintln(f.Writer, "No issues found.")
		return nil
	}

	if len(s.Warnings) > 0 {
		fmt.Fprintf(f.Writer, "Validation warnings (%d):\n", len(s.Warnings))
		for _, w := range s.Warnings {
			fmt.Fprintf(f.Writer, "  [line %d] %s\n", w.Line, w.Message)
		}
	}

	if len(s.Conflicts) > 0 {
		fmt.Fprintf(f.Writer, "Conflicts (%d):\n", len(s.Conflicts))
		for _, c := range s.Conflicts {
			fmt.Fprintf(f.Writer, "  schedule=%q command=%q hosts=%s\n",
				c.Schedule, c.Command, strings.Join(c.Hosts, ", "))
		}
	}
	return nil
}

func (f *Formatter) writeJSON(s Summary) error {
	warnings := make([]map[string]interface{}, len(s.Warnings))
	for i, w := range s.Warnings {
		warnings[i] = map[string]interface{}{"line": w.Line, "message": w.Message}
	}
	conflicts := make([]map[string]interface{}, len(s.Conflicts))
	for i, c := range s.Conflicts {
		conflicts[i] = map[string]interface{}{
			"schedule": c.Schedule,
			"command":  c.Command,
			"hosts":    c.Hosts,
		}
	}
	fmt.Fprintf(f.Writer, `{"host":%q,"warnings":%s,"conflicts":%s}\n`,
		s.Host, marshalList(warnings), marshalList(conflicts))
	return nil
}

func marshalList(items []map[string]interface{}) string {
	if len(items) == 0 {
		return "[]"
	}
	parts := make([]string, len(items))
	for i, m := range items {
		pairs := make([]string, 0, len(m))
		for k, v := range m {
			switch val := v.(type) {
			case string:
				pairs = append(pairs, fmt.Sprintf("%q:%q", k, val))
			case int:
				pairs = append(pairs, fmt.Sprintf("%q:%d", k, val))
			case []string:
				quoted := make([]string, len(val))
				for j, s := range val {
					quoted[j] = fmt.Sprintf("%q", s)
				}
				pairs = append(pairs, fmt.Sprintf("%q:[%s]", k, strings.Join(quoted, ",")))
			}
		}
		parts[i] = "{" + strings.Join(pairs, ",") + "}"
	}
	return "[" + strings.Join(parts, ",") + "]"
}

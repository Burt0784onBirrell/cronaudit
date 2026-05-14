package validator

import (
	"fmt"
	"strconv"
	"strings"

	"github.com/cronaudit/internal/parser"
)

// ValidationError represents a problem found in a crontab entry.
type ValidationError struct {
	Host    string
	Line    int
	Field   string
	Message string
}

func (e ValidationError) Error() string {
	return fmt.Sprintf("%s line %d [%s]: %s", e.Host, e.Line, e.Field, e.Message)
}

// ValidateEntries validates all parsed entries and returns a list of errors.
func ValidateEntries(result *parser.ParseResult) []ValidationError {
	var errs []ValidationError
	for _, entry := range result.Entries {
		if entry.IsComment || entry.IsEmpty {
			continue
		}
		if entry.IsShorthand {
			if !isValidShorthand(entry.Schedule) {
				errs = append(errs, ValidationError{
					Host:    result.Host,
					Line:    entry.LineNumber,
					Field:   "schedule",
					Message: fmt.Sprintf("unknown shorthand %q", entry.Schedule),
				})
			}
			continue
		}
		fields := strings.Fields(entry.Schedule)
		if len(fields) != 5 {
			errs = append(errs, ValidationError{
				Host:    result.Host,
				Line:    entry.LineNumber,
				Field:   "schedule",
				Message: "expected 5 time fields",
			})
			continue
		}
		ranges := []struct {
			name     string
			min, max int
		}{
			{"minute", 0, 59},
			{"hour", 0, 23},
			{"day-of-month", 1, 31},
			{"month", 1, 12},
			{"day-of-week", 0, 7},
		}
		for i, r := range ranges {
			if err := validateField(fields[i], r.min, r.max); err != nil {
				errs = append(errs, ValidationError{
					Host:    result.Host,
					Line:    entry.LineNumber,
					Field:   r.name,
					Message: err.Error(),
				})
			}
		}
	}
	return errs
}

func isValidShorthand(s string) bool {
	valid := map[string]bool{
		"@yearly": true, "@annually": true, "@monthly": true,
		"@weekly": true, "@daily": true, "@midnight": true,
		"@hourly": true, "@reboot": true,
	}
	return valid[strings.ToLower(s)]
}

func validateField(field string, min, max int) error {
	if field == "*" {
		return nil
	}
	if strings.HasPrefix(field, "*/") {
		step, err := strconv.Atoi(strings.TrimPrefix(field, "*/"))
		if err != nil || step < 1 {
			return fmt.Errorf("invalid step value %q", field)
		}
		return nil
	}
	for _, part := range strings.Split(field, ",") {
		if strings.Contains(part, "-") {
			bounds := strings.SplitN(part, "-", 2)
			lo, err1 := strconv.Atoi(bounds[0])
			hi, err2 := strconv.Atoi(bounds[1])
			if err1 != nil || err2 != nil || lo > hi || lo < min || hi > max {
				return fmt.Errorf("invalid range %q (allowed %d-%d)", part, min, max)
			}
		} else {
			v, err := strconv.Atoi(part)
			if err != nil || v < min || v > max {
				return fmt.Errorf("value %q out of range (%d-%d)", part, min, max)
			}
		}
	}
	return nil
}

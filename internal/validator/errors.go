package validator

import "fmt"

// ValidationError represents a single validation issue on a crontab line.
type ValidationError struct {
	Host    string
	Line    int
	Field   string
	Message string
}

// Error implements the error interface.
func (e ValidationError) Error() string {
	if e.Host != "" {
		return fmt.Sprintf("[%s] line %d (%s): %s", e.Host, e.Line, e.Field, e.Message)
	}
	return fmt.Sprintf("line %d (%s): %s", e.Line, e.Field, e.Message)
}

// ValidationErrors is a slice of ValidationError with helper methods.
type ValidationErrors []ValidationError

// HasErrors returns true if there are any validation errors.
func (ve ValidationErrors) HasErrors() bool {
	return len(ve) > 0
}

// FilterByHost returns only errors matching the given host.
func (ve ValidationErrors) FilterByHost(host string) ValidationErrors {
	var out ValidationErrors
	for _, e := range ve {
		if e.Host == host {
			out = append(out, e)
		}
	}
	return out
}

// FilterByField returns only errors for the given crontab field name.
func (ve ValidationErrors) FilterByField(field string) ValidationErrors {
	var out ValidationErrors
	for _, e := range ve {
		if e.Field == field {
			out = append(out, e)
		}
	}
	return out
}

// Summary returns a human-readable count summary.
func (ve ValidationErrors) Summary() string {
	if len(ve) == 0 {
		return "no validation errors"
	}
	return fmt.Sprintf("%d validation error(s)", len(ve))
}

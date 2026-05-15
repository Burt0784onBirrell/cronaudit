package diff

import "github.com/user/cronaudit/internal/parser"

// HostSnapshot maps a host name to its parsed crontab result.
type HostSnapshot map[string]*parser.ParseResult

// DiffHosts compares two full snapshots (e.g. from two audit runs) and
// returns all changes across every host present in either snapshot.
func DiffHosts(before, after HostSnapshot) []Change {
	var all []Change

	seen := make(map[string]bool)

	for host, afterResult := range after {
		seen[host] = true
		beforeResult, ok := before[host]
		if !ok {
			// Entire host is new — every entry is Added.
			beforeResult = &parser.ParseResult{Host: host}
		}
		all = append(all, Diff(beforeResult, afterResult)...)
	}

	for host, beforeResult := range before {
		if seen[host] {
			continue
		}
		// Entire host was removed — every entry is Removed.
		afterResult := &parser.ParseResult{Host: host}
		all = append(all, Diff(beforeResult, afterResult)...)
	}

	return all
}

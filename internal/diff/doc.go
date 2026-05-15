// Package diff provides utilities for comparing crontab snapshots.
//
// A snapshot is a collection of parsed crontab results (parser.ParseResult),
// typically captured at different points in time or from different audit runs.
//
// # Basic usage
//
// Compare two results for a single host:
//
//	 changes := diff.Diff(before, after)
//	 for _, c := range changes {
//	     fmt.Println(c)
//	 }
//
// Compare full multi-host snapshots:
//
//	 changes := diff.DiffHosts(beforeSnapshot, afterSnapshot)
//
// # Change kinds
//
// Each Change has a Kind field of type ChangeKind:
//   - Added   — entry present in after but not before
//   - Removed — entry present in before but not after
//   - Modified — same command, different schedule
//
// Entries are keyed by their Command string; if two entries share a command
// but differ in schedule they are reported as Modified.
package diff

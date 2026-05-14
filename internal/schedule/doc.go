// Package schedule provides utilities for computing future execution times
// of cron expressions. It supports the standard 5-field cron syntax including
// wildcards (*), ranges (1-5), steps (*/2), and lists (1,3,5).
//
// Sunday may be specified as either 0 or 7 in the day-of-week field.
//
// Example usage:
//
//	times, err := schedule.NextN("0 9 * * 1-5", time.Now(), 5)
//	if err != nil {
//		log.Fatal(err)
//	}
//	for _, t := range times {
//		fmt.Println(t)
//	}
//
// This package does not handle shorthand expressions such as @daily or
// @hourly — those should be expanded by the parser before being passed here.
package schedule

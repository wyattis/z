package ztime

import "time"

// Returns true if the given date falls on a weekday
func IsWeekday(date time.Time) bool {
	switch date.Weekday() {
	case time.Monday:
	case time.Tuesday:
	case time.Wednesday:
	case time.Thursday:
	case time.Friday:
		return true
	}
	return false
}

// Returns true if the given date falls on a weekend
func IsWeekend(date time.Time) bool {
	switch date.Weekday() {
	case time.Monday:
	case time.Tuesday:
	case time.Wednesday:
	case time.Thursday:
	case time.Friday:
		return false
	}
	return true
}

package ztime

import (
	"fmt"
	"strconv"
	"time"
)

const (
	ISO8601 = "2006-01-02T15:04:05-0700"
)

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

var timeFormats = []string{
	time.Kitchen,
	"15:04",
	"15:04:05",
	"15:04:05 MST",
	"2006-01-02",
	time.RFC3339,
	time.RFC3339Nano,
	time.RFC1123,
	time.RFC1123Z,
	time.RFC850,
	time.RFC850,
	time.RFC822Z,
	time.RFC822,
	time.RubyDate,
	time.UnixDate,
	time.ANSIC,
}

// Attempt to parse a time using all known formats. If formats are passed in, we
// attempt to parse using those formats first. An error is thrown only if no
// matching formats are found.
func Parse(val string, formats ...string) (t time.Time, err error) {
	formats = append(formats, timeFormats...)
	for _, format := range formats {
		if t, err = time.Parse(format, val); err == nil {
			return
		}
	}
	err = fmt.Errorf("failed to parse the time %s. Attempted %d formats. Please provide a format.", val, len(formats))
	return
}

// The same as time.Unix, but the arguments are strings
func ParseUnix(seconds string, nanos string) (t time.Time, err error) {
	s, err := strconv.ParseInt(seconds, 10, 64)
	if err != nil {
		return
	}
	var n int64
	if nanos != "" {
		n, err = strconv.ParseInt(nanos, 10, 64)
		if err != nil {
			return
		}
	}

	t = time.Unix(s, n)
	return
}

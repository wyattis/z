package ztime

import (
	"fmt"
	"strconv"
	"strings"
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
	if strings.ToLower(strings.TrimSpace(val)) == "now" {
		return time.Now(), nil
	}
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

// Parse successfully or panic
func MustParse(val string, formats ...string) (t time.Time) {
	t, err := Parse(val, formats...)
	if err != nil {
		panic(err)
	}
	return
}

// ParseUnix successfully or panic
func MustParseUnix(seconds string, nanos string) (t time.Time) {
	t, err := ParseUnix(seconds, nanos)
	if err != nil {
		panic(err)
	}
	return
}

// Returns true if the two times are within the given duration of each other
func EqualWithin(t1, t2 time.Time, within time.Duration) bool {
	d := t1.Sub(t2)
	return d < within && -d < within
}

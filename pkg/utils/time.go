package utils

import "time"

func IsTimeBetween(start, end, check time.Time) bool {
	start = time.Date(check.Year(), check.Month(), check.Day(), start.Hour(), start.Minute(), start.Second(), start.Nanosecond(), check.Location())
	end = time.Date(check.Year(), check.Month(), check.Day(), end.Hour(), end.Minute(), end.Second(), end.Nanosecond(), check.Location())
	return check.After(start) && check.Before(end)
}

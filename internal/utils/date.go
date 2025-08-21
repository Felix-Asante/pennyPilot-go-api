package utils

import "time"

func HasPassedMinutesAgo(date time.Time, minutes int) bool {
	return time.Since(date) > time.Duration(minutes)*time.Minute
}

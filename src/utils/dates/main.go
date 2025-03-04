package dates

import "time"

func IsValidDate(date string) bool {
	_, err := time.Parse("2006-01-02", date)
	return err == nil
}

func ParseDate(date string) (time.Time, error) {
	return time.Parse("2006-01-02", date)
}

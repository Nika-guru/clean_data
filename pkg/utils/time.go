package utils

import "time"

func Timestamp() string {
	return time.Now().Format("2006-01-02 15:04:05")
}

func StringToTimestamp(timestamp string) (time.Time, error) {
	return time.Parse("2006-01-02T15:04:05Z", timestamp)
}

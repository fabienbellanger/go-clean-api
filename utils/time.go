package utils

import "time"

const SqlDateTimeFormat = "2006-01-02 15:04:05"

// FormatToSqlDateTime formats a time.Time to a string in the format "2006-01-02 15:04:05"
func FormatToSqlDateTime(t time.Time) string {
	return t.Format(SqlDateTimeFormat)
}

// FormatToRFC3339 formats a time.Time to a string in the format "2006-01-02T15:04:05Z07:00"
func FormatToRFC3339(t time.Time) string {
	return t.Format(time.RFC3339)
}

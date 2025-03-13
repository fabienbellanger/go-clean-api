package utils

import "time"

const SqlDateTimeFormat = "2006-01-02 15:04:05"

func FormatToSqlDateTime(t time.Time) string {
	return t.Format(SqlDateTimeFormat)
}

func FormatToRFC3339(t time.Time) string {
	return t.Format(time.RFC3339)
}

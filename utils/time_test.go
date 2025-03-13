package utils

import (
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
)

func TestFormatToSqlDateTime(t *testing.T) {
	tzParis, _ := time.LoadLocation("Europe/Paris")
	tests := []struct {
		name string
		got  string
		want string
	}{
		{
			name: "Test FormatToSqlDateTime with UTC",
			got:  FormatToSqlDateTime(time.Date(2020, 1, 1, 0, 0, 0, 0, time.UTC)),
			want: "2020-01-01 00:00:00",
		},
		{
			name: "Test FormatToSqlDateTime with Europe/Paris timezone",
			got:  FormatToSqlDateTime(time.Date(2020, 1, 1, 0, 0, 0, 0, tzParis)),
			want: "2020-01-01 00:00:00",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			assert.Equal(t, tt.got, tt.want)
		})
	}
}

func TestFormatToRFC3339(t *testing.T) {
	tzParis, _ := time.LoadLocation("Europe/Paris")
	tests := []struct {
		name string
		got  string
		want string
	}{
		{
			name: "Test FormatToRFC3339 with UTC",
			got:  FormatToRFC3339(time.Date(2020, 1, 1, 0, 0, 0, 0, time.UTC)),
			want: "2020-01-01T00:00:00Z",
		},
		{
			name: "Test FormatToRFC3339 with Europe/Paris timezone",
			got:  FormatToRFC3339(time.Date(2020, 1, 1, 0, 0, 0, 0, tzParis)),
			want: "2020-01-01T00:00:00+01:00",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			assert.Equal(t, tt.got, tt.want)
		})
	}
}

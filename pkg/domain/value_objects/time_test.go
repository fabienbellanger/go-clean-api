package values_objects

import (
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
)

func TestNewTime(t *testing.T) {
	locParis, _ := time.LoadLocation("Europe/Paris")
	tests := []struct {
		name  string
		value time.Time
		loc   *time.Location
		want  Time
	}{
		{
			name:  "Test NewTime without location",
			value: time.Date(2020, 1, 1, 0, 0, 0, 0, time.UTC),
			loc:   nil,
			want:  Time{value: time.Date(2020, 1, 1, 0, 0, 0, 0, time.UTC)},
		},
		{
			name:  "Test NewTime with UTC",
			value: time.Date(2020, 1, 1, 0, 0, 0, 0, time.UTC),
			loc:   time.UTC,
			want:  Time{value: time.Date(2020, 1, 1, 0, 0, 0, 0, time.UTC)},
		},
		{
			name:  "Test NewTime with different location",
			value: time.Date(2020, 1, 1, 0, 0, 0, 0, time.UTC),
			loc:   locParis,
			want:  Time{value: time.Date(2020, 1, 1, 1, 0, 0, 0, locParis)},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			assert.Equal(t, NewTime(tt.value, tt.loc), tt.want)
		})
	}
}

func TestTimeSQL(t *testing.T) {
	tests := []struct {
		name  string
		value Time
		want  string
	}{
		{
			name:  "Test Time SQL",
			value: NewTime(time.Date(2020, 1, 1, 0, 0, 0, 0, time.UTC), nil),
			want:  "2020-01-01 00:00:00",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			assert.Equal(t, tt.value.SQL(), tt.want)
		})
	}
}

func TestTimeRFC3339(t *testing.T) {
	locParis, _ := time.LoadLocation("Europe/Paris")
	tests := []struct {
		name  string
		value Time
		want  string
	}{
		{
			name:  "Test Time RFC3339",
			value: NewTime(time.Date(2020, 1, 1, 0, 0, 0, 0, time.UTC), nil),
			want:  "2020-01-01T00:00:00Z",
		},
		{
			name:  "Test Time RFC3339 with different location",
			value: NewTime(time.Date(2020, 1, 1, 0, 0, 0, 0, time.UTC), locParis),
			want:  "2020-01-01T01:00:00+01:00",
		},
		{
			name:  "Test Time RFC3339 with the same location",
			value: NewTime(time.Date(2020, 1, 1, 0, 0, 0, 0, locParis), locParis),
			want:  "2020-01-01T00:00:00+01:00",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			assert.Equal(t, tt.value.RFC3339(), tt.want)
		})
	}
}

func TestTimeParseSQL(t *testing.T) {
	locParis, _ := time.LoadLocation("Europe/Paris")
	tests := []struct {
		name  string
		value string
		loc   *time.Location
		want  Time
		err   bool
	}{
		{
			name:  "Test Time Parse SQL",
			value: "2020-01-01 00:00:00",
			loc:   nil,
			want:  NewTime(time.Date(2020, 1, 1, 0, 0, 0, 0, time.UTC), nil),
			err:   true,
		},
		{
			name:  "Test Time Parse SQL with UTC",
			value: "2020-01-01 00:00:00",
			loc:   time.UTC,
			want:  NewTime(time.Date(2020, 1, 1, 0, 0, 0, 0, time.UTC), nil),
			err:   false,
		},
		{
			name:  "Test Time Parse SQL with location",
			value: "2020-01-01 01:00:00",
			loc:   locParis,
			want:  NewTime(time.Date(2020, 1, 1, 0, 0, 0, 0, time.UTC), locParis),
			err:   false,
		},
		{
			name:  "Test Time Parse SQL with different location",
			value: "2020-01-01 00:00:00",
			loc:   locParis,
			want:  NewTime(time.Date(2020, 1, 1, 0, 0, 0, 0, locParis), locParis),
			err:   false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			value, err := ParseSQL(tt.value, tt.loc)

			if err != nil {
				assert.True(t, tt.err)
			} else {
				assert.Equal(t, value, tt.want)
			}
		})
	}
}

func TestTimeParseRFC3339(t *testing.T) {
	locParis, _ := time.LoadLocation("Europe/Paris")
	tests := []struct {
		name  string
		value string
		loc   *time.Location
		want  Time
		err   bool
	}{
		{
			name:  "Test Time Parse SQL",
			value: "2020-01-01T00:00:00Z",
			loc:   nil,
			want:  NewTime(time.Date(2020, 1, 1, 0, 0, 0, 0, time.UTC), nil),
			err:   false,
		},
		{
			name:  "Test Time Parse SQL with UTC",
			value: "2020-01-01T00:00:00Z",
			loc:   time.UTC,
			want:  NewTime(time.Date(2020, 1, 1, 0, 0, 0, 0, time.UTC), nil),
			err:   false,
		},
		{
			name:  "Test Time Parse SQL with location",
			value: "2020-01-01T01:00:00+01:00",
			loc:   locParis,
			want:  NewTime(time.Date(2020, 1, 1, 0, 0, 0, 0, time.UTC), locParis),
			err:   false,
		},
		{
			name:  "Test Time Parse SQL with different location",
			value: "2020-01-01T00:00:00+01:00",
			loc:   locParis,
			want:  NewTime(time.Date(2020, 1, 1, 0, 0, 0, 0, locParis), locParis),
			err:   false,
		},
		{
			name:  "Test Time Parse SQL with different location",
			value: "2020-01-01T00:00:00+01:00",
			loc:   nil,
			want:  NewTime(time.Date(2020, 1, 1, 0, 0, 0, 0, locParis), nil),
			err:   false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			value, err := ParseRFC3339(tt.value, tt.loc)

			if err != nil {
				assert.True(t, tt.err)
			} else {
				assert.Equal(t, value, tt.want)
			}
		})
	}
}

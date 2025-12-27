package valueobjects

import (
	"strings"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestNewRecurrenceFrequency(t *testing.T) {
	tests := []struct {
		name    string
		value   string
		wantErr bool
	}{
		{
			name:    "valid daily",
			value:   "DAILY",
			wantErr: false,
		},
		{
			name:    "valid weekly",
			value:   "WEEKLY",
			wantErr: false,
		},
		{
			name:    "valid monthly",
			value:   "MONTHLY",
			wantErr: false,
		},
		{
			name:    "valid yearly",
			value:   "YEARLY",
			wantErr: false,
		},
		{
			name:    "lowercase daily",
			value:   "daily",
			wantErr: false,
		},
		{
			name:    "invalid frequency",
			value:   "INVALID",
			wantErr: true,
		},
		{
			name:    "empty string",
			value:   "",
			wantErr: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := NewRecurrenceFrequency(tt.value)
			if tt.wantErr {
				assert.Error(t, err)
			} else {
				assert.NoError(t, err)
				assert.Equal(t, strings.ToUpper(tt.value), got.Value())
			}
		})
	}
}

func TestRecurrenceFrequency_IsDaily(t *testing.T) {
	daily := DailyFrequency()
	assert.True(t, daily.IsDaily())
	assert.False(t, daily.IsWeekly())
	assert.False(t, daily.IsMonthly())
	assert.False(t, daily.IsYearly())
}

func TestRecurrenceFrequency_IsWeekly(t *testing.T) {
	weekly := WeeklyFrequency()
	assert.False(t, weekly.IsDaily())
	assert.True(t, weekly.IsWeekly())
	assert.False(t, weekly.IsMonthly())
	assert.False(t, weekly.IsYearly())
}

func TestRecurrenceFrequency_IsMonthly(t *testing.T) {
	monthly := MonthlyFrequency()
	assert.False(t, monthly.IsDaily())
	assert.False(t, monthly.IsWeekly())
	assert.True(t, monthly.IsMonthly())
	assert.False(t, monthly.IsYearly())
}

func TestRecurrenceFrequency_IsYearly(t *testing.T) {
	yearly := YearlyFrequency()
	assert.False(t, yearly.IsDaily())
	assert.False(t, yearly.IsWeekly())
	assert.False(t, yearly.IsMonthly())
	assert.True(t, yearly.IsYearly())
}

func TestRecurrenceFrequency_Equals(t *testing.T) {
	freq1 := DailyFrequency()
	freq2 := DailyFrequency()
	freq3 := WeeklyFrequency()

	assert.True(t, freq1.Equals(freq2))
	assert.False(t, freq1.Equals(freq3))
}

func TestRecurrenceFrequency_DisplayName(t *testing.T) {
	assert.Equal(t, "Daily", DailyFrequency().DisplayName())
	assert.Equal(t, "Weekly", WeeklyFrequency().DisplayName())
	assert.Equal(t, "Monthly", MonthlyFrequency().DisplayName())
	assert.Equal(t, "Yearly", YearlyFrequency().DisplayName())
}

func TestMustRecurrenceFrequency(t *testing.T) {
	freq := MustRecurrenceFrequency("DAILY")
	assert.Equal(t, Daily, freq.Value())

	assert.Panics(t, func() {
		MustRecurrenceFrequency("INVALID")
	})
}

func TestParseRecurrenceFrequency(t *testing.T) {
	freq, err := ParseRecurrenceFrequency("MONTHLY")
	assert.NoError(t, err)
	assert.Equal(t, Monthly, freq.Value())

	_, err = ParseRecurrenceFrequency("")
	assert.Error(t, err)
}

func TestAllRecurrenceFrequencies(t *testing.T) {
	frequencies := AllRecurrenceFrequencies()
	assert.Len(t, frequencies, 4)
	assert.Contains(t, []string{frequencies[0].Value(), frequencies[1].Value(), frequencies[2].Value(), frequencies[3].Value()}, Daily)
	assert.Contains(t, []string{frequencies[0].Value(), frequencies[1].Value(), frequencies[2].Value(), frequencies[3].Value()}, Weekly)
	assert.Contains(t, []string{frequencies[0].Value(), frequencies[1].Value(), frequencies[2].Value(), frequencies[3].Value()}, Monthly)
	assert.Contains(t, []string{frequencies[0].Value(), frequencies[1].Value(), frequencies[2].Value(), frequencies[3].Value()}, Yearly)
}

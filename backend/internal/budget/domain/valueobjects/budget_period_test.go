package valueobjects

import (
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
)

func TestNewBudgetPeriod(t *testing.T) {
	tests := []struct {
		name       string
		periodType BudgetPeriodType
		year       int
		month      *int
		wantErr    bool
	}{
		{
			name:       "valid monthly period",
			periodType: Monthly,
			year:       2025,
			month:      intPtr(12),
			wantErr:    false,
		},
		{
			name:       "valid yearly period",
			periodType: Yearly,
			year:       2025,
			month:      nil,
			wantErr:    false,
		},
		{
			name:       "invalid period type",
			periodType: "INVALID",
			year:       2025,
			month:      nil,
			wantErr:    true,
		},
		{
			name:       "monthly without month",
			periodType: Monthly,
			year:       2025,
			month:      nil,
			wantErr:    true,
		},
		{
			name:       "yearly with month",
			periodType: Yearly,
			year:       2025,
			month:      intPtr(12),
			wantErr:    true,
		},
		{
			name:       "invalid month",
			periodType: Monthly,
			year:       2025,
			month:      intPtr(13),
			wantErr:    true,
		},
		{
			name:       "invalid year too low",
			periodType: Monthly,
			year:       1800,
			month:      intPtr(1),
			wantErr:    true,
		},
		{
			name:       "invalid year too high",
			periodType: Monthly,
			year:       4000,
			month:      intPtr(1),
			wantErr:    true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := NewBudgetPeriod(tt.periodType, tt.year, tt.month)
			if tt.wantErr {
				assert.Error(t, err)
			} else {
				assert.NoError(t, err)
				assert.Equal(t, tt.periodType, got.PeriodType())
				assert.Equal(t, tt.year, got.Year())
			}
		})
	}
}

func TestNewMonthlyBudgetPeriod(t *testing.T) {
	period, err := NewMonthlyBudgetPeriod(2025, 12)
	assert.NoError(t, err)
	assert.True(t, period.IsMonthly())
	assert.Equal(t, 2025, period.Year())
	assert.NotNil(t, period.Month())
	assert.Equal(t, 12, *period.Month())
}

func TestNewYearlyBudgetPeriod(t *testing.T) {
	period, err := NewYearlyBudgetPeriod(2025)
	assert.NoError(t, err)
	assert.True(t, period.IsYearly())
	assert.Equal(t, 2025, period.Year())
	assert.Nil(t, period.Month())
}

func TestBudgetPeriod_StartDate(t *testing.T) {
	monthly, _ := NewMonthlyBudgetPeriod(2025, 12)
	start := monthly.StartDate()
	assert.Equal(t, 2025, start.Year())
	assert.Equal(t, time.December, start.Month())
	assert.Equal(t, 1, start.Day())

	yearly, _ := NewYearlyBudgetPeriod(2025)
	start = yearly.StartDate()
	assert.Equal(t, 2025, start.Year())
	assert.Equal(t, time.January, start.Month())
	assert.Equal(t, 1, start.Day())
}

func TestBudgetPeriod_EndDate(t *testing.T) {
	monthly, _ := NewMonthlyBudgetPeriod(2025, 12)
	end := monthly.EndDate()
	assert.Equal(t, 2025, end.Year())
	assert.Equal(t, time.December, end.Month())
	assert.Equal(t, 31, end.Day())

	yearly, _ := NewYearlyBudgetPeriod(2025)
	end = yearly.EndDate()
	assert.Equal(t, 2025, end.Year())
	assert.Equal(t, time.December, end.Month())
	assert.Equal(t, 31, end.Day())
}

func TestBudgetPeriod_Includes(t *testing.T) {
	monthly, _ := NewMonthlyBudgetPeriod(2025, 12)

	// Date within period
	date1 := time.Date(2025, 12, 15, 0, 0, 0, 0, time.UTC)
	assert.True(t, monthly.Includes(date1))

	// Date before period
	date2 := time.Date(2025, 11, 30, 0, 0, 0, 0, time.UTC)
	assert.False(t, monthly.Includes(date2))

	// Date after period
	date3 := time.Date(2026, 1, 1, 0, 0, 0, 0, time.UTC)
	assert.False(t, monthly.Includes(date3))

	// Start date
	date4 := monthly.StartDate()
	assert.True(t, monthly.Includes(date4))

	// End date
	date5 := monthly.EndDate()
	assert.True(t, monthly.Includes(date5))
}

func TestBudgetPeriod_Equals(t *testing.T) {
	period1, _ := NewMonthlyBudgetPeriod(2025, 12)
	period2, _ := NewMonthlyBudgetPeriod(2025, 12)
	period3, _ := NewMonthlyBudgetPeriod(2025, 11)
	period4, _ := NewYearlyBudgetPeriod(2025)

	assert.True(t, period1.Equals(period2))
	assert.False(t, period1.Equals(period3))
	assert.False(t, period1.Equals(period4))
}

func intPtr(i int) *int {
	return &i
}

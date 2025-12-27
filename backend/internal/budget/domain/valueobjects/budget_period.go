package valueobjects

import (
	"errors"
	"fmt"
	"time"
)

// BudgetPeriodType represents the type of budget period.
type BudgetPeriodType string

const (
	Monthly BudgetPeriodType = "MONTHLY"
	Yearly  BudgetPeriodType = "YEARLY"
)

// BudgetPeriod represents a budget period value object.
type BudgetPeriod struct {
	periodType BudgetPeriodType
	year       int
	month      *int // nil for yearly periods
}

// NewBudgetPeriod creates a new BudgetPeriod.
func NewBudgetPeriod(periodType BudgetPeriodType, year int, month *int) (BudgetPeriod, error) {
	// Validate year
	if year < 1900 || year > 3000 {
		return BudgetPeriod{}, errors.New("year must be between 1900 and 3000")
	}

	// Validate period type
	if periodType != Monthly && periodType != Yearly {
		return BudgetPeriod{}, fmt.Errorf("invalid period type: %s. Supported values: MONTHLY, YEARLY", periodType)
	}

	// For monthly periods, month is required
	if periodType == Monthly {
		if month == nil {
			return BudgetPeriod{}, errors.New("month is required for monthly periods")
		}
		if *month < 1 || *month > 12 {
			return BudgetPeriod{}, errors.New("month must be between 1 and 12")
		}
	} else {
		// For yearly periods, month should be nil
		if month != nil {
			return BudgetPeriod{}, errors.New("month must be nil for yearly periods")
		}
	}

	return BudgetPeriod{
		periodType: periodType,
		year:       year,
		month:      month,
	}, nil
}

// NewMonthlyBudgetPeriod creates a new monthly budget period.
func NewMonthlyBudgetPeriod(year int, month int) (BudgetPeriod, error) {
	return NewBudgetPeriod(Monthly, year, &month)
}

// NewYearlyBudgetPeriod creates a new yearly budget period.
func NewYearlyBudgetPeriod(year int) (BudgetPeriod, error) {
	return NewBudgetPeriod(Yearly, year, nil)
}

// PeriodType returns the period type (MONTHLY or YEARLY).
func (bp BudgetPeriod) PeriodType() BudgetPeriodType {
	return bp.periodType
}

// Year returns the year.
func (bp BudgetPeriod) Year() int {
	return bp.year
}

// Month returns the month (nil for yearly periods).
func (bp BudgetPeriod) Month() *int {
	return bp.month
}

// IsMonthly checks if the period is monthly.
func (bp BudgetPeriod) IsMonthly() bool {
	return bp.periodType == Monthly
}

// IsYearly checks if the period is yearly.
func (bp BudgetPeriod) IsYearly() bool {
	return bp.periodType == Yearly
}

// StartDate returns the start date of the period.
func (bp BudgetPeriod) StartDate() time.Time {
	if bp.IsMonthly() {
		return time.Date(bp.year, time.Month(*bp.month), 1, 0, 0, 0, 0, time.UTC)
	}
	return time.Date(bp.year, 1, 1, 0, 0, 0, 0, time.UTC)
}

// EndDate returns the end date of the period.
func (bp BudgetPeriod) EndDate() time.Time {
	if bp.IsMonthly() {
		// Last day of the month
		nextMonth := time.Date(bp.year, time.Month(*bp.month)+1, 1, 0, 0, 0, 0, time.UTC)
		return nextMonth.AddDate(0, 0, -1)
	}
	// Last day of the year
	return time.Date(bp.year, 12, 31, 23, 59, 59, 999999999, time.UTC)
}

// Includes checks if a date is within the period.
func (bp BudgetPeriod) Includes(date time.Time) bool {
	start := bp.StartDate()
	end := bp.EndDate()
	return (date.Equal(start) || date.After(start)) && (date.Equal(end) || date.Before(end))
}

// Equals checks if two BudgetPeriod values are equal.
func (bp BudgetPeriod) Equals(other BudgetPeriod) bool {
	if bp.periodType != other.periodType || bp.year != other.year {
		return false
	}
	if bp.month == nil && other.month == nil {
		return true
	}
	if bp.month == nil || other.month == nil {
		return false
	}
	return *bp.month == *other.month
}

// String returns a string representation of the period.
func (bp BudgetPeriod) String() string {
	if bp.IsMonthly() {
		return fmt.Sprintf("%s-%02d", time.Month(*bp.month).String()[:3], bp.year)
	}
	return fmt.Sprintf("%d", bp.year)
}

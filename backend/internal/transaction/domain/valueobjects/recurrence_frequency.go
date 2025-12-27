package valueobjects

import (
	"errors"
	"fmt"
	"strings"
)

// RecurrenceFrequency represents a recurrence frequency value object.
type RecurrenceFrequency struct {
	value string
}

// Valid recurrence frequency values
const (
	Daily   = "DAILY"   // Daily recurrence
	Weekly  = "WEEKLY"  // Weekly recurrence
	Monthly = "MONTHLY" // Monthly recurrence
	Yearly  = "YEARLY"  // Yearly recurrence
)

// ValidRecurrenceFrequencies is a map of all supported recurrence frequencies.
var ValidRecurrenceFrequencies = map[string]string{
	Daily:   "Daily",
	Weekly:  "Weekly",
	Monthly: "Monthly",
	Yearly:  "Yearly",
}

// NewRecurrenceFrequency creates a new RecurrenceFrequency value object.
func NewRecurrenceFrequency(value string) (RecurrenceFrequency, error) {
	value = strings.ToUpper(strings.TrimSpace(value))

	if !IsValidRecurrenceFrequency(value) {
		return RecurrenceFrequency{}, fmt.Errorf("invalid recurrence frequency: %s. Supported values: DAILY, WEEKLY, MONTHLY, YEARLY", value)
	}

	return RecurrenceFrequency{value: value}, nil
}

// MustRecurrenceFrequency creates a new RecurrenceFrequency value object and panics if the value is invalid.
func MustRecurrenceFrequency(value string) RecurrenceFrequency {
	rf, err := NewRecurrenceFrequency(value)
	if err != nil {
		panic(err)
	}
	return rf
}

// IsValidRecurrenceFrequency checks if a recurrence frequency value is valid.
func IsValidRecurrenceFrequency(value string) bool {
	value = strings.ToUpper(strings.TrimSpace(value))
	_, exists := ValidRecurrenceFrequencies[value]
	return exists
}

// Value returns the recurrence frequency value.
func (rf RecurrenceFrequency) Value() string {
	return rf.value
}

// String returns the recurrence frequency value as a string.
func (rf RecurrenceFrequency) String() string {
	return rf.value
}

// DisplayName returns the human-readable name of the recurrence frequency.
func (rf RecurrenceFrequency) DisplayName() string {
	if name, exists := ValidRecurrenceFrequencies[rf.value]; exists {
		return name
	}
	return rf.value
}

// IsDaily checks if the recurrence frequency is Daily.
func (rf RecurrenceFrequency) IsDaily() bool {
	return rf.value == Daily
}

// IsWeekly checks if the recurrence frequency is Weekly.
func (rf RecurrenceFrequency) IsWeekly() bool {
	return rf.value == Weekly
}

// IsMonthly checks if the recurrence frequency is Monthly.
func (rf RecurrenceFrequency) IsMonthly() bool {
	return rf.value == Monthly
}

// IsYearly checks if the recurrence frequency is Yearly.
func (rf RecurrenceFrequency) IsYearly() bool {
	return rf.value == Yearly
}

// Equals checks if two RecurrenceFrequency values are equal.
func (rf RecurrenceFrequency) Equals(other RecurrenceFrequency) bool {
	return rf.value == other.value
}

// DailyFrequency returns a RecurrenceFrequency for Daily.
func DailyFrequency() RecurrenceFrequency {
	return RecurrenceFrequency{value: Daily}
}

// WeeklyFrequency returns a RecurrenceFrequency for Weekly.
func WeeklyFrequency() RecurrenceFrequency {
	return RecurrenceFrequency{value: Weekly}
}

// MonthlyFrequency returns a RecurrenceFrequency for Monthly.
func MonthlyFrequency() RecurrenceFrequency {
	return RecurrenceFrequency{value: Monthly}
}

// YearlyFrequency returns a RecurrenceFrequency for Yearly.
func YearlyFrequency() RecurrenceFrequency {
	return RecurrenceFrequency{value: Yearly}
}

// ParseRecurrenceFrequency attempts to parse a recurrence frequency from a string.
func ParseRecurrenceFrequency(s string) (RecurrenceFrequency, error) {
	if s == "" {
		return RecurrenceFrequency{}, errors.New("recurrence frequency cannot be empty")
	}

	return NewRecurrenceFrequency(s)
}

// AllRecurrenceFrequencies returns all valid recurrence frequencies.
func AllRecurrenceFrequencies() []RecurrenceFrequency {
	return []RecurrenceFrequency{
		DailyFrequency(),
		WeeklyFrequency(),
		MonthlyFrequency(),
		YearlyFrequency(),
	}
}

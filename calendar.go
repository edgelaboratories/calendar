package calendar

import "github.com/fxtlabs/date"

// Calendar exposes functions to manipulate dates with respect to a calendar.
type Calendar struct {
	dayCounter dayCounter
	convention Convention
}

// New returns a calendar based on the specified input convention.
func New(convention Convention) *Calendar {
	return &Calendar{
		dayCounter: newDayCounter(convention),
		convention: convention,
	}
}

// IsActive returns true if the input date is active.
func (c *Calendar) IsActive(date date.Date) bool {
	return c.dayCounter.isActive(date)
}

// add adds an input number of active days to the input origin date.
// The days parameter is allowed to be negative.
// This method is idempotent for zero-days shifts.
func (c *Calendar) Add(origin date.Date, days int) date.Date {
	return c.dayCounter.add(origin, days)
}

// DaysBetween computes the number of active dates between
// from (excluded) and to (included).
func (c *Calendar) DaysBetween(from, to date.Date) int {
	return c.dayCounter.daysBetween(from, to)
}

// DaysInYear returns the calendar's standard year duration (in days).
func (c *Calendar) DaysInYear() int {
	return c.dayCounter.daysInYear()
}

// LatestBefore returns the latest date before or equal to
// an input date. As opposed to the input date, the output date
// belongs to the calendar by construction.
func (c *Calendar) LatestBefore(date date.Date) date.Date {
	// Obtain the most recent date by means of a zero-days shift.
	return c.Add(date, 0)
}

// Convention returns the calendar convention.
func (c *Calendar) Convention() Convention {
	return c.convention
}

package calendar

import "github.com/fxtlabs/date"

// Calendar exposes functions to manipulate dates with respect to a calendar.
type Calendar struct {
	dayCounter

	convention Convention
}

// New returns a calendar based on the specified input convention.
func New(convention Convention) *Calendar {
	return &Calendar{
		newDayCounter(convention),
		convention,
	}
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

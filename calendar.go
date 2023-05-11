package calendar

import "github.com/fxtlabs/date"

// dayCounter defines the properties of a calendar.
type dayCounter interface {
	// Convention returns the calendar convention.
	Convention() Convention
	// IsActive returns true if the input date is active.
	IsActive(date date.Date) bool
	// DaysInYear returns the calendar's standard year duration (in days).
	DaysInYear() int
	// DaysBetween computes the number of active dates between
	// from (excluded) and to (included).
	DaysBetween(from, to date.Date) int
	// add adds an input number of active days to the input origin date.
	// The days parameter is allowed to be negative.
	// This method is idempotent for zero-days shifts.
	Add(origin date.Date, days int) date.Date
}

// Calendar exposes functions to manipulate dates with respect to a calendar.
type Calendar struct {
	dayCounter
}

// New returns a calendar based on the specified input convention.
func New(convention Convention) *Calendar {
	switch convention {
	case CalendarDays:
		return &Calendar{newPhysicalCalendar()}
	case BusinessDays:
		fallthrough

	default:
		return &Calendar{newBusinessCalendar()}
	}
}

// LatestBefore returns the latest date before or equal to
// an input date. As opposed to the input date, the output date
// belongs to the calendar by construction.
func (c *Calendar) LatestBefore(date date.Date) date.Date {
	// Obtain the most recent date by means of a zero-days shift.
	return c.Add(date, 0)
}

// Next returns the next calendar date with respect to the input.
func (c *Calendar) Next(date date.Date) date.Date {
	return c.Add(date, 1)
}

// Previous returns the previous calendar date with respect to the input.
func (c *Calendar) Previous(date date.Date) date.Date {
	return c.Add(date, -1)
}

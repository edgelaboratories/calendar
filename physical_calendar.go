package calendar

import "github.com/fxtlabs/date"

// physicalCalendar defines a calendar in which all days
// (including weekends) are active.
type physicalCalendar struct{}

func newPhysicalCalendar() *physicalCalendar {
	return &physicalCalendar{}
}

// isActive returns whether the input date is active according
// to the physical calendar. By definition the result is true
// for every input date.
func (c physicalCalendar) IsActive(date date.Date) bool {
	return true
}

// yearDuration returns the standard year duration according to the
// physical-days calendar.
func (c physicalCalendar) DaysInYear() int {
	return 365
}

// Add adds an input number of active days to the input origin date.
// The days parameter is allowed to be negative.
func (c physicalCalendar) Add(origin date.Date, days int) date.Date {
	return origin.Add(days)
}

// DaysBetween computes the number of active dates between
// from (excluded) and to (included).
// The from parameter is supposed not to be after to.
func (c physicalCalendar) DaysBetween(from, to date.Date) int {
	return to.Sub(from)
}

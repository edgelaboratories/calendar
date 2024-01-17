package calendar

import (
	"time"

	"github.com/edgelaboratories/date"
)

// businessCalendar is a calendar whose active days are working ones.
// Weekends are not active, but no further holidays are taken into account.
type businessCalendar struct{}

func newBusinessCalendar() *businessCalendar {
	return &businessCalendar{}
}

// Convention returns the BusinessDays convention.
func (c businessCalendar) Convention() Convention {
	return BusinessDays
}

// IsActive returns true if the input date is active.
func (c businessCalendar) IsActive(date date.Date) bool {
	return c.isBusinessDay(date)
}

// DaysInYear returns the standard year duration according to the
// business-days calendar.
func (c businessCalendar) DaysInYear() int {
	return 252
}

// Add adds an input number of active days to the input origin date.
// The days parameter is allowed to be negative.
// This method is idempotent when a zero-days shift is requested.
func (c businessCalendar) Add(origin date.Date, days int) date.Date {
	// go back to last Friday if it is a weekend
	current := c.previousBusinessDay(origin)

	if days == 0 {
		return current
	}

	// Number of weeks to shift (forward or backward)
	nbWeeks := days / 5

	// Number of business days left after shifting by nbWeeks
	nbDaysLeft := days - (nbWeeks * 5)

	// If the nbDaysLeft does not fit in the current week, then add two days for the weekend.
	nbDaysLeft = addWeekendDays(current, nbDaysLeft, days)

	// Total number of days (this number will eventually be incremented below)
	nbDays := nbWeeks*7 + nbDaysLeft

	return current.Add(nbDays)
}

// DaysBetween computes the number of active dates between
// from (excluded) and to (included).
// This implies there are zero business days from a Friday to
// a weekend day, but there is one between a weekend day and
// a Monday.
func (c businessCalendar) DaysBetween(from, to date.Date) int {
	if from.After(to) {
		return -c.DaysBetween(to, from)
	}

	start, end := from, to

	// Shift start date to next Sunday if on Friday or Saturday.
	switch start.Weekday() {
	case time.Friday:
		start = start.Add(2)

	case time.Saturday:
		start = start.Add(1)

	case time.Monday,
		time.Tuesday,
		time.Wednesday,
		time.Thursday,
		time.Sunday:
	}

	// Shift end date to previous Friday if on Saturday or Sunday.
	switch end.Weekday() {
	case time.Saturday:
		end = end.Add(-1)

	case time.Sunday:
		end = end.Add(-2)

	case time.Monday,
		time.Tuesday,
		time.Wednesday,
		time.Thursday,
		time.Friday:
	}

	// Compute raw day difference.
	rawDaysDiff := end.Sub(start)
	if rawDaysDiff <= 0 {
		return 0
	}

	// Remove a supplementary weekend if start weekday is after end weekday.
	if start.Weekday() > end.Weekday() {
		rawDaysDiff -= 2
	}

	return rawDaysDiff/7*5 + rawDaysDiff%7
}

// isWeekend returns true if the input date is a non-business
// day according to the calendar.
func (c businessCalendar) isWeekend(date date.Date) bool {
	w := date.Weekday()
	return w == time.Saturday || w == time.Sunday
}

// isBusinessDay returns true if the input date is a business
// day according to the calendar.
func (c businessCalendar) isBusinessDay(date date.Date) bool {
	return !c.isWeekend(date)
}

func (c businessCalendar) nextBusinessDay(origin date.Date) date.Date {
	return c.closestBusinessDay(origin, true)
}

func (c businessCalendar) previousBusinessDay(origin date.Date) date.Date {
	return c.closestBusinessDay(origin, false)
}

// businessDayAfter finds the closest business day
// (forwards or backwards) to an origin date.
func (c businessCalendar) closestBusinessDay(origin date.Date, forwards bool) date.Date {
	shift := -1
	if forwards {
		shift = 1
	}

	for current := origin; ; current = current.Add(shift) {
		if c.isBusinessDay(current) {
			return current
		}
	}
}

// addWeekendDays adds or removes 2 days to the day count (nbDaysLeft) to skip the weekend days.
func addWeekendDays(current date.Date, nbDaysLeft, days int) int {
	if days > 0 && int(time.Friday-current.Weekday()) < nbDaysLeft {
		return nbDaysLeft + 2
	}

	if days < 0 && int(current.Weekday()-time.Monday) < -nbDaysLeft {
		return nbDaysLeft - 2
	}

	return nbDaysLeft
}

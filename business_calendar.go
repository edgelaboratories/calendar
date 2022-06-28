package calendar

import (
	"time"

	"github.com/fxtlabs/date"
)

// businessCalendar is a calendar whose active days are working ones.
// Weekends are not active, but no further holidays are taken into account.
type businessCalendar struct{}

func newBusinessCalendar() *businessCalendar {
	return &businessCalendar{}
}

// isActive returns true if the input date is active.
func (c businessCalendar) isActive(date date.Date) bool {
	return c.isBusinessDay(date)
}

// daysInYear returns the standard year duration according to the
// business-days calendar.
func (c businessCalendar) daysInYear() int {
	return 252
}

// add adds an input number of active days to the input origin date.
// The days parameter is allowed to be negative.
// This method is idempotent when a zero-days shift is requested.
func (c businessCalendar) add(origin date.Date, days int) date.Date {
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

// daysBetween computes the number of active dates between
// from (excluded) and to (included).
// This implies there are zero business days from a Friday to
// a weekend day, but there is one between a weekend day and
// a Monday.
func (c businessCalendar) daysBetween(from, to date.Date) int {
	if from.After(to) {
		return -c.daysBetween(to, from)
	}

	// Shift both dates to the closest Sunday and then add
	// the weekly shifts in business days.
	i := 0
	start, end := from, to

	if end.Weekday() == time.Saturday {
		end = end.Add(1)
	}
	for end.After(start) && !c.isWeekend(end) {
		end = end.Add(-1)
		i++
	}

	for end.After(start) && !c.isWeekend(start) {
		start = start.Add(1)
		if !c.isWeekend(start) {
			i++
		}
	}
	if start.Weekday() == time.Saturday {
		start = start.Add(1)
	}

	// The ratio between business and weekly days is 5/7.
	return i + end.Sub(start)*5/7
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

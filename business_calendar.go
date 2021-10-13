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
func (c businessCalendar) IsActive(date date.Date) bool {
	return c.isBusinessDay(date)
}

// yearDuration returns the standard year duration according to the
// business-days calendar.
func (c businessCalendar) DaysInYear() int {
	return 252
}

// Add adds an input number of active days to the input origin date.
// The days parameter is allowed to be negative.
// This method is idempotent when a zero-days shift is requested.
func (c businessCalendar) Add(origin date.Date, days int) date.Date {
	if days == 0 {
		// In order for the method to be idempotent, the same result must
		// be obtained by shifting forwards and backwards (or viceversa)
		// by a single day. If the origin is already a business day, no
		// shift is applied.
		return c.previousBusinessDay(origin)
	}

	if days > 0 {
		current := c.nextBusinessDay(origin)
		signedDays := days
		if c.isWeekend(origin) {
			signedDays--
		}

		// Count from the first day of the week and go back to the previous
		// business day when landing on a weekend.
		weekDay := int(current.Weekday())
		dayShift := signedDays%5 + weekDay
		weekShift := signedDays/5 + dayShift/5

		return c.previousBusinessDay(current.Add(7*weekShift + dayShift%5 - weekDay))
	}

	// The algorithm runs in linear time for negative shifts.
	// This should be improved at some point.
	current := c.previousBusinessDay(origin)
	for shiftDays := (-days) % 5; shiftDays > 0; {
		current = current.Add(-1)
		if c.isBusinessDay(current) {
			shiftDays--
		}
	}

	return current.Add((days / 5) * 7)
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

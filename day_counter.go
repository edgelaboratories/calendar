package calendar

import "github.com/fxtlabs/date"

// dayCounter defines the properties of a calendar convention:
// - active (working) date criterion;
// - standard year duration;
// - days between two dates;
// - addition of days to an origin date.
// The methods of the daycounter are implementation-dependent
// and allow to exploit the calendar properties.
type dayCounter interface {
	isActive(date date.Date) bool
	daysInYear() int
	daysBetween(from, to date.Date) int
	add(origin date.Date, days int) date.Date
}

func newDayCounter(convention Convention) dayCounter {
	switch convention {
	case BusinessDays:
		return newBusinessCalendar()

	case CalendarDays:
		return newPhysicalCalendar()

	default:
		return newBusinessCalendar()
	}
}

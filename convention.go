package calendar

// Convention defines the calendar convention.
type Convention string

const (
	// BusinessDays uses a no-holiday calendar.
	// Only working days are active, weekends are neglected.
	BusinessDays Convention = "BusinessDays"
	// CalendarDays defines a convention that uses a nominal ISO calendar.
	// All days, including weekends, are considered active.
	CalendarDays Convention = "CalendarDays"
)

package calendar

import (
	"testing"
	"time"

	"github.com/edgelaboratories/date"
	"github.com/stretchr/testify/assert"
)

func Test_New(t *testing.T) {
	t.Parallel()

	for _, tc := range []struct {
		name       string
		convention Convention
	}{
		{
			"business",
			BusinessDays,
		},
		{
			"physical",
			CalendarDays,
		},
		{
			"unspecified",
			BusinessDays,
		},
	} {
		tc := tc

		t.Run(tc.name, func(t *testing.T) {
			t.Parallel()

			assert.Equal(t, tc.convention, New(tc.convention).Convention())
		})
	}
}

func Test_Calendar_Convention(t *testing.T) {
	t.Parallel()

	for _, convention := range []Convention{
		CalendarDays,
		BusinessDays,
	} {
		assert.Equal(t, convention, New(convention).Convention())
	}
}

func Test_Calendar_IsActive(t *testing.T) {
	t.Parallel()

	for _, tc := range []struct {
		convention Convention
		date       date.Date
		expected   bool
	}{
		{
			BusinessDays,
			date.New(2020, time.January, 5),
			false,
		},
		{
			BusinessDays,
			date.New(2020, time.January, 6),
			true,
		},
		{
			CalendarDays,
			date.New(2020, time.April, 11),
			true,
		},
	} {
		assert.Equal(t, tc.expected, New(tc.convention).IsActive(tc.date))
	}
}

func Test_Calendar_DaysInYear(t *testing.T) {
	t.Parallel()

	assert.Equal(t, 252, New(BusinessDays).DaysInYear())
	assert.Equal(t, 365, New(CalendarDays).DaysInYear())
}

func Test_Calendar_Add(t *testing.T) {
	t.Parallel()

	for _, tc := range []struct {
		convention Convention
		origin     date.Date
		days       int
		expected   date.Date
	}{
		{
			CalendarDays,
			date.New(2021, time.October, 18),
			7,
			date.New(2021, time.October, 25),
		},
		{
			BusinessDays,
			date.New(2021, time.October, 18),
			5,
			date.New(2021, time.October, 25),
		},
	} {
		assert.Equal(t, tc.expected, New(tc.convention).Add(tc.origin, tc.days))
	}
}

func Test_Calendar_DaysBetween(t *testing.T) {
	t.Parallel()

	for _, tc := range []struct {
		convention Convention
		from       date.Date
		to         date.Date
		expected   int
	}{
		{
			CalendarDays,
			date.New(2021, time.October, 18),
			date.New(2021, time.October, 25),
			7,
		},
		{
			BusinessDays,
			date.New(2021, time.October, 18),
			date.New(2021, time.October, 25),
			5,
		},
	} {
		assert.Equal(t, tc.expected, New(tc.convention).DaysBetween(tc.from, tc.to))
	}
}

func Benchmark_Calendar_Add(b *testing.B) {
	var (
		origin = date.New(2021, time.January, 1)
		days   = 365
	)

	for _, bc := range []struct {
		name string
		c    *Calendar
	}{
		{
			"business",
			New(BusinessDays),
		},
		{
			"physical",
			New(CalendarDays),
		},
	} {
		bc := bc

		b.Run(bc.name, func(b *testing.B) {
			b.RunParallel(func(p *testing.PB) {
				for p.Next() {
					_ = bc.c.Add(origin, days)
				}
			})
		})
	}
}

func Benchmark_Calendar_DaysBetween(b *testing.B) {
	var (
		from = date.New(2021, time.January, 1)
		to   = date.New(2022, time.January, 1)
	)

	for _, bc := range []struct {
		name string
		c    *Calendar
	}{
		{
			"business",
			New(BusinessDays),
		},
		{
			"physical",
			New(CalendarDays),
		},
	} {
		bc := bc

		b.Run(bc.name, func(b *testing.B) {
			b.RunParallel(func(p *testing.PB) {
				for p.Next() {
					_ = bc.c.DaysBetween(from, to)
				}
			})
		})
	}
}

func Test_Calendar_LatestBefore_Next_Previous(t *testing.T) {
	t.Parallel()

	for _, tc := range []struct {
		name     string
		calendar *Calendar
		date     date.Date
		latest   date.Date
		previous date.Date
		next     date.Date
	}{
		{
			"calendar/business day",
			New(CalendarDays),
			date.New(2021, time.September, 30),
			date.New(2021, time.September, 30),
			date.New(2021, time.September, 29),
			date.New(2021, time.October, 1),
		},
		{
			"calendar/saturday",
			New(CalendarDays),
			date.New(2021, time.October, 2),
			date.New(2021, time.October, 2),
			date.New(2021, time.October, 1),
			date.New(2021, time.October, 3),
		},
		{
			"calendar/sunday",
			New(CalendarDays),
			date.New(2021, time.October, 3),
			date.New(2021, time.October, 3),
			date.New(2021, time.October, 2),
			date.New(2021, time.October, 4),
		},
		{
			"calendar/monday",
			New(CalendarDays),
			date.New(2021, time.October, 4),
			date.New(2021, time.October, 4),
			date.New(2021, time.October, 3),
			date.New(2021, time.October, 5),
		},
		{
			"business/business day",
			New(BusinessDays),
			date.New(2021, time.September, 30),
			date.New(2021, time.September, 30),
			date.New(2021, time.September, 29),
			date.New(2021, time.October, 1),
		},
		{
			"business/saturday",
			New(BusinessDays),
			date.New(2021, time.October, 2),
			date.New(2021, time.October, 1),
			date.New(2021, time.September, 30),
			date.New(2021, time.October, 4),
		},
		{
			"business/sunday",
			New(BusinessDays),
			date.New(2021, time.October, 3),
			date.New(2021, time.October, 1),
			date.New(2021, time.September, 30),
			date.New(2021, time.October, 4),
		},
		{
			"business/monday",
			New(BusinessDays),
			date.New(2021, time.October, 4),
			date.New(2021, time.October, 4),
			date.New(2021, time.October, 1),
			date.New(2021, time.October, 5),
		},
	} {
		tc := tc
		t.Run(tc.name, func(t *testing.T) {
			t.Parallel()

			assert.Equal(t, tc.latest, tc.calendar.LatestBefore(tc.date))
			assert.Equal(t, tc.previous, tc.calendar.Previous(tc.date))
			assert.Equal(t, tc.next, tc.calendar.Next(tc.date))
		})
	}
}

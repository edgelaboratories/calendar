package calendar

import (
	"testing"
	"time"

	"github.com/fxtlabs/date"
	"github.com/stretchr/testify/assert"
)

func Test_New(t *testing.T) {
	t.Run("business", func(t *testing.T) {
		_, ok := New(BusinessDays).dayCounter.(*businessCalendar)
		assert.True(t, ok)
	})

	t.Run("physical", func(t *testing.T) {
		_, ok := New(CalendarDays).dayCounter.(*physicalCalendar)
		assert.True(t, ok)
	})

	t.Run("unspecified", func(t *testing.T) {
		_, ok := New(Convention("unspecified")).dayCounter.(*businessCalendar)
		assert.True(t, ok)
	})
}

func Test_Calendar_Convention(t *testing.T) {
	for _, convention := range []Convention{
		CalendarDays,
		BusinessDays,
	} {
		assert.Equal(t, convention, New(convention).Convention())
	}
}

func Test_Calendar_IsActive(t *testing.T) {
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
	assert.Equal(t, 252, New(BusinessDays).DaysInYear())
	assert.Equal(t, 365, New(CalendarDays).DaysInYear())
}

func Benchmark_Calendar_Add(b *testing.B) {
	var (
		origin = date.New(2021, time.January, 1)
		days   = 365
	)

	b.Run("business", func(b *testing.B) {
		c := New(BusinessDays)

		b.ResetTimer()
		for i := 0; i < b.N; i++ {
			_ = c.Add(origin, days)
		}
	})

	b.Run("physical", func(b *testing.B) {
		c := New(CalendarDays)

		b.ResetTimer()
		for i := 0; i < b.N; i++ {
			_ = c.Add(origin, days)
		}
	})
}

func Benchmark_Calendar_DaysBetween(b *testing.B) {
	var (
		from = date.New(2021, time.January, 1)
		to   = date.New(2022, time.January, 1)
	)

	b.Run("business", func(b *testing.B) {
		c := New(BusinessDays)

		b.ResetTimer()
		for i := 0; i < b.N; i++ {
			_ = c.DaysBetween(from, to)
		}
	})

	b.Run("physical", func(b *testing.B) {
		c := New(CalendarDays)

		b.ResetTimer()
		for i := 0; i < b.N; i++ {
			_ = c.DaysBetween(from, to)
		}
	})
}

func Test_Calendar_LatestBefore(t *testing.T) {
	for _, tc := range []struct {
		name     string
		calendar *Calendar
		date     date.Date
		expected date.Date
	}{
		{
			"calendar/business day",
			New(CalendarDays),
			date.New(2021, time.September, 30),
			date.New(2021, time.September, 30),
		},
		{
			"calendar/saturday",
			New(CalendarDays),
			date.New(2021, time.October, 2),
			date.New(2021, time.October, 2),
		},
		{
			"calendar/sunday",
			New(CalendarDays),
			date.New(2021, time.October, 3),
			date.New(2021, time.October, 3),
		},
		{
			"calendar/monday",
			New(CalendarDays),
			date.New(2021, time.October, 4),
			date.New(2021, time.October, 4),
		},
		{
			"business/business day",
			New(BusinessDays),
			date.New(2021, time.September, 30),
			date.New(2021, time.September, 30),
		},
		{
			"business/saturday",
			New(BusinessDays),
			date.New(2021, time.October, 2),
			date.New(2021, time.October, 1),
		},
		{
			"business/sunday",
			New(BusinessDays),
			date.New(2021, time.October, 3),
			date.New(2021, time.October, 1),
		},
		{
			"business/monday",
			New(BusinessDays),
			date.New(2021, time.October, 4),
			date.New(2021, time.October, 4),
		},
	} {
		tc := tc
		t.Run(tc.name, func(t *testing.T) {
			assert.Equal(t, tc.expected, tc.calendar.LatestBefore(tc.date))
		})
	}
}

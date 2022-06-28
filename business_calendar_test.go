package calendar

import (
	"testing"
	"time"

	"github.com/fxtlabs/date"
	"github.com/stretchr/testify/assert"
)

func Test_businessCalendar_isActive(t *testing.T) {
	calendar := newBusinessCalendar()

	for _, tc := range []struct {
		date     date.Date
		expected bool
	}{
		{
			date.New(2020, time.January, 3),
			true,
		},
		{
			date.New(2020, time.January, 4),
			false,
		},
		{
			date.New(2020, time.January, 5),
			false,
		},
		{
			date.New(2020, time.January, 6),
			true,
		},
	} {
		assert.Equal(t, tc.expected, calendar.isActive(tc.date))
	}
}

func Test_businessCalendar_daysInYear(t *testing.T) {
	assert.Equal(t, 252, newBusinessCalendar().daysInYear())
}

func Test_businessCalendar_add(t *testing.T) {
	calendar := newBusinessCalendar()

	for _, tc := range []struct {
		origin   date.Date
		days     int
		expected date.Date
	}{
		{
			date.New(2020, time.September, 1),
			0,
			date.New(2020, time.September, 1),
		},
		{
			date.New(2020, time.August, 30),
			0,
			date.New(2020, time.August, 28),
		},
		{
			date.New(2020, time.September, 1),
			3,
			date.New(2020, time.September, 4),
		},
		{
			date.New(2020, time.August, 30),
			1,
			date.New(2020, time.August, 31),
		},
		{
			date.New(2020, time.September, 1),
			-1,
			date.New(2020, time.August, 31),
		},
		{
			date.New(2020, time.August, 30),
			-1,
			date.New(2020, time.August, 27),
		},
		{
			date.New(2020, time.September, 1),
			-3,
			date.New(2020, time.August, 27),
		},
		{
			date.New(2020, time.September, 1),
			4,
			date.New(2020, time.September, 7),
		},
	} {
		assert.Equal(t, tc.expected, calendar.add(tc.origin, tc.days))
	}
}

func Test_businessCalendar_daysBetween(t *testing.T) {
	calendar := newBusinessCalendar()

	for _, tc := range []struct {
		from     date.Date
		to       date.Date
		expected int
	}{
		{
			date.New(2020, time.January, 1),
			date.New(2020, time.January, 2),
			1,
		},
		{
			date.New(2020, time.August, 28),
			date.New(2020, time.September, 1),
			2,
		},
		{
			date.New(2020, time.August, 29),
			date.New(2020, time.August, 30),
			0,
		},
		{
			date.New(2020, time.August, 28),
			date.New(2020, time.August, 29),
			0,
		},
		{
			date.New(2020, time.August, 30),
			date.New(2020, time.August, 31),
			1,
		},
		{
			date.New(2020, time.January, 1),
			date.New(2020, time.January, 10),
			7,
		},
		{
			date.New(2020, time.February, 1),
			date.New(2020, time.March, 1),
			20,
		},
	} {
		assert.Equal(t, tc.expected, calendar.daysBetween(tc.from, tc.to))
		assert.Equal(t, -tc.expected, calendar.daysBetween(tc.to, tc.from))
	}
}

func Test_businessCalendar_DaysBetween_AcrossWeekend(t *testing.T) {
	// The start date is a Friday. Shift the current date
	// week by week until the end date is reached one year after.
	// Notice that the DaysBetween function includes the arrival
	// business day but excludes the starting one:
	// - Friday to Saturday/Sunday: 0 days;
	// - Friday to Monday: 1 day;
	// - Saturday/Sunday to Monday: 1 day.

	var (
		start    = date.New(2021, time.October, 1)
		end      = start.AddDate(1, 0, 0)
		calendar = newBusinessCalendar()
	)

	for friday := start; !friday.After(end); friday = friday.Add(7) {
		// Get the days across the weekend and check they're consistent.
		saturday := friday.Add(1)
		sunday := friday.Add(2)
		monday := calendar.add(friday, 1)
		assert.Equal(t, friday.Add(3), monday)

		// There's no business days from a business day to a non-business day.
		assert.Equal(t, 0, calendar.daysBetween(friday, saturday))
		assert.Equal(t, 0, calendar.daysBetween(friday, sunday))
		assert.Equal(t, 1, calendar.daysBetween(friday, monday))

		// There is one business day from a non-business day to a business day.
		assert.Equal(t, 1, calendar.daysBetween(saturday, monday))
		assert.Equal(t, 1, calendar.daysBetween(sunday, monday))
	}
}

func Test_businessCalendar_isWeekend(t *testing.T) {
	calendar := newBusinessCalendar()

	// Dates range from Friday to Monday.
	assert.False(t, calendar.isWeekend(date.New(2021, time.September, 24)))
	assert.True(t, calendar.isWeekend(date.New(2021, time.September, 25)))
	assert.True(t, calendar.isWeekend(date.New(2021, time.September, 26)))
	assert.False(t, calendar.isWeekend(date.New(2021, time.September, 27)))
}

func Test_businessCalendar_isBusinessDay(t *testing.T) {
	calendar := newBusinessCalendar()

	// Dates range from Friday to Monday.
	assert.True(t, calendar.isBusinessDay(date.New(2021, time.September, 24)))
	assert.False(t, calendar.isBusinessDay(date.New(2021, time.September, 25)))
	assert.False(t, calendar.isBusinessDay(date.New(2021, time.September, 26)))
	assert.True(t, calendar.isBusinessDay(date.New(2021, time.September, 27)))
}

func Test_businessCalendar_closestBusinessDay(t *testing.T) {
	calendar := newBusinessCalendar()

	for _, tc := range []struct {
		origin   date.Date
		forwards bool
		expected date.Date
	}{
		{
			date.New(2021, time.September, 24),
			true,
			date.New(2021, time.September, 24),
		},
		{
			date.New(2021, time.September, 24),
			false,
			date.New(2021, time.September, 24),
		},
		{
			date.New(2021, time.September, 25),
			true,
			date.New(2021, time.September, 27),
		},
		{
			date.New(2021, time.September, 25),
			false,
			date.New(2021, time.September, 24),
		},
		{
			date.New(2021, time.September, 26),
			true,
			date.New(2021, time.September, 27),
		},
		{
			date.New(2021, time.September, 26),
			false,
			date.New(2021, time.September, 24),
		},
	} {
		assert.Equal(t, tc.expected, calendar.closestBusinessDay(tc.origin, tc.forwards))
	}
}

func Test_businessCalendar_ConsistencyChecks(t *testing.T) {
	calendar := newBusinessCalendar()

	t.Run("complete week", func(t *testing.T) {
		origin := date.New(2017, time.January, 9)
		for j := 0; j < 7; j++ {
			current := origin.Add(j)

			for i := 1; i <= 256; i++ {
				to := calendar.add(current, i)
				from := calendar.add(to, -i)

				// Forward-shifted date by i days should not equal the current date.
				assert.False(t, current.Equal(to))

				// Ensure between current/from and to there are exactly i business days.
				assert.Equal(t, i, calendar.daysBetween(current, to))
				assert.Equal(t, i, calendar.daysBetween(from, to))

				// Ensure forward and backward shifts are consistent.
				assert.True(t, calendar.add(from, i).Equal(to))
				assert.True(t, calendar.add(to, -i).Equal(from))
			}
		}
	})

	t.Run("zero shift", func(t *testing.T) {
		var (
			thursday = date.New(2018, time.August, 30)
			friday   = thursday.Add(1)
			saturday = thursday.Add(2)
			sunday   = thursday.Add(3)
			calendar = newBusinessCalendar()
		)

		assert.True(t, calendar.add(thursday, 0).Equal(thursday))
		assert.True(t, calendar.add(friday, 0).Equal(friday))
		assert.True(t, calendar.add(saturday, 0).Equal(friday))
		assert.True(t, calendar.add(sunday, 0).Equal(friday))
	})

	t.Run("shift across weekend", func(t *testing.T) {
		var (
			from                = date.New(2018, time.August, 2)
			businessBeforeTo    = date.New(2018, time.August, 31)
			to                  = date.New(2018, time.September, 1)
			businessDaysBetween = 21
			calendar            = newBusinessCalendar()
		)

		assert.Equal(t, calendar.daysBetween(from, to), businessDaysBetween)
		assert.True(t, calendar.add(to, -businessDaysBetween).Equal(from))
		assert.True(t, calendar.previousBusinessDay(to).Equal(businessBeforeTo))
		assert.True(t, calendar.add(from, businessDaysBetween).Equal(businessBeforeTo))
	})

	t.Run("all shifts", func(t *testing.T) {
		origin := date.New(2017, time.January, 9)
		for i := 1; i <= 256; i++ {
			to := calendar.add(origin, i)
			from := calendar.add(to, -i)

			assert.False(t, origin.Equal(to))
			assert.Equal(t, i, calendar.daysBetween(origin, to))
			assert.Equal(t, i, calendar.daysBetween(from, to))
			assert.True(t, calendar.add(from, i).Equal(to))
			assert.True(t, calendar.add(to, -i).Equal(from))
		}
	})

	t.Run("closest business day", func(t *testing.T) {
		// The end date is set 1024 business days (public holidays excluded)
		// after start in this test.

		var (
			start = date.New(2017, time.January, 9)
			end   = date.New(2020, time.December, 11)
		)

		for current := start; !current.After(end); current = current.Add(1) {
			from := calendar.previousBusinessDay(current)
			to := calendar.nextBusinessDay(current)

			dayDiff := 1
			if calendar.isBusinessDay(current) {
				dayDiff = 0
			}

			// Current is a business day ==> both from and to should coincide with it.
			assert.True(t, !calendar.isBusinessDay(current) || (from.Equal(current) && to.Equal(current)))
			// When current is not a business day, it is assumed equal to the latest
			// business day before it, so it always coincides with from.
			assert.Equal(t, 0, calendar.daysBetween(from, current))
			// Daily difference between from/current and to is equivalent.
			assert.Equal(t, dayDiff, calendar.daysBetween(current, to))
			assert.Equal(t, dayDiff, calendar.daysBetween(from, to))
		}
	})
}

func BenchmarkShiftBusinessDays(b *testing.B) {
	calendar := newBusinessCalendar()
	startingDate := date.New(2017, time.November, 9)

	for _, bc := range []struct {
		name  string
		shift int
	}{
		{
			"add 17 days",
			17,
		},
		{
			"remove 17 days",
			-17,
		},
		{
			"add 7778 days",
			17,
		},
		{
			"remove 7778 days",
			17,
		},
	} {
		bc := bc
		b.Run(bc.name, func(b *testing.B) {
			for n := 0; n < b.N; n++ {
				_ = calendar.add(startingDate, bc.shift)
			}
		})
	}
}

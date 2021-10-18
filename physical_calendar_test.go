package calendar

import (
	"fmt"
	"testing"
	"time"

	"github.com/fxtlabs/date"
	"github.com/stretchr/testify/assert"
)

func Test_physicalCalendar_isActive(t *testing.T) {
	calendar := newPhysicalCalendar()

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
			true,
		},
		{
			date.New(2020, time.January, 5),
			true,
		},
		{
			date.New(2020, time.January, 6),
			true,
		},
	} {
		assert.Equal(t, tc.expected, calendar.isActive(tc.date))
	}
}

func Test_physicalCalendar_daysInYear(t *testing.T) {
	assert.Equal(t, 365, newPhysicalCalendar().daysInYear())
}

func Test_physicalCalendar_add(t *testing.T) {
	calendar := newPhysicalCalendar()

	for _, tc := range []struct {
		origin   date.Date
		days     int
		expected date.Date
	}{
		{
			date.New(2020, time.September, 1),
			3,
			date.New(2020, time.September, 4),
		},
		{
			date.New(2020, time.September, 1),
			-3,
			date.New(2020, time.August, 29),
		},
		{
			date.New(2020, time.September, 1),
			4,
			date.New(2020, time.September, 5),
		},
	} {
		tc := tc
		t.Run(fmt.Sprintf("from %s plus %d", tc.origin, tc.days), func(t *testing.T) {
			assert.Equal(t, tc.expected, calendar.add(tc.origin, tc.days))
		})
	}
}

func Test_physicalCalendar_daysBetween(t *testing.T) {
	calendar := newPhysicalCalendar()

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
			date.New(2020, time.January, 1),
			date.New(2020, time.January, 10),
			9,
		},
		{
			date.New(2020, time.February, 1),
			date.New(2020, time.March, 1),
			29,
		},
	} {
		assert.Equal(t, tc.expected, calendar.daysBetween(tc.from, tc.to))
	}
}

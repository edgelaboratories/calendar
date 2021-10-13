package calendar

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func Test_newDayCounter(t *testing.T) {
	t.Run("business", func(t *testing.T) {
		_, ok := newDayCounter(BusinessDays).(*businessCalendar)
		assert.True(t, ok)
	})

	t.Run("physical", func(t *testing.T) {
		_, ok := newDayCounter(CalendarDays).(*physicalCalendar)
		assert.True(t, ok)
	})

	t.Run("unspecified", func(t *testing.T) {
		_, ok := newDayCounter(Convention("Unspecified")).(*businessCalendar)
		assert.True(t, ok)
	})
}

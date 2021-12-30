package githubasana

import (
	"time"
)

// NextBusinessDay returns next n business day.
// If n is zero, it returns most recent business day including its day.
func NextBusinessDay(n int, base time.Time, holidays map[string]bool) time.Time {
	i := 0

	// move to first business day.
	for {
		if isHoliday(base, holidays) {
			base = base.AddDate(0, 0, 1)
			i = 1
		} else {
			break
		}
	}

	// start to count.
	for {
		if !isHoliday(base, holidays) {
			i++
		}

		if n < i {
			return base
		}

		base = base.AddDate(0, 0, 1)
	}
}

func isHoliday(day time.Time, holidays map[string]bool) bool {
	if v, ok := holidays[day.Format("2006-01-02")]; ok {
		return v
	}

	// Normally we do not work on the weekends, right? do we?
	return day.Weekday() == time.Saturday || day.Weekday() == time.Sunday
}

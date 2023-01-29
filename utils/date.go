package utils

import (
	"time"
)

type Date struct{}

var DateUtil Date

const (
	MONTH_LAYOUT = "2006-01"
	DATE_LAYOUT  = "2006-01-02"
)

func (d Date) TransferToTimeFormat(date string) (time.Time, error) {
	time, err := time.Parse(DATE_LAYOUT, date)
	return time, err
}

func (d Date) SliceByWeek(start, end time.Time) (res [][]time.Time) {
	week := []time.Time{}

	for start.Before(end.Add(time.Second)) {

		// Fisrt day in a week
		if len(week) == 0 {
			week = append(week, start)
		}

		// Last day in a week
		if start.Weekday() == time.Sunday && !start.Equal(week[0]) {
			week = append(week, start)
			res = append(res, week)
			week = []time.Time{}
		}

		start = start.AddDate(0, 0, 1)
	}

	if len(week) > 0 {
		week = append(week, start.Add(-time.Second))
		res = append(res, week)
	}

	return
}

func (d Date) DaysInMonth(month time.Month) int {
	return time.Date(0, month, 0, 0, 0, 0, 0, time.UTC).Day()
}

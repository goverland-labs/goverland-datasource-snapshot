package helpers

import (
	"github.com/samber/lo"
)

type Interval struct {
	From  int
	Limit int
}

func GenerateIntervals(from, limit int, skip []int) []Interval {
	intervals := make([]Interval, 0, len(skip)+1)

	shift := 0
	for _, index := range skip {
		intervals = append(intervals, Interval{
			From:  from + shift,
			Limit: index - shift - 1,
		})

		shift = index
	}

	intervals = append(intervals, Interval{
		From:  from + shift,
		Limit: limit - shift,
	})

	intervals = lo.Filter(intervals, func(item Interval, index int) bool {
		return item.Limit > 0
	})

	return intervals
}

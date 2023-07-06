package helpers

import (
	"testing"

	"github.com/smartystreets/goconvey/convey"
)

func TestGenerateIntervals(t *testing.T) {
	convey.Convey("generating intervals with skipped elements", t, func() {
		intervals := GenerateIntervals(5, 10, []int{2, 3, 5, 9})

		convey.So(intervals, convey.ShouldHaveLength, 4)
		convey.So(intervals[0].From, convey.ShouldEqual, 5)
		convey.So(intervals[0].Limit, convey.ShouldEqual, 1)
		convey.So(intervals[1].From, convey.ShouldEqual, 8)
		convey.So(intervals[1].Limit, convey.ShouldEqual, 1)
		convey.So(intervals[2].From, convey.ShouldEqual, 10)
		convey.So(intervals[2].Limit, convey.ShouldEqual, 3)
		convey.So(intervals[3].From, convey.ShouldEqual, 14)
		convey.So(intervals[3].Limit, convey.ShouldEqual, 1)
	})
}

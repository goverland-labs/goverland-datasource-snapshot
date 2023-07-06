package helpers

import (
	"testing"

	"github.com/smartystreets/goconvey/convey"
)

func TestZeroIfNil(t *testing.T) {
	convey.Convey("test unwrap zero value for nil pointer", t, func() {
		convey.So(ZeroIfNil(Ptr(300)), convey.ShouldEqual, 300)
		convey.So(ZeroIfNil(Ptr("value")), convey.ShouldEqual, "value")

		var a *string
		convey.So(ZeroIfNil(a), convey.ShouldEqual, "")
	})
}

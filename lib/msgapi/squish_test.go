package msgapi

import (
	. "github.com/smartystreets/goconvey/convey"
	"testing"
)

func TestSquishBufHash32(t *testing.T) {
	Convey("Check Squish bufHash32()", t, func() {
		Convey("calculates correct Sqush bufhash-32 of an empty string", func() {
			So(0x0, ShouldEqual, bufHash32(""))
		})
		Convey("calculates correct Sqush bufhash-32 of the string 'Alexander N. Skovpen'", func() {
			So(0x00efd7be, ShouldEqual, bufHash32("Alexander N. Skovpen"))
		})
	})
}

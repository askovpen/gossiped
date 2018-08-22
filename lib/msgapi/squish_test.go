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
		Convey("calculates correct Sqush bufhash-32 of the string 'Эдуардыч'[CP866]", func() {
			So(0x02debb97, ShouldEqual, bufHash32(string([]byte{'\x9D','\xA4', '\xE3', '\xA0', '\xE0', '\xA4', '\xEB', '\xE7'})))
		})
	})
}

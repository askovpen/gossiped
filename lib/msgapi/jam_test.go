package msgapi

import (
	. "github.com/smartystreets/goconvey/convey"
	"testing"
)

func TestJamCrc32r(t *testing.T) {
	Convey("Check Jam crc32r()", t, func() {
		Convey("calculates correct JAM CRC-32 of an empty string", func() {
			So(0xffffffff, ShouldEqual, crc32r(""))
		})
		Convey("calculates correct JAM CRC-32 of the string 'Alexander N. Skovpen'", func() {
			So(0x30222bd1, ShouldEqual, crc32r("Alexander N. Skovpen"))
		})
	})
}

package utils

import (
	"bytes"
	"errors"
	"testing"

	. "github.com/franela/goblin"
)

type TS struct {
	A uint8
	B string
	C [3]byte
}

func TestStructFromBuffer(t *testing.T) {
	g := Goblin(t)
	g.Describe("Check WriteStructToBuffer()", func() {
		g.It("Check uint", func() {
			buf := new(bytes.Buffer)
			testStruct := TS{A: 254}
			err := WriteStructToBuffer(buf, &testStruct)
			g.Assert(err).Equal(nil)
			g.Assert(buf.Bytes()).Equal([]byte{254, 0, 0, 0})
		})
		g.It("Check string", func() {
			buf := new(bytes.Buffer)
			testStruct := TS{B: "test"}
			err := WriteStructToBuffer(buf, &testStruct)
			g.Assert(err).Equal(nil)
			g.Assert(buf.Bytes()).Equal([]byte{0, 0x74, 0x65, 0x73, 0x74, 0, 0, 0})
		})
		g.It("Check array", func() {
			buf := new(bytes.Buffer)
			testStruct := TS{C: [3]byte{0, 128, 255}}
			err := WriteStructToBuffer(buf, &testStruct)
			g.Assert(err).Equal(nil)
			g.Assert(buf.Bytes()).Equal([]byte{0, 0, 128, 255})
		})
		g.It("Check invalid", func() {
			buf := new(bytes.Buffer)
			a := 1
			err := WriteStructToBuffer(buf, &a)
			g.Assert(err).Equal(errors.New("invaild type Not a struct"))
		})
	})
	g.Describe("Check ReadStructFromBuffer()", func() {
		g.It("Check uint", func() {
			buf := bytes.NewBuffer([]byte{254, 0, 0, 0, 0, 0, 0, 0})
			var testStruct TS
			err := ReadStructFromBuffer(buf, &testStruct)
			g.Assert(err).Equal(nil)
			g.Assert(testStruct).Equal(TS{A: 254, B: "\x00", C: [3]byte{0}})
		})
		g.It("Check string", func() {
			buf := bytes.NewBuffer([]byte{0, 0x74, 0x65, 0x73, 0x74, 0, 0, 0, 0})
			var testStruct TS
			err := ReadStructFromBuffer(buf, &testStruct)
			g.Assert(err).Equal(nil)
			g.Assert(testStruct).Equal(TS{A: 0, B: "test\x00", C: [3]byte{0}})
		})
		g.It("Check array", func() {
			buf := bytes.NewBuffer([]byte{0, 0, 0, 128, 255})
			var testStruct TS
			err := ReadStructFromBuffer(buf, &testStruct)
			g.Assert(err).Equal(nil)
			g.Assert(testStruct).Equal(TS{A: 0, B: "\x00", C: [3]byte{0, 128, 255}})
		})
		g.It("Check invalid", func() {
			buf := bytes.NewBuffer([]byte{0, 0, 128, 255})
			var testStruct TS
			err := ReadStructFromBuffer(buf, &testStruct)
			g.Assert(err).Equal(errors.New("unexpected EOF"))
		})
	})
}

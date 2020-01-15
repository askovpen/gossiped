package utils

import (
	"bytes"
	"errors"
	. "github.com/franela/goblin"
	"testing"
)

type TS struct {
	A uint8
	B string
	C [3]byte
}
func TestReadStructFromBuffer(t *testing.T) {
	g := Goblin(t)
	g.Describe("Check WriteStructToBuffer()", func() {
		g.It("Check uint", func() {
			buf := new(bytes.Buffer)
			testStruct:=TS{A:254}
			err := WriteStructToBuffer(buf, &testStruct)
			g.Assert(err).Equal(nil)
			g.Assert(buf.Bytes()).Equal([]byte{254,0,0,0})
		})
		g.It("Check string", func() {
			buf := new(bytes.Buffer)
			testStruct:=TS{B:"test"}
			err := WriteStructToBuffer(buf, &testStruct)
			g.Assert(err).Equal(nil)
			g.Assert(buf.Bytes()).Equal([]byte{0,0x74, 0x65, 0x73, 0x74,0,0,0})
		})
		g.It("Check array", func() {
			buf := new(bytes.Buffer)
			testStruct:=TS{C:[3]byte{0,128,255}}
			err := WriteStructToBuffer(buf, &testStruct)
			g.Assert(err).Equal(nil)
			g.Assert(buf.Bytes()).Equal([]byte{0,0,128,255})
		})
		g.It("Check invalid", func() {
			buf := new(bytes.Buffer)
			a:=1
			err := WriteStructToBuffer(buf, &a)
			g.Assert(err).Equal(errors.New("invaild type Not a struct"))
		})
	})
}

package msgapi

import (
	"github.com/askovpen/gossiped/pkg/types"
	. "github.com/franela/goblin"
	"os"
	"testing"
	"time"
	"unsafe"
)

func TestJamCrc32r(t *testing.T) {
	g := Goblin(t)
	g.Describe("Check Jam crc32r()", func() {
		g.It("calculates correct JAM CRC-32 of an empty string", func() {
			g.Assert(crc32r("")).Equal(uint32(0xffffffff))
		})
		g.It("calculates correct JAM CRC-32 of the string 'Alexander N. Skovpen'", func() {
			g.Assert(crc32r("Alexander N. Skovpen")).Equal(uint32(0x30222bd1))
		})
	})
}
func TestJam(t *testing.T) {
	Area := &JAM{
		AreaPath: "../../testdata/jamtest",
		AreaName: "test",
		AreaType: EchoAreaTypeEcho,
	}
	Areas = Areas[:0]
	Areas = append(Areas, Area)
	g := Goblin(t)
	g.Describe("Check JAM read/write", func() {
		m := &Message{
			AreaObject:  &Areas[0],
			From:        "SysOp",
			To:          "SysOp",
			Subject:     "Test",
			FromAddr:    types.AddrFromNum(2, 5020, 9696, 1),
			ToAddr:      types.AddrFromNum(2, 5020, 9696, 2),
			DateWritten: time.Now(),
			DateArrived: time.Now(),
			Body:        "Test\nBody\n * Origin: test",
			Kludges:     make(map[string]string),
		}
		m.MakeBody()
		g.It("create msg", func() {
			g.Assert(len(*Area.GetMessages())).Equal(0)
			g.Assert(Area.SaveMsg(m)).Equal(nil)
		})
		g.It("add msg", func() {
			g.Assert(Area.SaveMsg(m)).Equal(nil)
		})
		g.It("check num msgs", func() {
			g.Assert(Area.GetCount()).Equal(uint32(2))
		})
		g.It("read msg", func() {
			nm, err := Area.GetMsg(1)
			g.Assert(err).Equal(nil)
			g.Assert(nm.FromAddr).Equal(types.AddrFromNum(2, 5020, 9696, 1))
		})
		g.It("get/set last", func() {
			Area.SetLast(1)
			g.Assert(Area.GetLast()).Equal(uint32(1))
			g.Assert(len(*Area.GetMessages())).Equal(2)
		})
		g.It("del msg", func() {
			g.Assert(Area.DelMsg(2)).Equal(nil)
		})
	})
	os.Remove("../../testdata/jamtest.jdt")
	os.Remove("../../testdata/jamtest.jdx")
	os.Remove("../../testdata/jamtest.jhr")
	os.Remove("../../testdata/jamtest.jlr")
}

func BenchmarkJamCrc32r(b *testing.B) {
	b.SetBytes(20)
	for n := 0; n < b.N; n++ {
		crc32r("Alexander N. Skovpen")
	}
}

func BenchmarkJamGetMessages(b *testing.B) {
	Area := &JAM{
		AreaPath: "../../testdata/jamtest",
		AreaName: "test",
		AreaType: EchoAreaTypeEcho,
	}
	Areas = Areas[:0]
	Areas = append(Areas, Area)
	for n := 0; n < b.N; n++ {
		m := &Message{
			AreaObject:  &Areas[0],
			From:        "SysOp",
			To:          "SysOp",
			Subject:     "Test",
			FromAddr:    types.AddrFromNum(2, 5020, 9696, 1),
			ToAddr:      types.AddrFromNum(2, 5020, 9696, 2),
			DateWritten: time.Now(),
			DateArrived: time.Now(),
			Body:        "Test\nBody\n * Origin: test",
			Kludges:     make(map[string]string),
		}
		m.MakeBody()
		Area.SaveMsg(m)
	}
	m := &Message{
		AreaObject:  &Areas[0],
		From:        "SysOp",
		To:          "SysOp",
		Subject:     "Test",
		FromAddr:    types.AddrFromNum(2, 5020, 9696, 1),
		ToAddr:      types.AddrFromNum(2, 5020, 9696, 2),
		DateWritten: time.Now(),
		DateArrived: time.Now(),
		Body:        "Test\nBody\n * Origin: test",
		Kludges:     make(map[string]string),
	}
	m.MakeBody()
	b.SetBytes(int64(unsafe.Sizeof(m)))
	b.ResetTimer()
	//for n := 0; n < b.N; n++ {
	//      Area.GetMsg(uint32(n))
	//}
	Area.GetMessages()
	b.StopTimer()
	os.Remove("../../testdata/jamtest.jdt")
	os.Remove("../../testdata/jamtest.jdx")
	os.Remove("../../testdata/jamtest.jhr")
	os.Remove("../../testdata/jamtest.jlr")

}

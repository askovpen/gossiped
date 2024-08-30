package msgapi

import (
	"github.com/askovpen/gossiped/pkg/types"
	. "github.com/franela/goblin"
	"os"
	"testing"
	"time"
	"unsafe"
)

func TestSquishBufHash32(t *testing.T) {
	g := Goblin(t)
	g.Describe("Check Squish bufHash32()", func() {
		g.It("calculates correct Sqush bufhash-32 of an empty string", func() {
			g.Assert(bufHash32("")).Equal(uint32(0))
		})
		g.It("calculates correct Sqush bufhash-32 of the string 'Alexander N. Skovpen'", func() {
			g.Assert(bufHash32("Alexander N. Skovpen")).Equal(uint32(0x00efd7be))
		})
		g.It("calculates correct Sqush bufhash-32 of the string 'Эдуардыч'[CP866]", func() {
			g.Assert(bufHash32(string([]byte{'\x9D', '\xA4', '\xE3', '\xA0', '\xE0', '\xA4', '\xEB', '\xE7'}))).Equal(uint32(0x02debb97))
		})
		g.It("calculates correct Sqush bufhash-32 of the string 'Тест1'[CP866]", func() {
			g.Assert(bufHash32(string([]byte{'\x92', '\xA5', '\xE1', '\xE2', '\x32'}))).Equal(uint32(0x009D3F52))
		})
		g.It("calculates correct Sqush bufhash-32 of the string 'Юрий Григорьев'[CP866]", func() {
			g.Assert(bufHash32(string([]byte{'\x9E', '\xE0', '\xA8', '\xA9', '\x20', '\x83', '\xE0', '\xA8', '\xA3', '\xAE', '\xE0', '\xEC', '\xA5', '\xA2'}))).Equal(uint32(0x7c100ff2))
		})
	})
}

func TestSquish(t *testing.T) {
	Area := &Squish{
		AreaPath: "../../testdata/sqtest",
		AreaName: "test",
		AreaType: EchoAreaTypeEcho,
	}
	Areas = Areas[:0]
	Areas = append(Areas, Area)
	g := Goblin(t)
	g.Describe("Check Squish read/write", func() {
		m := &Message{
			AreaObject:  &Areas[0],
			From:        "SysOp",
			To:          "SysOp",
			Subject:     "Test",
			FromAddr:    types.AddrFromNum(2, 5020, 9696, 1),
			ToAddr:      types.AddrFromNum(2, 5020, 9696, 2),
			DateWritten: time.Now(),
			DateArrived: time.Now(),
			Body:        "Test\nBody",
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
			g.Assert(Area.DelMsg(1)).Equal(nil)
			g.Assert(len(*Area.GetMessages())).Equal(1)
		})
	})
	os.Remove("../../testdata/sqtest.sqd")
	os.Remove("../../testdata/sqtest.sqi")
	os.Remove("../../testdata/sqtest.sql")
}

func BenchmarkSquishBufHash32(b *testing.B) {
	b.SetBytes(20)
	for n := 0; n < b.N; n++ {
		bufHash32("Alexander N. Skovpen")
	}
}

func BenchmarkSquishGetMessages(b *testing.B) {
	Area := &Squish{
		AreaPath: "../../testdata/sqtest",
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
			Body:        "Test\nBody",
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
		Body:        "Test\nBody",
		Kludges:     make(map[string]string),
	}
	m.MakeBody()
	b.SetBytes(int64(unsafe.Sizeof(m)))
	b.ResetTimer()
	//for n := 0; n < b.N; n++ {
	//	Area.GetMsg(uint32(n))
	//}
	Area.GetMessages()
	b.StopTimer()
	os.Remove("../../testdata/sqtest.sqd")
	os.Remove("../../testdata/sqtest.sqi")
}

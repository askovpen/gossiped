package msgapi

import (
	"github.com/askovpen/goated/lib/types"
	. "github.com/franela/goblin"
	"os"
	"testing"
	"time"
)

func TestMSG(t *testing.T) {
	Area := &MSG{
		AreaPath: "../../testdata/test",
		AreaName: "test",
		AreaType: EchoAreaTypeNetmail,
	}
	Areas = append(Areas, Area)
	g := Goblin(t)
	g.Describe("Check MSG read/write", func() {
		m := &Message{
			AreaID:      0,
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
			g.Assert(Area.SaveMsg(m)).Equal(nil)
		})
		g.It("add msg", func() {
			g.Assert(Area.SaveMsg(m)).Equal(nil)
		})
	})
	os.RemoveAll("../../testdata/test")
}

package msgapi

import (
	"cmp"
	"slices"
	"strings"

	"github.com/askovpen/gossiped/pkg/config"
)

// EchoAreaMsgType Area msg base type
type EchoAreaMsgType string

// EchoAreaType Area type
type EchoAreaType uint8

const (
	AreasSortingDefault = "default"
	AreasSortingUnread  = "unread"
)

var (
	validAreaSortModes = map[string]bool{
		AreasSortingDefault: true,
		AreasSortingUnread:  true,
	}
	// Areas list
	Areas []AreaPrimitive
)

// types
const (
	EchoAreaMsgTypeJAM        EchoAreaMsgType = "JAM"
	EchoAreaMsgTypeMSG        EchoAreaMsgType = "MSG"
	EchoAreaMsgTypeSquish     EchoAreaMsgType = "Squish"
	EchoAreaMsgTypePasstrough EchoAreaMsgType = "Passtrough"
	EchoAreaTypeNetmail       EchoAreaType    = 0
	EchoAreaTypeEcho          EchoAreaType    = 3
	EchoAreaTypeLocal         EchoAreaType    = 4
	EchoAreaTypeDupe          EchoAreaType    = 2
	EchoAreaTypeBad           EchoAreaType    = 1
	EchoAreaTypeNone          EchoAreaType    = 5
)

// AreaPrimitive interface
type AreaPrimitive interface {
	Init()
	GetCount() uint32
	GetLast() uint32
	GetMsg(position uint32) (*Message, error)
	GetName() string
	GetMsgType() EchoAreaMsgType
	GetType() EchoAreaType
	SetChrs(string)
	GetChrs() string
	SetLast(uint32)
	DelMsg(uint32) error
	SaveMsg(*Message) error
	GetMessages() *[]MessageListItem
}

func AreaHasUnreadMessages(area *AreaPrimitive) bool {
	return (*area).GetCount()-(*area).GetLast() > 0
}

func SortAreas() {
	var configMode = AreasSortingDefault
	var configValue, _ = config.Config.Sorting["areas"]
	var match, okMode = validAreaSortModes[configValue]
	if okMode && match {
		configMode = configValue
	}
	slices.SortFunc(Areas, func(a AreaPrimitive, b AreaPrimitive) int {
		var n = 0
		if configMode == AreasSortingUnread {
			var aUnread = a.GetCount() - a.GetLast()
			var bUnread = b.GetCount() - b.GetLast()
			if n = cmp.Compare(bUnread, aUnread); n != 0 {
				return n
			}
		}
		if n = cmp.Compare(a.GetType(), b.GetType()); n != 0 {
			return n
		}
		n = strings.Compare(a.GetName(), b.GetName())
		return n
	})
}

// Lookup name->id
func Lookup(name string) int {
	for i, a := range Areas {
		if a.GetName() == name {
			return i
		}
	}
	return 0
}

// Search part name->id
func Search(name string) int {
	for i, a := range Areas {
		if strings.Contains(strings.ToLower(a.GetName()), strings.ToLower(name)) {
			return i + 1
		}
	}
	return 0
}

package msgapi

import (
	"strings"
)

// EchoAreaMsgType Area msg base type
type EchoAreaMsgType string

// EchoAreaType Area type
type EchoAreaType uint8

var (
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

// Lookup name->id
func Lookup(name string) int {
	for i, a := range Areas {
		if a.GetName() == name {
			return i
		}
	}
	return 0
}

func Search(name string) int {
	for i, a := range Areas {
		if strings.Contains(strings.ToLower(a.GetName()), strings.ToLower(name)) {
			return i + 1
		}
	}
	return 0
}

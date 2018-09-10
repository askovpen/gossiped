package msgapi

// EchoAreaMsgType Area msg base type
type EchoAreaMsgType string

// EchoAreaType Area type
type EchoAreaType uint8

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
	SaveMsg(*Message) error
	GetMessages() *[]MessageListItem
}

package msgapi

type EchoAreaMsgType string
type EchoAreaType uint8

const (
	EchoAreaMsgTypeJAM    EchoAreaMsgType = "JAM"
	EchoAreaMsgTypeMSG    EchoAreaMsgType = "MSG"
	EchoAreaMsgTypeSquish EchoAreaMsgType = "Squish"

	EchoAreaTypeNetmail EchoAreaType = 0
	EchoAreaTypeEcho    EchoAreaType = 3
	EchoAreaTypeLocal   EchoAreaType = 4
	EchoAreaTypeDupe    EchoAreaType = 2
	EchoAreaTypeBad     EchoAreaType = 1
)

type AreaPrimitive interface {
	Init()
	GetCount() uint32
	GetLast() uint32
	GetMsg(position uint32) (*Message, error)
	GetName() string
	GetMsgType() EchoAreaMsgType
	GetType() EchoAreaType
	SetLast(uint32)
	SaveMsg(*Message) error
}

type Area_s struct {
	AreaName string
	FileName string
	MsgBType EchoAreaMsgType
	AreaType EchoAreaType
	Count    uint32
	Position uint32
}

package msgapi

type EchoAreaMsgType string
type EchoAreaType string

const (
	EchoAreaMsgTypeJAM    EchoAreaMsgType = "JAM"
	EchoAreaMsgTypeMSG    EchoAreaMsgType = "MSG"
	EchoAreaMsgTypeSquish EchoAreaMsgType = "Squish"

	EchoAreaTypeNetmail EchoAreaType = "Netmail"
	EchoAreaTypeEcho    EchoAreaType = "Echo"
	EchoAreaTypeLocal   EchoAreaType = "Local"
	EchoAreaTypeDupe    EchoAreaType = "Dupe"
	EchoAreaTypeBad     EchoAreaType = "Bad"
)

type AreaPrimitive interface {
	Init()
	GetCount() uint32
	GetLast() uint32
	GetMsg(position uint32) (*Message, error)
	GetName() string
	GetType() EchoAreaMsgType
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

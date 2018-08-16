package msgapi

type EchoAreaType string

const (
	EchoAreaTypeJAM    EchoAreaType = "JAM"
	EchoAreaTypeMSG    EchoAreaType = "MSG"
	EchoAreaTypeSquish EchoAreaType = "Squish"
)

type AreaPrimitive interface {
	Init()
	GetCount() uint32
	GetLast() uint32
	GetMsg(position uint32) (*Message, error)
	GetName() string
	GetType() EchoAreaType
	SetLast(uint32)
	SaveMsg(*Message) error
}

type Area_s struct {
	AreaName string
	FileName string
	MsgBType EchoAreaType
	Count    uint32
	Position uint32
}

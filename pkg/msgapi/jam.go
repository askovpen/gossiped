package msgapi

import (
	"bufio"
	"bytes"
	"encoding/binary"
	"errors"

	"github.com/askovpen/gossiped/pkg/config"
	"github.com/askovpen/gossiped/pkg/types"
	"github.com/askovpen/gossiped/pkg/utils"

	//"hash/crc32"
	"io"

	//"log"
	"os"
	// "sort"
	//"strings"
	"time"
)

// JAM struct
type JAM struct {
	AreaPath, AreaName string
	AreaType           EchoAreaType
	Chrs               string
	indexStructure     []jamS
	lastRead           []jamL
	messages           []MessageListItem
	headerStructure    jhrS
}

type jamS struct {
	MessageNum uint32
	jamsh      jamSH
}

type jamSH struct {
	ToCRC, Offset uint32
}

type jhrS struct {
	Signature                           uint32
	DateCreated, ModCounter, ActiveMsgs uint32
	PasswordCRC, BaseMsgNum, Highwater  uint32
	RSRVD                               [996]byte
}
type jamH struct {
	Signature                               uint32
	Revision, ReservedWord                  uint16
	SubfieldLen, TimesRead, MSGIDcrc        uint32
	REPLYcrc, ReplyTo, Reply1st             uint32
	ReplyNext, DateWritten, DateReceived    uint32
	DateProcessed, MessageNumber, Attribute uint32
	Attribute2, Offset, TxtLen              uint32
	PasswordCRC, Cost                       uint32
}

type jamL struct {
	UserCRC, UserID, LastReadMsg, HighReadMsg uint32
}

func (j *JAM) getAttrs(a uint32) (attrs []string) {
	datr := []string{
		"Loc", "", "Pvt", "Rcv",
		"Snt", "", "", "",
		"", "", "", "",
		"", "", "", "",
		"", "", "", "",
		"", "", "", "",
		"", "", "", "",
		"", "", "", "[red]Del[silver]",
	}
	i := 0
	for a > 0 {
		if a&1 > 0 {
			if datr[i] != "" {
				attrs = append(attrs, datr[i])
			}
		}
		i++
		a >>= 1
	}
	return
}

func (j *JAM) getOffsetByNum(num uint32) (offset uint32) {
	for i, is := range j.indexStructure {
		if is.MessageNum == num {
			return uint32(i) + 1
		}
	}
	return 0
}

// GetMsg return msg
func (j *JAM) GetMsg(position uint32) (*Message, error) {
	if len(j.indexStructure) == 0 {
		//		return nil, errors.New("Empty Area")
		return nil, nil
	}
	if position == 0 {
		position = 1
	}
	fJhr, err := os.Open(j.AreaPath + ".jhr")
	if err != nil {
		return nil, err
	}
	defer fJhr.Close()
	_, err = fJhr.Seek(int64(j.indexStructure[position-1].jamsh.Offset), 0)
	if err != nil {
		return nil, err
	}
	header := make([]byte, 76)
	fJhr.Read(header)
	headerb := bytes.NewBuffer(header)
	var jamh jamH
	if err = utils.ReadStructFromBuffer(headerb, &jamh); err != nil {
		return nil, err
	}
	if jamh.Signature != 0x4d414a {
		return nil, errors.New("wrong message signature")
	}
	rm := &Message{
		Area:        j.AreaName,
		MsgNum:      position,
		MaxNum:      uint32(len(j.indexStructure)),
		DateWritten: time.Unix(int64(jamh.DateWritten), 0)}
	_, tofs := rm.DateWritten.Zone()
	if jamh.DateReceived > 0 {
		rm.DateArrived = time.Unix(int64(jamh.DateReceived), 0)
	} else {
		rm.DateArrived = time.Unix(int64(jamh.DateProcessed), 0)
	}
	rm.DateWritten = rm.DateWritten.Add(time.Duration(tofs) * -time.Second)
	rm.DateArrived = rm.DateArrived.Add(time.Duration(tofs) * -time.Second)
	rm.Attrs = j.getAttrs(jamh.Attribute)
	if jamh.ReplyTo > 0 {
		rm.ReplyTo = j.getOffsetByNum(jamh.ReplyTo)
	} else {
		jamh.ReplyTo = 0
	}
	if jamh.Reply1st > 0 {
		rm.Replies = append(rm.Replies, j.getOffsetByNum(jamh.Reply1st))
	}
	//if jamh.ReplyNext>0 {
	//  rm.Replies = append(rm.Replies,j.getOffsetByNum(jamh.ReplyNext))
	//}
	deleted := false
	if jamh.Attribute&0x80000000 > 0 {
		deleted = true
	}
	rm.Body += ""
	kl := make([]byte, jamh.SubfieldLen)
	fJhr.Read(kl)
	klb := bytes.NewBuffer(kl)
	afterBody := ""
	for {
		var LoID, HiID uint16
		var datLen uint32
		err = binary.Read(klb, binary.LittleEndian, &LoID)
		if err != nil {
			break
		}
		binary.Read(klb, binary.LittleEndian, &HiID)
		binary.Read(klb, binary.LittleEndian, &datLen)
		if datLen > 80 {
			datLen = 80
		}
		val := make([]byte, datLen)
		binary.Read(klb, binary.LittleEndian, &val)
		switch LoID {
		case 0:
			fr := types.AddrFromString(string(val))
			if fr != nil {
				rm.FromAddr = fr
			}
		case 1:
			if j.AreaType != EchoAreaTypeLocal && j.AreaType != EchoAreaTypeEcho {
				rm.ToAddr = types.AddrFromString(string(val))
			}
		case 2:
			rm.From = string(val)
		case 3:
			if !deleted {
				if crc32r(string(val)) != j.indexStructure[position-1].jamsh.ToCRC {
					rm.Corrupted = true
				}
			}
			rm.To = string(val)
		case 4:
			if crc32r(string(val)) != jamh.MSGIDcrc {
				rm.Corrupted = true
			}
			rm.Body += "\x01MSGID: " + string(val) + "\x0d"
		case 5:
			if crc32r(string(val)) != jamh.REPLYcrc {
				rm.Corrupted = true
			}
			rm.Body += "\x01REPLYID: " + string(val) + "\x0d"
		case 6:
			rm.Subject = string(val)
		case 7:
			rm.Body += "\x01PID: " + string(val) + "\x0d"
		case 8:
			afterBody += "\x01Via " + string(val) + "\x0d"
		case 2004:
			rm.Body += "\x01TZUTC: " + string(val) + "\x0d"
		case 2000:
			rm.Body += "\x01" + string(val) + "\x0d"
		case 2001:
			afterBody += "SEEN-BY: " + string(val) + "\x0d"
		case 2002:
			afterBody += "\x01PATH: " + string(val) + "\x0d"
		}
	}
	fJdt, err := os.Open(j.AreaPath + ".jdt")
	if err != nil {
		return nil, err
	}
	defer fJdt.Close()
	fJdt.Seek(int64(jamh.Offset), 0)
	txt := make([]byte, jamh.TxtLen)
	fJdt.Read(txt)
	rm.Body += string(txt)
	rm.Body += afterBody
	err = rm.ParseRaw()
	if err != nil {
		return nil, err
	}
	return rm, nil
}
func (j *JAM) readJDX() {
	if len(j.indexStructure) > 0 {
		return
	}
	fJhr, err := os.Open(j.AreaPath + ".jhr")
	if err != nil {
		return
	}
	defer fJhr.Close()
	header := make([]byte, 1024)
	fJhr.Read(header)
	headerb := bytes.NewBuffer(header)
	if err = utils.ReadStructFromBuffer(headerb, &j.headerStructure); err != nil {
		return
	}

	file, err := os.Open(j.AreaPath + ".jdx")
	if err != nil {
		return
	}
	defer file.Close()
	reader := bufio.NewReader(file)
	part := make([]byte, 16384)
	i := uint32(0)
	for {
		count, err := reader.Read(part)
		if err != nil {
			break
		}
		partb := bytes.NewBuffer(part[:count])
		for {
			var jam jamSH
			if err = utils.ReadStructFromBuffer(partb, &jam); err != nil {
				break
			}
			if jam.Offset != 0xffffffff {
				j.indexStructure = append(j.indexStructure, jamS{i + j.headerStructure.BaseMsgNum, jam})
			}
			i++
		}
	}
	// sort.Slice(j.indexStructure, func(a, b int) bool { return j.indexStructure[a].MessageNum < j.indexStructure[b].MessageNum })
}

func (j *JAM) readJLR() {
	if len(j.lastRead) > 0 {
		return
	}
	file, err := os.Open(j.AreaPath + ".jlr")
	if err != nil {
		return
	}
	defer file.Close()
	reader := bufio.NewReader(file)
	part := make([]byte, 16384)
	for {
		count, err := reader.Read(part)
		if err != nil {
			break
		}
		partb := bytes.NewBuffer(part[:count])
		for {
			var jaml jamL
			if err = utils.ReadStructFromBuffer(partb, &jaml); err != nil {
				break
			}
			j.lastRead = append(j.lastRead, jaml)
		}
	}
}
func (j *JAM) getPositionOfJamMsg(mID uint32) uint32 {
	for i, ji := range j.indexStructure {
		if mID == ji.MessageNum {
			return uint32(i)
		}
	}
	if mID != 0 {
		return uint32(len(j.indexStructure))
	}
	return 0
}

// GetLast return last message
func (j *JAM) GetLast() uint32 {
	j.readJDX()
	if len(j.indexStructure) == 0 {
		return 0
	}
	j.readJLR()
	for _, l := range j.lastRead {
		if l.UserCRC == crc32r(config.Config.Username) {
			if j.getPositionOfJamMsg(l.LastReadMsg)+1 > uint32(len(j.indexStructure)) {
				return uint32(len(j.indexStructure))
			}
			return j.getPositionOfJamMsg(l.LastReadMsg) + 1
		}
	}
	return 0
}

// GetCount return count messages
func (j *JAM) GetCount() uint32 {
	j.readJDX()
	return uint32(len(j.indexStructure))
}

// GetMsgType return msg base type
func (j *JAM) GetMsgType() EchoAreaMsgType {
	return EchoAreaMsgTypeJAM
}

// GetType return area type
func (j *JAM) GetType() EchoAreaType {
	return j.AreaType
}

// Init init
func (j *JAM) Init() {
}

// GetName return area name
func (j *JAM) GetName() string {
	return j.AreaName
}

var crc32tab = []uint32{ /* CRC polynomial 0xedb88320 */
	0x00000000, 0x77073096, 0xee0e612c, 0x990951ba, 0x076dc419, 0x706af48f, 0xe963a535, 0x9e6495a3,
	0x0edb8832, 0x79dcb8a4, 0xe0d5e91e, 0x97d2d988, 0x09b64c2b, 0x7eb17cbd, 0xe7b82d07, 0x90bf1d91,
	0x1db71064, 0x6ab020f2, 0xf3b97148, 0x84be41de, 0x1adad47d, 0x6ddde4eb, 0xf4d4b551, 0x83d385c7,
	0x136c9856, 0x646ba8c0, 0xfd62f97a, 0x8a65c9ec, 0x14015c4f, 0x63066cd9, 0xfa0f3d63, 0x8d080df5,
	0x3b6e20c8, 0x4c69105e, 0xd56041e4, 0xa2677172, 0x3c03e4d1, 0x4b04d447, 0xd20d85fd, 0xa50ab56b,
	0x35b5a8fa, 0x42b2986c, 0xdbbbc9d6, 0xacbcf940, 0x32d86ce3, 0x45df5c75, 0xdcd60dcf, 0xabd13d59,
	0x26d930ac, 0x51de003a, 0xc8d75180, 0xbfd06116, 0x21b4f4b5, 0x56b3c423, 0xcfba9599, 0xb8bda50f,
	0x2802b89e, 0x5f058808, 0xc60cd9b2, 0xb10be924, 0x2f6f7c87, 0x58684c11, 0xc1611dab, 0xb6662d3d,
	0x76dc4190, 0x01db7106, 0x98d220bc, 0xefd5102a, 0x71b18589, 0x06b6b51f, 0x9fbfe4a5, 0xe8b8d433,
	0x7807c9a2, 0x0f00f934, 0x9609a88e, 0xe10e9818, 0x7f6a0dbb, 0x086d3d2d, 0x91646c97, 0xe6635c01,
	0x6b6b51f4, 0x1c6c6162, 0x856530d8, 0xf262004e, 0x6c0695ed, 0x1b01a57b, 0x8208f4c1, 0xf50fc457,
	0x65b0d9c6, 0x12b7e950, 0x8bbeb8ea, 0xfcb9887c, 0x62dd1ddf, 0x15da2d49, 0x8cd37cf3, 0xfbd44c65,
	0x4db26158, 0x3ab551ce, 0xa3bc0074, 0xd4bb30e2, 0x4adfa541, 0x3dd895d7, 0xa4d1c46d, 0xd3d6f4fb,
	0x4369e96a, 0x346ed9fc, 0xad678846, 0xda60b8d0, 0x44042d73, 0x33031de5, 0xaa0a4c5f, 0xdd0d7cc9,
	0x5005713c, 0x270241aa, 0xbe0b1010, 0xc90c2086, 0x5768b525, 0x206f85b3, 0xb966d409, 0xce61e49f,
	0x5edef90e, 0x29d9c998, 0xb0d09822, 0xc7d7a8b4, 0x59b33d17, 0x2eb40d81, 0xb7bd5c3b, 0xc0ba6cad,
	0xedb88320, 0x9abfb3b6, 0x03b6e20c, 0x74b1d29a, 0xead54739, 0x9dd277af, 0x04db2615, 0x73dc1683,
	0xe3630b12, 0x94643b84, 0x0d6d6a3e, 0x7a6a5aa8, 0xe40ecf0b, 0x9309ff9d, 0x0a00ae27, 0x7d079eb1,
	0xf00f9344, 0x8708a3d2, 0x1e01f268, 0x6906c2fe, 0xf762575d, 0x806567cb, 0x196c3671, 0x6e6b06e7,
	0xfed41b76, 0x89d32be0, 0x10da7a5a, 0x67dd4acc, 0xf9b9df6f, 0x8ebeeff9, 0x17b7be43, 0x60b08ed5,
	0xd6d6a3e8, 0xa1d1937e, 0x38d8c2c4, 0x4fdff252, 0xd1bb67f1, 0xa6bc5767, 0x3fb506dd, 0x48b2364b,
	0xd80d2bda, 0xaf0a1b4c, 0x36034af6, 0x41047a60, 0xdf60efc3, 0xa867df55, 0x316e8eef, 0x4669be79,
	0xcb61b38c, 0xbc66831a, 0x256fd2a0, 0x5268e236, 0xcc0c7795, 0xbb0b4703, 0x220216b9, 0x5505262f,
	0xc5ba3bbe, 0xb2bd0b28, 0x2bb45a92, 0x5cb36a04, 0xc2d7ffa7, 0xb5d0cf31, 0x2cd99e8b, 0x5bdeae1d,
	0x9b64c2b0, 0xec63f226, 0x756aa39c, 0x026d930a, 0x9c0906a9, 0xeb0e363f, 0x72076785, 0x05005713,
	0x95bf4a82, 0xe2b87a14, 0x7bb12bae, 0x0cb61b38, 0x92d28e9b, 0xe5d5be0d, 0x7cdcefb7, 0x0bdbdf21,
	0x86d3d2d4, 0xf1d4e242, 0x68ddb3f8, 0x1fda836e, 0x81be16cd, 0xf6b9265b, 0x6fb077e1, 0x18b74777,
	0x88085ae6, 0xff0f6a70, 0x66063bca, 0x11010b5c, 0x8f659eff, 0xf862ae69, 0x616bffd3, 0x166ccf45,
	0xa00ae278, 0xd70dd2ee, 0x4e048354, 0x3903b3c2, 0xa7672661, 0xd06016f7, 0x4969474d, 0x3e6e77db,
	0xaed16a4a, 0xd9d65adc, 0x40df0b66, 0x37d83bf0, 0xa9bcae53, 0xdebb9ec5, 0x47b2cf7f, 0x30b5ffe9,
	0xbdbdf21c, 0xcabac28a, 0x53b39330, 0x24b4a3a6, 0xbad03605, 0xcdd70693, 0x54de5729, 0x23d967bf,
	0xb3667a2e, 0xc4614ab8, 0x5d681b02, 0x2a6f2b94, 0xb40bbe37, 0xc30c8ea1, 0x5a05df1b, 0x2d02ef8d,
}

func tolower(b byte) byte {
	if b >= 0x41 && b <= 0x5a {
		b += 32
	}
	return b
}

func crc32rieee(data []byte) uint32 {
	crcI := uint32(0xffffffff)
	for i := 0; i < len(data); i++ {
		crcI = (crcI >> 8) ^ crc32tab[byte(crcI)^tolower(data[i])]
	}
	return crcI
}

func crc32r(str string) uint32 {
	//bstr := []byte(strings.ToLower(str))
	return crc32rieee([]byte(str))
}

// SetLast set last message
func (j *JAM) SetLast(l uint32) {
	if l == 0 {
		l = 1
	}
	found := -1
	for i, lr := range j.lastRead {
		if lr.UserCRC == crc32r(config.Config.Username) {
			found = i
		}
	}
	if found == -1 {
		j.lastRead = append(j.lastRead, jamL{
			crc32r(config.Config.Username),
			crc32r(config.Config.Username),
			j.indexStructure[l-1].MessageNum,
			j.indexStructure[l-1].MessageNum})
	} else {
		j.lastRead[found].LastReadMsg = j.indexStructure[l-1].MessageNum // l
		if j.indexStructure[l-1].MessageNum > j.lastRead[found].HighReadMsg {
			j.lastRead[found].HighReadMsg = j.indexStructure[l-1].MessageNum
		}
	}
	buf := new(bytes.Buffer)
	err := binary.Write(buf, binary.LittleEndian, j.lastRead)
	if err != nil {
		return
	}
	err = os.WriteFile(j.AreaPath+".jlr", buf.Bytes(), 0644)
	if err != nil {
		return
	}
}
func packJamKludge(b io.Writer, loID uint16, hiID uint16, data []byte) {
	datLen := uint32(len(data))
	binary.Write(b, binary.LittleEndian, loID)
	binary.Write(b, binary.LittleEndian, hiID)
	binary.Write(b, binary.LittleEndian, datLen)
	binary.Write(b, binary.LittleEndian, data)
}
func packJamKludges(tm *Message) []byte {
	klb := new(bytes.Buffer)
	for kl, v := range tm.Kludges {
		switch kl {
		case "MSGID:":
			packJamKludge(klb, 4, 0, []byte(v))
		case "REPLYID:":
			packJamKludge(klb, 5, 0, []byte(v))
		case "PID:":
			packJamKludge(klb, 7, 0, []byte(v))
		default:
			packJamKludge(klb, 2000, 0, []byte(kl+" "+v))
		}
	}
	packJamKludge(klb, 0, 0, []byte(tm.FromAddr.String()))
	if tm.ToAddr != nil {
		packJamKludge(klb, 1, 0, []byte(tm.ToAddr.String()))
	}
	packJamKludge(klb, 2, 0, []byte(tm.From))
	packJamKludge(klb, 3, 0, []byte(tm.To))
	packJamKludge(klb, 6, 0, []byte(tm.Subject))
	return klb.Bytes()
}

// SaveMsg save message
func (j *JAM) SaveMsg(tm *Message) error {
	//	if len(j.indexStructure) == 0 {
	//		return errors.New("creating JAM area not implemented")
	//	}
	var jhr jhrS
	if len(j.indexStructure) == 0 {
		jhr.Signature = 0x4d414a
		jhr.PasswordCRC = 0xffffffff
		jhr.BaseMsgNum = 1
		jhr.DateCreated = uint32(time.Now().Unix())
	} else {
		jhr = j.headerStructure
	}

	jamh := jamH{Signature: 0x4d414a, Revision: 1, Attribute: 0x01000001}
	tm.Encode()
	kl := packJamKludges(tm)
	jamh.SubfieldLen = uint32(len(kl))
	jamh.MSGIDcrc = crc32r(tm.Kludges["MSGID:"])
	if val, ok := tm.Kludges["REPLYID:"]; ok {
		jamh.REPLYcrc = crc32r(val)
	} else {
		jamh.REPLYcrc = 0xffffffff
	}
	jamh.PasswordCRC = 0xffffffff
	jamh.DateWritten = uint32(tm.DateWritten.Unix())
	jamh.DateReceived = uint32(tm.DateArrived.Unix())
	jamh.DateProcessed = uint32(tm.DateArrived.Unix())
	jamh.TxtLen = uint32(len(tm.Body))
	jamh.MessageNumber = uint32(len(j.indexStructure)) + jhr.BaseMsgNum
	jam := jamSH{ToCRC: crc32r(tm.To)}
	f, err := os.OpenFile(j.AreaPath+".jdt", os.O_RDWR|os.O_CREATE, 0644)
	if err != nil {
		return err
	}
	defer f.Close()
	offset, _ := f.Seek(0, 2)
	jamh.Offset = uint32(offset)
	f.Write([]byte(tm.Body))
	f.Close()
	f, err = os.OpenFile(j.AreaPath+".jhr", os.O_RDWR|os.O_CREATE, 0644)
	if err != nil {
		return err
	}
	defer f.Close()
	jhr.ActiveMsgs++
	buf := new(bytes.Buffer)
	err = utils.WriteStructToBuffer(buf, &jhr)
	if err != nil {
		return err
	}
	j.headerStructure = jhr
	f.Seek(0, 0)
	f.Write(buf.Bytes())
	buf.Reset()
	offset, _ = f.Seek(0, 2)
	jam.Offset = uint32(offset)
	err = utils.WriteStructToBuffer(buf, &jamh)
	if err != nil {
		return err
	}
	f.Write(buf.Bytes())
	f.Write(kl)
	f.Close()
	buf.Reset()
	f, err = os.OpenFile(j.AreaPath+".jdx", os.O_RDWR|os.O_CREATE, 0644)
	if err != nil {
		return err
	}
	defer f.Close()
	if len(j.indexStructure) == 0 {
		f.Seek(0, 0)
	} else {
		f.Seek(0, 2)
	}
	err = utils.WriteStructToBuffer(buf, &jam)
	if err != nil {
		return err
	}
	f.Write(buf.Bytes())
	f.Close()
	j.indexStructure = append(j.indexStructure, jamS{jamh.MessageNumber, jam})
	return nil
}

// SetChrs set charset
func (j *JAM) SetChrs(c string) {
	j.Chrs = c
}

// GetChrs get charset
func (j *JAM) GetChrs() string {
	return j.Chrs
}

// GetMessages get headers
func (j *JAM) GetMessages() *[]MessageListItem {
	if len(j.messages) > 0 || len(j.indexStructure) == 0 {
		return &j.messages
	}
	for i := uint32(0); i < j.GetCount(); i++ {
		m, err := j.GetMsg(i + 1)
		if err != nil {
			continue
		}
		j.messages = append(j.messages, MessageListItem{
			MsgNum:      i + 1,
			From:        m.From,
			To:          m.To,
			Subject:     m.Subject,
			DateWritten: m.DateWritten,
		})
	}
	return &j.messages
}

// DelMsg remove msg
func (j *JAM) DelMsg(l uint32) error {
	if l == 0 {
		l = 1
	}
	fJhr, err := os.OpenFile(j.AreaPath+".jhr", os.O_RDWR, 0644)
	if err != nil {
		return err
	}
	defer fJhr.Close()
	_, err = fJhr.Seek(int64(j.indexStructure[l-1].jamsh.Offset), 0)
	if err != nil {
		return err
	}
	header := make([]byte, 76)
	fJhr.Read(header)
	headerb := bytes.NewBuffer(header)
	var jamh jamH
	if err = utils.ReadStructFromBuffer(headerb, &jamh); err != nil {
		return err
	}
	if jamh.Signature != 0x4d414a {
		return errors.New("wrong message signature")
	}
	jamh.Attribute |= 0x80000000
	err = utils.WriteStructToBuffer(headerb, &jamh)
	if err != nil {
		return err
	}
	_, err = fJhr.Seek(int64(j.indexStructure[l-1].jamsh.Offset), 0)
	if err != nil {
		return err
	}
	_, err = fJhr.Write(headerb.Bytes())
	if err != nil {
		return err
	}
	fJhr.Close()
	return nil
}

package msgapi

import (
	"bufio"
	"bytes"
	"encoding/binary"
	"errors"
	"github.com/askovpen/goated/lib/config"
	"github.com/askovpen/goated/lib/types"
	"github.com/askovpen/goated/lib/utils"
	"hash/crc32"
	"io/ioutil"
	// "log"
	"os"
	// "sort"
	"strings"
	"time"
)

// JAM struct
type JAM struct {
	AreaPath, AreaName string
	AreaType           EchoAreaType
	Chrs               string
	indexStructure     []jamS
	lastRead           []jamL
}

type jamS struct {
	MessageNum uint32
	jamsh      jamSH
}

type jamSH struct {
	ToCRC, Offset uint32
}

type jhrS struct {
	Signature                           [4]byte
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
		"Del", "", "", "",
		"", "", "", "",
	}
	i := 0
	for a > 0 {
		if a&1 > 0 {
			if datr[i] != "" {
				attrs = append(attrs, datr[i])
			}
		}
		i++
		a = a >> 1
	}
	return
}

// GetMsg return msg
func (j *JAM) GetMsg(position uint32) (*Message, error) {
	if len(j.indexStructure) == 0 {
		return nil, errors.New("Empty Area")
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
	var header []byte
	header = make([]byte, 76)
	fJhr.Read(header)
	headerb := bytes.NewBuffer(header)
	var jamh jamH
	if err = utils.ReadStructFromBuffer(headerb, &jamh); err != nil {
		return nil, err
	}
	if jamh.Signature != 0x4d414a {
		return nil, errors.New("wrong message signature")
	}
	rm := &Message{}
	rm.Area = j.AreaName
	rm.MsgNum = position
	rm.MaxNum = uint32(len(j.indexStructure))
	rm.DateWritten = time.Unix(int64(jamh.DateWritten), 0)
	_, tofs := rm.DateWritten.Zone()
	if jamh.DateReceived > 0 {
		rm.DateArrived = time.Unix(int64(jamh.DateReceived), 0)
	} else {
		rm.DateArrived = time.Unix(int64(jamh.DateProcessed), 0)
	}
	rm.DateWritten = rm.DateWritten.Add(time.Duration(tofs) * -time.Second)
	rm.DateArrived = rm.DateArrived.Add(time.Duration(tofs) * -time.Second)
	rm.Attrs = j.getAttrs(jamh.Attribute)
	deleted := false
	if jamh.Attribute&0x80000000 > 0 {
		deleted = true
	}
	rm.Body += ""
	var kl []byte
	kl = make([]byte, jamh.SubfieldLen)
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
		var val []byte
		val = make([]byte, datLen)
		binary.Read(klb, binary.LittleEndian, &val)
		switch LoID {
		case 0:
			fr := types.AddrFromString(string(val[:]))
			if fr != nil {
				rm.FromAddr = fr
			}
		case 1:
			if j.AreaType != EchoAreaTypeLocal && j.AreaType != EchoAreaTypeEcho {
				rm.ToAddr = types.AddrFromString(string(val[:]))
			}
		case 2:
			rm.From = string(val[:])
		case 3:
			if !deleted {
				if crc32r(string(val[:])) != j.indexStructure[position-1].jamsh.ToCRC {
					rm.Corrupted = true
				}
			}
			rm.To = string(val[:])
		case 4:
			if crc32r(string(val[:])) != jamh.MSGIDcrc {
				rm.Corrupted = true
			}
			rm.Body += "\x01MSGID: " + string(val[:]) + "\x0d"
		case 5:
			if crc32r(string(val[:])) != jamh.REPLYcrc {
				rm.Corrupted = true
			}
			rm.Body += "\x01REPLYID: " + string(val[:]) + "\x0d"
		case 6:
			rm.Subject = string(val[:])
		case 7:
			rm.Body += "\x01PID: " + string(val[:]) + "\x0d"
		case 8:
			afterBody += "\x01Via " + string(val[:]) + "\x0d"
		case 2004:
			rm.Body += "\x01TZUTC: " + string(val[:]) + "\x0d"
		case 2000:
			rm.Body += "\x01" + string(val[:]) + "\x0d"
		case 2001:
			afterBody += "SEEN-BY: " + string(val[:]) + "\x0d"
		case 2002:
			afterBody += "\x01PATH: " + string(val[:]) + "\x0d"
		}
	}
	fJdt, err := os.Open(j.AreaPath + ".jdt")
	if err != nil {
		return nil, err
	}
	defer fJdt.Close()
	fJdt.Seek(int64(jamh.Offset), 0)
	var txt []byte
	txt = make([]byte, jamh.TxtLen)
	fJdt.Read(txt)
	rm.Body += string(txt[:])
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
				j.indexStructure = append(j.indexStructure, jamS{i + 1, jam})
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
			// log.Printf("GetLast()->%d",j.getPositionOfJamMsg(l.LastReadMsg) + 1)
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

func crc32r(str string) uint32 {
	bstr := []byte(strings.ToLower(str))
	return 0xffffffff - crc32.ChecksumIEEE(bstr)
}

// SetLast set last message
func (j *JAM) SetLast(l uint32) {
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
	err = ioutil.WriteFile(j.AreaPath+".jlr", buf.Bytes(), 0644)
	if err != nil {
		return
	}
}
func packJamKludge(b *bytes.Buffer, LoID uint16, HiID uint16, data []byte) {
	datLen := uint32(len(data))
	binary.Write(b, binary.LittleEndian, LoID)
	binary.Write(b, binary.LittleEndian, HiID)
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
	// log.Printf("klb: %#v", klb.Bytes())
	return klb.Bytes()
}

// SaveMsg save message
func (j *JAM) SaveMsg(tm *Message) error {
	if len(j.indexStructure) == 0 {
		return errors.New("creating JAM area not implemented")
	}
	jamh := jamH{Signature: 0x4d414a, Revision: 1, Attribute: 0x02000001}
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
	jamh.MessageNumber = uint32(len(j.indexStructure)) + 1
	var jam jamSH
	jam.ToCRC = crc32r(tm.To)
	f, err := os.OpenFile(j.AreaPath+".jdt", os.O_RDWR, 0644)
	if err != nil {
		return err
	}
	defer f.Close()
	offset, _ := f.Seek(0, 2)
	jamh.Offset = uint32(offset)
	// log.Printf("offset: %d", offset)
	f.Write([]byte(tm.Body))
	f.Close()
	f, err = os.OpenFile(j.AreaPath+".jhr", os.O_RDWR, 0644)
	if err != nil {
		return err
	}
	defer f.Close()
	var jhr jhrS
	header := make([]byte, 1024)
	f.Read(header)
	headerb := bytes.NewBuffer(header)
	if err := utils.ReadStructFromBuffer(headerb, &jhr); err != nil {
		return err
	}
	jhr.ActiveMsgs++
	buf := new(bytes.Buffer)
	err = utils.WriteStructToBuffer(buf, &jhr)
	if err != nil {
		return err
	}
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
	f, err = os.OpenFile(j.AreaPath+".jdx", os.O_RDWR, 0644)
	if err != nil {
		return err
	}
	defer f.Close()
	f.Seek(0, 2)
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

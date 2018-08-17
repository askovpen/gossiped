package msgapi

import (
	"bufio"
	"bytes"
	"encoding/binary"
	"errors"
	"fmt"
	"hash/crc32"
	"io/ioutil"
	//  "log"
	"github.com/askovpen/goated/lib/config"
	"github.com/askovpen/goated/lib/types"
	"github.com/askovpen/goated/lib/utils"
	"os"
	"strings"
	"time"
)

type JAM struct {
	AreaPath       string
	AreaName       string
	indexStructure []jam_s
	lastRead       []jam_l
}

type jam_s struct {
	MessageNum uint32
	jamsh      jam_sh
}

type jam_sh struct {
	ToCRC  uint32
	Offset uint32
}

type jam_h struct {
	Signature     uint32
	Revision      uint16
	ReservedWord  uint16
	SubfieldLen   uint32
	TimesRead     uint32
	MSGIDcrc      uint32
	REPLYcrc      uint32
	ReplyTo       uint32
	Reply1st      uint32
	ReplyNext     uint32
	DateWritten   uint32
	DateReceived  uint32
	DateProcessed uint32
	MessageNumber uint32
	Attribute     uint32
	Attribute2    uint32
	Offset        uint32
	TxtLen        uint32
	PasswordCRC   uint32
	Cost          uint32
}

type jam_l struct {
	UserCRC, UserId, LastReadMsg, HighReadMsg uint32
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
		"", "Del", "", "",
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

func (j *JAM) GetMsg(position uint32) (*Message, error) {
	if len(j.indexStructure) == 0 {
		return nil, errors.New("Empty Area")
	}
	if position == 0 {
		position = 1
	}
	f, err := os.Open(j.AreaPath + ".jhr")
	if err != nil {
		return nil, err
	}
	defer f.Close()
	_, err = f.Seek(int64(j.indexStructure[position-1].jamsh.Offset), 0)
	if err != nil {
		return nil, err
	}
	var header []byte
	header = make([]byte, 76)
	//  reader := bufio.NewReader(f)
	f.Read(header)
	headerb := bytes.NewBuffer(header)
	var jamh jam_h
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
	//  _, tofs:=time.Now().Local().Zone()
	rm.DateWritten = time.Unix(int64(jamh.DateWritten), 0)
	_, tofs := rm.DateWritten.Zone()
	if jamh.DateReceived > 0 {
		rm.DateArrived = time.Unix(int64(jamh.DateReceived), 0)
	} else {
		rm.DateArrived = time.Unix(int64(jamh.DateProcessed), 0)
	}
	rm.DateWritten = rm.DateWritten.Add(time.Duration(tofs) * -time.Second)
	rm.DateArrived = rm.DateArrived.Add(time.Duration(tofs) * -time.Second)
	//  rm.Attr=jamh.Attribute
	rm.Attrs = j.getAttrs(jamh.Attribute)
	deleted := false
	if jamh.Attribute&0x80000000 > 0 {
		deleted = true
	}
	rm.Body += ""
	var kl []byte
	kl = make([]byte, jamh.SubfieldLen)
	f.Read(kl)
	//log.Printf("kl: %x", kl)
	klb := bytes.NewBuffer(kl)
	for {
		var LoID, HiID uint16
		var datLen uint32
		err = binary.Read(klb, binary.LittleEndian, &LoID)
		if err != nil {
			break
		}
		binary.Read(klb, binary.LittleEndian, &HiID)
		binary.Read(klb, binary.LittleEndian, &datLen)
		var val []byte
		val = make([]byte, datLen)
		binary.Read(klb, binary.LittleEndian, &val)
		//log.Printf("%d, %d (%d): %s",LoID, HiID, datLen, val)
		switch LoID {
		case 0:
			rm.FromAddr = types.AddrFromString(string(val[:]))
		case 1:
			rm.ToAddr = types.AddrFromString(string(val[:]))
		case 2:
			rm.From = string(val[:])
		case 3:
			if !deleted {
				if crc32r(string(val[:])) != j.indexStructure[position-1].jamsh.ToCRC {
					return nil, errors.New(fmt.Sprintf("'To' crc incorrect, got %08x, need %08x", crc32r(string(val[:])), j.indexStructure[position-1].jamsh.ToCRC))
				}
			}
			rm.To = string(val[:])
		case 4:
			if crc32r(string(val[:])) != jamh.MSGIDcrc {
				return nil, errors.New("crc incorrect")
			}
			rm.Body += "\x01MSGID: " + string(val[:]) + "\x0d"
		case 5:
			if crc32r(string(val[:])) != jamh.REPLYcrc {
				return nil, errors.New("crc incorrect")
			}
			rm.Body += "\x01REPLYID: " + string(val[:]) + "\x0d"
		case 6:
			rm.Subject = string(val[:])
		case 7:
			rm.Body += "\x01PID: " + string(val[:]) + "\x0d"
		default:
			rm.Body += "\x01" + string(val[:]) + "\x0d"
		}
	}
	f, err = os.Open(j.AreaPath + ".jdt")
	if err != nil {
		return nil, err
	}
	f.Seek(int64(jamh.Offset), 0)
	defer f.Close()
	var txt []byte
	txt = make([]byte, jamh.TxtLen)
	f.Read(txt)
	rm.Body += string(txt[:])
	err = rm.ParseRaw()
	if err != nil {
		return nil, err
	}
	//log.Printf("msgh: %#v", jamh)
	//log.Printf("rm: %#v", rm)
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
			var jam jam_sh
			if err = utils.ReadStructFromBuffer(partb, &jam); err != nil {
				break
			}
			if jam.Offset != 0xffffffff { //&& jam.ToCRC!=0xffffffff) {
				j.indexStructure = append(j.indexStructure, jam_s{i + 1, jam})
			}
			i++
		}
	}
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
			var jaml jam_l
			if err = utils.ReadStructFromBuffer(partb, &jaml); err != nil {
				break
			}
			j.lastRead = append(j.lastRead, jaml)
		}
	}
	//log.Printf("%#v", j.lastRead)
}
func (j *JAM) getPositionOfJamMsg(mId uint32) uint32 {
	//log.Printf("%d %#v",mId,j.indexStructure)
	for i, ji := range j.indexStructure {
		if mId == ji.MessageNum {
			return uint32(i)
		}
	}
	return 0
}

func (j *JAM) GetLast() uint32 {
	j.readJLR()
	for _, l := range j.lastRead {
		if l.UserCRC == crc32r(config.Config.Username) {
			return j.getPositionOfJamMsg(l.LastReadMsg) + 1
		}
	}
	return 0
}

func (j *JAM) GetCount() uint32 {
	j.readJDX()
	return uint32(len(j.indexStructure))
}

func (j *JAM) GetType() EchoAreaType {
	return EchoAreaTypeJAM
}

func (j *JAM) Init() {
}

func (j *JAM) GetName() string {
	return j.AreaName
}

func crc32r(str string) uint32 {
	bstr := []byte(strings.ToLower(str))
	return 0xffffffff - crc32.ChecksumIEEE(bstr)
}

func (j *JAM) SetLast(l uint32) {
	found := -1
	for i, lr := range j.lastRead {
		if lr.UserCRC == crc32r(config.Config.Username) {
			found = i
		}
	}
	if found == -1 {
		j.lastRead = append(j.lastRead, jam_l{
			crc32r(config.Config.Username),
			crc32r(config.Config.Username),
			j.indexStructure[l-1].MessageNum,
			j.indexStructure[l-1].MessageNum})
	} else {
		j.lastRead[found].LastReadMsg = l
		if j.indexStructure[l-1].MessageNum > j.lastRead[found].HighReadMsg {
			j.lastRead[found].HighReadMsg = j.indexStructure[l-1].MessageNum
		}
	}
	buf := new(bytes.Buffer)
	err := binary.Write(buf, binary.LittleEndian, j.lastRead)
	if err != nil {
		//log.Print(err)
		return
	}
	err = ioutil.WriteFile(j.AreaPath+".jlr", buf.Bytes(), 0644)
	if err != nil {
		//log.Print(err)
		return
	}
}

func (j *JAM) SaveMsg(tm *Message) error {
	return errors.New("not implemented")
}

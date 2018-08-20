package msgapi

import (
	"bufio"
	"bytes"
	"encoding/binary"
	"errors"
	"fmt"
	"github.com/askovpen/goated/lib/types"
	"github.com/askovpen/goated/lib/utils"
	"io/ioutil"
	"log"
	"os"
	"strings"
	"time"
)

type SquishAttrs uint32

const (
	SquishPRIVATE SquishAttrs = 0x0001
	SquishCRASH   SquishAttrs = 0x0002
	SquishREAD    SquishAttrs = 0x0004
	SquishSENT    SquishAttrs = 0x0008
	SquishFILE    SquishAttrs = 0x0010
	SquishFWD     SquishAttrs = 0x0020
	SquishORPHAN  SquishAttrs = 0x0040
	SquishKILL    SquishAttrs = 0x0080
	SquishLOCAL   SquishAttrs = 0x0100
	SquishHOLD    SquishAttrs = 0x0200
	SquishXX2     SquishAttrs = 0x0400
	SquishFRQ     SquishAttrs = 0x0800
	SquishRRQ     SquishAttrs = 0x1000
	SquishCPT     SquishAttrs = 0x2000
	SquishARQ     SquishAttrs = 0x4000
	SquishURQ     SquishAttrs = 0x8000
	SquishSCANNED SquishAttrs = 0x00010000
	SquishUID     SquishAttrs = 0x00020000
	SquishSEEN    SquishAttrs = 0x00080000
)

type Squish struct {
	AreaPath       string
	AreaName       string
	AreaType       EchoAreaType
	indexStructure []sqi_s
}

type sqi_s struct {
	Offset     uint32
	MessageNum uint32
	CRC        uint32
}

type sqd_s struct {
	Len, Rsvd1                                                uint16
	NumMsg, HighMsg, SkipMsg, HighWater, Uid                  uint32
	Base                                                      [80]byte
	BeginFrame, LastFrame, FreeFrame, LastFreeFrame, EndFrame uint32
	MaxMsg                                                    uint32
	KeepDays                                                  uint16
	SzSQHdr                                                   uint16
	Rsvd2                                                     [124]byte
}

type sqd_h struct {
	Id, NextFrame, PrevFrame, FrameLength, MsgLength, CLen uint32
	FrameType, Rsvd                                        uint16
	Attr                                                   uint32
	From, To                                               [36]byte
	Subject                                                [72]byte
	FromZone, FromNet, FromNode, FromPoint                 uint16
	ToZone, ToNet, ToNode, ToPoint                         uint16
	DateWritten, DateArrived                               uint32
	Utc                                                    uint16
	ReplyTo                                                uint32
	Replies                                                [9]uint32
	UMsgId                                                 uint32
	Date                                                   [20]byte
}

func (s *Squish) getAttrs(a uint32) (attrs []string) {
	datr := []string{
		"Pvt", "", "Rcv", "Snt",
		"", "Trs", "", "K/s",
		"Loc", "", "", "",
		"Rrq", "", "Arq", "",
		"Scn", "", "", "",
		"", "", "", "",
		"", "", "", "",
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

func readSQDH(headerb *bytes.Buffer) (sqd_h, error) {
	var sqdh sqd_h
	if err := utils.ReadStructFromBuffer(headerb, &sqdh); err != nil {
		return sqdh, err
	}
	if sqdh.Id != 0xafae4453 {
		return sqdh, errors.New(fmt.Sprintf("Wrong Squish header %08x", sqdh.Id))
	}
	return sqdh, nil
}

func (s *Squish) GetMsg(position uint32) (*Message, error) {
	if len(s.indexStructure) == 0 {
		return nil, errors.New("Empty Area")
	}
	if position == 0 {
		position = 1
	}
	f, err := os.Open(s.AreaPath + ".sqd")
	if err != nil {
		return nil, err
	}
	defer f.Close()
	f.Seek(int64(s.indexStructure[position-1].Offset), 0)
	var header []byte
	header = make([]byte, 266)
	f.Read(header)
	headerb := bytes.NewBuffer(header)
	sqdh, err := readSQDH(headerb)
	if err != nil {
		return nil, err
	}
	//	var sqdh sqd_h
	//	if err = utils.ReadStructFromBuffer(headerb, &sqdh); err != nil {
	//		return nil, err
	//	}
	//log.Printf("%#v", sqdh)
	//var body []byte
	body := make([]byte, sqdh.MsgLength+28-266)
	f.Read(body)
	//log.Printf("%s", body)
	toHash := bufHash32(string(sqdh.To[:]))
	if sqdh.Attr&uint32(SquishREAD) > 0 {
		toHash = toHash | 0x80000000
	}
	if s.indexStructure[position-1].CRC != toHash {
		return nil, errors.New(fmt.Sprintf("Wrong message CRC need 0x%08x, got 0x%08x for name %s", s.indexStructure[position-1].CRC, bufHash32(string(sqdh.To[:])), sqdh.To))
	}
	rm := &Message{}
	rm.From = strings.Trim(string(sqdh.From[:]), "\x00")
	rm.To = strings.Trim(string(sqdh.To[:]), "\x00")
	rm.FromAddr = types.AddrFromNum(sqdh.FromZone, sqdh.FromNet, sqdh.FromNode, sqdh.FromPoint)
	if s.AreaType != EchoAreaTypeLocal && s.AreaType != EchoAreaTypeEcho {
		rm.ToAddr = types.AddrFromNum(sqdh.ToZone, sqdh.ToNet, sqdh.ToNode, sqdh.ToPoint)
	}
	rm.Subject = strings.Trim(string(sqdh.Subject[:]), "\x00")
	rm.Attrs = s.getAttrs(sqdh.Attr)
	rm.Body = string(body[:])
	rm.DateWritten = getTime(sqdh.DateWritten)
	rm.DateArrived = getTime(sqdh.DateArrived)
	kla := strings.Split(rm.Body[1:sqdh.CLen], "\x01")
	for i := range kla {
		kla[i] = strings.Trim(kla[i], "\x00")
	}
	rm.Body = "\x01" + strings.Join(kla, "\x0d\x01") + "\x0d" + rm.Body[sqdh.CLen:]
	//log.Printf("body: %s",rm.Body)
	//log.Printf("after Kludges: %d",rm.Body[sqdh.CLen-1])
	if strings.Index(rm.Body, "\x00") != -1 {
		rm.Body = rm.Body[0:strings.Index(rm.Body, "\x00")]
	}
	err = rm.ParseRaw()
	if err != nil {
		return nil, err
	}
	//log.Printf("msg: %#v", rm)
	return rm, nil
}

func (s *Squish) readSQI() {
	if len(s.indexStructure) > 0 {
		return
	}
	file, err := os.Open(s.AreaPath + ".sqi")
	if err != nil {
		return
	}
	reader := bufio.NewReader(file)
	part := make([]byte, 12288)
	for {
		count, err := reader.Read(part)
		if err != nil {
			break
		}
		partb := bytes.NewBuffer(part[:count])
		for {
			var sqi sqi_s
			if err = utils.ReadStructFromBuffer(partb, &sqi); err != nil {
				break
			}
			if sqi.Offset != 0 {
				s.indexStructure = append(s.indexStructure, sqi)
			}
		}
	}
	//  log.Printf("%s %#v", s.AreaName, s.indexStructure)
}
func (s *Squish) GetLast() uint32 {
	s.readSQI()
	if len(s.indexStructure) == 0 {
		return 0
	}
	file, err := os.Open(s.AreaPath + ".sql")
	defer file.Close()
	if err != nil {
		return 0
	}
	var ret uint32
	err = binary.Read(file, binary.LittleEndian, &ret)
	if err != nil {
		return 0
	}
	for i, is := range s.indexStructure {
		if ret == is.MessageNum {
			//      log.Printf("ret, i: %d %d",ret, i)
			return uint32(i + 1)
		}
	}
	return 0
}

func (s *Squish) GetCount() uint32 {
	s.readSQI()
	return uint32(len(s.indexStructure))
}

func (s *Squish) GetMsgType() EchoAreaMsgType {
	return EchoAreaMsgTypeSquish
}

func (s *Squish) GetType() EchoAreaType {
	return s.AreaType
}

func (s *Squish) Init() {
}

func (s *Squish) GetName() string {
	return s.AreaName
}

func getTime(t uint32) time.Time {
	return time.Date(
		int(t>>9&127)+1980,
		time.Month(int(t>>5&15)),
		int(t&31),
		int(t>>27&31),
		int(t>>21&63),
		int(t>>16&31)*2,
		0,
		time.Local)
}
func setTime(t time.Time) (rt uint32) {
	rt = 0
	rt |= uint32(t.Day() & 31)
	rt |= uint32(t.Month() & 15 << 5)
	rt |= uint32((t.Year() - 1980) & 127 << 9)
	rt |= uint32((t.Second() / 2) & 31 << 16)
	rt |= uint32(t.Minute() & 63 << 21)
	rt |= uint32(t.Hour() & 31 << 27)
	return
}
func bufHash32(str string) (h uint32) {
	h = 0
	for _, b := range strings.ToLower(str) {
		if b == 0 {
			continue
		}
		h = (h << 4) + uint32(b)
		g := h & 0xF0000000
		if g != 0 {
			h |= g >> 24
			h |= g
		}
	}
	h = h & 0x7fffffff
	return
}
func (s *Squish) SetLast(l uint32) {
	if l == 0 {
		l = 1
	}
	r := s.indexStructure[l-1].MessageNum
	buf := new(bytes.Buffer)
	err := binary.Write(buf, binary.LittleEndian, r)
	if err != nil {
		log.Print(err)
	}
	err = ioutil.WriteFile(s.AreaPath+".sql", buf.Bytes(), 0644)
	if err != nil {
		log.Print(err)
	}
}

func (s *Squish) SaveMsg(tm *Message) error {
	if len(s.indexStructure) == 0 {
		return errors.New("creating Squish area not implemented")
	}
	lastIdx := len(s.indexStructure) - 1
	var sqdh sqd_h
	var sqi sqi_s
	kludges := ""
	for kl, v := range tm.Kludges {
		kludges += "\x01" + kl + " " + v
	}
	kludges += "\x00"
	sqdh.Id = 0xafae4453
	sqdh.NextFrame = 0
	sqdh.PrevFrame = s.indexStructure[lastIdx].Offset
	sqdh.Attr = uint32(SquishLOCAL | SquishSEEN)
	tm.Encode()
	copy(sqdh.From[:], tm.From)
	copy(sqdh.To[:], tm.To)
	copy(sqdh.Subject[:], tm.Subject)
	copy(sqdh.Date[:], tm.DateWritten.Format("02 Jan 06  15:04:05"))
	sqdh.DateWritten = setTime(tm.DateWritten)
	sqdh.DateArrived = setTime(tm.DateArrived)
	sqdh.FromZone, sqdh.FromNet, sqdh.FromNode, sqdh.FromPoint = tm.FromAddr.GetZone(), tm.FromAddr.GetNet(), tm.FromAddr.GetNode(), tm.FromAddr.GetPoint()
	if s.AreaType == EchoAreaTypeNetmail {
		sqdh.ToZone, sqdh.ToNet, sqdh.ToNode, sqdh.ToPoint = tm.ToAddr.GetZone(), tm.ToAddr.GetNet(), tm.ToAddr.GetNode(), tm.ToAddr.GetPoint()
	} else {
		sqdh.ToZone, sqdh.ToNet, sqdh.ToNode, sqdh.ToPoint = 0, 0, 0, 0
	}
	sqdh.UMsgId = s.indexStructure[lastIdx].MessageNum + 1
	sqdh.CLen = uint32(len(kludges))
	body := kludges + tm.Body + "\x00"
	sqdh.MsgLength = uint32(len(body)) + 266 - 28
	sqdh.FrameLength = uint32(len(body)) + 266 - 28
	sqi.CRC = bufHash32(tm.To)
	sqi.MessageNum = sqdh.UMsgId
	f, err := os.OpenFile(s.AreaPath+".sqd", os.O_RDWR, 0644)
	if err != nil {
		return err
	}
	defer f.Close()
	var header []byte
	header = make([]byte, 256)
	f.Read(header)
	headerb := bytes.NewBuffer(header)
	var sqd sqd_s
	if err := utils.ReadStructFromBuffer(headerb, &sqd); err != nil {
		return err
	}
	sqi.Offset = sqd.EndFrame
	//log.Printf("sqd header: %#v",sqd)
	//	_, err=f.Seek(int64(sqd.EndFrame), 0)
	//	if err != nil {
	//		return err
	//	}
	sqd.NumMsg += 1
	sqd.HighMsg += 1
	sqd.Uid += 1
	sqd.LastFrame = sqd.EndFrame
	sqd.EndFrame = sqd.LastFrame + sqdh.FrameLength + 28
	f.Seek(0, 0)
	buf := new(bytes.Buffer)
	err = utils.WriteStructToBuffer(buf, &sqd)
	if err != nil {
		return err
	}
	//log.Printf("len: %d", len(buf.Bytes()))
	f.Write(buf.Bytes())
	buf.Reset()
	f.Seek(int64(sqdh.PrevFrame), 0)
	//var header []byte
	header = make([]byte, 266)
	f.Read(header)
	headerb = bytes.NewBuffer(header)
	prevSqdh, err := readSQDH(headerb)
	prevSqdh.NextFrame = sqi.Offset
	err = utils.WriteStructToBuffer(buf, &prevSqdh)
	if err != nil {
		return err
	}
	f.Seek(int64(sqdh.PrevFrame), 0)
	f.Write(buf.Bytes())
	buf.Reset()
	err = utils.WriteStructToBuffer(buf, &sqdh)
	if err != nil {
		return err
	}
	f.Seek(int64(sqi.Offset), 0)
	f.Write(buf.Bytes())
	buf.Reset()
	f.Write([]byte(body))
	f.Close()
	f, err = os.OpenFile(s.AreaPath+".sqi", os.O_RDWR, 0644)
	if err != nil {
		return err
	}
	defer f.Close()
	err = utils.WriteStructToBuffer(buf, &sqi)
	f.Seek(0, 2)
	f.Write(buf.Bytes())
	f.Close()
	s.indexStructure = append(s.indexStructure, sqi)
	//	log.Printf("%d",r)
	//	log.Printf("sqdh: %#v, %x",sqdh, body)
	return nil
	//return errors.New("not implemented")
}

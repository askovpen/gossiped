package msgapi

import (
	"bufio"
	"bytes"
	"encoding/binary"
	"errors"
	"fmt"
	"log"
	"os"
	"sort"
	"strings"
	"time"

	"github.com/askovpen/gossiped/pkg/types"
	"github.com/askovpen/gossiped/pkg/utils"
	//"unicode"
)

// SquishAttrs Squish Attributes
type SquishAttrs uint32

// attrubutes
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

// Squish struct
type Squish struct {
	AreaPath       string
	AreaName       string
	AreaType       EchoAreaType
	Chrs           string
	indexStructure []sqiS
	messages       []MessageListItem
}

type sqiS struct {
	Offset     uint32
	MessageNum uint32
	CRC        uint32
}

type sqdS struct {
	Len, Rsvd1                                                uint16
	NumMsg, HighMsg, SkipMsg, HighWater, UID                  uint32
	Base                                                      [80]byte
	BeginFrame, LastFrame, FreeFrame, LastFreeFrame, EndFrame uint32
	MaxMsg                                                    uint32
	KeepDays                                                  uint16
	SzSQHdr                                                   uint16
	Rsvd2                                                     [124]byte
}

type sqdH struct {
	ID, NextFrame, PrevFrame, FrameLength, MsgLength, CLen uint32
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
	UMsgID                                                 uint32
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
		a >>= 1
	}
	return
}

func (s *Squish) getOffsetByNum(num uint32) (offset uint32) {
	for i, is := range s.indexStructure {
		if is.MessageNum == num {
			return uint32(i) + 1
		}
	}
	return 0
}

func readSQDH(headerb *bytes.Buffer) (sqdH, error) {
	var sqdh sqdH
	if err := utils.ReadStructFromBuffer(headerb, &sqdh); err != nil {
		return sqdh, err
	}
	if sqdh.ID != 0xafae4453 {
		return sqdh, fmt.Errorf("wrong Squish header %08x", sqdh.ID)
	}
	return sqdh, nil
}

// GetMsg return message
func (s *Squish) GetMsg(position uint32) (*Message, error) {
	if len(s.indexStructure) == 0 {
		//		return nil, errors.New("Empty Area")
		return nil, nil
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
	//var header []byte
	header := make([]byte, 266)
	f.Read(header)
	headerb := bytes.NewBuffer(header)
	sqdh, err := readSQDH(headerb)
	if err != nil {
		return nil, err
	}
	body := make([]byte, sqdh.MsgLength+28-266)
	f.Read(body)
	toHash := bufHash32(string(sqdh.To[:]))
	if sqdh.Attr&uint32(SquishREAD) > 0 {
		toHash |= 0x80000000
	}
	rm := &Message{Area: s.AreaName, MsgNum: position}
	if s.indexStructure[position-1].CRC != toHash {
		rm.Corrupted = true
	}
	rm.From = strings.Trim(string(sqdh.From[:]), "\x00")
	rm.To = strings.Trim(string(sqdh.To[:]), "\x00")
	rm.FromAddr = types.AddrFromNum(sqdh.FromZone, sqdh.FromNet, sqdh.FromNode, sqdh.FromPoint)
	if s.AreaType != EchoAreaTypeLocal && s.AreaType != EchoAreaTypeEcho {
		rm.ToAddr = types.AddrFromNum(sqdh.ToZone, sqdh.ToNet, sqdh.ToNode, sqdh.ToPoint)
	}
	rm.Subject = strings.Trim(string(sqdh.Subject[:]), "\x00")
	rm.Attrs = s.getAttrs(sqdh.Attr)
	rm.Body = string(body)
	rm.DateWritten = getTime(sqdh.DateWritten)
	rm.DateArrived = getTime(sqdh.DateArrived)
	if sqdh.ReplyTo > 0 {
		rm.ReplyTo = s.getOffsetByNum(sqdh.ReplyTo) //s.indexStructure[sqdh.ReplyTo].MessageNum - 1
	} else {
		rm.ReplyTo = 0
	}
	for _, reply := range sqdh.Replies {
		if reply > 0 {
			//rm.Replies = append(rm.Replies, s.indexStructure[reply-1].MessageNum)
			rm.Replies = append(rm.Replies, s.getOffsetByNum(reply))
		}
	}
	if sqdh.CLen > 0 {
		kla := strings.Split(rm.Body[1:sqdh.CLen], "\x01")
		for i := range kla {
			kla[i] = strings.Trim(kla[i], "\x00")
		}
		rm.Body = "\x01" + strings.Join(kla, "\x0d\x01") + "\x0d" + rm.Body[sqdh.CLen:]
	}
	if strings.Contains(rm.Body, "\x00") {
		rm.Body = rm.Body[0:strings.Index(rm.Body, "\x00")]
	}
	err = rm.ParseRaw()
	if err != nil {
		return nil, err
	}
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
			var sqi sqiS
			if err = utils.ReadStructFromBuffer(partb, &sqi); err != nil {
				break
			}
			if sqi.Offset != 0 {
				s.indexStructure = append(s.indexStructure, sqi)
			}
		}
	}
	sort.Slice(s.indexStructure, func(i, j int) bool { return s.indexStructure[i].MessageNum < s.indexStructure[j].MessageNum })
}

// GetLast get last message number
func (s *Squish) GetLast() uint32 {
	s.readSQI()
	if len(s.indexStructure) == 0 {
		return 0
	}
	file, err := os.Open(s.AreaPath + ".sql")
	if err != nil {
		return 0
	}
	defer file.Close()
	var ret uint32
	err = binary.Read(file, binary.LittleEndian, &ret)
	if err != nil {
		return 0
	}
	for i, is := range s.indexStructure {
		if ret == is.MessageNum {
			return uint32(i + 1)
		}
	}
	if ret != 0 {
		return uint32(len(s.indexStructure))
	}
	return 0
}

// GetCount get messages count
func (s *Squish) GetCount() uint32 {
	s.readSQI()
	return uint32(len(s.indexStructure))
}

// GetMsgType return area msg base type
func (s *Squish) GetMsgType() EchoAreaMsgType {
	return EchoAreaMsgTypeSquish
}

// GetType get area type
func (s *Squish) GetType() EchoAreaType {
	return s.AreaType
}

// Init for future
func (s *Squish) Init() {
}

// GetName return area name
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
	//str = strings.ToLower(str)
	strb := []byte(str)
	//strb=bytes.ToLower(strb)
	h = 0
	for i := range strb {
		if strb[i] == 0 {
			continue
		}
		//if strb[i] < 0x7f {
		//	strb[i] = byte(unicode.ToLower(rune(strb[i])))
		//}
		strb[i] = tolower(strb[i])
		h = (h << 4) + uint32(strb[i])
		g := h & 0xF0000000
		if g != 0 {
			h |= g >> 24
			h |= g
		}
	}
	h &= 0x7fffffff
	return
}

// SetLast set last message number
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
	err = os.WriteFile(s.AreaPath+".sql", buf.Bytes(), 0644)
	if err != nil {
		log.Print(err)
	}
}

// SaveMsg save message
func (s *Squish) SaveMsg(tm *Message) error {
	lastIdx := len(s.indexStructure) - 1
	if len(s.indexStructure) == 0 {
		lastIdx = 0
	}
	//var sqi sqiS
	kludges := ""
	tm.Encode()
	for kl, v := range tm.Kludges {
		kludges += "\x01" + kl + " " + v
	}
	kludges += "\x00"
	body := kludges + tm.Body + "\x00"
	sqdh := sqdH{ID: 0xafae4453,
		NextFrame:   0,
		Attr:        uint32(SquishLOCAL | SquishSEEN),
		DateWritten: setTime(tm.DateWritten),
		DateArrived: setTime(tm.DateArrived),
		FromZone:    tm.FromAddr.GetZone(),
		FromNet:     tm.FromAddr.GetNet(),
		FromNode:    tm.FromAddr.GetNode(),
		FromPoint:   tm.FromAddr.GetPoint(),
		CLen:        uint32(len(kludges)),
		MsgLength:   uint32(len(body)) + 266 - 28,
		FrameLength: uint32(len(body)) + 266 - 28}
	if len(s.indexStructure) > 0 {
		sqdh.PrevFrame = s.indexStructure[lastIdx].Offset
	}
	copy(sqdh.From[:], tm.From)
	copy(sqdh.To[:], tm.To)
	copy(sqdh.Subject[:], tm.Subject)
	copy(sqdh.Date[:], tm.DateWritten.Format("02 Jan 06  15:04:05"))
	if s.AreaType == EchoAreaTypeNetmail {
		sqdh.ToZone, sqdh.ToNet, sqdh.ToNode, sqdh.ToPoint = tm.ToAddr.GetZone(), tm.ToAddr.GetNet(), tm.ToAddr.GetNode(), tm.ToAddr.GetPoint()
	} else {
		sqdh.ToZone, sqdh.ToNet, sqdh.ToNode, sqdh.ToPoint = 0, 0, 0, 0
	}
	if len(s.indexStructure) == 0 {
		sqdh.UMsgID = 1
	} else {
		sqdh.UMsgID = s.indexStructure[lastIdx].MessageNum + 1
	}
	sqi := sqiS{CRC: bufHash32(tm.To), MessageNum: sqdh.UMsgID}
	f, err := os.OpenFile(s.AreaPath+".sqd", os.O_RDWR|os.O_CREATE, 0644)
	if err != nil {
		return err
	}
	defer f.Close()
	var header []byte
	var sqd sqdS
	var headerb *bytes.Buffer
	if len(s.indexStructure) == 0 {
		sqd.Len = 256
		sqd.EndFrame = 256
		sqd.UID = 1
		sqd.BeginFrame = 256
		sqd.SzSQHdr = 28
	} else {
		header = make([]byte, 256)
		f.Read(header)
		headerb = bytes.NewBuffer(header)
		if err := utils.ReadStructFromBuffer(headerb, &sqd); err != nil {
			return err
		}
	}
	sqi.Offset = sqd.EndFrame
	sqd.NumMsg++
	sqd.HighMsg++
	sqd.UID++
	sqd.LastFrame = sqd.EndFrame
	sqd.EndFrame = sqd.LastFrame + sqdh.FrameLength + 28
	f.Seek(0, 0)
	buf := new(bytes.Buffer)
	err = utils.WriteStructToBuffer(buf, &sqd)
	if err != nil {
		return err
	}
	f.Write(buf.Bytes())
	buf.Reset()
	if sqdh.PrevFrame > 0 {
		f.Seek(int64(sqdh.PrevFrame), 0)
		header = make([]byte, 266)
		f.Read(header)
		headerb = bytes.NewBuffer(header)
		prevSqdh, err := readSQDH(headerb)
		if err != nil {
			return err
		}
		prevSqdh.NextFrame = sqi.Offset
		err = utils.WriteStructToBuffer(buf, &prevSqdh)
		if err != nil {
			return err
		}
		f.Seek(int64(sqdh.PrevFrame), 0)
		f.Write(buf.Bytes())
		buf.Reset()
	}
	err = utils.WriteStructToBuffer(buf, &sqdh)
	if err != nil {
		return err
	}
	f.Seek(int64(sqi.Offset), 0)
	f.Write(buf.Bytes())
	buf.Reset()
	f.Write([]byte(body))
	f.Close()
	f, err = os.OpenFile(s.AreaPath+".sqi", os.O_RDWR|os.O_CREATE, 0644)
	if err != nil {
		return err
	}
	defer f.Close()
	utils.WriteStructToBuffer(buf, &sqi)
	if len(s.indexStructure) == 0 {
		f.Seek(0, 0)
	} else {
		f.Seek(0, 2)
	}
	f.Write(buf.Bytes())
	f.Close()
	s.indexStructure = append(s.indexStructure, sqi)
	return nil
}

// SetChrs set charset
func (s *Squish) SetChrs(c string) {
	s.Chrs = c
}

// GetChrs get charset
func (s *Squish) GetChrs() string {
	return s.Chrs
}

// GetMessages get headers
func (s *Squish) GetMessages() *[]MessageListItem {
	if len(s.messages) > 0 || len(s.indexStructure) == 0 {
		return &s.messages
	}

	for i := uint32(0); i < s.GetCount(); i++ {
		m, err := s.GetMsg(i + 1)
		if err != nil {
			continue
		}
		s.messages = append(s.messages, MessageListItem{
			MsgNum:      i + 1,
			From:        m.From,
			To:          m.To,
			Subject:     m.Subject,
			DateWritten: m.DateWritten,
		})
	}
	return &s.messages
}

// DelMsg remove msg
func (s *Squish) DelMsg(l uint32) error {
	if len(s.indexStructure) == 0 {
		return errors.New("empty Area")
	}
	if l == 0 {
		l = 1
	}
	f, err := os.OpenFile(s.AreaPath+".sqd", os.O_RDWR, 0644)
	if err != nil {
		return err
	}
	defer f.Close()
	f.Seek(int64(s.indexStructure[l-1].Offset), 0)
	//var header []byte
	header := make([]byte, 266)
	f.Read(header)
	headerb := bytes.NewBuffer(header)
	sqdh, err := readSQDH(headerb)
	if err != nil {
		return err
	}
	prev := sqdh.PrevFrame
	next := sqdh.NextFrame
	if l > 1 {
		f.Seek(int64(s.indexStructure[l-2].Offset), 0)
		f.Read(header)
		headerb := bytes.NewBuffer(header)
		sqdh, err = readSQDH(headerb)
		if err != nil {
			return err
		}
		sqdh.NextFrame = next
		err = utils.WriteStructToBuffer(headerb, &sqdh)
		if err != nil {
			return err
		}
		f.Seek(int64(s.indexStructure[l-2].Offset), 0)
		f.Write(headerb.Bytes())
	}
	if prev > 0 && int(l) < len(s.indexStructure)-1 {
		f.Seek(int64(s.indexStructure[l].Offset), 0)
		f.Read(header)
		headerb := bytes.NewBuffer(header)
		sqdh, err = readSQDH(headerb)
		if err != nil {
			return err
		}
		sqdh.PrevFrame = prev
		err = utils.WriteStructToBuffer(headerb, &sqdh)
		if err != nil {
			return err
		}
		f.Seek(int64(s.indexStructure[l].Offset), 0)
		f.Write(headerb.Bytes())
	}
	f.Seek(0, 0)
	var sqd sqdS
	header = make([]byte, 256)
	f.Read(header)
	headerb = bytes.NewBuffer(header)
	if err := utils.ReadStructFromBuffer(headerb, &sqd); err != nil {
		return err
	}
	sqd.NumMsg--
	sqd.HighMsg--
	sqd.UID--
	f.Seek(0, 0)
	err = utils.WriteStructToBuffer(headerb, &sqd)
	if err != nil {
		return err
	}
	f.Write(headerb.Bytes())

	f.Close()
	if len(s.messages) > 0 {
		s.messages = append(s.messages[:l-1], s.messages[l:]...)
	}
	s.indexStructure = append(s.indexStructure[:l-1], s.indexStructure[l:]...)
	f, err = os.OpenFile(s.AreaPath+".sqi", os.O_RDWR, 0644)
	if err != nil {
		return err
	}
	f.Truncate(0)
	f.Seek(0, 0)
	for _, is := range s.indexStructure {
		buf := new(bytes.Buffer)
		utils.WriteStructToBuffer(buf, &is)
		f.Write(buf.Bytes())
	}
	defer f.Close()

	return nil
}

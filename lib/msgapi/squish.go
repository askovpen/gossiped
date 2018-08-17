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

type Squish struct {
	AreaPath       string
	AreaName       string
	indexStructure []sqi_s
}

type sqi_s struct {
	Offset     uint32
	MessageNum uint32
	CRC        uint32
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
	var sqdh sqd_h
	if err = utils.ReadStructFromBuffer(headerb, &sqdh); err != nil {
		return nil, err
	}
	//log.Printf("%#v", sqdh)
	//var body []byte
	body := make([]byte, sqdh.MsgLength+28-266)
	f.Read(body)
	//log.Printf("%s", body)
	if s.indexStructure[position-1].CRC != bufHash32(string(sqdh.To[:])) && s.indexStructure[position-1].CRC != bufHash32(string(sqdh.To[:]))|0x80000000 {
		return nil, errors.New(fmt.Sprintf("Wrong message CRC need 0x%08x, got 0x%08x for name %s", s.indexStructure[position-1].CRC, bufHash32(string(sqdh.To[:])), sqdh.To))
	}
	rm := &Message{}
	rm.From = strings.Trim(string(sqdh.From[:]), "\x00")
	rm.To = strings.Trim(string(sqdh.To[:]), "\x00")
	rm.FromAddr = types.AddrFromNum(sqdh.FromZone, sqdh.FromNet, sqdh.FromNode, sqdh.FromPoint)
	rm.ToAddr = types.AddrFromNum(sqdh.ToZone, sqdh.ToNet, sqdh.ToNode, sqdh.ToPoint)
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

func (s *Squish) GetType() EchoAreaType {
	return EchoAreaTypeSquish
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
	return errors.New("not implemented")
}

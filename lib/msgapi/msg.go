package msgapi

import (
	"bytes"
	"encoding/binary"
	"errors"
	"github.com/askovpen/goated/lib/utils"
	"io/ioutil"
	"log"
	"os"
	"path/filepath"
	"sort"
	"strconv"
	"strings"
	"time"
)

type MSG struct {
	AreaPath    string
	AreaName    string
	lastreads   string
	messageNums []uint32
}

type msg_s struct {
	From        [36]byte
	To          [36]byte
	Subj        [72]byte
	Date        [20]byte
	Times       uint16
	DestNode    uint16
	OrigNode    uint16
	Cost        uint16
	OrigNet     uint16
	DestNet     uint16
	DateWritten uint32
	DateArrived uint32
	Reply       uint16
	Attr        uint16
	Up          uint16
	Body        string
}

func (m *MSG) Init() {
}

func (m *MSG) getAttrs(a uint16) (attrs []string) {
	datr := []string{
		"Pvt", "", "Rcv", "Snt",
		"", "Trs", "", "K/s",
		"Loc", "", "", "",
		"Rrq", "", "Arq", "",
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

func (m *MSG) GetMsg(position uint32) (*Message, error) {
	if len(m.messageNums) == 0 {
		return nil, errors.New("Empty Area")
	}
	if position == 0 {
		position = 1
	}
	f, err := os.Open(filepath.Join(m.AreaPath, strconv.FormatUint(uint64(m.messageNums[position-1]), 10)+".msg"))
	if err != nil {
		return nil, err
	}
	defer f.Close()
	fi, err := f.Stat()
	if err != nil {
		return nil, err
	}
	msg, err := ioutil.ReadAll(f)
	if err != nil {
		return nil, err
	}
	msgb := bytes.NewBuffer(msg)
	var msgm msg_s
	err = utils.ReadStructFromBuffer(msgb, &msgm)
	if err != nil {
		return nil, err
	}
	rm := &Message{}
	rm.Area = m.AreaName
	rm.MsgNum = position
	rm.MaxNum = uint32(len(m.messageNums))
	rm.From = strings.Trim(string(msgm.From[:]), "\x00")
	rm.To = strings.Trim(string(msgm.To[:]), "\x00")
	rm.Subject = strings.Trim(string(msgm.Subj[:]), "\x00")
	rm.Body = strings.Trim(string(msgm.Body[:]), "\x00")
	rm.DateWritten, err = time.Parse("02 Jan 06  15:04:05", strings.Trim(string(msgm.Date[:]), "\x00"))
	rm.DateArrived = fi.ModTime()
	rm.Attrs = m.getAttrs(msgm.Attr)
	//  rm.Attr=uint32(msgm.Attr)
	err = rm.ParseRaw()
	if err != nil {
		return nil, err
	}
	//  tBody:=strings.Trim(string(msgm.Body[:]),"\x00")
	return rm, nil
	//return nil, errors.New("not implemented")
}

func (m *MSG) GetName() string {
	return m.AreaName
}

func (m *MSG) GetCount() uint32 {
	m.readMN()
	return uint32(len(m.messageNums))
}

func (m *MSG) GetLast() uint32 {
	m.readMN()
	file, err := os.Open(filepath.Join(m.AreaPath, "lastread"))
	if err != nil {
		return 0
	}
	b, err := ioutil.ReadAll(file)
	if len(b) != 2 {
		return 0
	}
	msgNum := uint32(binary.LittleEndian.Uint16(b))
	for i, is := range m.messageNums {
		if is == msgNum {
			return uint32(i + 1)
		}
	}
	return 0
}

func (m *MSG) readMN() {
	if len(m.messageNums) > 0 {
		return
	}
	fp, err := filepath.Glob(filepath.Join(m.AreaPath, "*.msg"))
	if err != nil {
		return
	}
	for _, fn := range fp {
		num, err := strconv.ParseUint(strings.TrimSuffix(filepath.Base(fn), ".msg"), 10, 32)
		if err == nil {
			m.messageNums = append(m.messageNums, uint32(num))
		} else {
			log.Print(err)
		}
	}
	sort.Slice(m.messageNums, func(i, j int) bool { return m.messageNums[i] < m.messageNums[j] })
}

func (m *MSG) GetType() EchoAreaType {
	return EchoAreaTypeMSG
}

func (m *MSG) SetLast(l uint32) {
	if l == 0 {
		l = 1
	}
	r := m.messageNums[l-1]
	buf := new(bytes.Buffer)
	err := binary.Write(buf, binary.LittleEndian, uint16(r))
	if err != nil {
		log.Print(err)
		return
	}
	err = ioutil.WriteFile(filepath.Join(m.AreaPath, "lastread"), buf.Bytes(), 0644)
	if err != nil {
		log.Print(err)
		return
	}
}

func (m *MSG) SaveMsg(tm *Message) error {
	return errors.New("not implemented")
}

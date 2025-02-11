package msgapi

import (
	"bytes"
	"encoding/binary"
	"io"

	//"errors"

	"log"
	"os"
	"path/filepath"
	"sort"
	"strconv"
	"strings"
	"time"

	"github.com/askovpen/gossiped/pkg/utils"
)

// MSG struct
type MSG struct {
	AreaPath string
	AreaName string
	AreaType EchoAreaType
	Chrs     string
	//	lastreads   string
	messageNums []uint32
	messages    []MessageListItem
}

type msgS struct {
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
	Attr        MSGAttrs
	Up          uint16
	Body        string
}

// MSGAttrs MSG attributes
type MSGAttrs uint16

// attributes
const (
	MSGPRIVATE MSGAttrs = 0x0001
	MSGCRASH   MSGAttrs = 0x0002
	MSGREAD    MSGAttrs = 0x0004
	MSGSENT    MSGAttrs = 0x0008
	MSGFILE    MSGAttrs = 0x0010
	MSGFWD     MSGAttrs = 0x0020
	MSGORPHAN  MSGAttrs = 0x0040
	MSGKILL    MSGAttrs = 0x0080
	MSGLOCAL   MSGAttrs = 0x0100
	MSGHOLD    MSGAttrs = 0x0200
	MSGXX2     MSGAttrs = 0x0400
	MSGFRQ     MSGAttrs = 0x0800
	MSGRRQ     MSGAttrs = 0x1000
	MSGCPT     MSGAttrs = 0x2000
	MSGARQ     MSGAttrs = 0x4000
	MSGURQ     MSGAttrs = 0x8000
)

// Init for future
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
		a >>= 1
	}
	return
}

func parseDate(date string) (ret time.Time) {
	ret, _ = time.Parse("02 Jan 06  15:04:05", date)
	return ret
}

// GetMsg getmsg
func (m *MSG) GetMsg(position uint32) (*Message, error) {
	if len(m.messageNums) == 0 {
		return nil, nil
	}
	if position == 0 {
		position = 1
	}
	f, err := os.Open(filepath.Join(m.AreaPath, strconv.FormatUint(uint64(m.messageNums[position-1]), 10)+".msg"))
	if err != nil {
		return nil, err
	}
	defer f.Close()
	msg, err := io.ReadAll(f)
	if err != nil {
		return nil, err
	}
	msgb := bytes.NewBuffer(msg)
	var msgm msgS
	err = utils.ReadStructFromBuffer(msgb, &msgm)
	if err != nil {
		return nil, err
	}
	rm := &Message{Area: m.AreaName,
		MsgNum:      position,
		MaxNum:      uint32(len(m.messageNums)),
		From:        strings.Trim(string(msgm.From[:]), "\x00"),
		To:          strings.Trim(string(msgm.To[:]), "\x00"),
		Subject:     strings.Trim(string(msgm.Subj[:]), "\x00"),
		Body:        strings.Trim(msgm.Body, "\x00"),
		DateWritten: parseDate(strings.Trim(string(msgm.Date[:]), "\x00")),
		DateArrived: getTime(msgm.DateArrived),
		Attrs:       m.getAttrs(uint16(msgm.Attr))}
	err = rm.ParseRaw()
	if err != nil {
		return nil, err
	}
	return rm, nil
}

// GetName get areaname
func (m *MSG) GetName() string {
	return m.AreaName
}

// GetCount get msg count
func (m *MSG) GetCount() uint32 {
	m.readMN()
	return uint32(len(m.messageNums))
}

// GetLast get last msg number
func (m *MSG) GetLast() uint32 {
	m.readMN()
	file, err := os.Open(filepath.Join(m.AreaPath, "lastread"))
	if err != nil {
		return 0
	}
	b, _ := io.ReadAll(file)
	if len(b) != 2 {
		return 0
	}
	msgNum := uint32(binary.LittleEndian.Uint16(b))
	for i, is := range m.messageNums {
		if is == msgNum {
			return uint32(i + 1)
		}
	}
	if msgNum != 0 {
		return uint32(len(m.messageNums))
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

// GetMsgType return area msg base type
func (m *MSG) GetMsgType() EchoAreaMsgType {
	return EchoAreaMsgTypeMSG
}

// GetType get area type
func (m *MSG) GetType() EchoAreaType {
	return m.AreaType
}

// SetLast set last message num
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
	err = os.WriteFile(filepath.Join(m.AreaPath, "lastread"), buf.Bytes(), 0644)
	if err != nil {
		log.Print(err)
		return
	}
}

// SaveMsg save message
func (m *MSG) SaveMsg(tm *Message) error {
	if _, err := os.Stat(m.AreaPath); os.IsNotExist(err) {
		err = os.MkdirAll(m.AreaPath, 0755)
		if err != nil {
			return err
		}
	}
	tm.Encode()
	msgm := msgS{Attr: MSGLOCAL,
		DateWritten: setTime(tm.DateWritten),
		DateArrived: setTime(tm.DateArrived),
		DestNode:    tm.ToAddr.GetNode(),
		DestNet:     tm.ToAddr.GetNet(),
		OrigNode:    tm.FromAddr.GetNode(),
		OrigNet:     tm.FromAddr.GetNet(),
		Body:        tm.Body}
	copy(msgm.From[:], tm.From)
	copy(msgm.To[:], tm.To)
	copy(msgm.Subj[:], tm.Subject)
	copy(msgm.Date[:], tm.DateWritten.Format("02 Jan 06  15:04:05"))
	for kl, v := range tm.Kludges {
		msgm.Body = "\x01" + kl + " " + v + "\x0d" + msgm.Body
	}
	msgm.Body += "\x00"
	buf := new(bytes.Buffer)
	err := utils.WriteStructToBuffer(buf, &msgm)
	if err != nil {
		return err
	}
	if len(m.messageNums) == 0 {
		err = os.WriteFile(
			filepath.Join(m.AreaPath, "1.msg"),
			buf.Bytes(),
			0644)
	} else {
		err = os.WriteFile(
			filepath.Join(m.AreaPath, strconv.FormatUint(uint64(m.messageNums[len(m.messageNums)-1]+1), 10)+".msg"),
			buf.Bytes(),
			0644)
	}
	if err != nil {
		return err
	}
	if len(m.messageNums) == 0 {
		m.messageNums = append(m.messageNums, 1)
	} else {
		m.messageNums = append(m.messageNums, m.messageNums[len(m.messageNums)-1]+1)
	}
	return nil
}

// SetChrs set charset
func (m *MSG) SetChrs(s string) {
	m.Chrs = s
}

// GetChrs get charset
func (m *MSG) GetChrs() string {
	return m.Chrs
}

// GetMessages get headers
func (m *MSG) GetMessages() *[]MessageListItem {
	if len(m.messages) > 0 || len(m.messageNums) == 0 {
		return &m.messages
	}
	for i := uint32(0); i < m.GetCount(); i++ {
		mm, err := m.GetMsg(i + 1)
		if err != nil {
			continue
		}
		m.messages = append(m.messages, MessageListItem{
			MsgNum:      i + 1,
			From:        mm.From,
			To:          mm.To,
			Subject:     mm.Subject,
			DateWritten: mm.DateWritten,
		})
	}
	return &m.messages
}

// DelMsg remove msg
func (m *MSG) DelMsg(l uint32) error {
	if l == 0 {
		l = 1
	}
	err := os.Remove(filepath.Join(m.AreaPath, strconv.FormatUint(uint64(m.messageNums[l-1]), 10)+".msg"))
	if err != nil {
		return err
	}
	if len(m.messages) == len(m.messageNums) {
		m.messages = append(m.messages[:l-1], m.messages[l:]...)
	}
	m.messageNums = append(m.messageNums[:l-1], m.messageNums[l:]...)
	return nil
}

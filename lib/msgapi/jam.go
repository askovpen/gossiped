package msgapi

import (
  "bufio"
  "bytes"
  "encoding/binary"
  "errors"
  "hash/crc32"
  "log"
  "os"
  "github.com/askovpen/goated/lib/config"
  "github.com/askovpen/goated/lib/types"
  "github.com/askovpen/goated/lib/utils"
  "strings"
  "time"
)

type JAM struct {
  AreaPath string
  AreaName string
  indexStructure []jam_s
  lastRead  []jam_l
}

type jam_s struct {
  ToCRC uint32
  Offset uint32
}

type jam_h struct {
  Signature uint32
  Revision uint16
  ReservedWord uint16
  SubfieldLen uint32
  TimesRead uint32
  MSGIDcrc uint32
  REPLYcrc uint32
  ReplyTo uint32
  Reply1st uint32
  ReplyNext uint32
  DateWritten uint32
  DateReceived uint32
  DateProcessed uint32
  MessageNumber uint32
  Attribute uint32
  Attribute2 uint32
  Offset uint32
  TxtLen uint32
  PasswordCRC uint32
  Cost uint32
}
type jam_l struct {
  UserCRC, UserId, LastReadMsg, HighReadMsg uint32
}

func (j *JAM) GetMsg(position uint32) (*Message, error) {
  if len(j.indexStructure)==0 { return nil, errors.New("Empty Area") }
  if position==0 {
    position=1
  }
  f, err := os.Open(j.AreaPath+".jhr")
  if err!=nil {
    return nil, err
  }
  defer f.Close()
  _, err = f.Seek(int64(j.indexStructure[position-1].Offset),0)
  if err!=nil {
    return nil, err
  }
  var header []byte
  header=make([]byte,76)
//  reader := bufio.NewReader(f)
  f.Read(header)
  headerb:=bytes.NewBuffer(header)
  var jamh jam_h
  if err=utils.ReadStructFromBuffer(headerb, &jamh); err!=nil {
    return nil, err
  }
  if jamh.Signature!=0x4d414a {return nil, errors.New("wrong message signature")}
  rm:=&Message{}
  rm.Area=j.AreaName
  rm.MsgNum=position
  rm.MaxNum=uint32(len(j.indexStructure))
//  _, tofs:=time.Now().Local().Zone()
  rm.DateWritten=time.Unix(int64(jamh.DateWritten),0)
  _, tofs:=rm.DateWritten.Zone()
  rm.DateArrived=time.Unix(int64(jamh.DateReceived),0)
  rm.DateWritten=rm.DateWritten.Add(time.Duration(tofs)* -time.Second)
  rm.DateArrived=rm.DateArrived.Add(time.Duration(tofs)* -time.Second)
  rm.Attr=jamh.Attribute
  rm.Body+=""
  var kl []byte
  kl=make([]byte,jamh.SubfieldLen)
  f.Read(kl)
  log.Printf("kl: %x", kl)
  klb:=bytes.NewBuffer(kl)
  for {
    var LoID,HiID uint16
    var datLen uint32
    err=binary.Read(klb, binary.LittleEndian, &LoID)
    if err!=nil {break}
    binary.Read(klb, binary.LittleEndian, &HiID)
    binary.Read(klb, binary.LittleEndian, &datLen)
    var val []byte
    val=make([]byte,datLen)
    binary.Read(klb, binary.LittleEndian, &val)
    log.Printf("%d, %d (%d): %s",LoID, HiID, datLen, val)
    switch LoID {
      case 0:
        rm.FromAddr=types.AddrFromString(string(val[:]))
      case 1:
        rm.ToAddr=types.AddrFromString(string(val[:]))
      case 2:
        rm.From=string(val[:])
      case 3:
        if crc32r(string(val[:]))!=j.indexStructure[position-1].ToCRC {
          return nil, errors.New("crc incorrect")
        }
        rm.To=string(val[:])
      case 4:
        if crc32r(string(val[:]))!=jamh.MSGIDcrc {
          return nil, errors.New("crc incorrect")
        }
        rm.Body+="\x01MSGID: "+string(val[:])+"\x0d"
      case 5:
        if crc32r(string(val[:]))!=jamh.REPLYcrc {
          return nil, errors.New("crc incorrect")
        }
        rm.Body+="\x01REPLYID: "+string(val[:])+"\x0d"
      case 6:
        rm.Subject=string(val[:])
      case 7:
        rm.Body+="\x01PID: "+string(val[:])+"\x0d"
      default:
        rm.Body+="\x01"+string(val[:])+"\x0d"
    }
  }
  f, err = os.Open(j.AreaPath+".jdt")
  if err!=nil {
    return nil, err
  }
  f.Seek(int64(jamh.Offset), 0)
  defer f.Close()
  var txt []byte
  txt=make([]byte,jamh.TxtLen)
  f.Read(txt)
  rm.Body+=string(txt[:])
  err=rm.ParseRaw()
  if err!=nil {
    return nil, err
  }
  //log.Printf("msgh: %#v", jamh)
  //log.Printf("rm: %#v", rm)
  
  return rm, nil
}
func (j *JAM) readJDX() {
  if len(j.indexStructure)>0 {
    return
  }
  file, err := os.Open(j.AreaPath+".jdx")
  if err!=nil {
    return
  }
  defer file.Close()
  reader := bufio.NewReader(file)
  part := make([]byte, 16384)
  for {
    count, err := reader.Read(part);
    if err!=nil {
      break
    }
    partb:=bytes.NewBuffer(part[:count])
    for {
      var jam jam_s
      if err=utils.ReadStructFromBuffer(partb, &jam); err!=nil {
        break
      }
      if (jam.ToCRC!=0xffffffff || jam.Offset!=0xffffffff) {
        j.indexStructure=append(j.indexStructure,jam)
      }
    }
  }
}

func (j *JAM) readJLR() {
  if len(j.lastRead)>0 {
    return
  }
  file, err := os.Open(j.AreaPath+".jlr")
  if err!=nil {
    return
  }
  defer file.Close()
  reader := bufio.NewReader(file)
  part := make([]byte, 16384)
  for {
    count, err := reader.Read(part);
    if err!=nil {
      break
    }
    partb:=bytes.NewBuffer(part[:count])
    for {
      var jaml jam_l 
      if err=utils.ReadStructFromBuffer(partb, &jaml); err!=nil {
        break
      }
      j.lastRead=append(j.lastRead,jaml)
    }
  }
  log.Printf("%#v", j.lastRead)
}

func (j *JAM) GetLast() uint32 {
  j.readJLR()
  for _,l:=range j.lastRead {
    if l.UserCRC==crc32r(config.Config.Username) {
      return l.LastReadMsg
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
  bstr:=[]byte(strings.ToLower(str))
  return 0xffffffff-crc32.ChecksumIEEE(bstr)
}

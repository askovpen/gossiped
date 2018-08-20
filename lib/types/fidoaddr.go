package types

import (
	"errors"
	//  "fmt"
	"regexp"
	"strconv"
	//  "strings"
	//  "log"
)

type FidoAddr struct {
	zone  uint16
	net   uint16
	node  uint16
	point uint16
}

func (f *FidoAddr) Equal(fn FidoAddr) bool {
	if f.zone == fn.zone && f.net == fn.net && f.node == fn.node && f.point == fn.point {
		return true
	} else {
		return false
	}
}
func (f *FidoAddr) String() string {
	if f.zone == 0 {
		return ""
	}
	if f.point == 0 {
		return strconv.Itoa(int(f.zone)) + ":" + strconv.Itoa(int(f.net)) + "/" + strconv.Itoa(int(f.node))
	}
	return strconv.Itoa(int(f.zone)) + ":" + strconv.Itoa(int(f.net)) + "/" + strconv.Itoa(int(f.node)) + "." + strconv.Itoa(int(f.point))
}
func (f *FidoAddr) FQDN() (string, error) {
	if f.point > 0 {
		return "", errors.New("point")
	}
	return "f" + strconv.Itoa(int(f.node)) + ".n" + strconv.Itoa(int(f.net)) + ".z" + strconv.Itoa(int(f.zone)) + ".binkp.net", nil
}
func AddrFromString(s string) *FidoAddr {
	f := &FidoAddr{}
	res := regexp.MustCompile("(\\d+):(\\d+)/(\\d+)\\.?(\\d+)?(@.*)?").FindStringSubmatch(s)
	if len(res) == 0 {
		return nil
	}
	if len(res[1]) > 0 {
		zone, _ := strconv.Atoi(res[1])
		f.zone = uint16(zone)
	}
	if len(res[2]) > 0 {
		net, _ := strconv.Atoi(res[2])
		f.net = uint16(net)
	}
	if len(res[3]) > 0 {
		node, _ := strconv.Atoi(res[3])
		f.node = uint16(node)
	}
	if len(res[4]) > 0 {
		point, _ := strconv.Atoi(res[4])
		f.point = uint16(point)
	}
	return f
}
func AddrFromNum(zone uint16, net uint16, node uint16, point uint16) *FidoAddr {
	f := &FidoAddr{}
	f.zone = zone
	f.net = net
	f.node = node
	f.point = point
	return f
}

func (f *FidoAddr) UnmarshalYAML(unmarshal func(interface{}) error) (err error) {
	var fm string
	if err := unmarshal(&fm); err != nil {
		return err
	}
	tf := AddrFromString(fm)
	f.zone = tf.zone
	f.net = tf.net
	f.node = tf.node
	f.point = tf.point
	return nil
}
func (f FidoAddr) MarshalYAML() (interface{}, error) {
	return f.String(), nil
}


func (f *FidoAddr) GetZone() uint16 {
	return f.zone
}

func (f *FidoAddr) GetNode() uint16 {
	return f.node
}

func (f *FidoAddr) GetNet() uint16 {
	return f.net
}

func (f *FidoAddr) GetPoint() uint16 {
	return f.point
}

func (f *FidoAddr) SetPoint(p uint16) {
	f.point = p
}

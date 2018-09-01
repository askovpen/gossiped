package types

import (
	"errors"
	"regexp"
	"strconv"
)

// FidoAddr struct
type FidoAddr struct {
	zone  uint16
	net   uint16
	node  uint16
	point uint16
}

// Equal compare two *FidoAddr
func (f *FidoAddr) Equal(fn *FidoAddr) bool {
	if f.zone == fn.zone && f.net == fn.net && f.node == fn.node && f.point == fn.point {
		return true
	}
	return false
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

// FQDN return hostname
func (f *FidoAddr) FQDN() (string, error) {
	if f.point > 0 {
		return "", errors.New("point")
	}
	return "f" + strconv.Itoa(int(f.node)) + ".n" + strconv.Itoa(int(f.net)) + ".z" + strconv.Itoa(int(f.zone)) + ".binkp.net", nil
}

// AddrFromString return FidoAddr from string
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

// AddrFromNum return FidoAddr from digits
func AddrFromNum(zone uint16, net uint16, node uint16, point uint16) *FidoAddr {
	f := &FidoAddr{}
	f.zone = zone
	f.net = net
	f.node = node
	f.point = point
	return f
}

// UnmarshalYAML for UnmarshalYAML
func (f *FidoAddr) UnmarshalYAML(unmarshal func(interface{}) error) (err error) {
	var fm string
	if err := unmarshal(&fm); err != nil {
		return err
	}
	tf := AddrFromString(fm)
	if tf==nil {
		return errors.New("wrong address")
	}
	f.zone = tf.zone
	f.net = tf.net
	f.node = tf.node
	f.point = tf.point
	return nil
}

// MarshalYAML for MarshaYAML
func (f FidoAddr) MarshalYAML() (interface{}, error) {
	return f.String(), nil
}

// GetZone return zone
func (f *FidoAddr) GetZone() uint16 {
	return f.zone
}

// GetNode return node
func (f *FidoAddr) GetNode() uint16 {
	return f.node
}

// GetNet return net
func (f *FidoAddr) GetNet() uint16 {
	return f.net
}

// GetPoint return point
func (f *FidoAddr) GetPoint() uint16 {
	return f.point
}

// SetPoint set point
func (f *FidoAddr) SetPoint(p uint16) *FidoAddr {
	f.point = p
	return f
}

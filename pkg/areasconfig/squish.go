package areasconfig

import (
	"bufio"
	"io"
	"os"
	"regexp"
	"strings"

	"github.com/askovpen/gossiped/pkg/msgapi"
)

func squishConfigRead(fn string) error {
	re := regexp.MustCompile(`[^\s\t"']+|"([^"]*)"|'([^']*)`)
	file, err := os.Open(fn)
	if err != nil {
		return err
	}
	b, err := io.ReadAll(file)
	if err != nil {
		return err
	}
	scanner := bufio.NewScanner(strings.NewReader(string(b)))
	for scanner.Scan() {
		res := re.FindAllString(scanner.Text(), -1)
		if len(res) < 3 {
			continue
		}
		amType := getSquishAreaType(res)
		if amType == msgapi.EchoAreaMsgTypePasstrough {
			continue
		}
		var aType msgapi.EchoAreaType
		//aType := msgapi.EchoAreaTypeNone
		if strings.EqualFold(res[0], "echoarea") {
			aType = msgapi.EchoAreaTypeEcho
		} else if strings.EqualFold(res[0], "netarea") {
			aType = msgapi.EchoAreaTypeNetmail
		} else if strings.EqualFold(res[0], "badarea") {
			aType = msgapi.EchoAreaTypeBad
		} else if strings.EqualFold(res[0], "dupearea") {
			aType = msgapi.EchoAreaTypeDupe
		} else if strings.EqualFold(res[0], "localarea") {
			aType = msgapi.EchoAreaTypeLocal
		} else {
			continue
		}
		switch amType {
		case msgapi.EchoAreaMsgTypeSquish:
			area := &msgapi.Squish{AreaName: res[1], AreaPath: res[2], AreaType: aType}
			msgapi.Areas = append(msgapi.Areas, area)
		case msgapi.EchoAreaMsgTypeMSG:
			area := &msgapi.MSG{AreaName: res[1], AreaPath: res[2], AreaType: aType}
			msgapi.Areas = append(msgapi.Areas, area)
		}
	}
	return nil
}

func getSquishAreaType(tokens []string) msgapi.EchoAreaMsgType {
	for _, t := range tokens {
		if t == "-$" {
			return msgapi.EchoAreaMsgTypeSquish
		} else if t == "-0" {
			return msgapi.EchoAreaMsgTypePasstrough
		}
	}
	return msgapi.EchoAreaMsgTypeMSG
}

package areasconfig

import (
	"bufio"
	"github.com/askovpen/goated/pkg/msgapi"
	"io/ioutil"
	"os"
	"regexp"
	"strings"
)

func areasbbsConfigRead(fn string) error {
	re := regexp.MustCompile(`[^\s\t"']+|"([^"]*)"|'([^']*)`)
	file, err := os.Open(fn)
	if err != nil {
		return err
	}
	b, err := ioutil.ReadAll(file)
	if err != nil {
		return err
	}
	scanner := bufio.NewScanner(strings.NewReader(string(b[:])))
	for scanner.Scan() {
		res := re.FindAllString(scanner.Text(), -1)
		if len(res) < 2 {
			continue
		}
		if len(res[0]) < 3 {
			continue
		}
		switch res[0][0] {
		case ';':
			continue
		case '$':
			area := &msgapi.Squish{AreaName: res[1], AreaPath: res[0][1:], AreaType: msgapi.EchoAreaTypeEcho}
			msgapi.Areas = append(msgapi.Areas, area)
		case '!':
			area := &msgapi.JAM{AreaName: res[1], AreaPath: res[0][1:], AreaType: msgapi.EchoAreaTypeEcho}
			msgapi.Areas = append(msgapi.Areas, area)
		default:
			area := &msgapi.MSG{AreaName: res[1], AreaPath: res[0], AreaType: msgapi.EchoAreaTypeEcho}
			msgapi.Areas = append(msgapi.Areas, area)
		}
	}
	return nil
}

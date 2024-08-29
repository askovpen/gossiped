package areasconfig

import (
	"errors"
	"github.com/askovpen/gossiped/pkg/config"
	"github.com/askovpen/gossiped/pkg/msgapi"

	"strings"
)

// Read area configs
func Read() error {
	// log.Printf(config.Config.AreaFile.Type)
	var err error
	switch config.Config.AreaFile.Type {
	case "fidoconfig":
		err = fidoConfigRead(config.Config.AreaFile.Path)
	case "areas.bbs":
		err = areasbbsConfigRead(config.Config.AreaFile.Path)
	case "squish":
		err = squishConfigRead(config.Config.AreaFile.Path)
	case "crashmail":
		err = crashmailConfigRead(config.Config.AreaFile.Path)
	default:
		return errors.New("unknown AreasConfig.Type '" + config.Config.AreaFile.Type + "'")
	}
	if err != nil {
		return nil
	}
	for i := range config.Config.Areas {
		found := false
		for _, da := range msgapi.Areas {
			if config.Config.Areas[i].Name == da.GetName() {
				found = true
				if config.Config.Areas[i].Chrs != "" {
					da.SetChrs(config.Config.Areas[i].Chrs)
				}
			}
		}
		if !found {
			a, err := getArea(i)
			if err == nil {
				msgapi.Areas = append(msgapi.Areas, a)
			}
		}
	}

	if len(msgapi.Areas) == 0 {
		return errors.New("no Areas found")
	}
	return nil
}

func getArea(i int) (msgapi.AreaPrimitive, error) {
	switch config.Config.Areas[i].BaseType {
	case "msg":
		r := &msgapi.MSG{AreaName: config.Config.Areas[i].Name, AreaPath: config.Config.Areas[i].Path, AreaType: getType(config.Config.Areas[i].Type)}
		if config.Config.Areas[i].Chrs != "" {
			r.Chrs = config.Config.Areas[i].Chrs
		}
		return r, nil
	case "squish":
		r := &msgapi.Squish{AreaName: config.Config.Areas[i].Name, AreaPath: config.Config.Areas[i].Path, AreaType: getType(config.Config.Areas[i].Type)}
		if config.Config.Areas[i].Chrs != "" {
			r.Chrs = config.Config.Areas[i].Chrs
		}
		return r, nil
	case "jam":
		r := &msgapi.JAM{AreaName: config.Config.Areas[i].Name, AreaPath: config.Config.Areas[i].Path, AreaType: getType(config.Config.Areas[i].Type)}
		if config.Config.Areas[i].Chrs != "" {
			r.Chrs = config.Config.Areas[i].Chrs
		}
		return r, nil
	}
	return nil, errors.New("uknown type")
}
func getType(t string) msgapi.EchoAreaType {
	if strings.EqualFold(t, "echo") {
		return msgapi.EchoAreaTypeEcho
	} else if strings.EqualFold(t, "local") {
		return msgapi.EchoAreaTypeLocal
	} else if strings.EqualFold(t, "netmail") {
		return msgapi.EchoAreaTypeNetmail
	} else if strings.EqualFold(t, "dupe") {
		return msgapi.EchoAreaTypeDupe
	} else if strings.EqualFold(t, "bad") {
		return msgapi.EchoAreaTypeBad
	}
	return msgapi.EchoAreaTypeLocal
}

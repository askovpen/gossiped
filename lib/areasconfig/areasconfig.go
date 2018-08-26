package areasconfig

import (
	"errors"
	"github.com/askovpen/goated/lib/config"
)

// Read area configs
func Read() error {
	switch config.Config.AreaFile.Type {
	case "fidoconfig":
		return fidoConfigRead(config.Config.AreaFile.Path)
	}
	return errors.New("unknown AreasConfig.Type '" + config.Config.AreaFile.Type + "'")
}

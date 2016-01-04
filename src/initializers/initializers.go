package initializers

import (
	"path/filepath"
)

type Paths struct {
	ConfigFilePath  string
	LogFilePath     string
	SoaRegistryPath string
}

var ConfigPaths *Paths

func init() {
	ConfPath, _ := filepath.Abs("../gatekeeper/config/config.json")
	LogPath, _ := filepath.Abs("../gatekeeper/log/h-gatekeeper.log")
	SoaPath, _ := filepath.Abs("../gatekeeper/config/soa.json")

	ConfigPaths = &Paths{ConfPath, LogPath, SoaPath}
}

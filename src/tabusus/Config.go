package tabusus

import (
	hocon "github.com/btnguyen2k/configuration"
	"github.com/labstack/gommon/log"
	"os"
	"path"
)

// HoconConfig encapsulates application's configurations in HOCON format
type HoconConfig struct {
	File string        // config file
	Conf *hocon.Config // configurations
}

type FileInfoList []os.FileInfo

func (s FileInfoList) Len() int {
	return len(s)
}

func (s FileInfoList) Less(i, j int) bool {
	return s[i].Name() < s[j].Name()
}

func (s FileInfoList) Swap(i, j int) {
	s[i], s[j] = s[j], s[i]
}

func LoadAppConfig(file string) *HoconConfig {
	dir, err := os.Getwd()
	if err != nil {
		panic(err)
	}
	defer os.Chdir(dir)

	config := HoconConfig{}
	log.Info("Loading configurations in file [", file, "]")
	confDir, confFile := path.Split(file)
	os.Chdir(confDir)
	config.File = file
	config.Conf = hocon.LoadConfig(confFile)
	return &config
}

package conf

import (
	"os"

	"github.com/k0mmsussert0d/fukaeri/internal"
	"gopkg.in/yaml.v2"
)

type ConfigurationDB struct {
	Connstring string
	Name       string
}

type ConfigurationArchive struct {
	Boards []string
}

type Configuration struct {
	DB      ConfigurationDB
	Archive ConfigurationArchive
}

var _config *Configuration = nil

func GetConfig() *Configuration {
	if _config == nil {
		parseConfig()
	}

	return _config
}

func GetNewConfig() *Configuration {
	parseConfig()
	return _config
}

func parseConfig() {
	data, err := os.ReadFile("./conf.yml")
	internal.HandleError(err)

	err = yaml.Unmarshal(data, &_config)
	internal.HandleError(err)
}

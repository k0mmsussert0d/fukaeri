package conf

import (
	"log"
	"os"

	"gopkg.in/yaml.v2"
)

type ConfigurationDB struct {
	Connstring string
	Name       string
	Files      *string
}

type ConfigurationArchive struct {
	Boards []string
}

type Configuration struct {
	DB       ConfigurationDB
	Archive  ConfigurationArchive
	LogLevel string
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
	if err != nil {
		log.Panic(err.Error())
	}

	err = yaml.Unmarshal(data, &_config)
	if err != nil {
		log.Panic(err.Error())
	}
}

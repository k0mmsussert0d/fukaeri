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
	LogLevel string `yaml:"log_level"`
}

func GetConfig() *Configuration {
	var config *Configuration
	data, err := os.ReadFile("./conf.yml")
	if err != nil {
		log.Panic(err.Error())
	}

	err = yaml.Unmarshal(data, &config)
	if err != nil {
		log.Panic(err.Error())
	}

	return config
}

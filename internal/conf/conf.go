package conf

import (
	"log"
	"os"

	"gopkg.in/yaml.v2"
)

type ConfigurationDB struct {
	Connstring string
	Name       string
}

type Configuration struct {
	DB ConfigurationDB
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
	config := Configuration{}
	data, err := os.ReadFile("./conf.yml")
	check(err)

	err = yaml.Unmarshal(data, &config)
	check(err)
}

func check(e error) {
	if e != nil {
		log.Fatal(e)
	}
}

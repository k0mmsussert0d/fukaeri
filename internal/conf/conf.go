package conf

import (
	"fmt"

	"github.com/fsnotify/fsnotify"
	"github.com/spf13/viper"
)

type Configuration struct {
	DB struct {
		Connstring string
		Collection string
		Files      string
	}
	Archive struct {
		Boards []string
	}
}

var c Configuration

func Init() {
	viper.SetConfigName("conf")
	viper.SetConfigType("yml")
	viper.AddConfigPath(".")
	viper.SetDefault("db.collection", "fukaeri")
	viper.SetDefault("db.files", "fukaeri_files")

	err := viper.ReadInConfig()
	if err != nil {
		panic(fmt.Errorf("failed to read config file, %v", err))
	}
	viper.Unmarshal(&c)
	viper.OnConfigChange(func(in fsnotify.Event) {
		viper.Unmarshal(&c)
	})
	viper.WatchConfig()
}

func Get() Configuration {
	return c
}

package conf

import (
	"fmt"

	"github.com/fsnotify/fsnotify"
	"github.com/spf13/viper"
	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"
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
	Zap zap.Config
}

var c Configuration

func Init() {
	viper.SetConfigName("conf")
	viper.SetConfigType("yml")
	viper.AddConfigPath(".")
	viper.SetDefault("db.collection", "fukaeri")
	viper.SetDefault("db.files", "fukaeri_files")
	viper.SetDefault("zap", zap.Config{
		Level:            zap.NewAtomicLevelAt(zap.InfoLevel),
		Encoding:         "json",
		OutputPaths:      []string{"stdout"},
		ErrorOutputPaths: []string{"stderr"},
		EncoderConfig: zapcore.EncoderConfig{
			EncodeLevel:  zapcore.LowercaseLevelEncoder,
			LevelKey:     "level",
			EncodeTime:   zapcore.EpochTimeEncoder,
			TimeKey:      "time",
			EncodeCaller: zapcore.FullCallerEncoder,
			CallerKey:    "line",
			MessageKey:   "message",
		},
	})

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

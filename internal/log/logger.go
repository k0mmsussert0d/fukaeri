package log

import (
	"github.com/k0mmsussert0d/fukaeri/internal/conf"
	"go.uber.org/zap"
)

var logger zap.SugaredLogger

func Init() {
	cfg := conf.Get().Zap
	logger = *zap.Must(cfg.Build()).Sugar()
}

func Logger() *zap.SugaredLogger {
	return &logger
}

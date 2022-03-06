package log

import (
	"io"
	"log"
	"os"

	"github.com/k0mmsussert0d/fukaeri/internal/conf"
)

var (
	debug   *log.Logger
	info    *log.Logger
	warning *log.Logger
	err     *log.Logger
)

func Auto() {
	switch level := conf.GetConfig().LogLevel; level {
	case "ERROR":
		Init(io.Discard, io.Discard, io.Discard, os.Stderr)
	case "WARN":
		Init(io.Discard, io.Discard, os.Stdout, os.Stderr)
	case "INFO":
		Init(io.Discard, os.Stdout, os.Stdout, os.Stderr)
	case "DEBUG":
		Init(os.Stdout, os.Stdout, os.Stdout, os.Stderr)
	default:
		log.Panicf("Unrecognized log level option %v", level)
	}

}

func Init(
	debugHandle io.Writer,
	infoHandle io.Writer,
	warningHandle io.Writer,
	errHandle io.Writer,
) {

	flags := log.Ldate | log.Ltime | log.Lmicroseconds | log.Llongfile | log.Lmsgprefix
	debug = log.New(debugHandle, "DEBUG: ", flags)
	info = log.New(infoHandle, "INFO: ", flags)
	warning = log.New(warningHandle, "WARN: ", flags)
	err = log.New(errHandle, "ERROR: ", flags)
}

func Debug() *log.Logger {
	return debug
}

func Info() *log.Logger {
	return info
}

func Warn() *log.Logger {
	return warning
}

func Error() *log.Logger {
	return err
}

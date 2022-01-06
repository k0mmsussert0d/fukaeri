package log

import (
	"io"
	"log"
)

var (
	debug   *log.Logger
	info    *log.Logger
	warning *log.Logger
	err     *log.Logger
)

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

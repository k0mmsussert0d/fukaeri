package internal

import (
	"context"
	"errors"
	"fmt"
	"runtime/debug"

	"github.com/k0mmsussert0d/fukaeri/internal/log"
)

type FukaeriError struct {
	Inner      error
	Message    string
	StackTrace string
	Misc       map[string]interface{}
}

func (ferror FukaeriError) Error() string {
	return ferror.Message
}

func (ferror FukaeriError) Unwrap() error {
	return ferror.Inner
}

func HandleError(err error) {
	if err != nil {
		ferr := WrapError(err, err.Error())
		if errors.Is(err, context.Canceled) { // silence context cancelation
			log.Debug().Println(ferr.Message)
		} else {
			log.Error().Println(ferr.Message)
			log.Error().Println(ferr.StackTrace)
		}
		panic(ferr)
	}
}

func ConvertPanicToError(panicVal interface{}) error {
	switch x := panicVal.(type) {
	case string:
		return WrapError(errors.New(x), x)
	case error:
		return WrapError(x, x.Error())
	default:
		msg := "Unknown panic type"
		return WrapError(errors.New(msg), msg)
	}
}

func WrapError(err error, msg string, msgArgs ...interface{}) FukaeriError {
	switch x := err.(type) {
	case *FukaeriError:
		return FukaeriError{
			Inner:      x,
			Message:    fmt.Sprintf(msg, msgArgs...),
			StackTrace: fmt.Sprintf("%s\nCaused by inner error:\n%s", string(debug.Stack()), x.StackTrace),
			Misc:       make(map[string]interface{}),
		}
	default:
		return FukaeriError{
			Inner:      x,
			Message:    fmt.Sprintf(msg, msgArgs...),
			StackTrace: string(debug.Stack()),
			Misc:       make(map[string]interface{}),
		}
	}
}

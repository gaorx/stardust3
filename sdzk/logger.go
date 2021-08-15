package sdzk

import (
	"fmt"

	"github.com/gaorx/stardust3/sdlog"
)

type sdLogger struct {
}

type discardLogger struct {
}

var (
	SdLog   = sdLogger{}
	Discard = discardLogger{}
)

func (sdLogger) Printf(format string, args ...interface{}) {
	sdlog.Debug(fmt.Sprintf(format, args...))
}

func (discardLogger) Printf(format string, args ...interface{}) {
}

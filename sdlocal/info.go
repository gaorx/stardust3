package sdlocal

import (
	"os"
	"runtime"

	"github.com/gaorx/stardust3/sderr"
	"github.com/gaorx/stardust3/sdlog"
)

func Hostname() string {
	hn, err := os.Hostname()
	if err != nil {
		sdlog.WithError(sderr.WithStack(err)).Warn("get hostname error")
		return ""
	}
	return hn
}

// 返回当前操作系统名称 linux,windows,darwin,openbsd,freebsd,android...
func OS() string {
	return runtime.GOOS
}

func Arch() string {
	return runtime.GOARCH
}

func NumCPU() int {
	return runtime.NumCPU()
}

func NumGoroutine() int {
	return runtime.NumGoroutine()
}

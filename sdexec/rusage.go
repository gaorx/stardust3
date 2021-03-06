//go:build !windows
// +build !windows

package sdexec

import (
	"syscall"
)

type Rusage syscall.Rusage

func getRusage(ru interface{}) Rusage {
	if ru1 := ru.(*syscall.Rusage); ru != nil {
		rur := *ru1
		return Rusage(rur)
	} else {
		return Rusage{}
	}
}

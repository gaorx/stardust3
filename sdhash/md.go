package sdhash

import (
	"crypto/md5"

	"github.com/gaorx/stardust3/sdbytes"
)

func Md5(data []byte) sdbytes.Packet {
	sum := md5.Sum(data)
	return sum[:]
}

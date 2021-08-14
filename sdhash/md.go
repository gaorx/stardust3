package sdhash

import (
	"bytes"
	"crypto/md5"

	"github.com/gaorx/stardust3/sdencoding"
)

func Md5(data []byte) sdencoding.Bytes {
	sum := md5.Sum(data)
	return sum[:]
}

func ValidMd5(data, expected []byte) bool {
	sum := md5.Sum(data)
	return bytes.Equal(sum[:], expected)
}

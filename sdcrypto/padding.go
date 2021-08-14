package sdcrypto

import (
	"bytes"

	"github.com/gaorx/stardust3/sdbytes"
	"github.com/gaorx/stardust3/sderr"
)

type Padding func(data []byte, blockSize int) ([]byte, error)
type Unpadding func(data []byte, blockSize int) ([]byte, error)

func Pkcs5(data []byte, blockSize int) ([]byte, error) {
	if blockSize <= 0 {
		return nil, sderr.New("illegal block size")
	}
	padding := blockSize - len(data)%blockSize
	padded := bytes.Repeat([]byte{byte(padding)}, padding)
	return append(sdbytes.Copy(data), padded...), nil
}

func UnPkcs5(data []byte, blockSize int) ([]byte, error) {
	if blockSize <= 0 {
		return nil, sderr.New("illegal block size")
	}
	if len(data) < blockSize {
		return nil, sderr.New("data too short")
	}
	lastByte := int(data[len(data)-1])
	if lastByte <= 0 || lastByte > blockSize {
		return nil, sderr.New("illegal padding size")
	}
	return sdbytes.Copy(data[:len(data)-lastByte]), nil
}

func Zeros(data []byte, blockSize int) ([]byte, error) {
	if blockSize <= 0 {
		return nil, sderr.New("illegal block size")
	}
	padding := blockSize - len(data)%blockSize
	padded := bytes.Repeat([]byte{0}, padding)
	result := make([]byte, 0, len(data)+padding)
	result = append(result, data...)
	return append(result, padded...), nil
}

func UnZeros(data []byte, blockSize int) ([]byte, error) {
	return bytes.TrimRightFunc(data,
		func(r rune) bool {
			return r == 0
		}), nil
}

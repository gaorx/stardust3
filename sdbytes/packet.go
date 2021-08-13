package sdbytes

import (
	"encoding/base64"
	"encoding/hex"
	"strings"
)

type Packet []byte

func (p Packet) ToBytes() []byte {
	return []byte(p)
}

func (p Packet) HexL() string {
	if len(p) <= 0 {
		return ""
	}
	return hex.EncodeToString(p)
}

func (p Packet) HexU() string {
	if len(p) <= 0 {
		return ""
	}
	return strings.ToUpper(hex.EncodeToString(p))
}

func (p Packet) Base64Std() string {
	if len(p) <= 0 {
		return ""
	}
	return base64.StdEncoding.EncodeToString(p)
}

func (p Packet) Base64Url() string {
	if len(p) <= 0 {
		return ""
	}
	return base64.URLEncoding.EncodeToString(p)
}

func DecodeHex(s string) (Packet, error) {
	return hex.DecodeString(s)
}

func MustDecodeHex(s string) Packet {
	data, err := DecodeHex(s)
	if err != nil {
		panic(err)
	}
	return data
}

// TODO: Decode Base64

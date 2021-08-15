package sduuid

import (
	"github.com/gaorx/stardust3/sdencoding"
	uuid "github.com/satori/go.uuid"
)

type (
	UUID = uuid.UUID
)

func Encode(id UUID) sdencoding.Bytes {
	return id.Bytes()
}

func NewV1() sdencoding.Bytes {
	return uuid.NewV1().Bytes()
}

func NewV4() sdencoding.Bytes {
	return uuid.NewV4().Bytes()
}

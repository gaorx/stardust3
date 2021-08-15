package sdsqlx

import (
	"database/sql/driver"

	"github.com/gaorx/stardust3/sderr"
)

type JsonStrings []string

func (j JsonStrings) Value() (driver.Value, error) {
	return JsonValue(j, []byte("[]"))
}

func (j *JsonStrings) Scan(value interface{}) error {
	if value == nil {
		value = []byte("[]")
	}
	var a []string
	err := JsonScan(value, &a)
	if err != nil {
		return sderr.WithStack(err)
	}
	if a == nil {
		a = []string{}
	}
	*j = JsonStrings(a)
	return nil
}

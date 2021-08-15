package sdsqlx

import (
	"database/sql/driver"

	"github.com/gaorx/stardust3/sderr"
)

type JsonInts []int

func (j JsonInts) Value() (driver.Value, error) {
	return JsonValue(j, []byte("[]"))
}

func (j *JsonInts) Scan(value interface{}) error {
	if value == nil {
		value = []byte("[]")
	}
	var a []int
	err := JsonScan(value, &a)
	if err != nil {
		return sderr.WithStack(err)
	}
	if a == nil {
		a = []int{}
	}
	*j = JsonInts(a)
	return nil
}

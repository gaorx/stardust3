package sdsqlx

import (
	"database/sql/driver"

	"github.com/gaorx/stardust3/sderr"
	"github.com/gaorx/stardust3/sdjson"
)

type JsonArray sdjson.Array

func (j JsonArray) Value() (driver.Value, error) {
	return JsonValue(j, []byte("[]"))
}

func (j *JsonArray) Scan(value interface{}) error {
	if value == nil {
		value = []byte("[]")
	}
	var a sdjson.Array
	err := JsonScan(value, &a)
	if err != nil {
		return sderr.WithStack(err)
	}
	if a == nil {
		a = sdjson.Array{}
	}
	*j = JsonArray(a)
	return nil
}

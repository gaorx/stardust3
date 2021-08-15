package sdsqlx

import (
	"database/sql/driver"

	"github.com/gaorx/stardust3/sderr"
	"github.com/gaorx/stardust3/sdjson"
)

type JsonObject sdjson.Object

func (j JsonObject) Value() (driver.Value, error) {
	return JsonValue(j, []byte("{}"))
}

func (j *JsonObject) Scan(value interface{}) error {
	if value == nil {
		value = []byte("{}")
	}
	var o sdjson.Object
	err := JsonScan(value, &o)
	if err != nil {
		return sderr.WithStack(err)
	}
	if o == nil {
		o = sdjson.Object{}
	}
	*j = JsonObject(o)
	return nil
}

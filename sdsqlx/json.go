package sdsqlx

import (
	"database/sql/driver"
	"encoding/json"

	"github.com/gaorx/stardust3/sderr"
)

func JsonValue(v interface{}, nilAs []byte) (driver.Value, error) {
	if v == nil {
		return nilAs, nil
	}
	r, err := json.Marshal(v)
	if err != nil {
		return nil, sderr.WithStack(err)
	}
	return r, nil
}

func JsonScan(dbVal interface{}, toVal interface{}) error {
	buff, ok := dbVal.([]byte)
	if !ok {
		return sderr.Newf("failed to unmarshal jsonb value: %v", dbVal)
	}
	return sderr.WithStack(json.Unmarshal(buff, toVal))
}

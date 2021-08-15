package sdsqlx

import (
	"github.com/gaorx/stardust3/sderr"
	"github.com/jmoiron/sqlx"
)

type ColConverter func(col string, v interface{}) (interface{}, error)

type mapScan interface {
	MapScan(dest map[string]interface{}) error
}

func ToMap(row mapScan, dest map[string]interface{}, converter ColConverter) error {
	if row == nil {
		return sderr.New("nil row")
	}
	if row1, ok := row.(*sqlx.Row); ok && row1 == nil {
		return sderr.New("nil row")
	}
	if rows1, ok := row.(*sqlx.Rows); ok && rows1 == nil {
		return sderr.New("nil rows")
	}
	m := map[string]interface{}{}
	err := row.MapScan(m)
	if err != nil {
		return sderr.WithStack(err)
	}
	if converter == nil {
		converter = DefaultColConverter
	}
	m1 := make(map[string]interface{}, len(m))
	for k, v := range m {
		v1, err := converter(k, v)
		if err != nil {
			return sderr.Newf("convert column %s error", k)
		}
		m1[k] = v1
	}
	for k, v := range m1 {
		dest[k] = v
	}
	return nil
}

func DefaultColConverter(col string, v interface{}) (interface{}, error) {
	if v == nil {
		return nil, nil
	}
	switch v1 := v.(type) {
	case []byte:
		return string(v1), nil
	default:
		return v, nil
	}
}

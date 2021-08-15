package sdsqlx

import (
	"github.com/gaorx/stardust3/sderr"
	"github.com/jmoiron/sqlx"
)

type Address struct {
	Driver string `json:"driver" toml:"driver"`
	Addr   string `json:"addr" toml:"addr"`
}

func Dial(addr Address) (*sqlx.DB, error) {
	db, err := sqlx.Connect(addr.Driver, addr.Addr)
	if err != nil {
		return nil, sderr.WithStack(err)
	}
	return db.Unsafe(), nil
}

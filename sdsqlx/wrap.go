package sdsqlx

import (
	"database/sql"

	"github.com/jmoiron/sqlx"
)

func NewMysqlDB(db *sql.DB) *sqlx.DB {
	return sqlx.NewDb(db, "mysql")
}

package sdgoqu

import (
	"github.com/doug-martin/goqu/v9"
	_ "github.com/doug-martin/goqu/v9/dialect/mysql"
	_ "github.com/doug-martin/goqu/v9/dialect/postgres"
	_ "github.com/doug-martin/goqu/v9/dialect/sqlite3"
)

var (
	Mysql    = goqu.Dialect("mysql")
	Postgres = goqu.Dialect("postgres")
	Sqlite3  = goqu.Dialect("sqlite3")
)

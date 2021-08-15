package sdsqlx

import (
	"context"
	"database/sql"
	"github.com/jmoiron/sqlx"
)

type NamedBinder interface {
	BindNamed(query string, arg interface{}) (string, []interface{}, error)
}

type Queryer interface {
	sqlx.Queryer
	Select(dest interface{}, query string, args ...interface{}) error
	Get(dest interface{}, query string, args ...interface{}) error
	NamedQuery(query string, arg interface{}) (*sqlx.Rows, error)
}

type Execer interface {
	sqlx.Execer
	MustExec(query string, args ...interface{}) sql.Result
	NamedExec(query string, arg interface{}) (sql.Result, error)
}

type Ext interface {
	Queryer
	Execer
	NamedBinder
}

type QueryerContext interface {
	sqlx.QueryerContext
	SelectContext(ctx context.Context, dest interface{}, query string, args ...interface{}) error
	GetContext(ctx context.Context, dest interface{}, query string, args ...interface{}) error
	NamedQueryContext(ctx context.Context, query string, arg interface{}) (*sqlx.Rows, error)
}

type ExecerContext interface {
	sqlx.ExecerContext
	MustExecContext(ctx context.Context, query string, args ...interface{}) sql.Result
	NamedExecContext(ctx context.Context, query string, arg interface{}) (sql.Result, error)
}

type ExtContext interface {
	QueryerContext
	ExecerContext
	NamedBinder
}

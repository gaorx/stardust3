package sdsqlx

import (
	"context"
	"database/sql"

	"github.com/gaorx/stardust3/sderr"
	"github.com/jmoiron/sqlx"
)

type (
	TxFunc        func(*sqlx.Tx) error
	TxContextFunc func(context.Context, *sqlx.Tx) error
)

type (
	TxFunc2        func(*sqlx.Tx) (interface{}, error)
	TxContextFunc2 func(context.Context, *sqlx.Tx) (interface{}, error)
)

func Tx(db *sqlx.DB, fn TxFunc) error {
	tx, err := db.Beginx()
	defer func() {
		if r := recover(); r != nil {
			_ = tx.Rollback()
		}
	}()
	if err != nil {
		return sderr.WithStack(err)
	}
	err = fn(tx)
	if err != nil {
		_ = tx.Rollback()
		return sderr.WithStack(err)
	}
	return sderr.WithStack(tx.Commit())
}

func TxContext(ctx context.Context, db *sqlx.DB, opts *sql.TxOptions, fn TxContextFunc) error {
	tx, err := db.BeginTxx(ctx, opts)
	defer func() {
		if r := recover(); r != nil {
			_ = tx.Rollback()
		}
	}()
	if err != nil {
		return sderr.WithStack(err)
	}
	err = fn(ctx, tx)
	if err != nil {
		_ = tx.Rollback()
		return sderr.WithStack(err)
	}
	return sderr.WithStack(tx.Commit())
}

func (dbg *DBGroup) WithTx(fn TxFunc) error {
	return Tx(dbg.W(), fn)
}

func (dbg *DBGroup) WithTxContext(ctx context.Context, opts *sql.TxOptions, fn TxContextFunc) error {
	return TxContext(ctx, dbg.W(), opts, fn)
}

func Tx2(db *sqlx.DB, fn TxFunc2) (interface{}, error) {
	var r interface{}
	sqlErr := Tx(db, func(tx *sqlx.Tx) error {
		r0, err := fn(tx)
		if err != nil {
			return sderr.WithStack(err)
		}
		r = r0
		return nil
	})
	if sqlErr != nil {
		return nil, sderr.WithStack(sqlErr)
	}
	return r, nil
}

func TxContext2(ctx context.Context, db *sqlx.DB, opts *sql.TxOptions, fn TxContextFunc2) (interface{}, error) {
	var r interface{}
	sqlErr := TxContext(ctx, db, opts, func(ctx context.Context, tx *sqlx.Tx) error {
		r0, err := fn(ctx, tx)
		if err != nil {
			return sderr.WithStack(err)
		}
		r = r0
		return nil
	})
	if sqlErr != nil {
		return nil, sderr.WithStack(sqlErr)
	}
	return r, nil
}

func MultiExec(tx *sqlx.Tx, queries []string) error {
	var errs []error
	for _, q := range queries {
		_, err := tx.Exec(q)
		if err != nil {
			errs = append(errs, err)
		}
	}
	return sderr.Multi(errs)
}

func (dbg *DBGroup) WithTx2(fn TxFunc2) (interface{}, error) {
	return Tx2(dbg.W(), fn)
}

func (dbg *DBGroup) WithTxContext2(ctx context.Context, opts *sql.TxOptions, fn TxContextFunc2) (interface{}, error) {
	return TxContext2(ctx, dbg.W(), opts, fn)
}

func IsNoRowsErr(err error) bool {
	return sderr.Is(err, sql.ErrNoRows)
}

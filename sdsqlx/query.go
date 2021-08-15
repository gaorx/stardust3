package sdsqlx

import (
	"context"

	"github.com/gaorx/stardust3/sderr"
	"github.com/jmoiron/sqlx"
)

type Query struct {
	C    context.Context
	S    string
	Args []interface{}
}

func Q(s string, args ...interface{}) Query {
	return Query{
		C:    nil,
		S:    s,
		Args: args,
	}
}

func (q Query) WithContext(ctx context.Context) Query {
	q.C = ctx
	return q
}

func (q Query) Do(db *sqlx.DB) (*sqlx.Rows, error) {
	var err error
	var rows *sqlx.Rows
	if q.C == nil {
		rows, err = db.Queryx(q.S, q.Args...)
	} else {
		rows, err = db.QueryxContext(q.C, q.S, q.Args)
	}
	if err != nil {
		return nil, sderr.WithStack(err)
	}
	return rows, nil
}

func (q Query) For(db *sqlx.DB, fn func(*sqlx.Rows) error) error {
	rows, err := q.Do(db)
	if err != nil {
		return sderr.WithStack(err)
	}
	defer func() {
		_ = rows.Close()
	}()

	var firstErr error
	for rows.Next() {
		err := fn(rows)
		if err != nil && firstErr == nil {
			firstErr = err
		}
	}
	return sderr.WithStack(firstErr)
}

func (q Query) TxDo(tx *sqlx.Tx) (*sqlx.Rows, error) {
	var err error
	var rows *sqlx.Rows
	if q.C == nil {
		rows, err = tx.Queryx(q.S, q.Args...)
	} else {
		rows, err = tx.QueryxContext(q.C, q.S, q.Args)
	}
	if err != nil {
		return nil, sderr.WithStack(err)
	}
	return rows, nil
}

func (q Query) TxFor(tx *sqlx.Tx, fn func(*sqlx.Rows) error) error {
	rows, err := q.TxDo(tx)
	if err != nil {
		return sderr.WithStack(err)
	}
	defer func() {
		_ = rows.Close()
	}()

	var firstErr error
	for rows.Next() {
		err := fn(rows)
		if err != nil && firstErr == nil {
			firstErr = err
		}
	}
	return sderr.WithStack(firstErr)
}

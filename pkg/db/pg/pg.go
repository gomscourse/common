package pg

import (
	"context"
	"fmt"
	"github.com/georgysavva/scany/pgxscan"
	"github.com/gomscourse/common/pkg/db"
	"github.com/gomscourse/common/pkg/db/prettier"
	"github.com/jackc/pgconn"
	"github.com/jackc/pgx/v4"
	"github.com/jackc/pgx/v4/pgxpool"
	"log"
	"time"
)

const TxKey = "tx"

type pg struct {
	dbc *pgxpool.Pool
}

func (p pg) BeginTx(ctx context.Context, opts pgx.TxOptions) (pgx.Tx, error) {
	return p.dbc.BeginTx(ctx, opts)
}

func (p pg) ScanOneContext(ctx context.Context, dest interface{}, q db.Query, args ...interface{}) error {
	logQuery(ctx, q, args...)

	row, err := p.QueryContext(ctx, q, args...)
	if err != nil {
		return err
	}

	return pgxscan.ScanOne(dest, row)
}

func (p pg) ScanAllContext(ctx context.Context, dest interface{}, q db.Query, args ...interface{}) error {
	logQuery(ctx, q, args...)

	rows, err := p.QueryContext(ctx, q, args...)
	if err != nil {
		return err
	}

	return pgxscan.ScanAll(dest, rows)
}

func (p pg) ExecContext(ctx context.Context, q db.Query, args ...interface{}) (pgconn.CommandTag, error) {
	logQuery(ctx, q, args...)

	ctx, cancel := context.WithTimeout(ctx, time.Second*5)
	defer cancel()

	tx, ok := ctx.Value(TxKey).(pgx.Tx)
	if ok {
		return tx.Exec(ctx, q.QueryRow, args...)
	}

	return p.dbc.Exec(ctx, q.QueryRow, args...)
}

func (p pg) QueryContext(ctx context.Context, q db.Query, args ...interface{}) (pgx.Rows, error) {
	logQuery(ctx, q, args...)

	ctx, cancel := context.WithTimeout(ctx, time.Second*5)
	defer cancel()

	tx, ok := ctx.Value(TxKey).(pgx.Tx)
	if ok {
		return tx.Query(ctx, q.QueryRow, args...)
	}

	return p.dbc.Query(ctx, q.QueryRow, args...)
}

func (p pg) QueryRowContext(ctx context.Context, q db.Query, args ...interface{}) pgx.Row {
	logQuery(ctx, q, args...)

	ctx, cancel := context.WithTimeout(ctx, time.Second*5)
	defer cancel()

	tx, ok := ctx.Value(TxKey).(pgx.Tx)
	if ok {
		return tx.QueryRow(ctx, q.QueryRow, args...)
	}

	return p.dbc.QueryRow(ctx, q.QueryRow, args...)
}

func (p pg) Ping(ctx context.Context) error {
	return p.dbc.Ping(ctx)
}

func (p pg) Close() {
	p.dbc.Close()
}

func logQuery(ctx context.Context, q db.Query, args ...interface{}) {
	prettyQuery := prettier.Pretty(q.QueryRow, prettier.PlaceholderDollar, args...)
	log.Println(
		ctx,
		fmt.Sprintf("sql: %s", q.Name),
		fmt.Sprintf("query: %s", prettyQuery),
	)
}

func MakeContextTx(ctx context.Context, tx pgx.Tx) context.Context {
	return context.WithValue(ctx, TxKey, tx)
}

package db

import (
	"context"
	"github.com/jackc/pgconn"
	"github.com/jackc/pgx/v4"
)

type Client interface {
	DB() DB
	Close() error
}

type Handler func(ctx context.Context) error

type TxManager interface {
	ReadCommitted(ctx context.Context, handler Handler) error
}

type Transactor interface {
	BeginTx(ctx context.Context, opts pgx.TxOptions) (pgx.Tx, error)
}

type Query struct {
	Name     string
	QueryRow string
}

type Pinger interface {
	Ping(ctx context.Context) error
}

type NamedExecutor interface {
	ScanOneContext(ctx context.Context, dest interface{}, q Query, args ...interface{}) error
	ScanAllContext(ctx context.Context, dest interface{}, q Query, args ...interface{}) error
}

type QueryExecutor interface {
	ExecContext(ctx context.Context, q Query, args ...interface{}) (pgconn.CommandTag, error)
	QueryContext(ctx context.Context, q Query, args ...interface{}) (pgx.Rows, error)
	QueryRowContext(ctx context.Context, q Query, args ...interface{}) pgx.Row
}

type SQLExecutor interface {
	NamedExecutor
	QueryExecutor
}

type DB interface {
	SQLExecutor
	Pinger
	Transactor
	Close()
}

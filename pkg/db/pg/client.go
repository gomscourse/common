package pg

import (
	"context"
	"github.com/gomscourse/common/pkg/db"
	"github.com/jackc/pgx/v4/pgxpool"
	"github.com/pkg/errors"
)

type pgClient struct {
	masterDBC db.DB
}

func New(ctx context.Context, dsn string) (db.Client, error) {
	pool, err := pgxpool.Connect(ctx, dsn)
	if err != nil {
		return nil, errors.Wrap(err, "failed to connect to postgres")
	}
	return &pgClient{masterDBC: &pg{dbc: pool}}, nil
}

func (c *pgClient) DB() db.DB {
	return c.masterDBC
}

func (c *pgClient) Close() error {
	if c.masterDBC != nil {
		c.masterDBC.Close()
	}

	return nil
}

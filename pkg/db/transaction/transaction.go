package transaction

import (
	"context"
	"github.com/gomscourse/common/pkg/db"
	"github.com/gomscourse/common/pkg/db/pg"
	"github.com/jackc/pgx/v4"
	"github.com/pkg/errors"
)

type manager struct {
	db db.Transactor
}

func NewTransactionManager(db db.Transactor) db.TxManager {
	return &manager{db: db}
}

func (m *manager) transaction(ctx context.Context, opts pgx.TxOptions, handler db.Handler) (err error) {
	tx, ok := ctx.Value(pg.TxKey).(pgx.Tx)
	if ok {
		return handler(ctx)
	}

	tx, err = m.db.BeginTx(ctx, opts)
	if err != nil {
		return errors.Wrap(err, "failed to begin tx")
	}

	ctx = pg.MakeContextTx(ctx, tx)

	defer func() {
		if r := recover(); r != nil {
			err = errors.Errorf("panic recovered: %v", r)
		}

		if err != nil {
			if errRollback := tx.Rollback(ctx); errRollback != nil {
				err = errors.Wrapf(err, "errRollback: %v", errRollback)
			}

			return
		}

		if err == nil {
			if errCommit := tx.Commit(ctx); errCommit != nil {
				err = errors.Wrapf(err, "tx commit failed")
			}
		}
	}()

	if err = handler(ctx); err != nil {
		err = errors.Wrap(err, "failed executing code inside transaction")
	}

	return err
}

func (m *manager) ReadCommitted(ctx context.Context, handler db.Handler) error {
	opts := pgx.TxOptions{IsoLevel: pgx.ReadCommitted}
	return m.transaction(ctx, opts, handler)
}

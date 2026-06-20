package otp

import (
	"context"
	"database/sql"
	"time"

	"github.com/uptrace/bun"
)

type defaultRepository struct {
	db *bun.DB
}

func (r *defaultRepository) Create(ctx context.Context, token *oneTimeTokenModel) error {
	_, err := r.db.NewInsert().Model(token).Exec(ctx)

	return err
}

func (r *defaultRepository) Delete(ctx context.Context, id string) error {
	_, err := r.db.NewDelete().Model((*oneTimeTokenModel)(nil)).Where("id = ?", id).Exec(ctx)
	return err
}

func (r *defaultRepository) DeleteAll(ctx context.Context, kind CodeKind, principal string) error {
	_, err := r.db.NewDelete().Model((*oneTimeTokenModel)(nil)).Where("kind = ?", kind).Where("principal = ?", principal).Exec(ctx)
	return err
}

// FindByKindAndPrincipal finds the latest code of the specified kind and principal, querying the db removes all codes for the specified kind and principal
func (r *defaultRepository) FindByKindAndPrincipal(ctx context.Context, kind CodeKind, principal string) (*oneTimeTokenModel, error) {
	var record oneTimeTokenModel

	err := r.db.RunInTx(ctx, &sql.TxOptions{}, func(ctx context.Context, tx bun.Tx) error {
		err := tx.NewSelect().Model(&record).
			Where("principal = ?", principal).
			Where("kind = ?", kind).
			Where("expires_at > ?", time.Now().UTC()).
			Order("created_at DESC").
			Limit(1).
			Scan(ctx)

		if err != nil {
			return err
		}

		_, err = tx.NewDelete().Model((*oneTimeTokenModel)(nil)).Where("kind = ?", kind).Where("principal = ?", principal).Exec(ctx)

		return err
	})

	return &record, err
}

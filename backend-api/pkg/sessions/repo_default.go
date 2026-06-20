package sessions

import (
	"context"

	"github.com/uptrace/bun"
)

type DefaultRepository struct {
	db *bun.DB
}

func (r *DefaultRepository) Create(ctx context.Context, s *Session) error {
	_, err := r.db.NewInsert().Model(s).Exec(ctx)
	return err
}

func (r *DefaultRepository) FindByToken(ctx context.Context, token string) (*Session, error) {
	record := &Session{}

	err := r.db.NewSelect().
		Model(record).
		Where("id = ?", token).
		Scan(ctx)

	if err != nil {
		return nil, err
	}

	return record, nil
}

func (r *DefaultRepository) Revoke(ctx context.Context, token string) {
	_, _ = r.db.NewDelete().
		Model((*Session)(nil)).
		Where("id = ?", token).
		Exec(ctx)
}

func (r *DefaultRepository) RevokeByPrincipalID(ctx context.Context, principalID string) {
	_, _ = r.db.NewDelete().
		Model((*Session)(nil)).
		Where("principal = ?", principalID).
		Exec(ctx)
}

func (r *DefaultRepository) CleanupExpired(ctx context.Context) (int, error) {
	res, err := r.db.NewDelete().
		Model((*Session)(nil)).
		Where("expires_at < NOW()").
		Exec(ctx)

	if err != nil {
		return 0, err
	}

	n, _ := res.RowsAffected()
	return int(n), nil
}

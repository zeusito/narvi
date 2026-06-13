package sessions

import (
	"context"
	"database/sql"
	"time"

	"github.com/uptrace/bun"
)

type DefaultRepository struct {
	db *bun.DB
}

func NewDefaultRepository(db *bun.DB) *DefaultRepository {
	return &DefaultRepository{db: db}
}

func (r *DefaultRepository) Create(ctx context.Context, s *Session) error {
	_, err := r.db.NewInsert().Model(s).Exec(ctx)
	return err
}

func (r *DefaultRepository) FindByToken(ctx context.Context, token string) (*Session, error) {
	s := new(Session)
	err := r.db.NewSelect().Model(s).Where("id = ?", token).Scan(ctx)
	if err != nil {
		if err == sql.ErrNoRows {
			return nil, nil
		}
		return nil, err
	}
	return s, nil
}

func (r *DefaultRepository) Revoke(ctx context.Context, token string) error {
	now := time.Now()
	_, err := r.db.NewUpdate().
		Model((*Session)(nil)).
		Set("revoked_at = ?", now).
		Where("id = ?", token).
		Where("revoked_at IS NULL").
		Exec(ctx)
	return err
}

func (r *DefaultRepository) RevokeByPrincipalID(ctx context.Context, principalID string) error {
	now := time.Now()
	_, err := r.db.NewUpdate().
		Model((*Session)(nil)).
		Set("revoked_at = ?", now).
		Where("principal_id = ?", principalID).
		Where("revoked_at IS NULL").
		Exec(ctx)
	return err
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

func (r *DefaultRepository) TouchLastUsed(ctx context.Context, token string) error {
	now := time.Now()
	_, err := r.db.NewUpdate().
		Model((*Session)(nil)).
		Set("last_used_at = ?", now).
		Where("id = ?", token).
		Exec(ctx)
	return err
}

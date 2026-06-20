package sessions

import (
	"context"

	"github.com/uptrace/bun"
)

type Repository interface {
	Create(ctx context.Context, s *Session) error
	FindByToken(ctx context.Context, token string) (*Session, error)
	Revoke(ctx context.Context, token string)
	RevokeByPrincipalID(ctx context.Context, principalID string)
	CleanupExpired(ctx context.Context) (int, error)
}

func NewDefaultRepository(db *bun.DB) *DefaultRepository {
	return &DefaultRepository{db: db}
}

package sessions

import "context"

type Repository interface {
	Create(ctx context.Context, s *Session) error
	FindByToken(ctx context.Context, token string) (*Session, error)
	Revoke(ctx context.Context, token string) error
	RevokeByPrincipalID(ctx context.Context, principalID string) error
	CleanupExpired(ctx context.Context) (int, error)
	TouchLastUsed(ctx context.Context, token string) error
}

package sessions

import (
	"context"
	"time"
)

type Service interface {
	CreateSession(ctx context.Context, principalID string, ttl time.Duration, opts ...SessionOption) (*Session, error)
	GetSession(ctx context.Context, token string) (*Session, error)
	RevokeSession(ctx context.Context, token string) error
	RevokePrincipalSessions(ctx context.Context, principalID string) error
	CleanupExpired(ctx context.Context) (int, error)
}

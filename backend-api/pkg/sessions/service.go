package sessions

import (
	"context"
	"time"

	"github.com/zeusito/narvi/pkg/toolbox/hasher"
)

type Manager interface {
	Create(ctx context.Context, claims PrincipalClaims, remoteAddr string, userAgent string, ttl time.Duration) (string, error)
	GetAndVerify(ctx context.Context, token string) (*Session, error)
	Revoke(ctx context.Context, token string) error
	RevokeAll(ctx context.Context, principalID string) error
	Cleanup(ctx context.Context) (int, error)
}

func NewDefaultManager(repo Repository, tokenHasher hasher.Hasher) *DefaultManager {
	return &DefaultManager{repo: repo, tokenHasher: tokenHasher}
}

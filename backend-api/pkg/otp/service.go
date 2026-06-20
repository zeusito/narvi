package otp

import (
	"context"
	"time"

	"github.com/zeusito/narvi/pkg/toolbox/hasher"
)

type Manager interface {
	GenerateCode(ctx context.Context, length int, kind CodeKind, principal string, tenant string) (string, bool)
	VerifyCode(ctx context.Context, kind CodeKind, principal, suppliedCode string) bool
	Remove(ctx context.Context, kind CodeKind, principal string) bool
}

func NewDefaultManager(repo repository, codeHasher hasher.Hasher) Manager {
	return &defaultManager{
		repo:               repo,
		codeHasher:         codeHasher,
		expirationDuration: 5 * time.Minute,
	}
}

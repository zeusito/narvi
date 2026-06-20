package otp

import (
	"context"

	"github.com/uptrace/bun"
)

type repository interface {
	Create(ctx context.Context, token *oneTimeTokenModel) error
	Delete(ctx context.Context, id string) error
	DeleteAll(ctx context.Context, kind CodeKind, principal string) error
	FindByKindAndPrincipal(ctx context.Context, kind CodeKind, principal string) (*oneTimeTokenModel, error)
}

func newDefaultRepository(db *bun.DB) repository {
	return &defaultRepository{db: db}
}

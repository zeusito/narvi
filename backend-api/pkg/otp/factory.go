package otp

import (
	"github.com/uptrace/bun"
	"github.com/zeusito/narvi/pkg/toolbox/hasher"
)

type Module struct {
	Manager Manager
}

func NewModule(db *bun.DB, hasher hasher.Hasher) *Module {
	repo := newDefaultRepository(db)
	svc := NewDefaultManager(repo, hasher)

	return &Module{
		Manager: svc,
	}
}

package sessions

import "github.com/uptrace/bun"

type Module struct {
	Service Service
}

func NewModule(db *bun.DB) *Module {
	repo := NewDefaultRepository(db)
	svc := NewDefaultService(repo)
	return &Module{Service: svc}
}

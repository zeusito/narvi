package healthcheck

import (
	"github.com/go-chi/chi/v5"
)

type Module struct {
}

func NewModule(mux *chi.Mux) *Module {
	_ = newHealthController(mux)
	return &Module{}
}

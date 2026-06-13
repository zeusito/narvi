package healthcheck

import (
	"net/http"

	"github.com/zeusito/narvi/pkg/router"

	"github.com/go-chi/chi/v5"
)

type healthController struct{}

func newHealthController(mux *chi.Mux) *healthController {
	c := &healthController{}

	mux.Get("/health/readiness", c.handleReadiness)
	mux.Get("/health/liveness", c.handleLiveness)

	return c
}

func (c *healthController) handleReadiness(w http.ResponseWriter, req *http.Request) {
	router.RenderJSON(req.Context(), w, http.StatusOK, router.DefaultSuccessResponseBody())
}

func (c *healthController) handleLiveness(w http.ResponseWriter, req *http.Request) {
	router.RenderJSON(req.Context(), w, http.StatusOK, router.DefaultSuccessResponseBody())
}

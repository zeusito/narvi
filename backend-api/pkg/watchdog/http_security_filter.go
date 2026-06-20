package watchdog

import (
	"net/http"
	"strings"

	"github.com/rs/zerolog/log"
	"github.com/zeusito/narvi/pkg/sessions"
)

// AuthenticationFilter is a middleware that checks if the request has a valid token
func AuthenticationFilter(sessionManager sessions.Manager) func(http.Handler) http.Handler {
	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			token := r.Header.Get("Authorization")

			if token == "" || !strings.HasPrefix(token, "Bearer ") {
				log.Warn().Msg("no token provided")
				http.Error(w, "Unauthorized", http.StatusUnauthorized)
				return
			}

			// Remove the "Bearer " prefix
			token = token[7:]

			// Validate the token
			record, err := sessionManager.GetAndVerify(r.Context(), token)
			if err != nil {
				log.Warn().Msg("session not authenticated")
				http.Error(w, "Unauthorized", http.StatusUnauthorized)
				return
			}

			// Add claims to context
			ctx := sessions.AddClaimsToContext(r.Context(), record.Claims)

			next.ServeHTTP(w, r.WithContext(ctx))
		})
	}
}

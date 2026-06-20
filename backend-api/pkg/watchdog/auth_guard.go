package watchdog

import (
	"context"
	"slices"

	"github.com/zeusito/narvi/pkg/sessions"
)

func PreAuthorize(ctx context.Context, requiredAudience sessions.Audience, requiredRole string) (*sessions.PrincipalClaims, bool) {
	claims := sessions.ExtractClaimsFromContext(ctx)

	// If the principal is not authenticated, reject the action
	if !claims.IsAuthenticated {
		return nil, false
	}

	// Check audience
	if claims.Audience != requiredAudience {
		return nil, false
	}

	// Check if the required role is present in the claims
	if !slices.Contains(claims.Roles, requiredRole) {
		return nil, false
	}

	// If all checks pass, return the claims and allow the action
	return &claims, true
}

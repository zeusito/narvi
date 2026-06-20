package sessions

import (
	"context"
)

type ctxKeyAuthClaims int

type Audience string

const (
	PrincipalClaimsKey ctxKeyAuthClaims = 1

	AUD_NARVI Audience = "narvi"
)

// PrincipalClaims represents the claims of a principal, customize it as needed
type PrincipalClaims struct {
	IsAuthenticated bool     `json:"isAuthenticated"`
	Audience        Audience `json:"audience"`
	Subject         string   `json:"subject"`
	SubjectName     string   `json:"subjectName"`
	TenantID        string   `json:"tenantId"`
	TenantSlug      string   `json:"tenantSlug"`
	TenantName      string   `json:"tenantName"`
	Roles           []string `json:"roles"`
}

func AddClaimsToContext(ctx context.Context, claims PrincipalClaims) context.Context {
	return context.WithValue(ctx, PrincipalClaimsKey, claims)
}

func ExtractClaimsFromContext(ctx context.Context) PrincipalClaims {
	claims, ok := ctx.Value(PrincipalClaimsKey).(PrincipalClaims)
	if !ok {
		return PrincipalClaims{IsAuthenticated: false}
	}

	return claims
}

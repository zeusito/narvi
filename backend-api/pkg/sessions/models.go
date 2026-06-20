package sessions

import (
	"time"

	"github.com/uptrace/bun"
)

type Session struct {
	bun.BaseModel `bun:"table:principal_sessions,alias:ps"`
	ID            string          `bun:"id,pk"` // hashed ID
	Principal     string          `bun:"principal"`
	IPAddress     string          `bun:"ip_address,type:inet"`
	UserAgent     string          `bun:"user_agent"`
	Tenant        string          `bun:"tenant"`
	Claims        PrincipalClaims `bun:"metadata"`
	ExpiresAt     time.Time       `bun:"expires_at"`
	CreatedAt     time.Time       `bun:"created_at"`
}

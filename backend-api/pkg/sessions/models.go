package sessions

import (
	"time"

	"github.com/uptrace/bun"
)

type Session struct {
	bun.BaseModel `bun:"table:sessions"`

	ID          string     `bun:"id,pk" json:"id"`
	PrincipalID string     `bun:"principal_id,notnull" json:"principal_id"`
	IPAddress   string     `bun:"ip_address,notnull,default:''" json:"ip_address"`
	UserAgent   string     `bun:"user_agent,notnull,default:''" json:"user_agent"`
	ExpiresAt   time.Time  `bun:"expires_at,notnull" json:"expires_at"`
	CreatedAt   time.Time  `bun:"created_at,notnull,default:current_timestamp" json:"created_at"`
	LastUsedAt  *time.Time `bun:"last_used_at" json:"last_used_at,omitempty"`
	RevokedAt   *time.Time `bun:"revoked_at" json:"revoked_at,omitempty"`
	Metadata    []byte     `bun:"metadata,type:jsonb,default:'{}'" json:"metadata,omitempty"`
}

type SessionOption func(s *Session)

func WithIP(ip string) SessionOption {
	return func(s *Session) { s.IPAddress = ip }
}

func WithUserAgent(ua string) SessionOption {
	return func(s *Session) { s.UserAgent = ua }
}

func WithMetadata(md []byte) SessionOption {
	return func(s *Session) { s.Metadata = md }
}

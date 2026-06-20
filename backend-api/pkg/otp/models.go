package otp

import (
	"time"

	"github.com/uptrace/bun"
)

type CodeKind string

const (
	CodeKindUserPassword CodeKind = "user_password"
)

type oneTimeTokenModel struct {
	bun.BaseModel `bun:"table:principal_tokens,alias:pt"`
	ID            string    `bun:"id,pk"`
	Kind          CodeKind  `bun:"kind"`
	Principal     string    `bun:"principal"`
	Tenant        string    `bun:"tenant"`
	ExpiresAt     time.Time `bun:"expires_at"`
	CreatedAt     time.Time `bun:"created_at"`
}

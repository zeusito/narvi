package otp

import (
	"context"
	"time"

	"github.com/rs/zerolog/log"
	"github.com/zeusito/narvi/pkg/toolbox"
	"github.com/zeusito/narvi/pkg/toolbox/hasher"
)

type defaultManager struct {
	repo               repository
	codeHasher         hasher.Hasher
	expirationDuration time.Duration
}

// GenerateCode generates a random code of the specified length and kind
func (m *defaultManager) GenerateCode(ctx context.Context, length int, kind CodeKind, principal string, tenant string) (string, bool) {
	code := toolbox.SecureRandomString(length)
	now := time.Now().UTC()

	// hash the code
	hashedCode, err := m.codeHasher.Hash(code)
	if err != nil {
		log.Error().Err(err).Msg("failed to hash code")
		return "", false
	}

	// Persist the code
	err = m.repo.Create(ctx, &oneTimeTokenModel{
		ID:        hashedCode,
		Kind:      kind,
		Principal: principal,
		Tenant:    tenant,
		ExpiresAt: now.Add(m.expirationDuration),
		CreatedAt: now,
	})

	if err != nil {
		log.Error().Err(err).Msg("failed to insert OTP")
		return "", false
	}

	return code, true
}

// VerifyCode verifies the code of the specified kind and principal.
// By default, only the latest code from the combined kind and principal is valid.
// Expiration is checked at the storage level.
func (m *defaultManager) VerifyCode(ctx context.Context, kind CodeKind, principal, suppliedCode string) bool {
	record, err := m.repo.FindByKindAndPrincipal(ctx, kind, principal)
	if err != nil {
		return false
	}

	// check if the hashes match
	if !m.codeHasher.Verify(suppliedCode, record.ID) {
		log.Error().Msg("hashes do not match")
		return false
	}

	return true
}

// Remove removes the code from the storage. All codes for the specified kind and principal are removed.
func (m *defaultManager) Remove(ctx context.Context, kind CodeKind, principal string) bool {
	err := m.repo.DeleteAll(ctx, kind, principal)
	return err == nil
}

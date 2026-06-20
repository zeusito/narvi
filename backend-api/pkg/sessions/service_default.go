package sessions

import (
	"context"
	"time"

	"github.com/rs/zerolog/log"
	"github.com/zeusito/narvi/pkg/terrors"
	"github.com/zeusito/narvi/pkg/toolbox/hasher"
)

type DefaultManager struct {
	repo        Repository
	tokenHasher hasher.Hasher
}

// Create starts a new user session associated with a random opaque token.
// The token is returned to the caller and can be used to authenticate requests.
// The session will be automatically cleaned up after the specified TTL.
func (s *DefaultManager) Create(ctx context.Context, claims PrincipalClaims, remoteAddr string, userAgent string, ttl time.Duration) (string, error) {
	if ttl <= 0 {
		return "", terrors.PreconditionFailed("ttl must be positive")
	}

	if !claims.IsAuthenticated {
		return "", terrors.PreconditionFailed("principal must be authenticated")
	}

	now := time.Now().UTC()

	token, hashedToken, err := newOpaqueToken(s.tokenHasher, "session")
	if err != nil {
		log.Error().Err(err).Str("principal", claims.Subject).Msg("failed to generate token")
		return "", terrors.OperationFailed("failed to generate token")
	}

	record := &Session{
		ID:        hashedToken,
		Principal: claims.Subject,
		IPAddress: remoteAddr,
		UserAgent: userAgent,
		Tenant:    claims.TenantID,
		Claims:    claims,
		ExpiresAt: now.Add(ttl),
		CreatedAt: now,
	}

	if err := s.repo.Create(ctx, record); err != nil {
		log.Error().Err(err).Str("principal", claims.Subject).Msg("failed to create session")
		return "", terrors.OperationFailed("failed to create session")
	}

	return token, nil
}

// GetAndVerify retrieves and verifies a session by its token.
func (s *DefaultManager) GetAndVerify(ctx context.Context, token string) (*Session, error) {
	if token == "" {
		return nil, terrors.PreconditionFailed("token is required")
	}

	now := time.Now().UTC()

	// Hash the token before looking it up
	hashedToken, err := s.tokenHasher.Hash(token)
	if err != nil {
		log.Error().Err(err).Msg("failed to hash token")
		return nil, terrors.OperationFailed("failed to retrieve session")
	}

	session, err := s.repo.FindByToken(ctx, hashedToken)
	if err != nil {
		log.Error().Err(err).Msg("failed to find session")
		return nil, terrors.OperationFailed("failed to retrieve session")
	}

	if now.After(session.ExpiresAt) {
		return nil, terrors.RecordNotFound("session has expired")
	}

	return session, nil
}

// Revoke revokes a session immediately.
func (s *DefaultManager) Revoke(ctx context.Context, token string) error {
	if token == "" {
		return terrors.PreconditionFailed("token is required")
	}

	// Hash the token before looking it up
	hashedToken, err := s.tokenHasher.Hash(token)
	if err != nil {
		log.Error().Err(err).Msg("failed to hash token")
		return terrors.OperationFailed("failed to revoke session")
	}

	s.repo.Revoke(ctx, hashedToken)

	return nil
}

// RevokeAll revokes all sessions for a principal.
func (s *DefaultManager) RevokeAll(ctx context.Context, principalID string) error {
	if principalID == "" {
		return terrors.PreconditionFailed("principal_id is required")
	}

	s.repo.RevokeByPrincipalID(ctx, principalID)

	return nil
}

// Cleanup revokes all expired sessions.
func (s *DefaultManager) Cleanup(ctx context.Context) (int, error) {
	n, err := s.repo.CleanupExpired(ctx)
	if err != nil {
		log.Error().Err(err).Msg("failed to cleanup expired sessions")
		return 0, terrors.OperationFailed("failed to cleanup expired sessions")
	}

	return n, nil
}

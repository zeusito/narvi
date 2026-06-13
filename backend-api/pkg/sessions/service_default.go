package sessions

import (
	"context"
	"crypto/rand"
	"encoding/hex"
	"time"

	"github.com/rs/zerolog/log"
	"github.com/zeusito/narvi/pkg/terrors"
)

const tokenByteSize = 32

type DefaultService struct {
	repo Repository
}

func NewDefaultService(repo Repository) *DefaultService {
	return &DefaultService{repo: repo}
}

func (s *DefaultService) CreateSession(ctx context.Context, principalID string, ttl time.Duration, opts ...SessionOption) (*Session, error) {
	if principalID == "" {
		return nil, terrors.PreconditionFailed("principal_id is required")
	}
	if ttl <= 0 {
		return nil, terrors.PreconditionFailed("ttl must be positive")
	}

	token, err := generateToken()
	if err != nil {
		return nil, terrors.OperationFailed("failed to generate session token")
	}

	session := &Session{
		ID:          token,
		PrincipalID: principalID,
		ExpiresAt:   time.Now().UTC().Add(ttl),
		CreatedAt:   time.Now().UTC(),
	}

	for _, opt := range opts {
		opt(session)
	}

	if session.Metadata == nil {
		session.Metadata = []byte("{}")
	}

	if err := s.repo.Create(ctx, session); err != nil {
		log.Error().Err(err).Str("principal_id", principalID).Msg("failed to create session")
		return nil, terrors.OperationFailed("failed to create session")
	}

	return session, nil
}

func (s *DefaultService) GetSession(ctx context.Context, token string) (*Session, error) {
	if token == "" {
		return nil, terrors.PreconditionFailed("token is required")
	}

	session, err := s.repo.FindByToken(ctx, token)
	if err != nil {
		log.Error().Err(err).Str("token", token).Msg("failed to find session")
		return nil, terrors.OperationFailed("failed to retrieve session")
	}
	if session == nil {
		return nil, terrors.RecordNotFound("session not found")
	}
	if session.RevokedAt != nil {
		return nil, terrors.RecordNotFound("session has been revoked")
	}
	if time.Now().UTC().After(session.ExpiresAt) {
		return nil, terrors.RecordNotFound("session has expired")
	}

	if touchErr := s.repo.TouchLastUsed(ctx, token); touchErr != nil {
		log.Warn().Err(touchErr).Str("token", token).Msg("failed to update last_used_at")
	}

	return session, nil
}

func (s *DefaultService) RevokeSession(ctx context.Context, token string) error {
	if token == "" {
		return terrors.PreconditionFailed("token is required")
	}

	if err := s.repo.Revoke(ctx, token); err != nil {
		log.Error().Err(err).Str("token", token).Msg("failed to revoke session")
		return terrors.OperationFailed("failed to revoke session")
	}
	return nil
}

func (s *DefaultService) RevokePrincipalSessions(ctx context.Context, principalID string) error {
	if principalID == "" {
		return terrors.PreconditionFailed("principal_id is required")
	}

	if err := s.repo.RevokeByPrincipalID(ctx, principalID); err != nil {
		log.Error().Err(err).Str("principal_id", principalID).Msg("failed to revoke principal sessions")
		return terrors.OperationFailed("failed to revoke principal sessions")
	}
	return nil
}

func (s *DefaultService) CleanupExpired(ctx context.Context) (int, error) {
	n, err := s.repo.CleanupExpired(ctx)
	if err != nil {
		log.Error().Err(err).Msg("failed to cleanup expired sessions")
		return 0, terrors.OperationFailed("failed to cleanup expired sessions")
	}
	return n, nil
}

func generateToken() (string, error) {
	b := make([]byte, tokenByteSize)
	if _, err := rand.Read(b); err != nil {
		return "", err
	}
	return hex.EncodeToString(b), nil
}

package sessions

import (
	"context"
	"testing"
	"time"

	"github.com/rs/zerolog/log"
	"github.com/stretchr/testify/assert"
	"github.com/zeusito/narvi/pkg/configurer"
	"github.com/zeusito/narvi/pkg/db"
	"github.com/zeusito/narvi/pkg/toolbox/hasher"
	"github.com/zeusito/narvi/pkg/toolbox/testbox"
)

var testDB *db.DatabaseConnection

func TestMain(m *testing.M) {
	initScriptPath := []string{
		"../../db/schema.sql",
	}

	connData, closeFunc, err := testbox.InitPostgresqlContainer(context.Background(), initScriptPath)
	if err != nil {
		log.Fatal().Err(err).Msg("unable to create container")
	}
	defer closeFunc()

	testDB = db.MustCreatePooledConnection(configurer.DatabaseConfigurations{
		Enabled:  true,
		DbName:   connData.DBName,
		Username: connData.UserName,
		Password: connData.Password,
		Host:     connData.Host,
		Port:     connData.Port,
		PoolSize: 1,
	})

	// Run the tests
	m.Run()
}

func TestCreatePrincipalSessionWithUnauthenticatedPrincipal(t *testing.T) {
	ctx := t.Context()
	theHasher, err := hasher.NewHmacSHA256("dGVzdC1zZWNyZXQ=")
	repo := NewDefaultRepository(testDB.Conn)
	manager := &DefaultManager{
		repo:        repo,
		tokenHasher: theHasher,
	}

	// Make sure the hasher was created successfully
	assert.NoError(t, err)

	// Try to create a session with an unauthenticated principal (should fail)
	id, err := manager.Create(ctx, PrincipalClaims{IsAuthenticated: false}, "192.168.1.1", "user-agent", time.Hour)

	assert.Error(t, err)
	assert.Empty(t, id)
}

func TestCreatePrincipalSessionSuccess(t *testing.T) {
	ctx := t.Context()
	theHasher, err := hasher.NewHmacSHA256("dGVzdC1zZWNyZXQ=")
	repo := NewDefaultRepository(testDB.Conn)
	manager := &DefaultManager{
		repo:        repo,
		tokenHasher: theHasher,
	}

	// Make sure the hasher was created successfully
	assert.NoError(t, err)

	// Try to create a session with a valid principal (should succeed)
	token, err := manager.Create(ctx, PrincipalClaims{
		IsAuthenticated: true,
		Subject:         "user-123",
		Roles:           []string{"MY_ROLE"},
	}, "192.168.1.1", "user-agent", time.Hour)

	assert.NoError(t, err)
	assert.NotEmpty(t, token)
}

func TestVerifyPrincipalSessionWithValidToken(t *testing.T) {
	ctx := t.Context()
	theHasher, err := hasher.NewHmacSHA256("dGVzdC1zZWNyZXQ=")
	repo := NewDefaultRepository(testDB.Conn)
	manager := &DefaultManager{
		repo:        repo,
		tokenHasher: theHasher,
	}

	// Make sure the hasher was created successfully
	assert.NoError(t, err)

	// Try to create a session with a valid principal (should succeed)
	token, err := manager.Create(ctx, PrincipalClaims{
		IsAuthenticated: true,
		Subject:         "user-123",
		Roles:           []string{"MY_ROLE"},
	}, "192.168.1.1", "user-agent", time.Hour)

	assert.NoError(t, err)
	assert.NotEmpty(t, token)

	// Try to verify the session with the valid token (should succeed)
	verifiedSession, err := manager.GetAndVerify(ctx, token)

	assert.NoError(t, err)
	assert.Equal(t, "user-123", verifiedSession.Principal)
	assert.Equal(t, "MY_ROLE", verifiedSession.Claims.Roles[0])
	assert.Equal(t, "192.168.1.1", verifiedSession.IPAddress)
	assert.Equal(t, "user-agent", verifiedSession.UserAgent)
}

func TestVerifyPrincipalSessionWithExpiredToken(t *testing.T) {
	ctx := t.Context()
	theHasher, err := hasher.NewHmacSHA256("dGVzdC1zZWNyZXQ=")
	repo := NewDefaultRepository(testDB.Conn)
	manager := &DefaultManager{
		repo:        repo,
		tokenHasher: theHasher,
	}

	// Make sure the hasher was created successfully
	assert.NoError(t, err)

	// Try to create a session with a valid principal (should succeed)
	token, err := manager.Create(ctx, PrincipalClaims{
		IsAuthenticated: true,
		Subject:         "user-123",
		Roles:           []string{"MY_ROLE"},
	}, "192.168.1.1", "user-agent", time.Second)

	assert.NoError(t, err)
	assert.NotEmpty(t, token)

	// Wait for the session to expire
	time.Sleep(2 * time.Second)

	// Try to verify the session with the expired token (should fail)
	verifiedSession, err := manager.GetAndVerify(ctx, token)

	assert.Error(t, err)
	assert.Empty(t, verifiedSession)
}

func TestRevokePrincipalSession(t *testing.T) {
	ctx := t.Context()
	theHasher, err := hasher.NewHmacSHA256("dGVzdC1zZWNyZXQ=")
	repo := NewDefaultRepository(testDB.Conn)
	manager := &DefaultManager{
		repo:        repo,
		tokenHasher: theHasher,
	}

	// Make sure the hasher was created successfully
	assert.NoError(t, err)

	// Try to create a session with a valid principal (should succeed)
	token, err := manager.Create(ctx, PrincipalClaims{
		IsAuthenticated: true,
		Subject:         "user-456",
		Roles:           []string{"MY_ROLE"},
	}, "192.168.1.1", "user-agent", time.Hour)

	assert.NoError(t, err)
	assert.NotEmpty(t, token)

	// Try to revoke the session with the valid token (should succeed)
	err = manager.Revoke(ctx, token)
	assert.NoError(t, err)

	// Try to verify the session with the revoked token (should fail)
	verifiedSession, err := manager.GetAndVerify(ctx, token)

	assert.Error(t, err)
	assert.Empty(t, verifiedSession)
}

func TestRevokeAllPrincipalSessions(t *testing.T) {
	ctx := t.Context()
	theHasher, err := hasher.NewHmacSHA256("dGVzdC1zZWNyZXQ=")
	repo := NewDefaultRepository(testDB.Conn)
	manager := &DefaultManager{
		repo:        repo,
		tokenHasher: theHasher,
	}

	// Make sure the hasher was created successfully
	assert.NoError(t, err)

	// Try to create a session with a valid principal (should succeed)
	token1, err := manager.Create(ctx, PrincipalClaims{
		IsAuthenticated: true,
		Subject:         "user-bad",
		Roles:           []string{"MY_ROLE"},
	}, "192.168.1.1", "user-agent", time.Hour)

	assert.NoError(t, err)
	assert.NotEmpty(t, token1)

	// Try to create another session with a valid principal (should succeed)
	token2, err := manager.Create(ctx, PrincipalClaims{
		IsAuthenticated: true,
		Subject:         "user-bad",
		Roles:           []string{"MY_ROLE"},
	}, "192.168.1.1", "user-agent", time.Hour)

	assert.NoError(t, err)
	assert.NotEmpty(t, token2)

	// Try to revoke all sessions for the principal (should succeed)
	err = manager.RevokeAll(ctx, "user-bad")
	assert.NoError(t, err)

	// Try to verify the sessions with the revoked tokens (should fail)
	verifiedSession1, err := manager.GetAndVerify(ctx, token1)
	assert.Error(t, err)
	assert.Empty(t, verifiedSession1)

	verifiedSession2, err := manager.GetAndVerify(ctx, token2)
	assert.Error(t, err)
	assert.Empty(t, verifiedSession2)
}

func TestCleanupExpiredSessions(t *testing.T) {
	ctx := t.Context()
	theHasher, err := hasher.NewHmacSHA256("dGVzdC1zZWNyZXQ=")
	repo := NewDefaultRepository(testDB.Conn)
	manager := &DefaultManager{
		repo:        repo,
		tokenHasher: theHasher,
	}

	// Make sure the hasher was created successfully
	assert.NoError(t, err)

	// Try to create a session with a valid principal (should succeed)
	token1, err := manager.Create(ctx, PrincipalClaims{
		IsAuthenticated: true,
		Subject:         "user-exp",
		Roles:           []string{"MY_ROLE"},
	}, "192.168.1.1", "user-agent", time.Second)

	assert.NoError(t, err)
	assert.NotEmpty(t, token1)

	// Try to create another session with a valid principal (should succeed)
	token2, err := manager.Create(ctx, PrincipalClaims{
		IsAuthenticated: true,
		Subject:         "user-exp",
		Roles:           []string{"MY_ROLE"},
	}, "192.168.1.1", "user-agent", time.Second)

	assert.NoError(t, err)
	assert.NotEmpty(t, token2)

	// Wait for the sessions to expire
	time.Sleep(2 * time.Second)

	// Try to cleanup expired sessions (should succeed)
	count, err := manager.Cleanup(ctx)
	assert.NoError(t, err)
	assert.GreaterOrEqual(t, count, 2)
}

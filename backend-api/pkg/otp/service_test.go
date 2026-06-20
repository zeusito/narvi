package otp

import (
	"context"
	"testing"

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

func TestGenerateCodeSuccess(t *testing.T) {
	ctx := t.Context()
	theHasher, err := hasher.NewHmacSHA256("dGVzdC1zZWNyZXQ=")
	assert.NoError(t, err)

	repo := newDefaultRepository(testDB.Conn)
	manager := NewDefaultManager(repo, theHasher)

	code, ok := manager.GenerateCode(ctx, 6, CodeKindUserPassword, "user-123", "tenant-123")

	assert.True(t, ok)
	assert.NotEmpty(t, code)
	assert.Len(t, code, 6)
}

func TestVerifyCodeSuccess(t *testing.T) {
	ctx := t.Context()
	theHasher, err := hasher.NewHmacSHA256("dGVzdC1zZWNyZXQ=")
	assert.NoError(t, err)

	repo := newDefaultRepository(testDB.Conn)
	manager := NewDefaultManager(repo, theHasher)

	code, ok := manager.GenerateCode(ctx, 6, CodeKindUserPassword, "user-456", "tenant-123")
	assert.True(t, ok)
	assert.NotEmpty(t, code)

	verified := manager.VerifyCode(ctx, CodeKindUserPassword, "user-456", code)
	assert.True(t, verified)
}

func TestVerifyCodeInvalidCode(t *testing.T) {
	ctx := t.Context()
	theHasher, err := hasher.NewHmacSHA256("dGVzdC1zZWNyZXQ=")
	assert.NoError(t, err)

	repo := newDefaultRepository(testDB.Conn)
	manager := NewDefaultManager(repo, theHasher)

	code, ok := manager.GenerateCode(ctx, 6, CodeKindUserPassword, "user-789", "tenant-123")
	assert.True(t, ok)
	assert.NotEmpty(t, code)

	verified := manager.VerifyCode(ctx, CodeKindUserPassword, "user-789", "invalid-code")
	assert.False(t, verified)
}

func TestVerifyOnlyLatestCodeIsValid(t *testing.T) {
	ctx := t.Context()
	theHasher, err := hasher.NewHmacSHA256("dGVzdC1zZWNyZXQ=")
	assert.NoError(t, err)

	repo := newDefaultRepository(testDB.Conn)
	manager := NewDefaultManager(repo, theHasher)

	code, ok := manager.GenerateCode(ctx, 6, CodeKindUserPassword, "user-1011", "tenant-123")
	assert.True(t, ok)
	assert.NotEmpty(t, code)

	// Generate another code
	code2, ok := manager.GenerateCode(ctx, 6, CodeKindUserPassword, "user-1011", "tenant-123")
	assert.True(t, ok)
	assert.NotEmpty(t, code2)
	assert.NotEqual(t, code, code2)

	// Verify the second code first, during verification all codes are invalidated
	verified := manager.VerifyCode(ctx, CodeKindUserPassword, "user-1011", code2)
	assert.True(t, verified)

	// Verify the first code again, it should not be valid since it was invalidated during the previous verification
	verified = manager.VerifyCode(ctx, CodeKindUserPassword, "user-1011", code)
	assert.False(t, verified)
}

func TestVerifyCodeInvalidPrincipal(t *testing.T) {
	ctx := t.Context()
	theHasher, err := hasher.NewHmacSHA256("dGVzdC1zZWNyZXQ=")
	assert.NoError(t, err)

	repo := newDefaultRepository(testDB.Conn)
	manager := NewDefaultManager(repo, theHasher)

	code, ok := manager.GenerateCode(ctx, 6, CodeKindUserPassword, "userx-123", "tenant-123")
	assert.True(t, ok)
	assert.NotEmpty(t, code)

	verified := manager.VerifyCode(ctx, CodeKindUserPassword, "userx-999", code)
	assert.False(t, verified)
}

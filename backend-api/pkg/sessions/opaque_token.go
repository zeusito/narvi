package sessions

import (
	"crypto/rand"
	"encoding/base64"

	"github.com/zeusito/narvi/pkg/toolbox/hasher"
)

// newOpaqueToken generates a new opaque token following the spec:
// - Generate 256 bits of random data
// - Apply base64 URL-safe encoding
// - Add a descriptive prefix
// - Hash the full prefixed string for storage
func newOpaqueToken(th hasher.Hasher, prefix string) (token string, hashedToken string, err error) {
	// 1. Generate 32 bytes of raw entropy (256 bits)
	rawBytes := make([]byte, 32)
	if _, err := rand.Read(rawBytes); err != nil {
		return "", "", err
	}

	// 2. Base64URL encoding without padding
	encoded := base64.RawURLEncoding.EncodeToString(rawBytes)

	// 3. Add a descriptive prefix (if any)
	token = encoded
	if prefix != "" {
		token = prefix + "_" + encoded
	}

	// 4. Hash the full prefixed string for storage
	hashedToken, err = th.Hash(token)

	return token, hashedToken, err
}

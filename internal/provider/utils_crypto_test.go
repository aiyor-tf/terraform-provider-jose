// Copyright (c) Tze Liang
// SPDX-License-Identifier: MPL-2.0

package provider

import (
	"crypto/ecdsa"
	"crypto/ed25519"
	"crypto/rsa"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestGenerateKey(t *testing.T) {
	tests := []struct {
		name    string
		alg     string
		size    int
		wantErr bool
	}{
		{"RSA 2048", "RS256", 2048, false},
		{"RSA Default", "RS256", 0, false},
		{"ECDSA P-256", "ES256", 0, false},
		{"ECDSA P-384", "ES384", 0, false},
		{"ECDSA P-521", "ES512", 0, false},
		{"Ed25519", "EdDSA", 0, false},
		{"Invalid Alg", "INVALID", 0, true},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			priv, pub, err := generateKey(tt.alg, tt.size)
			if tt.wantErr {
				assert.Error(t, err)
				return
			}
			assert.NoError(t, err)
			assert.NotNil(t, priv)
			assert.NotNil(t, pub)

			switch tt.alg {
			case "RS256":
				assert.IsType(t, &rsa.PrivateKey{}, priv)
				assert.IsType(t, &rsa.PublicKey{}, pub)
			case "ES256", "ES384", "ES512":
				assert.IsType(t, &ecdsa.PrivateKey{}, priv)
				assert.IsType(t, &ecdsa.PublicKey{}, pub)
			case "EdDSA":
				assert.IsType(t, ed25519.PrivateKey{}, priv)
				assert.IsType(t, ed25519.PublicKey{}, pub)
			}
		})
	}
}

func TestBuildJWK(t *testing.T) {
	// Generate a key to use
	_, pub, _ := generateKey("RS256", 2048)

	jwk, err := buildJWK(pub, "RS256", "test-kid", "sig")
	assert.NoError(t, err)
	assert.Equal(t, "RS256", jwk.Algorithm)
	assert.Equal(t, "test-kid", jwk.KeyID)
	assert.Equal(t, "sig", jwk.Use)
}

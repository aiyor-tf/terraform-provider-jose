// Copyright (c) Tze Liang
// SPDX-License-Identifier: MPL-2.0

package provider

import (
	"crypto/ecdsa"
	"crypto/ed25519"
	"crypto/rsa"
	"crypto/x509"
	"encoding/pem"
	"errors"

	"github.com/golang-jwt/jwt/v5"
	"github.com/hashicorp/terraform-plugin-framework/types"
)

type PrivateKey interface {
	sign(claims jwt.Claims, kid string) (string, error)
}

// Alg is used for RSA signing algorithm.
// The choice of algorithm used for signing is only applicable to RSA key.
// The signing key can be of RS256, RS384 or RS512.
type RSAPrivateKey struct {
	*rsa.PrivateKey
	Alg types.String
}

// For ECDSA private key type, the signing function will auto-select the
// appropriate signing algorithm to use.  Therefore no 'Alg' is stored.
type ECDSAPrivateKey struct {
	*ecdsa.PrivateKey
}

// For ECDSA private key type, there is only one algorithm to use.
// Therefore no 'Alg' is stored.
type EdDSAPrivateKey struct {
	ed25519.PrivateKey
}

func (k *RSAPrivateKey) sign(claims jwt.Claims, kid string) (string, error) {
	var token *jwt.Token
	var signingMethod jwt.SigningMethod

	switch k.Alg.ValueString() {
	case "RS256":
		signingMethod = jwt.SigningMethodRS256
	case "RS384":
		signingMethod = jwt.SigningMethodRS384
	case "RS512":
		signingMethod = jwt.SigningMethodRS512
	default:
		signingMethod = jwt.SigningMethodRS256
	}

	token = jwt.NewWithClaims(signingMethod, claims)
	if kid != "" {
		token.Header["kid"] = kid
	}

	return token.SignedString(k.PrivateKey)
}

func (k *ECDSAPrivateKey) sign(claims jwt.Claims, kid string) (string, error) {
	var token *jwt.Token
	var signingMethod jwt.SigningMethod

	// The choice of signing method is directly related to the private key length
	keySize := k.PrivateKey.Params().BitSize
	switch keySize {
	case 256:
		signingMethod = jwt.SigningMethodES256
	case 384:
		signingMethod = jwt.SigningMethodES384
	default:
		signingMethod = jwt.SigningMethodES512
	}

	token = jwt.NewWithClaims(signingMethod, claims)
	if kid != "" {
		token.Header["kid"] = kid
	}

	return token.SignedString(k.PrivateKey)
}

func (k *EdDSAPrivateKey) sign(claims jwt.Claims, kid string) (string, error) {
	var token *jwt.Token

	signingMethod := jwt.SigningMethodEdDSA

	token = jwt.NewWithClaims(signingMethod, claims)
	token.Header["kid"] = kid

	return token.SignedString(k.PrivateKey)
}

func parsePrivateKey(key []byte, alg types.String) (PrivateKey, error) {
	var privateKey PrivateKey

	// Parse PEM type
	block, _ := pem.Decode(key)
	if block == nil {
		return nil, errors.New("failed to parse PEM block containing the key")
	}

	// Attempt to assign the correct supported key type.
	// Currently only supports RSA, ECDSA, and EdDSA types.
	if pKey, err := x509.ParsePKCS1PrivateKey(block.Bytes); err == nil {
		privateKey = &RSAPrivateKey{pKey, alg}
	} else if pKey, err := x509.ParseECPrivateKey(block.Bytes); err == nil {
		privateKey = &ECDSAPrivateKey{pKey}
	} else if pKey, err := x509.ParsePKCS8PrivateKey(block.Bytes); err == nil {
		ed25519Key, ok := pKey.(ed25519.PrivateKey) // Type assertion
		if !ok {
			return nil, errors.New("error type assertion: ed25519.PrivateKey")
		} else {
			privateKey = &EdDSAPrivateKey{ed25519Key}
		}
	} else {
		return nil, errors.New("unsupported private key type")
	}

	return privateKey, nil
}

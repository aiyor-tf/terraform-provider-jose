// Copyright (c) Tze Liang
// SPDX-License-Identifier: MPL-2.0

package provider

import (
	"crypto/ecdsa"
	"crypto/ed25519"
	"crypto/elliptic"
	"crypto/rand"
	"crypto/rsa"
	"crypto/x509"
	"encoding/pem"
	"errors"
	"fmt"

	"github.com/go-jose/go-jose/v4"
	"github.com/hashicorp/terraform-plugin-framework/types"
)

// parsePublicKey parses a public key from PEM or JWK format.
func parsePublicKey(keyStr string) (interface{}, error) {
	// Try PEM
	block, _ := pem.Decode([]byte(keyStr))
	if block != nil {
		pub, err := x509.ParsePKIXPublicKey(block.Bytes)
		if err != nil {
			return nil, fmt.Errorf("failed to parse PKIX public key: %w", err)
		}
		return pub, nil
	}

	// Try JWK
	var jwk jose.JSONWebKey
	if err := jwk.UnmarshalJSON([]byte(keyStr)); err == nil {
		return jwk.Key, nil
	}

	return nil, fmt.Errorf("failed to parse public key as PEM or JWK")
}

// parsePrivateKey parses a PEM-encoded private key.
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

func marshalPrivateKeyToPEM(key interface{}) ([]byte, error) {
	var pemBlock *pem.Block
	switch k := key.(type) {
	case *rsa.PrivateKey:
		pemBlock = &pem.Block{
			Type:  "RSA PRIVATE KEY",
			Bytes: x509.MarshalPKCS1PrivateKey(k),
		}
	case *ecdsa.PrivateKey:
		b, err := x509.MarshalECPrivateKey(k)
		if err != nil {
			return nil, err
		}
		pemBlock = &pem.Block{
			Type:  "EC PRIVATE KEY",
			Bytes: b,
		}
	case ed25519.PrivateKey:
		b, err := x509.MarshalPKCS8PrivateKey(k)
		if err != nil {
			return nil, err
		}
		pemBlock = &pem.Block{
			Type:  "PRIVATE KEY",
			Bytes: b,
		}
	default:
		return nil, fmt.Errorf("unsupported private key type")
	}
	return pem.EncodeToMemory(pemBlock), nil
}

func marshalPublicKeyToPEM(key interface{}) ([]byte, error) {
	b, err := x509.MarshalPKIXPublicKey(key)
	if err != nil {
		return nil, err
	}
	pemBlock := &pem.Block{
		Type:  "PUBLIC KEY",
		Bytes: b,
	}
	return pem.EncodeToMemory(pemBlock), nil
}

// generateKey generates a key pair based on the algorithm and size.
func generateKey(alg string, size int) (interface{}, interface{}, error) {
	var privKey interface{}
	var pubKey interface{}
	var err error

	switch alg {
	case "RS256", "RS384", "RS512", "RSA":
		if size == 0 {
			size = 2048
		}
		privKey, err = rsa.GenerateKey(rand.Reader, size)
		if err != nil {
			return nil, nil, fmt.Errorf("failed to generate RSA key: %w", err)
		}
		pubKey = &privKey.(*rsa.PrivateKey).PublicKey

	case "ES256":
		privKey, err = ecdsa.GenerateKey(elliptic.P256(), rand.Reader)
		if err != nil {
			return nil, nil, fmt.Errorf("failed to generate ECDSA P-256 key: %w", err)
		}
		pubKey = &privKey.(*ecdsa.PrivateKey).PublicKey

	case "ES384":
		privKey, err = ecdsa.GenerateKey(elliptic.P384(), rand.Reader)
		if err != nil {
			return nil, nil, fmt.Errorf("failed to generate ECDSA P-384 key: %w", err)
		}
		pubKey = &privKey.(*ecdsa.PrivateKey).PublicKey

	case "ES512":
		privKey, err = ecdsa.GenerateKey(elliptic.P521(), rand.Reader)
		if err != nil {
			return nil, nil, fmt.Errorf("failed to generate ECDSA P-521 key: %w", err)
		}
		pubKey = &privKey.(*ecdsa.PrivateKey).PublicKey

	case "EdDSA":
		pub, priv, err := ed25519.GenerateKey(rand.Reader)
		if err != nil {
			return nil, nil, fmt.Errorf("failed to generate Ed25519 key: %w", err)
		}
		privKey = priv
		pubKey = pub

	default:
		return nil, nil, fmt.Errorf("unsupported algorithm: %s", alg)
	}

	return privKey, pubKey, nil
}

// buildJWK creates a JSON Web Key from a public key and metadata.
func buildJWK(key interface{}, alg, kid, use string) (*jose.JSONWebKey, error) {
	jwk := &jose.JSONWebKey{
		Key:       key,
		KeyID:     kid,
		Algorithm: alg,
		Use:       use,
	}

	// Deduce algorithm if not provided
	if jwk.Algorithm == "" {
		switch k := key.(type) {
		case *rsa.PublicKey:
			// Default to RS256 if not specified, or leave empty?
			// The original code set it if data.Alg was present.
			// If we want to be helpful, we can default, but maybe better to leave it if not specified.
			// However, for ECDSA/EdDSA, the original code set it explicitly based on type.
		case *ecdsa.PublicKey:
			bitLength := k.Curve.Params().BitSize
			if bitLength == 256 {
				jwk.Algorithm = "ES256"
			} else if bitLength == 384 {
				jwk.Algorithm = "ES384"
			} else if bitLength == 521 {
				jwk.Algorithm = "ES512"
			}
		case ed25519.PublicKey:
			jwk.Algorithm = "EdDSA"
		}
	}

	return jwk, nil
}

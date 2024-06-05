// Copyright (c) Tze Liang
// SPDX-License-Identifier: MPL-2.0

package provider

import (
	"crypto/ecdsa"
	"crypto/ed25519"
	"crypto/rsa"
	"crypto/x509"
	"encoding/pem"

	"github.com/go-jose/go-jose/v4"
)

// Create JWK.
func createJWK(data joseJwkResourceModel) ([]byte, error) {
	var (
		bitLength int
		err       error
		jwk       jose.JSONWebKey
	)

	// Parse the PEM-encoded public key
	block, _ := pem.Decode([]byte(data.PublicKey.ValueString()))
	if block == nil {
		return nil, err
	}

	pubKey, err := x509.ParsePKIXPublicKey(block.Bytes)
	if err != nil {
		return nil, err
	} else {
		jwk.Key = pubKey
	}

	switch k := pubKey.(type) {
	case *rsa.PublicKey:
		if data.Alg.ValueString() != "" {
			jwk.Algorithm = data.Alg.ValueString()
		}
	case *ecdsa.PublicKey:
		bitLength = k.Curve.Params().BitSize
		if bitLength == 256 {
			jwk.Algorithm = "ES256"
		} else if bitLength == 384 {
			jwk.Algorithm = "ES384"
		} else if bitLength == 521 {
			jwk.Algorithm = "ES512"
		} else {
			return nil, err
		}
	case ed25519.PublicKey:
		jwk.Algorithm = "EdDSA"
	default:
		return nil, err
	}

	if data.KID.ValueString() != "" {
		jwk.KeyID = data.KID.ValueString()
	}
	if data.Use.ValueString() != "" {
		jwk.Use = data.Use.ValueString()
	}

	jwkJSON, err := jwk.MarshalJSON()
	if err != nil {
		return nil, err
	}

	return jwkJSON, nil
}

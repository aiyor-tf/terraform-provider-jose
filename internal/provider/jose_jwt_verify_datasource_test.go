// Copyright (c) Tze Liang
// SPDX-License-Identifier: MPL-2.0

package provider

import (
	"testing"

	"github.com/hashicorp/terraform-plugin-testing/helper/resource"
)

func TestAccJoseJwtVerifyDataSource(t *testing.T) {
	resource.Test(t, resource.TestCase{
		PreCheck:                 func() { testAccPreCheck(t) },
		ProtoV6ProviderFactories: testAccProtoV6ProviderFactories,
		Steps: []resource.TestStep{
			{
				Config: `
					resource "jose_key" "test" {
						alg = "RS256"
					}

					resource "jose_jwk" "test" {
						alg        = "RS256"
						public_key = jose_key.test.public_key_pem
						kid        = "test-kid"
						use        = "sig"
					}

					resource "jose_jwks" "test" {
						jwks_properties = [
							{
								alg        = "RS256"
								public_key = jose_key.test.public_key_pem
								kid        = "test-kid"
								use        = "sig"
							}
						]
					}

					resource "jose_jwt_sign" "test" {
						alg         = "RS256"
						kid         = "test-kid"
						private_key = jose_key.test.private_key_pem
						claims_json = jsonencode({
							iss = "test-issuer"
							sub = "test-subject"
						})
					}

					data "jose_jwt_verify" "test" {
						jwt  = jose_jwt_sign.test.jwt
						jwks = jose_jwks.test.jwks
					}
				`,
				Check: resource.ComposeAggregateTestCheckFunc(
					resource.TestCheckResourceAttr("data.jose_jwt_verify.test", "claims.iss", "test-issuer"),
					resource.TestCheckResourceAttr("data.jose_jwt_verify.test", "claims.sub", "test-subject"),
					resource.TestCheckResourceAttr("data.jose_jwt_verify.test", "header.alg", "RS256"),
					resource.TestCheckResourceAttr("data.jose_jwt_verify.test", "header.kid", "test-kid"),
				),
			},
		},
	})
}

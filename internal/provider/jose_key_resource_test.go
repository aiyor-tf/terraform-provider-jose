// Copyright (c) Tze Liang
// SPDX-License-Identifier: MPL-2.0

package provider

import (
	"testing"

	"github.com/hashicorp/terraform-plugin-testing/helper/resource"
)

func TestAccJoseKeyResource(t *testing.T) {
	resource.Test(t, resource.TestCase{
		PreCheck:                 func() { testAccPreCheck(t) },
		ProtoV6ProviderFactories: testAccProtoV6ProviderFactories,
		Steps: []resource.TestStep{
			// RSA Key
			{
				Config: `
					resource "jose_key" "rsa" {
						alg  = "RS256"
						size = 2048
					}
				`,
				Check: resource.ComposeAggregateTestCheckFunc(
					resource.TestCheckResourceAttr("jose_key.rsa", "alg", "RS256"),
					resource.TestCheckResourceAttrSet("jose_key.rsa", "public_key_pem"),
					resource.TestCheckResourceAttrSet("jose_key.rsa", "private_key_pem"),
					resource.TestCheckResourceAttrSet("jose_key.rsa", "public_key_jwk"),
					resource.TestCheckResourceAttrSet("jose_key.rsa", "private_key_jwk"),
					resource.TestCheckResourceAttrSet("jose_key.rsa", "id"),
				),
			},
			// ECDSA Key
			{
				Config: `
					resource "jose_key" "ecdsa" {
						alg = "ES256"
					}
				`,
				Check: resource.ComposeAggregateTestCheckFunc(
					resource.TestCheckResourceAttr("jose_key.ecdsa", "alg", "ES256"),
					resource.TestCheckResourceAttrSet("jose_key.ecdsa", "public_key_pem"),
				),
			},
			// EdDSA Key
			{
				Config: `
					resource "jose_key" "eddsa" {
						alg = "EdDSA"
					}
				`,
				Check: resource.ComposeAggregateTestCheckFunc(
					resource.TestCheckResourceAttr("jose_key.eddsa", "alg", "EdDSA"),
					resource.TestCheckResourceAttrSet("jose_key.eddsa", "public_key_pem"),
				),
			},
		},
	})
}

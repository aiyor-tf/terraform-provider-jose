// Copyright (c) Tze Liang
// SPDX-License-Identifier: MPL-2.0

package provider

import (
	"testing"

	"github.com/hashicorp/terraform-plugin-testing/helper/resource"
)

func TestAccJoseJwtSignResource(t *testing.T) {
	resource.Test(t, resource.TestCase{
		PreCheck:                 func() { testAccPreCheck(t) },
		ProtoV6ProviderFactories: testAccProtoV6ProviderFactories,
		Steps: []resource.TestStep{
			{
				Config: `
					resource "jose_key" "test" {
						alg = "RS256"
					}

					resource "jose_jwt_sign" "test" {
						alg         = "RS256"
						private_key = jose_key.test.private_key_pem
						claims_json = jsonencode({
							iss = "test-issuer"
							sub = "test-subject"
						})
					}
				`,
				Check: resource.ComposeAggregateTestCheckFunc(
					resource.TestCheckResourceAttr("jose_jwt_sign.test", "alg", "RS256"),
					resource.TestCheckResourceAttrSet("jose_jwt_sign.test", "jwt"),
					resource.TestCheckResourceAttrSet("jose_jwt_sign.test", "id"),
				),
			},
		},
	})
}

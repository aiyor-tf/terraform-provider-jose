// Copyright (c) HashiCorp, Inc.
// Copyright (c) Tze Liang
// SPDX-License-Identifier: MPL-2.0

package provider

import (
	"strings"
	"testing"

	"github.com/aiyor-tf/terraform-provider-jose/internal/provider/fixtures"
	"github.com/hashicorp/terraform-plugin-testing/helper/resource"
)

func TestAccJoseJwkResource(t *testing.T) {
	resource.Test(t, resource.TestCase{
		PreCheck:                 func() { testAccPreCheck(t) },
		ProtoV6ProviderFactories: testAccProtoV6ProviderFactories,
		Steps: []resource.TestStep{
			// Create and Read testing
			{
				Config: `
					resource "jose_jwk" "test" {
						kid        = "this-is-a-key-id-for-rsa-key"
						alg        = "RS256"
						public_key = file("./fixtures/rsa-pub.pem")
						use        = "sig" 
					}
				`,
				Check: resource.ComposeAggregateTestCheckFunc(
					resource.TestCheckResourceAttr("jose_jwk.test", "jwk_b64", strings.TrimSpace(fixtures.B64JWKRSA)),
				),
			},
			// Update and Read testing
			// {
			// 	Config: testAccJoseJwkResourceConfig("two"),
			// 	Check: resource.ComposeAggregateTestCheckFunc(
			// 		resource.TestCheckResourceAttr("jose_jwk.test", "configurable_attribute", "two"),
			// 	),
			// },
			// Delete testing automatically occurs in TestCase
		},
	})
}

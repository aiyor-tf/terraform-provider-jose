// Copyright (c) Tze Liang
// SPDX-License-Identifier: MPL-2.0

package provider

import (
	"testing"

	"github.com/hashicorp/terraform-plugin-testing/helper/resource"
)

func TestAccJoseJweResource(t *testing.T) {
	resource.Test(t, resource.TestCase{
		PreCheck:                 func() { testAccPreCheck(t) },
		ProtoV6ProviderFactories: testAccProtoV6ProviderFactories,
		Steps: []resource.TestStep{
			{
				Config: `
					resource "jose_key" "test" {
						alg = "RS256"
					}

					resource "jose_jwe" "test" {
						plaintext  = "Secret Message"
						public_key = jose_key.test.public_key_pem
						alg        = "RSA-OAEP"
						enc        = "A256GCM"
					}
				`,
				Check: resource.ComposeAggregateTestCheckFunc(
					resource.TestCheckResourceAttrSet("jose_jwe.test", "jwe"),
					resource.TestCheckResourceAttrSet("jose_jwe.test", "id"),
				),
			},
		},
	})
}

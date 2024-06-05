terraform {
  required_providers {
    jose = {
      source = "registry.terraform.io/aiyor-tf/jose"
    }
  }
}

provider "jose" {}

resource "jose_jwk" "example_rsa" {
  kid        = "this-is-a-key-id-for-rsa-key"
  alg        = "RS256" # Optional, defaults to "RS256". Only applicable for RSA keys.
  public_key = file("../../../internal/provider/fixtures/rsa-pub.pem")
  use        = "sig"
}

resource "jose_jwk" "example_ecdsa" {
  kid        = "this-is-a-key-id-for-ecdsa-key"
  public_key = file("../../../internal/provider/fixtures/ecdsa-pub.pem")
  use        = "sig"
}

resource "jose_jwk" "example_ed25519" {
  kid        = "this-is-a-key-id-for-ed25519-key"
  public_key = file("../../../internal/provider/fixtures/ed25519-pub.pem")
  use        = "sig"
}

output "jwk_rsa" {
  value = jsondecode(jose_jwk.example_rsa.jwk)
}

output "jwk_ecdsa" {
  value = jsondecode(jose_jwk.example_ecdsa.jwk)
}

output "jwk_ed25519" {
  value = jsondecode(jose_jwk.example_ed25519.jwk)
}

output "jwk_rsa_b64" {
  value = jose_jwk.example_rsa.jwk_b64
}

output "jwk_ecdsa_b64" {
  value = jose_jwk.example_ecdsa.jwk_b64
}

output "jwk_ed25519_b64" {
  value = jose_jwk.example_ed25519.jwk_b64
}

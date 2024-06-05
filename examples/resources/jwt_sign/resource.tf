terraform {
  required_providers {
    jose = {
      source = "registry.terraform.io/aiyor-tf/jose"
    }
  }
}

provider "jose" {}


locals {
  claims = {
    "iss" : "https://example.com",
    "sub" : "1234567890",
    "aud" : "https://example.com",
    "iat" : 1516239022,
    "exp" : 1516239022,
    "custom_list" : [
      "foo",
      "bar"
    ]
  }
}

resource "jose_jwt_sign" "rsa" {
  private_key = file("./rsa.key")
  alg         = "RS256" # Optional
  kid         = "this-is-a-key-id-for-rsa-key"
  claims_json = jsonencode(local.claims)
}

resource "jose_jwt_sign" "ecdsa" {
  private_key = file("./ecdsa.key")
  kid         = "this-is-a-key-id-for-ecdsa-key"
  claims_json = jsonencode(local.claims)
}

resource "jose_jwt_sign" "ed25519" {
  private_key = file("./ed25519.key")
  kid         = "this-is-a-key-id-for-ed25519-key"
  claims_json = jsonencode(local.claims)
}

output "rsa_jwt" {
  value     = jose_jwt_sign.rsa.jwt
  sensitive = true
}

output "ecdsa_jwt" {
  value     = jose_jwt_sign.ecdsa.jwt
  sensitive = true
}

output "ed25519_jwt" {
  value     = jose_jwt_sign.ed25519.jwt
  sensitive = true
}

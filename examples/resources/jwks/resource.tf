terraform {
  required_providers {
    jose = {
      source = "registry.terraform.io/aiyor-tf/jose"
    }

    # local = {
    #   source = "hashicorp/local"
    #   version = "2.5.1"
    # }
  }
}

provider "jose" {}
# provider "local" {}

resource "jose_jwks" "example" {
  jwks_properties = [
    {
      kid        = "this-is-a-key-id-for-rsa-key"
      public_key = file("../../../internal/provider/fixtures/rsa-pub.pem")
      use        = "sig"
    },
    {
      kid        = "this-is-a-key-id-for-ecdsa-key"
      public_key = file("../../../internal/provider/fixtures/ecdsa-pub.pem")
      use        = "sig"
    },
    {
      kid        = "this-is-a-key-id-for-eddsa-key"
      public_key = file("../../../internal/provider/fixtures/ed25519-pub.pem")
      use        = "sig"
    },
  ]
}

output "jwks" {
  value = jsondecode(jose_jwks.example.jwks)
}

output "jwks_b64" {
  value = jose_jwks.example.jwks_b64
}
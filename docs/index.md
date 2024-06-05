---
page_title: "Provider: JOSE"
description: |-
  he JOSE provider aims to provide utilities to manage and interact with JSON Object Signing and Encryption (JOSE) operations.
---

# JOSE Provider

The JOSE provider aims to provide utilities to manage and interact with JSON Object Signing and Encryption (JOSE) operations. Currently, the provider's resources can do the following operations,
* JSON Web Key (JWK) and JSON Web Key Set (JWKS):
    * `jose_jwk`: Generate a JWK with customizable parameters, such as key ID, Usage.
    * `jose_jwks`: Generate a JWKS containing multiple JWKs, ideal for managing key sets.
    * `jose_jwt_sign`: Creates a JWT that can be signed by supported private keys (RSA, ECDSA, and ECDSA).



## Compatibility

Compatibility table between this provider, the [Terraform Plugin Protocol](https://www.terraform.io/plugin/how-terraform-works#terraform-plugin-protocol)
version it implements, and Terraform:

| JOSE Provider | Terraform Plugin Protocol | Terraform |
|:-------------:|:-------------------------:|:---------:|
|   `>= 0.x`    |            `6.0`          | `>= 1.0`  |
|   `>= 1.x`    |            `6.0`          | `>= 1.0`  |


This provider plugin is developed based on the [Terraform Plugin Framework](https://developer.hashicorp.com/terraform/plugin/framework).


## Requirements

- [Terraform](https://developer.hashicorp.com/terraform/downloads) >= 1.0


## Example Usage

```terraform
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
```
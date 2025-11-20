# Terraform Provider: JOSE

The JOSE provider offers utilities to manage and interact with JSON Object Signing and Encryption (JOSE) standards, including JWK, JWKS, JWT, and JWE. It allows you to generate keys, sign and verify tokens, and encrypt data directly within Terraform.

## Features

-   **Key Generation**: Native generation of RSA, ECDSA, and Ed25519 keys (`jose_key`).
-   **JWK/JWKS Management**: Create JSON Web Keys and Sets from PEM encoded keys (`jose_jwk`, `jose_jwks`).
-   **JWT Operations**: Sign (`jose_jwt_sign`) and verify (`jose_jwt_verify`) JSON Web Tokens.
-   **JWE Operations**: Encrypt payloads using JWE compact serialization (`jose_jwe`).
-   **Utilities**: Convert JWK to PEM (`jose_jwk_to_pem`) and fetch remote JWKS (`jose_jwks_uri`).

## Installation

### Local Development

To build and install the provider locally:

1.  Clone the repository.
2.  Run `make install` (requires Go 1.21+).
3.  Configure `~/.terraformrc` to use the local build (see [Developer Guide](docs/guides/dev-setup.md)).

## Resources

### `jose_key`

Generates a cryptographic key pair.

```hcl
resource "jose_key" "example" {
  alg  = "RS256" # Options: RS256, RS384, RS512, ES256, ES384, ES512, EdDSA
  size = 2048    # Optional, for RSA only
}
```

### `jose_jwk`

Creates a JSON Web Key (JWK) from a public key.

```hcl
resource "jose_jwk" "example" {
  alg        = "RS256"
  public_key = file("public.pem")
  kid        = "my-key-id"
  use        = "sig"
}
```

### `jose_jwks`

Creates a JSON Web Key Set (JWKS).

```hcl
resource "jose_jwks" "example" {
  jwks_properties = [
    {
      alg        = "RS256"
      public_key = file("rsa.pem")
      kid        = "rsa-1"
      use        = "sig"
    }
  ]
}
```

### `jose_jwt_sign`

Signs a JWT using a private key.

```hcl
resource "jose_jwt_sign" "example" {
  alg         = "RS256"
  private_key = file("private.pem")
  kid         = "my-key-id"
  claims_json = jsonencode({
    iss = "my-issuer"
    sub = "my-subject"
    exp = 1999999999
  })
}
```

### `jose_jwe`

Encrypts a payload into a JWE.

```hcl
resource "jose_jwe" "example" {
  plaintext  = "Secret Message"
  public_key = file("public.pem")
  alg        = "RSA-OAEP"
  enc        = "A256GCM"
}
```

## Data Sources

### `jose_jwt_verify`

Verifies a JWT and extracts claims.

```hcl
data "jose_jwt_verify" "example" {
  jwt      = var.jwt_token
  jwks     = var.jwks_json
  issuer   = "expected-issuer"
  audience = ["expected-audience"]
}
```

### `jose_jwk_to_pem`

Converts a JWK JSON string to PEM format.

```hcl
data "jose_jwk_to_pem" "example" {
  jwk = var.jwk_json
}
```

### `jose_jwks_uri`

Fetches a JWK Set from a remote URL.

```hcl
data "jose_jwks_uri" "example" {
  uri = "https://www.googleapis.com/oauth2/v3/certs"
}
```

## License

[Mozilla Public License v2.0](./LICENSE)
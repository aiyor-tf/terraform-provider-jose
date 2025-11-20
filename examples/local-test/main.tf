terraform {
  required_providers {
    jose = {
      source = "aiyor-tf/jose"
    }
  }
}

provider "jose" {}

# ==========================================
# 1. Key Generation (RSA, ECDSA, Ed25519)
# ==========================================

resource "jose_key" "rsa" {
  alg  = "RS256"
  size = 2048
}

resource "jose_key" "ecdsa" {
  alg = "ES256"
}

resource "jose_key" "eddsa" {
  alg = "EdDSA"
}

# ==========================================
# 2. JWK Creation from PEM
# ==========================================

resource "jose_jwk" "from_pem" {
  alg        = "RS256"
  public_key = jose_key.rsa.public_key_pem
  kid        = "rsa-key-id"
  use        = "sig"
}

# ==========================================
# 3. JWK Set (JWKS) Creation
# ==========================================

resource "jose_jwks" "example" {
  jwks_properties = [
    {
      alg        = "RS256"
      public_key = jose_key.rsa.public_key_pem
      kid        = "rsa-key-id"
      use        = "sig"
    },
    {
      alg        = "ES256"
      public_key = jose_key.ecdsa.public_key_pem
      kid        = "ecdsa-key-id"
      use        = "sig"
    }
  ]
}

# ==========================================
# 4. JWT Signing
# ==========================================

resource "jose_jwt_sign" "rsa_token" {
  alg         = "RS256"
  private_key = jose_key.rsa.private_key_pem
  kid         = "rsa-key-id"
  claims_json = jsonencode({
    iss = "terraform-provider-jose"
    sub = "rsa-test"
    aud = ["developers"]
    exp = 1999999999
  })
}

resource "jose_jwt_sign" "eddsa_token" {
  alg         = "EdDSA"
  private_key = jose_key.eddsa.private_key_pem
  claims_json = jsonencode({
    iss = "terraform-provider-jose"
    sub = "eddsa-test"
  })
}

# ==========================================
# 5. JWT Verification
# ==========================================

data "jose_jwt_verify" "verify_rsa" {
  jwt = jose_jwt_sign.rsa_token.jwt
  jwks = jsonencode({
    keys = [
      jsondecode(jose_key.rsa.public_key_jwk)
    ]
  })
  issuer   = "terraform-provider-jose"
  audience = ["developers"]
}

# ==========================================
# 6. JWE Encryption
# ==========================================

resource "jose_jwe" "encrypt_rsa" {
  plaintext  = "This is a secret message encrypted with RSA"
  public_key = jose_key.rsa.public_key_pem
  alg        = "RSA-OAEP"
  enc        = "A256GCM"
}

# ==========================================
# 7. JWK to PEM Conversion
# ==========================================

data "jose_jwk_to_pem" "convert" {
  jwk = jose_key.ecdsa.public_key_jwk
}

# ==========================================
# Outputs
# ==========================================

output "rsa_public_key_pem" {
  value = jose_key.rsa.public_key_pem
}

output "verified_claims" {
  value = data.jose_jwt_verify.verify_rsa.claims
}

output "jwe_output" {
  value = jose_jwe.encrypt_rsa.jwe
}

output "converted_pem" {
  value = data.jose_jwk_to_pem.convert.pem
}

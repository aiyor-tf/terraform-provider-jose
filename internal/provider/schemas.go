// Copyright (c) Tze Liang
// SPDX-License-Identifier: MPL-2.0

package provider

import (
	"github.com/hashicorp/terraform-plugin-framework-validators/stringvalidator"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/int64default"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/int64planmodifier"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/planmodifier"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/stringdefault"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/stringplanmodifier"
	"github.com/hashicorp/terraform-plugin-framework/schema/validator"
	"github.com/hashicorp/terraform-plugin-framework/types"
)

var (
	jwkSchema = map[string]schema.Attribute{
		"kid": schema.StringAttribute{
			Optional: true,
			PlanModifiers: []planmodifier.String{
				stringplanmodifier.RequiresReplace(),
			},
			Description: "Key ID.",
		},
		"alg": schema.StringAttribute{
			Computed:    true,
			Optional:    true,
			Description: "The algorithm used to sign the JWT - Applicable only for RSA keys. Supported values: RS256, RS384, RS512. Default to RS256",
			PlanModifiers: []planmodifier.String{
				stringplanmodifier.RequiresReplace(),
			},
			Default: stringdefault.StaticString("RS256"),
			Validators: []validator.String{
				stringvalidator.OneOf(
					"RS256", "RS384", "RS512"),
			},
		},
		"use": schema.StringAttribute{
			Computed:    true,
			Optional:    true,
			Description: "The key usage. Supported values: sig, enc. Default to sig",
			PlanModifiers: []planmodifier.String{
				stringplanmodifier.RequiresReplace(),
			},
			Default: stringdefault.StaticString("sig"),
			Validators: []validator.String{
				stringvalidator.OneOf(
					"sig", "enc",
				),
			},
		},
		"public_key": schema.StringAttribute{
			Optional: true,
			PlanModifiers: []planmodifier.String{
				stringplanmodifier.RequiresReplace(),
			},
			Description: "Public key in PEM format.",
		},
		"jwk": schema.StringAttribute{ // This is a stub. Not used in this resource.
			Computed:    true,
			Description: "The resulting JWK Set in JSON format.",
			PlanModifiers: []planmodifier.String{
				stringplanmodifier.RequiresReplace(),
			},
		},
		"jwk_b64": schema.StringAttribute{ // This is a stub. Not used in this resource.
			Computed:    true,
			Description: "The resulting JWK Set in JSON format.",
			PlanModifiers: []planmodifier.String{
				stringplanmodifier.RequiresReplace(),
			},
		},
		"id": schema.StringAttribute{
			Computed: true,
			PlanModifiers: []planmodifier.String{
				stringplanmodifier.UseStateForUnknown(),
			},
		},
	}

	jwtSchema = map[string]schema.Attribute{
		"private_key": schema.StringAttribute{
			Required:            true,
			Sensitive:           true,
			MarkdownDescription: "Private key in PEM format for signing JWT.",
			PlanModifiers: []planmodifier.String{
				stringplanmodifier.RequiresReplace(),
			},
		},
		"alg": schema.StringAttribute{
			Computed:            true,
			Optional:            true,
			MarkdownDescription: "Algorithm to use for signing JWT. Only applicable to RSA keys.Defaults to \"RS256\".  Accepted values: \"RS256\", \"384\", \"RS512\".",
			PlanModifiers: []planmodifier.String{
				stringplanmodifier.RequiresReplace(),
			},
			Default: stringdefault.StaticString("RS256"),
			Validators: []validator.String{
				stringvalidator.OneOf(
					"RS256", "RS384", "RS512"),
			},
		},
		"kid": schema.StringAttribute{
			Optional:            true,
			MarkdownDescription: "Key ID, in the context of JWK Set, to identify the key used.",
			PlanModifiers: []planmodifier.String{
				stringplanmodifier.RequiresReplace(),
			},
		},
		"claims_json": schema.StringAttribute{
			Required:            true,
			MarkdownDescription: "Claims (in JSON format) to be included in the JWT.",
			PlanModifiers: []planmodifier.String{
				stringplanmodifier.RequiresReplace(),
			},
		},
		"jwt": schema.StringAttribute{
			Computed:            true,
			Sensitive:           true,
			MarkdownDescription: "The resulting signed JWT in Base64url format.",
		},
		"id": schema.StringAttribute{
			Computed: true,
			PlanModifiers: []planmodifier.String{
				stringplanmodifier.UseStateForUnknown(),
			},
		},
	}

	jwtVerifySchema = map[string]schema.Attribute{
		"jwt": schema.StringAttribute{
			Required:            true,
			MarkdownDescription: "The JWT to verify.",
		},
		"jwks": schema.StringAttribute{
			Optional:            true,
			MarkdownDescription: "The JWK Set to use for verification.",
		},
		"issuer": schema.StringAttribute{
			Optional:            true,
			MarkdownDescription: "The expected issuer of the JWT.",
		},
		"audience": schema.ListAttribute{
			ElementType:         types.StringType,
			Optional:            true,
			MarkdownDescription: "The expected audience of the JWT.",
		},
		"claims": schema.MapAttribute{
			ElementType:         types.StringType,
			Computed:            true,
			MarkdownDescription: "The claims from the verified JWT.",
		},
		"header": schema.MapAttribute{
			ElementType:         types.StringType,
			Computed:            true,
			MarkdownDescription: "The header from the verified JWT.",
		},
	}

	jweSchema = map[string]schema.Attribute{
		"plaintext": schema.StringAttribute{
			Required:            true,
			Sensitive:           true,
			MarkdownDescription: "The plaintext to encrypt.",
			PlanModifiers: []planmodifier.String{
				stringplanmodifier.RequiresReplace(),
			},
		},
		"public_key": schema.StringAttribute{
			Required:            true,
			MarkdownDescription: "The public key to use for encryption (PEM or JWK).",
			PlanModifiers: []planmodifier.String{
				stringplanmodifier.RequiresReplace(),
			},
		},
		"alg": schema.StringAttribute{
			Required:            true,
			MarkdownDescription: "The key management algorithm (e.g., RSA-OAEP, ECDH-ES).",
			PlanModifiers: []planmodifier.String{
				stringplanmodifier.RequiresReplace(),
			},
		},
		"enc": schema.StringAttribute{
			Required:            true,
			MarkdownDescription: "The content encryption algorithm (e.g., A256GCM).",
			PlanModifiers: []planmodifier.String{
				stringplanmodifier.RequiresReplace(),
			},
		},
		"aad": schema.StringAttribute{
			Optional:            true,
			MarkdownDescription: "Additional Authenticated Data.",
			PlanModifiers: []planmodifier.String{
				stringplanmodifier.RequiresReplace(),
			},
		},
		"jwe": schema.StringAttribute{
			Computed:            true,
			MarkdownDescription: "The resulting JWE compact serialization.",
		},
		"id": schema.StringAttribute{
			Computed: true,
			PlanModifiers: []planmodifier.String{
				stringplanmodifier.UseStateForUnknown(),
			},
		},
	}

	keySchema = map[string]schema.Attribute{
		"alg": schema.StringAttribute{
			Required:            true,
			MarkdownDescription: "The algorithm to generate key for (e.g., RS256, ES256, EdDSA).",
			PlanModifiers: []planmodifier.String{
				stringplanmodifier.RequiresReplace(),
			},
		},
		"size": schema.Int64Attribute{
			Optional:            true,
			MarkdownDescription: "Key size in bits (only for RSA). Default is 2048.",
			PlanModifiers: []planmodifier.Int64{
				int64planmodifier.RequiresReplace(),
			},
			Default: int64default.StaticInt64(2048),
		},
		"public_key_pem": schema.StringAttribute{
			Computed:            true,
			MarkdownDescription: "The generated public key in PEM format.",
		},
		"private_key_pem": schema.StringAttribute{
			Computed:            true,
			Sensitive:           true,
			MarkdownDescription: "The generated private key in PEM format.",
		},
		"public_key_jwk": schema.StringAttribute{
			Computed:            true,
			MarkdownDescription: "The generated public key in JWK JSON format.",
		},
		"private_key_jwk": schema.StringAttribute{
			Computed:            true,
			Sensitive:           true,
			MarkdownDescription: "The generated private key in JWK JSON format.",
		},
		"id": schema.StringAttribute{
			Computed: true,
			PlanModifiers: []planmodifier.String{
				stringplanmodifier.UseStateForUnknown(),
			},
		},
	}

	jwkToPemSchema = map[string]schema.Attribute{
		"jwk": schema.StringAttribute{
			Required:            true,
			MarkdownDescription: "The JWK JSON string to convert.",
		},
		"pem": schema.StringAttribute{
			Computed:            true,
			MarkdownDescription: "The resulting PEM encoded key.",
		},
	}

	jwksUriSchema = map[string]schema.Attribute{
		"uri": schema.StringAttribute{
			Required:            true,
			MarkdownDescription: "The URI to fetch the JWK Set from.",
		},
		"jwks": schema.StringAttribute{
			Computed:            true,
			MarkdownDescription: "The fetched JWK Set JSON.",
		},
	}
)

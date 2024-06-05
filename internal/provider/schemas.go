// Copyright (c) Tze Liang
// SPDX-License-Identifier: MPL-2.0

package provider

import (
	"github.com/hashicorp/terraform-plugin-framework-validators/stringvalidator"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/planmodifier"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/stringdefault"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/stringplanmodifier"
	"github.com/hashicorp/terraform-plugin-framework/schema/validator"
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
	}
)

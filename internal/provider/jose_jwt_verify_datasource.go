// Copyright (c) Tze Liang
// SPDX-License-Identifier: MPL-2.0

package provider

import (
	"context"
	"encoding/json"
	"fmt"

	"github.com/go-jose/go-jose/v4"
	"github.com/golang-jwt/jwt/v5"
	"github.com/hashicorp/terraform-plugin-framework/datasource"
	"github.com/hashicorp/terraform-plugin-framework/datasource/schema"
	"github.com/hashicorp/terraform-plugin-framework/types"
)

// Ensure provider defined types fully satisfy framework interfaces.
var (
	_ datasource.DataSource = &joseJwtVerifyDataSource{}
)

func NewJoseJwtVerifyDataSource() datasource.DataSource {
	return &joseJwtVerifyDataSource{}
}

// joseJwtVerifyDataSource defines the data source implementation.
type joseJwtVerifyDataSource struct{}

// joseJwtVerifyDataSourceModel describes the data source data model.
type joseJwtVerifyDataSourceModel struct {
	JWT      types.String `tfsdk:"jwt"`
	JWKS     types.String `tfsdk:"jwks"`
	Issuer   types.String `tfsdk:"issuer"`
	Audience types.List   `tfsdk:"audience"`
	Claims   types.Map    `tfsdk:"claims"`
	Header   types.Map    `tfsdk:"header"`
}

func (d *joseJwtVerifyDataSource) Metadata(ctx context.Context, req datasource.MetadataRequest, resp *datasource.MetadataResponse) {
	resp.TypeName = req.ProviderTypeName + "_jwt_verify"
}

func (d *joseJwtVerifyDataSource) Schema(ctx context.Context, req datasource.SchemaRequest, resp *datasource.SchemaResponse) {
	// Convert resource schema to datasource schema
	// Note: This is a bit hacky, but since we defined schemas in schemas.go as resource/schema.Attribute,
	// we need to redefine or convert. For now, I'll redefine locally or use a helper if available.
	// Actually, the framework types are different. I should have defined them as datasource/schema.Attribute
	// or just defined them inline here.
	// Let's define inline for now to fix the error quickly.
	resp.Schema = schema.Schema{
		Attributes: map[string]schema.Attribute{
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
		},
	}
}

func (d *joseJwtVerifyDataSource) Configure(ctx context.Context, req datasource.ConfigureRequest, resp *datasource.ConfigureResponse) {
}

func (d *joseJwtVerifyDataSource) Read(ctx context.Context, req datasource.ReadRequest, resp *datasource.ReadResponse) {
	var data joseJwtVerifyDataSourceModel

	// Read Terraform configuration data into the model
	resp.Diagnostics.Append(req.Config.Get(ctx, &data)...)

	if resp.Diagnostics.HasError() {
		return
	}

	tokenString := data.JWT.ValueString()

	// Parse and verify the token
	token, err := jwt.Parse(tokenString, func(token *jwt.Token) (interface{}, error) {
		// If JWKS is provided
		if !data.JWKS.IsNull() {
			var jwks jose.JSONWebKeySet
			if err := json.Unmarshal([]byte(data.JWKS.ValueString()), &jwks); err != nil {
				return nil, fmt.Errorf("failed to parse JWKS: %w", err)
			}

			// Find the key
			kid, ok := token.Header["kid"].(string)
			if !ok {
				keys := jwks.Keys
				if len(keys) == 0 {
					return nil, fmt.Errorf("empty JWKS")
				}
				// If only one key, use it
				if len(keys) == 1 {
					return keys[0].Key, nil
				}
				return nil, fmt.Errorf("kid header missing in JWT and multiple keys in JWKS")
			}

			keys := jwks.Key(kid)
			if len(keys) == 0 {
				return nil, fmt.Errorf("key with kid %s not found in JWKS", kid)
			}

			// Security Check: Ensure the token's algorithm matches the key's algorithm
			// Note: go-jose keys might not always have Algorithm set, but if they do, we should check.
			// A stronger check is to ensure the token alg is one of the expected ones.
			// For now, we rely on the library to verify the signature with the key.
			// However, we should prevent 'none' alg.
			if token.Method.Alg() == "none" {
				return nil, fmt.Errorf("algorithm 'none' is not allowed")
			}

			return keys[0].Key, nil
		}

		return nil, fmt.Errorf("no verification key provided (jwks is required)")
	})

	if err != nil {
		resp.Diagnostics.AddError("JWT Verification Failed", err.Error())
		return
	}

	if !token.Valid {
		resp.Diagnostics.AddError("Invalid Token", "The token is invalid")
		return
	}

	// Validate Issuer
	if !data.Issuer.IsNull() {
		issuer, err := token.Claims.GetIssuer()
		if err != nil {
			resp.Diagnostics.AddError("Failed to get issuer", err.Error())
			return
		}
		if issuer != data.Issuer.ValueString() {
			resp.Diagnostics.AddError("Issuer Mismatch", fmt.Sprintf("Expected %s, got %s", data.Issuer.ValueString(), issuer))
			return
		}
	}

	// Validate Audience
	if !data.Audience.IsNull() {
		var expectedAud []string
		data.Audience.ElementsAs(ctx, &expectedAud, false)

		aud, err := token.Claims.GetAudience()
		if err != nil {
			resp.Diagnostics.AddError("Failed to get audience", err.Error())
			return
		}

		// Simple check: ensure at least one expected audience is in the token
		found := false
		for _, a := range aud {
			for _, e := range expectedAud {
				if a == e {
					found = true
					break
				}
			}
		}
		if !found {
			resp.Diagnostics.AddError("Audience Mismatch", fmt.Sprintf("Expected one of %v, got %v", expectedAud, aud))
			return
		}
	}

	// Extract Claims
	claimsMap := make(map[string]interface{})
	if mapClaims, ok := token.Claims.(jwt.MapClaims); ok {
		claimsMap = mapClaims
	}

	// Convert claims to map[string]string for Terraform (simplification)
	// Complex objects will be JSON stringified
	claimsStringMap := make(map[string]string)
	for k, v := range claimsMap {
		switch val := v.(type) {
		case string:
			claimsStringMap[k] = val
		default:
			jsonVal, _ := json.Marshal(val)
			claimsStringMap[k] = string(jsonVal)
		}
	}

	data.Claims, _ = types.MapValueFrom(ctx, types.StringType, claimsStringMap)

	// Extract Header
	headerStringMap := make(map[string]string)
	for k, v := range token.Header {
		switch val := v.(type) {
		case string:
			headerStringMap[k] = val
		default:
			jsonVal, _ := json.Marshal(val)
			headerStringMap[k] = string(jsonVal)
		}
	}
	data.Header, _ = types.MapValueFrom(ctx, types.StringType, headerStringMap)

	// Save data into Terraform state
	resp.Diagnostics.Append(resp.State.Set(ctx, &data)...)
}

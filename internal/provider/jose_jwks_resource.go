// Copyright (c) HashiCorp, Inc.
// Copyright (c) Tze Liang
// SPDX-License-Identifier: MPL-2.0

package provider

import (
	"context"
	"encoding/base64"
	"encoding/json"

	"github.com/hashicorp/terraform-plugin-framework/path"
	"github.com/hashicorp/terraform-plugin-framework/resource"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema"
	"github.com/hashicorp/terraform-plugin-framework/types"
	"github.com/hashicorp/terraform-plugin-log/tflog"
)

// jwtResource defines the resource implementation.
type joseJwksResource struct{}

// jwtResourceModel describes the resource data model.
type joseJwksResourceModel struct {
	JWKSProperties []joseJwkResourceModel `tfsdk:"jwks_properties"`
	JWKS           types.String           `tfsdk:"jwks"`
	JWKSBase64     types.String           `tfsdk:"jwks_b64"`
}

// type joseJwksModel struct {
// 	KID       types.String `tfsdk:"kid"`
// 	PublicKey types.String `tfsdk:"public_key"`
// 	Use       types.String `tfsdk:"use"`
// }

type JWKSet struct {
	Keys []json.RawMessage `json:"keys"` // Store keys as raw JSON
}

// Ensure provider defined types fully satisfy framework interfaces.
var (
	_ resource.Resource                = &joseJwksResource{}
	_ resource.ResourceWithImportState = &joseJwksResource{}
)

func NewJoseJwksResource() resource.Resource {
	return &joseJwksResource{}
}

func (r *joseJwksResource) Metadata(ctx context.Context, req resource.MetadataRequest, resp *resource.MetadataResponse) {
	resp.TypeName = req.ProviderTypeName + "_jwks"
}

func (r *joseJwksResource) Schema(ctx context.Context, req resource.SchemaRequest, resp *resource.SchemaResponse) {
	resp.Schema = schema.Schema{
		// This description is used by the documentation generator and the language server.
		MarkdownDescription: "Create a JWK Set.",
		Attributes: map[string]schema.Attribute{
			"jwks_properties": schema.SetNestedAttribute{
				NestedObject: schema.NestedAttributeObject{
					Attributes: jwkSchema,
				},
				MarkdownDescription: "JWK Set configurable attribute.",
				Required:            true,
			},
			"jwks": schema.StringAttribute{
				Computed:    true,
				Description: "The resulting JWK Set in JSON format.",
			},
			"jwks_b64": schema.StringAttribute{
				Computed:    true,
				Description: "The resulting JWK Set encoded in base64.",
			},
		},
	}
}

func (r *joseJwksResource) Configure(ctx context.Context, req resource.ConfigureRequest, resp *resource.ConfigureResponse) {
}

func (r *joseJwksResource) Create(ctx context.Context, req resource.CreateRequest, resp *resource.CreateResponse) {
	var data joseJwksResourceModel

	// Create a JWKSet to store all keys
	jwkSet := JWKSet{Keys: make([]json.RawMessage, 0)}

	// Read plan data into the model
	// Now the model '&data' holds the plan data
	resp.Diagnostics.Append(req.Plan.Get(ctx, &data)...)

	if resp.Diagnostics.HasError() {
		return
	}

	//for _, item := range data.JWKSProperties {}
	for _, item := range data.JWKSProperties {
		jwkJSON, err := createJWK(item)

		if err != nil {
			resp.Diagnostics.AddError("Error creating JWK:", err.Error())
			return
		}

		// Append the raw JSON to the JWKSet.Keys
		jwkSet.Keys = append(jwkSet.Keys, jwkJSON)
	}

	// Marshal the JWKSet to JSON
	jwkSetJSON, err := json.Marshal(jwkSet)
	if err != nil {
		resp.Diagnostics.AddError("Error marshalling JWK Set to JSON:", err.Error())
		return
	}

	// Save the JWKSet result into JWKS
	data.JWKS = types.StringValue(string(jwkSetJSON))

	// Save jwkJSON as data.JWKBase64 encoded in Base64
	data.JWKSBase64 = types.StringValue(base64.StdEncoding.EncodeToString(jwkSetJSON))

	// Write logs using the tflog package
	// Documentation: https://terraform.io/plugin/log
	tflog.Trace(ctx, "created a resource")

	// Save data into Terraform state
	resp.Diagnostics.Append(resp.State.Set(ctx, &data)...)
}

func (r *joseJwksResource) Read(ctx context.Context, req resource.ReadRequest, resp *resource.ReadResponse) {
	var data joseJwksResourceModel

	// Read Terraform prior state data into the model
	resp.Diagnostics.Append(req.State.Get(ctx, &data)...)

	if resp.Diagnostics.HasError() {
		return
	}

	resp.Diagnostics.Append(resp.State.Set(ctx, &data)...)
}

func (r *joseJwksResource) Update(ctx context.Context, req resource.UpdateRequest, resp *resource.UpdateResponse) {
	var data joseJwksResourceModel

	// Read Terraform plan data into the model
	resp.Diagnostics.Append(req.Plan.Get(ctx, &data)...)

	if resp.Diagnostics.HasError() {
		return
	}

	resp.Diagnostics.Append(resp.State.Set(ctx, &data)...)
}

func (r *joseJwksResource) Delete(ctx context.Context, req resource.DeleteRequest, resp *resource.DeleteResponse) {
	var data joseJwksResourceModel

	// Read Terraform prior state data into the model
	resp.Diagnostics.Append(req.State.Get(ctx, &data)...)

	if resp.Diagnostics.HasError() {
		return
	}
}

func (r *joseJwksResource) ImportState(ctx context.Context, req resource.ImportStateRequest, resp *resource.ImportStateResponse) {
	resource.ImportStatePassthroughID(ctx, path.Root("jwks"), req, resp)
}

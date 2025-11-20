// Copyright (c) HashiCorp, Inc.
// Copyright (c) Tze Liang
// SPDX-License-Identifier: MPL-2.0

package provider

import (
	"context"
	"encoding/base64"

	"github.com/google/uuid"
	"github.com/hashicorp/terraform-plugin-framework/resource"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema"
	"github.com/hashicorp/terraform-plugin-framework/types"
	"github.com/hashicorp/terraform-plugin-log/tflog"
)

// jwtResource defines the resource implementation.
type joseJwkResource struct{}

// jwtResourceModel describes the resource data model.
type joseJwkResourceModel struct {
	PublicKey types.String `tfsdk:"public_key"`
	Alg       types.String `tfsdk:"alg"`
	KID       types.String `tfsdk:"kid"`
	Use       types.String `tfsdk:"use"`
	JWK       types.String `tfsdk:"jwk"`
	JWKBase64 types.String `tfsdk:"jwk_b64"`
	ID        types.String `tfsdk:"id"`
}

// Ensure provider defined types fully satisfy framework interfaces.
var (
	_ resource.Resource = &joseJwkResource{}
	_ resource.Resource = &joseJwkResource{}
)

func NewJoseJwkResource() resource.Resource {
	return &joseJwkResource{}
}

func (r *joseJwkResource) Metadata(ctx context.Context, req resource.MetadataRequest, resp *resource.MetadataResponse) {
	resp.TypeName = req.ProviderTypeName + "_jwk"
}

func (r *joseJwkResource) Schema(ctx context.Context, req resource.SchemaRequest, resp *resource.SchemaResponse) {
	resp.Schema = schema.Schema{
		MarkdownDescription: "Creates a JWK from a public key.",
		Attributes:          jwkSchema,
	}
}

func (r *joseJwkResource) Configure(ctx context.Context, req resource.ConfigureRequest, resp *resource.ConfigureResponse) {
}

func (r *joseJwkResource) Create(ctx context.Context, req resource.CreateRequest, resp *resource.CreateResponse) {
	var data joseJwkResourceModel

	// Read plan data into the model
	// Now the model '&data' holds the plan data
	resp.Diagnostics.Append(req.Plan.Get(ctx, &data)...)

	if resp.Diagnostics.HasError() {
		return
	}

	pubKey, err := parsePublicKey(data.PublicKey.ValueString())
	if err != nil {
		resp.Diagnostics.AddError("Invalid public key", err.Error())
		return
	}

	jwk, err := buildJWK(pubKey, data.Alg.ValueString(), data.KID.ValueString(), data.Use.ValueString())
	if err != nil {
		resp.Diagnostics.AddError("Error building JWK", err.Error())
		return
	}

	jwkJSON, err := jwk.MarshalJSON()
	if err != nil {
		resp.Diagnostics.AddError("Error marshaling JWK", err.Error())
		return
	}

	data.JWK = types.StringValue(string(jwkJSON))

	// Save jwkJSON as data.JWKBase64 encoded in Base64
	data.JWKBase64 = types.StringValue(base64.StdEncoding.EncodeToString(jwkJSON))

	data.ID = types.StringValue(uuid.NewString())

	// Write logs using the tflog package
	// Documentation: https://terraform.io/plugin/log
	tflog.Trace(ctx, "created a resource")

	// Save data into Terraform state
	resp.Diagnostics.Append(resp.State.Set(ctx, &data)...)
}

func (r *joseJwkResource) Read(ctx context.Context, req resource.ReadRequest, resp *resource.ReadResponse) {
	var data joseJwkResourceModel

	// Read Terraform prior state data into the model
	resp.Diagnostics.Append(req.State.Get(ctx, &data)...)

	if resp.Diagnostics.HasError() {
		return
	}

	resp.Diagnostics.Append(resp.State.Set(ctx, &data)...)
}

func (r *joseJwkResource) Update(ctx context.Context, req resource.UpdateRequest, resp *resource.UpdateResponse) {
	var data joseJwkResourceModel

	// Read Terraform plan data into the model
	resp.Diagnostics.Append(req.Plan.Get(ctx, &data)...)

	if resp.Diagnostics.HasError() {
		return
	}

	resp.Diagnostics.Append(resp.State.Set(ctx, &data)...)
}

func (r *joseJwkResource) Delete(ctx context.Context, req resource.DeleteRequest, resp *resource.DeleteResponse) {
	var data joseJwkResourceModel

	// Read Terraform prior state data into the model
	resp.Diagnostics.Append(req.State.Get(ctx, &data)...)

	if resp.Diagnostics.HasError() {
		return
	}

}

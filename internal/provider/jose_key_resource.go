// Copyright (c) Tze Liang
// SPDX-License-Identifier: MPL-2.0

package provider

import (
	"context"

	"github.com/google/uuid"
	"github.com/hashicorp/terraform-plugin-framework/resource"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema"
	"github.com/hashicorp/terraform-plugin-framework/types"
	"github.com/hashicorp/terraform-plugin-log/tflog"
)

// Ensure provider defined types fully satisfy framework interfaces.
var (
	_ resource.Resource = &joseKeyResource{}
	_ resource.Resource = &joseKeyResource{}
)

func NewJoseKeyResource() resource.Resource {
	return &joseKeyResource{}
}

// joseKeyResource defines the resource implementation.
type joseKeyResource struct{}

// joseKeyResourceModel describes the resource data model.
type joseKeyResourceModel struct {
	Alg           types.String `tfsdk:"alg"`
	Size          types.Int64  `tfsdk:"size"`
	PublicKeyPEM  types.String `tfsdk:"public_key_pem"`
	PrivateKeyPEM types.String `tfsdk:"private_key_pem"`
	PublicKeyJWK  types.String `tfsdk:"public_key_jwk"`
	PrivateKeyJWK types.String `tfsdk:"private_key_jwk"`
	ID            types.String `tfsdk:"id"`
}

func (r *joseKeyResource) Metadata(ctx context.Context, req resource.MetadataRequest, resp *resource.MetadataResponse) {
	resp.TypeName = req.ProviderTypeName + "_key"
}

func (r *joseKeyResource) Schema(ctx context.Context, req resource.SchemaRequest, resp *resource.SchemaResponse) {
	resp.Schema = schema.Schema{
		MarkdownDescription: "Generates a cryptographic key (RSA, ECDSA, Ed25519).",
		Attributes:          keySchema,
	}
}

func (r *joseKeyResource) Configure(ctx context.Context, req resource.ConfigureRequest, resp *resource.ConfigureResponse) {
}

func (r *joseKeyResource) Create(ctx context.Context, req resource.CreateRequest, resp *resource.CreateResponse) {
	var data joseKeyResourceModel

	// Read Terraform plan data into the model
	resp.Diagnostics.Append(req.Plan.Get(ctx, &data)...)

	if resp.Diagnostics.HasError() {
		return
	}

	alg := data.Alg.ValueString()
	size := int(data.Size.ValueInt64())
	privKey, pubKey, err := generateKey(alg, size)
	if err != nil {
		resp.Diagnostics.AddError("Failed to generate key", err.Error())
		return
	}

	// Convert to PEM
	privPEM, err := marshalPrivateKeyToPEM(privKey)
	if err != nil {
		resp.Diagnostics.AddError("Failed to marshal private key to PEM", err.Error())
		return
	}
	pubPEM, err := marshalPublicKeyToPEM(pubKey)
	if err != nil {
		resp.Diagnostics.AddError("Failed to marshal public key to PEM", err.Error())
		return
	}

	// Convert to JWK
	privJWK, err := buildJWK(privKey, alg, "", "sig")
	if err != nil {
		resp.Diagnostics.AddError("Failed to build private JWK", err.Error())
		return
	}
	privJWKJSON, err := privJWK.MarshalJSON()
	if err != nil {
		resp.Diagnostics.AddError("Failed to marshal private key to JWK", err.Error())
		return
	}

	pubJWK, err := buildJWK(pubKey, alg, "", "sig")
	if err != nil {
		resp.Diagnostics.AddError("Failed to build public JWK", err.Error())
		return
	}
	pubJWKJSON, err := pubJWK.MarshalJSON()
	if err != nil {
		resp.Diagnostics.AddError("Failed to marshal public key to JWK", err.Error())
		return
	}

	data.PrivateKeyPEM = types.StringValue(string(privPEM))
	data.PublicKeyPEM = types.StringValue(string(pubPEM))
	data.PrivateKeyJWK = types.StringValue(string(privJWKJSON))
	data.PublicKeyJWK = types.StringValue(string(pubJWKJSON))
	data.ID = types.StringValue(uuid.NewString())

	// Write logs using the tflog package
	tflog.Trace(ctx, "created a resource")

	// Save data into Terraform state
	resp.Diagnostics.Append(resp.State.Set(ctx, &data)...)
}

func (r *joseKeyResource) Read(ctx context.Context, req resource.ReadRequest, resp *resource.ReadResponse) {
	var data joseKeyResourceModel

	// Read Terraform prior state data into the model
	resp.Diagnostics.Append(req.State.Get(ctx, &data)...)

	if resp.Diagnostics.HasError() {
		return
	}

	// Save updated data into Terraform state
	resp.Diagnostics.Append(resp.State.Set(ctx, &data)...)
}

func (r *joseKeyResource) Update(ctx context.Context, req resource.UpdateRequest, resp *resource.UpdateResponse) {
	// All attributes force replacement
}

func (r *joseKeyResource) Delete(ctx context.Context, req resource.DeleteRequest, resp *resource.DeleteResponse) {
}

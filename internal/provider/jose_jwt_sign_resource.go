// Copyright (c) HashiCorp, Inc.
// Copyright (c) Tze Liang
// SPDX-License-Identifier: MPL-2.0

package provider

import (
	"context"
	"encoding/json"

	"github.com/golang-jwt/jwt/v5"
	"github.com/hashicorp/terraform-plugin-framework/path"
	"github.com/hashicorp/terraform-plugin-framework/resource"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema"
	"github.com/hashicorp/terraform-plugin-framework/types"
	"github.com/hashicorp/terraform-plugin-log/tflog"
)

// Ensure provider defined types fully satisfy framework interfaces.
var (
	_ resource.Resource                = &joseJwtSignResource{}
	_ resource.ResourceWithImportState = &joseJwtSignResource{}
)

func NewJoseJwtSignResource() resource.Resource {
	return &joseJwtSignResource{}
}

// jwtResource defines the resource implementation.
type joseJwtSignResource struct{}

// jwtResourceModel describes the resource data model.
type joseJwtSignResourceModel struct {
	PrivateKey types.String `tfsdk:"private_key"`
	Alg        types.String `tfsdk:"alg"`
	KID        types.String `tfsdk:"kid"`
	ClaimsJSON types.String `tfsdk:"claims_json"`
	JWT        types.String `tfsdk:"jwt"`
}

func (r *joseJwtSignResource) Metadata(ctx context.Context, req resource.MetadataRequest, resp *resource.MetadataResponse) {
	resp.TypeName = req.ProviderTypeName + "_jwt_sign"
}

func (r *joseJwtSignResource) Schema(ctx context.Context, req resource.SchemaRequest, resp *resource.SchemaResponse) {
	resp.Schema = schema.Schema{
		// This description is used by the documentation generator and the language server.
		MarkdownDescription: "Creates a JWT. Supports RSA, ECDSA and EdDSA keys.",
		Attributes:          jwtSchema,
	}
}

func (r *joseJwtSignResource) Configure(ctx context.Context, req resource.ConfigureRequest, resp *resource.ConfigureResponse) {
}

func (r *joseJwtSignResource) Create(ctx context.Context, req resource.CreateRequest, resp *resource.CreateResponse) {
	var (
		data joseJwtSignResourceModel
	)

	// Read Terraform plan data into the model
	resp.Diagnostics.Append(req.Plan.Get(ctx, &data)...)

	if resp.Diagnostics.HasError() {
		return
	}

	claims := jwt.MapClaims{}
	// Parse JSON claims
	if err := json.Unmarshal([]byte(data.ClaimsJSON.ValueString()), &claims); err != nil {
		resp.Diagnostics.AddError("Invalid claims JSON", err.Error())
		return
	}

	// Parse PEM into correct private key type
	privateKey, err := parsePrivateKey([]byte(data.PrivateKey.ValueString()), data.Alg)
	if err != nil {
		resp.Diagnostics.AddError("Invalid private key", err.Error())
		return
	}

	// Create the JWT token
	token, err := privateKey.sign(claims, data.KID.ValueString())
	if err != nil {
		resp.Diagnostics.AddError("Failed to sign JWT", err.Error())
		return
	}

	// For the purposes of this example code, hardcoding a response value to
	// save into the Terraform state.
	data.JWT = types.StringValue(token)

	// Write logs using the tflog package
	// Documentation: https://terraform.io/plugin/log
	tflog.Trace(ctx, "created a resource")

	// Save data into Terraform state
	resp.Diagnostics.Append(resp.State.Set(ctx, &data)...)
}

func (r *joseJwtSignResource) Read(ctx context.Context, req resource.ReadRequest, resp *resource.ReadResponse) {
	var data joseJwtSignResourceModel

	// Read Terraform prior state data into the model
	resp.Diagnostics.Append(req.State.Get(ctx, &data)...)

	if resp.Diagnostics.HasError() {
		return
	}

	// If applicable, this is a great opportunity to initialize any necessary
	// provider client data and make a call using it.
	// httpResp, err := r.client.Do(httpReq)
	// if err != nil {
	//     resp.Diagnostics.AddError("Client Error", fmt.Sprintf("Unable to read example, got error: %s", err))
	//     return
	// }

	// Save updated data into Terraform state
	resp.Diagnostics.Append(resp.State.Set(ctx, &data)...)
}

func (r *joseJwtSignResource) Update(ctx context.Context, req resource.UpdateRequest, resp *resource.UpdateResponse) {
	var data joseJwtSignResourceModel

	// Read Terraform plan data into the model
	resp.Diagnostics.Append(req.Plan.Get(ctx, &data)...)

	if resp.Diagnostics.HasError() {
		return
	}

	// Save updated data into Terraform state
	resp.Diagnostics.Append(resp.State.Set(ctx, &data)...)
}

func (r *joseJwtSignResource) Delete(ctx context.Context, req resource.DeleteRequest, resp *resource.DeleteResponse) {
	var data joseJwtSignResourceModel

	// Read Terraform prior state data into the model
	resp.Diagnostics.Append(req.State.Get(ctx, &data)...)

	if resp.Diagnostics.HasError() {
		return
	}
}

func (r *joseJwtSignResource) ImportState(ctx context.Context, req resource.ImportStateRequest, resp *resource.ImportStateResponse) {
	resource.ImportStatePassthroughID(ctx, path.Root("jwt"), req, resp)
}

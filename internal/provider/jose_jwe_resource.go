// Copyright (c) Tze Liang
// SPDX-License-Identifier: MPL-2.0

package provider

import (
	"context"

	"github.com/go-jose/go-jose/v4"
	"github.com/google/uuid"
	"github.com/hashicorp/terraform-plugin-framework/resource"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/planmodifier"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/stringplanmodifier"
	"github.com/hashicorp/terraform-plugin-framework/types"
	"github.com/hashicorp/terraform-plugin-log/tflog"
)

// Ensure provider defined types fully satisfy framework interfaces.
var (
	_ resource.Resource = &joseJweResource{}
	_ resource.Resource = &joseJweResource{}
)

func NewJoseJweResource() resource.Resource {
	return &joseJweResource{}
}

// joseJweResource defines the resource implementation.
type joseJweResource struct{}

// joseJweResourceModel describes the resource data model.
type joseJweResourceModel struct {
	Plaintext types.String `tfsdk:"plaintext"`
	PublicKey types.String `tfsdk:"public_key"`
	Alg       types.String `tfsdk:"alg"`
	Enc       types.String `tfsdk:"enc"`
	AAD       types.String `tfsdk:"aad"`
	JWE       types.String `tfsdk:"jwe"`
	ID        types.String `tfsdk:"id"`
}

func (r *joseJweResource) Metadata(ctx context.Context, req resource.MetadataRequest, resp *resource.MetadataResponse) {
	resp.TypeName = req.ProviderTypeName + "_jwe"
}

func (r *joseJweResource) Schema(ctx context.Context, req resource.SchemaRequest, resp *resource.SchemaResponse) {
	// Define schema inline to avoid type mismatch issues
	resp.Schema = schema.Schema{
		MarkdownDescription: "Creates a JWE (JSON Web Encryption) from a plaintext and public key.",
		Attributes: map[string]schema.Attribute{
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
		},
	}
}

func (r *joseJweResource) Configure(ctx context.Context, req resource.ConfigureRequest, resp *resource.ConfigureResponse) {
}

func (r *joseJweResource) Create(ctx context.Context, req resource.CreateRequest, resp *resource.CreateResponse) {
	var data joseJweResourceModel

	// Read Terraform plan data into the model
	resp.Diagnostics.Append(req.Plan.Get(ctx, &data)...)

	if resp.Diagnostics.HasError() {
		return
	}

	// Parse Public Key
	pubKey, err := parsePublicKey(data.PublicKey.ValueString())
	if err != nil {
		resp.Diagnostics.AddError("Invalid Public Key", err.Error())
		return
	}

	// Encrypt
	encrypter, err := jose.NewEncrypter(
		jose.ContentEncryption(data.Enc.ValueString()),
		jose.Recipient{Algorithm: jose.KeyAlgorithm(data.Alg.ValueString()), Key: pubKey},
		nil,
	)
	if err != nil {
		resp.Diagnostics.AddError("Failed to create encrypter", err.Error())
		return
	}

	object, err := encrypter.Encrypt([]byte(data.Plaintext.ValueString()))
	if err != nil {
		resp.Diagnostics.AddError("Encryption failed", err.Error())
		return
	}

	jweString, err := object.CompactSerialize()
	if err != nil {
		resp.Diagnostics.AddError("Serialization failed", err.Error())
		return
	}

	data.JWE = types.StringValue(jweString)
	data.ID = types.StringValue(uuid.NewString())

	// Write logs using the tflog package
	tflog.Trace(ctx, "created a resource")

	// Save data into Terraform state
	resp.Diagnostics.Append(resp.State.Set(ctx, &data)...)
}

func (r *joseJweResource) Read(ctx context.Context, req resource.ReadRequest, resp *resource.ReadResponse) {
	var data joseJweResourceModel

	// Read Terraform prior state data into the model
	resp.Diagnostics.Append(req.State.Get(ctx, &data)...)

	if resp.Diagnostics.HasError() {
		return
	}

	// Save updated data into Terraform state
	resp.Diagnostics.Append(resp.State.Set(ctx, &data)...)
}

func (r *joseJweResource) Update(ctx context.Context, req resource.UpdateRequest, resp *resource.UpdateResponse) {
	// Since all attributes require replace, Update shouldn't be called for changes,
	// but we implement it for safety.
	var data joseJweResourceModel

	// Read Terraform plan data into the model
	resp.Diagnostics.Append(req.Plan.Get(ctx, &data)...)

	if resp.Diagnostics.HasError() {
		return
	}

	resp.Diagnostics.Append(resp.State.Set(ctx, &data)...)
}

func (r *joseJweResource) Delete(ctx context.Context, req resource.DeleteRequest, resp *resource.DeleteResponse) {
	var data joseJweResourceModel

	// Read Terraform prior state data into the model
	resp.Diagnostics.Append(req.State.Get(ctx, &data)...)

	if resp.Diagnostics.HasError() {
		return
	}
}

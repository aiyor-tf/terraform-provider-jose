// Copyright (c) Tze Liang
// SPDX-License-Identifier: MPL-2.0

package provider

import (
	"context"

	"github.com/go-jose/go-jose/v4"
	"github.com/hashicorp/terraform-plugin-framework/datasource"
	"github.com/hashicorp/terraform-plugin-framework/datasource/schema"
	"github.com/hashicorp/terraform-plugin-framework/types"
)

// Ensure provider defined types fully satisfy framework interfaces.
var (
	_ datasource.DataSource = &joseJwkToPemDataSource{}
)

func NewJoseJwkToPemDataSource() datasource.DataSource {
	return &joseJwkToPemDataSource{}
}

// joseJwkToPemDataSource defines the data source implementation.
type joseJwkToPemDataSource struct{}

// joseJwkToPemDataSourceModel describes the data source data model.
type joseJwkToPemDataSourceModel struct {
	JWK types.String `tfsdk:"jwk"`
	PEM types.String `tfsdk:"pem"`
}

func (d *joseJwkToPemDataSource) Metadata(ctx context.Context, req datasource.MetadataRequest, resp *datasource.MetadataResponse) {
	resp.TypeName = req.ProviderTypeName + "_jwk_to_pem"
}

func (d *joseJwkToPemDataSource) Schema(ctx context.Context, req datasource.SchemaRequest, resp *datasource.SchemaResponse) {
	// Define schema inline to avoid type mismatch
	resp.Schema = schema.Schema{
		MarkdownDescription: "Converts a JWK JSON string to a PEM encoded key.",
		Attributes: map[string]schema.Attribute{
			"jwk": schema.StringAttribute{
				Required:            true,
				MarkdownDescription: "The JWK JSON string to convert.",
			},
			"pem": schema.StringAttribute{
				Computed:            true,
				MarkdownDescription: "The resulting PEM encoded key.",
			},
		},
	}
}

func (d *joseJwkToPemDataSource) Configure(ctx context.Context, req datasource.ConfigureRequest, resp *datasource.ConfigureResponse) {
}

func (d *joseJwkToPemDataSource) Read(ctx context.Context, req datasource.ReadRequest, resp *datasource.ReadResponse) {
	var data joseJwkToPemDataSourceModel

	// Read Terraform configuration data into the model
	resp.Diagnostics.Append(req.Config.Get(ctx, &data)...)

	if resp.Diagnostics.HasError() {
		return
	}

	jwkJSON := data.JWK.ValueString()
	var jwk jose.JSONWebKey
	if err := jwk.UnmarshalJSON([]byte(jwkJSON)); err != nil {
		resp.Diagnostics.AddError("Failed to parse JWK", err.Error())
		return
	}

	var pemBytes []byte
	var err error

	if jwk.IsPublic() {
		pemBytes, err = marshalPublicKeyToPEM(jwk.Key)
	} else {
		pemBytes, err = marshalPrivateKeyToPEM(jwk.Key)
	}

	if err != nil {
		resp.Diagnostics.AddError("Failed to marshal key to PEM", err.Error())
		return
	}

	data.PEM = types.StringValue(string(pemBytes))

	// Save data into Terraform state
	resp.Diagnostics.Append(resp.State.Set(ctx, &data)...)
}

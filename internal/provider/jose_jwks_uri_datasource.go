// Copyright (c) Tze Liang
// SPDX-License-Identifier: MPL-2.0

package provider

import (
	"context"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"time"

	"github.com/go-jose/go-jose/v4"
	"github.com/hashicorp/terraform-plugin-framework/datasource"
	"github.com/hashicorp/terraform-plugin-framework/datasource/schema"
	"github.com/hashicorp/terraform-plugin-framework/types"
)

// Ensure provider defined types fully satisfy framework interfaces.
var (
	_ datasource.DataSource = &joseJwksUriDataSource{}
)

func NewJoseJwksUriDataSource() datasource.DataSource {
	return &joseJwksUriDataSource{}
}

// joseJwksUriDataSource defines the data source implementation.
type joseJwksUriDataSource struct{}

// joseJwksUriDataSourceModel describes the data source data model.
type joseJwksUriDataSourceModel struct {
	URI  types.String `tfsdk:"uri"`
	JWKS types.String `tfsdk:"jwks"`
}

func (d *joseJwksUriDataSource) Metadata(ctx context.Context, req datasource.MetadataRequest, resp *datasource.MetadataResponse) {
	resp.TypeName = req.ProviderTypeName + "_jwks_uri"
}

func (d *joseJwksUriDataSource) Schema(ctx context.Context, req datasource.SchemaRequest, resp *datasource.SchemaResponse) {
	resp.Schema = schema.Schema{
		MarkdownDescription: "Fetches a JWK Set from a URI.",
		Attributes: map[string]schema.Attribute{
			"uri": schema.StringAttribute{
				Required:            true,
				MarkdownDescription: "The URI to fetch the JWK Set from.",
			},
			"jwks": schema.StringAttribute{
				Computed:            true,
				MarkdownDescription: "The fetched JWK Set JSON.",
			},
		},
	}
}

func (d *joseJwksUriDataSource) Configure(ctx context.Context, req datasource.ConfigureRequest, resp *datasource.ConfigureResponse) {
}

func (d *joseJwksUriDataSource) Read(ctx context.Context, req datasource.ReadRequest, resp *datasource.ReadResponse) {
	var data joseJwksUriDataSourceModel

	// Read Terraform configuration data into the model
	resp.Diagnostics.Append(req.Config.Get(ctx, &data)...)

	if resp.Diagnostics.HasError() {
		return
	}

	uri := data.URI.ValueString()

	client := &http.Client{
		Timeout: 10 * time.Second,
	}

	httpReq, err := http.NewRequestWithContext(ctx, "GET", uri, nil)
	if err != nil {
		resp.Diagnostics.AddError("Failed to create HTTP request", err.Error())
		return
	}

	httpResp, err := client.Do(httpReq)
	if err != nil {
		resp.Diagnostics.AddError("Failed to fetch JWKS", err.Error())
		return
	}
	defer httpResp.Body.Close()

	if httpResp.StatusCode != http.StatusOK {
		resp.Diagnostics.AddError("Failed to fetch JWKS", fmt.Sprintf("HTTP status %d", httpResp.StatusCode))
		return
	}

	// Limit response body to 10MB to prevent DoS
	body, err := io.ReadAll(io.LimitReader(httpResp.Body, 10*1024*1024))
	if err != nil {
		resp.Diagnostics.AddError("Failed to read response body", err.Error())
		return
	}

	// Validate it's a valid JWKS
	var jwks jose.JSONWebKeySet
	if err := json.Unmarshal(body, &jwks); err != nil {
		resp.Diagnostics.AddError("Invalid JWKS response", err.Error())
		return
	}

	// Normalize JSON
	normalized, _ := json.Marshal(jwks)

	data.JWKS = types.StringValue(string(normalized))

	// Save data into Terraform state
	resp.Diagnostics.Append(resp.State.Set(ctx, &data)...)
}

// Copyright (c) HashiCorp, Inc.
// Copyright (c) Tze Liang
// SPDX-License-Identifier: MPL-2.0

package provider

import (
	"context"

	"github.com/hashicorp/terraform-plugin-framework/datasource"
	"github.com/hashicorp/terraform-plugin-framework/provider"
	"github.com/hashicorp/terraform-plugin-framework/resource"
)

// Ensure jwtProvider satisfies various provider interfaces.
var (
	_ provider.Provider = &joseProvider{}
)

// New is a helper function to simplify provider server and testing implementation.
func New(version string) func() provider.Provider {
	return func() provider.Provider {
		return &joseProvider{
			version: version,
		}
	}
}

// jwtProvider defines the provider implementation.
type joseProvider struct {
	// version is set to the provider version on release, "dev" when the
	// provider is built and ran locally, and "test" when running acceptance
	// testing.
	version string
}

// Metadata returns the provider type name.
func (p *joseProvider) Metadata(ctx context.Context, req provider.MetadataRequest, resp *provider.MetadataResponse) {
	resp.TypeName = "jose"
	resp.Version = p.version
}

// Schema defines the provider-level schema for configuration data.
func (p *joseProvider) Schema(ctx context.Context, req provider.SchemaRequest, resp *provider.SchemaResponse) {
}

// Configure prepares a jwtProvider instance.
func (p *joseProvider) Configure(ctx context.Context, req provider.ConfigureRequest, resp *provider.ConfigureResponse) {
}

// Resource defines the resources implemented in the provider.
func (p *joseProvider) Resources(ctx context.Context) []func() resource.Resource {
	return []func() resource.Resource{
		NewJoseJwkResource,
		NewJoseJwksResource,
		NewJoseJwtSignResource,
	}
}

// DataSources defines the data sources implemented in the provider.
func (p *joseProvider) DataSources(ctx context.Context) []func() datasource.DataSource {
	return nil
}

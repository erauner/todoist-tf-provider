// Copyright (c) HashiCorp, Inc.
// SPDX-License-Identifier: MPL-2.0

package provider

import (
	"context"
	"os"

	client "github.com/andreaswwilson/terraform-provider-todoist/internal/client"
	"github.com/hashicorp/terraform-plugin-framework/datasource"
	"github.com/hashicorp/terraform-plugin-framework/path"
	"github.com/hashicorp/terraform-plugin-framework/provider"
	"github.com/hashicorp/terraform-plugin-framework/provider/schema"
	"github.com/hashicorp/terraform-plugin-framework/resource"
	"github.com/hashicorp/terraform-plugin-framework/types"
	"github.com/hashicorp/terraform-plugin-log/tflog"
)

// Ensure ScaffoldingProvider satisfies various provider interfaces.
var _ provider.Provider = &TodoistProvider{}

// TodoistProvider defines the provider implementation.
type TodoistProvider struct {
	// version is set to the provider version on release, "dev" when the
	// provider is built and ran locally, and "test" when running acceptance
	// testing.
	version string
}

// TodoistProviderModel describes the provider data model.
type TodoistProviderModel struct {
	Token types.String `tfsdk:"token"`
}

func (p *TodoistProvider) Metadata(ctx context.Context, req provider.MetadataRequest, resp *provider.MetadataResponse) {
	resp.TypeName = "todoist"
	resp.Version = p.version
}

func (p *TodoistProvider) Schema(ctx context.Context, req provider.SchemaRequest, resp *provider.SchemaResponse) {
	resp.Schema = schema.Schema{
		Attributes: map[string]schema.Attribute{
			"token": schema.StringAttribute{
				MarkdownDescription: "Todoist API token",
				Sensitive:           true,
				Optional:            true,
			},
		},
	}
}

func (p *TodoistProvider) Configure(ctx context.Context, req provider.ConfigureRequest, resp *provider.ConfigureResponse) {
	var data TodoistProviderModel

	resp.Diagnostics.Append(req.Config.Get(ctx, &data)...)

	if resp.Diagnostics.HasError() {
		return
	}

	if data.Token.IsUnknown() {
		resp.Diagnostics.AddAttributeError(
			path.Root("token"),
			"Unknown API token",
			"The provider cannot create the API client as there is an unknown configuration value for the API token. "+
				"Either target apply the source of the value first, set the value statically in the configuration, or use the TODOIST_TOKEN environment variable.",
		)
	}
	if resp.Diagnostics.HasError() {
		return
	}
	token := os.Getenv("TODOIST_TOKEN")
	if !data.Token.IsNull() {
		token = data.Token.ValueString()
	}

	if token == "" {
		resp.Diagnostics.AddAttributeError(
			path.Root("token"),
			"Missing API Token",
			"The provider cannot create the API client as there is a missing or empty value for the API Token. "+
				"Set the password value in the configuration or use the TODOIST_TOKEN environment variable. "+
				"If either is already set, ensure the value is not empty.",
		)
	}
	if resp.Diagnostics.HasError() {
		return
	}
	client, err := client.NewClient(token)
	tflog.Info(ctx, "Client setup ok!")
	if err != nil {
		resp.Diagnostics.AddError(
			"Unable to Create Todoist API Client",
			"An unexpected error occurred when creating the Todoist API client. "+
				"If the error is not clear, please contact the provider developers.\n\n"+
				"Todoist Client Error: "+err.Error(),
		)
		return
	}
	resp.DataSourceData = client
	resp.ResourceData = client
}

func (p *TodoistProvider) Resources(ctx context.Context) []func() resource.Resource {
	return []func() resource.Resource{
		NewExampleResource,
	}
}

func (p *TodoistProvider) DataSources(ctx context.Context) []func() datasource.DataSource {
	return []func() datasource.DataSource{
		NewProjectsDataSource,
	}
}

func New(version string) func() provider.Provider {
	return func() provider.Provider {
		return &TodoistProvider{
			version: version,
		}
	}
}

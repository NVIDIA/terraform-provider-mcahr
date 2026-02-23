// SPDX-FileCopyrightText: Copyright (c) 2025 NVIDIA CORPORATION & AFFILIATES. All rights reserved.
// SPDX-License-Identifier: Apache-2.0
//
// Licensed under the Apache License, Version 2.0 (the "License");
// you may not use this file except in compliance with the License.
// You may obtain a copy of the License at
//
// http://www.apache.org/licenses/LICENSE-2.0
//
// Unless required by applicable law or agreed to in writing, software
// distributed under the License is distributed on an "AS IS" BASIS,
// WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
// See the License for the specific language governing permissions and
// limitations under the License.

package provider

import (
	"context"
	"fmt"
	"os"

	"terraform/terraform-provider/provider/common"
	"terraform/terraform-provider/provider/common/log"
	"terraform/terraform-provider/provider/common/version"
	"terraform/terraform-provider/provider/external_api/client"
	"terraform/terraform-provider/provider/external_api/provider/backend_version"
	"terraform/terraform-provider/provider/tf/core/config"
	"terraform/terraform-provider/provider/tf/resource/action"
	"terraform/terraform-provider/provider/tf/resource/alarm"
	"terraform/terraform-provider/provider/tf/resource/bot"
	"terraform/terraform-provider/provider/tf/resource/dashboard"
	"terraform/terraform-provider/provider/tf/resource/file"
	"terraform/terraform-provider/provider/tf/resource/integration"
	"terraform/terraform-provider/provider/tf/resource/principal"
	"terraform/terraform-provider/provider/tf/resource/report_template"
	resourceres "terraform/terraform-provider/provider/tf/resource/resource"
	"terraform/terraform-provider/provider/tf/resource/runbook"
	"terraform/terraform-provider/provider/tf/resource/secret"
	"terraform/terraform-provider/provider/tf/resource/system_settings"
	"terraform/terraform-provider/provider/tf/resource/time_trigger"

	"github.com/hashicorp/terraform-plugin-framework/datasource"
	"github.com/hashicorp/terraform-plugin-framework/diag"
	"github.com/hashicorp/terraform-plugin-framework/function"
	"github.com/hashicorp/terraform-plugin-framework/provider"
	"github.com/hashicorp/terraform-plugin-framework/provider/schema"
	"github.com/hashicorp/terraform-plugin-framework/resource"
	"github.com/hashicorp/terraform-plugin-framework/types"
)

var (
	RenderedProviderName string // Will be set at build time
	ProviderShortName    string // Will be set at build time
	EnvVarsNamePrefix    string // Will be set at build time
)

const EMPTY_STRING = ""

// Default fallback version used when backend version cannot be fetched from API
const DEFAULT_FALLBACK_VERSION = "release-29.1.0"

// FrameworkProvider defines the provider implementation.
type FrameworkProvider struct {
	// version is set to the provider version on release, "dev" when the
	// provider is built and ran locally, and "test" when running acceptance
	// testing.
	version        string
	backendVersion *version.BackendVersion
}

var _ provider.Provider = &FrameworkProvider{}

// FrameworkProviderModel describes the provider data model.
type FrameworkProviderModel struct {
	URL        types.String `tfsdk:"url"`
	Token      types.String `tfsdk:"token"`
	Retries    types.Int64  `tfsdk:"retries"`
	MinVersion types.String `tfsdk:"min_version"`
}

func NewFrameworkProvider(versionStr string) func() provider.Provider {

	return func() provider.Provider {
		return &FrameworkProvider{
			version: versionStr,
			// backendVersion will be initialized in Configure method
		}
	}
}

func (p *FrameworkProvider) Metadata(ctx context.Context, req provider.MetadataRequest, resp *provider.MetadataResponse) {
	resp.TypeName = ProviderShortName
	resp.Version = p.version
}

func (p *FrameworkProvider) Schema(ctx context.Context, req provider.SchemaRequest, resp *provider.SchemaResponse) {
	resp.Schema = schema.Schema{
		Description: "",
		Attributes: map[string]schema.Attribute{
			"url": schema.StringAttribute{
				MarkdownDescription: "Customer-specific URL for the " + RenderedProviderName + " API server.",
				Optional:            true,
			},
			"token": schema.StringAttribute{
				MarkdownDescription: "Customer/user-specific authorization token for the " + RenderedProviderName + " API server. May be provided via `" + EnvVarsNamePrefix + "_TOKEN` env variable.",
				Optional:            true,
				Sensitive:           true,
			},
			"retries": schema.Int64Attribute{
				MarkdownDescription: "Number of retries for API calls, in case of e.g. transient network failures.",
				Optional:            true,
			},
			"min_version": schema.StringAttribute{
				MarkdownDescription: "Minimum version required on the " + RenderedProviderName + " backend (API server).",
				Optional:            true,
			},
		},
	}
}

func (p *FrameworkProvider) Resources(ctx context.Context) []func() resource.Resource {
	return []func() resource.Resource{
		func() resource.Resource {
			return action.NewActionResource()
		},
		func() resource.Resource {
			return runbook.NewRunbookResource("runbook")
		},
		func() resource.Resource {
			return runbook.NewRunbookResource("notebook")
		},
		func() resource.Resource {
			return time_trigger.NewTimeTriggerResource()
		},
		func() resource.Resource {
			return alarm.NewAlarmResource()
		},
		func() resource.Resource {
			return bot.NewBotResource()
		},
		func() resource.Resource {
			return integration.NewIntegrationResource()
		},
		func() resource.Resource {
			return file.NewFileResource()
		},
		func() resource.Resource {
			return principal.NewPrincipalResource()
		},
		func() resource.Resource {
			return secret.NewNVaultSecretResource()
		},
		func() resource.Resource {
			return system_settings.NewSystemSettingsResource()
		},
		func() resource.Resource {
			return report_template.NewReportTemplateResource()
		},
		func() resource.Resource {
			return dashboard.NewDashboardResource()
		},
		func() resource.Resource {
			return resourceres.NewResource()
		},
	}
}

func (p *FrameworkProvider) DataSources(ctx context.Context) []func() datasource.DataSource {
	return []func() datasource.DataSource{}
}

func (p *FrameworkProvider) Functions(ctx context.Context) []func() function.Function {
	return []func() function.Function{}
}

func (p *FrameworkProvider) Configure(ctx context.Context, req provider.ConfigureRequest, resp *provider.ConfigureResponse) {
	var data FrameworkProviderModel

	resp.Diagnostics.Append(req.Config.Get(ctx, &data)...)
	if resp.Diagnostics.HasError() {
		return
	}

	// Initialize global logging configuration based on environment variables only
	log.InitializeGlobalConfig()

	client, configDiags := p.initPlatformClient(&data)
	resp.Diagnostics.Append(configDiags...)
	if resp.Diagnostics.HasError() {
		return
	}

	// Initialize backend version dynamically from API
	backendVersion, err := p.initBackendVersion(ctx, client)
	if err != nil {
		resp.Diagnostics.AddError(
			"Backend Version Initialization Failed",
			fmt.Sprintf("Failed to initialize backend version: %s", err.Error()),
		)
		return
	}
	p.backendVersion = backendVersion

	providerData := config.FrameworkProviderData{
		Client:         client,
		BackendVersion: backendVersion,
	}

	resp.DataSourceData = &providerData
	resp.ResourceData = &providerData
}

func (p *FrameworkProvider) initPlatformClient(data *FrameworkProviderModel) (*client.PlatformClient, diag.Diagnostics) {
	var diags diag.Diagnostics

	// Setup URL
	url := data.URL.ValueString()
	if url == EMPTY_STRING {
		url = os.Getenv(EnvVarsNamePrefix + "_URL")
	}
	if url == EMPTY_STRING {
		diags.AddError(
			"Missing Configuration",
			"The provider requires a URL to be configured. "+
				"Set the 'url' attribute in the provider configuration or use the "+
				EnvVarsNamePrefix+"_URL environment variable.",
		)
		return nil, diags
	}

	// Setup TOKEN
	token := data.Token.ValueString()
	if token == EMPTY_STRING {
		token = os.Getenv(EnvVarsNamePrefix + "_TOKEN")
	}
	if token == EMPTY_STRING {
		diags.AddError(
			"Missing Configuration",
			"The provider requires a TOKEN to be configured. "+
				"Set the 'token' attribute in the provider configuration or use the "+
				EnvVarsNamePrefix+"_TOKEN environment variable.",
		)
		return nil, diags
	}

	platformClient := client.NewPlatformClient(url, token)

	return platformClient, diags
}

// initBackendVersion initializes the backend version by calling the API
func (p *FrameworkProvider) initBackendVersion(ctx context.Context, client *client.PlatformClient) (*version.BackendVersion, error) {
	requestContext := common.NewRequestContext(ctx).
		WithResourceType("BackendVersion").
		WithOperation(common.Read).
		WithAPIVersion(common.V1).
		// Adding a dummy backend version to avoid nil pointer dereference
		WithBackendVersion(version.NewBackendVersion("0.0.0"))

	log.LogInfo(requestContext, "Initializing backend version from API", nil)

	// Create backend version provider
	versionProvider := backend_version.NewBackendVersionProvider(client)

	// Fetch backend version with fallback
	backendVersion, err := versionProvider.FetchBackendVersionWithFallback(requestContext, DEFAULT_FALLBACK_VERSION)
	if err != nil {
		log.LogError(requestContext, "Failed to initialize backend version", map[string]any{"error": err.Error()})
		return nil, err
	}

	log.LogInfo(requestContext, "Backend version initialized", map[string]any{
		"version": backendVersion.Version,
		"major":   backendVersion.Major,
		"minor":   backendVersion.Minor,
		"patch":   backendVersion.Patch,
	})

	return backendVersion, nil
}

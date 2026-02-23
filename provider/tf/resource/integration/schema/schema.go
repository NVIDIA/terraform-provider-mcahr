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

package schema

import (
	"terraform/terraform-provider/provider/common/attribute"
	deprecation "terraform/terraform-provider/provider/tf/core/plan/modifiers/deprecation"
	nulls "terraform/terraform-provider/provider/tf/core/plan/modifiers/null"
	coreschema "terraform/terraform-provider/provider/tf/core/schema"
	schemabuilder "terraform/terraform-provider/provider/tf/core/schema"
	defaults "terraform/terraform-provider/provider/tf/resource/integration/plan/modifier/default"

	"github.com/hashicorp/terraform-plugin-framework-validators/int64validator"
	"github.com/hashicorp/terraform-plugin-framework/path"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/booldefault"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/planmodifier"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/stringplanmodifier"
	"github.com/hashicorp/terraform-plugin-framework/schema/validator"
	"github.com/hashicorp/terraform-plugin-framework/types"
)

type IntegrationSchema struct{}

var _ coreschema.ResourceSchema = &IntegrationSchema{}

// GetSchema returns the schema definition for the Integration resource.
//
// This schema defines all the attributes available for configuring a integration,
// including required fields (name), optional configuration fields,
// and computed fields.
func (s *IntegrationSchema) GetSchema() schema.Schema {

	builder := schemabuilder.NewSchemaBuilder()

	builder.AddMarkdownDescription("Integration resource for configuring integrations with external services")

	// Global Shared schema attributes
	addGlobalSharedSchema(builder)

	// Integration-specific schema attributes
	addCustomSharedSchema(builder)

	addAlertmanagerSchema(builder)
	addDatadogSchema(builder)
	addAzureActiveDirectorySchema(builder)
	addOktaSchema(builder)
	addGoogleCloudIdentitySchema(builder)
	addBcmSchema(builder)
	addBcmConnectivitySchema(builder)
	addElasticSchema(builder)
	addFluentbitElasticSchema(builder)
	addNvaultSchema(builder)

	return builder.Build()
}

func addGlobalSharedSchema(builder *schemabuilder.SchemaBuilder) {

	builder.AddAttribute("name", schema.StringAttribute{
		MarkdownDescription: "The name of the integration",
		Required:            true,
		PlanModifiers:       []planmodifier.String{stringplanmodifier.RequiresReplace()},
	})

	builder.AddAttribute("service_name", schema.StringAttribute{
		MarkdownDescription: "The integration type",
		Required:            true,
		PlanModifiers:       []planmodifier.String{stringplanmodifier.RequiresReplace()},
	})

	builder.AddAttribute("serial_number", schema.StringAttribute{
		MarkdownDescription: "The serial number of the integration",
		Required:            true,
	})

	builder.AddAttribute("enabled", schema.BoolAttribute{
		MarkdownDescription: "Whether the integration is enabled",
		Optional:            true,
		Computed:            true,
		Default:             booldefault.StaticBool(false),
	})

	builder.AddAttribute("permissions_user", schema.StringAttribute{
		MarkdownDescription: "The permissions user of the integration",
		Optional:            true,
		Computed:            true,
		// The default user may be dynamic on the backend side
		PlanModifiers: []planmodifier.String{stringplanmodifier.UseStateForUnknown()},
	})
}

func addCustomSharedSchema(builder *schemabuilder.SchemaBuilder) {

	builder.AddAttribute("api_url", schema.StringAttribute{
		MarkdownDescription: "The API URL of the integration",
		Optional:            true,
		Computed:            true,
		PlanModifiers: []planmodifier.String{
			nulls.NullStringIfUnknownModifier(),
			defaults.EmptyStringModifier(),
		},
	})

	builder.AddAttribute("api_key", schema.StringAttribute{
		MarkdownDescription: "The API key of the integration",
		Optional:            true,
		Computed:            true,
		PlanModifiers: []planmodifier.String{
			nulls.NullStringIfUnknownModifier(),
			defaults.EmptyStringModifier(),
		},
	})

	builder.AddAttribute("idp_name", schema.StringAttribute{
		MarkdownDescription: "The IDP name of the integration",
		Optional:            true,
		Computed:            true,
		PlanModifiers: []planmodifier.String{
			nulls.NullStringIfUnknownModifier(),
			defaults.EmptyStringModifier(),
		},
	})

	builder.AddAttribute("cache_ttl", schema.Int64Attribute{
		MarkdownDescription: "**Deprecated** Use `cache_ttl_ms` instead.",
		Optional:            true,
		Computed:            true,
		DeprecationMessage:  "use `cache_ttl_ms` instead.",
		PlanModifiers: []planmodifier.Int64{
			nulls.NullInt64IfUnknownModifier(),
			defaults.EmptyInt64Modifier(),
		},
		Validators: []validator.Int64{
			int64validator.ConflictsWith(path.MatchRoot("cache_ttl_ms")),
		},
	})

	builder.AddAttribute("cache_ttl_ms", schema.Int64Attribute{
		MarkdownDescription: "The cache TTL of the integration in milliseconds",
		Optional:            true,
		Computed:            true,
		PlanModifiers: []planmodifier.Int64{
			// Keep this order
			nulls.NullInt64IfUnknownModifier(),
			defaults.EmptyInt64Modifier(),
			deprecation.MaybeGetFromDeprecatedInt64Modifier("cache_ttl"),
		},
	})

	builder.AddAttribute("api_rate_limit", schema.Int64Attribute{
		MarkdownDescription: "The API rate limit of the integration",
		Optional:            true,
		Computed:            true,
		PlanModifiers: []planmodifier.Int64{
			nulls.NullInt64IfUnknownModifier(),
			defaults.EmptyInt64Modifier(),
		},
	})
}

func addAlertmanagerSchema(builder *schemabuilder.SchemaBuilder) {

	builder.AddAttribute("external_url", schema.StringAttribute{
		MarkdownDescription: "The external URL of the integration",
		Optional:            true,
		Computed:            true,
		PlanModifiers: []planmodifier.String{
			nulls.NullStringIfUnknownModifier(),
			defaults.EmptyStringModifier(),
		},
	})

	builder.AddAttribute("payload_paths", schema.ListAttribute{
		MarkdownDescription: "The payload paths of the integration",
		Optional:            true,
		Computed:            true,
		ElementType:         types.StringType,
		PlanModifiers: []planmodifier.List{
			nulls.NullListIfUnknownModifier(),
			defaults.EmptyStringListModifier(),
		},
	})
}

func addDatadogSchema(builder *schemabuilder.SchemaBuilder) {
	// NOTE: api_key and api_url are defined in addCustomSharedSchema

	builder.AddAttribute("site_url", schema.StringAttribute{
		MarkdownDescription: "The site URL of the integration",
		Optional:            true,
		Computed:            true,
		PlanModifiers: []planmodifier.String{
			nulls.NullStringIfUnknownModifier(),
			defaults.EmptyStringModifier(),
		},
	})

	builder.AddAttribute("app_key", schema.StringAttribute{
		MarkdownDescription: "The app key of the integration",
		Optional:            true,
		Computed:            true,
		PlanModifiers: []planmodifier.String{
			nulls.NullStringIfUnknownModifier(),
			defaults.EmptyStringModifier(),
		},
	})

	builder.AddAttribute("webhook_name", schema.StringAttribute{
		MarkdownDescription: "The webhook name of the integration",
		Optional:            true,
		Computed:            true,
		PlanModifiers: []planmodifier.String{
			nulls.NullStringIfUnknownModifier(),
			defaults.EmptyStringModifier(),
		},
	})
}

func addAzureActiveDirectorySchema(builder *schemabuilder.SchemaBuilder) {
	// NOTE: idp_name, cache_ttl_ms, api_rate_limit are defined in addCustomSharedSchema

	builder.AddAttribute("tenant_id", schema.StringAttribute{
		MarkdownDescription: "The tenant ID of the integration",
		Optional:            true,
		Computed:            true,
		PlanModifiers: []planmodifier.String{
			nulls.NullStringIfUnknownModifier(),
			defaults.EmptyStringModifier(),
		},
	})

	builder.AddAttribute("client_id", schema.StringAttribute{
		MarkdownDescription: "The client ID of the integration",
		Optional:            true,
		Computed:            true,
		PlanModifiers: []planmodifier.String{
			nulls.NullStringIfUnknownModifier(),
			defaults.EmptyStringModifier(),
		},
	})

	builder.AddAttribute("client_secret", schema.StringAttribute{
		MarkdownDescription: "The client secret of the integration",
		Optional:            true,
		Computed:            true,
		PlanModifiers: []planmodifier.String{
			nulls.NullStringIfUnknownModifier(),
			defaults.EmptyStringModifier(),
		},
	})
}

func addOktaSchema(builder *schemabuilder.SchemaBuilder) {
	// NOTE: api_key, api_url, idp_name, cache_ttl_ms, api_rate_limit are defined in addCustomSharedSchema
}

func addGoogleCloudIdentitySchema(builder *schemabuilder.SchemaBuilder) {
	// NOTE: idp_name, cache_ttl_ms, api_rate_limit are defined in addCustomSharedSchema

	builder.AddAttribute("subject", schema.StringAttribute{
		MarkdownDescription: "The subject of the integration",
		Optional:            true,
		Computed:            true,
		PlanModifiers: []planmodifier.String{
			nulls.NullStringIfUnknownModifier(),
			defaults.EmptyStringModifier(),
		},
	})

	builder.AddAttribute("credentials", schema.StringAttribute{
		MarkdownDescription: "The credentials of the integration",
		Optional:            true,
		Computed:            true,
		PlanModifiers: []planmodifier.String{
			nulls.NullStringIfUnknownModifier(),
			defaults.EmptyStringModifier(),
		},
	})

}

func addBcmSchema(builder *schemabuilder.SchemaBuilder) {
	// NOTE: idp_name, cache_ttl_ms, api_rate_limit are defined in addCustomSharedSchema
}

func addBcmConnectivitySchema(builder *schemabuilder.SchemaBuilder) {
	// NOTE: api_key is defined in shared schema

	builder.AddAttribute("api_certificate", schema.StringAttribute{
		MarkdownDescription: "The API certificate of the integration",
		Optional:            true,
		Computed:            true,
		PlanModifiers: []planmodifier.String{
			nulls.NullStringIfUnknownModifier(),
			defaults.EmptyStringModifier(),
		},
	})

}

func addElasticSchema(builder *schemabuilder.SchemaBuilder) {
	// NOTE: api_key and api_url are defined in shared schema
}

func addFluentbitElasticSchema(builder *schemabuilder.SchemaBuilder) {
	// NOTE: api_url are defined in shared schema
}

func addNvaultSchema(builder *schemabuilder.SchemaBuilder) {
	// !!! IMPORTANT !!!
	// Don't add defaults here. Use the adapter to set the defaults for custom attributes.

	builder.AddAttribute("address", schema.StringAttribute{
		MarkdownDescription: "The address of the integration",
		Optional:            true,
		Computed:            true,
		PlanModifiers: []planmodifier.String{
			nulls.NullStringIfUnknownModifier(),
			defaults.EmptyStringModifier(),
		},
	})

	builder.AddAttribute("namespace", schema.StringAttribute{
		MarkdownDescription: "The namespace of the integration",
		Optional:            true,
		Computed:            true,
		PlanModifiers: []planmodifier.String{
			nulls.NullStringIfUnknownModifier(),
			defaults.EmptyStringModifier(),
		},
	})

	builder.AddAttribute("role_name", schema.StringAttribute{
		MarkdownDescription: "The role name of the integration",
		Optional:            true,
		Computed:            true,
		PlanModifiers: []planmodifier.String{
			nulls.NullStringIfUnknownModifier(),
			defaults.EmptyStringModifier(),
		},
	})

	builder.AddAttribute("jwt_auth_path", schema.StringAttribute{
		MarkdownDescription: "The JWT auth path of the integration",
		Optional:            true,
		Computed:            true,
		PlanModifiers: []planmodifier.String{
			nulls.NullStringIfUnknownModifier(),
			defaults.EmptyStringModifier(),
		},
	})

}

func (s *IntegrationSchema) GetCompatibilityOptions() map[string]attribute.CompatibilityOptions {
	return map[string]attribute.CompatibilityOptions{
		"external_url": {
			MinVersion: "release-17.0.0",
		},
		"cache_ttl_ms": {
			MinVersion: "release-18.0.0",
		},
		"tenant_id": {
			MinVersion: "release-18.0.0",
		},
		"client_id": {
			MinVersion: "release-18.0.0",
		},
		"client_secret": {
			MinVersion: "release-18.0.0",
		},
		"subject": {
			MinVersion: "release-18.0.0",
		},
		"credentials": {
			MinVersion: "release-18.0.0",
		},
		"site_url": {
			MinVersion: "release-19.0.0",
		},
		"idp_name": {
			MinVersion: "release-22.0.0",
		},
		"api_certificate": {
			MinVersion: "release-28.1.0",
		},
		"payload_paths": {
			MinVersion: "release-28.4.0",
		},
		"address": {
			MinVersion: "release-29.0.0",
		},
		"namespace": {
			MinVersion: "release-29.0.0",
		},
		"role_name": {
			MinVersion: "release-29.0.0",
		},
		"jwt_auth_path": {
			MinVersion: "release-29.0.0",
		},
	}
}

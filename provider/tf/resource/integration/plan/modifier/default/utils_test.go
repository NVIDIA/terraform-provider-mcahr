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

package defaults

import (
	"testing"

	"github.com/hashicorp/terraform-plugin-framework/path"
	"github.com/hashicorp/terraform-plugin-framework/types"
	"github.com/stretchr/testify/assert"
)

func TestIsServiceNameCompatible(t *testing.T) {
	t.Parallel()

	tests := []struct {
		name          string
		serviceName   string
		attributeName string
		expected      bool
	}{
		// Alertmanager tests
		{
			name:          "Alertmanager - external_url compatible",
			serviceName:   "alertmanager",
			attributeName: "external_url",
			expected:      true,
		},
		{
			name:          "Alertmanager - payload_paths compatible",
			serviceName:   "alertmanager",
			attributeName: "payload_paths",
			expected:      true,
		},
		{
			name:          "Alertmanager - incompatible attribute",
			serviceName:   "alertmanager",
			attributeName: "api_key",
			expected:      false,
		},

		// Azure Active Directory tests
		{
			name:          "Azure AD - idp_name compatible",
			serviceName:   "azure_active_directory",
			attributeName: "idp_name",
			expected:      true,
		},
		{
			name:          "Azure AD - tenant_id compatible",
			serviceName:   "azure_active_directory",
			attributeName: "tenant_id",
			expected:      true,
		},
		{
			name:          "Azure AD - client_secret compatible",
			serviceName:   "azure_active_directory",
			attributeName: "client_secret",
			expected:      true,
		},
		{
			name:          "Azure AD - incompatible attribute",
			serviceName:   "azure_active_directory",
			attributeName: "external_url",
			expected:      false,
		},

		// BCM tests
		{
			name:          "BCM - idp_name compatible",
			serviceName:   "bcm",
			attributeName: "idp_name",
			expected:      true,
		},
		{
			name:          "BCM - cache_ttl_ms compatible",
			serviceName:   "bcm",
			attributeName: "cache_ttl_ms",
			expected:      true,
		},
		{
			name:          "BCM - api_rate_limit compatible",
			serviceName:   "bcm",
			attributeName: "api_rate_limit",
			expected:      true,
		},
		{
			name:          "BCM - incompatible attribute",
			serviceName:   "bcm",
			attributeName: "api_certificate",
			expected:      false,
		},

		// BCM Connectivity tests
		{
			name:          "BCM Connectivity - api_key compatible",
			serviceName:   "bcm_connectivity",
			attributeName: "api_key",
			expected:      true,
		},
		{
			name:          "BCM Connectivity - api_certificate compatible",
			serviceName:   "bcm_connectivity",
			attributeName: "api_certificate",
			expected:      true,
		},
		{
			name:          "BCM Connectivity - incompatible attribute",
			serviceName:   "bcm_connectivity",
			attributeName: "idp_name",
			expected:      false,
		},

		// Datadog tests
		{
			name:          "Datadog - api_key compatible",
			serviceName:   "datadog",
			attributeName: "api_key",
			expected:      true,
		},
		{
			name:          "Datadog - webhook_name compatible",
			serviceName:   "datadog",
			attributeName: "webhook_name",
			expected:      true,
		},
		{
			name:          "Datadog - incompatible attribute",
			serviceName:   "datadog",
			attributeName: "tenant_id",
			expected:      false,
		},

		// Fluentbit Elastic tests
		{
			name:          "Fluentbit Elastic - api_url compatible",
			serviceName:   "fluentbit_elastic",
			attributeName: "api_url",
			expected:      true,
		},
		{
			name:          "Fluentbit Elastic - incompatible attribute",
			serviceName:   "fluentbit_elastic",
			attributeName: "api_key",
			expected:      false,
		},

		// Google Cloud Identity tests
		{
			name:          "Google Cloud Identity - subject compatible",
			serviceName:   "google_cloud_identity",
			attributeName: "subject",
			expected:      true,
		},
		{
			name:          "Google Cloud Identity - credentials compatible",
			serviceName:   "google_cloud_identity",
			attributeName: "credentials",
			expected:      true,
		},
		{
			name:          "Google Cloud Identity - incompatible attribute",
			serviceName:   "google_cloud_identity",
			attributeName: "client_secret",
			expected:      false,
		},

		// NVault tests
		{
			name:          "NVault - address compatible",
			serviceName:   "nvault",
			attributeName: "address",
			expected:      true,
		},
		{
			name:          "NVault - jwt_auth_path compatible",
			serviceName:   "nvault",
			attributeName: "jwt_auth_path",
			expected:      true,
		},
		{
			name:          "NVault - incompatible attribute",
			serviceName:   "nvault",
			attributeName: "api_key",
			expected:      false,
		},

		// Okta tests
		{
			name:          "Okta - api_key compatible (TF field name)",
			serviceName:   "okta",
			attributeName: "api_key",
			expected:      true,
		},
		{
			name:          "Okta - api_url compatible (TF field name)",
			serviceName:   "okta",
			attributeName: "api_url",
			expected:      true,
		},
		{
			name:          "Okta - idp_name compatible",
			serviceName:   "okta",
			attributeName: "idp_name",
			expected:      true,
		},
		{
			name:          "Okta - incompatible attribute",
			serviceName:   "okta",
			attributeName: "api_token", // This is the data field name, not TF field name
			expected:      false,
		},

		// Unsupported service name tests
		{
			name:          "Unsupported service name",
			serviceName:   "unsupported_service",
			attributeName: "any_attribute",
			expected:      false,
		},
		{
			name:          "Empty service name",
			serviceName:   "",
			attributeName: "any_attribute",
			expected:      false,
		},
		{
			name:          "Nil-like service name",
			serviceName:   "null",
			attributeName: "any_attribute",
			expected:      false,
		},

		// Edge cases
		{
			name:          "Empty attribute name",
			serviceName:   "datadog",
			attributeName: "",
			expected:      false,
		},
		{
			name:          "Case sensitivity test - uppercase service",
			serviceName:   "DATADOG",
			attributeName: "api_key",
			expected:      false, // Service names should be lowercase
		},
		{
			name:          "Case sensitivity test - uppercase attribute",
			serviceName:   "datadog",
			attributeName: "API_KEY",
			expected:      false, // Attribute names should be lowercase
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := IsServiceNameCompatible(tt.serviceName, tt.attributeName)

			assert.Equal(t, tt.expected, result)
		})
	}
}

func TestGetAttributeName(t *testing.T) {
	t.Parallel()

	tests := []struct {
		name     string
		path     path.Path
		expected string
	}{
		{
			name:     "Simple attribute path",
			path:     path.Root("api_key"),
			expected: "api_key",
		},
		{
			name:     "Nested attribute path",
			path:     path.Root("integration").AtName("external_url"),
			expected: "external_url",
		},
		{
			name:     "Deeply nested attribute path",
			path:     path.Root("config").AtName("integration").AtName("datadog").AtName("webhook_name"),
			expected: "webhook_name",
		},
		{
			name:     "Path with list element",
			path:     path.Root("integrations").AtListIndex(0).AtName("api_url"),
			expected: "api_url",
		},
		{
			name:     "Path with set element",
			path:     path.Root("payload_paths").AtSetValue(types.StringValue("test")).AtName("value"),
			expected: "value",
		},
		{
			name:     "Path with map element",
			path:     path.Root("config").AtMapKey("integration").AtName("enabled"),
			expected: "enabled",
		},
		{
			name:     "Empty path",
			path:     path.Empty(),
			expected: "",
		},
		{
			name:     "Path ending with list index (no attribute)",
			path:     path.Root("integrations").AtListIndex(0),
			expected: "",
		},
		{
			name:     "Path ending with set value (no attribute)",
			path:     path.Root("payload_paths").AtSetValue(types.StringValue("test")),
			expected: "",
		},
		{
			name:     "Path ending with map key (no attribute)",
			path:     path.Root("config").AtMapKey("integration"),
			expected: "",
		},
		{
			name:     "Complex path with multiple step types",
			path:     path.Root("data").AtListIndex(1).AtName("integrations").AtMapKey("primary").AtName("cache_ttl_ms"),
			expected: "cache_ttl_ms",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := GetAttributeName(tt.path)

			assert.Equal(t, tt.expected, result)
		})
	}
}

func TestIsServiceNameCompatible_AllSupportedServices(t *testing.T) {
	t.Parallel()

	// Test that all known services are handled correctly
	supportedServices := []string{
		"alertmanager",
		"azure_active_directory",
		"bcm",
		"bcm_connectivity",
		"datadog",
		"fluentbit_elastic",
		"google_cloud_identity",
		"nvault",
		"okta",
	}

	for _, serviceName := range supportedServices {
		t.Run("Service "+serviceName+" should have at least one compatible attribute", func(t *testing.T) {
			// Each service should have at least one compatible attribute
			// Let's test with a common attribute that should exist
			hasCompatibleAttribute := false

			// Test common attributes that might exist
			testAttributes := []string{
				"api_key", "api_url", "idp_name", "cache_ttl_ms", "api_rate_limit",
				"external_url", "payload_paths", "tenant_id", "client_id", "client_secret",
				"site_url", "app_key", "webhook_name", "api_certificate", "subject",
				"credentials", "address", "namespace", "role_name", "jwt_auth_path",
			}

			for _, attr := range testAttributes {
				if IsServiceNameCompatible(serviceName, attr) {
					hasCompatibleAttribute = true
					break
				}
			}

			assert.True(t, hasCompatibleAttribute, "Service %s should have at least one compatible attribute", serviceName)
		})
	}
}

func TestGetAttributeName_PathStepTypes(t *testing.T) {
	t.Parallel()

	// Test various path step combinations to ensure robustness
	t.Run("Multiple attribute steps", func(t *testing.T) {
		// This tests the behavior when there are multiple attribute steps
		// The function should return the last one
		path := path.Root("first").AtName("second").AtName("third")
		result := GetAttributeName(path)
		assert.Equal(t, "third", result)
	})

	t.Run("Attribute after list index", func(t *testing.T) {
		path := path.Root("items").AtListIndex(5).AtName("attribute")
		result := GetAttributeName(path)
		assert.Equal(t, "attribute", result)
	})

	t.Run("Attribute after map key", func(t *testing.T) {
		path := path.Root("config").AtMapKey("database").AtName("host")
		result := GetAttributeName(path)
		assert.Equal(t, "host", result)
	})
}

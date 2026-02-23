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

package customattribute

import (
	"encoding/json"
	"terraform/terraform-provider/provider/common"
	commonstruct "terraform/terraform-provider/provider/common/struct"
)

var _ common.JsonConfigurable = &ConfigurationJson{}
var _ json.Marshaler = &ConfigurationJson{}
var _ json.Unmarshaler = &ConfigurationJson{}

// ConfigurationJson represents the dashboard configuration with custom JSON processing
type ConfigurationJson struct {
	// Embedded config to be used in marshal/unmarshal functions
	Config common.JsonConfig `json:"-" skip:"true"`

	Groups        []GroupJson `json:"groups"`
	Values        []ValueJson `json:"values"`
	Identifiers   []string    `json:"identifiers"`
	OtherTags     []string    `json:"other_tags"`
	ResourceQuery string      `json:"resource_query"`
}

// SetConfig sets the JSON configuration
func (c *ConfigurationJson) SetConfig(config common.JsonConfig) {
	c.Config = config
}

// GetConfig returns the JSON configuration
func (c *ConfigurationJson) GetConfig() common.JsonConfig {
	return c.Config
}

// MarshalJSON customizes the JSON marshaling
func (c *ConfigurationJson) MarshalJSON() ([]byte, error) {
	// Apply version-specific processing and defaults
	options := map[string]any{"backend_version": c.Config.BackendVersion}
	result := commonstruct.ApplyCustomStructTags(*c, options)

	return json.Marshal(result)
}

// UnmarshalJSON customizes the JSON unmarshaling with version-aware processing
func (c *ConfigurationJson) UnmarshalJSON(data []byte) error {
	// Unmarshal directly into the struct first
	type Alias ConfigurationJson
	aux := &Alias{
		Config: c.Config, // Preserve config
		// Set default values
		Groups:        []GroupJson{},
		Values:        []ValueJson{},
		Identifiers:   []string{},
		OtherTags:     []string{},
		ResourceQuery: "",
	}

	// Direct unmarshal - let Go handle the type conversion
	if err := json.Unmarshal(data, aux); err != nil {
		return err
	}

	// Apply custom struct tags processing if needed
	if c.Config.BackendVersion != nil {
		options := map[string]any{"backend_version": c.Config.BackendVersion}
		result := commonstruct.ApplyCustomStructTags(*aux, options)

		// Re-marshal and unmarshal to apply any filtering
		processedData, _ := json.Marshal(result)
		json.Unmarshal(processedData, aux)
	}

	// Copy the processed data to the original struct
	*c = ConfigurationJson(*aux)

	return nil
}

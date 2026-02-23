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

var _ common.JsonConfigurable = &GroupJson{}
var _ json.Marshaler = &GroupJson{}
var _ json.Unmarshaler = &GroupJson{}

// GroupJson represents a dashboard group configuration with custom JSON processing
type GroupJson struct {
	// Embedded config to be used in marshal/unmarshal functions
	Config common.JsonConfig `json:"-" skip:"true"`

	Name string   `json:"name"`
	Tags []string `json:"tags"`
}

// SetConfig sets the JSON configuration
func (g *GroupJson) SetConfig(config common.JsonConfig) {
	g.Config = config
}

// GetConfig returns the JSON configuration
func (g *GroupJson) GetConfig() common.JsonConfig {
	return g.Config
}

// MarshalJSON customizes the JSON marshaling
func (g *GroupJson) MarshalJSON() ([]byte, error) {
	// Apply version-specific processing and defaults
	options := map[string]any{"backend_version": g.Config.BackendVersion}
	result := commonstruct.ApplyCustomStructTags(*g, options)

	return json.Marshal(result)
}

// UnmarshalJSON customizes the JSON unmarshaling with version-aware processing
func (g *GroupJson) UnmarshalJSON(data []byte) error {
	// Unmarshal directly into the struct first
	type Alias GroupJson
	aux := &Alias{
		Config: g.Config, // Preserve config
		// Set default values
		Name: "",
		Tags: []string{},
	}

	// Direct unmarshal - let Go handle the type conversion
	if err := json.Unmarshal(data, aux); err != nil {
		return err
	}

	// Apply custom struct tags processing if needed
	if g.Config.BackendVersion != nil {
		options := map[string]any{"backend_version": g.Config.BackendVersion}
		result := commonstruct.ApplyCustomStructTags(*aux, options)

		// Re-marshal and unmarshal to apply any filtering
		processedData, _ := json.Marshal(result)
		json.Unmarshal(processedData, aux)
	}

	// Copy the processed data to the original struct
	*g = GroupJson(*aux)

	return nil
}

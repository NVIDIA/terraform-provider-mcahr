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

var _ common.JsonConfigurable = &ValueJson{}
var _ json.Marshaler = &ValueJson{}
var _ json.Unmarshaler = &ValueJson{}

// ValueJson represents a dashboard value configuration with custom JSON processing
type ValueJson struct {
	// Embedded config to be used in marshal/unmarshal functions
	Config common.JsonConfig `json:"-" skip:"true"`

	Color  string   `json:"color"`
	Values []string `json:"values"`
}

// SetConfig sets the JSON configuration
func (v *ValueJson) SetConfig(config common.JsonConfig) {
	v.Config = config
}

// GetConfig returns the JSON configuration
func (v *ValueJson) GetConfig() common.JsonConfig {
	return v.Config
}

// MarshalJSON customizes the JSON marshaling
func (v *ValueJson) MarshalJSON() ([]byte, error) {
	// Apply version-specific processing and defaults
	options := map[string]any{"backend_version": v.Config.BackendVersion}
	result := commonstruct.ApplyCustomStructTags(*v, options)

	return json.Marshal(result)
}

// UnmarshalJSON customizes the JSON unmarshaling with version-aware processing
func (v *ValueJson) UnmarshalJSON(data []byte) error {
	// Unmarshal directly into the struct first
	type Alias ValueJson
	aux := &Alias{
		Config: v.Config, // Preserve config
		// Set default values
		Color:  "",
		Values: []string{},
	}

	// Direct unmarshal - let Go handle the type conversion
	if err := json.Unmarshal(data, aux); err != nil {
		return err
	}

	// Apply custom struct tags processing if needed
	if v.Config.BackendVersion != nil {
		options := map[string]any{"backend_version": v.Config.BackendVersion}
		result := commonstruct.ApplyCustomStructTags(*aux, options)

		// Re-marshal and unmarshal to apply any filtering
		processedData, _ := json.Marshal(result)
		json.Unmarshal(processedData, aux)
	}

	// Copy the processed data to the original struct
	*v = ValueJson(*aux)

	return nil
}

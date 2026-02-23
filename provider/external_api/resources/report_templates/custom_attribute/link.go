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

var _ common.JsonConfigurable = &LinkJson{}
var _ json.Marshaler = &LinkJson{}
var _ json.Unmarshaler = &LinkJson{}

// LinkJson represents a link with custom attribute processing
type LinkJson struct {
	// Embedded config to be used in marshal/unmarshal functions
	Config common.JsonConfig `json:"-" skip:"true"`

	Label              string `json:"label"`
	ReportTemplateName string `json:"report_template_name"`
}

func (l *LinkJson) SetConfig(config common.JsonConfig) {
	l.Config = config
}

func (l *LinkJson) GetConfig() common.JsonConfig {
	return l.Config
}

func (l *LinkJson) MarshalJSON() ([]byte, error) {
	// Apply version-specific processing and defaults
	options := map[string]any{"backend_version": l.Config.BackendVersion}
	result := commonstruct.ApplyCustomStructTags(*l, options)

	return json.Marshal(result)
}

func (l *LinkJson) UnmarshalJSON(data []byte) error {
	type Alias LinkJson // Prevent recursion
	aux := &Alias{
		Config: common.JsonConfig{
			BackendVersion: l.Config.BackendVersion,
		},
		Label:              "",
		ReportTemplateName: "",
	} // defaults

	// Direct unmarshal - let Go handle the type conversion
	if err := json.Unmarshal(data, aux); err != nil {
		return err
	}

	// Apply custom struct tags processing if needed
	if l.Config.BackendVersion != nil {
		options := map[string]any{"backend_version": l.Config.BackendVersion}
		result := commonstruct.ApplyCustomStructTags(*aux, options)

		// Re-marshal and unmarshal to apply any filtering
		processedData, _ := json.Marshal(result)
		json.Unmarshal(processedData, aux)
	}

	// Copy the processed data to the original struct
	*l = LinkJson(*aux)

	return nil
}

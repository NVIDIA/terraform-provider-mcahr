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

var (
	DefaultExternalParamName        = ""
	DefaultExternalParamValue       = ""
	DefaultExternalParamSource      = ""
	DefaultExternalParamJsonPath    = ""
	DefaultExternalParamExport      = false
	DefaultExternalParamType        = "EXTERNAL"
	DefaultExternalParamDescription = ""

	ValidExternalParamSources = []string{"alertmanager"}
)

type ExternalParamJson struct {
	// Embedded config to be used in marshal/unmarshal functions
	Config common.JsonConfig `json:"-" skip:"true"`

	Name        string `json:"name"`
	Value       string `json:"value"`
	Source      string `json:"source"`
	JsonPath    string `json:"json_path"`
	Export      bool   `json:"export"`
	ParamType   string `json:"param_type"`
	Description string `json:"description" min_version:"release-28.4.0"`
}

var _ common.JsonConfigurable = &ExternalParamJson{}
var _ json.Marshaler = &ExternalParamJson{}
var _ json.Unmarshaler = &ExternalParamJson{}

func (p *ExternalParamJson) SetConfig(config common.JsonConfig) {
	p.Config = config
}

func (p *ExternalParamJson) GetConfig() common.JsonConfig {
	return p.Config
}

func (p *ExternalParamJson) MarshalJSON() ([]byte, error) {

	options := map[string]any{"backend_version": p.Config.BackendVersion}
	result := commonstruct.ApplyCustomStructTags(*p, options)

	return json.Marshal(result)
}

func (p *ExternalParamJson) UnmarshalJSON(b []byte) error {
	type Alias ExternalParamJson // Prevent recursion
	aux := &Alias{
		Config: common.JsonConfig{
			BackendVersion: p.Config.BackendVersion,
		},
		Name:        DefaultExternalParamName,
		Value:       DefaultExternalParamValue,
		Source:      DefaultExternalParamSource,
		JsonPath:    DefaultExternalParamJsonPath,
		Export:      DefaultExternalParamExport,
		ParamType:   DefaultExternalParamType,
		Description: DefaultExternalParamDescription,
	} // defaults

	// Unmarshal and replace the default values with the ones from the JSON
	if err := json.Unmarshal(b, aux); err != nil {
		return err
	}

	// Apply custom struct tags to the unmarshalled data
	options := map[string]any{"backend_version": p.Config.BackendVersion}
	result := commonstruct.ApplyCustomStructTags(*aux, options)

	// Set in the final struct the values from the processed data
	p.Name = stringOrKeep(result, "name", p.Name)

	p.Value = stringOrKeep(result, "value", p.Value)

	p.Source = stringOrKeep(result, "source", p.Source)

	p.JsonPath = stringOrKeep(result, "json_path", p.JsonPath)

	p.Export = boolOrKeep(result, "export", p.Export)

	p.ParamType = stringOrKeep(result, "param_type", p.ParamType)

	p.Description = stringOrKeep(result, "description", p.Description)

	return nil
}

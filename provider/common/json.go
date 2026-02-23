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

package common

import (
	"encoding/json"
	"terraform/terraform-provider/provider/common/version"
)

type JsonConfig struct {
	BackendVersion *version.BackendVersion `json:"backend_version"`
}

type JsonConfigurable interface {
	SetConfig(config JsonConfig)
	GetConfig() JsonConfig
}

// RemarshalWithConfig unmarshals the encoded data into a JsonConfigurable and then marshals it again
// The config contains metadata needed while marshalling and unmarshalling the data
//
// The main purpose of this function is to process JSON fields in the following way:
// 1. Unmarshal the JSON data into a JsonConfigurable structure
//   - Apply custom struct tags (like min_version, max_version, etc.)
//   - Set default values for the fields that are not present in the JSON
//   - Validation on fields
//
// 2. Marshal the JsonConfigurable structure back into a JSON string
func RemarshalWithConfig[T JsonConfigurable](encodedData string, config JsonConfig) (string, error) {

	var defaultValues T

	defaultValues.SetConfig(config)

	err := json.Unmarshal([]byte(encodedData), &defaultValues)
	if err != nil {
		return "", err
	}

	jsonData, err := json.Marshal(defaultValues)
	if err != nil {
		return "", err
	}

	return string(jsonData), nil
}

// Same as RemarshalWithConfig, but for lists
// The config cannot be set for Unmarshal, so it will only be available for Marshal
func RemarshalListWithConfig[T JsonConfigurable](encodedData string, config JsonConfig) (string, error) {

	var defaultValues []T

	err := json.Unmarshal([]byte(encodedData), &defaultValues)
	if err != nil {
		return "", err
	}

	for i := range defaultValues {
		defaultValues[i].SetConfig(config)
	}

	jsonData, err := json.Marshal(defaultValues)
	if err != nil {
		return "", err
	}
	return string(jsonData), nil
}

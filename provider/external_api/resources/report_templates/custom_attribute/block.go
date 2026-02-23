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

var _ common.JsonConfigurable = &BlockJson{}
var _ json.Marshaler = &BlockJson{}
var _ json.Unmarshaler = &BlockJson{}

// BreakdownTagValue represents breakdown tag values used in breakdown_tags_values
type BreakdownTagValue struct {
	Color  string   `json:"color"`
	Label  string   `json:"label"`
	Values []string `json:"values"`
}

// BreakdownValue represents individual breakdown values used in resources_breakdown
type BreakdownValue struct {
	Value string `json:"value"`
	Count int    `json:"count"`
}

// ResourcesBreakdown represents resources breakdown structure
type ResourcesBreakdown struct {
	GroupByValue    string           `json:"group_by_value"`
	BreakdownValues []BreakdownValue `json:"breakdown_values"`
}

// GroupByTagOrder represents the ordering configuration for group-by tags
type GroupByTagOrder struct {
	Type   string   `json:"type"`   // e.g., "DEFAULT", "CUSTOM", "ALPHABETICAL"
	Values []string `json:"values"` // custom ordering values
}

// BlockJson represents a block with custom attribute processing
type BlockJson struct {
	// Embedded config to be used in marshal/unmarshal functions
	Config common.JsonConfig `json:"-" skip:"true"`

	Title                           string               `json:"title"`
	ResourceQuery                   string               `json:"resource_query"`
	GroupByTag                      string               `json:"group_by_tag,omitempty"`
	GroupByTagOrder                 GroupByTagOrder      `json:"group_by_tag_order,omitempty"`
	BreakdownByTag                  string               `json:"breakdown_by_tag,omitempty"`
	ViewMode                        string               `json:"view_mode,omitempty"`
	IncludeResourcesWithoutGroupTag bool                 `json:"include_resources_without_group_tag,omitempty"`
	IncludeOtherBreakdownTagValues  bool                 `json:"include_other_breakdown_tag_values,omitempty"`
	OtherTagsToExport               []string             `json:"other_tags_to_export,omitempty"`
	BreakdownTagsValues             []BreakdownTagValue  `json:"breakdown_tags_values,omitempty"`
	ResourcesBreakdown              []ResourcesBreakdown `json:"resources_breakdown,omitempty"`
}

func (b *BlockJson) SetConfig(config common.JsonConfig) {
	b.Config = config
}

func (b *BlockJson) GetConfig() common.JsonConfig {
	return b.Config
}

func (b *BlockJson) MarshalJSON() ([]byte, error) {
	// Apply version-specific processing and defaults
	options := map[string]any{"backend_version": b.Config.BackendVersion}
	result := commonstruct.ApplyCustomStructTags(*b, options)

	return json.Marshal(result)
}

func (b *BlockJson) UnmarshalJSON(data []byte) error {
	// Unmarshal directly into the struct first
	type Alias BlockJson
	aux := &Alias{
		Config: b.Config, // Preserve config
		// Set default values
		ViewMode: "COUNT",
		GroupByTagOrder: GroupByTagOrder{
			Type:   "DEFAULT",
			Values: []string{},
		},
		OtherTagsToExport:               []string{},
		BreakdownTagsValues:             []BreakdownTagValue{},
		ResourcesBreakdown:              []ResourcesBreakdown{},
		IncludeResourcesWithoutGroupTag: false,
		IncludeOtherBreakdownTagValues:  false,
	}

	// Direct unmarshal - let Go handle the type conversion
	if err := json.Unmarshal(data, aux); err != nil {
		return err
	}

	// Apply custom struct tags processing if needed
	if b.Config.BackendVersion != nil {
		options := map[string]any{"backend_version": b.Config.BackendVersion}
		result := commonstruct.ApplyCustomStructTags(*aux, options)

		// Re-marshal and unmarshal to apply any filtering
		processedData, _ := json.Marshal(result)
		json.Unmarshal(processedData, aux)
	}

	// Copy the processed data to the original struct
	*b = BlockJson(*aux)

	return nil
}

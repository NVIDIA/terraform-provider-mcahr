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

	"github.com/hashicorp/terraform-plugin-framework/resource/schema"
)

// FieldComparisonBehavior defines how a field should be handled in API response comparisons
type FieldComparisonBehavior int

const (
	// CompareNormally - Compare the field normally (default behavior)
	CompareNormally FieldComparisonBehavior = iota
	// SkipComparison - Skip comparing this field entirely (differences are expected)
	SkipComparison
	// CustomComparison - Use a custom comparison function for this field
	CustomComparison
)

// FieldComparisonRule defines how a specific field should be compared
type FieldComparisonRule struct {
	Behavior FieldComparisonBehavior
	// CustomCompare is called when Behavior is CustomComparison
	// Returns true if values are considered equal, false if they differ
	// planValue and apiValue are the string representations of the field values
	CustomCompare func(fieldName, planValue, apiValue string) bool
	// Reason explains why this field has special comparison behavior (for documentation)
	Reason string
}

type ResourceSchema interface {
	GetSchema() schema.Schema
	GetCompatibilityOptions() map[string]attribute.CompatibilityOptions
	// GetFieldComparisonRules returns rules for how specific fields should be compared in API response warnings.
	// Map key is the field name (tfsdk tag), value defines the comparison behavior.
	GetFieldComparisonRules() map[string]FieldComparisonRule
}

// DefaultFieldComparisonRules provides a default empty implementation
// for resources that don't need special field comparison rules
func DefaultFieldComparisonRules() map[string]FieldComparisonRule {
	return map[string]FieldComparisonRule{}
}

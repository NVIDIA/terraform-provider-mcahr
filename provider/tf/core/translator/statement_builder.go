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

package translator

import (
	"fmt"
	"strings"
	"terraform/terraform-provider/provider/common/attribute"
	"terraform/terraform-provider/provider/common/version"

	"github.com/elliotchance/orderedmap/v3"
)

// StatementBuilder helps build statements with field names as keys
type StatementBuilder struct {
	attrCompatibilityChecker *attribute.CompatibilityChecker

	statementName string
	fields        *orderedmap.OrderedMap[string, any]
}

// NewStatementBuilder creates a new builder for statements
func NewStatementBuilder(statementName string, backendVersion *version.BackendVersion, compatibilityOptions map[string]attribute.CompatibilityOptions) *StatementBuilder {

	attrCompatibilityChecker := attribute.NewCompatibilityChecker(backendVersion, compatibilityOptions)

	return &StatementBuilder{
		attrCompatibilityChecker: attrCompatibilityChecker,
		statementName:            statementName,
		fields:                   orderedmap.NewOrderedMap[string, any](),
	}
}

// SetField adds a field with its value to the statement WITHOUT escaping it.
//
// This is the raw escape hatch. Build renders the value with %v straight into
// the statement, so a string passed here reaches the backend parser verbatim:
// a quote in it terminates the surrounding literal and everything after it is
// parsed as statement syntax.
//
// Use SetStringField for anything that is a string. SetField is only correct
// for values that cannot carry statement syntax:
//
//   - bools and numbers, which is what most callers pass;
//   - JSON produced by encoding/json, which escapes quotes and backslashes
//     itself -- integration params (adapter.TFDataToJSON), runbook
//     params_groups, and the nvault secret external_value all rely on this.
//
// Passing an unescaped operator-supplied string here is an injection bug.
func (b *StatementBuilder) SetField(apiFieldName string, value any, tfFieldName string) *StatementBuilder {

	if b.attrCompatibilityChecker.IsAttributeCompatible(tfFieldName) {
		b.fields.Set(apiFieldName, value)
	}
	return b
}

// SetStringField adds a string field with automatic escaping
func (b *StatementBuilder) SetStringField(apiFieldName string, value string, tfFieldName string) *StatementBuilder {

	if b.attrCompatibilityChecker.IsAttributeCompatible(tfFieldName) {
		b.fields.Set(apiFieldName, EscapeString(value))
	}
	return b
}

// SetStringField adds a string field with automatic escaping
func (b *StatementBuilder) MaybeSetStringField(apiFieldName string, value string, tfFieldName string, isKnown bool) *StatementBuilder {

	if b.attrCompatibilityChecker.IsAttributeCompatible(tfFieldName) && isKnown {
		b.fields.Set(apiFieldName, EscapeString(value))
	}
	return b
}

// SetCommandField adds a command field type string
func (b *StatementBuilder) SetCommandField(apiFieldName string, value string, tfFieldName string) *StatementBuilder {

	if b.attrCompatibilityChecker.IsAttributeCompatible(tfFieldName) {
		b.fields.Set(apiFieldName, EscapeString(value))
	}
	return b
}

// SetArrayField adds an array field with automatic OpLang conversion
func (b *StatementBuilder) SetArrayField(apiFieldName string, value []string, tfFieldName string) *StatementBuilder {

	if b.attrCompatibilityChecker.IsAttributeCompatible(tfFieldName) {
		b.fields.Set(apiFieldName, ArrayToOpLang(value))
	}
	return b
}

// SetBoolField adds a boolean field
func (b *StatementBuilder) SetBoolField(apiFieldName string, value bool, tfFieldName string) *StatementBuilder {

	if b.attrCompatibilityChecker.IsAttributeCompatible(tfFieldName) {
		b.fields.Set(apiFieldName, value)
	}
	return b
}

// Build constructs the final statement string.
//
// Values are rendered verbatim with %v. Escaping is the responsibility of the
// Set*Field method that stored them -- SetStringField, SetCommandField,
// MaybeSetStringField and SetArrayField escape; SetField does not.
func (b *StatementBuilder) Build() string {
	var parts []string

	for pair := b.fields.Front(); pair != nil; pair = pair.Next() {
		parts = append(parts, fmt.Sprintf("%s=%v", pair.Key, pair.Value))
	}

	return fmt.Sprintf("%s(%s)", b.statementName, strings.Join(parts, ", "))
}

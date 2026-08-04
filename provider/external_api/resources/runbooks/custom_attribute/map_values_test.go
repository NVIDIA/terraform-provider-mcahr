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
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestStringOrKeep(t *testing.T) {
	t.Parallel()

	data := map[string]any{
		"good":    "value",
		"number":  42.0,
		"null":    nil,
		"boolean": true,
		"object":  map[string]any{"nested": 1},
		"array":   []any{"a"},
	}

	assert.Equal(t, "value", stringOrKeep(data, "good", "fallback"))
	assert.Equal(t, "fallback", stringOrKeep(data, "absent", "fallback"))

	// Each of these panicked before: the setters asserted .(string) directly.
	for _, key := range []string{"number", "null", "boolean", "object", "array"} {
		assert.Equal(t, "fallback", stringOrKeep(data, key, "fallback"), "key %q", key)
	}
}

func TestBoolOrKeep(t *testing.T) {
	t.Parallel()

	data := map[string]any{
		"good":   true,
		"string": "true",
		"number": 1.0,
		"null":   nil,
	}

	assert.True(t, boolOrKeep(data, "good", false))
	assert.False(t, boolOrKeep(data, "absent", false))

	for _, key := range []string{"string", "number", "null"} {
		assert.False(t, boolOrKeep(data, key, false), "key %q", key)
	}
}

func TestCellJsonAPI_SetFromMap_MistypedBackendDataDoesNotPanic(t *testing.T) {
	t.Parallel()

	// given: every field arrives as the wrong type, as a backend regression or
	// a hostile response could produce
	cell := &CellJsonAPI{}
	input := map[string]interface{}{
		"name":         42.0,
		"content":      []any{"not", "a", "string"},
		"enabled":      "yes",
		"type":         nil,
		"secret_aware": 1.0,
		"description":  map[string]any{},
	}

	// when: this used to panic on the first assertion, taking the provider
	// process down with it
	assert.NotPanics(t, func() { cell.SetFromMap(input) })

	// then: defaults survive
	assert.Equal(t, DefaultCellName, cell.Name)
	assert.Equal(t, DefaultCellEnabled, cell.Enabled)
	assert.Equal(t, DefaultCellSecretAware, cell.SecretAware)
	assert.Equal(t, DefaultCellDescription, cell.Description)
	assert.Equal(t, "", cell.Content)
	assert.Equal(t, "", cell.CellType)
}

func TestCellJsonAPI_SetFromMap_FallsBackToCellType(t *testing.T) {
	t.Parallel()

	// given: "type" is absent, so the older "cell_type" spelling is used
	cell := &CellJsonAPI{}
	cell.SetFromMap(map[string]interface{}{"cell_type": OP_LANG_TYPE})
	assert.Equal(t, OP_LANG_TYPE, cell.CellType)

	// and: a mistyped "type" no longer shadows a usable "cell_type"
	other := &CellJsonAPI{}
	other.SetFromMap(map[string]interface{}{"type": 7.0, "cell_type": OP_LANG_TYPE})
	assert.Equal(t, OP_LANG_TYPE, other.CellType)
}

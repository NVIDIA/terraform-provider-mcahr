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
	"terraform/terraform-provider/provider/common/attribute"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestNewStatementBuilder(t *testing.T) {
	t.Parallel()

	tests := []struct {
		name          string
		statementName string
	}{
		{
			name:          "simple statement name",
			statementName: "define_action",
		},
		{
			name:          "empty statement name",
			statementName: "",
		},
		{
			name:          "statement name with underscores",
			statementName: "update_alarm_config",
		},
		{
			name:          "statement name with special characters",
			statementName: "test-statement_123",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			// given
			statementName := tt.statementName

			// when
			builder := NewStatementBuilder(statementName, nil, map[string]attribute.CompatibilityOptions{})

			// then
			assert.NotNil(t, builder)
			assert.Equal(t, tt.statementName, builder.statementName)
			assert.NotNil(t, builder.fields)
			assert.Equal(t, 0, builder.fields.Len())
		})
	}
}

func TestStatementBuilder_SetField(t *testing.T) {
	t.Parallel()

	tests := []struct {
		name      string
		fieldName string
		value     any
		expected  any
	}{
		{
			name:      "string value",
			fieldName: "name",
			value:     "test_value",
			expected:  "test_value",
		},
		{
			name:      "integer value",
			fieldName: "timeout",
			value:     30,
			expected:  30,
		},
		{
			name:      "boolean value",
			fieldName: "enabled",
			value:     true,
			expected:  true,
		},
		{
			name:      "float value",
			fieldName: "threshold",
			value:     3.14,
			expected:  3.14,
		},
		{
			name:      "nil value",
			fieldName: "optional_field",
			value:     nil,
			expected:  nil,
		},
		{
			name:      "slice value",
			fieldName: "items",
			value:     []string{"item1", "item2"},
			expected:  []string{"item1", "item2"},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			// given
			builder := NewStatementBuilder("test_statement", nil, map[string]attribute.CompatibilityOptions{})

			// when
			result := builder.SetField(tt.fieldName, tt.value, "")

			// then
			// Should return the same builder for chaining
			assert.Equal(t, builder, result)

			// Should store the value correctly
			storedValue, exists := builder.fields.Get(tt.fieldName)
			assert.True(t, exists)
			assert.Equal(t, tt.expected, storedValue)
		})
	}
}

func TestStatementBuilder_SetStringField(t *testing.T) {
	t.Parallel()

	tests := []struct {
		name      string
		fieldName string
		value     string
		expected  string
	}{
		{
			name:      "simple string",
			fieldName: "description",
			value:     "hello world",
			expected:  `"hello world"`,
		},
		{
			name:      "empty string",
			fieldName: "empty",
			value:     "",
			expected:  `""`,
		},
		{
			name:      "string with quotes",
			fieldName: "quoted",
			value:     `hello "world"`,
			expected:  `"hello \"world\""`,
		},
		{
			name:      "string with backslashes",
			fieldName: "path",
			value:     `C:\Program Files`,
			expected:  `"C:\\Program Files"`,
		},
		{
			name:      "string with newlines",
			fieldName: "multiline",
			value:     "line1\nline2",
			expected:  `"line1\nline2"`,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			// given
			builder := NewStatementBuilder("test_statement", nil, map[string]attribute.CompatibilityOptions{})

			// when
			result := builder.SetStringField(tt.fieldName, tt.value, "")

			// then
			// Should return the same builder for chaining
			assert.Equal(t, builder, result)

			// Should store the escaped value
			storedValue, exists := builder.fields.Get(tt.fieldName)
			assert.True(t, exists)
			assert.Equal(t, tt.expected, storedValue)
		})
	}
}

func TestStatementBuilder_SetCommandField(t *testing.T) {
	t.Parallel()

	tests := []struct {
		name      string
		fieldName string
		value     string
		expected  string
	}{
		{
			name:      "simple command",
			fieldName: "command",
			value:     "echo hello",
			expected:  `"echo hello"`,
		},
		{
			name:      "empty command",
			fieldName: "empty_cmd",
			value:     "",
			expected:  `""`,
		},
		{
			name:      "command with quotes (properly escaped)",
			fieldName: "quoted_cmd",
			value:     `echo "hello world"`,
			expected:  `"echo \"hello world\""`,
		},
		{
			name:      "command with backslashes (properly escaped)",
			fieldName: "path_cmd",
			value:     `find C:\Program Files`,
			expected:  `"find C:\\Program Files"`,
		},
		{
			name:      "command with newlines (properly escaped)",
			fieldName: "multiline_cmd",
			value:     "line1\nline2",
			expected:  `"line1\nline2"`,
		},
		{
			name:      "complex op language command",
			fieldName: "op_lang",
			value:     `(cpu_usage > 80 && memory < 1GB) | count > 5`,
			expected:  `"(cpu_usage > 80 && memory < 1GB) | count > 5"`,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			// given
			builder := NewStatementBuilder("test_statement", nil, map[string]attribute.CompatibilityOptions{})

			// when
			result := builder.SetCommandField(tt.fieldName, tt.value, "")

			// then
			// Should return the same builder for chaining
			assert.Equal(t, builder, result)

			// Should store the quoted but not escaped value
			storedValue, exists := builder.fields.Get(tt.fieldName)
			assert.True(t, exists)
			assert.Equal(t, tt.expected, storedValue)
		})
	}
}

func TestStatementBuilder_SetArrayField(t *testing.T) {
	t.Parallel()

	tests := []struct {
		name      string
		fieldName string
		value     []string
		expected  string
	}{
		{
			name:      "empty array",
			fieldName: "empty_list",
			value:     []string{},
			expected:  "[]",
		},
		{
			name:      "nil array",
			fieldName: "nil_list",
			value:     nil,
			expected:  "[]",
		},
		{
			name:      "single element",
			fieldName: "single_item",
			value:     []string{"item1"},
			expected:  `["item1"]`,
		},
		{
			name:      "multiple elements",
			fieldName: "multi_items",
			value:     []string{"item1", "item2", "item3"},
			expected:  `["item1", "item2", "item3"]`,
		},
		{
			name:      "elements with special characters",
			fieldName: "special_items",
			value:     []string{`hello "world"`, "item\nwith\nnewlines"},
			expected:  `["hello \"world\"", "item\nwith\nnewlines"]`,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			// given
			builder := NewStatementBuilder("test_statement", nil, map[string]attribute.CompatibilityOptions{})

			// when
			result := builder.SetArrayField(tt.fieldName, tt.value, "")

			// then
			// Should return the same builder for chaining
			assert.Equal(t, builder, result)

			// Should store the OpLang formatted array
			storedValue, exists := builder.fields.Get(tt.fieldName)
			assert.True(t, exists)
			assert.Equal(t, tt.expected, storedValue)
		})
	}
}

func TestStatementBuilder_Build(t *testing.T) {
	t.Parallel()

	tests := []struct {
		name     string
		setup    func(*StatementBuilder)
		expected string
	}{
		{
			name: "empty statement",
			setup: func(b *StatementBuilder) {
				// No fields added
			},
			expected: "test_statement()",
		},
		{
			name: "single field",
			setup: func(b *StatementBuilder) {
				b.SetField("name", "test_value", "")
			},
			expected: "test_statement(name=test_value)",
		},
		{
			name: "multiple fields with different types",
			setup: func(b *StatementBuilder) {
				b.SetField("name", "test_action", "").
					SetField("enabled", true, "").
					SetField("timeout", 30, "")
			},
			expected: "test_statement(name=test_action, enabled=true, timeout=30)",
		},
		{
			name: "mixed field types (string, command, array)",
			setup: func(b *StatementBuilder) {
				b.SetStringField("description", "Test description", "").
					SetCommandField("command", "echo hello", "").
					SetArrayField("params", []string{"param1", "param2"}, "")
			},
			expected: `test_statement(description="Test description", command="echo hello", params=["param1", "param2"])`,
		},
		{
			name: "field order is preserved",
			setup: func(b *StatementBuilder) {
				b.SetField("first", "1", "").
					SetField("second", "2", "").
					SetField("third", "3", "")
			},
			expected: "test_statement(first=1, second=2, third=3)",
		},
		{
			name: "overwrite existing field",
			setup: func(b *StatementBuilder) {
				b.SetField("name", "original", "").
					SetField("name", "updated", "")
			},
			expected: "test_statement(name=updated)",
		},
		{
			name: "complex realistic example",
			setup: func(b *StatementBuilder) {
				b.SetStringField("action_name", "cpu_monitor", "").
					SetCommandField("command", "cpu_usage > 80", "").
					SetField("enabled", true, "").
					SetField("timeout", 60, "").
					SetStringField("description", `Monitor CPU with "alerts"`, "").
					SetArrayField("tags", []string{"monitoring", "cpu"}, "")
			},
			expected: `test_statement(action_name="cpu_monitor", command="cpu_usage > 80", enabled=true, timeout=60, description="Monitor CPU with \"alerts\"", tags=["monitoring", "cpu"])`,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			// given
			builder := NewStatementBuilder("test_statement", nil, map[string]attribute.CompatibilityOptions{})
			tt.setup(builder)

			// when
			result := builder.Build()

			// then
			assert.Equal(t, tt.expected, result)
		})
	}
}

func TestStatementBuilder_MethodChaining(t *testing.T) {
	t.Parallel()

	t.Run("all methods return builder for chaining", func(t *testing.T) {
		// given
		builder := NewStatementBuilder("chain_test", nil, map[string]attribute.CompatibilityOptions{})

		// when
		result := builder.
			SetField("field1", "value1", "").
			SetStringField("field2", "value2", "").
			SetCommandField("field3", "value3", "").
			SetArrayField("field4", []string{"item1", "item2"}, "")

		// then
		// All methods should return the same builder instance
		assert.Equal(t, builder, result)

		// Final statement should include all fields
		statement := builder.Build()
		expected := `chain_test(field1=value1, field2="value2", field3="value3", field4=["item1", "item2"])`
		assert.Equal(t, expected, statement)
	})
}

func TestStatementBuilder_EdgeCases(t *testing.T) {
	t.Parallel()

	t.Run("field name with special characters", func(t *testing.T) {
		// given
		builder := NewStatementBuilder("test", nil, map[string]attribute.CompatibilityOptions{})

		// when
		builder.SetField("field-with_special.chars", "value", "")
		result := builder.Build()

		// then
		assert.Equal(t, "test(field-with_special.chars=value)", result)
	})

	t.Run("empty field name", func(t *testing.T) {
		// given
		builder := NewStatementBuilder("test", nil, map[string]attribute.CompatibilityOptions{})

		// when
		builder.SetField("", "value", "")
		result := builder.Build()

		// then
		assert.Equal(t, "test(=value)", result)
	})

	t.Run("field value with equals sign", func(t *testing.T) {
		// given
		builder := NewStatementBuilder("test", nil, map[string]attribute.CompatibilityOptions{})

		// when
		builder.SetStringField("equation", "x=y+z", "")
		result := builder.Build()

		// then
		assert.Equal(t, `test(equation="x=y+z")`, result)
	})

	t.Run("field value with commas", func(t *testing.T) {
		// given
		builder := NewStatementBuilder("test", nil, map[string]attribute.CompatibilityOptions{})

		// when
		builder.SetStringField("csv", "a,b,c", "")
		result := builder.Build()

		// then
		assert.Equal(t, `test(csv="a,b,c")`, result)
	})

	t.Run("field value with parentheses", func(t *testing.T) {
		// given
		builder := NewStatementBuilder("test", nil, map[string]attribute.CompatibilityOptions{})

		// when
		builder.SetStringField("expression", "(a+b)*c", "")
		result := builder.Build()

		// then
		assert.Equal(t, `test(expression="(a+b)*c")`, result)
	})
}

// Test the difference between SetStringField and SetCommandField
func TestStatementBuilder_StringVsCommandFields(t *testing.T) {
	t.Parallel()

	t.Run("string field escapes quotes", func(t *testing.T) {
		// given
		builder := NewStatementBuilder("test", nil, map[string]attribute.CompatibilityOptions{})

		// when
		builder.SetStringField("description", `Action with "quotes"`, "")
		result := builder.Build()

		// then
		expected := `test(description="Action with \"quotes\"")`
		assert.Equal(t, expected, result)
	})

	t.Run("command field properly escapes quotes", func(t *testing.T) {
		// given
		builder := NewStatementBuilder("test", nil, map[string]attribute.CompatibilityOptions{})

		// when
		builder.SetCommandField("query", `host.name == "web-server"`, "")
		result := builder.Build()

		// then
		expected := `test(query="host.name == \"web-server\"")`
		assert.Equal(t, expected, result)
	})

	t.Run("comparison of escaping behavior", func(t *testing.T) {
		// given
		input := `value with "quotes" and \backslashes`

		// when
		stringBuilder := NewStatementBuilder("string_test", nil, map[string]attribute.CompatibilityOptions{})
		stringBuilder.SetStringField("field", input, "")
		stringResult := stringBuilder.Build()

		commandBuilder := NewStatementBuilder("command_test", nil, map[string]attribute.CompatibilityOptions{})
		commandBuilder.SetCommandField("field", input, "")
		commandResult := commandBuilder.Build()

		// then
		// Both string and command fields should escape the quotes and backslashes for security
		expectedStringResult := `string_test(field="value with \"quotes\" and \\backslashes")`
		assert.Equal(t, expectedStringResult, stringResult)

		expectedCommandResult := `command_test(field="value with \"quotes\" and \\backslashes")`
		assert.Equal(t, expectedCommandResult, commandResult)

		// Both methods should now produce the same secure output
		assert.Equal(t, expectedStringResult, stringResult)
		assert.Equal(t, expectedCommandResult, commandResult)
	})
}

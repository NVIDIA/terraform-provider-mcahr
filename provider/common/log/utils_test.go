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

package log

import (
	"os"
	"testing"

	"github.com/hashicorp/go-hclog"
	"github.com/stretchr/testify/assert"
)

func TestParseLogLevel(t *testing.T) {
	t.Parallel()

	tests := []struct {
		name     string
		input    string
		expected hclog.Level
	}{
		{"TRACE uppercase", "TRACE", hclog.Trace},
		{"DEBUG uppercase", "DEBUG", hclog.Debug},
		{"INFO uppercase", "INFO", hclog.Info},
		{"WARN uppercase", "WARN", hclog.Warn},
		{"ERROR uppercase", "ERROR", hclog.Error},
		{"trace lowercase", "trace", hclog.Trace},
		{"debug lowercase", "debug", hclog.Debug},
		{"info lowercase", "info", hclog.Info},
		{"warn lowercase", "warn", hclog.Warn},
		{"error lowercase", "error", hclog.Error},
		{"INVALID value", "INVALID", hclog.NoLevel},
		{"empty string", "", hclog.NoLevel},
		{"Trace mixed case", "Trace", hclog.Trace},
		{"Debug mixed case", "Debug", hclog.Debug},
		{"Info mixed case", "Info", hclog.Info},
		{"Warning wrong keyword", "Warning", hclog.NoLevel},
		{"Off invalid level", "Off", hclog.NoLevel},
		{"DEBUG with whitespace", " DEBUG ", hclog.NoLevel},
		{"numeric value", "123", hclog.NoLevel},
	}

	for _, tt := range tests {
		tt := tt // capture range variable
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()

			// given
			input := tt.input

			// when
			result := parseLogLevel(input)

			// then
			assert.Equal(t, tt.expected, result)
		})
	}
}

func TestParseLogLevelFromEnv(t *testing.T) {
	t.Parallel()

	tests := []struct {
		name     string
		envVar   string
		envValue string
		expected hclog.Level
	}{
		{
			name:     "parses TRACE from env var",
			envVar:   "TEST_LOG_LEVEL_TRACE",
			envValue: "TRACE",
			expected: hclog.Trace,
		},
		{
			name:     "parses DEBUG from env var",
			envVar:   "TEST_LOG_LEVEL_DEBUG",
			envValue: "DEBUG",
			expected: hclog.Debug,
		},
		{
			name:     "parses INFO from env var",
			envVar:   "TEST_LOG_LEVEL_INFO",
			envValue: "INFO",
			expected: hclog.Info,
		},
		{
			name:     "parses WARN from env var",
			envVar:   "TEST_LOG_LEVEL_WARN",
			envValue: "WARN",
			expected: hclog.Warn,
		},
		{
			name:     "parses ERROR from env var",
			envVar:   "TEST_LOG_LEVEL_ERROR",
			envValue: "ERROR",
			expected: hclog.Error,
		},
		{
			name:     "returns NoLevel when env var not set",
			envVar:   "TEST_LOG_LEVEL_NOT_SET",
			envValue: "",
			expected: hclog.NoLevel,
		},
		{
			name:     "returns NoLevel for invalid value",
			envVar:   "TEST_LOG_LEVEL_INVALID",
			envValue: "INVALID",
			expected: hclog.NoLevel,
		},
		{
			name:     "case insensitive parsing lowercase",
			envVar:   "TEST_LOG_LEVEL_LOWERCASE",
			envValue: "info",
			expected: hclog.Info,
		},
		{
			name:     "case insensitive parsing mixed case",
			envVar:   "TEST_LOG_LEVEL_MIXEDCASE",
			envValue: "Warn",
			expected: hclog.Warn,
		},
		{
			name:     "handles whitespace in env value",
			envVar:   "TEST_LOG_LEVEL_WHITESPACE",
			envValue: " ERROR ",
			expected: hclog.NoLevel,
		},
	}

	for _, tt := range tests {
		tt := tt // capture range variable
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()

			// given
			original := os.Getenv(tt.envVar)
			os.Unsetenv(tt.envVar)
			if tt.envValue != "" {
				os.Setenv(tt.envVar, tt.envValue)
			}

			defer func() {
				os.Unsetenv(tt.envVar)
				if original != "" {
					os.Setenv(tt.envVar, original)
				}
			}()

			// when
			result := parseLogLevelFromEnv(tt.envVar)

			// then
			assert.Equal(t, tt.expected, result)
		})
	}
}

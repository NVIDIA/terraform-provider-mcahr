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

func TestSubsystemConfig_GetLogLevelForResource(t *testing.T) {

	tests := []struct {
		name         string
		config       *SubsystemConfig
		resourceName string
		envVars      map[string]string
		expected     hclog.Level
	}{
		{
			name:         "uses environment variable override",
			config:       NewSubsystemConfig(hclog.Error),
			resourceName: "action_test1",
			envVars: map[string]string{
				"TF_LOG_PROVIDER_RESOURCE_ACTION_TEST1": "DEBUG",
			},
			expected: hclog.Debug,
		},
		{
			name:         "uses global default env var",
			config:       NewSubsystemConfig(hclog.Error),
			resourceName: "unknown_test2",
			envVars: map[string]string{
				"TF_LOG_PROVIDER_RESOURCE_ALL": "WARN",
			},
			expected: hclog.Warn,
		},
		{
			name:         "uses config default as fallback",
			config:       NewSubsystemConfig(hclog.Info),
			resourceName: "unknown_test3",
			envVars:      map[string]string{},
			expected:     hclog.Info,
		},
		{
			name:         "case insensitive environment variable",
			config:       NewSubsystemConfig(hclog.Error),
			resourceName: "bot_test4",
			envVars: map[string]string{
				"TF_LOG_PROVIDER_RESOURCE_BOT_TEST4": "TRACE",
			},
			expected: hclog.Trace,
		},
	}

	for _, tt := range tests {
		tt := tt // capture range variable
		t.Run(tt.name, func(t *testing.T) {

			// given
			originalEnv := make(map[string]string)
			for key, value := range tt.envVars {
				originalEnv[key] = os.Getenv(key)
				os.Setenv(key, value)
			}

			defer func() {
				for key := range tt.envVars {
					if originalValue, exists := originalEnv[key]; exists {
						os.Setenv(key, originalValue)
					} else {
						os.Unsetenv(key)
					}
				}
			}()

			// when
			result := tt.config.GetLogLevelForResource(tt.resourceName)

			// then
			assert.Equal(t, tt.expected, result)
		})
	}
}

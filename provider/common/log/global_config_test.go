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
	"testing"

	"github.com/hashicorp/go-hclog"
	"github.com/stretchr/testify/assert"
)

func TestInitializeGlobalConfig(t *testing.T) {
	// Reset global state
	GlobalSubsystemConfig = nil

	// when
	InitializeGlobalConfig()

	// then
	assert.NotNil(t, GlobalSubsystemConfig, "expected GlobalSubsystemConfig to be initialized")
	assert.Equal(t, DefaultLogLevel, GlobalSubsystemConfig.DefaultLevel)
}

func TestGetLogLevelForResourceWhenConfigIsNil(t *testing.T) {
	// given
	GlobalSubsystemConfig = nil

	// when
	level := GetLogLevelForResource("test_resource")

	// then
	assert.Equal(t, DefaultLogLevel, level, "expected default level when config is nil")
}

func TestGetLogLevelForResourceWhenConfigExists(t *testing.T) {
	// given
	customLevel := hclog.Debug
	GlobalSubsystemConfig = NewSubsystemConfig(customLevel)

	// when
	level := GetLogLevelForResource("test_resource")

	// then - should use the config's GetLogLevelForResource method
	// which will return the custom level if no env vars are set
	assert.Equal(t, customLevel, level)

	// cleanup
	GlobalSubsystemConfig = nil
}

func TestGetLogLevelForResourceMultipleResources(t *testing.T) {
	// given
	InitializeGlobalConfig()
	resources := []string{"action", "user", "policy", "organization"}

	// when/then - should work for all resource types
	for _, resource := range resources {
		level := GetLogLevelForResource(resource)
		assert.Equal(t, DefaultLogLevel, level, "expected default level for resource %s", resource)
	}

	// cleanup
	GlobalSubsystemConfig = nil
}

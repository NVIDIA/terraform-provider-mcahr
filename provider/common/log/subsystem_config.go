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
	"strings"

	"github.com/hashicorp/go-hclog"
)

// SubsystemConfig provides centralized logging configuration for all resources
// Configuration is driven entirely by environment variables
type SubsystemConfig struct {
	DefaultLevel hclog.Level
}

// NewSubsystemConfig creates a new subsystem configuration
func NewSubsystemConfig(defaultLevel hclog.Level) *SubsystemConfig {
	return &SubsystemConfig{
		DefaultLevel: defaultLevel,
	}
}

// GetLogLevelForResource returns the appropriate log level for a given resource
func (c *SubsystemConfig) GetLogLevelForResource(resourceName string) hclog.Level {
	// Check for resource-specific environment variable first (highest priority)
	envVar := "TF_LOG_PROVIDER_RESOURCE_" + strings.ToUpper(resourceName)
	if level := parseLogLevelFromEnv(envVar); level != hclog.NoLevel {
		return level
	}

	// Check for global provider default environment variable
	if level := parseLogLevelFromEnv("TF_LOG_PROVIDER_RESOURCE_ALL"); level != hclog.NoLevel {
		return level
	}

	// Fall back to configured default
	return c.DefaultLevel
}

// SetDefaultLevel sets the default log level for all resources
func (c *SubsystemConfig) SetDefaultLevel(level hclog.Level) *SubsystemConfig {
	c.DefaultLevel = level
	return c
}

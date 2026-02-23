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
	"github.com/hashicorp/go-hclog"
)

// Global subsystem configuration instance
var GlobalSubsystemConfig *SubsystemConfig

// NoLevel means that will inherit the log level from the parent context(i.e see TF_LOG, TF_LOG_PROVIDER)
var DefaultLogLevel = hclog.NoLevel

// InitializeGlobalConfig initializes the global subsystem configuration with default values
// This is a fallback if provider arguments are not used
func InitializeGlobalConfig() {
	GlobalSubsystemConfig = NewSubsystemConfig(DefaultLogLevel)
}

// GetLogLevelForResource is a convenience function to get log level for a resource
func GetLogLevelForResource(resourceName string) hclog.Level {
	if GlobalSubsystemConfig == nil {
		return DefaultLogLevel
	}
	return GlobalSubsystemConfig.GetLogLevelForResource(resourceName)
}

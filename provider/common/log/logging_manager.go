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
	"context"
	"strings"
	"terraform/terraform-provider/provider/common"

	"github.com/hashicorp/go-hclog"
	"github.com/hashicorp/terraform-plugin-log/tflog"
)

// LoggingManager manages the root context for creating subsystems on-demand
type LoggingManager struct {
}

// Global logging manager instance
var globalLoggingManager = &LoggingManager{}

// IsDebugEnabled checks if debug logging is enabled at global level
// Use this in order to avoid expensive debug log computations
func IsDebugEnabled(requestCtx *common.RequestContext) bool {

	// Check general provider log level
	if level := parseLogLevelFromEnv("TF_LOG_PROVIDER"); isDebugOrTrace(level) {
		return true
	}

	// Check global TF_LOG
	if level := parseLogLevelFromEnv("TF_LOG"); isDebugOrTrace(level) {
		return true
	}

	return false
}

// isDebugOrTrace checks if a log level is Debug or Trace
func isDebugOrTrace(level hclog.Level) bool {
	return level != hclog.NoLevel && level <= hclog.Debug
}

// CreateResourceLogContextWithFields creates a subsystem context with persistent fields
func CreateResourceLogContextWithFields(baseCtx context.Context, resourceType string, fields map[string]any) context.Context {
	// Create subsystem context for this resource type with appropriate log level
	logLevel := GetLogLevelForResource(resourceType)
	subsystemCtx := tflog.NewSubsystem(baseCtx, strings.ToUpper(resourceType), tflog.WithLevel(logLevel))

	// Add resource_type as a default persistent field that will appear in ALL logs for this subsystem
	subsystemName := strings.ToUpper(resourceType)
	subsystemCtx = tflog.SubsystemSetField(subsystemCtx, subsystemName, "resource_type", resourceType)

	// Add additional persistent fields that will appear in ALL logs for this subsystem
	for key, value := range fields {
		subsystemCtx = tflog.SubsystemSetField(subsystemCtx, subsystemName, key, value)
	}

	return subsystemCtx
}

// LogInfo logs an informational message to the appropriate subsystem
func LogInfo(requestCtx *common.RequestContext, message string, fields map[string]any) {
	subsystemName := strings.ToUpper(requestCtx.ResourceType)
	var logFields map[string]any
	if fields != nil {
		logFields = fields
	}
	tflog.SubsystemInfo(requestCtx.Context, subsystemName, message, logFields)
}

// LogError logs an error message to the appropriate subsystem
func LogError(requestCtx *common.RequestContext, message string, fields map[string]any) {
	subsystemName := strings.ToUpper(requestCtx.ResourceType)
	var logFields map[string]any
	if fields != nil {
		logFields = fields
	}
	tflog.SubsystemError(requestCtx.Context, subsystemName, message, logFields)
}

// LogDebug logs a debug message to the appropriate subsystem
func LogDebug(requestCtx *common.RequestContext, message string, fields map[string]any) {
	subsystemName := strings.ToUpper(requestCtx.ResourceType)
	var logFields map[string]any
	if fields != nil {
		logFields = fields
	}
	tflog.SubsystemDebug(requestCtx.Context, subsystemName, message, logFields)
}

// LogWarn logs a warning message to the appropriate subsystem
func LogWarn(requestCtx *common.RequestContext, message string, fields map[string]any) {
	subsystemName := strings.ToUpper(requestCtx.ResourceType)
	var logFields map[string]any
	if fields != nil {
		logFields = fields
	}
	tflog.SubsystemWarn(requestCtx.Context, subsystemName, message, logFields)
}

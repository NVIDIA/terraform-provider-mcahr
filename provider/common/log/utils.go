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
	"strings"

	"github.com/hashicorp/go-hclog"
)

// parseLogLevel parses a log level string to hclog.Level
func parseLogLevel(logLevelStr string) hclog.Level {
	switch strings.ToUpper(logLevelStr) {
	case "TRACE":
		return hclog.Trace
	case "DEBUG":
		return hclog.Debug
	case "INFO":
		return hclog.Info
	case "WARN":
		return hclog.Warn
	case "ERROR":
		return hclog.Error
	default:
		return hclog.NoLevel // Invalid level
	}
}

// parseLogLevelFromEnv parses log level from environment variable
func parseLogLevelFromEnv(envVar string) hclog.Level {
	if logLevelStr := os.Getenv(envVar); logLevelStr != "" {
		return parseLogLevel(logLevelStr)
	}
	return hclog.NoLevel
}

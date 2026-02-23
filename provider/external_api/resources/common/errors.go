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

package common

import "strings"

// Error represents an error in V2 API responses
type Error struct {
	Message string `json:"message"`
	Type    string `json:"type"`
}

// FormatErrors formats V2 API errors into a string representation
func FormatErrors(status string, errors []Error) string {
	if len(errors) == 0 {
		return ""
	}

	result := "Status: " + status + "; Errors: "
	errorStrings := make([]string, len(errors))
	for i, err := range errors {
		errorStrings[i] = err.Type + ": " + err.Message
	}
	result += strings.Join(errorStrings, ", ")
	return result
}

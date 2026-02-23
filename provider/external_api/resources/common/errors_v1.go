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

// ValidationError represents a single validation error with field and message
type ValidationError struct {
	Field   string `json:"field,omitempty"`
	Message string `json:"message"`
}

// ErrorV1 represents error information in V1 API responses
type ErrorV1 struct {
	Message          string            `json:"message"`
	Type             string            `json:"type"`
	ValidationErrors []ValidationError `json:"validation_errors,omitempty"`
}

// SyntaxErrorsV1 represents top-level syntax error structure in V1 API responses
type SyntaxErrorsV1 struct {
	Root []string `json:"$"`
}

// V1ErrorContainer interface for containers that hold V1 errors
type V1ErrorContainer interface {
	GetNestedError() ErrorV1
	GetDirectErrors() []string
}

// FormatV1Error formats a V1 API error struct into a string representation
func FormatV1Error(err ErrorV1) string {
	if err.Type == "OK" || err.Type == "" {
		return ""
	}

	result := "Error Type: " + err.Type
	if err.Message != "" {
		result += "; Message: " + err.Message
	}

	// Add validation errors if present
	if len(err.ValidationErrors) > 0 {
		validationStr := formatValidationErrors(err.ValidationErrors)
		if validationStr != "" {
			result += "; Validation Errors: " + validationStr
		}
	}

	return result
}

// FormatSyntaxErrors formats top-level syntax errors (e.g., JSON parsing errors)
func FormatSyntaxErrors(errors []string) string {
	if len(errors) == 0 {
		return ""
	}

	result := "Syntax Errors: "
	for i, err := range errors {
		if i > 0 {
			result += "; "
		}
		result += err
	}
	return result
}

// FormatDirectValidationErrors formats direct validation error strings
func FormatDirectValidationErrors(errors []string) string {
	if len(errors) == 0 {
		return ""
	}

	result := "Validation Errors: "
	for i, err := range errors {
		if i > 0 {
			result += "; "
		}
		result += err
	}
	return result
}

// FormatV1ErrorsWithPriority formats V1 API errors with priority checking:
// Priority 1: Top-level syntax errors (if provided)
// Priority 2: Container-level direct validation errors (if provided)
// Priority 3: Standard nested error structure
func FormatV1ErrorsWithPriority(syntaxErrors *SyntaxErrorsV1, directErrors []string, nestedError ErrorV1) string {
	// Priority 1: Check for top-level syntax errors (highest priority)
	if syntaxErrors != nil {
		if syntaxError := FormatSyntaxErrors(syntaxErrors.Root); syntaxError != "" {
			return syntaxError
		}
	}

	// Priority 2: Check for direct validation errors array
	if directError := FormatDirectValidationErrors(directErrors); directError != "" {
		return directError
	}

	// Priority 3: Check nested error structure
	return FormatV1Error(nestedError)
}

// FormatV1ErrorsFromContainer formats V1 API errors from a container with priority checking:
// Priority 1: Top-level syntax errors (if provided)
// Priority 2: Container-level direct validation errors
// Priority 3: Standard nested error structure
func FormatV1ErrorsFromContainer(syntaxErrors *SyntaxErrorsV1, container V1ErrorContainer) string {
	return FormatV1ErrorsWithPriority(syntaxErrors, container.GetDirectErrors(), container.GetNestedError())
}

// formatValidationErrors converts validation errors to a string representation
func formatValidationErrors(validationErrors []ValidationError) string {
	var errorStrings []string
	for _, err := range validationErrors {
		var formatted string
		if err.Field != "" {
			formatted = err.Field + ": " + err.Message
		} else {
			formatted = err.Message
		}
		if formatted != "" {
			errorStrings = append(errorStrings, formatted)
		}
	}
	return strings.Join(errorStrings, ", ")
}

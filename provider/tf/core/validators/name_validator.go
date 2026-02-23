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

package validators

import (
	"context"
	"fmt"
	"regexp"

	"github.com/hashicorp/terraform-plugin-framework/schema/validator"
)

// nameValidator validates that a string follows entity naming rules.
type nameValidator struct{}

// NameValidator returns a validator which ensures that the string value:
// - Only contains alphanumeric characters and underscores
// - Starts with a letter or underscore
//
// This validator is intended for entity names (action, alarm, resource, etc.)
func NameValidator() validator.String {
	return nameValidator{}
}

func (v nameValidator) Description(_ context.Context) string {
	return "value must be an alphanumeric/underscore string that starts with a letter or underscore"
}

func (v nameValidator) MarkdownDescription(ctx context.Context) string {
	return v.Description(ctx)
}

func (v nameValidator) ValidateString(ctx context.Context, request validator.StringRequest, response *validator.StringResponse) {
	if request.ConfigValue.IsNull() || request.ConfigValue.IsUnknown() {
		return
	}

	value := request.ConfigValue.ValueString()

	// Check if the string contains only alphanumeric characters and underscores
	alphanumericPattern := regexp.MustCompile("^[a-zA-Z0-9_]*$")
	if !alphanumericPattern.MatchString(value) {
		response.Diagnostics.AddAttributeError(
			request.Path,
			"Invalid Name Format",
			fmt.Sprintf("Name must contain only alphanumeric characters and underscores, got: '%s'", value),
		)
		return
	}

	// Check if the string starts with a letter or underscore
	startsWithLetterOrUnderscore := regexp.MustCompile("^[a-zA-Z_]")
	if !startsWithLetterOrUnderscore.MatchString(value) {
		response.Diagnostics.AddAttributeError(
			request.Path,
			"Invalid Name Start",
			fmt.Sprintf("Name must start with a letter or underscore, got: '%s'", value),
		)
		return
	}
}

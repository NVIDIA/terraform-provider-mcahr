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

	"github.com/hashicorp/terraform-plugin-framework/schema/validator"
)

// exactValueValidator validates that a string matches exactly the required value.
type exactValueValidator struct {
	requiredValue string
}

// ExactValueValidator returns a validator which ensures that the string value
// matches exactly the provided required value.
//
// This validator is intended for singleton resources or any exact value requirements.
func ExactValueValidator(requiredValue string) validator.String {
	return exactValueValidator{
		requiredValue: requiredValue,
	}
}

func (v exactValueValidator) Description(_ context.Context) string {
	return fmt.Sprintf("value must be exactly '%s'", v.requiredValue)
}

func (v exactValueValidator) MarkdownDescription(ctx context.Context) string {
	return v.Description(ctx)
}

func (v exactValueValidator) ValidateString(ctx context.Context, request validator.StringRequest, response *validator.StringResponse) {

	value := request.ConfigValue.ValueString()

	if value != v.requiredValue {
		response.Diagnostics.AddAttributeError(
			request.Path,
			"Invalid Value",
			fmt.Sprintf("The value must be exactly '%s', got: '%s'", v.requiredValue, value),
		)
	}
}

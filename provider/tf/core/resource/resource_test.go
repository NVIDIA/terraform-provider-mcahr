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

package resource

import (
	"context"
	"testing"

	"terraform/terraform-provider/provider/common/attribute"
	coreschema "terraform/terraform-provider/provider/tf/core/schema"
	"terraform/terraform-provider/provider/tf/core/validators"

	"github.com/hashicorp/terraform-plugin-framework/resource"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema"
	"github.com/hashicorp/terraform-plugin-framework/schema/validator"
	"github.com/stretchr/testify/assert"
)

// stubSchema stands in for a resource schema with a configurable set of
// validators on the `name` attribute.
type stubSchema struct {
	nameValidators []validator.String
	omitName       bool
}

var _ coreschema.ResourceSchema = &stubSchema{}

func (s *stubSchema) GetSchema() schema.Schema {
	if s.omitName {
		return schema.Schema{Attributes: map[string]schema.Attribute{}}
	}

	return schema.Schema{
		Attributes: map[string]schema.Attribute{
			"name": schema.StringAttribute{
				Required:   true,
				Validators: s.nameValidators,
			},
		},
	}
}

func (s *stubSchema) GetCompatibilityOptions() map[string]attribute.CompatibilityOptions {
	return map[string]attribute.CompatibilityOptions{}
}

func (s *stubSchema) GetFieldComparisonRules() map[string]coreschema.FieldComparisonRule {
	return coreschema.DefaultFieldComparisonRules()
}

func TestValidateImportID(t *testing.T) {
	t.Parallel()

	validated := &stubSchema{nameValidators: []validator.String{validators.NameValidator()}}
	unvalidated := &stubSchema{}

	tests := []struct {
		name           string
		importID       string
		resourceSchema coreschema.ResourceSchema
		wantValid      bool
	}{
		{
			name:           "plain name on a validated resource",
			importID:       "my_resource",
			resourceSchema: validated,
			wantValid:      true,
		},
		{
			name:           "plain name on an unvalidated resource",
			importID:       "my_resource",
			resourceSchema: unvalidated,
			wantValid:      true,
		},
		{
			// integration and runbook declare no validator on `name`, so a
			// dashed name is importable there and must stay importable.
			name:           "dashed name on an unvalidated resource stays importable",
			importID:       "my-integration-2024",
			resourceSchema: unvalidated,
			wantValid:      true,
		},
		{
			name:           "dashed name rejected where the schema forbids it",
			importID:       "my-resource",
			resourceSchema: validated,
			wantValid:      false,
		},
		{
			// The injection vector: ImportStatePassthroughID wrote this
			// straight into `name`, from where it reached statement builders.
			name:           "quote rejected on an unvalidated resource",
			importID:       `x") delete_integration(integration_name="y`,
			resourceSchema: unvalidated,
			wantValid:      false,
		},
		{
			name:           "backslash rejected on an unvalidated resource",
			importID:       `back\slash`,
			resourceSchema: unvalidated,
			wantValid:      false,
		},
		{
			name:           "control character rejected on an unvalidated resource",
			importID:       "line\nbreak",
			resourceSchema: unvalidated,
			wantValid:      false,
		},
		{
			name:           "resource without a name attribute falls back to the baseline check",
			importID:       "anything-goes",
			resourceSchema: &stubSchema{omitName: true},
			wantValid:      true,
		},
		{
			name:           "resource without a name attribute still rejects quotes",
			importID:       `has"quote`,
			resourceSchema: &stubSchema{omitName: true},
			wantValid:      false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()

			resp := &resource.ImportStateResponse{}

			got := validateImportID(context.Background(), tt.importID, tt.resourceSchema, resp)

			assert.Equal(t, tt.wantValid, got)
			assert.Equal(t, !tt.wantValid, resp.Diagnostics.HasError())
		})
	}
}

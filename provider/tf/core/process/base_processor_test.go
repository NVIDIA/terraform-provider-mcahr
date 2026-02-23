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

package process

import (
	"context"
	"testing"

	"terraform/terraform-provider/provider/common"
	"terraform/terraform-provider/provider/common/version"
	model "terraform/terraform-provider/provider/tf/core/model"

	"github.com/hashicorp/terraform-plugin-framework/diag"
	"github.com/hashicorp/terraform-plugin-framework/path"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

// TestTFModel is a test implementation of TFModel interface
type TestTFModel struct {
	BackendVersion *version.BackendVersion `json:"-" tfsdk:"-"`

	ID          string
	Name        string
	Description string
	Enabled     bool
}

var _ model.TFModel = &TestTFModel{}

func (m *TestTFModel) GetName() string {
	return m.Name
}

// Ensure TestTFModel implements the TFModel interface
var _ model.TFModel = &TestTFModel{}

// MockGetter implements the Getter interface for testing
type MockGetter struct {
	GetFunc          func(ctx context.Context, target interface{}) diag.Diagnostics
	GetAttributeFunc func(ctx context.Context, path path.Path, target interface{}) diag.Diagnostics
}

func (m *MockGetter) Get(ctx context.Context, target interface{}) diag.Diagnostics {
	if m.GetFunc != nil {
		return m.GetFunc(ctx, target)
	}
	return diag.Diagnostics{}
}

func (m *MockGetter) GetAttribute(ctx context.Context, path path.Path, target interface{}) diag.Diagnostics {
	if m.GetAttributeFunc != nil {
		return m.GetAttributeFunc(ctx, path, target)
	}
	return diag.Diagnostics{}
}

// TestBasePreProcessor_ExtractFrom_Success tests successful extraction
func TestBasePreProcessor_ExtractFrom_Success(t *testing.T) {
	t.Parallel()
	// given
	processor := &BasePreProcessor[*TestTFModel]{}
	requestContext := common.NewRequestContext(context.Background())
	expectedModel := TestTFModel{
		ID:          "test-123",
		Name:        "test-name",
		Description: "test description",
		Enabled:     true,
	}

	mockGetter := &MockGetter{
		GetFunc: func(ctx context.Context, target interface{}) diag.Diagnostics {
			if model, ok := target.(*TestTFModel); ok {
				*model = expectedModel
			}
			return diag.Diagnostics{}
		},
	}

	// when
	result, err := processor.ExtractFrom(requestContext, mockGetter, &TestTFModel{})

	// then
	require.NoError(t, err)
	assert.NotNil(t, result)
	assert.Equal(t, expectedModel.ID, result.ID)
	assert.Equal(t, expectedModel.Name, result.Name)
	assert.Equal(t, expectedModel.Description, result.Description)
	assert.Equal(t, expectedModel.Enabled, result.Enabled)
}

// TestBasePreProcessor_ExtractFrom_WithDiagnosticErrors tests extraction with diagnostic errors
func TestBasePreProcessor_ExtractFrom_WithDiagnosticErrors(t *testing.T) {
	t.Parallel()
	// given
	processor := &BasePreProcessor[*TestTFModel]{}
	requestContext := common.NewRequestContext(context.Background())

	mockGetter := &MockGetter{
		GetFunc: func(ctx context.Context, target interface{}) diag.Diagnostics {
			diags := diag.Diagnostics{}
			diags.AddError("Extraction Error", "Failed to extract data from TF source")
			diags.AddWarning("Extraction Warning", "This is a warning")
			return diags
		},
	}

	// when
	result, err := processor.ExtractFrom(requestContext, mockGetter, &TestTFModel{})

	// then
	assert.Error(t, err)
	assert.Nil(t, result)
	assert.Contains(t, err.Error(), "failed to get data from TF source")
	assert.Contains(t, err.Error(), "Extraction Error")
}

// TestBasePreProcessor_ExtractFrom_WithAttributeErrors tests extraction with attribute-specific errors
func TestBasePreProcessor_ExtractFrom_WithAttributeErrors(t *testing.T) {
	t.Parallel()
	// given
	processor := &BasePreProcessor[*TestTFModel]{}
	requestContext := common.NewRequestContext(context.Background())

	mockGetter := &MockGetter{
		GetFunc: func(ctx context.Context, target interface{}) diag.Diagnostics {
			diags := diag.Diagnostics{}
			diags.AddAttributeError(path.Root("name"), "Name Error", "Name is invalid")
			diags.AddAttributeError(path.Root("enabled"), "Enabled Error", "Enabled value is invalid")
			return diags
		},
	}

	// when
	result, err := processor.ExtractFrom(requestContext, mockGetter, &TestTFModel{})

	// then
	assert.Error(t, err)
	assert.Nil(t, result)
	assert.Contains(t, err.Error(), "failed to get data from TF source")
	assert.Contains(t, err.Error(), "Name Error")
	assert.Contains(t, err.Error(), "Enabled Error")
}

// TestBasePreProcessor_ExtractFrom_WithMixedDiagnostics tests extraction with mixed error types
func TestBasePreProcessor_ExtractFrom_WithMixedDiagnostics(t *testing.T) {
	t.Parallel()
	// given
	processor := &BasePreProcessor[*TestTFModel]{}
	requestContext := common.NewRequestContext(context.Background())

	mockGetter := &MockGetter{
		GetFunc: func(ctx context.Context, target interface{}) diag.Diagnostics {
			diags := diag.Diagnostics{}
			diags.AddError("General Error", "General extraction failed")
			diags.AddAttributeError(path.Root("id"), "ID Error", "ID is required")
			diags.AddWarning("General Warning", "This is just a warning")
			return diags
		},
	}

	// when
	result, err := processor.ExtractFrom(requestContext, mockGetter, &TestTFModel{})

	// then
	assert.Error(t, err)
	assert.Nil(t, result)
	errorMessage := err.Error()
	assert.Contains(t, errorMessage, "failed to get data from TF source")
	assert.Contains(t, errorMessage, "General Error")
	assert.Contains(t, errorMessage, "ID Error")
}

// TestBasePreProcessor_ExtractFrom_WithWarningsOnly tests extraction with warnings but no errors
func TestBasePreProcessor_ExtractFrom_WithWarningsOnly(t *testing.T) {
	t.Parallel()
	// given
	processor := &BasePreProcessor[*TestTFModel]{}
	requestContext := common.NewRequestContext(context.Background())
	expectedModel := TestTFModel{
		ID:   "warning-test",
		Name: "test-with-warnings",
	}

	mockGetter := &MockGetter{
		GetFunc: func(ctx context.Context, target interface{}) diag.Diagnostics {
			if model, ok := target.(*TestTFModel); ok {
				*model = expectedModel
			}
			diags := diag.Diagnostics{}
			diags.AddWarning("Warning 1", "This is a warning")
			diags.AddWarning("Warning 2", "This is another warning")
			return diags
		},
	}

	// when
	result, err := processor.ExtractFrom(requestContext, mockGetter, &TestTFModel{})

	// then
	require.NoError(t, err) // Warnings don't cause errors
	assert.NotNil(t, result)
	assert.Equal(t, expectedModel.ID, result.ID)
	assert.Equal(t, expectedModel.Name, result.Name)
}

// TestBasePreProcessor_ExtractFrom_EmptyModel tests extraction with empty model
func TestBasePreProcessor_ExtractFrom_EmptyModel(t *testing.T) {
	t.Parallel()
	// given
	processor := &BasePreProcessor[*TestTFModel]{}
	requestContext := common.NewRequestContext(context.Background())

	mockGetter := &MockGetter{
		GetFunc: func(ctx context.Context, target interface{}) diag.Diagnostics {
			// Don't populate the target, leaving it with zero values
			return diag.Diagnostics{}
		},
	}

	// when
	result, err := processor.ExtractFrom(requestContext, mockGetter, &TestTFModel{})

	// then
	require.NoError(t, err)
	assert.NotNil(t, result)
	// Verify zero values
	assert.Equal(t, "", result.ID)
	assert.Equal(t, "", result.Name)
	assert.Equal(t, "", result.Description)
	assert.Equal(t, false, result.Enabled)
}

// TestBasePreProcessor_ExtractFrom_PartialModel tests extraction with partially populated model
func TestBasePreProcessor_ExtractFrom_PartialModel(t *testing.T) {
	t.Parallel()
	// given
	processor := &BasePreProcessor[*TestTFModel]{}
	requestContext := common.NewRequestContext(context.Background())

	mockGetter := &MockGetter{
		GetFunc: func(ctx context.Context, target interface{}) diag.Diagnostics {
			if model, ok := target.(*TestTFModel); ok {
				// Only populate some fields
				model.ID = "partial-123"
				model.Name = "partial-name"
				// Leave Description empty and Enabled as false
			}
			return diag.Diagnostics{}
		},
	}

	// when
	result, err := processor.ExtractFrom(requestContext, mockGetter, &TestTFModel{})

	// then
	require.NoError(t, err)
	assert.NotNil(t, result)
	assert.Equal(t, "partial-123", result.ID)
	assert.Equal(t, "partial-name", result.Name)
	assert.Equal(t, "", result.Description) // Zero value
	assert.Equal(t, false, result.Enabled)  // Zero value
}

// TestBasePreProcessor_ExtractFrom_NilContext tests extraction with nil context
func TestBasePreProcessor_ExtractFrom_NilContext(t *testing.T) {
	t.Parallel()
	// given
	processor := &BasePreProcessor[*TestTFModel]{}
	expectedModel := TestTFModel{
		ID:   "nil-context-test",
		Name: "test-nil-context",
	}

	mockGetter := &MockGetter{
		GetFunc: func(ctx context.Context, target interface{}) diag.Diagnostics {
			// Verify context is nil as passed through requestContext
			assert.Nil(t, ctx)
			if model, ok := target.(*TestTFModel); ok {
				*model = expectedModel
			}
			return diag.Diagnostics{}
		},
	}

	// Create RequestContext with nil context
	requestContext := &common.RequestContext{Context: nil}

	// when
	result, err := processor.ExtractFrom(requestContext, mockGetter, &TestTFModel{})

	// then
	require.NoError(t, err)
	assert.NotNil(t, result)
	assert.Equal(t, expectedModel.ID, result.ID)
	assert.Equal(t, expectedModel.Name, result.Name)
}

// TestBasePreProcessor_ExtractFrom_GetterPanic tests handling when getter panics
func TestBasePreProcessor_ExtractFrom_GetterPanic(t *testing.T) {
	t.Parallel()
	// given
	processor := &BasePreProcessor[*TestTFModel]{}
	requestContext := common.NewRequestContext(context.Background())

	mockGetter := &MockGetter{
		GetFunc: func(ctx context.Context, target interface{}) diag.Diagnostics {
			panic("getter panicked")
		},
	}

	// when/then - This should panic since we don't handle panics in the function
	assert.Panics(t, func() {
		processor.ExtractFrom(requestContext, mockGetter, &TestTFModel{})
	})
}

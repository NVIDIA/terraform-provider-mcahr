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
	"terraform/terraform-provider/provider/common"
	"terraform/terraform-provider/provider/common/log"
	api "terraform/terraform-provider/provider/external_api/resources"
	model "terraform/terraform-provider/provider/tf/core/model"

	"github.com/hashicorp/terraform-plugin-framework/diag"
)

// createSubsystemContext creates a subsystem context with standard persistent fields for resource operations
func createSubsystemContext[TF model.TFModel, API_V1 api.APIModel, API_V2 api.APIModel](ctx context.Context, params *CRUDOperationParams[TF, API_V1, API_V2], operation common.CrudOperation, resourceName string) context.Context {
	fields := map[string]interface{}{
		"resource_type": params.ResourceType,
		"operation":     operation.String(),
	}

	// Add resource name to MDC if provided
	if resourceName != "" {
		fields["resource_name"] = resourceName
	}

	return log.CreateResourceLogContextWithFields(ctx, params.ResourceType, fields)
}

// createSimpleSubsystemContext creates a subsystem context for operations without CRUDOperationParams (like import)
func createSimpleSubsystemContext(ctx context.Context, resourceType string, operation string, resourceName string) context.Context {
	fields := map[string]interface{}{
		"resource_type": resourceType,
		"operation":     operation,
	}

	// Add resource name to MDC if provided
	if resourceName != "" {
		fields["resource_name"] = resourceName
	}

	return log.CreateResourceLogContextWithFields(ctx, resourceType, fields)
}

// extractResourceName extracts the resource name from any Terraform data source (config, state, plan)
func extractResourceName[TF model.TFModel](ctx context.Context, getter interface {
	Get(context.Context, interface{}) diag.Diagnostics
}) string {
	var data TF
	if diags := getter.Get(ctx, &data); !diags.HasError() {
		return data.GetName()
	}
	return ""
}

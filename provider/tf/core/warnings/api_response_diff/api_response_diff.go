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

package apiresponsediff

import (
	"reflect"
	"terraform/terraform-provider/provider/common"
	corecommon "terraform/terraform-provider/provider/tf/core/common"
	model "terraform/terraform-provider/provider/tf/core/model"
	"terraform/terraform-provider/provider/tf/core/process"
	coremodel "terraform/terraform-provider/provider/tf/core/schema"
)

type fieldDifference struct {
	FieldName     string
	PlanValue     string
	ResponseValue string
}

// CheckPlanVsApiResponseDelta checks for API response differences and warns the user.
// This detects when the API normalized or modified user input.
// Uses the schema's GetFieldComparisonRules to handle fields with custom comparison logic.
func CheckPlanVsApiResponseDelta[TF model.TFModel](requestContext *common.RequestContext, processData *process.ProcessData, schema coremodel.ResourceSchema, apiResponseTfModel TF) error {

	planTfModel, err := getPlanTfModel[TF](requestContext, processData)
	if err != nil {
		return err
	}

	if common.IsNil(planTfModel) || common.IsNil(apiResponseTfModel) {
		return nil
	}

	// Get comparison rules from schema
	comparisonRules := schema.GetFieldComparisonRules()

	differences := compareModels(planTfModel, apiResponseTfModel, comparisonRules)

	if len(differences) > 0 {
		addWarningsToDiagnostics(requestContext, processData, differences)
	}

	return nil
}

// getPlanTfModel extracts the Terraform model from the plan
// based on the current operation (Create or Update).
// Uses Plan instead of Config to include computed fields like _full variants,
// allowing detection of API modifications to those fields.
func getPlanTfModel[TF model.TFModel](requestContext *common.RequestContext, processData *process.ProcessData) (TF, error) {

	// Create a new instance of TF, properly handling pointer types
	var tfModel TF
	tfType := reflect.TypeOf(tfModel)

	// If TF is a pointer type, we need to allocate the underlying struct
	if tfType != nil && tfType.Kind() == reflect.Ptr {
		// Create a new instance of the underlying type
		tfModel = reflect.New(tfType.Elem()).Interface().(TF)
	}

	switch requestContext.Operation {
	case common.Create:
		return corecommon.ExtractFromTfSource[TF](requestContext, processData.CreateRequest.Plan, tfModel)
	case common.Update:
		return corecommon.ExtractFromTfSource[TF](requestContext, processData.UpdateRequest.Plan, tfModel)
	}

	var nilTF TF
	return nilTF, nil
}

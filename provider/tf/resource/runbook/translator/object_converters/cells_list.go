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

package converters

import (
	"context"
	"fmt"
	"terraform/terraform-provider/provider/common"
	customattribute "terraform/terraform-provider/provider/external_api/resources/runbooks/custom_attribute"

	"github.com/hashicorp/terraform-plugin-framework/attr"
	"github.com/hashicorp/terraform-plugin-framework/diag"
	"github.com/hashicorp/terraform-plugin-framework/types"
)

// CellTFModel is the Terraform-side representation of a single cell in cells_list.
// Each cell must have exactly one of Op or Md set (enforced by schema validators and the converter).
// Fields mirror the CellJson internal model from custom_attribute/cell.go.
type CellTFModel struct {
	Op          types.String `tfsdk:"op"`
	Md          types.String `tfsdk:"md"`
	Name        types.String `tfsdk:"name"`
	Enabled     types.Bool   `tfsdk:"enabled"`
	SecretAware types.Bool   `tfsdk:"secret_aware"`
	Description types.String `tfsdk:"description"`
}

// CellsListAttrTypes defines the attribute type map for the nested cell object.
// Must stay in sync with the cells_list schema definition in schema.go and CellTFModel above.
var CellsListAttrTypes = map[string]attr.Type{
	"op":           types.StringType,
	"md":           types.StringType,
	"name":         types.StringType,
	"enabled":      types.BoolType,
	"secret_aware": types.BoolType,
	"description":  types.StringType,
}

// CellsListObjectType is the types.ObjectType for a single cell element in cells_list.
// Used when constructing types.List values for cells_list (e.g. types.ListValue(CellsListObjectType, ...)).
var CellsListObjectType = types.ObjectType{AttrTypes: CellsListAttrTypes}

// CellsListToInternalCells converts a Terraform cells_list (types.List of cell objects) into
// the internal []CellJson representation used by the translator and API layers.
//
// Returns nil, nil when the list is null or unknown (user did not set cells_list).
func CellsListToInternalCells(ctx context.Context, tfCellsList types.List) ([]customattribute.CellJson, error) {
	if tfCellsList.IsNull() || tfCellsList.IsUnknown() {
		return nil, nil
	}

	var cellModels []CellTFModel
	diags := tfCellsList.ElementsAs(ctx, &cellModels, false)
	if diags.HasError() {
		return nil, fmt.Errorf("failed to extract cells from list: %s", diags.Errors())
	}

	cells := make([]customattribute.CellJson, len(cellModels))
	for i, cm := range cellModels {
		cell, err := cellTFModelToInternal(&cm, i)
		if err != nil {
			return nil, err
		}
		cells[i] = *cell
	}

	return cells, nil
}

// CellsListFromAPICells converts API response cells ([]CellJsonAPI) into a Terraform types.List
// suitable for storing in the RunbookTFModel.CellsList field.
//
// Returns a null list when apiCells is nil (e.g. API returned no cells field).
func CellsListFromAPICells(apiCells []customattribute.CellJsonAPI) (types.List, diag.Diagnostics) {
	if apiCells == nil {
		return types.ListNull(CellsListObjectType), nil
	}

	cellObjects := make([]attr.Value, len(apiCells))
	for i, apiCell := range apiCells {
		obj, diags := internalCellToTFObject(apiCell.ToInternalModel())
		if diags.HasError() {
			return types.ListNull(CellsListObjectType), diags
		}
		cellObjects[i] = obj
	}

	result, diags := types.ListValue(CellsListObjectType, cellObjects)
	if diags.HasError() {
		return types.ListNull(CellsListObjectType), diags
	}
	return result, nil
}

// cellTFModelToInternal converts a single CellTFModel to the internal CellJson.
func cellTFModelToInternal(cm *CellTFModel, index int) (*customattribute.CellJson, error) {
	cell := &customattribute.CellJson{
		Name:        cm.Name.ValueString(),
		Enabled:     cm.Enabled.ValueBool(),
		SecretAware: cm.SecretAware.ValueBool(),
		Description: cm.Description.ValueString(),
		Op:          common.NewOptionalUnset[string](),
		Md:          common.NewOptionalUnset[string](),
	}

	hasOp := !cm.Op.IsNull() && !cm.Op.IsUnknown()
	hasMd := !cm.Md.IsNull() && !cm.Md.IsUnknown()

	if hasOp && hasMd {
		return nil, fmt.Errorf("cell at index %d has both op and md set; only one is allowed", index)
	}
	if !hasOp && !hasMd {
		return nil, fmt.Errorf("cell at index %d must have either op or md set", index)
	}

	if hasOp {
		cell.Op = common.NewOptional(cm.Op.ValueString())
	} else {
		cell.Md = common.NewOptional(cm.Md.ValueString())
	}

	return cell, nil
}

// internalCellToTFObject converts an internal CellJson to a Terraform types.Object.
func internalCellToTFObject(internal *customattribute.CellJson) (types.Object, diag.Diagnostics) {
	attrs := map[string]attr.Value{
		"name":         types.StringValue(internal.Name),
		"enabled":      types.BoolValue(internal.Enabled),
		"secret_aware": types.BoolValue(internal.SecretAware),
		"description":  types.StringValue(internal.Description),
	}

	if internal.Op.IsSet {
		attrs["op"] = types.StringValue(internal.Op.Get())
		attrs["md"] = types.StringNull()
	} else if internal.Md.IsSet {
		attrs["op"] = types.StringNull()
		attrs["md"] = types.StringValue(internal.Md.Get())
	} else {
		attrs["op"] = types.StringNull()
		attrs["md"] = types.StringNull()
	}

	return types.ObjectValue(CellsListAttrTypes, attrs)
}

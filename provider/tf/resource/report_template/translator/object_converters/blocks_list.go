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
	customattribute "terraform/terraform-provider/provider/external_api/resources/report_templates/custom_attribute"

	"github.com/hashicorp/terraform-plugin-framework/attr"
	"github.com/hashicorp/terraform-plugin-framework/diag"
	"github.com/hashicorp/terraform-plugin-framework/types"
)

// --- Nested sub-type definitions ---

var BreakdownTagValueAttrTypes = map[string]attr.Type{
	"color":  types.StringType,
	"label":  types.StringType,
	"values": types.ListType{ElemType: types.StringType},
}

var BreakdownTagValueObjectType = types.ObjectType{AttrTypes: BreakdownTagValueAttrTypes}

var BreakdownValueAttrTypes = map[string]attr.Type{
	"value": types.StringType,
	"count": types.Int64Type,
}

var BreakdownValueObjectType = types.ObjectType{AttrTypes: BreakdownValueAttrTypes}

var ResourcesBreakdownAttrTypes = map[string]attr.Type{
	"group_by_value":   types.StringType,
	"breakdown_values": types.ListType{ElemType: BreakdownValueObjectType},
}

var ResourcesBreakdownObjectType = types.ObjectType{AttrTypes: ResourcesBreakdownAttrTypes}

var GroupByTagOrderAttrTypes = map[string]attr.Type{
	"type":   types.StringType,
	"values": types.ListType{ElemType: types.StringType},
}

var GroupByTagOrderObjectType = types.ObjectType{AttrTypes: GroupByTagOrderAttrTypes}

// --- Block type definition ---

var BlocksListAttrTypes = map[string]attr.Type{
	"title":                               types.StringType,
	"resource_query":                      types.StringType,
	"group_by_tag":                        types.StringType,
	"breakdown_by_tag":                    types.StringType,
	"view_mode":                           types.StringType,
	"include_resources_without_group_tag": types.BoolType,
	"include_other_breakdown_tag_values":  types.BoolType,
	"other_tags_to_export":                types.ListType{ElemType: types.StringType},
	"group_by_tag_order":                  GroupByTagOrderObjectType,
	"breakdown_tags_values":               types.ListType{ElemType: BreakdownTagValueObjectType},
	"resources_breakdown":                 types.ListType{ElemType: ResourcesBreakdownObjectType},
}

var BlocksListObjectType = types.ObjectType{AttrTypes: BlocksListAttrTypes}

// --- ToInternal conversion ---

func BlocksListToInternal(ctx context.Context, tfList types.List) ([]customattribute.BlockJson, error) {
	if tfList.IsNull() || tfList.IsUnknown() {
		return nil, nil
	}

	elements := tfList.Elements()
	blocks := make([]customattribute.BlockJson, len(elements))

	for i, elem := range elements {
		obj, ok := elem.(types.Object)
		if !ok {
			return nil, fmt.Errorf("block at index %d is not an object", i)
		}
		attrs := obj.Attributes()

		block := customattribute.BlockJson{
			Title:                           stringAttr(attrs, "title"),
			ResourceQuery:                   stringAttr(attrs, "resource_query"),
			GroupByTag:                      stringAttr(attrs, "group_by_tag"),
			BreakdownByTag:                  stringAttr(attrs, "breakdown_by_tag"),
			ViewMode:                        stringAttr(attrs, "view_mode"),
			IncludeResourcesWithoutGroupTag: boolAttr(attrs, "include_resources_without_group_tag"),
			IncludeOtherBreakdownTagValues:  boolAttr(attrs, "include_other_breakdown_tag_values"),
		}

		block.OtherTagsToExport = stringListAttr(ctx, attrs, "other_tags_to_export")
		block.GroupByTagOrder = groupByTagOrderFromAttr(ctx, attrs["group_by_tag_order"])
		block.BreakdownTagsValues = breakdownTagValuesFromAttr(ctx, attrs["breakdown_tags_values"])
		block.ResourcesBreakdown = resourcesBreakdownFromAttr(ctx, attrs["resources_breakdown"])

		blocks[i] = block
	}
	return blocks, nil
}

// --- FromAPI conversion ---

func BlocksListFromAPI(apiBlocks []customattribute.BlockJson) (types.List, diag.Diagnostics) {
	if apiBlocks == nil {
		return types.ListNull(BlocksListObjectType), nil
	}

	objects := make([]attr.Value, len(apiBlocks))
	for i, b := range apiBlocks {
		obj, diags := blockToTFObject(b)
		if diags.HasError() {
			return types.ListNull(BlocksListObjectType), diags
		}
		objects[i] = obj
	}

	result, diags := types.ListValue(BlocksListObjectType, objects)
	if diags.HasError() {
		return types.ListNull(BlocksListObjectType), diags
	}
	return result, nil
}

func blockToTFObject(b customattribute.BlockJson) (types.Object, diag.Diagnostics) {
	ctx := context.Background()

	otherTagsVal, diags := types.ListValueFrom(ctx, types.StringType, b.OtherTagsToExport)
	if diags.HasError() {
		return types.ObjectNull(BlocksListAttrTypes), diags
	}

	gboVal, diags := groupByTagOrderToTF(b.GroupByTagOrder)
	if diags.HasError() {
		return types.ObjectNull(BlocksListAttrTypes), diags
	}

	btvVal, diags := breakdownTagValuesToTF(b.BreakdownTagsValues)
	if diags.HasError() {
		return types.ObjectNull(BlocksListAttrTypes), diags
	}

	rbVal, diags := resourcesBreakdownToTF(b.ResourcesBreakdown)
	if diags.HasError() {
		return types.ObjectNull(BlocksListAttrTypes), diags
	}

	return types.ObjectValue(BlocksListAttrTypes, map[string]attr.Value{
		"title":                               types.StringValue(b.Title),
		"resource_query":                      types.StringValue(b.ResourceQuery),
		"group_by_tag":                        types.StringValue(b.GroupByTag),
		"breakdown_by_tag":                    types.StringValue(b.BreakdownByTag),
		"view_mode":                           types.StringValue(b.ViewMode),
		"include_resources_without_group_tag": types.BoolValue(b.IncludeResourcesWithoutGroupTag),
		"include_other_breakdown_tag_values":  types.BoolValue(b.IncludeOtherBreakdownTagValues),
		"other_tags_to_export":                otherTagsVal,
		"group_by_tag_order":                  gboVal,
		"breakdown_tags_values":               btvVal,
		"resources_breakdown":                 rbVal,
	})
}

// --- Nested type helpers: GroupByTagOrder ---

func groupByTagOrderToTF(gbo customattribute.GroupByTagOrder) (types.Object, diag.Diagnostics) {
	vals, diags := types.ListValueFrom(context.Background(), types.StringType, gbo.Values)
	if diags.HasError() {
		return types.ObjectNull(GroupByTagOrderAttrTypes), diags
	}
	return types.ObjectValue(GroupByTagOrderAttrTypes, map[string]attr.Value{
		"type":   types.StringValue(gbo.Type),
		"values": vals,
	})
}

func groupByTagOrderFromAttr(ctx context.Context, val attr.Value) customattribute.GroupByTagOrder {
	obj, ok := val.(types.Object)
	if !ok || obj.IsNull() || obj.IsUnknown() {
		return customattribute.GroupByTagOrder{Type: "DEFAULT", Values: []string{}}
	}
	attrs := obj.Attributes()
	return customattribute.GroupByTagOrder{
		Type:   stringAttr(attrs, "type"),
		Values: stringListAttr(ctx, attrs, "values"),
	}
}

// --- Nested type helpers: BreakdownTagsValues ---

func breakdownTagValuesToTF(btvs []customattribute.BreakdownTagValue) (types.List, diag.Diagnostics) {
	if btvs == nil {
		return types.ListValueMust(BreakdownTagValueObjectType, []attr.Value{}), nil
	}
	objects := make([]attr.Value, len(btvs))
	for i, btv := range btvs {
		vals, diags := types.ListValueFrom(context.Background(), types.StringType, btv.Values)
		if diags.HasError() {
			return types.ListNull(BreakdownTagValueObjectType), diags
		}
		obj, diags := types.ObjectValue(BreakdownTagValueAttrTypes, map[string]attr.Value{
			"color":  types.StringValue(btv.Color),
			"label":  types.StringValue(btv.Label),
			"values": vals,
		})
		if diags.HasError() {
			return types.ListNull(BreakdownTagValueObjectType), diags
		}
		objects[i] = obj
	}
	return types.ListValue(BreakdownTagValueObjectType, objects)
}

func breakdownTagValuesFromAttr(ctx context.Context, val attr.Value) []customattribute.BreakdownTagValue {
	list, ok := val.(types.List)
	if !ok || list.IsNull() || list.IsUnknown() {
		return []customattribute.BreakdownTagValue{}
	}
	result := make([]customattribute.BreakdownTagValue, len(list.Elements()))
	for i, elem := range list.Elements() {
		obj := elem.(types.Object)
		attrs := obj.Attributes()
		result[i] = customattribute.BreakdownTagValue{
			Color:  stringAttr(attrs, "color"),
			Label:  stringAttr(attrs, "label"),
			Values: stringListAttr(ctx, attrs, "values"),
		}
	}
	return result
}

// --- Nested type helpers: ResourcesBreakdown ---

func resourcesBreakdownToTF(rbs []customattribute.ResourcesBreakdown) (types.List, diag.Diagnostics) {
	if rbs == nil {
		return types.ListValueMust(ResourcesBreakdownObjectType, []attr.Value{}), nil
	}
	objects := make([]attr.Value, len(rbs))
	for i, rb := range rbs {
		bvObjs := make([]attr.Value, len(rb.BreakdownValues))
		for j, bv := range rb.BreakdownValues {
			obj, diags := types.ObjectValue(BreakdownValueAttrTypes, map[string]attr.Value{
				"value": types.StringValue(bv.Value),
				"count": types.Int64Value(int64(bv.Count)),
			})
			if diags.HasError() {
				return types.ListNull(ResourcesBreakdownObjectType), diags
			}
			bvObjs[j] = obj
		}
		bvList, diags := types.ListValue(BreakdownValueObjectType, bvObjs)
		if diags.HasError() {
			return types.ListNull(ResourcesBreakdownObjectType), diags
		}
		obj, diags := types.ObjectValue(ResourcesBreakdownAttrTypes, map[string]attr.Value{
			"group_by_value":   types.StringValue(rb.GroupByValue),
			"breakdown_values": bvList,
		})
		if diags.HasError() {
			return types.ListNull(ResourcesBreakdownObjectType), diags
		}
		objects[i] = obj
	}
	return types.ListValue(ResourcesBreakdownObjectType, objects)
}

func resourcesBreakdownFromAttr(ctx context.Context, val attr.Value) []customattribute.ResourcesBreakdown {
	list, ok := val.(types.List)
	if !ok || list.IsNull() || list.IsUnknown() {
		return []customattribute.ResourcesBreakdown{}
	}
	result := make([]customattribute.ResourcesBreakdown, len(list.Elements()))
	for i, elem := range list.Elements() {
		obj := elem.(types.Object)
		attrs := obj.Attributes()
		bvList, _ := attrs["breakdown_values"].(types.List)
		bvs := make([]customattribute.BreakdownValue, len(bvList.Elements()))
		for j, bvElem := range bvList.Elements() {
			bvObj := bvElem.(types.Object)
			bvAttrs := bvObj.Attributes()
			bvs[j] = customattribute.BreakdownValue{
				Value: stringAttr(bvAttrs, "value"),
				Count: int(int64Attr(bvAttrs, "count")),
			}
		}
		result[i] = customattribute.ResourcesBreakdown{
			GroupByValue:    stringAttr(attrs, "group_by_value"),
			BreakdownValues: bvs,
		}
	}
	return result
}

// --- Primitive helpers ---

func stringAttr(attrs map[string]attr.Value, key string) string {
	if v, ok := attrs[key].(types.String); ok && !v.IsNull() && !v.IsUnknown() {
		return v.ValueString()
	}
	return ""
}

func boolAttr(attrs map[string]attr.Value, key string) bool {
	if v, ok := attrs[key].(types.Bool); ok && !v.IsNull() && !v.IsUnknown() {
		return v.ValueBool()
	}
	return false
}

func int64Attr(attrs map[string]attr.Value, key string) int64 {
	if v, ok := attrs[key].(types.Int64); ok && !v.IsNull() && !v.IsUnknown() {
		return v.ValueInt64()
	}
	return 0
}

func stringListAttr(ctx context.Context, attrs map[string]attr.Value, key string) []string {
	if v, ok := attrs[key].(types.List); ok && !v.IsNull() && !v.IsUnknown() {
		var result []string
		v.ElementsAs(ctx, &result, false)
		return result
	}
	return []string{}
}

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

package utils

import (
	"terraform/terraform-provider/provider/common"
	"terraform/terraform-provider/provider/common/log"

	"github.com/hashicorp/terraform-plugin-framework/attr"
	"github.com/hashicorp/terraform-plugin-framework/types"
)

func GetStringOrEmpty(requestContext *common.RequestContext, integrationData map[string]interface{}, key string) types.String {
	if value, ok := integrationData[key]; ok {
		return types.StringValue(value.(string))
	} else {
		return types.StringValue("")
	}
}

func GetInt64OrZero(requestContext *common.RequestContext, integrationData map[string]interface{}, key string) types.Int64 {
	if value, ok := integrationData[key]; ok {
		// Handle various numeric types that may come from JSON or other sources
		switch v := value.(type) {
		case int64:
			return types.Int64Value(v)
		case int:
			return types.Int64Value(int64(v))
		case int32:
			return types.Int64Value(int64(v))
		case int16:
			return types.Int64Value(int64(v))
		case int8:
			return types.Int64Value(int64(v))
		case float64:
			return types.Int64Value(int64(v))
		case float32:
			return types.Int64Value(int64(v))
		default:
			return types.Int64Value(0)
		}
	} else {
		return types.Int64Value(0)
	}
}

func GetStringListOrEmpty(requestContext *common.RequestContext, integrationData map[string]interface{}, key string) types.List {
	if value, ok := integrationData[key]; ok {
		return StringListFromMap(requestContext, value)
	} else {
		return types.ListValueMust(types.StringType, []attr.Value{})
	}
}

func StringListFromMap(requestContext *common.RequestContext, value interface{}) types.List {

	interfaceList, ok := value.([]interface{})
	if !ok {
		log.LogError(requestContext, "failed to handle string list (not a list): ", map[string]interface{}{"string_list": value})

		return types.ListNull(types.StringType)
	}

	stringList := make([]string, len(interfaceList))
	for index, item := range interfaceList {
		if itemStr, ok := item.(string); ok {
			stringList[index] = itemStr

		} else {
			log.LogError(requestContext, "failed to handle string list (not a list of strings): ", map[string]interface{}{"string_list": value, "item": item})
			return types.ListNull(types.StringType)
		}
	}

	attrList := make([]attr.Value, len(stringList))
	for i, path := range stringList {
		attrList[i] = types.StringValue(path)
	}

	return types.ListValueMust(types.StringType, attrList)
}

func StringListTFModel(requestContext *common.RequestContext, value types.List) []string {

	stringList := make([]string, len(value.Elements()))
	for i, path := range value.Elements() {
		stringList[i] = path.(types.String).ValueString()
	}

	return stringList
}

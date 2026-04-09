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

package nulls

import (
	"context"
	"terraform/terraform-provider/provider/common"

	"github.com/hashicorp/terraform-plugin-framework/resource/schema/planmodifier"
	"github.com/hashicorp/terraform-plugin-framework/types"
)

func NullListIfUnknownModifier() planmodifier.List {
	return nullListIfUnknownModifier{}
}

type nullListIfUnknownModifier struct {
}

func (m nullListIfUnknownModifier) Description(ctx context.Context) string {
	return m.MarkdownDescription(ctx)
}

func (m nullListIfUnknownModifier) MarkdownDescription(_ context.Context) string {
	return "Set the value to null if it is not known"
}

func (m nullListIfUnknownModifier) PlanModifyList(ctx context.Context, req planmodifier.ListRequest, resp *planmodifier.ListResponse) {

	// If the config value are known, do nothing.
	if common.IsAttrKnown(req.ConfigValue) {
		return
	}

	resp.PlanValue = types.ListNull(req.PlanValue.ElementType(ctx))
}

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

package privatestate

import (
	"context"
	"terraform/terraform-provider/provider/common"

	"terraform/terraform-provider/provider/common/log"

	"github.com/hashicorp/terraform-plugin-framework/diag"
)

// Cannot import internal package privatestate, so we need to define an interface here
type PrivateState interface {
	GetKey(ctx context.Context, key string) ([]byte, diag.Diagnostics)
	SetKey(ctx context.Context, key string, value []byte) diag.Diagnostics
}

func GetKey(ctx *common.RequestContext, privateState PrivateState, key string) string {
	value, diags := privateState.GetKey(ctx.Context, key)
	if diags.HasError() {
		log.LogError(ctx, "failed to get key from private", map[string]interface{}{
			"key":   key,
			"error": diags.Errors(),
		})
	}

	return string(value)
}

func SetKey(ctx *common.RequestContext, privateState PrivateState, key string, value string) {
	diags := privateState.SetKey(ctx.Context, key, []byte(value))
	if diags.HasError() {
		log.LogError(ctx, "failed to set key in private", map[string]interface{}{
			"key":   key,
			"error": diags.Errors(),
		})
	}
}

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

package corecommon

import (
	"fmt"
	"terraform/terraform-provider/provider/common"
	"terraform/terraform-provider/provider/tf/core/model"
)

// ExtractFromTfSource extracts a Terraform model from a source (plan, config, or state)
// using the provided getter. Returns an error if the extraction fails.
func ExtractFromTfSource[TF model.TFModel](requestContext *common.RequestContext, getter Getter, tfModel TF) (TF, error) {
	diags := getter.Get(requestContext.Context, tfModel)

	if diags.HasError() {
		// If there is an error, return the zero value of TF
		var nilTF TF // if TF is a pointer type, it is initialized to nil
		return nilTF, fmt.Errorf("failed to get data from TF source: %s", diags.Errors())
	}

	return tfModel, nil
}

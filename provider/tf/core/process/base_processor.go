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
	"fmt"
	"terraform/terraform-provider/provider/common"
	model "terraform/terraform-provider/provider/tf/core/model"
)

// BasePreProcessor provides common preprocessing functionality to reduce code duplication
type BasePreProcessor[TF model.TFModel] struct{}

// ExtractFromPlan extracts the TF model from plan data (used for Create/Update operations)
func (b *BasePreProcessor[TF]) ExtractFrom(requestContext *common.RequestContext, getter Getter, tfModel TF) (TF, error) {
	diags := getter.Get(requestContext.Context, tfModel)

	if diags.HasError() {
		// If there is an error, return the zero value of TF (which is nil)
		var nilTF TF // TF is a pointer type, so it is initialized to nil
		return nilTF, fmt.Errorf("failed to get data from TF source: %s", diags.Errors())
	}

	return tfModel, nil
}

func (b *BasePreProcessor[TF]) ExtractForDelete(requestContext *common.RequestContext, data *ProcessData, tfModel TF) (TF, error) {

	if data.DeleteRequest != nil {
		return b.ExtractFrom(requestContext, data.DeleteRequest.State, tfModel)
	}

	if data.CreateRequest != nil {
		return b.ExtractFrom(requestContext, data.CreateRequest.Config, tfModel)
	}

	if data.UpdateRequest != nil {
		return b.ExtractFrom(requestContext, data.UpdateRequest.Config, tfModel)
	}

	var nilTf TF
	return nilTf, fmt.Errorf("no data to extract for delete operation")
}

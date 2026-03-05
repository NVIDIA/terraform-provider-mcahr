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
	corecommon "terraform/terraform-provider/provider/tf/core/common"
	model "terraform/terraform-provider/provider/tf/core/model"
)

// BasePreProcessor provides common preprocessing functionality to reduce code duplication
type BasePreProcessor[TF model.TFModel] struct{}

// ExtractFromPlan extracts the TF model from plan data (used for Create/Update operations)
func (b *BasePreProcessor[TF]) ExtractFrom(requestContext *common.RequestContext, getter corecommon.Getter, tfModel TF) (TF, error) {
	return corecommon.ExtractFromTfSource(requestContext, getter, tfModel)
}

func (b *BasePreProcessor[TF]) ExtractForDelete(requestContext *common.RequestContext, data *ProcessData, tfModel TF) (TF, error) {

	if data.DeleteRequest != nil {
		return b.ExtractFrom(requestContext, data.DeleteRequest.State, tfModel)
	}

	// In case of a failure after API call during create, the cleanups will call delete to remove the resource from the remote platform
	// In this case, we need to extract the resource from the create request config
	if data.CreateRequest != nil {
		return b.ExtractFrom(requestContext, data.CreateRequest.Config, tfModel)
	}

	// In case of a failure after API call during update, the cleanups might call the delete flow
	// In this case, we need to extract the resource from the update request config
	if data.UpdateRequest != nil {
		return b.ExtractFrom(requestContext, data.UpdateRequest.Config, tfModel)
	}

	var nilTf TF
	return nilTf, fmt.Errorf("no data to extract for delete operation")
}

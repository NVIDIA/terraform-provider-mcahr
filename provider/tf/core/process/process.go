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
	"terraform/terraform-provider/provider/common"
	systemdefer "terraform/terraform-provider/provider/common/systemdefer"
	"terraform/terraform-provider/provider/external_api/client"
	model "terraform/terraform-provider/provider/tf/core/model"

	"github.com/hashicorp/terraform-plugin-framework/resource"
)

// ProcessData contains the data for the processing flow.
// It can also be used to pass data between the preprocessor and the postprocessor.
type ProcessData struct {
	Client *client.PlatformClient

	DeferFunctionList *systemdefer.DeferFunctionList

	StringArgs map[string]string

	CreateRequest  *resource.CreateRequest
	CreateResponse *resource.CreateResponse

	ReadRequest  *resource.ReadRequest
	ReadResponse *resource.ReadResponse

	UpdateRequest  *resource.UpdateRequest
	UpdateResponse *resource.UpdateResponse

	DeleteRequest  *resource.DeleteRequest
	DeleteResponse *resource.DeleteResponse
}

type PreProcessor[TF model.TFModel] interface {
	PreProcessCreate(*common.RequestContext, *ProcessData) (TF, error)
	PreProcessRead(*common.RequestContext, *ProcessData) (TF, error)
	PreProcessUpdate(*common.RequestContext, *ProcessData) (TF, error)
	PreProcessDelete(*common.RequestContext, *ProcessData) (TF, error)
}

type PostProcessor[TF model.TFModel] interface {
	PostProcessCreate(*common.RequestContext, *ProcessData, TF) error
	PostProcessRead(*common.RequestContext, *ProcessData, TF) error
	PostProcessUpdate(*common.RequestContext, *ProcessData, TF) error
	PostProcessDelete(*common.RequestContext, *ProcessData, TF) error
}

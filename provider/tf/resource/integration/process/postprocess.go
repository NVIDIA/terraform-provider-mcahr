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
	"terraform/terraform-provider/provider/tf/core/process"
	integrationtf "terraform/terraform-provider/provider/tf/resource/integration/model"
)

type IntegrationPostProcessor struct{}

var _ process.PostProcessor[*integrationtf.IntegrationTFModel] = &IntegrationPostProcessor{}

func (p *IntegrationPostProcessor) PostProcessCreate(requestContext *common.RequestContext, data *process.ProcessData, tfModel *integrationtf.IntegrationTFModel) error {
	return setDeprecatedCacheTTL(requestContext, data.CreateRequest.Config, tfModel)
}

func (p *IntegrationPostProcessor) PostProcessRead(requestContext *common.RequestContext, data *process.ProcessData, tfModel *integrationtf.IntegrationTFModel) error {
	return setDeprecatedCacheTTL(requestContext, data.ReadRequest.State, tfModel)
}

func (p *IntegrationPostProcessor) PostProcessUpdate(requestContext *common.RequestContext, data *process.ProcessData, tfModel *integrationtf.IntegrationTFModel) error {
	return setDeprecatedCacheTTL(requestContext, data.UpdateRequest.Config, tfModel)
}

func (p *IntegrationPostProcessor) PostProcessDelete(requestContext *common.RequestContext, data *process.ProcessData, tfModel *integrationtf.IntegrationTFModel) error {
	return nil
}

func setDeprecatedCacheTTL(requestContext *common.RequestContext, config corecommon.Getter, tfModel *integrationtf.IntegrationTFModel) error {

	// Get the original values from config
	var configModel integrationtf.IntegrationTFModel
	diags := config.Get(requestContext.Context, &configModel)
	if diags.HasError() {
		return fmt.Errorf("failed to get original config model: %s", diags.Errors())
	}

	if common.IsAttrKnown(configModel.CacheTTL) {
		// Set it to the value returned by the API (the adapters are setting it in TF model)
		tfModel.CacheTTL = tfModel.CacheTTLMs
	}

	return nil
}

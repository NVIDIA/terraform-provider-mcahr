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

package adapter

import (
	"encoding/json"
	"fmt"
	"terraform/terraform-provider/provider/common"
	"terraform/terraform-provider/provider/tf/resource/integration/adapter/integrations"
	adapterinterface "terraform/terraform-provider/provider/tf/resource/integration/adapter/interface"

	integrationtf "terraform/terraform-provider/provider/tf/resource/integration/model"
)

func GetIntegrationDataAdapter(serviceName string) adapterinterface.IntegrationDataAdapter {
	switch serviceName {
	case "alertmanager":
		return &integrations.AlertmanagerDataAdapter{}
	case "azure_active_directory":
		return &integrations.AzureActiveDirectoryDataAdapter{}
	case "okta":
		return &integrations.OktaDataAdapter{}
	case "google_cloud_identity":
		return &integrations.GoogleCloudIdentityDataAdapter{}
	case "bcm":
		return &integrations.BcmDataAdapter{}
	case "bcm_connectivity":
		return &integrations.BcmConnectivityDataAdapter{}
	case "nvault":
		return &integrations.NvaultDataAdapter{}
	default:
		return nil
	}
}

func TFDataToMap(requestContext *common.RequestContext, adapterOptions *adapterinterface.IntegrationDataAdapterOptions, tfModel *integrationtf.IntegrationTFModel) (map[string]interface{}, error) {
	adapter := GetIntegrationDataAdapter(tfModel.ServiceName.ValueString())
	if adapter == nil {
		return nil, fmt.Errorf("unsupported integration service name: %s", tfModel.ServiceName.ValueString())
	}

	return adapter.TFModelToMap(requestContext, adapterOptions, tfModel), nil
}

func MapToTFData(requestContext *common.RequestContext, adapterOptions *adapterinterface.IntegrationDataAdapterOptions, integrationData map[string]interface{}, tfModel *integrationtf.IntegrationTFModel) error {
	adapter := GetIntegrationDataAdapter(tfModel.ServiceName.ValueString())
	if adapter == nil {
		return fmt.Errorf("unsupported integration service name: %s", tfModel.ServiceName.ValueString())
	}

	adapter.MapToTFModel(requestContext, adapterOptions, integrationData, tfModel)
	return nil
}

func TFDataToJSON(requestContext *common.RequestContext, adapterOptions *adapterinterface.IntegrationDataAdapterOptions, tfModel *integrationtf.IntegrationTFModel) (string, error) {
	params, err := TFDataToMap(requestContext, adapterOptions, tfModel)
	if err != nil {
		return "", err
	}

	jsonParams, err := json.Marshal(params)
	if err != nil {
		return "", err
	}
	return string(jsonParams), nil
}

func JSONToTFData(requestContext *common.RequestContext, adapterOptions *adapterinterface.IntegrationDataAdapterOptions, integrationData string, tfModel *integrationtf.IntegrationTFModel) error {

	var integrationDataMap map[string]interface{}
	err := json.Unmarshal([]byte(integrationData), &integrationDataMap)
	if err != nil {
		return err
	}
	return MapToTFData(requestContext, adapterOptions, integrationDataMap, tfModel)
}

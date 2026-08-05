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

package translator

import (
	"bytes"
	"encoding/json"
	"fmt"
	"strings"

	"terraform/terraform-provider/provider/common"
	"terraform/terraform-provider/provider/external_api/resources/statement"
	"terraform/terraform-provider/provider/tf/core/translator"
	utils "terraform/terraform-provider/provider/tf/core/translator"
	secrettf "terraform/terraform-provider/provider/tf/resource/secret/model"
)

// NVaultSecretTranslatorCommon provides common functionality for nvault secret translators across API versions
type NVaultSecretTranslatorCommon struct{}

// nvaultExternalValue is the JSON object sent as the secret's external_value.
// Field order matches the previously hand-built string so the wire format is
// unchanged.
type nvaultExternalValue struct {
	IntegrationName string `json:"integration_name"`
	VaultSecretPath string `json:"vault_secret_path"`
	VaultSecretKey  string `json:"vault_secret_key"`
}

// encodeExternalValue marshals the external_value object with HTML escaping
// disabled. encoding/json escapes <, > and & to <, > and & by
// default, which would alter the bytes sent for vault paths containing those
// characters; the previous hand-built string did not. Escaping of quotes and
// backslashes -- the part that matters here -- is unaffected.
func encodeExternalValue(value nvaultExternalValue) (string, error) {
	var buf bytes.Buffer
	encoder := json.NewEncoder(&buf)
	encoder.SetEscapeHTML(false)

	if err := encoder.Encode(value); err != nil {
		return "", err
	}

	// Encode terminates the object with a newline.
	return strings.TrimRight(buf.String(), "\n"), nil
}

// ToAPIModelWithVersion converts a TF model to an API model with specified backend version
func (t *NVaultSecretTranslatorCommon) ToAPIModelWithVersion(requestContext *common.RequestContext, translationData *translator.TranslationData, tfModel *secrettf.NVaultSecretTFModel) (*statement.InputAPIModel, error) {
	var stmt string
	var err error

	switch requestContext.Operation {
	case common.Create:
		stmt, err = t.buildCreateStatement(requestContext, translationData, tfModel)
	case common.Read:
		stmt = t.buildReadStatement(tfModel)
	case common.Update:
		stmt, err = t.buildUpdateStatement(requestContext, translationData, tfModel)
	case common.Delete:
		stmt = t.buildDeleteStatement(tfModel)
	default:
		return nil, fmt.Errorf("unsupported operation: %v", requestContext.Operation)
	}
	if err != nil {
		return nil, err
	}

	apiModel := &statement.InputAPIModel{
		Statement:  stmt,
		APIVersion: requestContext.APIVersion,
	}

	return apiModel, nil
}

func (t *NVaultSecretTranslatorCommon) buildCreateStatement(requestContext *common.RequestContext, translationData *translator.TranslationData, tfModel *secrettf.NVaultSecretTFModel) (string, error) {
	return t.buildSecretStatement(requestContext, translationData, "define_secret", tfModel)
}

func (t *NVaultSecretTranslatorCommon) buildReadStatement(tfModel *secrettf.NVaultSecretTFModel) string {
	name := tfModel.Name.ValueString()
	return fmt.Sprintf("get_secret(secret_name=%s)", utils.EscapeString(name))
}

func (t *NVaultSecretTranslatorCommon) buildUpdateStatement(requestContext *common.RequestContext, translationData *translator.TranslationData, tfModel *secrettf.NVaultSecretTFModel) (string, error) {
	return t.buildSecretStatement(requestContext, translationData, "update_secret", tfModel)
}

func (t *NVaultSecretTranslatorCommon) buildDeleteStatement(tfModel *secrettf.NVaultSecretTFModel) string {
	name := tfModel.Name.ValueString()
	return fmt.Sprintf("delete_secret(secret_name=%s)", utils.EscapeString(name))
}

func (t *NVaultSecretTranslatorCommon) buildSecretStatement(requestContext *common.RequestContext, translationData *translator.TranslationData, statementName string, tfModel *secrettf.NVaultSecretTFModel) (string, error) {
	builder := utils.NewStatementBuilder(statementName, requestContext.BackendVersion, translationData.CompatibilityOptions)

	// Add required fields
	name := tfModel.Name.ValueString()
	builder.SetStringField("secret_name", name, "name")

	// Build external_value as a JSON object for the secret configuration.
	//
	// This must be marshalled, not concatenated. vault_secret_path and
	// vault_secret_key carry no schema validator, so a quote character in
	// either one used to close the JSON string and then the surrounding
	// op-lang literal, letting HCL append statement syntax the backend would
	// execute. json.Marshal escapes the values while preserving the exact
	// wire format (an unquoted JSON object), which SetField passes through
	// verbatim.
	externalValue, err := encodeExternalValue(nvaultExternalValue{
		IntegrationName: tfModel.IntegrationName.ValueString(),
		VaultSecretPath: tfModel.VaultSecretPath.ValueString(),
		VaultSecretKey:  tfModel.VaultSecretKey.ValueString(),
	})
	if err != nil {
		return "", fmt.Errorf("failed to encode secret external_value: %w", err)
	}

	builder.SetField("external_value", externalValue, "")

	return builder.Build(), nil
}

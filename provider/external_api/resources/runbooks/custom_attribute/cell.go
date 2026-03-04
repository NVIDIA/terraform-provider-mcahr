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

package customattribute

import (
	"encoding/json"
	"fmt"
	"terraform/terraform-provider/provider/common"
	commonstruct "terraform/terraform-provider/provider/common/struct"
)

const (
	OP_LANG_TYPE  = "OP_LANG"
	MARKDOWN_TYPE = "MARKDOWN"
)

var (
	DefaultCellContent     = ""
	DefaultCellType        = "OP_LANG"
	DefaultCellName        = "unnamed"
	DefaultCellEnabled     = true
	DefaultCellSecretAware = false
	DefaultCellDescription = ""
)

// CellJsonAPI is the API model for a cell
type CellJsonAPI struct {
	// Two type fields because of a backend bug
	Type     string `json:"type"`                // output to API
	CellType string `json:"cell_type,omitempty"` // input from API
	Content  string `json:"content"`

	Name        string `json:"name"`
	Enabled     bool   `json:"enabled"`
	SecretAware bool   `json:"secret_aware" min_version:"release-28.1.0" skip_version_prefixes:"release-28.3"`
	Description string `json:"description" min_version:"release-29.0.1"`
}

func (c *CellJsonAPI) ToInternalModel() *CellJson {

	cell := &CellJson{
		Name:        c.Name,
		Enabled:     c.Enabled,
		SecretAware: c.SecretAware,
		Description: c.Description,
	}

	var typeValue string
	if c.Type != "" {
		typeValue = c.Type
	} else if c.CellType != "" {
		typeValue = c.CellType
	}

	switch typeValue {
	case OP_LANG_TYPE:
		cell.Op = common.NewOptional(c.Content)
	case MARKDOWN_TYPE:
		cell.Md = common.NewOptional(c.Content)
	}

	return cell
}

func (c *CellJsonAPI) SetFromMap(cell map[string]interface{}) {

	// Apply defaults for fields that should have them
	c.Name = DefaultCellName
	c.Enabled = DefaultCellEnabled
	c.SecretAware = DefaultCellSecretAware
	c.Description = DefaultCellDescription

	// Override with values from map if present
	if name, ok := cell["name"]; ok {
		c.Name = name.(string)
	}

	if content, ok := cell["content"]; ok {
		c.Content = content.(string)
	}

	if enabled, ok := cell["enabled"]; ok {
		c.Enabled = enabled.(bool)
	}

	if cellType, ok := cell["type"]; ok {
		c.CellType = cellType.(string)
	} else if cellType, ok := cell["cell_type"]; ok {
		c.CellType = cellType.(string)
	}

	if secretAware, ok := cell["secret_aware"]; ok {
		c.SecretAware = secretAware.(bool)
	}

	if description, ok := cell["description"]; ok {
		c.Description = description.(string)
	}
}

// CellJson is the internal model for a cell
type CellJson struct {
	// Embedded config to be used in marshal/unmarshal functions
	Config common.JsonConfig `json:"-" skip:"true"`

	Op common.Optional[string] `json:"op,omitempty"`
	Md common.Optional[string] `json:"md,omitempty"`

	Name        string `json:"name"`
	Enabled     bool   `json:"enabled"`
	SecretAware bool   `json:"secret_aware"`
	Description string `json:"description" min_version:"release-29.0.1"`
}

var _ common.JsonConfigurable = &CellJson{}
var _ json.Marshaler = &CellJson{}
var _ json.Unmarshaler = &CellJson{}

func (c *CellJson) SetConfig(config common.JsonConfig) {
	c.Config = config
}

func (c *CellJson) GetConfig() common.JsonConfig {
	return c.Config
}

func (c *CellJson) ToAPIModel() *CellJsonAPI {
	apiModel := &CellJsonAPI{
		Name:        c.Name,
		Enabled:     c.Enabled,
		SecretAware: c.SecretAware,
		Description: c.Description,
	}

	if c.Op.IsSet {
		apiModel.Content = c.Op.Get()
		apiModel.Type = OP_LANG_TYPE
	} else if c.Md.IsSet {
		apiModel.Content = c.Md.Get()
		apiModel.Type = MARKDOWN_TYPE
	} else {
		apiModel.Content = DefaultCellContent
		apiModel.Type = DefaultCellType
	}

	return apiModel
}

func (c *CellJson) MarshalJSON() ([]byte, error) {

	if err := validateOpAndMd(c.Op, c.Md); err != nil {
		return nil, err
	}

	options := map[string]any{"backend_version": c.Config.BackendVersion}
	result := commonstruct.ApplyCustomStructTags(*c, options)

	// Remove unset optional fields to respect omitempty behavior
	if !c.Op.IsSet {
		delete(result, "op")
	}
	if !c.Md.IsSet {
		delete(result, "md")
	}

	return json.Marshal(result)
}

func (c *CellJson) UnmarshalJSON(b []byte) error {

	type Alias CellJson // Prevent recursion
	aux := &Alias{
		Config: common.JsonConfig{
			BackendVersion: c.Config.BackendVersion,
		},
		Op:          common.NewOptionalUnset[string](),
		Md:          common.NewOptionalUnset[string](),
		Name:        DefaultCellName,
		Enabled:     DefaultCellEnabled,
		SecretAware: DefaultCellSecretAware,
		Description: DefaultCellDescription,
	} // defaults

	// Unmarshal and replace the default values with the ones from the JSON
	// The Optional[string] fields will have their UnmarshalJSON called automatically
	if err := json.Unmarshal(b, aux); err != nil {
		return err
	}

	// Validate Op and Md before applying custom struct tags
	if err := validateOpAndMd(aux.Op, aux.Md); err != nil {
		return err
	}

	// Apply custom struct tags to the unmarshalled data
	options := map[string]any{"backend_version": c.Config.BackendVersion}
	result := commonstruct.ApplyCustomStructTags(*aux, options)

	// Set in the final struct the values from the processed data
	if name, ok := result["name"]; ok {
		c.Name = name.(string)
	}

	// Use aux values directly since they're already properly typed
	c.Op = aux.Op
	c.Md = aux.Md

	if enabled, ok := result["enabled"]; ok {
		c.Enabled = enabled.(bool)
	}

	if description, ok := result["description"]; ok {
		c.Description = description.(string)
	}

	if secretAware, ok := result["secret_aware"]; ok {
		c.SecretAware = secretAware.(bool)
	}

	return nil
}

func validateOpAndMd(op common.Optional[string], md common.Optional[string]) error {

	hasOp := op.IsSet
	hasMd := md.IsSet

	if hasOp && hasMd {
		return fmt.Errorf("runbook cell cannot have both op and md")
	}
	if !hasOp && !hasMd {
		return fmt.Errorf("runbook cell must have either op or md set")
	}

	return nil
}

func MapCellsToAPIModel(encodedCells string) (string, error) {

	var cells []CellJson
	err := json.Unmarshal([]byte(encodedCells), &cells)
	if err != nil {
		return "", err
	}

	apiModels := make([]CellJsonAPI, len(cells))
	for i, cell := range cells {
		apiModels[i] = *cell.ToAPIModel()
	}

	marshaledCells, err := json.Marshal(apiModels)
	if err != nil {
		return "", err
	}

	return common.EncodeBase64(string(marshaledCells)), nil
}

func MapCellsToInternalModel(apiCells []CellJsonAPI) (string, error) {

	internalCells := make([]CellJson, len(apiCells))
	for i, cell := range apiCells {
		internalCells[i] = *cell.ToInternalModel()
	}

	marshaledCells, err := json.Marshal(internalCells)
	if err != nil {
		return "", err
	}

	return string(marshaledCells), nil
}

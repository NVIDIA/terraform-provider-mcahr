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

package model

import (
	model "terraform/terraform-provider/provider/tf/core/model"

	"github.com/hashicorp/terraform-plugin-framework/types"
)

var _ model.TFModel = &FileTFModel{}

// FileTFModel represents the Terraform model for file resources
type FileTFModel struct {
	Name            types.String `tfsdk:"name"`
	DestinationPath types.String `tfsdk:"destination_path"`
	Description     types.String `tfsdk:"description"`
	ResourceQuery   types.String `tfsdk:"resource_query"`
	Enabled         types.Bool   `tfsdk:"enabled"`
	InputFile       types.String `tfsdk:"input_file"`
	InlineData      types.String `tfsdk:"inline_data"`
	FileData        types.String `tfsdk:"file_data"`
	FileLength      types.Int64  `tfsdk:"file_length"`
	Checksum        types.String `tfsdk:"checksum"`
	MD5             types.String `tfsdk:"md5"`
	Mode            types.String `tfsdk:"mode"`
	Owner           types.String `tfsdk:"owner"`
}

// GetName returns the name of the file resource
func (f FileTFModel) GetName() string {
	return f.Name.ValueString()
}

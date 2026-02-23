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
	core "terraform/terraform-provider/provider/tf/core/model"

	"github.com/hashicorp/terraform-plugin-framework/types"
)

var _ core.TFModel = &BotTFModel{}

// BotTFModel represents the Terraform model for bot resources
type BotTFModel struct {
	Name        types.String `tfsdk:"name"`
	Command     types.String `tfsdk:"command"`
	Description types.String `tfsdk:"description"`
	Enabled     types.Bool   `tfsdk:"enabled"`
	Family      types.String `tfsdk:"family"`

	EventType     types.String `tfsdk:"event_type"`
	TriggerSource types.String `tfsdk:"trigger_source"`

	MonitorID types.String `tfsdk:"monitor_id"`
	TriggerID types.String `tfsdk:"trigger_id"`

	AlarmResourceQuery     types.String `tfsdk:"alarm_resource_query"`
	CommunicationWorkspace types.String `tfsdk:"communication_workspace"`
	CommunicationChannel   types.String `tfsdk:"communication_channel"`
	IntegrationName        types.String `tfsdk:"integration_name"`
}

// GetName returns the name of the bot resource
func (b *BotTFModel) GetName() string {
	return b.Name.ValueString()
}

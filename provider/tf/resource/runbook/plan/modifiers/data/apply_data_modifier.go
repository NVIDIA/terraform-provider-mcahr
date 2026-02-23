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

package data

import (
	"context"
	"terraform/terraform-provider/provider/tf/resource/runbook/model"

	data "terraform/terraform-provider/provider/tf/resource/runbook/data_attribute"
)

func ApplyDataModifier(ctx context.Context, config *model.RunbookTFModel) (*model.RunbookTFModel, error) {
	// Create a copy of the config
	modelCopy := config.Copy()

	err := data.ApplyDataJSONValues(ctx, modelCopy)
	if err != nil {
		return nil, err
	}

	return modelCopy, nil
}

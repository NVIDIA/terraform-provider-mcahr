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

package defaults

import (
	"terraform/terraform-provider/provider/tf/resource/integration/adapter"

	"github.com/hashicorp/terraform-plugin-framework/path"
)

func IsServiceNameCompatible(serviceName string, attributeName string) bool {

	integrationAdapter := adapter.GetIntegrationDataAdapter(serviceName)
	if integrationAdapter == nil {
		return false
	}

	for _, fieldName := range integrationAdapter.TFModelFieldNames() {
		if fieldName == attributeName {
			return true
		}
	}
	return false

}

func GetAttributeName(p path.Path) string {
	steps := p.Steps()
	if len(steps) == 0 {
		return ""
	}
	if attr, ok := steps[len(steps)-1].(path.PathStepAttributeName); ok {
		return string(attr)
	}
	return ""
}

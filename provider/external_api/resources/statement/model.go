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

package statement

import "terraform/terraform-provider/provider/common"

type InputAPIModel struct {
	Statement  string            `json:"statement"`
	ApiPayload string            `json:"-"` // this is used when the resource cannot be created via op statements (`define_<resource>`, ...)
	APIVersion common.APIVersion `json:"-"`
}

type StatementResponse struct {
	Output string `json:"output"`
}

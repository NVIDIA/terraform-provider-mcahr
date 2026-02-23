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

package systemdefer

type DeferFunctionList struct {
	deferList []func()
}

func NewDeferFunctionList() *DeferFunctionList {
	return &DeferFunctionList{
		deferList: []func(){},
	}
}

func (d *DeferFunctionList) AddDefer(deferFunc func()) {
	d.deferList = append(d.deferList, deferFunc)
}

func (d *DeferFunctionList) ExecuteAll() {
	for _, deferFunc := range d.deferList {
		deferFunc()
	}
}

// Size returns the number of deferred functions in the list
func (d *DeferFunctionList) Size() int {
	return len(d.deferList)
}

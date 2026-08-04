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

// The setters in this package read values out of map[string]any documents
// decoded from backend JSON. A bare type assertion there panics the whole
// provider process when a field arrives as a number, null or object instead of
// the expected type -- the backend's response shape is an assumption, not a
// guarantee, and a panic in a plugin surfaces to the operator as an opaque
// crash rather than a diagnostic.
//
// These helpers keep the caller's current value when the key is absent *or*
// holds a different type, so a mistyped field degrades to the existing default
// instead of taking the process down.

// stringOrKeep returns the string stored at key, or keep if the key is missing
// or does not hold a string.
func stringOrKeep(data map[string]any, key string, keep string) string {
	if value, ok := data[key]; ok {
		if typed, ok := value.(string); ok {
			return typed
		}
	}
	return keep
}

// boolOrKeep returns the bool stored at key, or keep if the key is missing or
// does not hold a bool.
func boolOrKeep(data map[string]any, key string, keep bool) bool {
	if value, ok := data[key]; ok {
		if typed, ok := value.(bool); ok {
			return typed
		}
	}
	return keep
}

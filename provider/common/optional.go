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

package common

import (
	"bytes"
	"encoding/json"
	"fmt"
)

// Optional represents a value that may or may not be present in JSON.
// It distinguishes between three states:
// - Not present in JSON (IsSet=false)
// - Present but null (IsSet=true, IsNull=true)
// - Present with a value (IsSet=true, IsNull=false)
type Optional[T any] struct {
	IsSet  bool // true if field was present in JSON
	IsNull bool // true if field value is null
	Value  T    // the actual value
}

// NewOptional creates an Optional with a value
func NewOptional[T any](value T) Optional[T] {
	return Optional[T]{
		IsSet:  true,
		IsNull: false,
		Value:  value,
	}
}

// NewOptionalEmpty creates an Optional that's set with the zero value
func NewOptionalEmpty[T any]() Optional[T] {
	var zero T
	return Optional[T]{
		IsSet:  true,
		IsNull: false,
		Value:  zero,
	}
}

// NewOptionalUnset creates an Optional that's not set
func NewOptionalUnset[T any]() Optional[T] {
	var zero T
	return Optional[T]{
		IsSet:  false,
		IsNull: false,
		Value:  zero,
	}
}

// Get returns the value or zero value if null or not set
func (o Optional[T]) Get() T {
	if o.IsSet && !o.IsNull {
		return o.Value
	}
	var zero T
	return zero
}

// IsEmpty returns true if set with empty/zero value (for strings, checks if "")
func (o Optional[T]) IsEmpty() bool {
	if !o.IsSet || o.IsNull {
		return false
	}
	// For strings, check if empty
	if str, ok := any(o.Value).(string); ok {
		return str == ""
	}
	return false
}

// HasValue returns true if set with non-empty value (for strings, checks if non-empty)
func (o Optional[T]) HasValue() bool {
	if !o.IsSet || o.IsNull {
		return false
	}
	// For strings, check if non-empty
	if str, ok := any(o.Value).(string); ok {
		return str != ""
	}
	return true
}

// UnmarshalJSON implements json.Unmarshaler
func (o *Optional[T]) UnmarshalJSON(data []byte) error {
	o.IsSet = true // UnmarshalJSON called means field exists
	o.IsNull = false

	if bytes.Equal(data, []byte("null")) {
		o.IsNull = true
		var zero T
		o.Value = zero
		return nil
	}

	var v T
	if err := json.Unmarshal(data, &v); err != nil {
		return err
	}
	o.Value = v
	return nil
}

// MarshalJSON implements json.Marshaler
func (o Optional[T]) MarshalJSON() ([]byte, error) {
	if o.IsNull {
		return []byte("null"), nil
	}
	return json.Marshal(o.Value)
}

// String returns a string representation for debugging
func (o Optional[T]) String() string {
	if !o.IsSet {
		return "<unset>"
	}
	if o.IsNull {
		return "<null>"
	}
	return fmt.Sprintf("%v", o.Value)
}

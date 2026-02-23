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
	"encoding/base64"
	"fmt"
	"reflect"
	"strings"
	"unicode"

	"github.com/hashicorp/terraform-plugin-framework/attr"
)

type CrudOperation int

const (
	Create CrudOperation = iota
	Read
	Update
	Delete
	Import
)

// String returns the string representation of the CrudOperation
func (op CrudOperation) String() string {
	switch op {
	case Create:
		return "Create"
	case Read:
		return "Read"
	case Update:
		return "Update"
	case Delete:
		return "Delete"
	case Import:
		return "Import"
	default:
		return "Unknown"
	}
}

// APIVersion represents the version of the API output version.
type APIVersion string

const (
	V1 APIVersion = "V1"
	V2 APIVersion = "V2"
)

func EncodeBase64(value string) string {
	return WrapInQuotes(base64.StdEncoding.EncodeToString([]byte(value)))
}

func WrapInQuotes(value string) string {
	return fmt.Sprintf("\"%s\"", value)
}

func IsAttrKnown(val attr.Value) bool {
	return !val.IsNull() && !val.IsUnknown()
}

func SnakeToCamelCase(snake string) string {
	if snake == "" {
		return snake
	}

	var camel []rune
	capitalizeNext := false

	for i, r := range snake {
		if r == '_' {
			capitalizeNext = true
		} else {
			if capitalizeNext && i > 0 {
				camel = append(camel, unicode.ToUpper(r))
				capitalizeNext = false
			} else {
				camel = append(camel, r)
			}
		}
	}

	return string(camel)
}

// CamelToSnakeCase converts camelCase to snake_case
func CamelToSnakeCase(camel string) string {
	if camel == "" {
		return camel
	}

	var result []rune
	for i, r := range camel {
		if i > 0 && r >= 'A' && r <= 'Z' {
			result = append(result, '_')
		}
		result = append(result, r)
	}
	return strings.ToLower(string(result))
}

func IsNil(t any) bool {
	return t == nil || reflect.ValueOf(t).IsNil()
}

func HasErrorOrNil(err error, pointerVar any) bool {
	return err != nil || IsNil(pointerVar)
}

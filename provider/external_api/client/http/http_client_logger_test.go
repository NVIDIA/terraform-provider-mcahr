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

package http_client

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestRedactSensitiveBody(t *testing.T) {
	t.Parallel()

	tests := []struct {
		name        string
		body        string
		mustContain []string
		mustNotHave []string
	}{
		{
			name:        "empty body",
			body:        "",
			mustContain: []string{},
		},
		{
			name:        "integration credentials are masked",
			body:        `{"name":"my_integration","api_key":"AKIA-real-key","client_secret":"hunter2"}`,
			mustContain: []string{`"name":"my_integration"`, redactedPlaceholder},
			mustNotHave: []string{"AKIA-real-key", "hunter2"},
		},
		{
			name:        "presigned storage URL is masked",
			body:        `{"get_file_attribute":"ok","presigned_put":"https://bucket/obj?X-Amz-Signature=deadbeef"}`,
			mustNotHave: []string{"X-Amz-Signature", "deadbeef"},
		},
		{
			name:        "secret external_value is masked",
			body:        `{"external_value":"{\"vault_secret_key\":\"password\"}"}`,
			mustNotHave: []string{"vault_secret_key"},
		},
		{
			name:        "nested credentials are masked",
			body:        `{"data":{"items":[{"smtp_password":"pw","sender":"a@b.c"}]}}`,
			mustContain: []string{"a@b.c", redactedPlaceholder},
			mustNotHave: []string{`"pw"`},
		},
		{
			name:        "camelCase spellings are masked",
			body:        `{"apiKey":"secret-value","idpName":"okta"}`,
			mustContain: []string{"okta"},
			mustNotHave: []string{"secret-value"},
		},
		{
			name:        "non-JSON body is dropped wholesale",
			body:        "api_key=AKIA-real-key&user=bob",
			mustContain: []string{redactedPlaceholder},
			mustNotHave: []string{"AKIA-real-key", "bob"},
		},
		{
			name:        "benign fields survive",
			body:        `{"name":"my_action","enabled":true,"count":3}`,
			mustContain: []string{`"name":"my_action"`, `"enabled":true`},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()

			got := redactSensitiveBody([]byte(tt.body))

			for _, want := range tt.mustContain {
				assert.Contains(t, got, want)
			}
			for _, unwanted := range tt.mustNotHave {
				assert.NotContains(t, got, unwanted)
			}
		})
	}
}

func TestRedactSensitiveBodyEmptyReturnsEmpty(t *testing.T) {
	t.Parallel()

	assert.Equal(t, "", redactSensitiveBody(nil))
	assert.Equal(t, "", redactSensitiveBody([]byte{}))
}

func TestIsSensitiveBodyField(t *testing.T) {
	t.Parallel()

	sensitive := []string{
		"api_key", "apiKey", "API_KEY", "client_secret", "credentials",
		"password", "smtp_password", "token", "external_value",
		"presigned_put", "api_certificate", "vault_secret_key",
	}
	for _, key := range sensitive {
		assert.True(t, isSensitiveBodyField(key), "expected %q to be treated as sensitive", key)
	}

	benign := []string{"name", "enabled", "serial_number", "sender", "idp_name", "description"}
	for _, key := range benign {
		assert.False(t, isSensitiveBodyField(key), "expected %q to be treated as benign", key)
	}
}

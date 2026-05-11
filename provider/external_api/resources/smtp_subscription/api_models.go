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

package smtp_subscription

type SmtpSubscriptionFilter struct {
	Category string `json:"category"`         // Required: Event category
	Type     string `json:"type"`             // Required: Event type
	Status   string `json:"status,omitempty"` // Optional: Event status
}

// SmtpSubscriptionUpdateRequest is sent for create/read/update/delete operations.
// Recipients/Filters/Enabled use omitempty so partial updates only carry changed fields.
type SmtpSubscriptionUpdateRequest struct {
	ID              string                   `json:"id"`
	Name            string                   `json:"name"`
	IntegrationName string                   `json:"integration_name"`
	Recipients      []string                 `json:"recipients,omitempty"`
	Filters         []SmtpSubscriptionFilter `json:"filters,omitempty"`
	Enabled         *bool                    `json:"enabled,omitempty"`
}

// SmtpSubscriptionResponseAPIModel represents the REST API response for SMTP subscriptions
// This is the actual response format from /api/v1/integrations/smtp/subscriptions endpoint
type SmtpSubscriptionResponseAPIModel struct {
	ID              string                   `json:"id"`
	Name            string                   `json:"name"`
	IntegrationName string                   `json:"integration_name"`
	Recipients      []string                 `json:"recipients"`
	Filters         []SmtpSubscriptionFilter `json:"filters"`
	Enabled         bool                     `json:"enabled"`
	CreatedBy       string                   `json:"created_by"`
	UpdatedBy       string                   `json:"updated_by"`
	CreatedTimeMs   int64                    `json:"created_time_ms"`
	UpdatedTimeMs   int64                    `json:"updated_time_ms"`
}

// GetErrors implements the APIModel interface
// REST API errors are typically returned via HTTP status codes, not in the response body
func (s SmtpSubscriptionResponseAPIModel) GetErrors() string {
	return ""
}

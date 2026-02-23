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

package client

import (
	"fmt"
	"net/http"
	"time"
)

// HTTP status code ranges for categorizing responses
const (
	// Success status codes (2xx)
	HTTPStatusSuccessMin = 200
	HTTPStatusSuccessMax = 300

	// Client error status codes (4xx)
	HTTPStatusClientErrorMin = 400
	HTTPStatusClientErrorMax = 500

	// Server error status codes (5xx)
	HTTPStatusServerErrorMin = 500
	HTTPStatusServerErrorMax = 600

	ExponentialBackoffUpperBound = 60 * time.Second
)

// HTTPError represents an HTTP error
type HTTPError struct {
	StatusCode int
	Status     string
}

func (e *HTTPError) Error() string {
	return fmt.Sprintf("HTTP %s", e.Status)
}

// ClientError for 4xx errors
type ClientError struct {
	*HTTPError
}

// ServerError for 5xx errors
type ServerError struct {
	*HTTPError
}

var _ RetryableError = &ServerError{} // check that ServerError implements RetryableError

func (e *ServerError) ShouldRetry(retryAttempt int) bool { return retryAttempt < maxRetries }
func (e *ServerError) RetryAfter(retryAttempt int) time.Duration {
	return exponentialBackoff(retryAttempt)
}

// AuthenticationError for 401 - missing/invalid credentials
type AuthenticationError struct {
	*HTTPError
}

// AuthorizationError for 403 - insufficient permissions
type AuthorizationError struct {
	*HTTPError
}

// RateLimitError for 429
type RateLimitError struct {
	*HTTPError
	RetryDelay time.Duration
}

var _ RetryableError = &RateLimitError{} // check that RateLimitError implements RetryableError

func (e *RateLimitError) ShouldRetry(retryAttempt int) bool { return retryAttempt < maxRetries }
func (e *RateLimitError) RetryAfter(retryAttempt int) time.Duration {
	return e.RetryDelay
}

// TimeoutError for request timeouts
type TimeoutError struct {
	*HTTPError
}

var _ RetryableError = &TimeoutError{} // check that TimeoutError implements RetryableError

func (e *TimeoutError) ShouldRetry(retryAttempt int) bool { return retryAttempt < maxRetries }
func (e *TimeoutError) RetryAfter(retryAttempt int) time.Duration {
	return exponentialBackoff(retryAttempt)
}

// Other functions

func exponentialBackoff(retryAttempt int) time.Duration {
	// Exponential backoff: 1s, 2s, 4s
	backoff := time.Duration(1<<retryAttempt) * time.Second
	if backoff > ExponentialBackoffUpperBound {
		return ExponentialBackoffUpperBound
	}

	return backoff
}

// maybeMapStatusCodeError maps the status code to the appropriate error type
func maybeMapStatusCodeError(resp *http.Response) error {
	statusCode := resp.StatusCode

	// Success status codes (2xx)
	if statusCode >= HTTPStatusSuccessMin && statusCode < HTTPStatusSuccessMax {
		return nil
	}

	httpErr := &HTTPError{
		StatusCode: statusCode,
		Status:     resp.Status,
	}

	// Handle specific status codes
	switch statusCode {
	case http.StatusUnauthorized: // 401
		return &AuthenticationError{HTTPError: httpErr}
	case http.StatusForbidden: // 403
		return &AuthorizationError{HTTPError: httpErr}
	case http.StatusTooManyRequests: // 429
		return &RateLimitError{HTTPError: httpErr, RetryDelay: time.Duration(rateLimitDelaySeconds) * time.Second}
	}

	if statusCode >= HTTPStatusClientErrorMin && statusCode < HTTPStatusClientErrorMax {
		// Client errors (4xx)
		return &ClientError{HTTPError: httpErr}
	}

	if statusCode >= HTTPStatusServerErrorMin && statusCode < HTTPStatusServerErrorMax {
		// Server errors (5xx)
		return &ServerError{HTTPError: httpErr}
	}

	// Other non-success status codes
	return httpErr
}

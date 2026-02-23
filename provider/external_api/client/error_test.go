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
	"testing"
	"time"
)

func TestHTTPError(t *testing.T) {
	t.Parallel()

	// given
	err := &HTTPError{
		StatusCode: 404,
		Status:     "Not Found",
	}

	// when
	result := err.Error()

	// then
	expected := "HTTP Not Found"
	if result != expected {
		t.Errorf("expected error message %q, got %q", expected, result)
	}
}

func TestClientError(t *testing.T) {
	t.Parallel()

	// given
	httpErr := &HTTPError{
		StatusCode: 400,
		Status:     "Bad Request",
	}
	clientErr := &ClientError{
		HTTPError: httpErr,
	}

	// when
	result := clientErr.Error()

	// then
	expected := "HTTP Bad Request"
	if result != expected {
		t.Errorf("expected error message %q, got %q", expected, result)
	}
}

func TestServerErrorShouldRetry(t *testing.T) {
	t.Parallel()

	// given
	httpErr := &HTTPError{
		StatusCode: 500,
		Status:     "Internal Server Error",
	}
	serverErr := &ServerError{
		HTTPError: httpErr,
	}
	testCases := []struct {
		retryAttempt        int
		expectedShouldRetry bool
	}{
		{0, true},
		{1, true},
		{2, true},
		{3, true},
		{4, true},
		{5, false}, // maxRetries is 5 in the code
		{6, false},
	}

	for _, tc := range testCases {
		// when
		result := serverErr.ShouldRetry(tc.retryAttempt)

		// then
		if result != tc.expectedShouldRetry {
			t.Errorf("ShouldRetry(%d) = %v, expected %v", tc.retryAttempt, result, tc.expectedShouldRetry)
		}
	}
}

func TestServerErrorRetryAfter(t *testing.T) {
	t.Parallel()

	// given
	httpErr := &HTTPError{
		StatusCode: 500,
		Status:     "Internal Server Error",
	}
	serverErr := &ServerError{
		HTTPError: httpErr,
	}
	testCases := []struct {
		retryAttempt int
		expected     time.Duration
	}{
		{0, 1 * time.Second},  // 2^0 = 1
		{1, 2 * time.Second},  // 2^1 = 2
		{2, 4 * time.Second},  // 2^2 = 4
		{3, 8 * time.Second},  // 2^3 = 8
		{4, 16 * time.Second}, // 2^4 = 16
	}

	for _, tc := range testCases {
		// when
		result := serverErr.RetryAfter(tc.retryAttempt)

		// then
		if result != tc.expected {
			t.Errorf("RetryAfter(%d) = %v, expected %v", tc.retryAttempt, result, tc.expected)
		}
	}
}

func TestAuthenticationError(t *testing.T) {
	t.Parallel()

	// given
	httpErr := &HTTPError{
		StatusCode: 401,
		Status:     "Unauthorized",
	}
	authErr := &AuthenticationError{
		HTTPError: httpErr,
	}

	// when
	result := authErr.Error()

	// then
	expected := "HTTP Unauthorized"
	if result != expected {
		t.Errorf("expected error message %q, got %q", expected, result)
	}
}

func TestAuthorizationError(t *testing.T) {
	t.Parallel()

	// given
	httpErr := &HTTPError{
		StatusCode: 403,
		Status:     "Forbidden",
	}
	authzErr := &AuthorizationError{
		HTTPError: httpErr,
	}

	// when
	result := authzErr.Error()

	// then
	expected := "HTTP Forbidden"
	if result != expected {
		t.Errorf("expected error message %q, got %q", expected, result)
	}
}

func TestRateLimitErrorShouldRetry(t *testing.T) {
	t.Parallel()

	// given
	httpErr := &HTTPError{
		StatusCode: 429,
		Status:     "Too Many Requests",
	}
	rateLimitErr := &RateLimitError{
		HTTPError:  httpErr,
		RetryDelay: 3 * time.Second,
	}
	testCases := []struct {
		retryAttempt int
		shouldRetry  bool
	}{
		{0, true},
		{1, true},
		{2, true},
		{3, true},
		{4, true},
		{5, false}, // maxRetries is 5
		{6, false},
	}

	for _, tc := range testCases {
		// when
		result := rateLimitErr.ShouldRetry(tc.retryAttempt)

		// then
		if result != tc.shouldRetry {
			t.Errorf("ShouldRetry(%d) = %v, expected %v", tc.retryAttempt, result, tc.shouldRetry)
		}
	}
}

func TestRateLimitErrorRetryAfter(t *testing.T) {
	t.Parallel()

	// given
	httpErr := &HTTPError{
		StatusCode: 429,
		Status:     "Too Many Requests",
	}
	rateLimitErr := &RateLimitError{
		HTTPError:  httpErr,
		RetryDelay: 3 * time.Second,
	}
	expected := 3 * time.Second
	testCases := []int{0, 1, 2, 3, 4, 5}

	for _, retryAttempt := range testCases {
		// when
		result := rateLimitErr.RetryAfter(retryAttempt)

		// then
		if result != expected {
			t.Errorf("RetryAfter(%d) = %v, expected %v", retryAttempt, result, expected)
		}
	}
}

func TestTimeoutErrorShouldRetry(t *testing.T) {
	t.Parallel()

	// given
	httpErr := &HTTPError{
		StatusCode: 408,
		Status:     "Request Timeout",
	}
	timeoutErr := &TimeoutError{
		HTTPError: httpErr,
	}
	testCases := []struct {
		retryAttempt int
		shouldRetry  bool
	}{
		{0, true},
		{1, true},
		{2, true},
		{3, true},
		{4, true},
		{5, false}, // maxRetries is 5
		{6, false},
	}

	for _, tc := range testCases {
		// when
		result := timeoutErr.ShouldRetry(tc.retryAttempt)

		// then
		if result != tc.shouldRetry {
			t.Errorf("ShouldRetry(%d) = %v, expected %v", tc.retryAttempt, result, tc.shouldRetry)
		}
	}
}

func TestTimeoutErrorRetryAfter(t *testing.T) {
	t.Parallel()

	// given
	httpErr := &HTTPError{
		StatusCode: 408,
		Status:     "Request Timeout",
	}
	timeoutErr := &TimeoutError{
		HTTPError: httpErr,
	}

	// Timeout errors should use exponential backoff
	testCases := []struct {
		retryAttempt int
		expected     time.Duration
	}{
		{0, 1 * time.Second},  // 2^0 = 1
		{1, 2 * time.Second},  // 2^1 = 2
		{2, 4 * time.Second},  // 2^2 = 4
		{3, 8 * time.Second},  // 2^3 = 8
		{4, 16 * time.Second}, // 2^4 = 16
	}

	for _, tc := range testCases {
		// when
		result := timeoutErr.RetryAfter(tc.retryAttempt)

		// then
		if result != tc.expected {
			t.Errorf("RetryAfter(%d) = %v, expected %v", tc.retryAttempt, result, tc.expected)
		}
	}
}

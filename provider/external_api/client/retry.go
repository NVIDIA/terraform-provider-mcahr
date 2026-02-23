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

import "time"

// Retryable interface for errors that support retry logic
type RetryableError interface {
	ShouldRetry(retryAttempt int) bool
	RetryAfter(retryAttempt int) time.Duration
}

func handleRetries(retryAttempt int, errIn error) (retry bool, errOut error) {
	retryable, ok := errIn.(RetryableError)
	if !ok {
		return false, errIn
	}

	if retryable.ShouldRetry(retryAttempt) {
		time.Sleep(retryable.RetryAfter(retryAttempt))
		return true, nil
	}

	return false, errIn
}

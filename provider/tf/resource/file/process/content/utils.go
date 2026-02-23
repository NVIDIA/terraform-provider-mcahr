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

package content

import (
	"crypto/md5"
	"fmt"
	"io"
	"os"
)

func ContentMd5(content []byte) (string, error) {
	// streaming md5sum
	hash := md5.New()
	_, err := hash.Write(content)
	if err != nil {
		return "", err
	}
	md5Sum := fmt.Sprintf("%x", hash.Sum(nil))
	return md5Sum, nil
}

func FileMd5(filename string) (string, error) {
	file, err := os.Open(filename)
	if err != nil {
		return "", err
	}
	defer file.Close()
	// streaming md5sum
	hash := md5.New()
	_, err = io.Copy(hash, file)
	if err != nil {
		return "", err
	}
	md5Sum := fmt.Sprintf("%x", hash.Sum(nil))
	return md5Sum, nil
}

func ContentSize(content []byte) int64 {
	return int64(len(content))
}

func FileSize(filename string) (int64, error) {
	fstat, err := os.Stat(filename)
	if err != nil {
		return 0, err
	}
	return fstat.Size(), nil
}

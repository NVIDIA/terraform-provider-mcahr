#!/bin/sh
# SPDX-FileCopyrightText: Copyright (c) 2025 NVIDIA CORPORATION & AFFILIATES. All rights reserved.
# SPDX-License-Identifier: Apache-2.0
#
# Licensed under the Apache License, Version 2.0 (the "License");
# you may not use this file except in compliance with the License.
# You may obtain a copy of the License at
#
# http://www.apache.org/licenses/LICENSE-2.0
#
# Unless required by applicable law or agreed to in writing, software
# distributed under the License is distributed on an "AS IS" BASIS,
# WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
# See the License for the specific language governing permissions and
# limitations under the License.

set -e

# Uploads a file to NGC if the tag does not exist.
# Usage: ./uploadVersion.sh <resource:tag> --source <file> [other-args...]

log() {
    echo "[$(date +'%Y-%m-%d %H:%M:%S')] $*" >&2
}

setup_ngc_cli() {
    if command -v ngc > /dev/null 2>&1; then
        log "NGC CLI already available"
        return 0
    fi
    
    log "Setting up NGC CLI..."
    
    log "Installing dependencies..."
    export DEBIAN_FRONTEND=noninteractive
    sudo apt-get update -qq && sudo apt-get install -y -qq curl wget unzip
    
    log "Downloading NGC CLI..."
    wget -q --content-disposition \
        'https://api.ngc.nvidia.com/v2/resources/nvidia/ngc-apps/ngc_cli/versions/3.63.0/files/ngccli_linux.zip' \
        -O ngccli_linux.zip
    unzip -o -q ngccli_linux.zip
    chmod +x ngc-cli/ngc
    export PATH="${PATH}:$(pwd)/ngc-cli"
    
    log "NGC CLI setup complete"
}

configure_ngc_cli() {
    if [ -z "$NGC_API_KEY_CI" ]; then
        log "Error: NGC_API_KEY_CI environment variable is not set"
        exit 1
    fi
    
    export NGC_CLI_API_KEY="${NGC_API_KEY_CI}"
    
    if [ -f ~/.ngc/config ]; then
        log "NGC CLI already configured"
        return 0
    fi
    
    log "Configuring NGC CLI..."
    mkdir -p ~/.ngc
    cat > ~/.ngc/config <<EOF
[CURRENT]
apikey = ${NGC_API_KEY_CI}
format_type = ascii
org = nvidian
EOF
    
    log "NGC CLI configured"
}

if [ $# -eq 0 ]; then
    log "Error: No arguments provided"
    log "Usage: $0 <resource:tag> --source <file> [other-args...]"
    exit 1
fi

setup_ngc_cli
configure_ngc_cli

resource_ref="$1"

if echo "$resource_ref" | grep -q '^nvidian/'; then
    log "Checking if resource version already exists..."
    
    if ngc registry resource info "$resource_ref" > /dev/null 2>&1; then
        log "Resource version already exists: $resource_ref"
        log "Skipping upload"
        exit 0
    else
        log "Resource version does not exist: $resource_ref"
        log "Proceeding with upload"
    fi
fi

log "Executing: ngc registry resource upload-version $*"
exec ngc registry resource upload-version "$@"

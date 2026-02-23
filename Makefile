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

# Import environment variables from .env file if it exists
ifneq (,$(wildcard ./.env))
    include .env
    export
endif

default: install

REPODIR=/tmp/tf-repo/providers

BINARY=terraform-provider-mcahr
VERSION=1.15.51

# IMPORTANT: When changing these build env vars, make sure to update the github release workflows files as well (the goreleaser env vars)
BUILD_ENV_VARS=-ldflags " -X 'main.ProviderPath=registry.opentofu.org/NVIDIA/mcahr' -X 'terraform/terraform-provider/provider.RenderedProviderName=\"NVIDIA Mission Control autonomous hardware recovery\"' -X 'terraform/terraform-provider/provider.ProviderShortName=mcahr' -X 'terraform/terraform-provider/provider.EnvVarsNamePrefix=MCAHR'"


// NOTE: this only works for 64 bit linux and MacOs ("darwin")
OS=$(shell uname | tr 'A-Z' 'a-z')
SUBPATH=registry.opentofu.org/NVIDIA/mcahr/local/mcahr/$(VERSION)/$(OS)_amd64

SCHEMA_FILE=provider_schema.json

generate_schema:
	@YELLOW='\033[1;33m'; GREEN='\033[0;32m'; RED='\033[0;31m'; NC='\033[0m'; \
	echo "$$YELLOW==> Generating provider schema...$$NC"; \
	\
	if [ -f ~/.terraformrc ] && grep -q "registry.opentofu.org/NVIDIA/mcahr" ~/.terraformrc 2>/dev/null; then \
		echo "$$GREEN✓ Local override already configured$$NC"; \
		WAS_LOCAL=true; \
	else \
		echo "$$YELLOW==> Setting up local override for schema generation$$NC"; \
		$(MAKE) use_local; \
		WAS_LOCAL=false; \
	fi; \
	\
	TEMP_DIR=$$(mktemp -d); \
	cp schema_init.tf "$$TEMP_DIR/main.tf"; \
	\
	if (cd "$$TEMP_DIR" && tofu providers schema -json > temp_schema.json 2>&1); then \
		PROVIDER_PATH_LOWER=$$(echo "registry.opentofu.org/NVIDIA/mcahr" | tr '[:upper:]' '[:lower:]'); \
		sed "s|\"$$PROVIDER_PATH_LOWER\"|\"mcahr\"|g" "$$TEMP_DIR/temp_schema.json" > "$(PWD)/$(SCHEMA_FILE)"; \
		rm -rf "$$TEMP_DIR"; \
		echo "$$GREEN✓ Schema generated: $(SCHEMA_FILE)$$NC"; \
	else \
		echo "$$RED✗ Schema generation failed!$$NC"; \
		echo "$$RED--- Error output ---$$NC"; \
		cat "$$TEMP_DIR/temp_schema.json"; \
		rm -rf "$$TEMP_DIR"; \
		if [ "$$WAS_LOCAL" = "false" ]; then $(MAKE) use_registry; fi; \
		exit 1; \
	fi; \
	\
	if [ "$$WAS_LOCAL" = "false" ]; then \
		echo "$$YELLOW==> Restoring registry configuration$$NC"; \
		$(MAKE) use_registry; \
	fi


generate_docs: generate_schema
	@YELLOW='\033[1;33m'; GREEN='\033[0;32m'; RED='\033[0;31m'; NC='\033[0m'; \
	echo "$$YELLOW\n==> Generating docs (using schema from $(SCHEMA_FILE)) ...\n$$NC"; \
	go run github.com/hashicorp/terraform-plugin-docs/cmd/tfplugindocs --provider-name=terraform-provider-mcahr --rendered-provider-name="NVIDIA Mission Control autonomous hardware recovery" --providers-schema=$(PWD)/$(SCHEMA_FILE); \
	GEN_STATUS=$$?; \
	if [ $$GEN_STATUS -eq 0 ]; then \
		echo "$$GREEN\n==> Docs generated successfully!\n$$NC"; \
		else \
		echo "$$RED\n==> Docs generation failed!\n$$NC"; \
		exit $$GEN_STATUS; \
	fi


go_generate:
	@YELLOW='\033[1;33m'; GREEN='\033[0;32m'; RED='\033[0;31m'; NC='\033[0m'; \
	echo "$$YELLOW\n==> Running go generate (format examples) ...\n$$NC"; \
	go generate $(BUILD_ENV_VARS); \
	GEN_STATUS=$$?; \
	if [ $$GEN_STATUS -eq 0 ]; then \
		echo "$$GREEN\n==> go generate succeeded!\n$$NC"; \
	else \
		echo "$$RED\n==> go generate failed!\n$$NC"; \
		exit $$GEN_STATUS; \
	fi

build: format go_generate
	@YELLOW='\033[1;33m'; GREEN='\033[0;32m'; RED='\033[0;31m'; NC='\033[0m'; \
	echo "$$YELLOW\n==> Building $(BINARY) ...\n$$NC"; \
	go build $(BUILD_ENV_VARS) -o ./$(BINARY); \
	BUILD_STATUS=$$?; \
	if [ $$BUILD_STATUS -eq 0 ]; then \
		echo "$$GREEN\n==> Build succeeded!\n$$NC"; \
	else \
		echo "$$RED\n==> Build failed!\n$$NC"; \
		exit $$BUILD_STATUS; \
	fi
	
test:
	echo unit-tests...

check:
	gofmt -l .

format:
	gofmt -w .

# NOTE: This relies on your ~/.terraformrc pointing to /tmp/tf-repo.
#   See terraformrc in the current dir
install: build
	@YELLOW='\033[1;33m'; GREEN='\033[0;32m'; RED='\033[0;31m'; NC='\033[0m'; \
	echo "$$YELLOW\n==> Installing provider binary to $(REPODIR)/$(SUBPATH) ...\n$$NC"; \
	rm -rf $(REPODIR)/*; \
	mkdir -p $(REPODIR)/$(SUBPATH); \
	cp $(BINARY) $(REPODIR)/$(SUBPATH)/$(BINARY); \
	if [ -f "$(REPODIR)/$(SUBPATH)/$(BINARY)" ]; then \
		echo "$$GREEN\n==> Provider binary installed successfully!\n$$NC"; \
	else \
		echo "$$RED\n==> Failed to install provider binary!\n$$NC"; \
		exit 1; \
	fi; \
	$(MAKE) generate_docs

# This sets up your ~/.terraformrc (NOTE: need to re-run when the version changes)
use_local: 
	@echo 'Setting up local overrides for terraform provider in ~/.terraformrc'
	@echo 'NOTE: You need to re-run "make use_local" when the version changes."'
	@echo 'provider_installation { dev_overrides { "registry.opentofu.org/NVIDIA/mcahr" = "$(REPODIR)/$(SUBPATH)" } }' > ${HOME}/.terraformrc

use_registry: 
	@echo 'Removing ~/.terraformrc, to use the terraform registry again'
	@rm ${HOME}/.terraformrc

release: 
	GOOS=darwin  GOARCH=amd64 go build $(BUILD_ENV_VARS) -o ./bin/$(BINARY)_$(VERSION)_darwin_amd64
	GOOS=darwin  GOARCH=arm64 go build $(BUILD_ENV_VARS) -o ./bin/$(BINARY)_$(VERSION)_darwin_arm64
	GOOS=linux   GOARCH=amd64 go build $(BUILD_ENV_VARS) -o ./bin/$(BINARY)_$(VERSION)_linux_amd64
	GOOS=linux   GOARCH=arm64 go build $(BUILD_ENV_VARS) -o ./bin/$(BINARY)_$(VERSION)_linux_arm64
	GOOS=linux   GOARCH=arm   go build $(BUILD_ENV_VARS) -o ./bin/$(BINARY)_$(VERSION)_linux_arm
	GOOS=openbsd GOARCH=amd64 go build $(BUILD_ENV_VARS) -o ./bin/$(BINARY)_$(VERSION)_openbsd_amd64
	GOOS=windows GOARCH=amd64 go build $(BUILD_ENV_VARS) -o ./bin/$(BINARY)_$(VERSION)_windows_amd64

version: 
	@echo "version: ${VERSION}\nTo create a release run: \n  git tag v${VERSION}\n  git push origin v${VERSION}"
	@git tag -l | grep '^v${VERSION}$$' >/dev/null && echo "WARNING: Release already exists" || true

# Run acceptance tests
.PHONY: testacc
testacc: 
	@ TF_ACC=1 go test ./... -v $(TESTARGS) -timeout 600s

# no checked in files should contain tokens
scan:
	find . -type f | xargs grep -l -e '[e]yJhb' || echo "scan is clean"


EXAMPLES_ROOT_PATH=./examples/resources/_root

init_ex:
	MCAHR_URL=$(MCAHR_URL) MCAHR_TOKEN=$(MCAHR_TOKEN) MCAHR_DEBUG=$(MCAHR_DEBUG) tofu -chdir=$(EXAMPLES_ROOT_PATH) init

apply_ex:
	MCAHR_URL=$(MCAHR_URL) MCAHR_TOKEN=$(MCAHR_TOKEN) MCAHR_DEBUG=$(MCAHR_DEBUG) tofu -chdir=$(EXAMPLES_ROOT_PATH) apply --auto-approve

apply_ex_na:
	MCAHR_URL=$(MCAHR_URL) MCAHR_TOKEN=$(MCAHR_TOKEN) MCAHR_DEBUG=$(MCAHR_DEBUG) tofu -chdir=$(EXAMPLES_ROOT_PATH) apply

destroy_ex:
	MCAHR_URL=$(MCAHR_URL) MCAHR_TOKEN=$(MCAHR_TOKEN) MCAHR_DEBUG=$(MCAHR_DEBUG) tofu -chdir=$(EXAMPLES_ROOT_PATH) destroy --auto-approve

destroy_ex_na:
	MCAHR_URL=$(MCAHR_URL) MCAHR_TOKEN=$(MCAHR_TOKEN) MCAHR_DEBUG=$(MCAHR_DEBUG) tofu -chdir=$(EXAMPLES_ROOT_PATH) destroy

plan_ex:
	MCAHR_URL=$(MCAHR_URL) MCAHR_TOKEN=$(MCAHR_TOKEN) MCAHR_DEBUG=$(MCAHR_DEBUG) tofu -chdir=$(EXAMPLES_ROOT_PATH) plan

.PHONY: distclean_ex
distclean_ex:
	rm -rf $(EXAMPLES_ROOT_PATH)/terraform.tfstate $(EXAMPLES_ROOT_PATH)/terraform.tfstate.backup $(EXAMPLES_ROOT_PATH)/.terraform $(EXAMPLES_ROOT_PATH)/.terraform.lock.hcl
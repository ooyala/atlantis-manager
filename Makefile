## Copyright 2014 Ooyala, Inc. All rights reserved.
##
## This file is licensed under the Apache License, Version 2.0 (the "License"); you may not use this file
## except in compliance with the License. You may obtain a copy of the License at
## http://www.apache.org/licenses/LICENSE-2.0
##
## Unless required by applicable law or agreed to in writing, software distributed under the License is
## distributed on an "AS IS" BASIS, WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
## See the License for the specific language governing permissions and limitations under the License.

PROJECT_ROOT := $(shell pwd)
PROJECT_NAME := $(shell pwd | xargs basename)
ifeq ($(shell pwd | xargs dirname | xargs basename),lib)
	LIB_PATH := $(shell pwd | xargs dirname)
	VENDOR_PATH := $(shell pwd | xargs dirname | xargs dirname)/vendor
else
	LIB_PATH := $(PROJECT_ROOT)/lib
	VENDOR_PATH := $(PROJECT_ROOT)/vendor
endif
ATLANTIS_PATH := $(LIB_PATH)/atlantis
SUPERVISOR_PATH := $(LIB_PATH)/atlantis-supervisor
ROUTER_PATH := $(LIB_PATH)/atlantis-router
BUILDER_PATH := $(LIB_PATH)/atlantis-builder

DEB_STAGING := $(PROJECT_ROOT)/staging
PKG_BIN_DIR := $(DEB_STAGING)/opt/atlantis-manager/bin

PKG := $(PROJECT_ROOT)/pkg
DEB := $(PROJECT_ROOT)/deb

DEB_INSTALL_DIR := $(PKG)/$(PROJECT_NAME)/opt/atlantis/manager
CLIENT_BIN_NAME := client
CLIENT_DEB_BIN_DIR := $(PKG)/atlantis-client/opt/atlantis/manager/bin

ifndef DEB_VERSION
	DEB_VERSION := "0.1.0"
endif

GOPATH := $(PROJECT_ROOT):$(VENDOR_PATH):$(ATLANTIS_PATH):$(SUPERVISOR_PATH):$(ROUTER_PATH):$(BUILDER_PATH)
export GOPATH

GOM := $(VENDOR_PATH)/bin/gom
GOM_VENDOR_NAME := vendor
export GOM_VENDOR_NAME

build: example

deb: clean build
	@cp -a $(DEB) $(PKG)
	@mkdir -p $(DEB_INSTALL_DIR)
	@mkdir -p $(DEB_INSTALL_DIR)/bin
	@cp -a example/manager $(DEB_INSTALL_DIR)/bin/atlantis-managerd
	@perl -p -i -e "s/__VERSION__/$(DEB_VERSION)/g" $(PKG)/$(PROJECT_NAME)/DEBIAN/control
	@mkdir -p $(CLIENT_DEB_BIN_DIR)
	@cp -a example/$(CLIENT_BIN_NAME) $(CLIENT_DEB_BIN_DIR)/atlantis-manager
	@perl -p -i -e "s/__VERSION__/$(DEB_VERSION)/g" $(PKG)/atlantis-client/DEBIAN/control
	@cd $(PKG) && dpkg --build $(PROJECT_NAME) ../pkg
	@cd $(PKG) && dpkg --build atlantis-client ../pkg

clean:
	@rm -rf $(ATLANTIS_PATH)/src/atlantis/crypto/key.go $(PROJECT_ROOT)/src/atlantis/manager/crypto/cert.go
	@rm -f example/client example/manager
	@rm -rf $(DEB_STAGING) atlantis-manager_*.deb
	@rm -rf ${PKG}
	@rm -rf $(VENDOR_PATH) $(LIB_PATH)

copy-key:
	@mkdir -p $(ATLANTIS_PATH)/src/atlantis/crypto
	@cp $(ATLANTIS_SECRET_DIR)/atlantis_key.go $(ATLANTIS_PATH)/src/atlantis/crypto/key.go
	@mkdir -p $(PROJECT_ROOT)/src/atlantis/manager/crypto
	@cp $(ATLANTIS_SECRET_DIR)/manager_cert.go $(PROJECT_ROOT)/src/atlantis/manager/crypto/cert.go

$(VENDOR_PATH):
	@echo "Installing Dependencies..."
	@mkdir -p $(VENDOR_PATH) || exit 2
	@GOPATH=$(VENDOR_PATH) go get github.com/ghao-ooyala/gom
	$(GOM) install
	@echo "Done."

test: clean copy-key
ifdef TEST_PACKAGE
	@echo "Testing $$TEST_PACKAGE..."
	@go test $$TEST_PACKAGE $$VERBOSE $$RACE
else
	@for p in `find ./src -type f -name "*.go" |sed 's-\./src/\(.*\)/.*-\1-' |sort -u`; do \
		echo "Testing $$p..."; \
		go test $$p || exit 1; \
	done
	@echo
	@echo "ok."
endif

.PHONY: example
example: $(VENDOR_PATH) copy-key
	@go build -o example/manager example/manager.go
	@go build -o example/client example/client.go

fmt:
	@find src -name \*.go -exec gofmt -l -w {} \;

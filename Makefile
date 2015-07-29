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

ifndef VERSION
	VERSION := "0.1.0"
endif

GOPATH := $(PROJECT_ROOT):$(VENDOR_PATH):$(ATLANTIS_PATH):$(SUPERVISOR_PATH):$(ROUTER_PATH):$(BUILDER_PATH)
export GOPATH

build: install-deps example

deb: clean build
	@cp -a $(PROJECT_ROOT)/deb $(DEB_STAGING)
	@mkdir -p $(PKG_BIN_DIR)

	@cp example/manager $(PKG_BIN_DIR)
	@cp example/client $(PKG_BIN_DIR)

	@sed -ri "s/__VERSION__/$(VERSION)/" $(DEB_STAGING)/DEBIAN/control
	@sed -ri "s/__PACKAGE__/atlantis-manager/" $(DEB_STAGING)/DEBIAN/control
	@dpkg -b $(DEB_STAGING) .

clean:
	@rm -rf $(ATLANTIS_PATH)/src/atlantis/crypto/key.go $(PROJECT_ROOT)/src/atlantis/manager/crypto/cert.go
	@rm -f example/client example/manager
	@rm -rf $(DEB_STAGING) atlantis-manager_*.deb

copy-key: clean
	@mkdir -p $(ATLANTIS_PATH)/src/atlantis/crypto
	@cp $(ATLANTIS_SECRET_DIR)/atlantis_key.go $(ATLANTIS_PATH)/src/atlantis/crypto/key.go
	@mkdir -p $(PROJECT_ROOT)/src/atlantis/manager/crypto
	@cp $(ATLANTIS_SECRET_DIR)/manager_cert.go $(PROJECT_ROOT)/src/atlantis/manager/crypto/cert.go

install-deps:
	@echo "Installing Dependencies..."
	@sudo apt-get install -y bzr
	@rm -rf $(VENDOR_PATH)
	@mkdir -p $(VENDOR_PATH) || exit 2
	@GOPATH=$(VENDOR_PATH) go get github.com/jigish/go-flags
	@GOPATH=$(VENDOR_PATH) go get github.com/jigish/gozk-recipes
	@GOPATH=$(VENDOR_PATH) go get github.com/BurntSushi/toml
	@GOPATH=$(VENDOR_PATH) go get github.com/cespare/go-apachelog
	@GOPATH=$(VENDOR_PATH) go get github.com/gorilla/context
	@GOPATH=$(VENDOR_PATH) go get github.com/gorilla/mux
	@GOPATH=$(VENDOR_PATH) go get github.com/mavricknz/asn1-ber
	@GOPATH=$(VENDOR_PATH) go get github.com/mavricknz/ldap
	@GOPATH=$(VENDOR_PATH) go get github.com/mewpkg/gopass
	@GOPATH=$(VENDOR_PATH) go get github.com/ooyala/go-jenkins-cli
	@git clone https://github.com/jigish/route53.git $(VENDOR_PATH)/src/github.com/jigish/route53
	@mkdir -p $(VENDOR_PATH)/src/github.com/crowdmob && git clone https://github.com/crowdmob/goamz.git $(VENDOR_PATH)/src/github.com/crowdmob/goamz
	@GOPATH=$(VENDOR_PATH) go get code.google.com/p/gographviz
	@GOPATH=$(VENDOR_PATH) go get launchpad.net/gocheck
	@echo "Done."

test: clean copy-key
ifdef TEST_PACKAGE
	@echo "Testing $$TEST_PACKAGE..."
	@go test $$TEST_PACKAGE $$VERBOSE $$RACE
else
ifneq ($(path),)
	@echo "Testing $(path)..."
	@go test $(path) || exit 1;
	@echo
	@echo "ok."
else
	@for p in `find ./src -type f -name "*.go" |sed 's-\./src/\(.*\)/.*-\1-' |sort -u`; do \
		echo "Testing $$p..."; \
		go test $$p || exit 1; \
	done
	@echo
	@echo "ok."
endif
endif

.PHONY: example
example: copy-key
	@go build -o example/manager example/manager.go
	@go build -o example/client example/client.go

fmt:
	@find src -name \*.go -exec gofmt -l -w {} \;

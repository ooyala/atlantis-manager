PROJECT_ROOT := $(shell pwd)
VENDOR_PATH  := $(PROJECT_ROOT)/vendor
LIB_PATH := $(PROJECT_ROOT)/lib
ATLANTIS_PATH := $(LIB_PATH)/atlantis
SUPERVISOR_PATH := $(LIB_PATH)/atlantis-supervisor
ROUTER_PATH := $(LIB_PATH)/atlantis-router

GOPATH := $(PROJECT_ROOT):$(VENDOR_PATH):$(ATLANTIS_PATH):$(SUPERVISOR_PATH):$(ROUTER_PATH)
export GOPATH

all: test

clean:
	@rm -rf $(ATLANTIS_PATH)/src/atlantis/crypto/key.go $(PROJECT_ROOT)/src/atlantis/manager/crypto/cert.go
	@rm example/client example/manager

copy-key: clean
	@cp $(ATLANTIS_SECRET_DIR)/atlantis_key.go $(ATLANTIS_PATH)/src/atlantis/crypto/key.go
	@cp $(ATLANTIS_SECRET_DIR)/manager_cert.go $(PROJECT_ROOT)/src/atlantis/manager/crypto/cert.go

install:
	@echo "Installing Dependencies..."
	@rm -rf $(LIB_PATH) $(VENDOR_PATH)
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
	@GOPATH=$(VENDOR_PATH) go get code.google.com/p/gographviz
	@GOPATH=$(VENDOR_PATH) go get launchpad.net/gocheck
	@git clone ssh://git@github.com/ooyala/atlantis $(ATLANTIS_PATH)
	@git clone ssh://git@github.com/ooyala/atlantis-supervisor $(SUPERVISOR_PATH)
	@git clone ssh://git@github.com/ooyala/atlantis-router $(ROUTER_PATH)
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
example: copy-key
	@go build -o example/manager example/manager.go
	@go build -o example/client example/client.go

fmt:
	@find src -name \*.go -exec gofmt -l -w {} \;

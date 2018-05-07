# Package configuration
PROJECT := gitbase-playground
COMMANDS := cmd/gitbase-playground
DOCKER_ORG := src-d
DEPENDENCIES := \
	github.com/golang/dep/cmd/dep \
	github.com/jteeuwen/go-bindata \
	github.com/golang/lint/golint
GO_LINTABLE_PACKAGES := $(shell go list ./... | grep -v '/vendor/')
GO_BUILD_ENV := CGO_ENABLED=0

# Tools
GODEP := dep
GOLINT := golint
GOVET := go vet
BINDATA := go-bindata
DIFF := diff

# Default rule
all:

# Including ci Makefile
CI_REPOSITORY ?= https://github.com/src-d/ci.git
CI_BRANCH ?= v1
CI_PATH ?= $(shell pwd)/.ci
MAKEFILE := $(CI_PATH)/Makefile.main
$(MAKEFILE):
	@git clone --quiet --depth 1 -b $(CI_BRANCH) $(CI_REPOSITORY) $(CI_PATH);
-include $(MAKEFILE)

# Makefile.main::dependencies -> Makefile.main::$(DEPENDENCIES) -> this::dependencies
# The `exit` is needed to prevent running `Makefile.main::dependencies` commands.
dependencies: | front-dependencies back-dependencies exit

# Makefile.main::test -> this::test
test: front-test

# this::build -> Makefile.main::build -> Makefile.main::$(COMMANDS)
# The @echo forces this prerequisites to be run before `Makefile.main::build` ones.
build: front-build back-build
	@echo

coverage: | test-coverage codecov

lint: back-lint front-lint

validate-commit: | \
	back-dependencies \
	back-ensure-assets-proxy \
	front-fix-lint-errors \
	no-changes-in-commit

exit:
	exit 0;


## Compiles the assets, and serve the tool through its API

serve: | front-build back-start


# Backend

assets := ./server/assets/asset.go
assets_back := $(assets).bak

back-dependencies:
	$(GODEP) ensure

back-build: back-bindata

back-bindata:
	$(BINDATA) \
		-pkg assets \
		-o $(assets) \
		build/public/*

back-lint: $(GO_LINTABLE_PACKAGES)
$(GO_LINTABLE_PACKAGES):
	$(GOLINT) $@
	$(GOVET) $@

back-start:
	GITBASEPG_ENV=dev go run cmd/gitbase-playground/main.go

back-ensure-assets-proxy:
	$(DIFF) $(assets) $(assets_back) || exit 1



# Frontend

front-dependencies:
	echo 'SKIP. no frontend dependencies to install'

front-test:
	echo 'SKIP. no frontend tests to run'

front-lint:
	echo 'SKIP. no frontend linters to run'

front-build:
	mkdir -p build/public
	cp public/index.html build/public/index.html
	cp public/secondary.html build/public/secondary.html

front-fix-lint-errors:
	echo 'SKIP. no fixable code'



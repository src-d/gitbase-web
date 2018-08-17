# Package configuration
PROJECT := gitbase-web
COMMANDS := cmd/gitbase-web
DEPENDENCIES := \
	github.com/jteeuwen/go-bindata \
	github.com/golang/lint/golint
DEPENDENCIES_DIRECTORY := ./vendor

PKG_OS = linux

GO_LINTABLE_PACKAGES := $(shell go list ./... | grep -v '/vendor/')
FRONTEND_PATH := ./frontend
FRONTEND_BUILD_PATH := $(FRONTEND_PATH)/build

# Tools
GOLINT := golint
GOVET := go vet
BINDATA := go-bindata
YARN := yarn --cwd $(FRONTEND_PATH)
REMOVE := rm -rf
MOVE := mv
MKDIR := mkdir -p
COMPOSE := docker-compose

# To be used as -tags
GO_BINDATA_TAG := bindata

# Environment and arguments to use in `go run` calls.
GO_RUN_ENV := GITBASEPG_ENV=dev
GO_RUN_ARGS += -tags "$(GO_BINDATA_TAG)"

GORUN = $(GOCMD) run $(GO_RUN_ARGS)

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

# Override Makefile.main defaults for arguments to be used in `go` commands.
GO_BUILD_ARGS := -ldflags "$(LD_FLAGS)" -tags "$(GO_BINDATA_TAG)"

# TODO: remove when https://github.com/src-d/ci/pull/84 is merged
.PHONY: godep
GODEP ?= $(CI_PATH)/dep
godep:
	export INSTALL_DIRECTORY=$(CI_PATH) ; \
	test -f $(GODEP) || \
		curl https://raw.githubusercontent.com/golang/dep/master/install.sh | bash ; \
	$(GODEP) ensure -v

# Makefile.main::dependencies -> Makefile.main::$(DEPENDENCIES) -> this::dependencies
# The `exit` is needed to prevent running `Makefile.main::dependencies` commands.
dependencies: | front-dependencies exit

# Makefile.main::test -> this::test
test: front-test

# this::build -> Makefile.main::build -> Makefile.main::$(COMMANDS)
# The @echo forces this prerequisites to be run before `Makefile.main::build` ones.
build: front-build back-build
	@echo

coverage: | test-coverage codecov

lint: back-lint front-lint

validate-commit: | \
	godep \
	front-fix-lint-errors \
	no-changes-in-commit

exit:
	exit 0;

clean: front-clean

build-path:
	$(MKDIR) $(BUILD_PATH)

.PHONY: version
version:
	@echo $(VERSION)

# Compiles the assets, and serves the tool through its API

serve: | front-dependencies front-build back-build back-start

compose-serve-latest:
	$(COMPOSE) pull && \
	GITBASEPG_REPOS_FOLDER=./repos $(COMPOSE) up --force-recreate

compose-serve: | require-repos-folder front-dependencies build
	GITBASEPG_REPOS_FOLDER=${GITBASEPG_REPOS_FOLDER} \
		$(COMPOSE) -f docker-compose.yml -f docker-compose.build.yml up --force-recreate --build

require-repos-folder:
	@if [[ -z "$(GITBASEPG_REPOS_FOLDER)" ]]; then \
		echo "error. undefined 'GITBASEPG_REPOS_FOLDER' to be served under gitbase"; \
		exit 1; \
	fi
	$(MKDIR) $(GITBASEPG_REPOS_FOLDER)

back-build: back-bindata

back-bindata:
	$(BINDATA) \
		-pkg assets \
		-o ./server/assets/asset.go \
		build/public/...

back-lint: $(GO_LINTABLE_PACKAGES)
$(GO_LINTABLE_PACKAGES):
	$(GOLINT) $@
	$(GOVET) $@

back-start:
	$(GO_RUN_ENV) $(GORUN) cmd/gitbase-web/main.go

# Frontend

front-dependencies:
	$(YARN) install

front-test:
	$(YARN) test

front-lint:
	$(YARN) lint

front-build: build-path
	$(YARN) build
	$(REMOVE) $(BUILD_PATH)/public
	$(MOVE) $(FRONTEND_BUILD_PATH) $(BUILD_PATH)/public

front-fix-lint-errors:
	$(YARN) format

front-clean:
	$(REMOVE) $(FRONTEND_PATH)/node_modules
	$(REMOVE) $(FRONTEND_BUILD_PATH)

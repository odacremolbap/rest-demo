GITHUB_REPO ?= github.com/odacremolbap/rest-demo
ALL_BIN ?= todolist
VERSION ?= $(shell git describe --tags --always --dirty)
OUTPUT_DIR := _output
CGO_ENABLED := 0
ALL_ARCH := amd64 arm64
ARCH ?= amd64
ALL_OS := linux darwin
OS ?= $(shell go env GOOS)

DATE:=$(shell TZ=UTC date +'%y.%m.%d %H:%M:%S')
GO_LDFLAGS:=-X $(REPO)/pkg/version.Version=${VERSION} -X \"$(REPO)/pkg/version.Date=${DATE}\"

.PHONY: build
build: \
			dep \
			clean \
			prebuild-bin-$(ARCH)-$(OS)

.PHONY: all
all: \
			dep \
			clean \
			test \
			prebuild-arch

.PHONY: release
release: \
			dep-clean \
			all

prebuild-arch: $(foreach arch, $(ALL_ARCH), prebuild-os-$(arch))
	$(NOOP)
prebuild-os-%: $(foreach os, $(ALL_OS), prebuild-bin-%-$(os))
	$(NOOP)
prebuild-bin-%: $(foreach bin, $(ALL_BIN), prebuild-launch-%-$(bin))
		$(NOOP)

prebuild-launch-%:
	$(eval STR = $(subst -, ,$@))
	$(eval ARCH = $(word 3, $(STR)))
	$(eval OS = $(word 4, $(STR)))
	$(eval BINARY = $(word 5, $(STR)))

	@$(MAKE) --no-print-directory BINARY=$(BINARY) ARCH=$(ARCH) OS=$(OS) launch-build

.PHONY: build
launch-build:
	$(eval OUTPUT_BIN_DIR = $(OUTPUT_DIR)/$(OS)/$(ARCH))

ifneq ($(ARCH)-$(OS),arm64-darwin)
	@mkdir -p $(OUTPUT_BIN_DIR)
	$(info Building $(BINARY) for $(OS)/$(ARCH))
	docker run -ti --rm  \
					-v "$$(pwd):/go/src/$(GITHUB_REPO)" \
					-e "CGO_ENABLED=$(CGO_ENABLED)" \
					-e "GOOS=$(OS)" \
					-e "GOARCH=$(ARCH)" \
					golang:latest \
					go build $(GO_FLAGS) -ldflags "$(GO_LDFLAGS)" -o /go/src/$(GITHUB_REPO)/$(OUTPUT_BIN_DIR)/$(BINARY) $(GITHUB_REPO)/cmd/$(BINARY)
else
	@rm -rf $(OUTPUT_BIN_DIR)
endif

.PHONY: test
test:
	@mkdir -p $(OUTPUT_DIR)/test
	docker run -ti --rm  \
					-v "$$(pwd):/go/src/$(GITHUB_REPO)" \
					-e "CGO_ENABLED=$(CGO_ENABLED)" \
					golang:latest \
					go test $(GITHUB_REPO)/... -coverprofile /go/src/$(GITHUB_REPO)/$(OUTPUT_DIR)/test/cover.out

.PHONY: dep
dep:
	dep ensure

.PHONY: clean
clean:
	rm -rf _output

.PHONY: dep-clean
dep-clean:
	rm -rf ./vendor/

.PHONY: db
db:
	if [ "$$(docker ps -q -f name=postgres)" ]; then \
		docker stop postgres; \
		docker rm postgres; \
	fi; \
	docker run --name postgres -p 5432:5432 -d postgres;
	# TODO replace with tcp pings
	sleep 3;
	psql -h localhost -U postgres -p 5432 -f assets/deployment/database/schema.sql;

.PHONY: run
run:
	@./assets/run/run.sh

.PHONY: run-container
run-container:
	@./assets/run/run-container.sh

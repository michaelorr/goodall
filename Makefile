# Copyright 2016 The Kubernetes Authors.
#
# Licensed under the Apache License, Version 2.0 (the "License");
# you may not use this file except in compliance with the License.
# You may obtain a copy of the License at
#
#     http://www.apache.org/licenses/LICENSE-2.0
#
# Unless required by applicable law or agreed to in writing, software
# distributed under the License is distributed on an "AS IS" BASIS,
# WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
# See the License for the specific language governing permissions and
# limitations under the License.

# The binary to build (just the basename).
BIN := goodall

# This repo's root import path (under GOPATH).
PKG := github.com/michaelorr/goodall

# Which architecture to build - see $(ALL_ARCH) for options.
ARCH ?= amd64

# Which OS to target when building.
# Common choices are `darwin`, `linux` or `windows`.
# Complete list of options: https://golang.org/doc/install/source#environment
OS ?= darwin

# This version-strategy uses git tags to set the version string
# VERSION := $(shell git describe --always --dirty)
#
# This version-strategy uses a manual value to set the version string
VERSION := 0.1.0

###
### These variables should not need tweaking.
###

SRC_DIRS := cmd pkg # directories which hold app source (not vendored)

ALL_ARCH := amd64 arm arm64 ppc64le

BUILD_IMAGE ?= golang:1.7-alpine

# If you want to build all binaries, see the 'all-build' rule.
# If you want to build all containers, see the 'all-container' rule.
# If you want to build AND push all containers, see the 'all-push' rule.
all: build

build-%:
	@$(MAKE) --no-print-directory ARCH=$* build

all-build: $(addprefix build-, $(ALL_ARCH))

build: bin/$(OS)_$(ARCH)/$(BIN)

bin/$(OS)_$(ARCH)/$(BIN): build-dirs
	@echo "building: $@"
	@docker run                                                            \
		-ti                                                                \
		-u $$(id -u):$$(id -g)                                             \
		-v $$(pwd)/.go:/go                                                 \
		-v $$(pwd):/go/src/$(PKG)                                          \
		-v $$(pwd)/bin/$(OS)_$(ARCH):/go/bin                               \
		-v $$(pwd)/.go/std/$(ARCH):/usr/local/go/pkg/linux_$(ARCH)_static  \
		-w /go/src/$(PKG)                                                  \
		$(BUILD_IMAGE)                                                     \
		/bin/sh -c "                                                       \
			OS=$(OS)													   \
			ARCH=$(ARCH)                                                   \
			VERSION=$(VERSION)                                             \
			PKG=$(PKG)                                                     \
			./build/build.sh                                               \
		"

version:
	@echo $(VERSION)

test-verbose:
	@$(MAKE) --no-print-directory VERBOSE="-v" test

test: build-dirs
	@docker run                                                            \
		-ti                                                                \
		-u $$(id -u):$$(id -g)                                             \
		-v $$(pwd)/.go:/go                                                 \
		-v $$(pwd):/go/src/$(PKG)                                          \
		-v $$(pwd)/bin/$(OS)_$(ARCH):/go/bin                               \
		-v $$(pwd)/.go/std/$(ARCH):/usr/local/go/pkg/linux_$(ARCH)_static  \
		-w /go/src/$(PKG)                                                  \
		$(BUILD_IMAGE)                                                     \
		/bin/sh -c "                                                       \
			VERBOSE=$(VERBOSE)											   \
			./build/test.sh $(SRC_DIRS)                                    \
		"

build-dirs:
	@mkdir -p bin/$(OS)_$(ARCH)
	@mkdir -p .go/src/$(PKG) .go/pkg .go/bin .go/std/$(ARCH)

clean: bin-clean

bin-clean:
	rm -rf .go bin

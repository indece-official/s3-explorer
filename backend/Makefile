# S3 Explorer
# Copyright (C) 2022  indece UG (haftungsbeschränkt)
#
# This program is free software: you can redistribute it and/or modify
# it under the terms of the GNU General Public License as published by
# the Free Software Foundation, either version 3 of the License or any
# later version.
#
# This program is distributed in the hope that it will be useful,
# but WITHOUT ANY WARRANTY; without even the implied warranty of
# MERCHANTABILITY or FITNESS FOR A PARTICULAR PURPOSE. See the
# GNU General Public License for more details.
#
# You should have received a copy of the GNU General Public License
# along with this program. If not, see <https://www.gnu.org/licenses/>.

# Go parameters
PROJECT_NAME ?= $(shell basename $(dir $(abspath $(lastword $(MAKEFILE_LIST)))))
GOCMD=go
GOPATH=$(shell $(GOCMD) env GOPATH))
GOBUILD=$(GOCMD) build
GOGENERATE=$(GOCMD) generate
GOCLEAN=$(GOCMD) clean
GOTEST=$(GOCMD) test
GOGET=$(GOCMD) get
DIR_SOURCE=./src
DIR_GENERATED_RESOURCES=$(DIR_SOURCE)/generated/resources
DIR_GENERATED_MODEL=$(DIR_SOURCE)/generated/model
DIR_DIST=./dist
BINARY_NAME=$(DIR_DIST)/bin/s3-explorer
BUILD_DATE=$(shell date +%Y%m%d.%H%M%S)
BUILD_VERSION ?= $(shell git rev-parse --short HEAD)
LDFLAGS := 
LDFLAGS := $(LDFLAGS) -X main.ProjectName=$(PROJECT_NAME)
LDFLAGS := $(LDFLAGS) -X main.BuildDate=$(BUILD_DATE)
LDFLAGS := $(LDFLAGS) -X main.BuildVersion=$(BUILD_VERSION)

all: generate test prepare_compile compile_linux crosscompile_windows

generate:
	mkdir -p $(DIR_GENERATED_MODEL)
	rm -rf $(DIR_GENERATED_MODEL)/*
	$(GOGENERATE) -tags=bindata ./...

prepare_compile:
	mkdir -p $(DIR_DIST)/bin

compile_linux:
	$(GOBUILD) -ldflags "$(LDFLAGS)" -o $(BINARY_NAME) -tags=prod -v $(DIR_SOURCE)/main

compile_windows:
	GOOS=windows \
	GOARCH=amd64 \
	CGO_ENABLED=1 \
	$(GOBUILD) -ldflags "-H=windowsgui $(LDFLAGS)" -o $(BINARY_NAME)_win64.exe -tags=prod -v $(DIR_SOURCE)/main

crosscompile_windows:
	# Make resources
	/usr/bin/x86_64-w64-mingw32-windres $(DIR_SOURCE)/main/resources.rc $(DIR_SOURCE)/main/resources.syso
	# Compile & link go files & resources
	PKG_CONFIG_PATH=/usr/x86_64-w64-mingw32/lib/pkgconfig \
	CGO_ENABLED=1 \
	CC=x86_64-w64-mingw32-gcc \
	CXX=x86_64-w64-mingw32-g++ \
	GOOS=windows \
	GOARCH=amd64 \
	$(GOBUILD) -ldflags "-H=windowsgui" -o $(BINARY_NAME)_win64-unsigned.exe -v $(DIR_SOURCE)/main

sign:
	test -n "${CODESIGN_CERT}"
	@test -n "${CODESIGN_PASSWORD}"
	@osslsigncode sign \
	-pkcs12 "${CODESIGN_CERT}" \
	-pass "${CODESIGN_PASSWORD}" \
	-n "S3 Explorer" \
	-i https://www.indece.com/ \
	-in $(BINARY_NAME)_win64-unsigned.exe \
	-out $(BINARY_NAME)_win64.exe

test:
	mkdir -p $(DIR_DIST)
ifeq ($(OUTPUT),json)
	$(GOTEST) -v ./...  -cover -coverprofile $(DIR_DIST)/cover.out -json > $(DIR_DIST)/test.json
else
	$(GOTEST) -v ./...  -cover
endif

clean:
	#$(GOCLEAN)
	rm -rf $(DIR_OUT)

run:
	$(GOBUILD) -o $(BINARY_NAME) -v ./...
	./$(BINARY_NAME)

deps:
	echo test
	#$(GOGET) -d -v ./...
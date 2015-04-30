
SHELL = /bin/bash

BUILD_DIR = build/metrilyx-cacher
BIN_DEST = /opt/metrilyx/bin

clean:
	rm -rf ./bin
	rm -rf ./build

build:
	go get -d -v ./...
	go build -v metrilyx-cacher.go 

install:
	[ -d "$(BUILD_DIR)" ] || mkdir -p $(BUILD_DIR)/$(BIN_DEST)
	mv metrilyx-cacher $(BUILD_DIR)/$(BIN_DEST)/
	cp -a etc $(BUILD_DIR)/

all: clean build

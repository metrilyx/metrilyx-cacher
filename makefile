
SHELL = /bin/bash

BUILD_DIR = build/metrilyx-cacher

clean:
	rm -rf ./bin
	rm -rf ./build

build:
	[ -d "$(BUILD_DIR)" ] || mkdir -p $(BUILD_DIR)/opt/metrilyx/bin
	go get -d -v ./...
	go build -v metrilyx-cacher.go 

install:
	mv metrilyx-cacher $(BUILD_DIR)/opt/metrilyx/bin/
	cp -a etc $(BUILD_DIR)/

all: clean build
